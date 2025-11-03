package persistence_test

import (
	"context"
	"errors"
	"ev-warranty-go/pkg/apperror"
	"ev-warranty-go/pkg/mocks"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/internal/infrastructure/persistence"
)

var _ = Describe("ClaimAttachmentRepository", func() {
	var (
		mock       sqlmock.Sqlmock
		db         *gorm.DB
		repository repositories.ClaimAttachmentRepository
		ctx        context.Context
	)

	BeforeEach(func() {
		mock, db = SetupMockDB()
		repository = persistence.NewClaimAttachmentRepository(db)
		ctx = context.Background()
	})

	AfterEach(func() {
		CleanupMockDB(mock)
	})

	Describe("Create", func() {
		var attachment *entity.ClaimAttachment

		BeforeEach(func() {
			attachment = newClaimAttachment()
		})

		Context("when attachment is created successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockSuccessfulInsert(mock, "claim_attachments", attachment.ID)

				err := repository.Create(mockTx, attachment)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a duplicate key constraint", func() {
			It("should return DBDuplicateKeyError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockDuplicateKeyError(mock, "claim_attachments", "claim_attachments_unique_key")

				err := repository.Create(mockTx, attachment)

				ExpectAppError(err, apperror.ErrorCodeDuplicateKey)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockInsertError(mock, "claim_attachments")

				err := repository.Create(mockTx, attachment)

				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})

		Context("boundary cases for attachment types", func() {
			It("should handle image type", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				attachment.Type = entity.AttachmentTypeImage
				MockSuccessfulInsert(mock, "claim_attachments", attachment.ID)

				err := repository.Create(mockTx, attachment)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should handle video type", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				attachment.Type = entity.AttachmentTypeVideo
				MockSuccessfulInsert(mock, "claim_attachments", attachment.ID)

				err := repository.Create(mockTx, attachment)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("boundary cases for URL", func() {
			It("should handle empty URL", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				attachment.URL = ""
				MockSuccessfulInsert(mock, "claim_attachments", attachment.ID)

				err := repository.Create(mockTx, attachment)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should handle very long URL", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				attachment.URL = "https://example.com/" + string(make([]byte, 1000))
				MockSuccessfulInsert(mock, "claim_attachments", attachment.ID)

				err := repository.Create(mockTx, attachment)

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("HardDelete", func() {
		var attachmentID uuid.UUID

		BeforeEach(func() {
			attachmentID = uuid.New()
		})

		Context("when attachment is hard deleted successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "claim_attachments" WHERE id = $1`)).
					WithArgs(attachmentID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := repository.HardDelete(mockTx, attachmentID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "claim_attachments" WHERE id = $1`)).
					WithArgs(attachmentID).
					WillReturnError(errors.New("database connection failed"))
				mock.ExpectRollback()

				err := repository.HardDelete(mockTx, attachmentID)

				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("SoftDeleteByClaimID", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when attachments are soft deleted successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "claim_attachments" SET "deleted_at"=$1 WHERE claim_id = $2`)).
					WithArgs(sqlmock.AnyArg(), claimID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := repository.SoftDeleteByClaimID(mockTx, claimID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "claim_attachments" SET "deleted_at"=$1 WHERE claim_id = $2`)).
					WithArgs(sqlmock.AnyArg(), claimID).
					WillReturnError(errors.New("database connection failed"))
				mock.ExpectRollback()

				err := repository.SoftDeleteByClaimID(mockTx, claimID)

				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("FindByID", func() {
		var attachmentID uuid.UUID

		BeforeEach(func() {
			attachmentID = uuid.New()
		})

		Context("when attachment is found", func() {
			It("should return the attachment", func() {
				expected := newClaimAttachment()
				expected.ID = attachmentID
				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "type", "url", "created_at", "deleted_at",
				}).AddRow(
					expected.ID, expected.ClaimID, expected.Type,
					expected.URL, expected.CreatedAt, expected.DeletedAt,
				)

				MockFindByID(mock, "claim_attachments", attachmentID, rows)

				attachment, err := repository.FindByID(ctx, attachmentID)

				Expect(err).NotTo(HaveOccurred())
				Expect(attachment).NotTo(BeNil())
				Expect(attachment.ID).To(Equal(expected.ID))
				Expect(attachment.ClaimID).To(Equal(expected.ClaimID))
				Expect(attachment.Type).To(Equal(expected.Type))
				Expect(attachment.URL).To(Equal(expected.URL))
			})
		})

		Context("when attachment is not found", func() {
			It("should return ClaimAttachmentNotFound error", func() {
				MockNotFound(mock, "claim_attachments", attachmentID)

				attachment, err := repository.FindByID(ctx, attachmentID)

				Expect(attachment).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeClaimAttachmentNotFound)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockQueryError(mock, `SELECT * FROM "claim_attachments" WHERE id = $1`)

				attachment, err := repository.FindByID(ctx, attachmentID)

				Expect(attachment).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("FindByClaimID", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when attachments are found", func() {
			It("should return all attachments for claim", func() {
				attachmentID1 := uuid.New()
				attachmentID2 := uuid.New()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "type", "url", "created_at", "deleted_at",
				}).AddRow(
					attachmentID1, claimID, entity.AttachmentTypeImage,
					"https://example.com/image1.jpg", time.Now(), nil,
				).AddRow(
					attachmentID2, claimID, entity.AttachmentTypeVideo,
					"https://example.com/video1.mp4", time.Now(), nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_attachments" WHERE claim_id = $1 AND "claim_attachments"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WithArgs(claimID).
					WillReturnRows(rows)

				attachments, err := repository.FindByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(attachments).To(HaveLen(2))
				Expect(attachments[0].ClaimID).To(Equal(claimID))
				Expect(attachments[1].ClaimID).To(Equal(claimID))
			})
		})

		Context("when no attachments are found", func() {
			It("should return empty slice", func() {
				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "type", "url", "created_at", "deleted_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_attachments" WHERE claim_id = $1 AND "claim_attachments"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WithArgs(claimID).
					WillReturnRows(rows)

				attachments, err := repository.FindByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(attachments).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_attachments" WHERE claim_id = $1 AND "claim_attachments"."deleted_at" IS NULL`)).
					WithArgs(claimID).
					WillReturnError(errors.New("database connection failed"))

				attachments, err := repository.FindByClaimID(ctx, claimID)

				Expect(attachments).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("CountByClaimID", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when count is successful", func() {
			It("should return the count", func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "claim_attachments" WHERE claim_id = $1 AND "claim_attachments"."deleted_at" IS NULL`)).
					WithArgs(claimID).
					WillReturnRows(rows)

				count, err := repository.CountByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(int64(5)))
			})
		})

		Context("when no attachments exist", func() {
			It("should return zero", func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "claim_attachments" WHERE claim_id = $1 AND "claim_attachments"."deleted_at" IS NULL`)).
					WithArgs(claimID).
					WillReturnRows(rows)

				count, err := repository.CountByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(int64(0)))
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "claim_attachments" WHERE claim_id = $1 AND "claim_attachments"."deleted_at" IS NULL`)).
					WithArgs(claimID).
					WillReturnError(errors.New("database connection failed"))

				count, err := repository.CountByClaimID(ctx, claimID)

				Expect(count).To(Equal(int64(0)))
				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})

		Context("boundary cases for count", func() {
			It("should handle very large counts", func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(999999)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "claim_attachments" WHERE claim_id = $1 AND "claim_attachments"."deleted_at" IS NULL`)).
					WithArgs(claimID).
					WillReturnRows(rows)

				count, err := repository.CountByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(int64(999999)))
			})
		})
	})

	Describe("FindByType", func() {
		var claimID uuid.UUID
		var attachmentType string

		BeforeEach(func() {
			claimID = uuid.New()
			attachmentType = entity.AttachmentTypeImage
		})

		Context("when attachments of type are found", func() {
			It("should return filtered attachments", func() {
				attachmentID1 := uuid.New()
				attachmentID2 := uuid.New()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "type", "url", "created_at", "deleted_at",
				}).AddRow(
					attachmentID1, claimID, entity.AttachmentTypeImage,
					"https://example.com/image1.jpg", time.Now(), nil,
				).AddRow(
					attachmentID2, claimID, entity.AttachmentTypeImage,
					"https://example.com/image2.jpg", time.Now(), nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_attachments" WHERE (claim_id = $1 AND attachment_type = $2) AND "claim_attachments"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WithArgs(claimID, attachmentType).
					WillReturnRows(rows)

				attachments, err := repository.FindByType(ctx, claimID, attachmentType)

				Expect(err).NotTo(HaveOccurred())
				Expect(attachments).To(HaveLen(2))
				Expect(attachments[0].Type).To(Equal(entity.AttachmentTypeImage))
				Expect(attachments[1].Type).To(Equal(entity.AttachmentTypeImage))
			})
		})

		Context("when no attachments of type are found", func() {
			It("should return empty slice", func() {
				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "type", "url", "created_at", "deleted_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_attachments" WHERE (claim_id = $1 AND attachment_type = $2) AND "claim_attachments"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WithArgs(claimID, attachmentType).
					WillReturnRows(rows)

				attachments, err := repository.FindByType(ctx, claimID, attachmentType)

				Expect(err).NotTo(HaveOccurred())
				Expect(attachments).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_attachments" WHERE (claim_id = $1 AND attachment_type = $2) AND "claim_attachments"."deleted_at" IS NULL`)).
					WithArgs(claimID, attachmentType).
					WillReturnError(errors.New("database connection failed"))

				attachments, err := repository.FindByType(ctx, claimID, attachmentType)

				Expect(attachments).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})

		Context("boundary cases for attachment type", func() {
			It("should filter by video type", func() {
				attachmentType = entity.AttachmentTypeVideo
				attachmentID := uuid.New()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "type", "url", "created_at", "deleted_at",
				}).AddRow(
					attachmentID, claimID, entity.AttachmentTypeVideo,
					"https://example.com/video.mp4", time.Now(), nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_attachments" WHERE (claim_id = $1 AND attachment_type = $2) AND "claim_attachments"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WithArgs(claimID, attachmentType).
					WillReturnRows(rows)

				attachments, err := repository.FindByType(ctx, claimID, attachmentType)

				Expect(err).NotTo(HaveOccurred())
				Expect(attachments).To(HaveLen(1))
				Expect(attachments[0].Type).To(Equal(entity.AttachmentTypeVideo))
			})

			It("should handle empty type string", func() {
				attachmentType = ""
				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "type", "url", "created_at", "deleted_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_attachments" WHERE (claim_id = $1 AND attachment_type = $2) AND "claim_attachments"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WithArgs(claimID, attachmentType).
					WillReturnRows(rows)

				attachments, err := repository.FindByType(ctx, claimID, attachmentType)

				Expect(err).NotTo(HaveOccurred())
				Expect(attachments).To(BeEmpty())
			})
		})
	})
})

func newClaimAttachment() *entity.ClaimAttachment {
	return &entity.ClaimAttachment{
		ID:        uuid.New(),
		ClaimID:   uuid.New(),
		Type:      entity.AttachmentTypeImage,
		URL:       "https://example.com/test-image.jpg",
		CreatedAt: time.Now(),
		DeletedAt: nil,
	}
}
