package cloudinary

import (
	"context"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/infrastructure/config"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService interface {
	UploadFile(ctx context.Context, file multipart.File, resourceType string) (string, error)
	DeleteFile(ctx context.Context, publicID string, resourceType string) error
	DeleteFileByURL(ctx context.Context, fileURL string) error
}

type cloudinaryService struct {
	cld          *cloudinary.Cloudinary
	uploadFolder string
}

func NewCloudinaryService(cfg *config.CloudinaryConfig) (CloudinaryService, error) {
	cld, err := cloudinary.NewFromURL(cfg.URL)
	if err != nil {
		return nil, apperrors.NewFailedInitializeCloudinary()
	}
	cld.Config.URL.Secure = true

	return &cloudinaryService{
		cld:          cld,
		uploadFolder: cfg.UploadFolder,
	}, nil
}

func (s *cloudinaryService) UploadFile(ctx context.Context, file multipart.File, resourceType string) (string, error) {
	uploadParams := uploader.UploadParams{
		ResourceType: resourceType,
	}

	if s.uploadFolder != "" {
		uploadParams.Folder = s.uploadFolder
	}

	resp, err := s.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", apperrors.NewFailedUploadCloudinary()
	}

	return resp.SecureURL, nil
}

func (s *cloudinaryService) DeleteFile(ctx context.Context, publicID string, resourceType string) error {
	if publicID == "" {
		return apperrors.NewEmptyCloudinaryParameter("publicID")
	}

	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: resourceType,
	})
	if err != nil {
		return apperrors.NewFailedDeleteCloudinary()
	}
	return nil
}

func (s *cloudinaryService) DeleteFileByURL(ctx context.Context, fileURL string) error {
	if fileURL == "" {
		return apperrors.NewEmptyCloudinaryParameter("fileURL")
	}

	publicID, resourceType, err := parseCloudinaryURL(fileURL)
	if err != nil {
		return err
	}

	return s.DeleteFile(ctx, publicID, resourceType)
}

func parseCloudinaryURL(fileURL string) (publicID, resourceType string, err error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", "", apperrors.NewInvalidCloudinaryURL()
	}

	parts := strings.Split(strings.TrimPrefix(parsedURL.Path, "/"), "/")

	if len(parts) < 5 {
		return "", "", apperrors.NewInvalidCloudinaryURL()
	}

	resourceType = parts[1]

	pathWithVersion := strings.Join(parts[3:], "/")

	pathWithoutVersion := pathWithVersion
	if strings.HasPrefix(pathWithVersion, "v") {
		versionParts := strings.SplitN(pathWithVersion, "/", 2)
		if len(versionParts) == 2 {
			pathWithoutVersion = versionParts[1]
		}
	}

	publicID = strings.TrimSuffix(pathWithoutVersion, filepath.Ext(pathWithoutVersion))

	return publicID, resourceType, nil
}

func DetermineResourceType(mimeType string) string {
	if strings.HasPrefix(mimeType, "image/") {
		return "image"
	}
	if strings.HasPrefix(mimeType, "video/") {
		return "video"
	}
	return "auto"
}
