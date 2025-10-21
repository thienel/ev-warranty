package persistence_test

import (
	"context"
	"errors"
	"ev-warranty-go/pkg/mocks"
	"regexp"
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

var _ = Describe("ClaimRepository", func() {
	var (
		mock       sqlmock.Sqlmock
		db         *gorm.DB
		repository repositories.ClaimRepository
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
		var claim *entities.Claim

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

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDuplicateKey))
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				MockInsertError(mock, "claims")

				err := repository.Create(mockTx, claim)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})
	})

	Describe("Update", func() {
		var claim *entities.Claim

		BeforeEach(func() {
			claim = newClaim()
			claim.Description = "Updated description"
			claim.Status = entities.ClaimStatusApproved
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

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
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

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
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

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})
	})

	Describe("UpdateStatus", func() {
		var claimID uuid.UUID
		var status string

		BeforeEach(func() {
			claimID = uuid.New()
			status = entities.ClaimStatusApproved
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

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})

		Context("boundary cases for status", func() {
			It("should handle all valid status values", func() {
				statuses := []string{
					entities.ClaimStatusDraft,
					entities.ClaimStatusSubmitted,
					entities.ClaimStatusReviewing,
					entities.ClaimStatusApproved,
					entities.ClaimStatusRejected,
					entities.ClaimStatusCancelled,
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

				Expect(err).To(HaveOccurred())
				Expect(claim).To(BeNil())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeClaimNotFound))
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockQueryError(mock, `SELECT * FROM "claims" WHERE id = $1`)

				claim, err := repository.FindByID(ctx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(claim).To(BeNil())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})
	})

	Describe("FindAll", func() {
		var filters repositories.ClaimFilters
		var pagination repositories.Pagination

		BeforeEach(func() {
			filters = repositories.ClaimFilters{}
			pagination = repositories.Pagination{
				Page:     1,
				PageSize: 10,
				SortBy:   "created_at",
				SortDir:  "DESC",
			}
		})

		Context("when claims are found without filters", func() {
			It("should return all claims with pagination", func() {
				claimID1 := uuid.New()
				claimID2 := uuid.New()
				vehicleID := uuid.New()
				customerID := uuid.New()

				countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
				mock.ExpectQuery(`SELECT count\(\*\) FROM "claims" WHERE "claims"."deleted_at" IS NULL`).
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "customer_id", "description", "status",
					"total_cost", "approved_by", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					claimID1, vehicleID, customerID, "Claim 1", entities.ClaimStatusDraft,
					1000.0, nil, time.Now(), time.Now(), nil,
				).AddRow(
					claimID2, vehicleID, customerID, "Claim 2", entities.ClaimStatusSubmitted,
					2000.0, nil, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "claims" WHERE "claims"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT \$1`).
					WithArgs(10).
					WillReturnRows(rows)

				claims, total, err := repository.FindAll(ctx, filters, pagination)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(HaveLen(2))
				Expect(total).To(Equal(int64(2)))
			})
		})

		Context("when filtering by customer ID", func() {
			It("should return filtered claims", func() {
				customerID := uuid.New()
				filters.CustomerID = &customerID

				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(`SELECT count\(\*\) FROM "claims" WHERE customer_id = \$1 AND "claims"."deleted_at" IS NULL`).
					WithArgs(customerID).
					WillReturnRows(countRows)

				claimID := uuid.New()
				vehicleID := uuid.New()
				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "customer_id", "description", "status",
					"total_cost", "approved_by", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					claimID, vehicleID, customerID, "Claim", entities.ClaimStatusDraft,
					1000.0, nil, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "claims" WHERE customer_id = \$1 AND "claims"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT \$2`).
					WithArgs(customerID, 10).
					WillReturnRows(rows)

				claims, total, err := repository.FindAll(ctx, filters, pagination)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(HaveLen(1))
				Expect(total).To(Equal(int64(1)))
				Expect(claims[0].CustomerID).To(Equal(customerID))
			})
		})

		Context("when filtering by vehicle ID", func() {
			It("should return filtered claims", func() {
				vehicleID := uuid.New()
				filters.VehicleID = &vehicleID

				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(`SELECT count\(\*\) FROM "claims" WHERE vehicle_id = \$1 AND "claims"."deleted_at" IS NULL`).
					WithArgs(vehicleID).
					WillReturnRows(countRows)

				claimID := uuid.New()
				customerID := uuid.New()
				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "customer_id", "description", "status",
					"total_cost", "approved_by", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					claimID, vehicleID, customerID, "Claim", entities.ClaimStatusDraft,
					1000.0, nil, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "claims" WHERE vehicle_id = \$1 AND "claims"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT \$2`).
					WithArgs(vehicleID, 10).
					WillReturnRows(rows)

				claims, total, err := repository.FindAll(ctx, filters, pagination)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(HaveLen(1))
				Expect(total).To(Equal(int64(1)))
				Expect(claims[0].VehicleID).To(Equal(vehicleID))
			})
		})

		Context("when filtering by status", func() {
			It("should return filtered claims", func() {
				status := entities.ClaimStatusApproved
				filters.Status = &status

				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(`SELECT count\(\*\) FROM "claims" WHERE status = \$1 AND "claims"."deleted_at" IS NULL`).
					WithArgs(status).
					WillReturnRows(countRows)

				claimID := uuid.New()
				vehicleID := uuid.New()
				customerID := uuid.New()
				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "customer_id", "description", "status",
					"total_cost", "approved_by", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					claimID, vehicleID, customerID, "Claim", status,
					1000.0, nil, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "claims" WHERE status = \$1 AND "claims"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT \$2`).
					WithArgs(status, 10).
					WillReturnRows(rows)

				claims, total, err := repository.FindAll(ctx, filters, pagination)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(HaveLen(1))
				Expect(total).To(Equal(int64(1)))
				Expect(claims[0].Status).To(Equal(status))
			})
		})

		Context("when filtering by date range", func() {
			It("should return claims within date range", func() {
				fromDate := time.Now().AddDate(0, 0, -7)
				toDate := time.Now()
				filters.FromDate = &fromDate
				filters.ToDate = &toDate

				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(`SELECT count\(\*\) FROM "claims" WHERE created_at >= \$1 AND created_at <= \$2 AND "claims"."deleted_at" IS NULL`).
					WithArgs(fromDate, toDate).
					WillReturnRows(countRows)

				claimID := uuid.New()
				vehicleID := uuid.New()
				customerID := uuid.New()
				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "customer_id", "description", "status",
					"total_cost", "approved_by", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					claimID, vehicleID, customerID, "Claim", entities.ClaimStatusDraft,
					1000.0, nil, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "claims" WHERE created_at >= \$1 AND created_at <= \$2 AND "claims"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT \$3`).
					WithArgs(fromDate, toDate, 10).
					WillReturnRows(rows)

				claims, total, err := repository.FindAll(ctx, filters, pagination)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(HaveLen(1))
				Expect(total).To(Equal(int64(1)))
			})
		})

		Context("boundary cases for pagination", func() {
			It("should handle page size of 0 (no pagination)", func() {
				pagination.PageSize = 0

				countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(`SELECT count\(\*\) FROM "claims" WHERE "claims"."deleted_at" IS NULL`).
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "customer_id", "description", "status",
					"total_cost", "approved_by", "created_at", "updated_at", "deleted_at",
				})

				mock.ExpectQuery(`SELECT \* FROM "claims" WHERE "claims"."deleted_at" IS NULL ORDER BY created_at DESC`).
					WillReturnRows(rows)

				claims, total, err := repository.FindAll(ctx, filters, pagination)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(BeEmpty())
				Expect(total).To(Equal(int64(0)))
			})

			It("should handle large page numbers", func() {
				pagination.Page = 100
				offset := (100 - 1) * 10

				countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(`SELECT count\(\*\) FROM "claims" WHERE "claims"."deleted_at" IS NULL`).
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{
					"id", "vehicle_id", "customer_id", "description", "status",
					"total_cost", "approved_by", "created_at", "updated_at", "deleted_at",
				})

				mock.ExpectQuery(`SELECT \* FROM "claims" WHERE "claims"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
					WithArgs(10, offset).
					WillReturnRows(rows)

				claims, total, err := repository.FindAll(ctx, filters, pagination)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(BeEmpty())
				Expect(total).To(Equal(int64(0)))
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError on count", func() {
				mock.ExpectQuery(`SELECT count\(\*\) FROM "claims" WHERE "claims"."deleted_at" IS NULL`).
					WillReturnError(errors.New("database connection failed"))

				claims, total, err := repository.FindAll(ctx, filters, pagination)

				Expect(err).To(HaveOccurred())
				Expect(claims).To(BeNil())
				Expect(total).To(Equal(int64(0)))
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})

			It("should return DBOperationError on find", func() {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(`SELECT count\(\*\) FROM "claims" WHERE "claims"."deleted_at" IS NULL`).
					WillReturnRows(countRows)

				mock.ExpectQuery(`SELECT \* FROM "claims" WHERE "claims"."deleted_at" IS NULL`).
					WillReturnError(errors.New("database connection failed"))

				claims, total, err := repository.FindAll(ctx, filters, pagination)

				Expect(err).To(HaveOccurred())
				Expect(claims).To(BeNil())
				Expect(total).To(Equal(int64(0)))
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
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
					claimID1, vehicleID, customerID, "Claim 1", entities.ClaimStatusDraft,
					1000.0, nil, time.Now(), time.Now(), nil,
				).AddRow(
					claimID2, vehicleID, customerID, "Claim 2", entities.ClaimStatusSubmitted,
					2000.0, nil, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "claims" WHERE customer_id = \$1 AND "claims"."deleted_at" IS NULL ORDER BY created_at DESC`).
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

				mock.ExpectQuery(`SELECT \* FROM "claims" WHERE customer_id = \$1 AND "claims"."deleted_at" IS NULL ORDER BY created_at DESC`).
					WithArgs(customerID).
					WillReturnRows(rows)

				claims, err := repository.FindByCustomerID(ctx, customerID)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(`SELECT \* FROM "claims" WHERE customer_id = \$1 AND "claims"."deleted_at" IS NULL`).
					WithArgs(customerID).
					WillReturnError(errors.New("database connection failed"))

				claims, err := repository.FindByCustomerID(ctx, customerID)

				Expect(err).To(HaveOccurred())
				Expect(claims).To(BeNil())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
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
					claimID1, vehicleID, customerID, "Claim 1", entities.ClaimStatusDraft,
					1000.0, nil, time.Now(), time.Now(), nil,
				).AddRow(
					claimID2, vehicleID, customerID, "Claim 2", entities.ClaimStatusSubmitted,
					2000.0, nil, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "claims" WHERE vehicle_id = \$1 AND "claims"."deleted_at" IS NULL ORDER BY created_at DESC`).
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

				mock.ExpectQuery(`SELECT \* FROM "claims" WHERE vehicle_id = \$1 AND "claims"."deleted_at" IS NULL ORDER BY created_at DESC`).
					WithArgs(vehicleID).
					WillReturnRows(rows)

				claims, err := repository.FindByVehicleID(ctx, vehicleID)

				Expect(err).NotTo(HaveOccurred())
				Expect(claims).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(`SELECT \* FROM "claims" WHERE vehicle_id = \$1 AND "claims"."deleted_at" IS NULL`).
					WithArgs(vehicleID).
					WillReturnError(errors.New("database connection failed"))

				claims, err := repository.FindByVehicleID(ctx, vehicleID)

				Expect(err).To(HaveOccurred())
				Expect(claims).To(BeNil())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})
	})
})

func newClaim() *entities.Claim {
	return &entities.Claim{
		ID:          uuid.New(),
		VehicleID:   uuid.New(),
		CustomerID:  uuid.New(),
		Description: "Test claim description",
		Status:      entities.ClaimStatusDraft,
		TotalCost:   1000.0,
		ApprovedBy:  nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
}
