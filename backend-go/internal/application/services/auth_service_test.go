package services_test

import (
	"context"
	"errors"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/internal/infrastructure/oauth/providers"
	"ev-warranty-go/internal/security"
	"ev-warranty-go/pkg/mocks"
)

var _ = Describe("AuthService", func() {
	var (
		mockUserRepo  *mocks.UserRepository
		mockTokenServ *mocks.TokenService
		service       services.AuthService
		ctx           context.Context
	)

	BeforeEach(func() {
		mockUserRepo = mocks.NewUserRepository(GinkgoT())
		mockTokenServ = mocks.NewTokenService(GinkgoT())
		service = services.NewAuthService(mockUserRepo, mockTokenServ)
		ctx = context.Background()
	})

	Describe("Login", func() {
		var (
			email    string
			password string
		)

		BeforeEach(func() {
			email = "user@example.com"
			password = "Password123!"
		})

		Context("when login is successful", func() {
			It("should return access and refresh tokens", func() {
				passwordHash, _ := security.HashPassword(password)
				user := &entities.User{
					ID:           uuid.New(),
					Email:        email,
					PasswordHash: passwordHash,
					IsActive:     true,
				}

				accessToken := "access.token.jwt"
				refreshToken := "refresh_token_string"

				mockUserRepo.EXPECT().FindByEmail(ctx, email).Return(user, nil).Once()
				mockTokenServ.EXPECT().GenerateAccessToken(user.ID).Return(accessToken, nil).Once()
				mockTokenServ.EXPECT().GenerateRefreshToken(ctx, user.ID).Return(refreshToken, nil).Once()

				at, rt, err := service.Login(ctx, email, password)

				Expect(err).NotTo(HaveOccurred())
				Expect(at).To(Equal(accessToken))
				Expect(rt).To(Equal(refreshToken))
			})
		})

		Context("when email has leading/trailing spaces", func() {
			It("should trim spaces and login successfully", func() {
				emailWithSpaces := "  user@example.com  "
				passwordHash, _ := security.HashPassword(password)
				user := &entities.User{
					ID:           uuid.New(),
					Email:        email,
					PasswordHash: passwordHash,
					IsActive:     true,
				}

				accessToken := "access.token.jwt"
				refreshToken := "refresh_token_string"

				mockUserRepo.EXPECT().FindByEmail(ctx, email).Return(user, nil).Once()
				mockTokenServ.EXPECT().GenerateAccessToken(user.ID).Return(accessToken, nil).Once()
				mockTokenServ.EXPECT().GenerateRefreshToken(ctx, user.ID).Return(refreshToken, nil).Once()

				at, rt, err := service.Login(ctx, emailWithSpaces, password)

				Expect(err).NotTo(HaveOccurred())
				Expect(at).To(Equal(accessToken))
				Expect(rt).To(Equal(refreshToken))
			})
		})

		Context("when user is not found", func() {
			It("should return UserNotFound error", func() {
				notFoundErr := apperrors.NewUserNotFound()
				mockUserRepo.EXPECT().FindByEmail(ctx, email).Return(nil, notFoundErr).Once()

				at, rt, err := service.Login(ctx, email, password)

				Expect(at).To(BeEmpty())
				Expect(rt).To(BeEmpty())
				Expect(err).To(Equal(notFoundErr))
			})
		})

		Context("when user is inactive", func() {
			It("should return UserInactive error", func() {
				passwordHash, _ := security.HashPassword(password)
				user := &entities.User{
					ID:           uuid.New(),
					Email:        email,
					PasswordHash: passwordHash,
					IsActive:     false,
				}

				mockUserRepo.EXPECT().FindByEmail(ctx, email).Return(user, nil).Once()

				at, rt, err := service.Login(ctx, email, password)

				Expect(at).To(BeEmpty())
				Expect(rt).To(BeEmpty())
				ExpectAppError(err, apperrors.ErrorCodeUserInactive)
			})
		})

		Context("when password is incorrect", func() {
			It("should return UserPasswordInvalid error", func() {
				passwordHash, _ := security.HashPassword("DifferentPassword123!")
				user := &entities.User{
					ID:           uuid.New(),
					Email:        email,
					PasswordHash: passwordHash,
					IsActive:     true,
				}

				mockUserRepo.EXPECT().FindByEmail(ctx, email).Return(user, nil).Once()

				at, rt, err := service.Login(ctx, email, password)

				Expect(at).To(BeEmpty())
				Expect(rt).To(BeEmpty())
				ExpectAppError(err, apperrors.ErrorCodeUserPasswordInvalid)
			})
		})

		Context("when FindByEmail returns error", func() {
			It("should return the error", func() {
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))
				mockUserRepo.EXPECT().FindByEmail(ctx, email).Return(nil, dbErr).Once()

				at, rt, err := service.Login(ctx, email, password)

				Expect(err).To(HaveOccurred())
				Expect(at).To(BeEmpty())
				Expect(rt).To(BeEmpty())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when GenerateAccessToken fails", func() {
			It("should return the error", func() {
				passwordHash, _ := security.HashPassword(password)
				user := &entities.User{
					ID:           uuid.New(),
					Email:        email,
					PasswordHash: passwordHash,
					IsActive:     true,
				}

				tokenErr := apperrors.New(500, apperrors.ErrorCodeFailedSignAccessToken, errors.New("token error"))
				mockUserRepo.EXPECT().FindByEmail(ctx, email).Return(user, nil).Once()
				mockTokenServ.EXPECT().GenerateAccessToken(user.ID).Return("", tokenErr).Once()

				at, rt, err := service.Login(ctx, email, password)

				Expect(err).To(HaveOccurred())
				Expect(at).To(BeEmpty())
				Expect(rt).To(BeEmpty())
				Expect(err).To(Equal(tokenErr))
			})
		})

		Context("when GenerateRefreshToken fails", func() {
			It("should return the error", func() {
				passwordHash, _ := security.HashPassword(password)
				user := &entities.User{
					ID:           uuid.New(),
					Email:        email,
					PasswordHash: passwordHash,
					IsActive:     true,
				}

				accessToken := "access.token.jwt"
				tokenErr := apperrors.New(500, apperrors.ErrorCodeFailedGenerateRefreshToken, errors.New("token error"))

				mockUserRepo.EXPECT().FindByEmail(ctx, email).Return(user, nil).Once()
				mockTokenServ.EXPECT().GenerateAccessToken(user.ID).Return(accessToken, nil).Once()
				mockTokenServ.EXPECT().GenerateRefreshToken(ctx, user.ID).Return("", tokenErr).Once()

				at, rt, err := service.Login(ctx, email, password)

				Expect(err).To(HaveOccurred())
				Expect(at).To(BeEmpty())
				Expect(rt).To(BeEmpty())
				Expect(err).To(Equal(tokenErr))
			})
		})
	})

	Describe("Logout", func() {
		var token string

		BeforeEach(func() {
			token = "refresh_token_string"
		})

		Context("when logout is successful", func() {
			It("should return nil error", func() {
				mockTokenServ.EXPECT().RevokeRefreshToken(ctx, token).Return(nil).Once()

				err := service.Logout(ctx, token)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when RevokeRefreshToken fails", func() {
			It("should return the error", func() {
				tokenErr := apperrors.New(404, apperrors.ErrorCodeRefreshTokenNotFound, errors.New("token not found"))
				mockTokenServ.EXPECT().RevokeRefreshToken(ctx, token).Return(tokenErr).Once()

				err := service.Logout(ctx, token)

				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(tokenErr))
			})
		})
	})

	Describe("HandleOAuthUser", func() {
		var userInfo *providers.UserInfo

		BeforeEach(func() {
			userInfo = &providers.UserInfo{
				Provider:   "google",
				ProviderID: "google_12345",
				Email:      "user@example.com",
				Name:       "John Doe",
			}
		})

		Context("when OAuth user exists", func() {
			It("should return access and refresh tokens", func() {
				user := &entities.User{
					ID:            uuid.New(),
					Email:         userInfo.Email,
					OAuthProvider: &userInfo.Provider,
					OAuthID:       &userInfo.ProviderID,
					IsActive:      true,
				}

				accessToken := "access.token.jwt"
				refreshToken := "refresh_token_string"

				mockUserRepo.EXPECT().FindByOAuth(ctx, userInfo.Provider, userInfo.ProviderID).Return(user, nil).Once()
				mockTokenServ.EXPECT().GenerateAccessToken(user.ID).Return(accessToken, nil).Once()
				mockTokenServ.EXPECT().GenerateRefreshToken(ctx, user.ID).Return(refreshToken, nil).Once()

				at, rt, err := service.HandleOAuthUser(ctx, userInfo)

				Expect(err).NotTo(HaveOccurred())
				Expect(at).To(Equal(accessToken))
				Expect(rt).To(Equal(refreshToken))
			})
		})

		Context("when DB operation error", func() {
			It("should return DBOperationError", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))

				mockUserRepo.EXPECT().FindByOAuth(ctx, userInfo.Provider, userInfo.ProviderID).Return(nil, dbErr).Once()

				at, rt, err := service.HandleOAuthUser(ctx, userInfo)

				Expect(at).To(BeEmpty())
				Expect(rt).To(BeEmpty())
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when OAuth user not found but email exists", func() {
			It("should link OAuth and return tokens", func() {
				user := &entities.User{
					ID:       uuid.New(),
					Email:    userInfo.Email,
					IsActive: true,
				}

				accessToken := "access.token.jwt"
				refreshToken := "refresh_token_string"

				mockUserRepo.EXPECT().FindByOAuth(ctx, userInfo.Provider, userInfo.ProviderID).
					Return(nil, apperrors.NewUserNotFound()).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, userInfo.Email).Return(user, nil).Once()
				mockUserRepo.EXPECT().Update(ctx, mock.MatchedBy(func(u *entities.User) bool {
					return u.ID == user.ID &&
						u.OAuthProvider != nil &&
						*u.OAuthProvider == userInfo.Provider &&
						u.OAuthID != nil &&
						*u.OAuthID == userInfo.ProviderID
				})).Return(nil).Once()
				mockTokenServ.EXPECT().GenerateAccessToken(user.ID).Return(accessToken, nil).Once()
				mockTokenServ.EXPECT().GenerateRefreshToken(ctx, user.ID).Return(refreshToken, nil).Once()

				at, rt, err := service.HandleOAuthUser(ctx, userInfo)

				Expect(err).NotTo(HaveOccurred())
				Expect(at).To(Equal(accessToken))
				Expect(rt).To(Equal(refreshToken))
			})
		})

		Context("when FindByEmail fails", func() {
			It("should return the error", func() {
				dbErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("database error"))

				mockUserRepo.EXPECT().FindByOAuth(mock.Anything, userInfo.Provider, userInfo.ProviderID).
					Return(nil, apperrors.NewUserNotFound()).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, userInfo.Email).Return(nil, dbErr).Once()

				at, rt, err := service.HandleOAuthUser(ctx, userInfo)

				Expect(err).To(HaveOccurred())
				Expect(at).To(BeEmpty())
				Expect(rt).To(BeEmpty())
				Expect(err).To(Equal(dbErr))
			})
		})

		Context("when Update fails after linking OAuth", func() {
			It("should return the error", func() {
				user := &entities.User{
					ID:       uuid.New(),
					Email:    userInfo.Email,
					IsActive: true,
				}

				updateErr := apperrors.New(500, apperrors.ErrorCodeDBOperation, errors.New("update error"))

				mockUserRepo.EXPECT().FindByOAuth(ctx, userInfo.Provider, userInfo.ProviderID).
					Return(nil, apperrors.NewUserNotFound()).Once()
				mockUserRepo.EXPECT().FindByEmail(ctx, userInfo.Email).Return(user, nil).Once()
				mockUserRepo.EXPECT().Update(ctx, mock.AnythingOfType("*entities.User")).Return(updateErr).Once()

				at, rt, err := service.HandleOAuthUser(ctx, userInfo)

				Expect(err).To(HaveOccurred())
				Expect(at).To(BeEmpty())
				Expect(rt).To(BeEmpty())
				Expect(err).To(Equal(updateErr))
			})
		})

		Context("when GenerateAccessToken fails", func() {
			It("should return the error", func() {
				user := &entities.User{
					ID:            uuid.New(),
					Email:         userInfo.Email,
					OAuthProvider: &userInfo.Provider,
					OAuthID:       &userInfo.ProviderID,
					IsActive:      true,
				}

				tokenErr := apperrors.New(500, apperrors.ErrorCodeFailedSignAccessToken, errors.New("token error"))

				mockUserRepo.EXPECT().FindByOAuth(ctx, userInfo.Provider, userInfo.ProviderID).Return(user, nil).Once()
				mockTokenServ.EXPECT().GenerateAccessToken(user.ID).Return("", tokenErr).Once()

				at, rt, err := service.HandleOAuthUser(ctx, userInfo)

				Expect(err).To(HaveOccurred())
				Expect(at).To(BeEmpty())
				Expect(rt).To(BeEmpty())
				Expect(err).To(Equal(tokenErr))
			})
		})

		Context("when GenerateRefreshToken fails", func() {
			It("should return the error", func() {
				user := &entities.User{
					ID:            uuid.New(),
					Email:         userInfo.Email,
					OAuthProvider: &userInfo.Provider,
					OAuthID:       &userInfo.ProviderID,
					IsActive:      true,
				}

				accessToken := "access.token.jwt"
				tokenErr := apperrors.New(500, apperrors.ErrorCodeFailedGenerateRefreshToken, errors.New("token error"))

				mockUserRepo.EXPECT().FindByOAuth(ctx, userInfo.Provider, userInfo.ProviderID).Return(user, nil).Once()
				mockTokenServ.EXPECT().GenerateAccessToken(user.ID).Return(accessToken, nil).Once()
				mockTokenServ.EXPECT().GenerateRefreshToken(ctx, user.ID).Return("", tokenErr).Once()

				at, rt, err := service.HandleOAuthUser(ctx, userInfo)

				Expect(err).To(HaveOccurred())
				Expect(at).To(BeEmpty())
				Expect(rt).To(BeEmpty())
				Expect(err).To(Equal(tokenErr))
			})
		})
	})
})
