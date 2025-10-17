package persistence_test

import (
	"context"
	"errors"
	"ev-warranty-go/pkg/mocks"
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

var _ = Describe("ClaimItemRepository", func() {
	var (
		mock       sqlmock.Sqlmock
		db         *gorm.DB
		repository repositories.ClaimItemRepository
		ctx        context.Context
	)

	BeforeEach(func() {
		mock, db = SetupMockDB()
		repository = persistence.NewClaimItemRepository(db)
		ctx = context.Background()
	})

	AfterEach(func() {
		CleanupMockDB(mock)
	})

	Describe("Create", func() {
		var item *entities.ClaimItem

		BeforeEach(func() {
			item = newClaimItem()
		})

		Context("when claim item is created successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				MockSuccessfulInsert(mock, "claim_items", item.ID)

				err := repository.Create(mockTx, item)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a duplicate key constraint", func() {
			It("should return DBDuplicateKeyError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				MockDuplicateKeyError(mock, "claim_items", "claim_items_unique_key")

				err := repository.Create(mockTx, item)

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
				MockInsertError(mock, "claim_items")

				err := repository.Create(mockTx, item)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})

		Context("boundary cases for item types", func() {
			It("should handle replacement type", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				item.Type = entities.ClaimItemTypeReplacement
				MockSuccessfulInsert(mock, "claim_items", item.ID)

				err := repository.Create(mockTx, item)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should handle repair type", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				item.Type = entities.ClaimItemTypeRepair
				MockSuccessfulInsert(mock, "claim_items", item.ID)

				err := repository.Create(mockTx, item)

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("Update", func() {
		var item *entities.ClaimItem

		BeforeEach(func() {
			item = newClaimItem()
			item.IssueDescription = "Updated issue description"
			item.Cost = 2000.0
		})

		Context("when claim item is updated successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				MockSuccessfulUpdate(mock, "claim_items")

				err := repository.Update(mockTx, item)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				MockUpdateError(mock, "claim_items")

				err := repository.Update(mockTx, item)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})

		Context("boundary cases for cost", func() {
			It("should handle zero cost", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				item.Cost = 0.0
				MockSuccessfulUpdate(mock, "claim_items")

				err := repository.Update(mockTx, item)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should handle negative cost", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				item.Cost = -100.0
				MockSuccessfulUpdate(mock, "claim_items")

				err := repository.Update(mockTx, item)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should handle very large cost", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				item.Cost = 999999999.99
				MockSuccessfulUpdate(mock, "claim_items")

				err := repository.Update(mockTx, item)

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("HardDelete", func() {
		var itemID uuid.UUID

		BeforeEach(func() {
			itemID = uuid.New()
		})

		Context("when claim item is hard deleted successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "claim_items" WHERE id = \$1`).
					WithArgs(itemID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := repository.HardDelete(mockTx, itemID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "claim_items" WHERE id = \$1`).
					WithArgs(itemID).
					WillReturnError(errors.New("database connection failed"))
				mock.ExpectRollback()

				err := repository.HardDelete(mockTx, itemID)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})
	})

	Describe("SoftDeleteByClaimID", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when claim items are soft deleted successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "claim_items" SET "deleted_at"=\$1 WHERE claim_id = \$2`).
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
				mockTx.On("GetTx").Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "claim_items" SET "deleted_at"=\$1 WHERE claim_id = \$2`).
					WithArgs(sqlmock.AnyArg(), claimID).
					WillReturnError(errors.New("database connection failed"))
				mock.ExpectRollback()

				err := repository.SoftDeleteByClaimID(mockTx, claimID)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})
	})

	Describe("UpdateStatus", func() {
		var itemID uuid.UUID
		var status string

		BeforeEach(func() {
			itemID = uuid.New()
			status = entities.ClaimItemStatusApproved
		})

		Context("when status is updated successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "claim_items" SET "line_status"=\$1,"updated_at"=\$2 WHERE id = \$3 AND "claim_items"."deleted_at" IS NULL`).
					WithArgs(status, sqlmock.AnyArg(), itemID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := repository.UpdateStatus(mockTx, itemID, status)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "claim_items" SET "line_status"=\$1,"updated_at"=\$2 WHERE id = \$3 AND "claim_items"."deleted_at" IS NULL`).
					WithArgs(status, sqlmock.AnyArg(), itemID).
					WillReturnError(errors.New("database connection failed"))
				mock.ExpectRollback()

				err := repository.UpdateStatus(mockTx, itemID, status)

				Expect(err).To(HaveOccurred())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})

		Context("boundary cases for status", func() {
			It("should handle all valid status values", func() {
				statuses := []string{
					entities.ClaimItemStatusPending,
					entities.ClaimItemStatusApproved,
					entities.ClaimItemStatusRejected,
				}

				for _, s := range statuses {
					mockTx := mocks.NewTx(GinkgoT())
					mockTx.On("GetTx").Return(db)
					mock.ExpectBegin()
					mock.ExpectExec(`UPDATE "claim_items" SET "line_status"=\$1,"updated_at"=\$2 WHERE id = \$3 AND "claim_items"."deleted_at" IS NULL`).
						WithArgs(s, sqlmock.AnyArg(), itemID).
						WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()

					err := repository.UpdateStatus(mockTx, itemID, s)
					Expect(err).NotTo(HaveOccurred())
				}
			})
		})
	})

	Describe("SumCostByClaimID", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when sum is calculated successfully", func() {
			It("should return total cost", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				rows := sqlmock.NewRows([]string{"sum"}).AddRow(5000.0)

				mock.ExpectQuery(`SELECT COALESCE\(SUM\(cost\), 0\) FROM "claim_items" WHERE \(claim_id = \$1 AND status = 'APPROVED'\) AND "claim_items"."deleted_at" IS NULL`).
					WithArgs(claimID).
					WillReturnRows(rows)

				totalCost, err := repository.SumCostByClaimID(mockTx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(totalCost).To(Equal(5000.0))
			})
		})

		Context("when no approved items exist", func() {
			It("should return zero", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				rows := sqlmock.NewRows([]string{"sum"}).AddRow(0.0)

				mock.ExpectQuery(`SELECT COALESCE\(SUM\(cost\), 0\) FROM "claim_items" WHERE \(claim_id = \$1 AND status = 'APPROVED'\) AND "claim_items"."deleted_at" IS NULL`).
					WithArgs(claimID).
					WillReturnRows(rows)

				totalCost, err := repository.SumCostByClaimID(mockTx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(totalCost).To(Equal(0.0))
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				mock.ExpectQuery(`SELECT COALESCE\(SUM\(cost\), 0\) FROM "claim_items" WHERE \(claim_id = \$1 AND status = 'APPROVED'\) AND "claim_items"."deleted_at" IS NULL`).
					WithArgs(claimID).
					WillReturnError(errors.New("database connection failed"))

				totalCost, err := repository.SumCostByClaimID(mockTx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(totalCost).To(Equal(0.0))
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})

		Context("boundary cases for sum", func() {
			It("should handle very large sums", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				rows := sqlmock.NewRows([]string{"sum"}).AddRow(999999999.99)

				mock.ExpectQuery(`SELECT COALESCE\(SUM\(cost\), 0\) FROM "claim_items" WHERE \(claim_id = \$1 AND status = 'APPROVED'\) AND "claim_items"."deleted_at" IS NULL`).
					WithArgs(claimID).
					WillReturnRows(rows)

				totalCost, err := repository.SumCostByClaimID(mockTx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(totalCost).To(Equal(999999999.99))
			})
		})
	})

	Describe("FindByID", func() {
		var itemID uuid.UUID

		BeforeEach(func() {
			itemID = uuid.New()
		})

		Context("when claim item is found", func() {
			It("should return the claim item", func() {
				expected := newClaimItem()
				expected.ID = itemID
				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "part_category_id", "faulty_part_id", "replacement_part_id",
					"issue_description", "status", "type", "cost", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					expected.ID, expected.ClaimID, expected.PartCategoryID, expected.FaultyPartID,
					expected.ReplacementPartID, expected.IssueDescription, expected.Status,
					expected.Type, expected.Cost, expected.CreatedAt, expected.UpdatedAt, expected.DeletedAt,
				)

				MockFindByID(mock, "claim_items", itemID, rows)

				item, err := repository.FindByID(ctx, itemID)

				Expect(err).NotTo(HaveOccurred())
				Expect(item).NotTo(BeNil())
				Expect(item.ID).To(Equal(expected.ID))
				Expect(item.ClaimID).To(Equal(expected.ClaimID))
				Expect(item.IssueDescription).To(Equal(expected.IssueDescription))
				Expect(item.Cost).To(Equal(expected.Cost))
			})
		})

		Context("when claim item is not found", func() {
			It("should return ClaimItemNotFound error", func() {
				MockNotFound(mock, "claim_items", itemID)

				item, err := repository.FindByID(ctx, itemID)

				Expect(err).To(HaveOccurred())
				Expect(item).To(BeNil())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeClaimItemNotFound))
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockQueryError(mock, `SELECT * FROM "claim_items" WHERE id = $1`)

				item, err := repository.FindByID(ctx, itemID)

				Expect(err).To(HaveOccurred())
				Expect(item).To(BeNil())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})
	})

	Describe("FindByClaimID", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when claim items are found", func() {
			It("should return all items for claim", func() {
				itemID1 := uuid.New()
				itemID2 := uuid.New()
				partID := uuid.New()
				replacementPartID := uuid.New()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "part_category_id", "faulty_part_id", "replacement_part_id",
					"issue_description", "status", "type", "cost", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					itemID1, claimID, 1, partID, &replacementPartID, "Issue 1",
					entities.ClaimItemStatusPending, entities.ClaimItemTypeReplacement,
					1000.0, time.Now(), time.Now(), nil,
				).AddRow(
					itemID2, claimID, 2, partID, nil, "Issue 2",
					entities.ClaimItemStatusApproved, entities.ClaimItemTypeRepair,
					2000.0, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "claim_items" WHERE claim_id = \$1 AND "claim_items"."deleted_at" IS NULL ORDER BY created_at ASC`).
					WithArgs(claimID).
					WillReturnRows(rows)

				items, err := repository.FindByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(items).To(HaveLen(2))
				Expect(items[0].ClaimID).To(Equal(claimID))
				Expect(items[1].ClaimID).To(Equal(claimID))
			})
		})

		Context("when no claim items are found", func() {
			It("should return empty slice", func() {
				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "part_category_id", "faulty_part_id", "replacement_part_id",
					"issue_description", "status", "type", "cost", "created_at", "updated_at", "deleted_at",
				})

				mock.ExpectQuery(`SELECT \* FROM "claim_items" WHERE claim_id = \$1 AND "claim_items"."deleted_at" IS NULL ORDER BY created_at ASC`).
					WithArgs(claimID).
					WillReturnRows(rows)

				items, err := repository.FindByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(items).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(`SELECT \* FROM "claim_items" WHERE claim_id = \$1 AND "claim_items"."deleted_at" IS NULL`).
					WithArgs(claimID).
					WillReturnError(errors.New("database connection failed"))

				items, err := repository.FindByClaimID(ctx, claimID)

				Expect(err).To(HaveOccurred())
				Expect(items).To(BeNil())
				var appErr *apperrors.AppError
				Expect(errors.As(err, &appErr)).To(BeTrue())
				Expect(appErr.ErrorCode).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})

		Context("boundary cases", func() {
			It("should handle items without replacement parts", func() {
				itemID := uuid.New()
				partID := uuid.New()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "part_category_id", "faulty_part_id", "replacement_part_id",
					"issue_description", "status", "type", "cost", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					itemID, claimID, 1, partID, nil, "Repair only",
					entities.ClaimItemStatusPending, entities.ClaimItemTypeRepair,
					500.0, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(`SELECT \* FROM "claim_items" WHERE claim_id = \$1 AND "claim_items"."deleted_at" IS NULL ORDER BY created_at ASC`).
					WithArgs(claimID).
					WillReturnRows(rows)

				items, err := repository.FindByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(items).To(HaveLen(1))
				Expect(items[0].ReplacementPartID).To(BeNil())
			})
		})
	})
})

func newClaimItem() *entities.ClaimItem {
	replacementPartID := uuid.New()
	return &entities.ClaimItem{
		ID:                uuid.New(),
		ClaimID:           uuid.New(),
		PartCategoryID:    1,
		FaultyPartID:      uuid.New(),
		ReplacementPartID: &replacementPartID,
		IssueDescription:  "Test issue description",
		Status:            entities.ClaimItemStatusPending,
		Type:              entities.ClaimItemTypeReplacement,
		Cost:              1000.0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		DeletedAt:         nil,
	}
}
