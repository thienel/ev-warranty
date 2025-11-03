package handler_test

import (
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/internal/interface/api/dto"
	"ev-warranty-go/internal/interface/api/handler"
	apperrors2 "ev-warranty-go/pkg/apperror"
	"ev-warranty-go/pkg/mocks"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("AuthHandler", func() {
	var (
		mockLogger   *mocks.Logger
		mockAuthSvc  *mocks.AuthService
		mockTokenSvc *mocks.TokenService
		mockUserSvc  *mocks.UserService
		authHandler  handler.AuthHandler
		r            *gin.Engine
		w            *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		mockLogger, r, w = SetupMock(GinkgoT())
		mockAuthSvc = mocks.NewAuthService(GinkgoT())
		mockTokenSvc = mocks.NewTokenService(GinkgoT())
		mockUserSvc = mocks.NewUserService(GinkgoT())
		authHandler = handler.NewAuthHandler(mockLogger, mockAuthSvc, mockTokenSvc, mockUserSvc)
	})

	Describe("Login", func() {
		var (
			loginRequest = dto.LoginRequest{Email: "test@example.com", Password: "password123"}
			userID       = uuid.New()
			accessToken  = "access-token-123"
			refreshToken = "refresh-token-456"
			user         *entities.User
		)

		BeforeEach(func() {
			r.POST("/auth/login", authHandler.Login)
			user = entities.NewUser("Test User", "test@example.com", entities.UserRoleAdmin, "password_hash", true, uuid.New())
			user.ID = userID
		})

		It("should return access token and user info with refresh token cookie", func() {
			claims := &services.CustomClaims{UserID: userID.String()}

			mockAuthSvc.EXPECT().Login(mock.Anything, loginRequest.Email, loginRequest.Password).
				Return(accessToken, refreshToken, nil).Once()
			mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
				Return(claims, nil).Once()
			mockUserSvc.EXPECT().GetByID(mock.Anything, userID).Return(user, nil).Once()

			SendRequest(r, http.MethodPost, "/auth/login", w, loginRequest)

			Expect(w.Code).To(Equal(http.StatusOK))
			ExpectCookieRefreshToken(w, refreshToken)
			ExpectResponseNotNil(w, http.StatusOK)
		})

		DescribeTable("should handle errors",
			func(setupMocks func(), request interface{}, expectedCode int, expectedError string) {
				if setupMocks != nil {
					setupMocks()
				}
				SendRequest(r, http.MethodPost, "/auth/login", w, request)
				ExpectErrorCode(w, expectedCode, expectedError)
			},
			Entry("invalid JSON", nil, "invalid json", http.StatusBadRequest, apperrors2.ErrorCodeInvalidJsonRequest),
			Entry("auth service login fails", func() {
				mockAuthSvc.EXPECT().Login(mock.Anything, loginRequest.Email, loginRequest.Password).
					Return("", "", apperrors2.NewInvalidCredentials()).Once()
			}, loginRequest, http.StatusUnauthorized, apperrors2.ErrorCodeInvalidCredentials),
			Entry("token validation fails", func() {
				mockAuthSvc.EXPECT().Login(mock.Anything, loginRequest.Email, loginRequest.Password).
					Return(accessToken, refreshToken, nil).Once()
				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(nil, apperrors2.NewInvalidAccessToken()).Once()
			}, loginRequest, http.StatusUnauthorized, apperrors2.ErrorCodeInvalidAccessToken),
			Entry("user not found", func() {
				claims := &services.CustomClaims{UserID: userID.String()}
				mockAuthSvc.EXPECT().Login(mock.Anything, loginRequest.Email, loginRequest.Password).
					Return(accessToken, refreshToken, nil).Once()
				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(claims, nil).Once()
				mockUserSvc.EXPECT().GetByID(mock.Anything, userID).
					Return(nil, apperrors2.NewUserNotFound()).Once()
			}, loginRequest, http.StatusNotFound, apperrors2.ErrorCodeUserNotFound),
		)
	})

	Describe("Logout", func() {
		var refreshToken = "refresh-token-456"

		BeforeEach(func() {
			r.POST("/auth/logout", authHandler.Logout)
		})

		It("should logout user and return success", func() {
			mockAuthSvc.EXPECT().Logout(mock.Anything, refreshToken).Return(nil).Once()

			req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
			req.AddCookie(&http.Cookie{Name: "refreshToken", Value: refreshToken})
			r.ServeHTTP(w, req)

			cookies := w.Result().Cookies()
			Expect(cookies).To(HaveLen(1))
			Expect(cookies[0].Name).To(Equal("refreshToken"))
			Expect(cookies[0].Value).To(Equal(""))
			Expect(w.Code).To(Equal(http.StatusNoContent))
		})

		DescribeTable("should handle errors",
			func(setupCookie bool, setupMocks func(), expectedCode int, expectedError string) {
				if setupMocks != nil {
					setupMocks()
				}
				req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
				if setupCookie {
					req.AddCookie(&http.Cookie{Name: "refreshToken", Value: refreshToken})
				}
				r.ServeHTTP(w, req)
				ExpectErrorCode(w, expectedCode, expectedError)
			},
			Entry("missing refresh token", false, nil, http.StatusNotFound, apperrors2.ErrorCodeRefreshTokenNotFound),
			Entry("auth service logout fails", true, func() {
				mockAuthSvc.EXPECT().Logout(mock.Anything, refreshToken).
					Return(apperrors2.NewInvalidRefreshToken()).Once()
			}, http.StatusUnauthorized, apperrors2.ErrorCodeInvalidRefreshToken),
		)
	})

	Describe("RefreshToken", func() {
		var (
			refreshToken   = "refresh-token-456"
			newAccessToken = "new-access-token-789"
		)

		BeforeEach(func() {
			r.POST("/auth/refresh", authHandler.RefreshToken)
		})

		It("should return new access token", func() {
			mockTokenSvc.EXPECT().RefreshAccessToken(mock.Anything, refreshToken).
				Return(newAccessToken, nil).Once()

			req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
			req.AddCookie(&http.Cookie{Name: "refreshToken", Value: refreshToken})
			r.ServeHTTP(w, req)

			ExpectResponseNotNil(w, http.StatusOK)
		})

		DescribeTable("should handle errors",
			func(setupCookie bool, setupMocks func(), expectedCode int, expectedError string) {
				if setupMocks != nil {
					setupMocks()
				}
				req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
				if setupCookie {
					req.AddCookie(&http.Cookie{Name: "refreshToken", Value: refreshToken})
				}
				r.ServeHTTP(w, req)
				ExpectErrorCode(w, expectedCode, expectedError)
			},
			Entry("missing refresh token", false, nil, http.StatusNotFound, apperrors2.ErrorCodeRefreshTokenNotFound),
			Entry("token service refresh fails", true, func() {
				mockTokenSvc.EXPECT().RefreshAccessToken(mock.Anything, refreshToken).
					Return("", apperrors2.NewInvalidRefreshToken()).Once()
			}, http.StatusUnauthorized, apperrors2.ErrorCodeInvalidRefreshToken),
		)
	})

	Describe("ValidateToken", func() {
		var (
			accessToken = "valid-access-token"
			userID      = uuid.New()
			user        *entities.User
		)

		BeforeEach(func() {
			r.POST("/auth/validate", authHandler.ValidateToken)
			user = &entities.User{
				ID: userID, Name: "Test User", Email: "test@example.com",
				Role: entities.UserRoleAdmin, IsActive: true, OfficeID: uuid.New(),
			}
		})

		It("should return valid token response with user info and headers", func() {
			claims := &services.CustomClaims{UserID: userID.String()}

			mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
				Return(claims, nil).Once()
			mockUserSvc.EXPECT().GetByID(mock.Anything, userID).Return(user, nil).Once()

			req, _ := http.NewRequest(http.MethodPost, "/auth/validate", nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)
			r.ServeHTTP(w, req)

			ExpectResponseNotNil(w, http.StatusOK)
			Expect(w.Header().Get("X-User-ID")).To(Equal(userID.String()))
			Expect(w.Header().Get("X-User-Role")).To(Equal(user.Role))
		})

		DescribeTable("should handle errors",
			func(authHeader string, setupMocks func(), expectedCode int, expectedError string) {
				if setupMocks != nil {
					setupMocks()
				}
				req, _ := http.NewRequest(http.MethodPost, "/auth/validate", nil)
				if authHeader != "" {
					req.Header.Set("Authorization", authHeader)
				}
				r.ServeHTTP(w, req)
				ExpectErrorCode(w, expectedCode, expectedError)
			},
			Entry("missing authorization header", "", nil, http.StatusUnauthorized, apperrors2.ErrorCodeInvalidAuthHeader),
			Entry("invalid authorization header format", "InvalidFormat "+accessToken, nil, http.StatusUnauthorized, apperrors2.ErrorCodeInvalidAuthHeader),
			Entry("token validation fails", "Bearer "+accessToken, func() {
				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(nil, apperrors2.NewInvalidAccessToken()).Once()
			}, http.StatusUnauthorized, apperrors2.ErrorCodeInvalidAccessToken),
			Entry("invalid user ID in token", "Bearer "+accessToken, func() {
				claims := &services.CustomClaims{UserID: "invalid-uuid"}
				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(claims, nil).Once()
			}, http.StatusBadRequest, apperrors2.ErrorCodeInvalidUserID),
			Entry("user not found", "Bearer "+accessToken, func() {
				claims := &services.CustomClaims{UserID: userID.String()}
				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(claims, nil).Once()
				mockUserSvc.EXPECT().GetByID(mock.Anything, userID).
					Return(nil, apperrors2.NewUserNotFound()).Once()
			}, http.StatusNotFound, apperrors2.ErrorCodeUserNotFound),
		)
	})
})
