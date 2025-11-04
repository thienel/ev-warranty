package service_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"ev-warranty-go/pkg/apperror"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"ev-warranty-go/internal/application/service"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/pkg/mocks"
)

var _ = Describe(
	"TokenService", func() {
		var (
			mockRefreshTokenRepo *mocks.RefreshTokenRepository
			tokenService         service.TokenService
			ctx                  context.Context
			privateKey           *rsa.PrivateKey
			publicKey            *rsa.PublicKey
			accessTTL            time.Duration
			refreshTTL           time.Duration
		)

		BeforeEach(
			func() {
				mockRefreshTokenRepo = mocks.NewRefreshTokenRepository(GinkgoT())
				ctx = context.Background()

				var err error
				privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
				Expect(err).NotTo(HaveOccurred())
				publicKey = &privateKey.PublicKey

				accessTTL = 15 * time.Minute
				refreshTTL = 7 * 24 * time.Hour

				tokenService = service.NewTokenService(
					mockRefreshTokenRepo, accessTTL, refreshTTL, privateKey, publicKey)
			})

		Describe("GenerateAccessToken", func() {
			var userID uuid.UUID
			BeforeEach(func() { userID = uuid.New() })

			Context("when token is generated successfully", func() {
				It("should return a valid JWT token", func() {
					token, err := tokenService.GenerateAccessToken(userID)

					Expect(err).NotTo(HaveOccurred())
					Expect(token).NotTo(BeEmpty())

					parsedToken, parseErr := jwt.ParseWithClaims(
						token, &service.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
							return publicKey, nil
						},
					)
					Expect(parseErr).NotTo(HaveOccurred())
					Expect(parsedToken.Valid).To(BeTrue())

					claims, ok := parsedToken.Claims.(*service.CustomClaims)
					Expect(ok).To(BeTrue())
					Expect(claims.UserID).To(Equal(userID.String()))
					Expect(claims.Subject).To(Equal(userID.String()))
					Expect(claims.Issuer).To(Equal("auth-service"))
				})
			})
		})

		Describe("GenerateRefreshToken", func() {
			var userID uuid.UUID

			BeforeEach(func() { userID = uuid.New() })

			Context("when refresh token is generated successfully", func() {
				It("should return a token string and create it in repository", func() {
					mockRefreshTokenRepo.EXPECT().Create(ctx, MatchRefreshToken(userID)).Return(nil).Once()

					token, err := tokenService.GenerateRefreshToken(ctx, userID)

					Expect(err).NotTo(HaveOccurred())
					Expect(token).NotTo(BeEmpty())
				})
			})

			Context("when repository create fails", func() {
				It("should return FailedGenerateRefreshToken error", func() {
					dbErr := apperror.ErrDBOperation
					mockRefreshTokenRepo.EXPECT().Create(
						ctx, MatchRefreshToken(userID),
					).Return(dbErr).Once()

					token, err := tokenService.GenerateRefreshToken(ctx, userID)

					Expect(token).To(BeEmpty())
					ExpectAppError(err, apperror.ErrFailedGenerateRefreshToken.ErrorCode)
				})
			})
		})

		Describe("ValidateAccessToken", func() {
			var userID uuid.UUID

			BeforeEach(func() { userID = uuid.New() })

			Context("when token is valid", func() {
				It("should return claims", func() {
					token, err := tokenService.GenerateAccessToken(userID)
					Expect(err).NotTo(HaveOccurred())

					claims, err := tokenService.ValidateAccessToken(ctx, token)

					Expect(err).NotTo(HaveOccurred())
					Expect(claims).NotTo(BeNil())
					Expect(claims.UserID).To(Equal(userID.String()))
				})
			})

			Context("when token is malformed", func() {
				It("should return InvalidAccessToken error", func() {
					claims, err := tokenService.ValidateAccessToken(ctx, "invalid.token.here")

					Expect(claims).To(BeNil())
					ExpectAppError(err, apperror.ErrInvalidAccessToken.ErrorCode)
				})
			})

			Context("when token is expired", func() {
				It("should return ExpiredAccessToken error", func() {
					shortService := service.NewTokenService(
						mockRefreshTokenRepo, -1*time.Hour, refreshTTL, privateKey, publicKey,
					)
					token, err := shortService.GenerateAccessToken(userID)
					Expect(err).NotTo(HaveOccurred())

					claims, err := tokenService.ValidateAccessToken(ctx, token)

					Expect(claims).To(BeNil())
					ExpectAppError(err, apperror.ErrExpiredAccessToken.ErrorCode)
				})
			})

			Context("when token is not yet valid", func() {
				It("should return InvalidAccessToken error", func() {
					claims := service.CustomClaims{
						UserID: userID.String(),
						RegisteredClaims: jwt.RegisteredClaims{
							ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
							IssuedAt:  jwt.NewNumericDate(time.Now()),
							NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour)),
							Subject:   userID.String(),
						},
					}
					token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
					tokenString, _ := token.SignedString(privateKey)

					validatedClaims, err := tokenService.ValidateAccessToken(ctx, tokenString)

					Expect(validatedClaims).To(BeNil())
					ExpectAppError(err, apperror.ErrInvalidAccessToken.ErrorCode)
				})
			})

			Context("when token has wrong signing method", func() {
				It("should return InvalidAccessToken error", func() {
					claims := service.CustomClaims{
						UserID: userID.String(),
						RegisteredClaims: jwt.RegisteredClaims{
							ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
							IssuedAt:  jwt.NewNumericDate(time.Now()),
							Subject:   userID.String(),
						},
					}
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
					tokenString, _ := token.SignedString([]byte("secret"))

					validatedClaims, err := tokenService.ValidateAccessToken(ctx, tokenString)

					Expect(err).To(HaveOccurred())
					Expect(validatedClaims).To(BeNil())
				})
			})

			Context("when token has empty UserID", func() {
				It("should return InvalidAccessToken error", func() {
					claims := service.CustomClaims{
						UserID: "",
						RegisteredClaims: jwt.RegisteredClaims{
							ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
							IssuedAt:  jwt.NewNumericDate(time.Now()),
							Subject:   userID.String(),
						},
					}
					token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
					tokenString, _ := token.SignedString(privateKey)

					validatedClaims, err := tokenService.ValidateAccessToken(ctx, tokenString)

					Expect(validatedClaims).To(BeNil())
					ExpectAppError(err, apperror.ErrInvalidAccessToken.ErrorCode)
				})
			})
		})

		Describe("ValidateRefreshToken", func() {
			var (
				userID       uuid.UUID
				refreshToken string
			)

			BeforeEach(func() { userID = uuid.New(); refreshToken = "test_refresh_token" })

			Context("when refresh token is valid", func() {
				It("should return the refresh token entity", func() {
					expectedToken := &entity.RefreshToken{
						ID:        uuid.New(),
						UserID:    userID,
						Token:     "hashed_token",
						ExpiresAt: time.Now().Add(24 * time.Hour),
						IsRevoked: false,
					}

					mockRefreshTokenRepo.EXPECT().Find(ctx, MatchHashedToken()).Return(
						expectedToken, nil,
					).Once()

					token, err := tokenService.ValidateRefreshToken(ctx, refreshToken)

					Expect(err).NotTo(HaveOccurred())
					Expect(token).NotTo(BeNil())
					Expect(token.UserID).To(Equal(userID))
				})
			})

			Context("when refresh token is not found", func() {
				It("should return InvalidRefreshToken error", func() {
					dbErr := apperror.ErrNotFoundError
					mockRefreshTokenRepo.EXPECT().Find(ctx, MatchHashedToken()).Return(nil, dbErr).Once()

					token, err := tokenService.ValidateRefreshToken(ctx, refreshToken)

					Expect(token).To(BeNil())
					ExpectAppError(err, apperror.ErrInvalidRefreshToken.ErrorCode)
				})
			})

			Context("when refresh token is expired", func() {
				It("should return ExpiredRefreshToken error", func() {
					expiredToken := &entity.RefreshToken{
						ID:        uuid.New(),
						UserID:    userID,
						Token:     "hashed_token",
						ExpiresAt: time.Now().Add(-24 * time.Hour),
						IsRevoked: false,
					}

					mockRefreshTokenRepo.EXPECT().Find(ctx, MatchHashedToken()).Return(expiredToken,
						nil).Once()

					token, err := tokenService.ValidateRefreshToken(ctx, refreshToken)

					Expect(token).To(BeNil())
					ExpectAppError(err, apperror.ErrExpiredRefreshToken.ErrorCode)
				})
			})

			Context("when refresh token is revoked", func() {
				It("should return RevokedRefreshToken error", func() {
					revokedToken := &entity.RefreshToken{
						ID:        uuid.New(),
						UserID:    userID,
						Token:     "hashed_token",
						ExpiresAt: time.Now().Add(24 * time.Hour),
						IsRevoked: true,
					}

					mockRefreshTokenRepo.EXPECT().Find(ctx, MatchHashedToken()).Return(revokedToken,
						nil).Once()

					token, err := tokenService.ValidateRefreshToken(ctx, refreshToken)

					Expect(token).To(BeNil())
					ExpectAppError(err, apperror.ErrRevokedRefreshToken.ErrorCode)
				})
			})
		})

		Describe("RevokeRefreshToken", func() {
			var refreshToken string
			BeforeEach(func() { refreshToken = "test_refresh_token" })

			Context("when token is revoked successfully", func() {
				It("should return nil error", func() {
					mockRefreshTokenRepo.EXPECT().Revoke(ctx, MatchHashedToken()).Return(nil).Once()

					err := tokenService.RevokeRefreshToken(ctx, refreshToken)

					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when repository revoke fails", func() {
				It("should return the error", func() {
					dbErr := apperror.ErrDBOperation
					mockRefreshTokenRepo.EXPECT().Revoke(ctx, MatchHashedToken()).Return(dbErr).Once()

					err := tokenService.RevokeRefreshToken(ctx, refreshToken)

					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(dbErr))
				})
			})
		})

		Describe("RefreshAccessToken", func() {
			var (
				userID       uuid.UUID
				refreshToken string
			)

			BeforeEach(func() { userID = uuid.New(); refreshToken = "test_refresh_token" })

			Context("when refresh is successful", func() {
				It("should return a new access token", func() {
					validToken := &entity.RefreshToken{
						ID:        uuid.New(),
						UserID:    userID,
						Token:     "hashed_token",
						ExpiresAt: time.Now().Add(24 * time.Hour),
						IsRevoked: false,
					}

					mockRefreshTokenRepo.EXPECT().Find(ctx, MatchHashedToken()).Return(
						validToken, nil,
					).Once()

					accessToken, err := tokenService.RefreshAccessToken(ctx, refreshToken)

					Expect(err).NotTo(HaveOccurred())
					Expect(accessToken).NotTo(BeEmpty())

					claims, err := tokenService.ValidateAccessToken(ctx, accessToken)
					Expect(err).NotTo(HaveOccurred())
					Expect(claims.UserID).To(Equal(userID.String()))
				})
			})

			Context("when refresh token is invalid", func() {
				It("should return InvalidRefreshToken error", func() {
					dbErr := apperror.ErrNotFoundError
					mockRefreshTokenRepo.EXPECT().Find(ctx, MatchHashedToken()).Return(nil, dbErr).Once()

					accessToken, err := tokenService.RefreshAccessToken(ctx, refreshToken)

					Expect(accessToken).To(BeEmpty())
					ExpectAppError(err, apperror.ErrInvalidRefreshToken.ErrorCode)
				})
			})
			Context("when refresh token is expired", func() {
				It("should return ExpiredRefreshToken error", func() {
					expiredToken := &entity.RefreshToken{ID: uuid.New(), UserID: userID, Token: "hashed_token",
						ExpiresAt: time.Now().Add(-24 * time.Hour), IsRevoked: false,
					}

					mockRefreshTokenRepo.EXPECT().Find(ctx, MatchHashedToken()).Return(expiredToken, nil).Once()
					accessToken, err := tokenService.RefreshAccessToken(ctx, refreshToken)
					Expect(accessToken).To(BeEmpty())
					ExpectAppError(err, apperror.ErrExpiredRefreshToken.ErrorCode)
				})
			})
		})
	})

func MatchRefreshToken(userID uuid.UUID) interface{} {
	return mock.MatchedBy(
		func(rt *entity.RefreshToken) bool {
			return rt.UserID == userID
		},
	)
}

func MatchHashedToken() interface{} {
	return mock.AnythingOfType("string")
}
