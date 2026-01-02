package main

import (
	"errors"
	"mime/multipart"
	"net/http"

	"api.etin.dev/internal/assets"
	"api.etin.dev/internal/data"
)

const maxAssetUploadBytes = 10 << 20 // 10 MiB

type assetSaver interface {
	Insert(*data.Asset) error
}

func (app *application) getCreateAssetsHandler(w http.ResponseWriter, r *http.Request) {
	if !app.isRequestAuthenticated(r) {
		app.writeError(w, http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAssetUploadBytes)
	if err := r.ParseMultipartForm(maxAssetUploadBytes); err != nil {
		if app.logger != nil {
			app.logger.Printf("parse multipart form: %v", err)
		}
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			app.writeError(w, http.StatusRequestEntityTooLarge)
			return
		}
		if errors.Is(err, http.ErrNotMultipart) {
			app.writeError(w, http.StatusBadRequest)
			return
		}

		if errors.Is(err, multipart.ErrMessageTooLarge) {
			app.writeError(w, http.StatusRequestEntityTooLarge)
			return
		}

		app.writeError(w, http.StatusBadRequest)
		return
	}
	defer func() {
		if r.MultipartForm != nil {
			r.MultipartForm.RemoveAll()
		}
	}()

	file, _, err := r.FormFile("file")
	if err != nil {
		if app.logger != nil {
			app.logger.Printf("read file part: %v", err)
		}
		if errors.Is(err, http.ErrMissingFile) {
			app.writeError(w, http.StatusBadRequest)
			return
		}

		app.writeError(w, http.StatusBadRequest)
		return
	}
	defer file.Close()

	result, err := app.assets.Upload(r.Context(), file, assets.UploadOptions{})
	if err != nil {
		if app.logger != nil {
			app.logger.Printf("upload asset: %v", err)
		}
		app.writeError(w, http.StatusBadGateway)
		return
	}

	asset := &data.Asset{
		URL:          result.URL,
		SecureURL:    result.SecureURL,
		PublicID:     result.PublicID,
		Format:       result.Format,
		ResourceType: result.ResourceType,
		Bytes:        result.Bytes,
		Width:        result.Width,
		Height:       result.Height,
	}

	if err := app.assetModel.Insert(asset); err != nil {
		if app.logger != nil {
			app.logger.Printf("persist asset: %v", err)
		}
		app.writeError(w, http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"asset": asset})
	return
}
