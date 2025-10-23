package handlers_test

import (
	"context"
	"errors"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/internal/interfaces/api/dtos"
	"ev-warranty-go/internal/interfaces/api/handlers"
	"ev-warranty-go/pkg/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"

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
		handler        handlers.ClaimHandler
		r              *gin.Engine
		w              *httptest.ResponseRecorder
		userID         uuid.UUID
		claimID        uuid.UUID
		sampleClaim    *entities.Claim
		validCreateReq dtos.CreateClaimRequest
		validUpdateReq dtos.UpdateClaimRequest
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
		handler = handlers.NewClaimHandler(mockLogger, mockTxManager, mockService)

		userID = uuid.New()
		claimID = uuid.New()
		vehicleID := uuid.New()
		customerID := uuid.New()

		validCreateReq = dtos.CreateClaimRequest{
			VehicleID:   vehicleID,
			CustomerID:  customerID,
			Description: "Test claim description for warranty issue",
		}

		validUpdateReq = dtos.UpdateClaimRequest{
			Description: "Updated claim description for warranty issue",
		}

		sampleClaim = entities.NewClaim(vehicleID, customerID, validCreateReq.Description, entities.ClaimStatusDraft, nil)
		sampleClaim.ID = claimID
	})

	Describe("GetByID", func() {
		BeforeEach(func() {
			r.GET("/claims/:id", handler.GetByID)
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

			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
		})

		It("should handle claim not found", func() {
			mockService.EXPECT().GetByID(mock.Anything, claimID).
				Return(nil, apperrors.NewClaimNotFound()).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String(), nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeClaimNotFound)
		})

		It("should handle service errors", func() {
			mockService.EXPECT().GetByID(mock.Anything, claimID).
				Return(nil, errors.New("database error")).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String(), nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeInternalServerError)
		})
	})

	Describe("GetAll", func() {
		BeforeEach(func() {
			r.GET("/claims", handler.GetAll)
		})

		It("should retrieve claims successfully with default pagination", func() {
			expectedResult := &services.ClaimListResult{
				Claims:     []*entities.Claim{sampleClaim},
				Total:      1,
				Page:       1,
				PageSize:   10,
				TotalPages: 1,
			}
			mockService.EXPECT().GetAll(mock.Anything, mock.AnythingOfType("services.ClaimFilters"), mock.AnythingOfType("services.Pagination")).
				Return(expectedResult, nil).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims", nil)
			r.ServeHTTP(w, req)

			ExpectResponseNotNil(w, http.StatusOK)
		})

		It("should handle query parameters correctly", func() {
			expectedResult := &services.ClaimListResult{
				Claims:     []*entities.Claim{},
				Total:      0,
				Page:       2,
				PageSize:   5,
				TotalPages: 1,
			}

			mockService.EXPECT().GetAll(mock.Anything, mock.MatchedBy(func(filters services.ClaimFilters) bool {
				return filters.Status != nil && *filters.Status == entities.ClaimStatusSubmitted &&
					filters.CustomerID != nil && *filters.CustomerID == sampleClaim.CustomerID
			}), mock.MatchedBy(func(pagination services.Pagination) bool {
				return pagination.Page == 2 && pagination.PageSize == 5
			})).Return(expectedResult, nil).Once()

			queryParams := url.Values{}
			queryParams.Add("status", entities.ClaimStatusSubmitted)
			queryParams.Add("customer_id", sampleClaim.CustomerID.String())
			queryParams.Add("page", "2")
			queryParams.Add("page_size", "5")

			req, _ := http.NewRequest(http.MethodGet, "/claims?"+queryParams.Encode(), nil)
			r.ServeHTTP(w, req)

			ExpectResponseNotNil(w, http.StatusOK)
		})

		DescribeTable("should handle invalid query parameters",
			func(queryParam, value, expectedError string) {
				queryParams := url.Values{}
				queryParams.Add(queryParam, value)

				req, _ := http.NewRequest(http.MethodGet, "/claims?"+queryParams.Encode(), nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusBadRequest, expectedError)
			},
			Entry("invalid customer_id UUID", "customer_id", "invalid-uuid", apperrors.ErrorCodeInvalidQueryParameter),
			Entry("invalid vehicle_id UUID", "vehicle_id", "invalid-uuid", apperrors.ErrorCodeInvalidQueryParameter),
			Entry("invalid status", "status", "INVALID_STATUS", apperrors.ErrorCodeInvalidQueryParameter),
			Entry("invalid from_date", "from_date", "invalid-date", apperrors.ErrorCodeInvalidQueryParameter),
			Entry("invalid to_date", "to_date", "invalid-date", apperrors.ErrorCodeInvalidQueryParameter),
		)

		It("should handle service errors", func() {
			mockService.EXPECT().GetAll(mock.Anything, mock.Anything, mock.Anything).
				Return(nil, errors.New("database error")).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims", nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeInternalServerError)
		})
	})

	Describe("Create", func() {
		Context("when authorized as SC_TECHNICIAN", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims", entities.UserRoleScTechnician, handler.Create)
			})

			It("should create claim successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().Create(mockTx, mock.MatchedBy(func(cmd *services.CreateClaimCommand) bool {
						return cmd.VehicleID == validCreateReq.VehicleID &&
							cmd.CustomerID == validCreateReq.CustomerID &&
							cmd.CreatorID == userID &&
							cmd.Description == validCreateReq.Description
					})).Return(sampleClaim, nil).Once()
				})

				SendRequest(r, http.MethodPost, "/claims", w, validCreateReq)
				ExpectResponseNotNil(w, http.StatusCreated)
			})

			It("should handle invalid JSON", func() {
				SendRequest(r, http.MethodPost, "/claims", w, "invalid json")
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidJsonRequest)
			})

			DescribeTable("should handle validation errors",
				func(modifyReq func(*dtos.CreateClaimRequest), expectedError string) {
					req := validCreateReq
					modifyReq(&req)

					SendRequest(r, http.MethodPost, "/claims", w, req)
					ExpectErrorCode(w, http.StatusBadRequest, expectedError)
				},
				Entry("empty description", func(req *dtos.CreateClaimRequest) {
					req.Description = ""
				}, apperrors.ErrorCodeInvalidJsonRequest),
				Entry("description too short", func(req *dtos.CreateClaimRequest) {
					req.Description = "short"
				}, apperrors.ErrorCodeInvalidJsonRequest),
			)

			It("should handle service errors", func() {
				dbError := errors.New("database error")
				setupTxMockWithError(func() {
					mockService.EXPECT().Create(mockTx, mock.Anything).
						Return(nil, dbError).Once()
				}, dbError)

				SendRequest(r, http.MethodPost, "/claims", w, validCreateReq)
				ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeInternalServerError)
			})

			It("should handle transaction errors", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Return(errors.New("transaction error")).Once()

				SendRequest(r, http.MethodPost, "/claims", w, validCreateReq)
				ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeInternalServerError)
			})
		})

		Context("when authorized as SC_STAFF", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims", entities.UserRoleScStaff, handler.Create)
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
			setupRoute("POST", "/claims", entities.UserRoleEvmStaff, handler.Create)
			SendRequest(r, http.MethodPost, "/claims", w, validCreateReq)
			ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
		})

		It("should handle missing user ID header", func() {
			r.POST("/claims", func(c *gin.Context) {
				SetHeaderRole(c, entities.UserRoleScTechnician)
				SetContentTypeJSON(c)
				handler.Create(c)
			})

			SendRequest(r, http.MethodPost, "/claims", w, validCreateReq)
			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeMissingUserID)
		})

		It("should handle invalid user ID header", func() {
			r.POST("/claims", func(c *gin.Context) {
				SetHeaderRole(c, entities.UserRoleScTechnician)
				c.Request.Header.Set("X-User-ID", "invalid-uuid")
				SetContentTypeJSON(c)
				handler.Create(c)
			})

			SendRequest(r, http.MethodPost, "/claims", w, validCreateReq)
			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUserID)
		})
	})

	Describe("Update", func() {
		Context("when authorized as SC_STAFF", func() {
			BeforeEach(func() {
				setupRoute("PUT", "/claims/:id", entities.UserRoleScStaff, handler.Update)
			})

			It("should update claim successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().Update(mockTx, claimID, mock.MatchedBy(func(cmd *services.UpdateClaimCommand) bool {
						return cmd.Description == validUpdateReq.Description
					})).Return(nil).Once()
				})

				SendRequest(r, http.MethodPut, "/claims/"+claimID.String(), w, validUpdateReq)
				ExpectResponseNotNil(w, http.StatusNoContent)
			})

			It("should handle invalid UUID", func() {
				SendRequest(r, http.MethodPut, "/claims/invalid-uuid", w, validUpdateReq)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})

			It("should handle invalid JSON", func() {
				SendRequest(r, http.MethodPut, "/claims/"+claimID.String(), w, "invalid json")
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidJsonRequest)
			})

			It("should handle service errors", func() {
				notFoundError := apperrors.NewClaimNotFound()
				setupTxMockWithError(func() {
					mockService.EXPECT().Update(mockTx, claimID, mock.Anything).
						Return(notFoundError).Once()
				}, notFoundError)

				SendRequest(r, http.MethodPut, "/claims/"+claimID.String(), w, validUpdateReq)
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeClaimNotFound)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("PUT", "/claims/:id", entities.UserRoleScTechnician, handler.Update)
			SendRequest(r, http.MethodPut, "/claims/"+claimID.String(), w, validUpdateReq)
			ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
		})
	})

	Describe("Delete", func() {
		Context("when authorized as SC_STAFF", func() {
			BeforeEach(func() {
				setupRoute("DELETE", "/claims/:id", entities.UserRoleScStaff, handler.Delete)
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
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})

			It("should handle service errors", func() {
				notFoundError := apperrors.NewClaimNotFound()
				setupTxMockWithError(func() {
					mockService.EXPECT().HardDelete(mockTx, claimID).
						Return(notFoundError).Once()
				}, notFoundError)

				SendRequest(r, http.MethodDelete, "/claims/"+claimID.String(), w, nil)
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeClaimNotFound)
			})
		})

		Context("when authorized as EVM_STAFF", func() {
			BeforeEach(func() {
				setupRoute("DELETE", "/claims/:id", entities.UserRoleEvmStaff, handler.Delete)
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
			setupRoute("DELETE", "/claims/:id", entities.UserRoleScTechnician, handler.Delete)
			SendRequest(r, http.MethodDelete, "/claims/"+claimID.String(), w, nil)
			ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
		})
	})

	Describe("Submit", func() {
		Context("when authorized as SC_STAFF", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims/:id/submit", entities.UserRoleScStaff, handler.Submit)
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
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})

			It("should handle service errors", func() {
				invalidActionError := apperrors.NewInvalidClaimAction()
				setupTxMockWithError(func() {
					mockService.EXPECT().Submit(mockTx, claimID, userID).
						Return(invalidActionError).Once()
				}, invalidActionError)

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/submit", w, nil)
				ExpectErrorCode(w, http.StatusConflict, apperrors.ErrorCodeInvalidClaimAction)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("POST", "/claims/:id/submit", entities.UserRoleScTechnician, handler.Submit)
			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/submit", w, nil)
			ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
		})

		It("should handle missing user ID header", func() {
			r.POST("/claims/:id/submit", func(c *gin.Context) {
				SetHeaderRole(c, entities.UserRoleScStaff)
				SetContentTypeJSON(c)
				handler.Submit(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/submit", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeMissingUserID)
		})

		It("should handle invalid user ID header", func() {
			r.POST("/claims/:id/submit", func(c *gin.Context) {
				SetHeaderRole(c, entities.UserRoleScStaff)
				c.Request.Header.Set("X-User-ID", "invalid-uuid")
				SetContentTypeJSON(c)
				handler.Submit(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/submit", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUserID)
		})
	})

	Describe("Review", func() {
		Context("when authorized as EVM_STAFF", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims/:id/review", entities.UserRoleEvmStaff, handler.Review)
			})

			It("should start review successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusReviewing, userID).Return(nil).Once()
				})

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/review", w, nil)
				ExpectResponseNotNil(w, http.StatusNoContent)
			})

			It("should handle invalid UUID", func() {
				SendRequest(r, http.MethodPost, "/claims/invalid-uuid/review", w, nil)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})

			It("should handle service errors", func() {
				invalidActionError := apperrors.NewInvalidClaimAction()
				setupTxMockWithError(func() {
					mockService.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusReviewing, userID).
						Return(invalidActionError).Once()
				}, invalidActionError)

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/review", w, nil)
				ExpectErrorCode(w, http.StatusConflict, apperrors.ErrorCodeInvalidClaimAction)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("POST", "/claims/:id/review", entities.UserRoleScStaff, handler.Review)
			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/review", w, nil)
			ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
		})

		It("should handle missing user ID header", func() {
			r.POST("/claims/:id/review", func(c *gin.Context) {
				SetHeaderRole(c, entities.UserRoleEvmStaff)
				SetContentTypeJSON(c)
				handler.Review(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/review", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeMissingUserID)
		})

		It("should handle invalid user ID header", func() {
			r.POST("/claims/:id/review", func(c *gin.Context) {
				SetHeaderRole(c, entities.UserRoleEvmStaff)
				c.Request.Header.Set("X-User-ID", "invalid-uuid")
				SetContentTypeJSON(c)
				handler.Review(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/review", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUserID)
		})
	})

	Describe("RequestInfo", func() {
		Context("when authorized as EVM_STAFF", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims/:id/request-info", entities.UserRoleEvmStaff, handler.RequestInfo)
			})

			It("should request info successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusRequestInfo, userID).Return(nil).Once()
				})

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/request-info", w, nil)
				ExpectResponseNotNil(w, http.StatusNoContent)
			})

			It("should handle invalid UUID", func() {
				SendRequest(r, http.MethodPost, "/claims/invalid-uuid/request-info", w, nil)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})

			It("should handle service errors", func() {
				invalidActionError := apperrors.NewInvalidClaimAction()
				setupTxMockWithError(func() {
					mockService.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusRequestInfo, userID).
						Return(invalidActionError).Once()
				}, invalidActionError)

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/request-info", w, nil)
				ExpectErrorCode(w, http.StatusConflict, apperrors.ErrorCodeInvalidClaimAction)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("POST", "/claims/:id/request-info", entities.UserRoleScStaff, handler.RequestInfo)
			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/request-info", w, nil)
			ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
		})

		It("should handle missing user ID header", func() {
			r.POST("/claims/:id/request-info", func(c *gin.Context) {
				SetHeaderRole(c, entities.UserRoleEvmStaff)
				SetContentTypeJSON(c)
				handler.RequestInfo(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/request-info", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeMissingUserID)
		})

		It("should handle invalid user ID header", func() {
			r.POST("/claims/:id/request-info", func(c *gin.Context) {
				SetHeaderRole(c, entities.UserRoleEvmStaff)
				c.Request.Header.Set("X-User-ID", "invalid-uuid")
				SetContentTypeJSON(c)
				handler.RequestInfo(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/request-info", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUserID)
		})
	})

	Describe("Cancel", func() {
		Context("when authorized as SC_STAFF", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims/:id/cancel", entities.UserRoleScStaff, handler.Cancel)
			})

			It("should cancel claim successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusCancelled, userID).Return(nil).Once()
				})

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/cancel", w, nil)
				ExpectResponseNotNil(w, http.StatusNoContent)
			})

			It("should handle invalid UUID", func() {
				SendRequest(r, http.MethodPost, "/claims/invalid-uuid/cancel", w, nil)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})

			It("should handle service errors", func() {
				invalidActionError := apperrors.NewInvalidClaimAction()
				setupTxMockWithError(func() {
					mockService.EXPECT().UpdateStatus(mockTx, claimID, entities.ClaimStatusCancelled, userID).
						Return(invalidActionError).Once()
				}, invalidActionError)

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/cancel", w, nil)
				ExpectErrorCode(w, http.StatusConflict, apperrors.ErrorCodeInvalidClaimAction)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("POST", "/claims/:id/cancel", entities.UserRoleEvmStaff, handler.Cancel)
			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/cancel", w, nil)
			ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
		})

		It("should handle missing user ID header", func() {
			r.POST("/claims/:id/cancel", func(c *gin.Context) {
				SetHeaderRole(c, entities.UserRoleScStaff)
				SetContentTypeJSON(c)
				handler.Cancel(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/cancel", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeMissingUserID)
		})

		It("should handle invalid user ID header", func() {
			r.POST("/claims/:id/cancel", func(c *gin.Context) {
				SetHeaderRole(c, entities.UserRoleScStaff)
				c.Request.Header.Set("X-User-ID", "invalid-uuid")
				SetContentTypeJSON(c)
				handler.Cancel(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/cancel", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUserID)
		})
	})

	Describe("Complete", func() {
		Context("when authorized as EVM_STAFF", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims/:id/complete", entities.UserRoleEvmStaff, handler.Complete)
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
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})

			It("should handle service errors", func() {
				invalidActionError := apperrors.NewInvalidClaimAction()
				setupTxMockWithError(func() {
					mockService.EXPECT().Complete(mockTx, claimID, userID).
						Return(invalidActionError).Once()
				}, invalidActionError)

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/complete", w, nil)
				ExpectErrorCode(w, http.StatusConflict, apperrors.ErrorCodeInvalidClaimAction)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("POST", "/claims/:id/complete", entities.UserRoleScStaff, handler.Complete)
			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/complete", w, nil)
			ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
		})

		It("should handle missing user ID header", func() {
			r.POST("/claims/:id/complete", func(c *gin.Context) {
				SetHeaderRole(c, entities.UserRoleEvmStaff)
				SetContentTypeJSON(c)
				handler.Complete(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/complete", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeMissingUserID)
		})

		It("should handle invalid user ID header", func() {
			r.POST("/claims/:id/complete", func(c *gin.Context) {
				SetHeaderRole(c, entities.UserRoleEvmStaff)
				c.Request.Header.Set("X-User-ID", "invalid-uuid")
				SetContentTypeJSON(c)
				handler.Complete(c)
			})

			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/complete", w, nil)
			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUserID)
		})
	})

	Describe("History", func() {
		BeforeEach(func() {
			r.GET("/claims/:id/history", handler.History)
		})

		It("should get claim history successfully", func() {
			sampleHistory := []*entities.ClaimHistory{
				entities.NewClaimHistory(claimID, entities.ClaimStatusSubmitted, uuid.New()),
			}
			mockService.EXPECT().GetHistory(mock.Anything, claimID).Return(sampleHistory, nil).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/history", nil)
			r.ServeHTTP(w, req)

			ExpectResponseNotNil(w, http.StatusOK)
		})

		It("should handle invalid UUID", func() {
			req, _ := http.NewRequest(http.MethodGet, "/claims/invalid-uuid/history", nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
		})

		It("should handle claim not found", func() {
			mockService.EXPECT().GetHistory(mock.Anything, claimID).
				Return(nil, apperrors.NewClaimNotFound()).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/history", nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeClaimNotFound)
		})

		It("should handle service errors", func() {
			mockService.EXPECT().GetHistory(mock.Anything, claimID).
				Return(nil, errors.New("database error")).Once()

			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/history", nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeInternalServerError)
		})
	})
})
