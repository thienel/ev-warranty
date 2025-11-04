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

var _ = Describe("ClaimItemRepository", func() {
	var (
		mock       sqlmock.Sqlmock
		db         *gorm.DB
		repository repository.ClaimItemRepository
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
		var item *entity.ClaimItem

		BeforeEach(func() {
			item = newClaimItem()
		})

		Context("when claim item is created successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockSuccessfulInsert(mock, "claim_items", item.ID)

				err := repository.Create(mockTx, item)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a duplicate key constraint", func() {
			It("should return DBDuplicateKeyError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockDuplicateKeyError(mock, "claim_items", "claim_items_unique_key")

				err := repository.Create(mockTx, item)

				ExpectAppError(err, apperror.ErrDuplicateKey.ErrorCode)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockInsertError(mock, "claim_items")

				err := repository.Create(mockTx, item)

				ExpectAppError(err, apperror.ErrDBOperation.ErrorCode)
			})
		})

		Context("boundary cases for item types", func() {
			It("should handle replacement type", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				item.Type = entity.ClaimItemTypeReplacement
				MockSuccessfulInsert(mock, "claim_items", item.ID)

				err := repository.Create(mockTx, item)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should handle repair type", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				item.Type = entity.ClaimItemTypeRepair
				MockSuccessfulInsert(mock, "claim_items", item.ID)

				err := repository.Create(mockTx, item)

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("Update", func() {
		var item *entity.ClaimItem

		BeforeEach(func() {
			item = newClaimItem()
			item.IssueDescription = "Updated issue description"
			item.Cost = 2000.0
		})

		Context("when claim item is updated successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockSuccessfulUpdate(mock, "claim_items")

				err := repository.Update(mockTx, item)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockUpdateError(mock, "claim_items")

				err := repository.Update(mockTx, item)

				ExpectAppError(err, apperror.ErrDBOperation.ErrorCode)
			})
		})

		Context("boundary cases for cost", func() {
			It("should handle zero cost", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				item.Cost = 0.0
				MockSuccessfulUpdate(mock, "claim_items")

				err := repository.Update(mockTx, item)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should handle negative cost", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
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
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "claim_items" WHERE id = $1`)).
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
				mockTx.EXPECT().GetTx().Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "claim_items" WHERE id = $1`)).
					WithArgs(itemID).
					WillReturnError(errors.New("database connection failed"))
				mock.ExpectRollback()

				err := repository.HardDelete(mockTx, itemID)

				ExpectAppError(err, apperror.ErrDBOperation.ErrorCode)
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
				mockTx.EXPECT().GetTx().Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "claim_items" SET "deleted_at"=$1 WHERE claim_id = $2`)).
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
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "claim_items" SET "deleted_at"=$1 WHERE claim_id = $2`)).
					WithArgs(sqlmock.AnyArg(), claimID).
					WillReturnError(errors.New("database connection failed"))
				mock.ExpectRollback()

				err := repository.SoftDeleteByClaimID(mockTx, claimID)

				ExpectAppError(err, apperror.ErrDBOperation.ErrorCode)
			})
		})
	})

	Describe("UpdateStatus", func() {
		var itemID uuid.UUID
		var status string

		BeforeEach(func() {
			itemID = uuid.New()
			status = entity.ClaimItemStatusApproved
		})

		Context("when status is updated successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "claim_items" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND "claim_items"."deleted_at" IS NULL`)).
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
				mockTx.EXPECT().GetTx().Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "claim_items" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND "claim_items"."deleted_at" IS NULL`)).
					WithArgs(status, sqlmock.AnyArg(), itemID).
					WillReturnError(errors.New("database connection failed"))
				mock.ExpectRollback()

				err := repository.UpdateStatus(mockTx, itemID, status)

				ExpectAppError(err, apperror.ErrDBOperation.ErrorCode)
			})
		})

		Context("boundary cases for status", func() {
			It("should handle all valid status values", func() {
				statuses := []string{
					entity.ClaimItemStatusPending,
					entity.ClaimItemStatusApproved,
					entity.ClaimItemStatusRejected,
				}

				for _, s := range statuses {
					mockTx := mocks.NewTx(GinkgoT())
					mockTx.EXPECT().GetTx().Return(db)
					mock.ExpectBegin()
					mock.ExpectExec(regexp.QuoteMeta(`UPDATE "claim_items" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND "claim_items"."deleted_at" IS NULL`)).
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
				mockTx.EXPECT().GetTx().Return(db)
				rows := sqlmock.NewRows([]string{"sum"}).AddRow(5000.0)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(SUM(cost), 0) FROM "claim_items" WHERE (claim_id = $1 AND status = 'APPROVED') AND "claim_items"."deleted_at" IS NULL`)).
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
				mockTx.EXPECT().GetTx().Return(db)
				rows := sqlmock.NewRows([]string{"sum"}).AddRow(0.0)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(SUM(cost), 0) FROM "claim_items" WHERE (claim_id = $1 AND status = 'APPROVED') AND "claim_items"."deleted_at" IS NULL`)).
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
				mockTx.EXPECT().GetTx().Return(db)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(SUM(cost), 0) FROM "claim_items" WHERE (claim_id = $1 AND status = 'APPROVED') AND "claim_items"."deleted_at" IS NULL`)).
					WithArgs(claimID).
					WillReturnError(errors.New("database connection failed"))

				totalCost, err := repository.SumCostByClaimID(mockTx, claimID)

				Expect(totalCost).To(Equal(0.0))
				ExpectAppError(err, apperror.ErrDBOperation.ErrorCode)
			})
		})

		Context("boundary cases for sum", func() {
			It("should handle very large sums", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.On("GetTx").Return(db)
				rows := sqlmock.NewRows([]string{"sum"}).AddRow(999999999.99)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT COALESCE(SUM(cost), 0) FROM "claim_items" WHERE (claim_id = $1 AND status = 'APPROVED') AND "claim_items"."deleted_at" IS NULL`)).
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

				Expect(item).To(BeNil())
				ExpectAppError(err, apperror.ErrNotFoundError.ErrorCode)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockQueryError(mock, `SELECT * FROM "claim_items" WHERE id = $1`)

				item, err := repository.FindByID(ctx, itemID)

				Expect(item).To(BeNil())
				ExpectAppError(err, apperror.ErrDBOperation.ErrorCode)
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
				cateID := uuid.New()
				replacementPartID := uuid.New()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "part_category_id", "faulty_part_id", "replacement_part_id",
					"issue_description", "status", "type", "cost", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					itemID1, claimID, cateID, partID, &replacementPartID, "Issue 1",
					entity.ClaimItemStatusPending, entity.ClaimItemTypeReplacement,
					1000.0, time.Now(), time.Now(), nil,
				).AddRow(
					itemID2, claimID, cateID, partID, nil, "Issue 2",
					entity.ClaimItemStatusApproved, entity.ClaimItemTypeRepair,
					2000.0, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_items" WHERE claim_id = $1 AND "claim_items"."deleted_at" IS NULL ORDER BY created_at ASC`)).
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

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_items" WHERE claim_id = $1 AND "claim_items"."deleted_at" IS NULL ORDER BY created_at ASC`)).
					WithArgs(claimID).
					WillReturnRows(rows)

				items, err := repository.FindByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(items).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_items" WHERE claim_id = $1 AND "claim_items"."deleted_at" IS NULL`)).
					WithArgs(claimID).
					WillReturnError(errors.New("database connection failed"))

				items, err := repository.FindByClaimID(ctx, claimID)

				Expect(items).To(BeNil())
				ExpectAppError(err, apperror.ErrDBOperation.ErrorCode)
			})
		})

		Context("boundary cases", func() {
			It("should handle items without replacement parts", func() {
				itemID := uuid.New()
				partID := uuid.New()
				cateID := uuid.New()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "part_category_id", "faulty_part_id", "replacement_part_id",
					"issue_description", "status", "type", "cost", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					itemID, claimID, cateID, partID, nil, "Repair only",
					entity.ClaimItemStatusPending, entity.ClaimItemTypeRepair,
					500.0, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_items" WHERE claim_id = $1 AND "claim_items"."deleted_at" IS NULL ORDER BY created_at ASC`)).
					WithArgs(claimID).
					WillReturnRows(rows)

				items, err := repository.FindByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(items).To(HaveLen(1))
				Expect(items[0].ReplacementPartID).To(BeNil())
			})
		})
	})

	Describe("CountByClaimID", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when count is successful", func() {
			It("should return the correct count", func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "claim_items" WHERE claim_id = $1 AND "claim_items"."deleted_at" IS NULL`)).
					WithArgs(claimID).
					WillReturnRows(rows)

				count, err := repository.CountByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(int64(5)))
			})
		})

		Context("when no items exist", func() {
			It("should return zero", func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "claim_items" WHERE claim_id = $1 AND "claim_items"."deleted_at" IS NULL`)).
					WithArgs(claimID).
					WillReturnRows(rows)

				count, err := repository.CountByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(count).To(Equal(int64(0)))
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "claim_items" WHERE claim_id = $1 AND "claim_items"."deleted_at" IS NULL`)).
					WithArgs(claimID).
					WillReturnError(errors.New("database connection failed"))

				count, err := repository.CountByClaimID(ctx, claimID)

				Expect(count).To(Equal(int64(0)))
				ExpectAppError(err, apperror.ErrDBOperation.ErrorCode)
			})
		})
	})

	Describe("FindByStatus", func() {
		var claimID uuid.UUID
		var status string

		BeforeEach(func() {
			claimID = uuid.New()
			status = entity.ClaimItemStatusApproved
		})

		Context("when items are found with specified status", func() {
			It("should return all items with that status", func() {
				itemID1 := uuid.New()
				itemID2 := uuid.New()
				partID := uuid.New()
				cateID := uuid.New()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "part_category_id", "faulty_part_id", "replacement_part_id",
					"issue_description", "status", "type", "cost", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					itemID1, claimID, cateID, partID, nil, "Issue 1",
					entity.ClaimItemStatusApproved, entity.ClaimItemTypeReplacement,
					1000.0, time.Now(), time.Now(), nil,
				).AddRow(
					itemID2, claimID, cateID, partID, nil, "Issue 2",
					entity.ClaimItemStatusApproved, entity.ClaimItemTypeRepair,
					2000.0, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_items" WHERE (claim_id = $1 AND status = $2) AND "claim_items"."deleted_at" IS NULL ORDER BY created_at ASC`)).
					WithArgs(claimID, status).
					WillReturnRows(rows)

				items, err := repository.FindByStatus(ctx, claimID, status)

				Expect(err).NotTo(HaveOccurred())
				Expect(items).To(HaveLen(2))
				Expect(items[0].Status).To(Equal(entity.ClaimItemStatusApproved))
				Expect(items[1].Status).To(Equal(entity.ClaimItemStatusApproved))
			})
		})

		Context("when no items match the status", func() {
			It("should return empty slice", func() {
				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "part_category_id", "faulty_part_id", "replacement_part_id",
					"issue_description", "status", "type", "cost", "created_at", "updated_at", "deleted_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_items" WHERE (claim_id = $1 AND status = $2) AND "claim_items"."deleted_at" IS NULL ORDER BY created_at ASC`)).
					WithArgs(claimID, status).
					WillReturnRows(rows)

				items, err := repository.FindByStatus(ctx, claimID, status)

				Expect(err).NotTo(HaveOccurred())
				Expect(items).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_items" WHERE (claim_id = $1 AND status = $2) AND "claim_items"."deleted_at" IS NULL ORDER BY created_at ASC`)).
					WithArgs(claimID, status).
					WillReturnError(errors.New("database connection failed"))

				items, err := repository.FindByStatus(ctx, claimID, status)

				Expect(items).To(BeNil())
				ExpectAppError(err, apperror.ErrDBOperation.ErrorCode)
			})
		})

		Context("boundary cases for different statuses", func() {
			It("should handle pending status", func() {
				itemID := uuid.New()
				partID := uuid.New()
				cateID := uuid.New()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "part_category_id", "faulty_part_id", "replacement_part_id",
					"issue_description", "status", "type", "cost", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					itemID, claimID, cateID, partID, nil, "Pending item",
					entity.ClaimItemStatusPending, entity.ClaimItemTypeRepair,
					500.0, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_items" WHERE (claim_id = $1 AND status = $2) AND "claim_items"."deleted_at" IS NULL ORDER BY created_at ASC`)).
					WithArgs(claimID, entity.ClaimItemStatusPending).
					WillReturnRows(rows)

				items, err := repository.FindByStatus(ctx, claimID, entity.ClaimItemStatusPending)

				Expect(err).NotTo(HaveOccurred())
				Expect(items).To(HaveLen(1))
				Expect(items[0].Status).To(Equal(entity.ClaimItemStatusPending))
			})

			It("should handle rejected status", func() {
				itemID := uuid.New()
				partID := uuid.New()
				cateID := uuid.New()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "part_category_id", "faulty_part_id", "replacement_part_id",
					"issue_description", "status", "type", "cost", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					itemID, claimID, cateID, partID, nil, "Rejected item",
					entity.ClaimItemStatusRejected, entity.ClaimItemTypeRepair,
					500.0, time.Now(), time.Now(), nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_items" WHERE (claim_id = $1 AND status = $2) AND "claim_items"."deleted_at" IS NULL ORDER BY created_at ASC`)).
					WithArgs(claimID, entity.ClaimItemStatusRejected).
					WillReturnRows(rows)

				items, err := repository.FindByStatus(ctx, claimID, entity.ClaimItemStatusRejected)

				Expect(err).NotTo(HaveOccurred())
				Expect(items).To(HaveLen(1))
				Expect(items[0].Status).To(Equal(entity.ClaimItemStatusRejected))
			})
		})
	})
})

func newClaimItem() *entity.ClaimItem {
	replacementPartID := uuid.New()
	return &entity.ClaimItem{
		ID:                uuid.New(),
		ClaimID:           uuid.New(),
		PartCategoryID:    uuid.New(),
		FaultyPartID:      uuid.New(),
		ReplacementPartID: &replacementPartID,
		IssueDescription:  "Test issue description",
		Status:            entity.ClaimItemStatusPending,
		Type:              entity.ClaimItemTypeReplacement,
		Cost:              1000.0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		DeletedAt:         nil,
	}
}
