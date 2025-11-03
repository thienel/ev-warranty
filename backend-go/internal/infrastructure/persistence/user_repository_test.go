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

var _ = Describe("UserRepository", func() {
	var (
		mock       sqlmock.Sqlmock
		db         *gorm.DB
		repository repository.UserRepository
		ctx        context.Context
	)

	BeforeEach(func() {
		mock, db = SetupMockDB()
		repository = persistence.NewUserRepository(db)
		ctx = context.Background()
	})

	AfterEach(func() {
		CleanupMockDB(mock)
	})

	Describe("Create", func() {
		var user *entity.User

		BeforeEach(func() {
			user = newUser()
		})

		Context("when user is created successfully", func() {
			It("should return nil error", func() {
				MockSuccessfulInsert(mock, "users", user.ID)

				err := repository.Create(ctx, user)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a duplicate email constraint", func() {
			It("should return DBDuplicateKeyError", func() {
				MockDuplicateKeyError(mock, "users", "users_email_key")

				err := repository.Create(ctx, user)

				ExpectAppError(err, apperror.ErrorCodeDuplicateKey)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockInsertError(mock, "users")

				err := repository.Create(ctx, user)

				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("FindByID", func() {
		var userID uuid.UUID

		BeforeEach(func() {
			userID = uuid.New()
		})

		Context("when user is found", func() {
			It("should return the user", func() {
				expected := newUser()
				expected.ID = userID
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "role", "password_hash", "is_active",
					"office_id", "oauth_provider", "oauth_id", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					expected.ID, expected.Name, expected.Email, expected.Role,
					expected.PasswordHash, expected.IsActive, expected.OfficeID,
					expected.OAuthProvider, expected.OAuthID, expected.CreatedAt,
					expected.UpdatedAt, expected.DeletedAt,
				)

				MockFindByID(mock, "users", userID, rows)

				user, err := repository.FindByID(ctx, userID)

				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.ID).To(Equal(expected.ID))
				Expect(user.Name).To(Equal(expected.Name))
				Expect(user.Email).To(Equal(expected.Email))
				Expect(user.Role).To(Equal(expected.Role))
				Expect(user.IsActive).To(Equal(expected.IsActive))
			})
		})

		Context("when user is not found", func() {
			It("should return UserNotFound error", func() {
				MockNotFound(mock, "users", userID)

				user, err := repository.FindByID(ctx, userID)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeUserNotFound)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockQueryError(mock, `SELECT * FROM "users" WHERE id = $1`)

				user, err := repository.FindByID(ctx, userID)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("FindByEmail", func() {
		var email string

		BeforeEach(func() {
			email = "test@example.com"
		})

		Context("when user is found", func() {
			It("should return the user", func() {
				expected := newUser()
				expected.Email = email
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "role", "password_hash", "is_active",
					"office_id", "oauth_provider", "oauth_id", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					expected.ID, expected.Name, expected.Email, expected.Role,
					expected.PasswordHash, expected.IsActive, expected.OfficeID,
					expected.OAuthProvider, expected.OAuthID, expected.CreatedAt,
					expected.UpdatedAt, expected.DeletedAt,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
					WithArgs(email, 1).
					WillReturnRows(rows)

				user, err := repository.FindByEmail(ctx, email)

				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.Email).To(Equal(email))
			})
		})

		Context("when user is not found", func() {
			It("should return UserNotFound error", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
					WithArgs(email, 1).
					WillReturnError(gorm.ErrRecordNotFound)

				user, err := repository.FindByEmail(ctx, email)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeUserNotFound)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
					WithArgs(email, 1).
					WillReturnError(errors.New("database connection failed"))

				user, err := repository.FindByEmail(ctx, email)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})

		Context("boundary cases for email", func() {
			It("should handle empty email", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
					WithArgs("", 1).
					WillReturnError(gorm.ErrRecordNotFound)

				user, err := repository.FindByEmail(ctx, "")

				Expect(err).To(HaveOccurred())
				Expect(user).To(BeNil())
			})

			It("should handle very long email", func() {
				longEmail := "very.long.email.address.that.exceeds.normal.length@verylongdomainname.com"
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
					WithArgs(longEmail, 1).
					WillReturnError(gorm.ErrRecordNotFound)

				user, err := repository.FindByEmail(ctx, longEmail)

				Expect(err).To(HaveOccurred())
				Expect(user).To(BeNil())
			})
		})
	})

	Describe("FindAll", func() {
		Context("when users are found", func() {
			It("should return all users", func() {
				userID1 := uuid.New()
				userID2 := uuid.New()
				officeID := uuid.New()
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "role", "password_hash", "is_active",
					"office_id", "oauth_provider", "oauth_id", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					userID1, "User 1", "user1@test.com", entity.UserRoleAdmin,
					"hash1", true, officeID, nil, nil, time.Now(), time.Now(), nil,
				).AddRow(
					userID2, "User 2", "user2@test.com", entity.UserRoleEvmStaff,
					"hash2", false, officeID, nil, nil, time.Now(), time.Now(), nil,
				)

				MockFindAll(mock, "users", rows)

				users, err := repository.FindAll(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(users).To(HaveLen(2))
				Expect(users[0].ID).To(Equal(userID1))
				Expect(users[0].Name).To(Equal("User 1"))
				Expect(users[1].ID).To(Equal(userID2))
				Expect(users[1].Name).To(Equal("User 2"))
			})
		})

		Context("when no users are found", func() {
			It("should return an empty slice", func() {
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "role", "password_hash", "is_active",
					"office_id", "oauth_provider", "oauth_id", "created_at", "updated_at", "deleted_at",
				})

				MockFindAll(mock, "users", rows)

				users, err := repository.FindAll(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(users).To(BeEmpty())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockQueryError(mock, `SELECT * FROM "users"`)

				users, err := repository.FindAll(ctx)

				Expect(users).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("Update", func() {
		var user *entity.User

		BeforeEach(func() {
			user = newUser()
			user.Name = "Updated Name"
			user.Email = "updated@test.com"
		})

		Context("when user is updated successfully", func() {
			It("should return nil error", func() {
				MockSuccessfulUpdate(mock, "users")

				err := repository.Update(ctx, user)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockUpdateError(mock, "users")

				err := repository.Update(ctx, user)

				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})

		Context("boundary cases for update", func() {
			It("should handle updating is_active to false", func() {
				user.IsActive = false
				MockSuccessfulUpdate(mock, "users")

				err := repository.Update(ctx, user)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should handle updating role", func() {
				user.Role = entity.UserRoleScTechnician
				MockSuccessfulUpdate(mock, "users")

				err := repository.Update(ctx, user)

				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("SoftDelete", func() {
		var userID uuid.UUID

		BeforeEach(func() {
			userID = uuid.New()
		})

		Context("when user is soft deleted successfully", func() {
			It("should return nil error", func() {
				MockSoftDelete(mock, "users", userID)

				err := repository.SoftDelete(ctx, userID)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				MockDeleteError(mock, "users")

				err := repository.SoftDelete(ctx, userID)

				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})
	})

	Describe("FindByOAuth", func() {
		var provider, oauthID string

		BeforeEach(func() {
			provider = "google"
			oauthID = "oauth123456"
		})

		Context("when user is found", func() {
			It("should return the user", func() {
				expected := newUserWithOAuth(provider, oauthID)
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "role", "password_hash", "is_active",
					"office_id", "oauth_provider", "oauth_id", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					expected.ID, expected.Name, expected.Email, expected.Role,
					expected.PasswordHash, expected.IsActive, expected.OfficeID,
					expected.OAuthProvider, expected.OAuthID, expected.CreatedAt,
					expected.UpdatedAt, expected.DeletedAt,
				)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (oauth_provider = $1 AND oauth_id = $2) AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $3`)).
					WithArgs(provider, oauthID, 1).
					WillReturnRows(rows)

				user, err := repository.FindByOAuth(ctx, provider, oauthID)

				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(*user.OAuthProvider).To(Equal(provider))
				Expect(*user.OAuthID).To(Equal(oauthID))
			})
		})

		Context("when user is not found", func() {
			It("should return UserNotFound error", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (oauth_provider = $1 AND oauth_id = $2) AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $3`)).
					WithArgs(provider, oauthID, 1).
					WillReturnError(gorm.ErrRecordNotFound)

				user, err := repository.FindByOAuth(ctx, provider, oauthID)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeUserNotFound)
			})
		})

		Context("when there is a database error", func() {
			It("should return DBOperationError", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (oauth_provider = $1 AND oauth_id = $2) AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $3`)).
					WithArgs(provider, oauthID, 1).
					WillReturnError(errors.New("database connection failed"))

				user, err := repository.FindByOAuth(ctx, provider, oauthID)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeDBOperation)
			})
		})

		Context("boundary cases for OAuth", func() {
			It("should handle empty provider", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (oauth_provider = $1 AND oauth_id = $2) AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $3`)).
					WithArgs("", oauthID, 1).
					WillReturnError(gorm.ErrRecordNotFound)

				user, err := repository.FindByOAuth(ctx, "", oauthID)

				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeUserNotFound)
			})

			It("should handle empty oauth_id", func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (oauth_provider = $1 AND oauth_id = $2) AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $3`)).
					WithArgs(provider, "", 1).
					WillReturnError(gorm.ErrRecordNotFound)

				user, err := repository.FindByOAuth(ctx, provider, "")

				Expect(user).To(BeNil())
				ExpectAppError(err, apperror.ErrorCodeUserNotFound)
			})
		})
	})
})

func newUser() *entity.User {
	officeID := uuid.New()
	return &entity.User{
		ID:           uuid.New(),
		Name:         "Test User",
		Email:        "test@example.com",
		Role:         entity.UserRoleAdmin,
		PasswordHash: "hashedpassword",
		IsActive:     true,
		OfficeID:     officeID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		DeletedAt:    gorm.DeletedAt{},
	}
}

func newUserWithOAuth(provider, oauthID string) *entity.User {
	user := newUser()
	user.OAuthProvider = &provider
	user.OAuthID = &oauthID
	return user
}
