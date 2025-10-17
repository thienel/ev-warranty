package persistence_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/repositories"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/internal/infrastructure/persistence"
)

func TestOfficeRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OfficeRepository Suite")
}

var _ = Describe("OfficeRepository", func() {
	var (
		mock       sqlmock.Sqlmock
		db         *gorm.DB
		repository repositories.OfficeRepository
		ctx        context.Context
	)

	BeforeEach(func() {
		mock, db = persistence.SetupMockDB()
		repository = persistence.NewOfficeRepository(db)
		ctx = context.Background()
	})

	AfterEach(func() {
		CleanupMockDB(mock)
	})

	Describe("Create", func() {
		var office *entities.Office

		BeforeEach(func() {
			office = entities.NewOffice("Test Office", entities.OfficeTypeEVM, "123 Test St", true)
		})

		Context("when office is created successfully", func() {
			It("should return nil error", func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "offices"`)).
					WithArgs(
						office.OfficeName,
						office.OfficeType,
						office.Address,
						office.IsActive,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						office.ID,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(office.ID))
				mock.ExpectCommit()

				err := repository.Create(ctx, office)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a duplicate key constraint", func() {
			It("should return DBDuplicateKeyError", func() {
				MockDuplicateKeyError(mock, "offices", "office name")
				err := repository.Create(ctx, office)
				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDuplicateKey))
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockDBError(mock, `INSERT INTO "offices"`)
				err := repository.Create(ctx, office)
				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})
	})

	Describe("FindByID", func() {
		var officeID uuid.UUID

		BeforeEach(func() {
			officeID = uuid.New()
		})

		Context("when office is found", func() {
			It("should return the office", func() {
				rows := sqlmock.NewRows([]string{
					"id", "office_name", "office_type", "address", "is_active", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					officeID,
					"Test Office",
					entities.OfficeTypeEVM,
					"123 Test St",
					true,
					time.Now(),
					time.Now(),
					nil,
				)

				MockFindByID(mock, "offices", officeID, rows)

				office, err := repository.FindByID(ctx, officeID)
				Expect(err).NotTo(HaveOccurred())
				Expect(office).NotTo(BeNil())
				Expect(office.ID).To(Equal(officeID))
				Expect(office.OfficeName).To(Equal("Test Office"))
				Expect(office.OfficeType).To(Equal(entities.OfficeTypeEVM))
				Expect(office.Address).To(Equal("123 Test St"))
				Expect(office.IsActive).To(BeTrue())
			})
		})

		Context("when office is not found", func() {
			It("should return OfficeNotFound error", func() {
				MockNotFound(mock, "offices", officeID)
				office, err := repository.FindByID(ctx, officeID)
				Expect(err).To(HaveOccurred())
				Expect(office).To(BeNil())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeOfficeNotFound))
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockDBError(mock, `SELECT * FROM "offices" WHERE id = $1`)
				office, err := repository.FindByID(ctx, officeID)
				Expect(err).To(HaveOccurred())
				Expect(office).To(BeNil())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})
	})

	Describe("FindAll", func() {
		Context("when offices are found", func() {
			It("should return all offices", func() {
				officeID1 := uuid.New()
				officeID2 := uuid.New()
				rows := sqlmock.NewRows([]string{
					"id", "office_name", "office_type", "address", "is_active", "created_at", "updated_at", "deleted_at",
				}).AddRow(officeID1, "Office 1", entities.OfficeTypeEVM, "123 Test St", true, time.Now(), time.Now(), nil).
					AddRow(officeID2, "Office 2", entities.OfficeTypeSC, "456 Main St", false, time.Now(), time.Now(), nil)

				MockFindAll(mock, "offices", rows)
				offices, err := repository.FindAll(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(offices).To(HaveLen(2))
				Expect(offices[0].ID).To(Equal(officeID1))
				Expect(offices[0].OfficeName).To(Equal("Office 1"))
				Expect(offices[1].ID).To(Equal(officeID2))
				Expect(offices[1].OfficeName).To(Equal("Office 2"))
			})
		})

		Context("when no offices are found", func() {
			It("should return an empty slice", func() {
				rows := sqlmock.NewRows([]string{
					"id", "office_name", "office_type", "address", "is_active", "created_at", "updated_at", "deleted_at",
				})
				MockFindAll(mock, "offices", rows)
				offices, err := repository.FindAll(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(offices).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockDBError(mock, `SELECT * FROM "offices"`)
				offices, err := repository.FindAll(ctx)
				Expect(err).To(HaveOccurred())
				Expect(offices).To(BeNil())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})
	})

	Describe("Update", func() {
		var office *entities.Office

		BeforeEach(func() {
			office = entities.NewOffice("Updated Office", entities.OfficeTypeSC, "789 New St", false)
		})

		Context("when office is updated successfully", func() {
			It("should return nil error", func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "offices" SET`)).
					WithArgs(
						office.OfficeName,
						office.OfficeType,
						office.Address,
						office.IsActive,
						sqlmock.AnyArg(),
						office.ID,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := repository.Update(ctx, office)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockDBError(mock, `UPDATE "offices" SET`)

				err := repository.Update(ctx, office)
				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})
	})

	Describe("SoftDelete", func() {
		var officeID uuid.UUID

		BeforeEach(func() {
			officeID = uuid.New()
		})

		Context("when office is soft deleted successfully", func() {
			It("should return nil error", func() {
				MockSoftDelete(mock, "offices", officeID)
				err := repository.SoftDelete(ctx, officeID)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockDBError(mock, `UPDATE "offices" SET "deleted_at"=$1 WHERE id = $2`)
				err := repository.SoftDelete(ctx, officeID)
				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})
	})
})
