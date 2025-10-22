package handlers_test

import (
	"errors"
	"ev-warranty-go/internal/infrastructure/oauth/providers"
	"ev-warranty-go/internal/interfaces/api/handlers"
	"ev-warranty-go/pkg/mocks"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("OAuthHandler", func() {
	var (
		mockLogger      *mocks.Logger
		mockOAuthSvc    *mocks.OAuthService
		mockAuthSvc     *mocks.AuthService
		handler         handlers.OAuthHandler
		r               *gin.Engine
		w               *httptest.ResponseRecorder
		frontendBaseURL string
	)

	BeforeEach(func() {
		mockLogger, r, w = SetupMock(GinkgoT())
		mockOAuthSvc = mocks.NewOAuthService(GinkgoT())
		mockAuthSvc = mocks.NewAuthService(GinkgoT())
		frontendBaseURL = "http://localhost:3000"
		handler = handlers.NewOAuthHandler(mockLogger, frontendBaseURL, mockOAuthSvc, mockAuthSvc)
	})

	Describe("InitiateOAuth", func() {
		Context("when auth URL is generated successfully", func() {
			It("should redirect to the OAuth provider", func() {
				authURL := "https://accounts.google.com/o/oauth2/auth?client_id=test"
				mockOAuthSvc.EXPECT().GenerateAuthURL().Return(authURL, nil).Once()

				r.GET("/oauth/login", handler.InitiateOAuth)
				req, _ := http.NewRequest(http.MethodGet, "/oauth/login", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusFound))
				Expect(w.Header().Get("Location")).To(Equal(authURL))
			})
		})

		Context("when auth URL generation fails", func() {
			It("should redirect to frontend login with error", func() {
				mockOAuthSvc.EXPECT().GenerateAuthURL().Return("", errors.New("state generation failed")).Once()

				r.GET("/oauth/login", handler.InitiateOAuth)
				req, _ := http.NewRequest(http.MethodGet, "/oauth/login", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusFound))
				expectedLocation := fmt.Sprintf("%s/login?error=%s", frontendBaseURL, "Error while login with Google, please try again!")
				Expect(w.Header().Get("Location")).To(Equal(expectedLocation))
			})
		})
	})

	Describe("HandleCallback", func() {
		var (
			validCode  string
			validState string
		)

		BeforeEach(func() {
			validCode = "test-auth-code"
			validState = "test-state-token"
		})

		Context("when OAuth callback is successful", func() {
			It("should set refresh token cookie and redirect with access token", func() {
				userInfo := &providers.UserInfo{
					Provider:   "google",
					ProviderID: "123456789",
					Email:      "test@example.com",
					Name:       "Test User",
				}
				accessToken := "access-token-123"
				refreshToken := "refresh-token-456"

				mockOAuthSvc.EXPECT().HandleCallback(mock.Anything, validCode, validState).Return(userInfo, nil).Once()
				mockAuthSvc.EXPECT().HandleOAuthUser(mock.Anything, userInfo).Return(accessToken, refreshToken, nil).Once()

				r.GET("/oauth/callback", handler.HandleCallback)
				req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/oauth/callback?code=%s&state=%s", validCode, validState), nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusFound))
				expectedLocation := fmt.Sprintf("%s/auth/callback?token=%s", frontendBaseURL, accessToken)
				Expect(w.Header().Get("Location")).To(Equal(expectedLocation))
				ExpectCookieRefreshToken(w, refreshToken)
			})
		})

		Context("when state is missing", func() {
			It("should redirect to frontend login with error", func() {
				r.GET("/oauth/callback", handler.HandleCallback)
				req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/oauth/callback?code=%s", validCode), nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusFound))
				expectedLocation := fmt.Sprintf("%s/login?error=%s", frontendBaseURL, "Error while login with Google, please try again!")
				Expect(w.Header().Get("Location")).To(Equal(expectedLocation))
			})
		})

		Context("when code is missing", func() {
			It("should redirect to frontend login with error", func() {
				r.GET("/oauth/callback", handler.HandleCallback)
				req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/oauth/callback?state=%s", validState), nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusFound))
				expectedLocation := fmt.Sprintf("%s/login?error=%s", frontendBaseURL, "Error while login with Google, please try again!")
				Expect(w.Header().Get("Location")).To(Equal(expectedLocation))
			})
		})

		Context("when OAuth provider returns error", func() {
			It("should redirect to frontend login with error", func() {
				r.GET("/oauth/callback", handler.HandleCallback)
				req, _ := http.NewRequest(http.MethodGet, "/oauth/callback?error=access_denied&state=test-state", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusFound))
				expectedLocation := fmt.Sprintf("%s/login?error=%s", frontendBaseURL, "Error while login with Google, please try again!")
				Expect(w.Header().Get("Location")).To(Equal(expectedLocation))
			})
		})

		Context("when OAuth service HandleCallback fails", func() {
			It("should redirect to frontend login with error", func() {
				mockOAuthSvc.EXPECT().HandleCallback(mock.Anything, validCode, validState).
					Return(nil, errors.New("invalid state")).Once()

				r.GET("/oauth/callback", handler.HandleCallback)
				req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/oauth/callback?code=%s&state=%s", validCode, validState), nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusFound))
				expectedLocation := fmt.Sprintf("%s/login?error=%s", frontendBaseURL, "Error while login with Google, please try again!")
				Expect(w.Header().Get("Location")).To(Equal(expectedLocation))
			})
		})

		Context("when auth service HandleOAuthUser fails", func() {
			It("should redirect to frontend login with error", func() {
				userInfo := &providers.UserInfo{
					Provider:   "google",
					ProviderID: "123456789",
					Email:      "test@example.com",
					Name:       "Test User",
				}

				mockOAuthSvc.EXPECT().HandleCallback(mock.Anything, validCode, validState).Return(userInfo, nil).Once()
				mockAuthSvc.EXPECT().HandleOAuthUser(mock.Anything, userInfo).
					Return("", "", errors.New("user creation failed")).Once()

				r.GET("/oauth/callback", handler.HandleCallback)
				req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/oauth/callback?code=%s&state=%s", validCode, validState), nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusFound))
				expectedLocation := fmt.Sprintf("%s/login?error=%s", frontendBaseURL, "Error while login with Google, please try again!")
				Expect(w.Header().Get("Location")).To(Equal(expectedLocation))
			})
		})
	})
})
