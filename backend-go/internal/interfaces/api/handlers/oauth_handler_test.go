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

const errPathMessage = "/login?error=Error while login with Google, please try again!"

var _ = Describe("OAuthHandler", func() {
	var (
		mockLogger      *mocks.Logger
		mockOAuthSvc    *mocks.OAuthService
		mockAuthSvc     *mocks.AuthService
		handler         handlers.OAuthHandler
		r               *gin.Engine
		w               *httptest.ResponseRecorder
		frontendBaseURL string
		sampleUserInfo  *providers.UserInfo
	)

	BeforeEach(func() {
		mockLogger, r, w = SetupMock(GinkgoT())
		mockOAuthSvc = mocks.NewOAuthService(GinkgoT())
		mockAuthSvc = mocks.NewAuthService(GinkgoT())
		frontendBaseURL = "http://localhost:3000"
		handler = handlers.NewOAuthHandler(mockLogger, frontendBaseURL, mockOAuthSvc, mockAuthSvc)

		sampleUserInfo = &providers.UserInfo{
			Provider:   "google",
			ProviderID: "123456789",
			Email:      "test@example.com",
			Name:       "Test User",
		}
	})

	Describe("InitiateOAuth", func() {
		BeforeEach(func() {
			r.GET("/oauth/login", handler.InitiateOAuth)
		})

		DescribeTable("should handle different scenarios",
			func(setupMock func(), expectedLocation string) {
				setupMock()
				req, _ := http.NewRequest(http.MethodGet, "/oauth/login", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusFound))
				Expect(w.Header().Get("Location")).To(Equal(expectedLocation))
			},
			Entry("successful auth URL generation",
				func() {
					authURL := "https://accounts.google.com/o/oauth2/auth?client_id=test"
					mockOAuthSvc.EXPECT().GenerateAuthURL().Return(authURL, nil).Once()
				},
				"https://accounts.google.com/o/oauth2/auth?client_id=test"),
			Entry("auth URL generation failure",
				func() {
					mockOAuthSvc.EXPECT().GenerateAuthURL().Return("", errors.New("state generation failed")).Once()
				},
				errPathMessage),
		)
	})

	Describe("HandleCallback", func() {
		var (
			validCode  string
			validState string
		)

		BeforeEach(func() {
			validCode = "test-auth-code"
			validState = "test-state-token"
			r.GET("/oauth/callback", handler.HandleCallback)
		})

		It("should handle successful OAuth callback", func() {
			accessToken := "access-token-123"
			refreshToken := "refresh-token-456"

			mockOAuthSvc.EXPECT().HandleCallback(mock.Anything, validCode, validState).Return(sampleUserInfo, nil).Once()
			mockAuthSvc.EXPECT().HandleOAuthUser(mock.Anything, sampleUserInfo).Return(accessToken, refreshToken, nil).Once()

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/oauth/callback?code=%s&state=%s", validCode, validState), nil)
			r.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusFound))
			expectedLocation := fmt.Sprintf("%s/auth/callback?token=%s", frontendBaseURL, accessToken)
			Expect(w.Header().Get("Location")).To(Equal(expectedLocation))
			ExpectCookieRefreshToken(w, refreshToken)
		})

		DescribeTable("should handle error scenarios",
			func(setupRequest func() *http.Request, setupMock func(), expectedLocation string) {
				if setupMock != nil {
					setupMock()
				}
				req := setupRequest()
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusFound))
				Expect(w.Header().Get("Location")).To(Equal(expectedLocation))
			},
			Entry("missing state parameter",
				func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/oauth/callback?code=%s", validCode), nil)
					return req
				},
				nil,
				errPathMessage),
			Entry("missing code parameter",
				func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/oauth/callback?state=%s", validState), nil)
					return req
				},
				nil,
				errPathMessage),
			Entry("OAuth provider error",
				func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "/oauth/callback?error=access_denied&state=test-state", nil)
					return req
				},
				nil,
				errPathMessage),
			Entry("OAuth service HandleCallback failure",
				func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/oauth/callback?code=%s&state=%s", validCode, validState), nil)
					return req
				},
				func() {
					mockOAuthSvc.EXPECT().HandleCallback(mock.Anything, validCode, validState).
						Return(nil, errors.New("invalid state")).Once()
				},
				errPathMessage),
			Entry("auth service HandleOAuthUser failure",
				func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/oauth/callback?code=%s&state=%s", validCode, validState), nil)
					return req
				},
				func() {
					mockOAuthSvc.EXPECT().HandleCallback(mock.Anything, validCode, validState).Return(sampleUserInfo, nil).Once()
					mockAuthSvc.EXPECT().HandleOAuthUser(mock.Anything, sampleUserInfo).
						Return("", "", errors.New("user creation failed")).Once()
				},
				errPathMessage),
		)
	})
})
