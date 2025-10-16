package assets

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Uploader defines the behaviour required to upload binary data to an asset store.
type Uploader interface {
	Upload(ctx context.Context, file io.Reader, options UploadOptions) (*UploadResult, error)
}

// UploadOptions configures how an asset should be stored.
type UploadOptions struct {
	Folder       string
	PublicID     string
	ResourceType string
	Overwrite    *bool
}

// UploadResult captures the relevant metadata returned by Cloudinary after an upload.
type UploadResult struct {
	AssetID      string
	PublicID     string
	Format       string
	ResourceType string
	URL          string
	SecureURL    string
	Version      int
	Bytes        int64
	Width        int
	Height       int
	CreatedAt    time.Time
}

type cloudinaryUploader struct {
	client        *cloudinary.Cloudinary
	defaultFolder string
}

// NewCloudinaryUploader initialises an uploader backed by Cloudinary.
func NewCloudinaryUploader(cloudName, apiKey, apiSecret, defaultFolder string) (Uploader, error) {
	client, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("create cloudinary client: %w", err)
	}

	return &cloudinaryUploader{client: client, defaultFolder: defaultFolder}, nil
}

// Upload pushes a file to Cloudinary and returns the metadata describing the stored asset.
func (u *cloudinaryUploader) Upload(ctx context.Context, file io.Reader, options UploadOptions) (*UploadResult, error) {
	params := uploader.UploadParams{Folder: u.defaultFolder}

	if options.Folder != "" {
		params.Folder = options.Folder
	}

	if options.PublicID != "" {
		params.PublicID = options.PublicID
	}

	if options.ResourceType != "" {
		params.ResourceType = options.ResourceType
	}

	if options.Overwrite != nil {
		params.Overwrite = options.Overwrite
	}

	result, err := u.client.Upload.Upload(ctx, file, params)
	if err != nil {
		return nil, fmt.Errorf("upload asset to cloudinary: %w", err)
	}

	return &UploadResult{
		AssetID:      result.AssetID,
		PublicID:     result.PublicID,
		Format:       result.Format,
		ResourceType: result.ResourceType,
		URL:          result.URL,
		SecureURL:    result.SecureURL,
		Version:      result.Version,
		Bytes:        int64(result.Bytes),
		Width:        result.Width,
		Height:       result.Height,
		CreatedAt:    result.CreatedAt,
	}, nil
}
