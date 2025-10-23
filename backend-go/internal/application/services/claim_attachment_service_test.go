package services_test

import (
	"bytes"
	"context"
	"errors"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/pkg/mocks"
	"io"
	"mime/multipart"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("ClaimAttachmentService", func() {
	var (
		mockLogger     *mocks.Logger
		mockClaimRepo  *mocks.ClaimRepository
		mockAttachRepo *mocks.ClaimAttachmentRepository
		mockCloudServ  *mocks.CloudinaryService
		mockTx         *mocks.Tx
		service        services.ClaimAttachmentService
		ctx            context.Context
	)

	BeforeEach(func() {
		mockLogger = mocks.NewLogger(GinkgoT())
		mockClaimRepo = mocks.NewClaimRepository(GinkgoT())
		mockAttachRepo = mocks.NewClaimAttachmentRepository(GinkgoT())
		mockCloudServ = mocks.NewCloudinaryService(GinkgoT())
		mockTx = mocks.NewTx(GinkgoT())
		service = services.NewClaimAttachmentService(mockLogger, mockClaimRepo, mockAttachRepo, mockCloudServ)
		ctx = context.Background()
	})

	Describe("GetByID", func() {
		var attachmentID uuid.UUID

		BeforeEach(func() {
			attachmentID = uuid.New()
		})

		Context("when attachment is found", func() {
			It("should return the attachment", func() {
				expectedAttachment := &entities.ClaimAttachment{
					ID:      attachmentID,
					ClaimID: uuid.New(),
					Type:    "image",
					URL:     "https://example.com/image.jpg",
				}

				mockAttachRepo.EXPECT().FindByID(ctx, attachmentID).Return(expectedAttachment, nil).Once()

				attachment, err := service.GetByID(ctx, attachmentID)

				Expect(err).NotTo(HaveOccurred())
				Expect(attachment).NotTo(BeNil())
				Expect(attachment.ID).To(Equal(expectedAttachment.ID))
				Expect(attachment.URL).To(Equal(expectedAttachment.URL))
			})
		})

		Context("when attachment is not found", func() {
			It("should return ClaimAttachmentNotFound error", func() {
				notFoundErr := apperrors.New(404, apperrors.ErrorCodeClaimAttachmentNotFound, errors.New("attachment not found"))
				mockAttachRepo.EXPECT().FindByID(ctx, attachmentID).Return(nil, notFoundErr).Once()

				attachment, err := service.GetByID(ctx, attachmentID)

				Expect(attachment).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeClaimAttachmentNotFound)
			})
		})
	})

	Describe("GetByClaimID", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when attachments are found", func() {
			It("should return all attachments for the claim", func() {
				expectedAttachments := []*entities.ClaimAttachment{
					{
						ID:      uuid.New(),
						ClaimID: claimID,
						Type:    "image",
						URL:     "https://example.com/image1.jpg",
					},
					{
						ID:      uuid.New(),
						ClaimID: claimID,
						Type:    "video",
						URL:     "https://example.com/video1.mp4",
					},
				}

				mockAttachRepo.EXPECT().FindByClaimID(ctx, claimID).Return(expectedAttachments, nil).Once()

				attachments, err := service.GetByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(attachments).NotTo(BeNil())
				Expect(attachments).To(HaveLen(2))
				Expect(attachments[0].ClaimID).To(Equal(claimID))
				Expect(attachments[1].ClaimID).To(Equal(claimID))
			})
		})

		Context("when no attachments are found", func() {
			It("should return an empty slice", func() {
				mockAttachRepo.EXPECT().FindByClaimID(ctx, claimID).Return([]*entities.ClaimAttachment{}, nil).Once()

				attachments, err := service.GetByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(attachments).NotTo(BeNil())
				Expect(attachments).To(BeEmpty())
			})
		})

		Context("when repository returns error", func() {
			It("should return the error", func() {
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))
				mockAttachRepo.EXPECT().FindByClaimID(ctx, claimID).Return(nil, dbErr).Once()

				attachments, err := service.GetByClaimID(ctx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(attachments).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("Create", func() {
		var (
			claimID uuid.UUID
			file    multipart.File
		)

		BeforeEach(func() {
			claimID = uuid.New()
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when attachment is created successfully with image", func() {
			It("should upload and create attachment", func() {
				jpegHeader := []byte{0xFF, 0xD8, 0xFF}
				fileContent := append(jpegHeader, make([]byte, 509)...)
				file = &mockFile{Reader: bytes.NewReader(fileContent)}

				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockCloudServ.EXPECT().UploadFile(ctx, file, "image").Return("https://example.com/image.jpg", nil).Once()
				mockAttachRepo.EXPECT().Create(mockTx, mock.MatchedBy(func(a *entities.ClaimAttachment) bool {
					return a.ClaimID == claimID &&
						a.Type == "image" &&
						a.URL == "https://example.com/image.jpg"
				})).Return(nil).Once()

				attachment, err := service.Create(mockTx, claimID, file)

				Expect(err).NotTo(HaveOccurred())
				Expect(attachment).NotTo(BeNil())
				Expect(attachment.ClaimID).To(Equal(claimID))
				Expect(attachment.Type).To(Equal("image"))
			})
		})

		Context("when attachment is created successfully with PNG image", func() {
			It("should upload and create PNG attachment", func() {
				pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
				fileContent := append(pngHeader, make([]byte, 504)...)
				file = &mockFile{Reader: bytes.NewReader(fileContent)}

				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockCloudServ.EXPECT().UploadFile(ctx, file, "image").Return("https://example.com/image.png", nil).Once()
				mockAttachRepo.EXPECT().Create(mockTx, mock.AnythingOfType("*entities.ClaimAttachment")).Return(nil).Once()

				attachment, err := service.Create(mockTx, claimID, file)

				Expect(err).NotTo(HaveOccurred())
				Expect(attachment).NotTo(BeNil())
			})
		})

		Context("when claim is not found", func() {
			It("should return ClaimNotFound error", func() {
				file = &mockFile{Reader: bytes.NewReader([]byte("test"))}
				notFoundErr := apperrors.New(404, apperrors.ErrorCodeClaimNotFound, errors.New("claim not found"))
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				attachment, err := service.Create(mockTx, claimID, file)

				Expect(err).To(HaveOccurred())
				Expect(attachment).To(BeNil())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when file read fails", func() {
			It("should return error", func() {
				file = &mockFile{Reader: &errorReader{}}
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				attachment, err := service.Create(mockTx, claimID, file)

				Expect(err).To(HaveOccurred())
				Expect(attachment).To(BeNil())
			})
		})

		Context("when file seek fails", func() {
			It("should return error", func() {
				file = &mockFile{Reader: bytes.NewReader([]byte("test")), seekError: true}
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				attachment, err := service.Create(mockTx, claimID, file)

				Expect(err).To(HaveOccurred())
				Expect(attachment).To(BeNil())
			})
		})

		Context("when attachment type is invalid", func() {
			It("should return InvalidAttachmentType error", func() {
				textContent := []byte("plain text file")
				fileContent := append(textContent, make([]byte, 497)...)
				file = &mockFile{Reader: bytes.NewReader(fileContent)}

				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				attachment, err := service.Create(mockTx, claimID, file)

				Expect(attachment).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeInvalidAttachmentType)
			})
		})

		Context("when cloud upload fails", func() {
			It("should return error", func() {
				jpegHeader := []byte{0xFF, 0xD8, 0xFF}
				fileContent := append(jpegHeader, make([]byte, 509)...)
				file = &mockFile{Reader: bytes.NewReader(fileContent)}

				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				cloudErr := errors.New("cloud upload failed")

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockCloudServ.EXPECT().UploadFile(ctx, file, "image").Return("", cloudErr).Once()

				attachment, err := service.Create(mockTx, claimID, file)

				Expect(err).To(HaveOccurred())
				Expect(attachment).To(BeNil())
				Expect(err).To(Equal(cloudErr))
			})
		})

		Context("when repository create fails", func() {
			It("should return error", func() {
				jpegHeader := []byte{0xFF, 0xD8, 0xFF}
				fileContent := append(jpegHeader, make([]byte, 509)...)
				file = &mockFile{Reader: bytes.NewReader(fileContent)}

				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockCloudServ.EXPECT().UploadFile(ctx, file, "image").Return("https://example.com/image.jpg", nil).Once()
				mockAttachRepo.EXPECT().Create(mockTx, mock.AnythingOfType("*entities.ClaimAttachment")).Return(dbErr).Once()

				attachment, err := service.Create(mockTx, claimID, file)

				Expect(err).To(HaveOccurred())
				Expect(attachment).To(BeNil())
				Expect(err).To(Equal(dbErr))
			})
		})
	})

	Describe("HardDelete", func() {
		var (
			claimID      uuid.UUID
			attachmentID uuid.UUID
		)

		BeforeEach(func() {
			claimID = uuid.New()
			attachmentID = uuid.New()
			mockTx.EXPECT().GetCtx().Return(ctx).Maybe()
		})

		Context("when attachment is deleted successfully", func() {
			It("should delete from repository and cloud storage", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				attachment := &entities.ClaimAttachment{
					ID:      attachmentID,
					ClaimID: claimID,
					URL:     "https://example.com/image.jpg",
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockAttachRepo.EXPECT().FindByID(ctx, attachmentID).Return(attachment, nil).Once()
				mockAttachRepo.EXPECT().HardDelete(mockTx, attachmentID).Return(nil).Once()
				mockCloudServ.EXPECT().DeleteFileByURL(ctx, attachment.URL).Return(nil).Maybe()
				mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Maybe()

				err := service.HardDelete(mockTx, claimID, attachmentID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when claim status is not draft", func() {
			It("should return NotAllowDeleteClaim error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusSubmitted,
				}

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()

				err := service.HardDelete(mockTx, claimID, attachmentID)

				ExpectAppError(err, apperrors.ErrorCodeClaimStatusNotAllowedDelete)
			})
		})

		Context("when attachment is not found", func() {
			It("should return ClaimAttachmentNotFound error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				notFoundErr := apperrors.New(404, apperrors.ErrorCodeClaimAttachmentNotFound, errors.New("attachment not found"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockAttachRepo.EXPECT().FindByID(ctx, attachmentID).Return(nil, notFoundErr).Once()

				err := service.HardDelete(mockTx, claimID, attachmentID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when repository delete fails", func() {
			It("should return the error", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				attachment := &entities.ClaimAttachment{
					ID:      attachmentID,
					ClaimID: claimID,
					URL:     "https://example.com/image.jpg",
				}
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockAttachRepo.EXPECT().FindByID(ctx, attachmentID).Return(attachment, nil).Once()
				mockAttachRepo.EXPECT().HardDelete(mockTx, attachmentID).Return(dbErr).Once()

				err := service.HardDelete(mockTx, claimID, attachmentID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when cloud delete fails", func() {
			It("should log error but still return nil", func() {
				claim := &entities.Claim{
					ID:     claimID,
					Status: entities.ClaimStatusDraft,
				}
				attachment := &entities.ClaimAttachment{
					ID:      attachmentID,
					ClaimID: claimID,
					URL:     "https://example.com/image.jpg",
				}
				cloudErr := errors.New("cloud storage error")

				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(claim, nil).Once()
				mockAttachRepo.EXPECT().FindByID(ctx, attachmentID).Return(attachment, nil).Once()
				mockAttachRepo.EXPECT().HardDelete(mockTx, attachmentID).Return(nil).Once()
				mockCloudServ.EXPECT().DeleteFileByURL(ctx, attachment.URL).Return(cloudErr).Once()
				mockLogger.EXPECT().Error("[Cloudinary] Failed to delete file when hard delete claim attachment",
					"error", cloudErr).Once()

				err := service.HardDelete(mockTx, claimID, attachmentID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when claim is not found", func() {
			It("should return ClaimNotFound error", func() {
				notFoundErr := apperrors.New(404, apperrors.ErrorCodeClaimNotFound, errors.New("claim not found"))
				mockClaimRepo.EXPECT().FindByID(ctx, claimID).Return(nil, notFoundErr).Once()

				err := service.HardDelete(mockTx, claimID, attachmentID)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(notFoundErr))
			})
		})
	})
})

type mockFile struct {
	io.Reader
	seekError bool
}

func (m *mockFile) Read(p []byte) (n int, err error) {
	return m.Reader.Read(p)
}

func (m *mockFile) Seek(offset int64, whence int) (int64, error) {
	if m.seekError {
		return 0, errors.New("seek error")
	}
	if seeker, ok := m.Reader.(io.Seeker); ok {
		return seeker.Seek(offset, whence)
	}
	return 0, nil
}

func (m *mockFile) Close() error {
	return nil
}

func (m *mockFile) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, nil
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}
