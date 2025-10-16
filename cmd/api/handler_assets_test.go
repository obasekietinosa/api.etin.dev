package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"api.etin.dev/internal/assets"
	"api.etin.dev/internal/data"
)

type stubUploader struct {
	result *assets.UploadResult
	err    error
}

func (s *stubUploader) Upload(_ context.Context, file io.Reader, _ assets.UploadOptions) (*assets.UploadResult, error) {
	if s.err != nil {
		return nil, s.err
	}

	// Drain the reader to mimic Cloudinary consumption and ensure callers can stream.
	io.Copy(io.Discard, file)

	if s.result == nil {
		return &assets.UploadResult{}, nil
	}

	clone := *s.result
	return &clone, nil
}

type stubAssetSaver struct {
	saved  []*data.Asset
	err    error
	nextID int64
}

func (s *stubAssetSaver) Insert(asset *data.Asset) error {
	if s.err != nil {
		return s.err
	}

	s.nextID++
	asset.ID = s.nextID

	copy := *asset
	s.saved = append(s.saved, &copy)
	return nil
}

func newAuthenticatedApp(t *testing.T, uploader assets.Uploader, saver assetSaver) (*application, string) {
	t.Helper()

	sm := newSessionManager(time.Hour)
	token, expiry, err := sm.create()
	if err != nil {
		t.Fatalf("failed to seed session token: %v", err)
	}

	sm.sessions[token] = expiry.Add(time.Hour)

	app := &application{
		logger:     log.New(io.Discard, "", 0),
		assets:     uploader,
		assetModel: saver,
		sessions:   sm,
	}

	return app, token
}

func createMultipartRequest(t *testing.T, token string, payload []byte) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if payload != nil {
		part, err := writer.CreateFormFile("file", "test.bin")
		if err != nil {
			t.Fatalf("create form file: %v", err)
		}
		if _, err := part.Write(payload); err != nil {
			t.Fatalf("write payload: %v", err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/assets", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	return req, rr
}

func TestGetCreateAssetsHandler_Success(t *testing.T) {
	uploader := &stubUploader{result: &assets.UploadResult{
		AssetID:      "asset123",
		PublicID:     "public123",
		Format:       "png",
		ResourceType: "image",
		URL:          "http://example.com/resource.png",
		SecureURL:    "https://example.com/resource.png",
		Version:      1,
		Bytes:        512,
		Width:        10,
		Height:       10,
	}}

	saver := &stubAssetSaver{}

	app, token := newAuthenticatedApp(t, uploader, saver)

	req, rr := createMultipartRequest(t, token, []byte("hello world"))

	app.getCreateAssetsHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d; got %d", http.StatusCreated, rr.Code)
	}

	if len(saver.saved) != 1 {
		t.Fatalf("expected asset to be saved")
	}

	var payload struct {
		Asset data.Asset `json:"asset"`
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.Asset.ID == 0 {
		t.Fatalf("expected asset ID to be populated")
	}

	if payload.Asset.SecureURL != uploader.result.SecureURL {
		t.Fatalf("expected secure URL to match upload result")
	}
}

func TestGetCreateAssetsHandler_MissingFile(t *testing.T) {
	uploader := &stubUploader{result: &assets.UploadResult{}}
	saver := &stubAssetSaver{}
	app, token := newAuthenticatedApp(t, uploader, saver)

	req, rr := createMultipartRequest(t, token, nil)

	app.getCreateAssetsHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d; got %d", http.StatusBadRequest, rr.Code)
	}

	if len(saver.saved) != 0 {
		t.Fatalf("expected no assets to be saved")
	}
}

func TestGetCreateAssetsHandler_FileTooLarge(t *testing.T) {
	uploader := &stubUploader{result: &assets.UploadResult{}}
	saver := &stubAssetSaver{}
	app, token := newAuthenticatedApp(t, uploader, saver)

	largePayload := bytes.Repeat([]byte("a"), int(maxAssetUploadBytes)+1)

	req, rr := createMultipartRequest(t, token, largePayload)

	app.getCreateAssetsHandler(rr, req)

	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d; got %d", http.StatusRequestEntityTooLarge, rr.Code)
	}

	if len(saver.saved) != 0 {
		t.Fatalf("expected no assets to be saved")
	}
}

func TestGetCreateAssetsHandler_UploadFailure(t *testing.T) {
	uploader := &stubUploader{err: errors.New("upload failed")}
	saver := &stubAssetSaver{}
	app, token := newAuthenticatedApp(t, uploader, saver)

	req, rr := createMultipartRequest(t, token, []byte("hello world"))

	app.getCreateAssetsHandler(rr, req)

	if rr.Code != http.StatusBadGateway {
		t.Fatalf("expected status %d; got %d", http.StatusBadGateway, rr.Code)
	}

	if len(saver.saved) != 0 {
		t.Fatalf("expected no assets to be saved")
	}
}
