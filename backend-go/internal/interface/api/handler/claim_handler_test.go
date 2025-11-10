package handler_test

import (
	"context"
	"errors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/service"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/internal/interface/api/dto"
	"ev-warranty-go/internal/interface/api/handler"
	"ev-warranty-go/pkg/apperror"
	"ev-warranty-go/pkg/mocks"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("ClaimHandler", func() {
	var (
		mockLogger     *mocks.Logger
		mockTxManager  *mocks.TxManager
		mockService    *mocks.ClaimService
		mockTx         *mocks.Tx
		claimHandler   handler.ClaimHandler
		r              *gin.Engine
		w              *httptest.ResponseRecorder
		userID         uuid.UUID
		claimID        uuid.UUID
		sampleClaim    *entity.Claim
		validCreateReq dto.CreateClaimRequest
		validUpdateReq dto.UpdateClaimRequest
	)

	setupRoute := func(method, path string, role string, handlerFunc gin.HandlerFunc) {
		r.Handle(method, path, func(c *gin.Context) {
			if role != "" {
				SetHeaderRole(c, role)
				SetHeaderID(c, userID)
			}
			SetContentTypeJSON(c)
			handlerFunc(c)
		})
	}

	setupTxMock := func(serviceMockFn func()) {
		mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
			Run(func(ctx context.Context, fn func(application.Tx) error) {
				serviceMockFn()
				_ = fn(mockTx)
			}).Return(nil).Once()
	}

	setupTxMockWithError := func(serviceMockFn func(), expectedError error) {
		mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
			Run(func(ctx context.Context, fn func(application.Tx) error) {
				serviceMockFn()
				_ = fn(mockTx)
			}).Return(expectedError).Once()
	}

	BeforeEach(func() {
		mockLogger, r, w = SetupMock(GinkgoT())
		mockTxManager = mocks.NewTxManager(GinkgoT())
		mockService = mocks.NewClaimService(GinkgoT())
		mockTx = mocks.NewTx(GinkgoT())
		claimHandler = handler.NewClaimHandler(mockLogger, mockTxManager, mockService)

		userID = uuid.New()
		claimID = uuid.New()
		vehicleID := uuid.New()
		kilometers := 1
		technicianID := uuid.New()
		customerID := uuid.New()

		validCreateReq = dto.CreateClaimRequest{
			VehicleID:    vehicleID,
			CustomerID:   customerID,
			TechnicianID: technicianID,
			Description:  "Test claim description for warranty issue",
		}

		validUpdateReq = dto.UpdateClaimRequest{
			Description: "Updated claim description for warranty issue",
		}

		sampleClaim = entity.NewClaim(vehicleID, customerID, kilometers, validCreateReq.Description, userID, technicianID)
		sampleClaim.ID = claimID
	})

	Describe("GetByID", func() {
		BeforeEach(func() {
			r.GET("/claims/:id", claimHandler.GetByID)
		})

		It("should retrieve claim successfully", func() {
			mockService.EXPECT().GetByID(mock.Anything, claimID).Return(sampleClaim, nil).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String(), nil)
			r.ServeHTTP(w, req)

			ExpectResponseNotNil(w, http.StatusOK)
		})

		It("should handle invalid UUID", func() {
			req, _ := http.NewRequest(http.MethodGet, "/claims/invalid-uuid", nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidParams.ErrorCode)
		})

		It("should handle claim not found", func() {
			mockService.EXPECT().GetByID(mock.Anything, claimID).
				Return(nil, apperror.ErrNotFoundError).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String(), nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusNotFound, apperror.ErrNotFoundError.ErrorCode)
		})

		It("should handle service errors", func() {
			mockService.EXPECT().GetByID(mock.Anything, claimID).
				Return(nil, errors.New("database error")).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String(), nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusInternalServerError, apperror.ErrInternalServerError.ErrorCode)
		})
	})

	Describe("GetAll", func() {
		BeforeEach(func() {
			r.GET("/claims", claimHandler.GetAll)
		})

		It("should retrieve all claims successfully", func() {
			expectedClaims := []*entity.Claim{sampleClaim}
			mockService.EXPECT().GetAll(mock.Anything).Return(expectedClaims, nil).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims", nil)
			r.ServeHTTP(w, req)

			ExpectResponseNotNil(w, http.StatusOK)
		})

		It("should handle empty claims list", func() {
			mockService.EXPECT().GetAll(mock.Anything).Return([]*entity.Claim{}, nil).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims", nil)
			r.ServeHTTP(w, req)

			ExpectResponseNotNil(w, http.StatusOK)
		})

		It("should handle service errors", func() {
			mockService.EXPECT().GetAll(mock.Anything).
				Return(nil, errors.New("database error")).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims", nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusInternalServerError, apperror.ErrInternalServerError.ErrorCode)
		})
	})

	Describe("Create", func() {
		Context("when authorized as SC_TECHNICIAN", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims", entity.UserRoleScTechnician, claimHandler.Create)
			})

			It("should create claim successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().Create(mockTx,
						mock.MatchedBy(func(cmd *service.CreateClaimCommand) bool {
							return cmd.VehicleID == validCreateReq.VehicleID &&
								cmd.CustomerID == validCreateReq.CustomerID &&
								cmd.StaffID == userID &&
								cmd.Description == validCreateReq.Description
						})).Return(sampleClaim, nil).Once()
				})

				SendRequest(r, http.MethodPost, "/claims", w, validCreateReq)
				ExpectResponseNotNil(w, http.StatusCreated)
			})

			It("should handle invalid JSON", func() {
				SendRequest(r, http.MethodPost, "/claims", w, "invalid json")
				ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidJsonRequest.ErrorCode)
			})

			DescribeTable("should handle validation errors",
				func(modifyReq func(*dto.CreateClaimRequest), expectedError string) {
					req := validCreateReq
					modifyReq(&req)

					SendRequest(r, http.MethodPost, "/claims", w, req)
					ExpectErrorCode(w, http.StatusBadRequest, expectedError)
				},
				Entry("empty description", func(req *dto.CreateClaimRequest) {
					req.Description = ""
				}, apperror.ErrInvalidJsonRequest),
				Entry("description too short", func(req *dto.CreateClaimRequest) {
					req.Description = "short"
				}, apperror.ErrInvalidJsonRequest),
			)

			It("should handle service errors", func() {
				dbError := errors.New("database error")
				setupTxMockWithError(func() {
					mockService.EXPECT().Create(mockTx, mock.Anything).
						Return(nil, dbError).Once()
				}, dbError)

				SendRequest(r, http.MethodPost, "/claims", w, validCreateReq)
				ExpectErrorCode(w, http.StatusInternalServerError, apperror.ErrInternalServerError.ErrorCode)
			})

			It("should handle transaction errors", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Return(errors.New("transaction error")).Once()

				SendRequest(r, http.MethodPost, "/claims", w, validCreateReq)
				ExpectErrorCode(w, http.StatusInternalServerError, apperror.ErrInternalServerError.ErrorCode)
			})
		})

		Context("when authorized as SC_STAFF", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims", entity.UserRoleScStaff, claimHandler.Create)
			})

			It("should create claim successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().Create(mockTx, mock.Anything).Return(sampleClaim, nil).Once()
				})

				SendRequest(r, http.MethodPost, "/claims", w, validCreateReq)
				ExpectResponseNotNil(w, http.StatusCreated)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("POST", "/claims", entity.UserRoleEvmStaff, claimHandler.Create)
			SendRequest(r, http.MethodPost, "/claims", w, validCreateReq)
			ExpectErrorCode(w, http.StatusForbidden, apperror.ErrUnauthorizedRole.ErrorCode)
		})

		It("should handle missing user ID header", func() {
			r.POST("/claims", func(c *gin.Context) {
				SetHeaderRole(c, entity.UserRoleScTechnician)
				SetContentTypeJSON(c)
				claimHandler.Create(c)
			})

			SendRequest(r, http.MethodPost, "/claims", w, validCreateReq)
			ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrMissingUserID.ErrorCode)
		})

		It("should handle invalid user ID header", func() {
			r.POST("/claims", func(c *gin.Context) {
				SetHeaderRole(c, entity.UserRoleScTechnician)
				c.Request.Header.Set("X-User-ID", "invalid-uuid")
				SetContentTypeJSON(c)
				claimHandler.Create(c)
			})

			SendRequest(r, http.MethodPost, "/claims", w, validCreateReq)
			ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidUserID.ErrorCode)
		})
	})

	Describe("Update", func() {
		Context("when authorized as SC_STAFF", func() {
			BeforeEach(func() {
				setupRoute("PUT", "/claims/:id", entity.UserRoleScStaff, claimHandler.Update)
			})

			It("should update claim successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().Update(mockTx, claimID,
						mock.MatchedBy(func(cmd *service.UpdateClaimCommand) bool {
							return cmd.Description == validUpdateReq.Description
						})).Return(nil).Once()
				})

				SendRequest(r, http.MethodPut, "/claims/"+claimID.String(), w, validUpdateReq)
				ExpectResponseNotNil(w, http.StatusNoContent)
			})

			It("should handle invalid UUID", func() {
				SendRequest(r, http.MethodPut, "/claims/invalid-uuid", w, validUpdateReq)
				ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidParams.ErrorCode)
			})

			It("should handle invalid JSON", func() {
				SendRequest(r, http.MethodPut, "/claims/"+claimID.String(), w, "invalid json")
				ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidJsonRequest.ErrorCode)
			})

			It("should handle service errors", func() {
				notFoundError := apperror.ErrNotFoundError
				setupTxMockWithError(func() {
					mockService.EXPECT().Update(mockTx, claimID, mock.Anything).
						Return(notFoundError).Once()
				}, notFoundError)

				SendRequest(r, http.MethodPut, "/claims/"+claimID.String(), w, validUpdateReq)
				ExpectErrorCode(w, http.StatusNotFound, apperror.ErrNotFoundError.ErrorCode)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("PUT", "/claims/:id", entity.UserRoleScTechnician, claimHandler.Update)
			SendRequest(r, http.MethodPut, "/claims/"+claimID.String(), w, validUpdateReq)
			ExpectErrorCode(w, http.StatusForbidden, apperror.ErrUnauthorizedRole.ErrorCode)
		})
	})

	Describe("Delete", func() {
		Context("when authorized as SC_STAFF", func() {
			BeforeEach(func() {
				setupRoute("DELETE", "/claims/:id", entity.UserRoleScStaff, claimHandler.Delete)
			})

			It("should perform hard delete successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().HardDelete(mockTx, claimID).Return(nil).Once()
				})

				SendRequest(r, http.MethodDelete, "/claims/"+claimID.String(), w, nil)
				ExpectResponseNotNil(w, http.StatusNoContent)
			})

			It("should handle invalid UUID", func() {
				SendRequest(r, http.MethodDelete, "/claims/invalid-uuid", w, nil)
				ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidParams.ErrorCode)
			})

			It("should handle service errors", func() {
				notFoundError := apperror.ErrNotFoundError
				setupTxMockWithError(func() {
					mockService.EXPECT().HardDelete(mockTx, claimID).
						Return(notFoundError).Once()
				}, notFoundError)

				SendRequest(r, http.MethodDelete, "/claims/"+claimID.String(), w, nil)
				ExpectErrorCode(w, http.StatusNotFound, apperror.ErrNotFoundError.ErrorCode)
			})
		})

		Context("when authorized as EVM_STAFF", func() {
			BeforeEach(func() {
				setupRoute("DELETE", "/claims/:id", entity.UserRoleEvmStaff, claimHandler.Delete)
			})

			It("should perform soft delete successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().SoftDelete(mockTx, claimID).Return(nil).Once()
				})

				SendRequest(r, http.MethodDelete, "/claims/"+claimID.String(), w, nil)
				ExpectResponseNotNil(w, http.StatusNoContent)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("DELETE", "/claims/:id", entity.UserRoleScTechnician, claimHandler.Delete)
			SendRequest(r, http.MethodDelete, "/claims/"+claimID.String(), w, nil)
			ExpectErrorCode(w, http.StatusForbidden, apperror.ErrUnauthorizedRole.ErrorCode)
		})
	})

	Describe("Submit", func() {
		Context("when authorized as SC_STAFF", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims/:id/submit", entity.UserRoleScStaff, claimHandler.Submit)
			})

			It("should submit claim successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().Submit(mockTx, claimID, userID).Return(nil).Once()
				})

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/submit", w, nil)
				ExpectResponseNotNil(w, http.StatusNoContent)
			})

			It("should handle invalid UUID", func() {
				SendRequest(r, http.MethodPost, "/claims/invalid-uuid/submit", w, nil)
				ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidParams.ErrorCode)
			})

			It("should handle service errors", func() {
				invalidActionError := apperror.ErrInvalidClaimAction
				setupTxMockWithError(func() {
					mockService.EXPECT().Submit(mockTx, claimID, userID).
						Return(invalidActionError).Once()
				}, invalidActionError)

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/submit", w, nil)
				ExpectErrorCode(w, http.StatusConflict, apperror.ErrInvalidClaimAction.ErrorCode)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("POST", "/claims/:id/submit", entity.UserRoleScTechnician, claimHandler.Submit)
			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/submit", w, nil)
			ExpectErrorCode(w, http.StatusForbidden, apperror.ErrUnauthorizedRole.ErrorCode)
		})

		It("should handle missing user ID header", func() {
			r.POST("/claims/:id/submit", func(c *gin.Context) {
				SetHeaderRole(c, entity.UserRoleScStaff)
				SetContentTypeJSON(c)
				claimHandler.Submit(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/submit", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrMissingUserID.ErrorCode)
		})

		It("should handle invalid user ID header", func() {
			r.POST("/claims/:id/submit", func(c *gin.Context) {
				SetHeaderRole(c, entity.UserRoleScStaff)
				c.Request.Header.Set("X-User-ID", "invalid-uuid")
				SetContentTypeJSON(c)
				claimHandler.Submit(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/submit", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidUserID.ErrorCode)
		})
	})

	Describe("Review", func() {
		Context("when authorized as EVM_STAFF", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims/:id/review", entity.UserRoleEvmStaff, claimHandler.Review)
			})

			It("should start review successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().UpdateStatus(mockTx, claimID, entity.ClaimStatusReviewing,
						userID).Return(nil).Once()
				})

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/review", w, nil)
				ExpectResponseNotNil(w, http.StatusNoContent)
			})

			It("should handle invalid UUID", func() {
				SendRequest(r, http.MethodPost, "/claims/invalid-uuid/review", w, nil)
				ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidParams.ErrorCode)
			})

			It("should handle service errors", func() {
				invalidActionError := apperror.ErrInvalidClaimAction
				setupTxMockWithError(func() {
					mockService.EXPECT().UpdateStatus(mockTx, claimID, entity.ClaimStatusReviewing, userID).
						Return(invalidActionError).Once()
				}, invalidActionError)

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/review", w, nil)
				ExpectErrorCode(w, http.StatusConflict, apperror.ErrInvalidClaimAction.ErrorCode)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("POST", "/claims/:id/review", entity.UserRoleScStaff, claimHandler.Review)
			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/review", w, nil)
			ExpectErrorCode(w, http.StatusForbidden, apperror.ErrUnauthorizedRole.ErrorCode)
		})

		It("should handle missing user ID header", func() {
			r.POST("/claims/:id/review", func(c *gin.Context) {
				SetHeaderRole(c, entity.UserRoleEvmStaff)
				SetContentTypeJSON(c)
				claimHandler.Review(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/review", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrMissingUserID.ErrorCode)
		})

		It("should handle invalid user ID header", func() {
			r.POST("/claims/:id/review", func(c *gin.Context) {
				SetHeaderRole(c, entity.UserRoleEvmStaff)
				c.Request.Header.Set("X-User-ID", "invalid-uuid")
				SetContentTypeJSON(c)
				claimHandler.Review(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/review", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidUserID.ErrorCode)
		})
	})

	Describe("Cancel", func() {
		Context("when authorized as SC_STAFF", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims/:id/cancel", entity.UserRoleScStaff, claimHandler.Cancel)
			})

			It("should cancel claim successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().UpdateStatus(mockTx, claimID, entity.ClaimStatusCancelled,
						userID).Return(nil).Once()
				})

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/cancel", w, nil)
				ExpectResponseNotNil(w, http.StatusNoContent)
			})

			It("should handle invalid UUID", func() {
				SendRequest(r, http.MethodPost, "/claims/invalid-uuid/cancel", w, nil)
				ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidParams.ErrorCode)
			})

			It("should handle service errors", func() {
				invalidActionError := apperror.ErrInvalidClaimAction
				setupTxMockWithError(func() {
					mockService.EXPECT().UpdateStatus(mockTx, claimID, entity.ClaimStatusCancelled, userID).
						Return(invalidActionError).Once()
				}, invalidActionError)

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/cancel", w, nil)
				ExpectErrorCode(w, http.StatusConflict, apperror.ErrInvalidClaimAction.ErrorCode)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("POST", "/claims/:id/cancel", entity.UserRoleEvmStaff, claimHandler.Cancel)
			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/cancel", w, nil)
			ExpectErrorCode(w, http.StatusForbidden, apperror.ErrUnauthorizedRole.ErrorCode)
		})

		It("should handle missing user ID header", func() {
			r.POST("/claims/:id/cancel", func(c *gin.Context) {
				SetHeaderRole(c, entity.UserRoleScStaff)
				SetContentTypeJSON(c)
				claimHandler.Cancel(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/cancel", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrMissingUserID.ErrorCode)
		})

		It("should handle invalid user ID header", func() {
			r.POST("/claims/:id/cancel", func(c *gin.Context) {
				SetHeaderRole(c, entity.UserRoleScStaff)
				c.Request.Header.Set("X-User-ID", "invalid-uuid")
				SetContentTypeJSON(c)
				claimHandler.Cancel(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/cancel", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidUserID.ErrorCode)
		})
	})

	Describe("DoneReview", func() {
		Context("when authorized as EVM_STAFF", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims/:id/complete", entity.UserRoleEvmStaff, claimHandler.DoneReview)
			})

			It("should complete claim successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().Complete(mockTx, claimID, userID).Return(nil).Once()
				})

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/complete", w, nil)
				ExpectResponseNotNil(w, http.StatusNoContent)
			})

			It("should handle invalid UUID", func() {
				SendRequest(r, http.MethodPost, "/claims/invalid-uuid/complete", w, nil)
				ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidParams.ErrorCode)
			})

			It("should handle service errors", func() {
				invalidActionError := apperror.ErrInvalidClaimAction
				setupTxMockWithError(func() {
					mockService.EXPECT().Complete(mockTx, claimID, userID).
						Return(invalidActionError).Once()
				}, invalidActionError)

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/complete", w, nil)
				ExpectErrorCode(w, http.StatusConflict, apperror.ErrInvalidClaimAction.ErrorCode)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("POST", "/claims/:id/complete", entity.UserRoleScStaff, claimHandler.DoneReview)
			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/complete", w, nil)
			ExpectErrorCode(w, http.StatusForbidden, apperror.ErrUnauthorizedRole.ErrorCode)
		})

		It("should handle missing user ID header", func() {
			r.POST("/claims/:id/complete", func(c *gin.Context) {
				SetHeaderRole(c, entity.UserRoleEvmStaff)
				SetContentTypeJSON(c)
				claimHandler.DoneReview(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/complete", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrMissingUserID.ErrorCode)
		})

		It("should handle invalid user ID header", func() {
			r.POST("/claims/:id/complete", func(c *gin.Context) {
				SetHeaderRole(c, entity.UserRoleEvmStaff)
				c.Request.Header.Set("X-User-ID", "invalid-uuid")
				SetContentTypeJSON(c)
				claimHandler.DoneReview(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/complete", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidUserID.ErrorCode)
		})
	})

	Describe("History", func() {
		BeforeEach(func() {
			r.GET("/claims/:id/history", claimHandler.History)
		})

		It("should get claim history successfully", func() {
			sampleHistory := []*entity.ClaimHistory{
				entity.NewClaimHistory(claimID, entity.ClaimStatusSubmitted, uuid.New()),
			}
			mockService.EXPECT().GetHistory(mock.Anything, claimID).Return(sampleHistory, nil).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/history", nil)
			r.ServeHTTP(w, req)

			ExpectResponseNotNil(w, http.StatusOK)
		})

		It("should handle invalid UUID", func() {
			req, _ := http.NewRequest(http.MethodGet, "/claims/invalid-uuid/history", nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusBadRequest, apperror.ErrInvalidParams.ErrorCode)
		})

		It("should handle claim not found", func() {
			mockService.EXPECT().GetHistory(mock.Anything, claimID).
				Return(nil, apperror.ErrNotFoundError).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/history", nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusNotFound, apperror.ErrNotFoundError.ErrorCode)
		})

		It("should handle service errors", func() {
			mockService.EXPECT().GetHistory(mock.Anything, claimID).
				Return(nil, errors.New("database error")).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/history", nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusInternalServerError, apperror.ErrInternalServerError.ErrorCode)
		})
	})
})
