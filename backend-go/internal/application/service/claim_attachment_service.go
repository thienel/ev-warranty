package service

import (
	"context"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/internal/infrastructure/cloudinary"
	"ev-warranty-go/pkg/apperror"
	"ev-warranty-go/pkg/logger"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
)

type ClaimAttachmentService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ClaimAttachment, error)
	GetByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimAttachment, error)

	Create(tx application.Tx, claimID uuid.UUID, file multipart.File) (*entity.ClaimAttachment, error)
	HardDelete(tx application.Tx, claimID, attachmentID uuid.UUID) error
}

type claimAttachmentService struct {
	log          logger.Logger
	claimRepo    repositories.ClaimRepository
	attachRepo   repositories.ClaimAttachmentRepository
	cloudService cloudinary.CloudinaryService
}

func NewClaimAttachmentService(log logger.Logger, claimRepo repositories.ClaimRepository,
	attachRepo repositories.ClaimAttachmentRepository, cloudService cloudinary.CloudinaryService) ClaimAttachmentService {
	return &claimAttachmentService{
		log:          log,
		claimRepo:    claimRepo,
		attachRepo:   attachRepo,
		cloudService: cloudService,
	}
}

func (s *claimAttachmentService) GetByID(ctx context.Context, id uuid.UUID) (*entity.ClaimAttachment, error) {
	claimAttachment, err := s.attachRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return claimAttachment, nil
}

func (s *claimAttachmentService) GetByClaimID(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimAttachment, error) {
	claimAttachments, err := s.attachRepo.FindByClaimID(ctx, claimID)
	if err != nil {
		return nil, err
	}

	return claimAttachments, nil
}

func (s *claimAttachmentService) Create(tx application.Tx, claimID uuid.UUID, file multipart.File) (*entity.ClaimAttachment, error) {
	_, err := s.claimRepo.FindByID(tx.GetCtx(), claimID)
	if err != nil {
		return nil, err
	}

	mimeType, err := GetMimeType(file)
	if err != nil {
		return nil, err
	}

	attachType := cloudinary.DetermineResourceType(mimeType)
	if !entity.IsValidAttachmentType(attachType) {
		return nil, apperror.NewInvalidAttachmentType()
	}
	attachURL, err := s.cloudService.UploadFile(tx.GetCtx(), file, attachType)
	if err != nil {
		return nil, err
	}

	attachment := entity.NewClaimAttachment(claimID, attachType, attachURL)
	err = s.attachRepo.Create(tx, attachment)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func (s *claimAttachmentService) HardDelete(tx application.Tx, claimID, attachmentID uuid.UUID) error {
	claim, err := s.claimRepo.FindByID(tx.GetCtx(), claimID)
	if err != nil {
		return err
	}

	if claim.Status != entity.ClaimStatusDraft {
		return apperror.NewNotAllowDeleteClaim()
	}
	attach, err := s.attachRepo.FindByID(tx.GetCtx(), attachmentID)
	if err != nil {
		return err
	}

	err = s.attachRepo.HardDelete(tx, attachmentID)
	if err == nil {
		if cloudErr := s.cloudService.DeleteFileByURL(tx.GetCtx(), attach.URL); cloudErr != nil {
			s.log.Error("[Cloudinary] Failed to delete file when hard delete claim attachment", "error", cloudErr)
		}
	}

	return err
}

func GetMimeType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", apperror.NewInternalServerError(err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", apperror.NewInternalServerError(err)
	}

	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}
