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

var _ = Describe("ClaimHistoryRepository", func() {
	var (
		mock       sqlmock.Sqlmock
		db         *gorm.DB
		repository repositories.ClaimHistoryRepository
		ctx        context.Context
	)

	BeforeEach(func() {
		mock, db = SetupMockDB()
		repository = persistence.NewClaimHistoryRepository(db)
		ctx = context.Background()
	})

	AfterEach(func() {
		CleanupMockDB(mock)
	})

	Describe("Create", func() {
		var history *entities.ClaimHistory

		BeforeEach(func() {
			history = newClaimHistory()
		})

		Context("when claim history is created successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockSuccessfulInsert(mock, "claim_histories", history.ID)

				err := repository.Create(mockTx, history)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a duplicate key constraint", func() {
			It("should return DBDuplicateKeyError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockDuplicateKeyError(mock, "claim_histories", "claim_histories_unique_key")

				err := repository.Create(mockTx, history)

				ExpectAppError(err, apperrors.ErrorCodeDuplicateKey)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				MockInsertError(mock, "claim_histories")

				err := repository.Create(mockTx, history)

				ExpectAppError(err, apperrors.ErrorCodeDBOperation)
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
					history.Status = s
					MockSuccessfulInsert(mock, "claim_histories", history.ID)

					err := repository.Create(mockTx, history)
					Expect(err).NotTo(HaveOccurred())
				}
			})
		})
	})

	Describe("SoftDeleteByClaimID", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when claim histories are soft deleted successfully", func() {
			It("should return nil error", func() {
				mockTx := mocks.NewTx(GinkgoT())
				mockTx.EXPECT().GetTx().Return(db)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "claim_histories" SET "deleted_at"=$1 WHERE claim_id = $2`)).
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
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "claim_histories" SET "deleted_at"=$1 WHERE claim_id = $2`)).
					WithArgs(sqlmock.AnyArg(), claimID).
					WillReturnError(errors.New("database connection failed"))
				mock.ExpectRollback()

				err := repository.SoftDeleteByClaimID(mockTx, claimID)

				ExpectAppError(err, apperrors.ErrorCodeDBOperation)
			})
		})
	})

	Describe("FindByClaimID", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when claim histories are found", func() {
			It("should return all histories for claim", func() {
				historyID1 := uuid.New()
				historyID2 := uuid.New()
				changedBy := uuid.New()
				now := time.Now()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "status", "changed_by", "changed_at", "deleted_at",
				}).AddRow(
					historyID1, claimID, entities.ClaimStatusDraft, changedBy, now, nil,
				).AddRow(
					historyID2, claimID, entities.ClaimStatusSubmitted, changedBy, now.Add(-1*time.Hour), nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE claim_id = $1 AND "claim_histories"."deleted_at" IS NULL ORDER BY changed_at DESC`)).
					WithArgs(claimID).
					WillReturnRows(rows)

				histories, err := repository.FindByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(histories).To(HaveLen(2))
				Expect(histories[0].ClaimID).To(Equal(claimID))
				Expect(histories[1].ClaimID).To(Equal(claimID))
			})
		})

		Context("when no claim histories are found", func() {
			It("should return empty slice", func() {
				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "status", "changed_by", "changed_at", "deleted_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE claim_id = $1 AND "claim_histories"."deleted_at" IS NULL ORDER BY changed_at DESC`)).
					WithArgs(claimID).
					WillReturnRows(rows)

				histories, err := repository.FindByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(histories).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE claim_id = $1 AND "claim_histories"."deleted_at" IS NULL`)).
					WithArgs(claimID).
					WillReturnError(errors.New("database connection failed"))

				histories, err := repository.FindByClaimID(ctx, claimID)

				Expect(histories).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeDBOperation)
			})
		})

		Context("boundary cases", func() {
			It("should handle single history entry", func() {
				historyID := uuid.New()
				changedBy := uuid.New()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "status", "changed_by", "changed_at", "deleted_at",
				}).AddRow(
					historyID, claimID, entities.ClaimStatusDraft, changedBy, time.Now(), nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE claim_id = $1 AND "claim_histories"."deleted_at" IS NULL ORDER BY changed_at DESC`)).
					WithArgs(claimID).
					WillReturnRows(rows)

				histories, err := repository.FindByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(histories).To(HaveLen(1))
			})

			It("should handle many history entries", func() {
				changedBy := uuid.New()
				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "status", "changed_by", "changed_at", "deleted_at",
				})

				for i := 0; i < 100; i++ {
					rows.AddRow(
						uuid.New(), claimID, entities.ClaimStatusDraft, changedBy, time.Now(), nil,
					)
				}

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE claim_id = $1 AND "claim_histories"."deleted_at" IS NULL ORDER BY changed_at DESC`)).
					WithArgs(claimID).
					WillReturnRows(rows)

				histories, err := repository.FindByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(histories).To(HaveLen(100))
			})
		})
	})

	Describe("FindLatestByClaimID", func() {
		var claimID uuid.UUID

		BeforeEach(func() {
			claimID = uuid.New()
		})

		Context("when latest history is found", func() {
			It("should return the most recent history", func() {
				historyID := uuid.New()
				changedBy := uuid.New()
				now := time.Now()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "status", "changed_by", "changed_at", "deleted_at",
				}).AddRow(
					historyID, claimID, entities.ClaimStatusApproved, changedBy, now, nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE claim_id = $1 AND "claim_histories"."deleted_at" IS NULL ORDER BY changed_at DESC,"claim_histories"."id" LIMIT $2`)).
					WithArgs(claimID, 1).
					WillReturnRows(rows)

				history, err := repository.FindLatestByClaimID(ctx, claimID)

				Expect(err).NotTo(HaveOccurred())
				Expect(history).NotTo(BeNil())
				Expect(history.ID).To(Equal(historyID))
				Expect(history.ClaimID).To(Equal(claimID))
				Expect(history.Status).To(Equal(entities.ClaimStatusApproved))
			})
		})

		Context("when no history is found", func() {
			It("should return ClaimHistoryNotFound error", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE claim_id = $1 AND "claim_histories"."deleted_at" IS NULL ORDER BY changed_at DESC`)).
					WithArgs(claimID, 1).
					WillReturnError(gorm.ErrRecordNotFound)

				history, err := repository.FindLatestByClaimID(ctx, claimID)

				Expect(history).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeClaimHistoryNotFound)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE claim_id = $1 AND "claim_histories"."deleted_at" IS NULL ORDER BY changed_at DESC`)).
					WithArgs(claimID, 1).
					WillReturnError(errors.New("database connection failed"))

				history, err := repository.FindLatestByClaimID(ctx, claimID)

				Expect(history).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeDBOperation)
			})
		})
	})

	Describe("FindByDateRange", func() {
		var claimID uuid.UUID
		var startDate, endDate time.Time

		BeforeEach(func() {
			claimID = uuid.New()
			startDate = time.Now().AddDate(0, 0, -7)
			endDate = time.Now()
		})

		Context("when histories within date range are found", func() {
			It("should return filtered histories", func() {
				historyID1 := uuid.New()
				historyID2 := uuid.New()
				changedBy := uuid.New()
				date1 := time.Now().AddDate(0, 0, -5)
				date2 := time.Now().AddDate(0, 0, -3)

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "status", "changed_by", "changed_at", "deleted_at",
				}).AddRow(
					historyID1, claimID, entities.ClaimStatusDraft, changedBy, date1, nil,
				).AddRow(
					historyID2, claimID, entities.ClaimStatusSubmitted, changedBy, date2, nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE (claim_id = $1 AND changed_at BETWEEN $2 AND $3) AND "claim_histories"."deleted_at" IS NULL ORDER BY changed_at DESC`)).
					WithArgs(claimID, startDate, endDate).
					WillReturnRows(rows)

				histories, err := repository.FindByDateRange(ctx, claimID, startDate, endDate)

				Expect(err).NotTo(HaveOccurred())
				Expect(histories).To(HaveLen(2))
				Expect(histories[0].ClaimID).To(Equal(claimID))
				Expect(histories[1].ClaimID).To(Equal(claimID))
			})
		})

		Context("when no histories within date range are found", func() {
			It("should return empty slice", func() {
				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "status", "changed_by", "changed_at", "deleted_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE (claim_id = $1 AND changed_at BETWEEN $2 AND $3) AND "claim_histories"."deleted_at" IS NULL ORDER BY changed_at DESC`)).
					WithArgs(claimID, startDate, endDate).
					WillReturnRows(rows)

				histories, err := repository.FindByDateRange(ctx, claimID, startDate, endDate)

				Expect(err).NotTo(HaveOccurred())
				Expect(histories).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE (claim_id = $1 AND changed_at BETWEEN $2 AND $3) AND "claim_histories"."deleted_at" IS NULL`)).
					WithArgs(claimID, startDate, endDate).
					WillReturnError(errors.New("database connection failed"))

				histories, err := repository.FindByDateRange(ctx, claimID, startDate, endDate)

				Expect(histories).To(BeNil())
				ExpectAppError(err, apperrors.ErrorCodeDBOperation)
			})
		})

		Context("boundary cases for date range", func() {
			It("should handle same start and end date", func() {
				sameDate := time.Now()
				historyID := uuid.New()
				changedBy := uuid.New()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "status", "changed_by", "changed_at", "deleted_at",
				}).AddRow(
					historyID, claimID, entities.ClaimStatusDraft, changedBy, sameDate, nil,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE (claim_id = $1 AND changed_at BETWEEN $2 AND $3) AND "claim_histories"."deleted_at" IS NULL ORDER BY changed_at DESC`)).
					WithArgs(claimID, sameDate, sameDate).
					WillReturnRows(rows)

				histories, err := repository.FindByDateRange(ctx, claimID, sameDate, sameDate)

				Expect(err).NotTo(HaveOccurred())
				Expect(histories).To(HaveLen(1))
			})

			It("should handle very large date range", func() {
				startDate = time.Now().AddDate(-10, 0, 0)
				endDate = time.Now()

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "status", "changed_by", "changed_at", "deleted_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE (claim_id = $1 AND changed_at BETWEEN $2 AND $3) AND "claim_histories"."deleted_at" IS NULL ORDER BY changed_at DESC`)).
					WithArgs(claimID, startDate, endDate).
					WillReturnRows(rows)

				histories, err := repository.FindByDateRange(ctx, claimID, startDate, endDate)

				Expect(err).NotTo(HaveOccurred())
				Expect(histories).To(BeEmpty())
			})

			It("should handle future date range", func() {
				startDate = time.Now().AddDate(0, 0, 1)
				endDate = time.Now().AddDate(0, 0, 7)

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "status", "changed_by", "changed_at", "deleted_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE (claim_id = $1 AND changed_at BETWEEN $2 AND $3) AND "claim_histories"."deleted_at" IS NULL ORDER BY changed_at DESC`)).
					WithArgs(claimID, startDate, endDate).
					WillReturnRows(rows)

				histories, err := repository.FindByDateRange(ctx, claimID, startDate, endDate)

				Expect(err).NotTo(HaveOccurred())
				Expect(histories).To(BeEmpty())
			})

			It("should handle very old date range", func() {
				startDate = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
				endDate = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)

				rows := sqlmock.NewRows([]string{
					"id", "claim_id", "status", "changed_by", "changed_at", "deleted_at",
				})

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "claim_histories" WHERE (claim_id = $1 AND changed_at BETWEEN $2 AND $3) AND "claim_histories"."deleted_at" IS NULL ORDER BY changed_at DESC`)).
					WithArgs(claimID, startDate, endDate).
					WillReturnRows(rows)

				histories, err := repository.FindByDateRange(ctx, claimID, startDate, endDate)

				Expect(err).NotTo(HaveOccurred())
				Expect(histories).To(BeEmpty())
			})
		})
	})
})

func newClaimHistory() *entities.ClaimHistory {
	return &entities.ClaimHistory{
		ID:        uuid.New(),
		ClaimID:   uuid.New(),
		Status:    entities.ClaimStatusDraft,
		ChangedBy: uuid.New(),
		ChangedAt: time.Now(),
		DeletedAt: nil,
	}
}
