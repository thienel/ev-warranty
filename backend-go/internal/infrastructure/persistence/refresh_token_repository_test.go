package persistence_test

import (
	"context"
	"errors"
	"ev-warranty-go/pkg/apperror"
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

var _ = Describe("RefreshTokenRepository", func() {
	var (
		mock       sqlmock.Sqlmock
		db         *gorm.DB
		repository repository.RefreshTokenRepository
		ctx        context.Context
	)

	BeforeEach(func() {
		mock, db = SetupMockDB()
		repository = persistence.NewTokenRepository(db)
		ctx = context.Background()
	})

	AfterEach(func() {
		CleanupMockDB(mock)
	})

	Describe("Create", func() {
		var token *entity.RefreshToken

		BeforeEach(func() {
			token = newRefreshToken()
		})

		Context("when token is created successfully", func() {
			It("should return nil error", func() {
				MockSuccessfulInsert(mock, "refresh_tokens", token.ID)

				err := repository.Create(ctx, token)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a duplicate key constraint", func() {
			It("should return DBDuplicateKeyError", func() {
				MockDuplicateKeyError(mock, "refresh_tokens", "refresh_tokens_token_key")

				err := repository.Create(ctx, token)

				ExpectAppError(err, apperror.ErrorCodeDuplicateKey)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockInsertError(mock, "refresh_tokens")

				err := repository.Create(ctx, token)

				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})

		Context("boundary cases for token", func() {
			It("should handle very long token string", func() {
				token.Token = string(make([]byte, 1000))
				MockSuccessfulInsert(mock, "refresh_tokens", token.ID)

				err := repository.Create(ctx, token)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should handle empty token string", func() {
				token.Token = ""
				MockSuccessfulInsert(mock, "refresh_tokens", token.ID)

				err := repository.Create(ctx, token)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("boundary cases for expiration", func() {
			It("should handle already expired token", func() {
				token.ExpiresAt = time.Now().Add(-1 * time.Hour)
				MockSuccessfulInsert(mock, "refresh_tokens", token.ID)

				err := repository.Create(ctx, token)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should handle far future expiration", func() {
				token.ExpiresAt = time.Now().AddDate(10, 0, 0)
				MockSuccessfulInsert(mock, "refresh_tokens", token.ID)

				err := repository.Create(ctx, token)

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("Update", func() {
		var token *entity.RefreshToken

		BeforeEach(func() {
			token = newRefreshToken()
			token.IsRevoked = true
		})

		Context("when token is updated successfully", func() {
			It("should return nil error", func() {
				MockSuccessfulUpdate(mock, "refresh_tokens")

				err := repository.Update(ctx, token)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockUpdateError(mock, "refresh_tokens")

				err := repository.Update(ctx, token)

				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})

		Context("boundary cases", func() {
			It("should handle setting is_revoked to false", func() {
				token.IsRevoked = false
				MockSuccessfulUpdate(mock, "refresh_tokens")

				err := repository.Update(ctx, token)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should handle setting is_revoked to true", func() {
				token.IsRevoked = true
				MockSuccessfulUpdate(mock, "refresh_tokens")

				err := repository.Update(ctx, token)

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("Find", func() {
		var tokenStr string

		BeforeEach(func() {
			tokenStr = "test-refresh-token-12345"
		})

		Context("when token is found", func() {
			It("should return the token", func() {
				expected := newRefreshToken()
				expected.Token = tokenStr
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "token", "expires_at", "is_revoked", "created_at", "updated_at",
				}).AddRow(
					expected.ID, expected.UserID, expected.Token, expected.ExpiresAt,
					expected.IsRevoked, expected.CreatedAt, expected.UpdatedAt,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "refresh_tokens" WHERE token = $1`)).
					WithArgs(tokenStr, 1).
					WillReturnRows(rows)

				token, err := repository.Find(ctx, tokenStr)

				Expect(err).NotTo(HaveOccurred())
				Expect(token).NotTo(BeNil())
				Expect(token.Token).To(Equal(tokenStr))
				Expect(token.UserID).To(Equal(expected.UserID))
			})
		})

		Context("when token is not found", func() {
			It("should return RefreshTokenNotFound error", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "refresh_tokens" WHERE token = $1`)).
					WithArgs(tokenStr, 1).
					WillReturnError(gorm.ErrRecordNotFound)

				token, err := repository.Find(ctx, tokenStr)

				Expect(token).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeRefreshTokenNotFound)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "refresh_tokens" WHERE token = $1`)).
					WithArgs(tokenStr, 1).
					WillReturnError(errors.New("database connection failed"))

				token, err := repository.Find(ctx, tokenStr)

				Expect(token).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})

		Context("boundary cases for token string", func() {
			It("should handle empty token string", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "refresh_tokens" WHERE token = $1`)).
					WithArgs("", 1).
					WillReturnError(gorm.ErrRecordNotFound)

				token, err := repository.Find(ctx, "")

				Expect(err).To(HaveOccurred())
				Expect(token).To(BeNil())
			})

			It("should handle very long token string", func() {
				longToken := string(make([]byte, 1000))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "refresh_tokens" WHERE token = $1`)).
					WithArgs(longToken, 1).
					WillReturnError(gorm.ErrRecordNotFound)

				token, err := repository.Find(ctx, longToken)

				Expect(err).To(HaveOccurred())
				Expect(token).To(BeNil())
			})

			It("should handle token with special characters", func() {
				specialToken := "token!@#$%^&*()_+-={}[]|:;<>?,./"
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "refresh_tokens" WHERE token = $1`)).
					WithArgs(specialToken, 1).
					WillReturnError(gorm.ErrRecordNotFound)

				token, err := repository.Find(ctx, specialToken)

				Expect(err).To(HaveOccurred())
				Expect(token).To(BeNil())
			})
		})

		Context("equivalence partitioning for token state", func() {
			It("should find non-revoked token", func() {
				expected := newRefreshToken()
				expected.Token = tokenStr
				expected.IsRevoked = false
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "token", "expires_at", "is_revoked", "created_at", "updated_at",
				}).AddRow(
					expected.ID, expected.UserID, expected.Token, expected.ExpiresAt,
					expected.IsRevoked, expected.CreatedAt, expected.UpdatedAt,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "refresh_tokens" WHERE token = $1`)).
					WithArgs(tokenStr, 1).
					WillReturnRows(rows)

				token, err := repository.Find(ctx, tokenStr)

				Expect(err).NotTo(HaveOccurred())
				Expect(token.IsRevoked).To(BeFalse())
			})

			It("should find revoked token", func() {
				expected := newRefreshToken()
				expected.Token = tokenStr
				expected.IsRevoked = true
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "token", "expires_at", "is_revoked", "created_at", "updated_at",
				}).AddRow(
					expected.ID, expected.UserID, expected.Token, expected.ExpiresAt,
					expected.IsRevoked, expected.CreatedAt, expected.UpdatedAt,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "refresh_tokens" WHERE token = $1`)).
					WithArgs(tokenStr, 1).
					WillReturnRows(rows)

				token, err := repository.Find(ctx, tokenStr)

				Expect(err).NotTo(HaveOccurred())
				Expect(token.IsRevoked).To(BeTrue())
			})

			It("should find expired token", func() {
				expected := newRefreshToken()
				expected.Token = tokenStr
				expected.ExpiresAt = time.Now().Add(-1 * time.Hour)
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "token", "expires_at", "is_revoked", "created_at", "updated_at",
				}).AddRow(
					expected.ID, expected.UserID, expected.Token, expected.ExpiresAt,
					expected.IsRevoked, expected.CreatedAt, expected.UpdatedAt,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "refresh_tokens" WHERE token = $1`)).
					WithArgs(tokenStr, 1).
					WillReturnRows(rows)

				token, err := repository.Find(ctx, tokenStr)

				Expect(err).NotTo(HaveOccurred())
				Expect(token.IsExpired()).To(BeTrue())
			})

			It("should find non-expired token", func() {
				expected := newRefreshToken()
				expected.Token = tokenStr
				expected.ExpiresAt = time.Now().Add(24 * time.Hour)
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "token", "expires_at", "is_revoked", "created_at", "updated_at",
				}).AddRow(
					expected.ID, expected.UserID, expected.Token, expected.ExpiresAt,
					expected.IsRevoked, expected.CreatedAt, expected.UpdatedAt,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "refresh_tokens" WHERE token = $1`)).
					WithArgs(tokenStr, 1).
					WillReturnRows(rows)

				token, err := repository.Find(ctx, tokenStr)

				Expect(err).NotTo(HaveOccurred())
				Expect(token.IsExpired()).To(BeFalse())
			})
		})
	})

	Describe("Revoke", func() {
		var tokenStr string

		BeforeEach(func() {
			tokenStr = "test-refresh-token-12345"
		})

		Context("when token is revoked successfully", func() {
			It("should return nil error", func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "refresh_tokens" SET "is_revoked"=$1,"updated_at"=$2 WHERE token = $3`)).
					WithArgs(true, sqlmock.AnyArg(), tokenStr).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := repository.Revoke(ctx, tokenStr)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "refresh_tokens" SET "is_revoked"=$1,"updated_at"=$2 WHERE token = $3`)).
					WithArgs(true, sqlmock.AnyArg(), tokenStr).
					WillReturnError(errors.New("database connection failed"))
				mock.ExpectRollback()

				err := repository.Revoke(ctx, tokenStr)

				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})

		Context("boundary cases", func() {
			It("should handle revoking non-existent token", func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "refresh_tokens" SET "is_revoked"=$1,"updated_at"=$2 WHERE token = $3`)).
					WithArgs(true, sqlmock.AnyArg(), tokenStr).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()

				err := repository.Revoke(ctx, tokenStr)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should handle empty token string", func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "refresh_tokens" SET "is_revoked"=$1,"updated_at"=$2 WHERE token = $3`)).
					WithArgs(true, sqlmock.AnyArg(), "").
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectCommit()

				err := repository.Revoke(ctx, "")

				Expect(err).NotTo(HaveOccurred())
			})

			It("should handle revoking already revoked token", func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "refresh_tokens" SET "is_revoked"=$1,"updated_at"=$2 WHERE token = $3`)).
					WithArgs(true, sqlmock.AnyArg(), tokenStr).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				err := repository.Revoke(ctx, tokenStr)

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})

func newRefreshToken() *entity.RefreshToken {
	return &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "test-refresh-token-" + uuid.New().String(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IsRevoked: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
