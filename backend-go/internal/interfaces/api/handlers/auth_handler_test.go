package handlers_test

import (
	"bytes"
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
			loginRequest dtos.LoginRequest
			userID       uuid.UUID
			accessToken  string
			refreshToken string
			user         *entities.User
		)

		BeforeEach(func() {
			loginRequest = dtos.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			}
			userID = uuid.New()
			accessToken = "access-token-123"
			refreshToken = "refresh-token-456"
			user = &entities.User{
				ID:       userID,
				Name:     "Test User",
				Email:    "test@example.com",
				Role:     entities.UserRoleAdmin,
				IsActive: true,
				OfficeID: uuid.New(),
			}
		})

		Context("when login is successful", func() {
			It("should return access token and user info with refresh token cookie", func() {
				claims := &services.CustomClaims{
					UserID: userID.String(),
				}

				mockAuthSvc.EXPECT().Login(mock.Anything, loginRequest.Email, loginRequest.Password).
					Return(accessToken, refreshToken, nil).Once()
				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(claims, nil).Once()
				mockUserSvc.EXPECT().GetByID(mock.Anything, userID).
					Return(user, nil).Once()

				r.POST("/auth/login", handler.Login)
				body, _ := json.Marshal(loginRequest)
				req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				ExpectCookieRefreshToken(w, refreshToken)

				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())

				loginResp, ok := response.Data.(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(loginResp["token"]).To(Equal(accessToken))
				Expect(loginResp["user"]).NotTo(BeNil())
			})
		})

		Context("when request body is invalid", func() {
			It("should return bad request error", func() {
				r.POST("/auth/login", handler.Login)
				req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte("invalid json")))
				req.Header.Set("Content-Type", "application/json")
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidJsonRequest)
			})
		})

		Context("when auth service login fails", func() {
			It("should return authentication error", func() {
				mockAuthSvc.EXPECT().Login(mock.Anything, loginRequest.Email, loginRequest.Password).
					Return("", "", apperrors.NewInvalidCredentials()).Once()

				r.POST("/auth/login", handler.Login)
				body, _ := json.Marshal(loginRequest)
				req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidCredentials)
			})
		})

		Context("when token validation fails", func() {
			It("should return token validation error", func() {
				mockAuthSvc.EXPECT().Login(mock.Anything, loginRequest.Email, loginRequest.Password).
					Return(accessToken, refreshToken, nil).Once()
				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(nil, apperrors.NewInvalidAccessToken()).Once()

				r.POST("/auth/login", handler.Login)
				body, _ := json.Marshal(loginRequest)
				req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusUnauthorized))
				ExpectErrorCode(w, http.StatusUnauthorized, apperrors.ErrorCodeInvalidAccessToken)
			})
		})

		Context("when user service GetByID fails", func() {
			It("should return user not found error", func() {
				claims := &services.CustomClaims{
					UserID: userID.String(),
				}

				mockAuthSvc.EXPECT().Login(mock.Anything, loginRequest.Email, loginRequest.Password).
					Return(accessToken, refreshToken, nil).Once()
				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(claims, nil).Once()
				mockUserSvc.EXPECT().GetByID(mock.Anything, userID).
					Return(nil, apperrors.NewUserNotFound()).Once()

				r.POST("/auth/login", handler.Login)
				body, _ := json.Marshal(loginRequest)
				req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeUserNotFound)
			})
		})
	})

	Describe("Logout", func() {
		var refreshToken string

		BeforeEach(func() {
			refreshToken = "refresh-token-456"
		})

		Context("when logout is successful", func() {
			It("should logout user and return success", func() {
				mockAuthSvc.EXPECT().Logout(mock.Anything, refreshToken).Return(nil).Once()

				r.POST("/auth/logout", handler.Logout)
				req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
				req.AddCookie(&http.Cookie{
					Name:  "refreshToken",
					Value: refreshToken,
				})
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})

		Context("when refresh token cookie is missing", func() {
			It("should return refresh token not found error", func() {
				r.POST("/auth/logout", handler.Logout)
				req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeRefreshTokenNotFound)
			})
		})

		Context("when auth service logout fails", func() {
			It("should return logout error", func() {
				mockAuthSvc.EXPECT().Logout(mock.Anything, refreshToken).
					Return(apperrors.NewInvalidRefreshToken()).Once()

				r.POST("/auth/logout", handler.Logout)
				req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
				req.AddCookie(&http.Cookie{
					Name:  "refreshToken",
					Value: refreshToken,
				})
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusUnauthorized))
				ExpectErrorCode(w, http.StatusUnauthorized, apperrors.ErrorCodeInvalidRefreshToken)
			})
		})
	})

	Describe("RefreshToken", func() {
		var (
			refreshToken   string
			newAccessToken string
		)

		BeforeEach(func() {
			refreshToken = "refresh-token-456"
			newAccessToken = "new-access-token-789"
		})

		Context("when token refresh is successful", func() {
			It("should return new access token", func() {
				mockTokenSvc.EXPECT().RefreshAccessToken(mock.Anything, refreshToken).
					Return(newAccessToken, nil).Once()

				r.POST("/auth/refresh", handler.RefreshToken)
				req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
				req.AddCookie(&http.Cookie{
					Name:  "refreshToken",
					Value: refreshToken,
				})
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))

				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())

				tokenResp, ok := response.Data.(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(tokenResp["access_token"]).To(Equal(newAccessToken))
			})
		})

		Context("when refresh token cookie is missing", func() {
			It("should return refresh token not found error", func() {
				r.POST("/auth/refresh", handler.RefreshToken)
				req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeRefreshTokenNotFound)
			})
		})

		Context("when token service refresh fails", func() {
			It("should return refresh token error", func() {
				mockTokenSvc.EXPECT().RefreshAccessToken(mock.Anything, refreshToken).
					Return("", apperrors.NewInvalidRefreshToken()).Once()

				r.POST("/auth/refresh", handler.RefreshToken)
				req, _ := http.NewRequest(http.MethodPost, "/auth/refresh", nil)
				req.AddCookie(&http.Cookie{
					Name:  "refreshToken",
					Value: refreshToken,
				})
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusUnauthorized))
				ExpectErrorCode(w, http.StatusUnauthorized, apperrors.ErrorCodeInvalidRefreshToken)
			})
		})
	})

	Describe("ValidateToken", func() {
		var (
			accessToken string
			userID      uuid.UUID
			user        *entities.User
		)

		BeforeEach(func() {
			accessToken = "valid-access-token"
			userID = uuid.New()
			user = &entities.User{
				ID:       userID,
				Name:     "Test User",
				Email:    "test@example.com",
				Role:     entities.UserRoleAdmin,
				IsActive: true,
				OfficeID: uuid.New(),
			}
		})

		Context("when token validation is successful", func() {
			It("should return valid token response with user info and headers", func() {
				claims := &services.CustomClaims{
					UserID: userID.String(),
				}

				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(claims, nil).Once()
				mockUserSvc.EXPECT().GetByID(mock.Anything, userID).
					Return(user, nil).Once()

				r.POST("/auth/validate", handler.ValidateToken)
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
		})

		Context("when authorization header is missing", func() {
			It("should return invalid auth header error", func() {
				r.POST("/auth/validate", handler.ValidateToken)
				req, _ := http.NewRequest(http.MethodPost, "/auth/validate", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusUnauthorized))
				ExpectErrorCode(w, http.StatusUnauthorized, apperrors.ErrorCodeInvalidAuthHeader)
			})
		})

		Context("when authorization header format is invalid", func() {
			It("should return invalid auth header error", func() {
				r.POST("/auth/validate", handler.ValidateToken)
				req, _ := http.NewRequest(http.MethodPost, "/auth/validate", nil)
				req.Header.Set("Authorization", "InvalidFormat "+accessToken)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusUnauthorized))
				ExpectErrorCode(w, http.StatusUnauthorized, apperrors.ErrorCodeInvalidAuthHeader)
			})
		})

		Context("when token validation fails", func() {
			It("should return token validation error", func() {
				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(nil, apperrors.NewInvalidAccessToken()).Once()

				r.POST("/auth/validate", handler.ValidateToken)
				req, _ := http.NewRequest(http.MethodPost, "/auth/validate", nil)
				req.Header.Set("Authorization", "Bearer "+accessToken)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusUnauthorized))
				ExpectErrorCode(w, http.StatusUnauthorized, apperrors.ErrorCodeInvalidAccessToken)
			})
		})

		Context("when user ID in token is invalid", func() {
			It("should return invalid user ID error", func() {
				claims := &services.CustomClaims{
					UserID: "invalid-uuid",
				}

				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(claims, nil).Once()

				r.POST("/auth/validate", handler.ValidateToken)
				req, _ := http.NewRequest(http.MethodPost, "/auth/validate", nil)
				req.Header.Set("Authorization", "Bearer "+accessToken)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUserID)
			})
		})

		Context("when user service GetByID fails", func() {
			It("should return user not found error", func() {
				claims := &services.CustomClaims{
					UserID: userID.String(),
				}

				mockTokenSvc.EXPECT().ValidateAccessToken(mock.Anything, accessToken).
					Return(claims, nil).Once()
				mockUserSvc.EXPECT().GetByID(mock.Anything, userID).
					Return(nil, apperrors.NewUserNotFound()).Once()

				r.POST("/auth/validate", handler.ValidateToken)
				req, _ := http.NewRequest(http.MethodPost, "/auth/validate", nil)
				req.Header.Set("Authorization", "Bearer "+accessToken)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeUserNotFound)
			})
		})
	})
})
