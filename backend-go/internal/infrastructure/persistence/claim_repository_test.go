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

	"ev-warranty-go/internal/application/repository"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/internal/infrastructure/persistence"
)

var _ = Describe("ClaimRepository", func() {
	var (
		mock       sqlmock.Sqlmock
		db         *gorm.DB
		repository repository.ClaimRepository
		ctx        context.Context
	)

	BeforeEach(func() {
		mock, db = SetupMockDB()
		repository = persistence.NewClaimRepository(db)
		ctx = context.Background()
	})

	AfterEach(func() {
		CleanupMockDB(mock)
	})

	Describe("Create", func() {
		var claim *entity.Claim

		BeforeEach(func() {
			claim = newClaim()
		})

		Context("when claim is created successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockSuccessfulInsert(mock, "claims", claim.ID)

				err := repository.Create(mockTx, claim)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a duplicate key constraint", func() {
			It("should return DBDuplicateKeyError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockDuplicateKeyError(mock, "claims", "claims_unique_key")

				err := repository.Create(mockTx, claim)

				ExpectAppError(err, apperror.ErrorCodeDuplicateKey)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				MockInsertError(mock, "claims")

				err := repository.Create(mockTx, claim)

				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("Update", func() {
		var claim *entity.Claim

		BeforeEach(func() {
			claim = newClaim()
			claim.Description = "Updated description"
			claim.Status = entity.ClaimStatusApproved
		})

		Context("when claim is updated successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				MockSuccessfulUpdate(mock, "claims")

				err := repository.Update(mockTx, claim)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				MockUpdateError(mock, "claims")

				err := repository.Update(mockTx, claim)

				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("HardDelete", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when claim is hard deleted successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "claims" WHERE id = $1`)).
					WithArgs(claimID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := repository.HardDelete(mockTx, claimID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "claims" WHERE id = $1`)).
					WithArgs(claimID).
					WillReturnError(errors.New("database connection failed"))
				mock.ExpectRollback()

				err := repository.HardDelete(mockTx, claimID)

				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("SoftDelete", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when claim is soft deleted successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockSoftDelete(mock, "claims", claimID)

				err := repository.SoftDelete(mockTx, claimID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockDeleteError(mock, "claims")

				err := repository.SoftDelete(mockTx, claimID)

				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("UpdateStatus", func() {
		var claimID uuid.UUID
		var status string

		BeforeEach(func() {
			claimID = uuid.New()
			status = entity.ClaimStatusApproved
		})

		Context("when status is updated successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "claims" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND "claims"."deleted_at" IS NULL`)).
					WithArgs(status, sqlmock.AnyArg(), claimID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := repository.UpdateStatus(mockTx, claimID, status)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "claims" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND "claims"."deleted_at" IS NULL`)).
					WithArgs(status, sqlmock.AnyArg(), claimID).
					WillReturnError(errors.New("database connection failed"))
				mock.ExpectRollback()

				err := repository.UpdateStatus(mockTx, claimID, status)

				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})

		Context("boundary cases for status", func() {
			It("should handle all valid status values", func() {
				statuses := []string{
					entity.ClaimStatusDraft,
					entity.ClaimStatusSubmitted,
					entity.ClaimStatusReviewing,
					entity.ClaimStatusApproved,
					entity.ClaimStatusRejected,
					entity.ClaimStatusCancelled,
				}

				for _, s := range statuses {
					mockTx := mocks.NewTx(GinkgoT())
					mockTx.EXPECT().GetTx().Return(db)
					mock.ExpectBegin()
					mock.ExpectExec(regexp.QuoteMeta(`UPDATE "claims" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND "claims"."deleted_at" IS NULL`)).
						WithArgs(s, sqlmock.AnyArg(), claimID).
						WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()

					err := repository.UpdateStatus(mockTx, claimID, s)
					Expect(err).NotTo(HaveOccurred())
				}
			})
		})
	})

	Describe("FindByID", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when claim is found", func() {
			It("should return the claim", func() {
				expected := newClaim()
				expected.ID = claimID
				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "customer_id", "description", "status",
					"total_cost", "approved_by", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					expected.ID, expected.VehicleID, expected.CustomerID,
					expected.Description, expected.Status, expected.TotalCost,
					expected.ApprovedBy, expected.CreatedAt, expected.UpdatedAt, expected.DeletedAt,
				)

				MockFindByID(mock, "claims", claimID, rows)

				claim, err := repository.FindByID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(claim).NotTo(BeNil())
				Expect(claim.ID).To(Equal(expected.ID))
				Expect(claim.VehicleID).To(Equal(expected.VehicleID))
				Expect(claim.CustomerID).To(Equal(expected.CustomerID))
				Expect(claim.Description).To(Equal(expected.Description))
				Expect(claim.Status).To(Equal(expected.Status))
			})
		})

		Context("when claim is not found", func() {
			It("should return ClaimNotFound error", func() {
				MockNotFound(mock, "claims", claimID)

				claim, err := repository.FindByID(ctx, claimID)

				Expect(claim).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeClaimNotFound)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockQueryError(mock, `SELECT * FROM "claims" WHERE id = $1`)

				claim, err := repository.FindByID(ctx, claimID)

				Expect(claim).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("FindAll", func() {
		Context("when claims are found", func() {
			It("should return all claims", func() {
				claimID1 := uuid.New()
				claimID2 := uuid.New()
				vehicleID := uuid.New()
				customerID := uuid.New()

				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "customer_id", "description", "status",
					"total_cost", "approved_by", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					claimID1, vehicleID, customerID, "Claim 1", entity.ClaimStatusDraft,
					1000.0, nil, time.Now(), time.Now(), nil,
				).AddRow(
					claimID2, vehicleID, customerID, "Claim 2", entity.ClaimStatusSubmitted,
					2000.0, nil, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claims" WHERE "claims"."deleted_at" IS NULL`)).
					WillReturnRows(rows)

				claims, err := repository.FindAll(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(HaveLen(2))
			})
		})

		Context("when no claims are found", func() {
			It("should return empty slice", func() {
				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "customer_id", "description", "status",
					"total_cost", "approved_by", "created_at", "updated_at", "deleted_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claims" WHERE "claims"."deleted_at" IS NULL`)).
					WillReturnRows(rows)

				claims, err := repository.FindAll(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claims" WHERE "claims"."deleted_at" IS NULL`)).
					WillReturnError(errors.New("database connection failed"))

				claims, err := repository.FindAll(ctx)

				Expect(claims).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("FindByCustomerID", func() {
		var customerID uuid.UUID

		BeforeEach(func() {
			customerID = uuid.New()
		})

		Context("when claims are found", func() {
			It("should return all claims for customer", func() {
				claimID1 := uuid.New()
				claimID2 := uuid.New()
				vehicleID := uuid.New()
				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "customer_id", "description", "status",
					"total_cost", "approved_by", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					claimID1, vehicleID, customerID, "Claim 1", entity.ClaimStatusDraft,
					1000.0, nil, time.Now(), time.Now(), nil,
				).AddRow(
					claimID2, vehicleID, customerID, "Claim 2", entity.ClaimStatusSubmitted,
					2000.0, nil, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claims" WHERE customer_id = $1 AND "claims"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WithArgs(customerID).
					WillReturnRows(rows)

				claims, err := repository.FindByCustomerID(ctx, customerID)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(HaveLen(2))
				Expect(claims[0].CustomerID).To(Equal(customerID))
				Expect(claims[1].CustomerID).To(Equal(customerID))
			})
		})

		Context("when no claims are found", func() {
			It("should return empty slice", func() {
				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "customer_id", "description", "status",
					"total_cost", "approved_by", "created_at", "updated_at", "deleted_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claims" WHERE customer_id = $1 AND "claims"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WithArgs(customerID).
					WillReturnRows(rows)

				claims, err := repository.FindByCustomerID(ctx, customerID)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claims" WHERE customer_id = $1 AND "claims"."deleted_at" IS NULL`)).
					WithArgs(customerID).
					WillReturnError(errors.New("database connection failed"))

				claims, err := repository.FindByCustomerID(ctx, customerID)

				Expect(claims).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("FindByVehicleID", func() {
		var vehicleID uuid.UUID

		BeforeEach(func() {
			vehicleID = uuid.New()
		})

		Context("when claims are found", func() {
			It("should return all claims for vehicle", func() {
				claimID1 := uuid.New()
				claimID2 := uuid.New()
				customerID := uuid.New()
				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "customer_id", "description", "status",
					"total_cost", "approved_by", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					claimID1, vehicleID, customerID, "Claim 1", entity.ClaimStatusDraft,
					1000.0, nil, time.Now(), time.Now(), nil,
				).AddRow(
					claimID2, vehicleID, customerID, "Claim 2", entity.ClaimStatusSubmitted,
					2000.0, nil, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claims" WHERE vehicle_id = $1 AND "claims"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WithArgs(vehicleID).
					WillReturnRows(rows)

				claims, err := repository.FindByVehicleID(ctx, vehicleID)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(HaveLen(2))
				Expect(claims[0].VehicleID).To(Equal(vehicleID))
				Expect(claims[1].VehicleID).To(Equal(vehicleID))
			})
		})

		Context("when no claims are found", func() {
			It("should return empty slice", func() {
				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "customer_id", "description", "status",
					"total_cost", "approved_by", "created_at", "updated_at", "deleted_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claims" WHERE vehicle_id = $1 AND "claims"."deleted_at" IS NULL ORDER BY created_at DESC`)).
					WithArgs(vehicleID).
					WillReturnRows(rows)

				claims, err := repository.FindByVehicleID(ctx, vehicleID)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claims" WHERE vehicle_id = $1 AND "claims"."deleted_at" IS NULL`)).
					WithArgs(vehicleID).
					WillReturnError(errors.New("database connection failed"))

				claims, err := repository.FindByVehicleID(ctx, vehicleID)

				Expect(claims).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})
})

func newClaim() *entity.Claim {
	return &entity.Claim{
		ID:          uuid.New(),
		VehicleID:   uuid.New(),
		CustomerID:  uuid.New(),
		Description: "Test claim description",
		Status:      entity.ClaimStatusDraft,
		TotalCost:   1000.0,
		ApprovedBy:  nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
}
