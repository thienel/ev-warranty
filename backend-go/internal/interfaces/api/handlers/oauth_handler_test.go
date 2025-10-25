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
		frontendBaseURL = "http://localhost:3000"
		errorMsg        = "Error while login with Google, please try again!"
	)

	BeforeEach(func() {
		mockLogger, r, w = SetupMock(GinkgoT())
		mockOAuthSvc = mocks.NewOAuthService(GinkgoT())
		mockAuthSvc = mocks.NewAuthService(GinkgoT())
		handler = handlers.NewOAuthHandler(mockLogger, frontendBaseURL, mockOAuthSvc, mockAuthSvc)
	})

	Describe("InitiateOAuth", func() {
		BeforeEach(func() {
			r.GET("/oauth/login", handler.InitiateOAuth)
		})

		It("should redirect to OAuth provider on success", func() {
			authURL := "https://accounts.google.com/o/oauth2/auth?client_id=test"
			mockOAuthSvc.EXPECT().GenerateAuthURL().Return(authURL, nil).Once()

			req, _ := http.NewRequest(http.MethodGet, "/oauth/login", nil)
			r.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusFound))
			Expect(w.Header().Get("Location")).To(Equal(authURL))
		})

		It("should redirect to frontend on error", func() {
			mockOAuthSvc.EXPECT().GenerateAuthURL().Return("", errors.New("failed")).Once()

			req, _ := http.NewRequest(http.MethodGet, "/oauth/login", nil)
			r.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusFound))
			Expect(w.Header().Get("Location")).To(Equal(fmt.Sprintf("%s/login?error=%s", frontendBaseURL, errorMsg)))
		})
	})

	Describe("HandleCallback", func() {
		var (
			validCode    = "test-auth-code"
			validState   = "test-state-token"
			userInfo     *providers.UserInfo
			accessToken  = "access-token-123"
			refreshToken = "refresh-token-456"
		)

		BeforeEach(func() {
			r.GET("/oauth/callback", handler.HandleCallback)
			userInfo = &providers.UserInfo{
				Provider:   "google",
				ProviderID: "123456789",
				Email:      "test@example.com",
				Name:       "Test User",
			}
		})

		It("should handle successful OAuth callback", func() {
			mockOAuthSvc.EXPECT().HandleCallback(mock.Anything, validCode, validState).Return(userInfo, nil).Once()
			mockAuthSvc.EXPECT().HandleOAuthUser(mock.Anything, userInfo).Return(accessToken, refreshToken, nil).Once()

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/oauth/callback?code=%s&state=%s", validCode, validState), nil)
			r.ServeHTTP(w, req)

			Expect(w.Code).To(Equal(http.StatusFound))
			Expect(w.Header().Get("Location")).To(Equal(fmt.Sprintf("%s/callback?token=%s", frontendBaseURL, accessToken)))
			ExpectCookieRefreshToken(w, refreshToken)
		})

		DescribeTable("should redirect to frontend with error",
			func(url string, setupMocks func()) {
				if setupMocks != nil {
					setupMocks()
				}
				req, _ := http.NewRequest(http.MethodGet, url, nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusFound))
				Expect(w.Header().Get("Location")).To(Equal(fmt.Sprintf("%s/login?error=%s", frontendBaseURL, errorMsg)))
			},
			Entry("when state is missing", fmt.Sprintf("/oauth/callback?code=%s", validCode), nil),
			Entry("when code is missing", fmt.Sprintf("/oauth/callback?state=%s", validState), nil),
			Entry("when OAuth provider returns error", "/oauth/callback?error=access_denied&state=test-state", nil),
			Entry("when OAuth service fails", fmt.Sprintf("/oauth/callback?code=%s&state=%s", validCode, validState), func() {
				mockOAuthSvc.EXPECT().HandleCallback(mock.Anything, validCode, validState).
					Return(nil, errors.New("invalid state")).Once()
			}),
			Entry("when auth service fails", fmt.Sprintf("/oauth/callback?code=%s&state=%s", validCode, validState), func() {
				mockOAuthSvc.EXPECT().HandleCallback(mock.Anything, validCode, validState).Return(userInfo, nil).Once()
				mockAuthSvc.EXPECT().HandleOAuthUser(mock.Anything, userInfo).
					Return("", "", errors.New("user creation failed")).Once()
			}),
		)
	})
})
