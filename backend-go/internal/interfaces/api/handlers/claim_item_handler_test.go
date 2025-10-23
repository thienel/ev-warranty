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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("ClaimItemHandler", func() {
	var (
		mockLogger      *mocks.Logger
		mockTxManager   *mocks.TxManager
		mockService     *mocks.ClaimItemService
		mockTx          *mocks.Tx
		handler         handlers.ClaimItemHandler
		r               *gin.Engine
		w               *httptest.ResponseRecorder
		claimID         uuid.UUID
		itemID          uuid.UUID
		sampleClaimItem *entities.ClaimItem
		validReq        dtos.CreateClaimItemRequest
	)

	setupTxMock := func(serviceMockFn func()) {
		mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
			Run(func(ctx context.Context, fn func(application.Tx) error) {
				serviceMockFn()
				_ = fn(mockTx)
			}).Return(nil).Once()
	}

	setupRoute := func(method, path string, role string, handlerFunc gin.HandlerFunc) {
		r.Handle(method, path, func(c *gin.Context) {
			if role != "" {
				SetHeaderRole(c, role)
			}
			SetContentTypeJSON(c)
			handlerFunc(c)
		})
	}

	BeforeEach(func() {
		mockLogger, r, w = SetupMock(GinkgoT())
		mockTxManager = mocks.NewTxManager(GinkgoT())
		mockService = mocks.NewClaimItemService(GinkgoT())
		mockTx = mocks.NewTx(GinkgoT())
		handler = handlers.NewClaimItemHandler(mockLogger, mockTxManager, mockService)

		claimID = uuid.New()
		itemID = uuid.New()
		validReq = dtos.CreateClaimItemRequest{
			PartCategoryID:    1,
			FaultyPartID:      uuid.New(),
			ReplacementPartID: func() *uuid.UUID { id := uuid.New(); return &id }(),
			IssueDescription:  "Test issue description for replacement",
			Type:              entities.ClaimItemTypeReplacement,
			Cost:              100.50,
		}
		sampleClaimItem = &entities.ClaimItem{
			ID:                itemID,
			ClaimID:           claimID,
			PartCategoryID:    validReq.PartCategoryID,
			FaultyPartID:      validReq.FaultyPartID,
			ReplacementPartID: validReq.ReplacementPartID,
			IssueDescription:  validReq.IssueDescription,
			Status:            entities.ClaimItemStatusPending,
			Type:              validReq.Type,
			Cost:              validReq.Cost,
		}
	})

	Describe("GetByID", func() {
		BeforeEach(func() {
			r.GET("/claims/:id/items/:itemID", handler.GetByID)
		})

		It("should handle successful retrieval", func() {
			mockService.EXPECT().GetByID(mock.Anything, itemID).Return(sampleClaimItem, nil).Once()
			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/items/"+itemID.String(), nil)
			r.ServeHTTP(w, req)
			ExpectResponseNotNil(w, http.StatusOK)
		})

		It("should handle invalid item UUID", func() {
			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/items/invalid-uuid", nil)
			r.ServeHTTP(w, req)
			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
		})

		It("should handle item not found", func() {
			mockService.EXPECT().GetByID(mock.Anything, itemID).
				Return(nil, apperrors.NewClaimItemNotFound()).Once()
			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/items/"+itemID.String(), nil)
			r.ServeHTTP(w, req)
			ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeClaimItemNotFound)
		})

		It("should handle service error", func() {
			mockService.EXPECT().GetByID(mock.Anything, itemID).
				Return(nil, errors.New("database error")).Once()
			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/items/"+itemID.String(), nil)
			r.ServeHTTP(w, req)
			ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeInternalServerError)
		})
	})

	Describe("GetByClaimID", func() {
		BeforeEach(func() {
			r.GET("/claims/:id/items", handler.GetByClaimID)
		})

		It("should handle successful retrieval with items", func() {
			items := []*entities.ClaimItem{sampleClaimItem}
			mockService.EXPECT().GetByClaimID(mock.Anything, claimID).Return(items, nil).Once()
			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/items", nil)
			r.ServeHTTP(w, req)
			ExpectResponseNotNil(w, http.StatusOK)
		})

		It("should handle invalid claim UUID", func() {
			req, _ := http.NewRequest(http.MethodGet, "/claims/invalid-uuid/items", nil)
			r.ServeHTTP(w, req)
			ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
		})

		It("should handle service error", func() {
			dbError := apperrors.NewDBOperationError(errors.New("database error"))
			mockService.EXPECT().GetByClaimID(mock.Anything, claimID).Return(nil, dbError).Once()
			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/items", nil)
			r.ServeHTTP(w, req)
			ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeDBOperation)
		})

		It("should handle empty results", func() {
			mockService.EXPECT().GetByClaimID(mock.Anything, claimID).Return([]*entities.ClaimItem{}, nil).Once()
			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/items", nil)
			r.ServeHTTP(w, req)
			ExpectResponseNotNil(w, http.StatusOK)
		})
	})

	Describe("Create", func() {
		Context("when authorized as SC_STAFF", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims/:id/items", entities.UserRoleScStaff, handler.Create)
			})

			It("should create claim item successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().Create(mockTx, claimID, mock.MatchedBy(func(cmd *services.CreateClaimItemCommand) bool {
						return cmd.PartCategoryID == validReq.PartCategoryID && cmd.Type == validReq.Type
					})).Return(sampleClaimItem, nil).Once()
				})

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/items", w, validReq)
				ExpectResponseNotNil(w, http.StatusCreated)
			})

			DescribeTable("should handle error scenarios",
				func(setupMock func(), url string, reqBody interface{}, expectedStatus int, expectedError string) {
					if setupMock != nil {
						setupMock()
					}
					SendRequest(r, http.MethodPost, url, w, reqBody)
					ExpectErrorCode(w, expectedStatus, expectedError)
				},
				Entry("invalid claim UUID",
					nil,
					"/claims/invalid-uuid/items", validReq, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID),
				Entry("invalid JSON",
					nil,
					"/claims/"+claimID.String()+"/items", "invalid json", http.StatusBadRequest, apperrors.ErrorCodeInvalidJsonRequest),
			)

			It("should handle invalid claim item type", func() {
				invalidReq := validReq
				invalidReq.Type = "INVALID_TYPE"
				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/items", w, invalidReq)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidClaimItemType)
			})

			It("should handle service error during creation", func() {
				setupTxMock(func() {
					mockService.EXPECT().Create(mockTx, claimID, mock.Anything).
						Return(nil, apperrors.NewClaimNotFound()).Once()
				})
				mockTxManager.ExpectedCalls[0].ReturnArguments = []interface{}{apperrors.NewClaimNotFound()}

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/items", w, validReq)
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeClaimNotFound)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("POST", "/claims/:id/items", entities.UserRoleEvmStaff, handler.Create)
			SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/items", w, validReq)
			ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
		})
	})

	Describe("Delete", func() {
		Context("when authorized as SC_STAFF", func() {
			BeforeEach(func() {
				setupRoute("DELETE", "/claims/:id/items/:itemID", entities.UserRoleScStaff, handler.Delete)
			})

			It("should delete claim item successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().HardDelete(mockTx, claimID, itemID).Return(nil).Once()
				})

				req, _ := http.NewRequest(http.MethodDelete, "/claims/"+claimID.String()+"/items/"+itemID.String(), nil)
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})

			DescribeTable("should handle error scenarios",
				func(setupMock func(), url string, expectedStatus int, expectedError string) {
					if setupMock != nil {
						setupMock()
					}
					req, _ := http.NewRequest(http.MethodDelete, url, nil)
					r.ServeHTTP(w, req)
					ExpectErrorCode(w, expectedStatus, expectedError)
				},
				Entry("invalid claim UUID",
					nil,
					"/claims/invalid-uuid/items/"+itemID.String(), http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID),
				Entry("invalid item UUID",
					nil,
					"/claims/"+claimID.String()+"/items/invalid-uuid", http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID),
			)

			It("should handle service error during deletion", func() {
				setupTxMock(func() {
					mockService.EXPECT().HardDelete(mockTx, claimID, itemID).
						Return(apperrors.NewClaimItemNotFound()).Once()
				})
				mockTxManager.ExpectedCalls[0].ReturnArguments = []interface{}{apperrors.NewClaimItemNotFound()}

				req, _ := http.NewRequest(http.MethodDelete, "/claims/"+claimID.String()+"/items/"+itemID.String(), nil)
				r.ServeHTTP(w, req)
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeClaimItemNotFound)
			})

			It("should handle transaction manager error", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Return(apperrors.NewDBOperationError(errors.New("transaction failed"))).Once()

				req, _ := http.NewRequest(http.MethodDelete, "/claims/"+claimID.String()+"/items/"+itemID.String(), nil)
				r.ServeHTTP(w, req)
				ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeDBOperation)
			})

		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("DELETE", "/claims/:id/items/:itemID", entities.UserRoleScTechnician, handler.Delete)
			req, _ := http.NewRequest(http.MethodDelete, "/claims/"+claimID.String()+"/items/"+itemID.String(), nil)
			r.ServeHTTP(w, req)
			ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
		})
	})

	Describe("Approve", func() {
		Context("when authorized as EVM_STAFF", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims/:id/items/:itemID/approve", entities.UserRoleEvmStaff, handler.Approve)
			})

			It("should approve claim item successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().Approve(mockTx, claimID, itemID).Return(nil).Once()
				})

				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/"+itemID.String()+"/approve", nil)
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})

			DescribeTable("should handle error scenarios",
				func(setupMock func(), url string, expectedStatus int, expectedError string) {
					if setupMock != nil {
						setupMock()
					}
					req, _ := http.NewRequest(http.MethodPost, url, nil)
					r.ServeHTTP(w, req)
					ExpectErrorCode(w, expectedStatus, expectedError)
				},
				Entry("invalid claim UUID",
					nil,
					"/claims/invalid-uuid/items/"+itemID.String()+"/approve", http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID),
				Entry("invalid item UUID",
					nil,
					"/claims/"+claimID.String()+"/items/invalid-uuid/approve", http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID),
			)

			It("should handle service error during approval", func() {
				setupTxMock(func() {
					mockService.EXPECT().Approve(mockTx, claimID, itemID).
						Return(apperrors.NewClaimItemNotFound()).Once()
				})
				mockTxManager.ExpectedCalls[0].ReturnArguments = []interface{}{apperrors.NewClaimItemNotFound()}

				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/"+itemID.String()+"/approve", nil)
				r.ServeHTTP(w, req)
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeClaimItemNotFound)
			})

			It("should handle transaction manager error", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Return(apperrors.NewDBOperationError(errors.New("transaction failed"))).Once()

				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/"+itemID.String()+"/approve", nil)
				r.ServeHTTP(w, req)
				ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeDBOperation)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("POST", "/claims/:id/items/:itemID/approve", entities.UserRoleScStaff, handler.Approve)
			req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/"+itemID.String()+"/approve", nil)
			r.ServeHTTP(w, req)
			ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
		})
	})

	Describe("Reject", func() {
		Context("when authorized as EVM_STAFF", func() {
			BeforeEach(func() {
				setupRoute("POST", "/claims/:id/items/:itemID/reject", entities.UserRoleEvmStaff, handler.Reject)
			})

			It("should reject claim item successfully", func() {
				setupTxMock(func() {
					mockService.EXPECT().Reject(mockTx, claimID, itemID).Return(nil).Once()
				})

				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/"+itemID.String()+"/reject", nil)
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})

			DescribeTable("should handle error scenarios",
				func(setupMock func(), url string, expectedStatus int, expectedError string) {
					if setupMock != nil {
						setupMock()
					}
					req, _ := http.NewRequest(http.MethodPost, url, nil)
					r.ServeHTTP(w, req)
					ExpectErrorCode(w, expectedStatus, expectedError)
				},
				Entry("invalid claim UUID",
					nil,
					"/claims/invalid-uuid/items/"+itemID.String()+"/reject", http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID),
				Entry("invalid item UUID",
					nil,
					"/claims/"+claimID.String()+"/items/invalid-uuid/reject", http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID),
			)

			It("should handle service error during rejection", func() {
				setupTxMock(func() {
					mockService.EXPECT().Reject(mockTx, claimID, itemID).
						Return(apperrors.NewClaimItemNotFound()).Once()
				})
				mockTxManager.ExpectedCalls[0].ReturnArguments = []interface{}{apperrors.NewClaimItemNotFound()}

				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/"+itemID.String()+"/reject", nil)
				r.ServeHTTP(w, req)
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeClaimItemNotFound)
			})

			It("should handle transaction manager error", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Return(apperrors.NewDBOperationError(errors.New("transaction failed"))).Once()

				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/"+itemID.String()+"/reject", nil)
				r.ServeHTTP(w, req)
				ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeDBOperation)
			})
		})

		It("should deny access for unauthorized roles", func() {
			setupRoute("POST", "/claims/:id/items/:itemID/reject", entities.UserRoleScStaff, handler.Reject)
			req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/"+itemID.String()+"/reject", nil)
			r.ServeHTTP(w, req)
			ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
		})
	})
})
