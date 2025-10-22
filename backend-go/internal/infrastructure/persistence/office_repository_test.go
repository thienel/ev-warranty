package persistence_test

import (
	"context"
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

var _ = Describe("OfficeRepository", func() {
	var (
		mock       sqlmock.Sqlmock
		db         *gorm.DB
		repository repositories.OfficeRepository
		ctx        context.Context
	)

	BeforeEach(func() {
		mock, db = SetupMockDB()
		repository = persistence.NewOfficeRepository(db)
		ctx = context.Background()
	})

	AfterEach(func() {
		CleanupMockDB(mock)
	})

	Describe("Create", func() {
		var office *entities.Office

		BeforeEach(func() {
			office = newOffice()
		})

		Context("when office is created successfully", func() {
			It("should return nil error", func() {
				MockSuccessfulInsert(mock, "offices", office.ID)

				err := repository.Create(ctx, office)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a duplicate key constraint", func() {
			It("should return DBDuplicateKeyError", func() {
				MockDuplicateKeyError(mock, "offices", "offices_office_name_key")

				err := repository.Create(ctx, office)

				ExpectAppError(err, apperrors.ErrorCodeDuplicateKey)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockInsertError(mock, "offices")

				err := repository.Create(ctx, office)

				ExpectAppError(err, apperrors.ErrorCodeDBOperation)
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
				expected := newOffice()
				expected.ID = officeID
				rows := sqlmock.NewRows([]string{
					"id", "office_name", "office_type", "address", "is_active",
					"created_at", "updated_at", "deleted_at",
				}).AddRow(
					expected.ID, expected.OfficeName, expected.OfficeType,
					expected.Address, expected.IsActive, expected.CreatedAt,
					expected.UpdatedAt, expected.DeletedAt,
				)

				MockFindByID(mock, "offices", officeID, rows)

				office, err := repository.FindByID(ctx, officeID)

				Expect(err).NotTo(HaveOccurred())
				Expect(office).NotTo(BeNil())
				Expect(office.ID).To(Equal(expected.ID))
				Expect(office.OfficeName).To(Equal(expected.OfficeName))
				Expect(office.OfficeType).To(Equal(expected.OfficeType))
				Expect(office.Address).To(Equal(expected.Address))
				Expect(office.IsActive).To(Equal(expected.IsActive))
			})
		})

		Context("when office is not found", func() {
			It("should return OfficeNotFound error", func() {
				MockNotFound(mock, "offices", officeID)

				office, err := repository.FindByID(ctx, officeID)

				Expect(office).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeOfficeNotFound)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockQueryError(mock, `SELECT * FROM "offices" WHERE id = $1`)

				office, err := repository.FindByID(ctx, officeID)

				Expect(office).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeDBOperation)
			})
		})
	})

	Describe("FindAll", func() {
		Context("when offices are found", func() {
			It("should return all offices", func() {
				officeID1 := uuid.New()
				officeID2 := uuid.New()
				rows := sqlmock.NewRows([]string{
					"id", "office_name", "office_type", "address", "is_active",
					"created_at", "updated_at", "deleted_at",
				}).AddRow(
					officeID1, "Office 1", entities.OfficeTypeEVM, "123 Test St",
					true, time.Now(), time.Now(), nil,
				).AddRow(
					officeID2, "Office 2", entities.OfficeTypeSC, "456 Main St",
					false, time.Now(), time.Now(), nil,
				)

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
					"id", "office_name", "office_type", "address", "is_active",
					"created_at", "updated_at", "deleted_at",
				})

				MockFindAll(mock, "offices", rows)

				offices, err := repository.FindAll(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(offices).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockQueryError(mock, `SELECT * FROM "offices"`)

				offices, err := repository.FindAll(ctx)

				Expect(offices).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeDBOperation)
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
				MockSuccessfulUpdate(mock, "offices")

				err := repository.Update(ctx, office)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockUpdateError(mock, "offices")

				err := repository.Update(ctx, office)

				ExpectAppError(err, apperrors.ErrorCodeDBOperation)
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
				MockDeleteError(mock, "offices")

				err := repository.SoftDelete(ctx, officeID)

				ExpectAppError(err, apperrors.ErrorCodeDBOperation)
			})
		})
	})
})

func newOffice() *entities.Office {
	return &entities.Office{
		ID:         uuid.New(),
		OfficeName: "Test Office",
		OfficeType: entities.OfficeTypeEVM,
		Address:    "123 Test St",
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		DeletedAt:  nil,
	}
}
