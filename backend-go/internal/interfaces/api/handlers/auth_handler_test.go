package handlers_test

import (
	"encoding/json"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/internal/interfaces/api/dtos"
	"ev-warranty-go/internal/interfaces/api/handlers"
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
		handler      handlers.AuthHandler
		r            *gin.Engine
		w            *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		mockLogger, r, w = SetupMock(GinkgoT())
		mockAuthSvc = mocks.NewAuthService(GinkgoT())
		mockTokenSvc = mocks.NewTokenService(GinkgoT())
		mockUserSvc = mocks.NewUserService(GinkgoT())
		handler = handlers.NewAuthHandler(mockLogger, mockAuthSvc, mockTokenSvc, mockUserSvc)
	})

	Describe("Login", func() {
		var (
			loginRequest = dtos.LoginRequest{Email: "test@example.com", Password: "password123"}
			userID       = uuid.New()
			accessToken  = "access-token-123"
			refreshToken = "refresh-token-456"
			user         *entities.User
		)

		BeforeEach(func() {
			r.POST("/auth/login", handler.Login)
			user = &entities.User{
				ID: userID, Name: "Test User", Email: "test@example.com",
				Role: entities.UserRoleAdmin, IsActive: true, OfficeID: uuid.New(),
			}
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
			Entry("invalid JSON", nil, "invalid json", http.StatusBadRequest, apperrors.ErrorCodeInvalidJsonRequest),
			Entry("auth service login fails", func() {
				mockAuthSvc.EXPECT().Login(mock.Anything, loginRequest.Email, loginRequest.Password).
					Return("", "", apperrors.NewInvalidCredentials()).Once()
			}, loginRequest, http.StatusUnauthorized, apperrors.ErrorCodeInvalidCredentials),
			Entry("token validation fails", func() {
				mockAuthSvc.EXPECT().Login(mock.Anything, loginRequest.Email, loginRequest.Password).
					Return(accessToken, refreshToken, nil).Once()
				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(nil, apperrors.NewInvalidAccessToken()).Once()
			}, loginRequest, http.StatusUnauthorized, apperrors.ErrorCodeInvalidAccessToken),
			Entry("user not found", func() {
				claims := &services.CustomClaims{UserID: userID.String()}
				mockAuthSvc.EXPECT().Login(mock.Anything, loginRequest.Email, loginRequest.Password).
					Return(accessToken, refreshToken, nil).Once()
				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(claims, nil).Once()
				mockUserSvc.EXPECT().GetByID(mock.Anything, userID).
					Return(nil, apperrors.NewUserNotFound()).Once()
			}, loginRequest, http.StatusNotFound, apperrors.ErrorCodeUserNotFound),
		)
	})

	Describe("Logout", func() {
		var refreshToken = "refresh-token-456"

		BeforeEach(func() {
			r.POST("/auth/logout", handler.Logout)
		})

		It("should logout user and return success", func() {
			mockAuthSvc.EXPECT().Logout(mock.Anything, refreshToken).Return(nil).Once()

			req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
			req.AddCookie(&http.Cookie{Name: "refreshToken", Value: refreshToken})
			r.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
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
			Entry("missing refresh token", false, nil, http.StatusNotFound, apperrors.ErrorCodeRefreshTokenNotFound),
			Entry("auth service logout fails", true, func() {
				mockAuthSvc.EXPECT().Logout(mock.Anything, refreshToken).
					Return(apperrors.NewInvalidRefreshToken()).Once()
			}, http.StatusUnauthorized, apperrors.ErrorCodeInvalidRefreshToken),
		)
	})

	Describe("RefreshToken", func() {
		var (
			refreshToken   = "refresh-token-456"
			newAccessToken = "new-access-token-789"
		)

		BeforeEach(func() {
			r.POST("/auth/refresh", handler.RefreshToken)
		})

		It("should return new access token", func() {
			mockTokenSvc.EXPECT().RefreshAccessToken(mock.Anything, refreshToken).
				Return(newAccessToken, nil).Once()

			req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
			req.AddCookie(&http.Cookie{Name: "refreshToken", Value: refreshToken})
			r.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))

			var response dtos.APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())

			tokenResp, ok := response.Data.(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(tokenResp["access_token"]).To(Equal(newAccessToken))
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
			Entry("missing refresh token", false, nil, http.StatusNotFound, apperrors.ErrorCodeRefreshTokenNotFound),
			Entry("token service refresh fails", true, func() {
				mockTokenSvc.EXPECT().RefreshAccessToken(mock.Anything, refreshToken).
					Return("", apperrors.NewInvalidRefreshToken()).Once()
			}, http.StatusUnauthorized, apperrors.ErrorCodeInvalidRefreshToken),
		)
	})

	Describe("ValidateToken", func() {
		var (
			accessToken = "valid-access-token"
			userID      = uuid.New()
			user        *entities.User
		)

		BeforeEach(func() {
			r.POST("/auth/validate", handler.ValidateToken)
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

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("X-User-ID")).To(Equal(userID.String()))
			Expect(w.Header().Get("X-User-Role")).To(Equal(user.Role))

			var response dtos.APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			Expect(err).NotTo(HaveOccurred())

			validationResp, ok := response.Data.(map[string]interface{})
			Expect(ok).To(BeTrue())
			Expect(validationResp["valid"]).To(BeTrue())
			Expect(validationResp["user"]).NotTo(BeNil())
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
			Entry("missing authorization header", "", nil, http.StatusUnauthorized, apperrors.ErrorCodeInvalidAuthHeader),
			Entry("invalid authorization header format", "InvalidFormat "+accessToken, nil, http.StatusUnauthorized, apperrors.ErrorCodeInvalidAuthHeader),
			Entry("token validation fails", "Bearer "+accessToken, func() {
				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(nil, apperrors.NewInvalidAccessToken()).Once()
			}, http.StatusUnauthorized, apperrors.ErrorCodeInvalidAccessToken),
			Entry("invalid user ID in token", "Bearer "+accessToken, func() {
				claims := &services.CustomClaims{UserID: "invalid-uuid"}
				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(claims, nil).Once()
			}, http.StatusBadRequest, apperrors.ErrorCodeInvalidUserID),
			Entry("user not found", "Bearer "+accessToken, func() {
				claims := &services.CustomClaims{UserID: userID.String()}
				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(claims, nil).Once()
				mockUserSvc.EXPECT().GetByID(mock.Anything, userID).
					Return(nil, apperrors.NewUserNotFound()).Once()
			}, http.StatusNotFound, apperrors.ErrorCodeUserNotFound),
		)
	})
})
