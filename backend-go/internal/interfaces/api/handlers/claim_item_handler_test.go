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
		mockLogger        *mocks.Logger
		mockTxManager     *mocks.TxManager
		mockService       *mocks.ClaimItemService
		mockTx            *mocks.Tx
		handler           handlers.ClaimItemHandler
		r                 *gin.Engine
		w                 *httptest.ResponseRecorder
		claimID           uuid.UUID
		itemID            uuid.UUID
		replacementPartID uuid.UUID
		faultyPartID      uuid.UUID
	)

	BeforeEach(func() {
		mockLogger, r, w = SetupMock(GinkgoT())
		mockTxManager = mocks.NewTxManager(GinkgoT())
		mockService = mocks.NewClaimItemService(GinkgoT())
		mockTx = mocks.NewTx(GinkgoT())
		handler = handlers.NewClaimItemHandler(mockLogger, mockTxManager, mockService)

		claimID = uuid.New()
		itemID = uuid.New()
		replacementPartID = uuid.New()
		faultyPartID = uuid.New()
	})

	Describe("GetByID", func() {
		Context("when claim item is found successfully", func() {
			It("should return the claim item", func() {
				claimItem := &entities.ClaimItem{
					ID:                itemID,
					ClaimID:           claimID,
					PartCategoryID:    1,
					FaultyPartID:      faultyPartID,
					ReplacementPartID: &replacementPartID,
					IssueDescription:  "Test issue description",
					Status:            entities.ClaimItemStatusPending,
					Type:              entities.ClaimItemTypeReplacement,
					Cost:              100.50,
				}

				mockService.EXPECT().GetByID(mock.Anything, itemID).Return(claimItem, nil).Once()

				r.GET("/claims/:id/items/:itemID", handler.GetByID)
				req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/items/"+itemID.String(), nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				ExpectResponseNotNil(w, http.StatusOK)
			})
		})

		Context("when item ID is invalid UUID", func() {
			It("should return invalid UUID error", func() {
				r.GET("/claims/:id/items/:itemID", handler.GetByID)
				req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/items/invalid-uuid", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when claim item is not found", func() {
			It("should return not found error", func() {
				mockService.EXPECT().GetByID(mock.Anything, itemID).
					Return(nil, apperrors.NewClaimItemNotFound()).Once()

				r.GET("/claims/:id/items/:itemID", handler.GetByID)
				req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/items/"+itemID.String(), nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeClaimItemNotFound)
			})
		})

		Context("when service returns error", func() {
			It("should return internal server error", func() {
				mockService.EXPECT().GetByID(mock.Anything, itemID).
					Return(nil, errors.New("database error")).Once()

				r.GET("/claims/:id/items/:itemID", handler.GetByID)
				req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/items/"+itemID.String(), nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeInternalServerError)
			})
		})
	})

	Describe("GetByClaimID", func() {
		Context("when claim items are found successfully", func() {
			It("should return the claim items", func() {
				claimItems := []*entities.ClaimItem{
					{
						ID:                itemID,
						ClaimID:           claimID,
						PartCategoryID:    1,
						FaultyPartID:      faultyPartID,
						ReplacementPartID: &replacementPartID,
						IssueDescription:  "Test issue description",
						Status:            entities.ClaimItemStatusPending,
						Type:              entities.ClaimItemTypeReplacement,
						Cost:              100.50,
					},
				}

				mockService.EXPECT().GetByClaimID(mock.Anything, claimID).Return(claimItems, nil).Once()

				r.GET("/claims/:id/items", handler.GetByClaimID)
				req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/items", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				ExpectResponseNotNil(w, http.StatusOK)
			})
		})

		Context("when claim ID is invalid UUID", func() {
			It("should return invalid UUID error", func() {
				r.GET("/claims/:id/items", handler.GetByClaimID)
				req, _ := http.NewRequest(http.MethodGet, "/claims/invalid-uuid/items", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when no claim items are found", func() {
			It("should return empty array", func() {
				mockService.EXPECT().GetByClaimID(mock.Anything, claimID).Return([]*entities.ClaimItem{}, nil).Once()

				r.GET("/claims/:id/items", handler.GetByClaimID)
				req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/items", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				ExpectResponseNotNil(w, http.StatusOK)
			})
		})

		Context("when service returns error", func() {
			It("should return internal server error", func() {
				mockService.EXPECT().GetByClaimID(mock.Anything, claimID).
					Return(nil, errors.New("database error")).Once()

				r.GET("/claims/:id/items", handler.GetByClaimID)
				req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/items", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeInternalServerError)
			})
		})
	})

	Describe("Create", func() {
		var validRequest dtos.CreateClaimItemRequest

		BeforeEach(func() {
			validRequest = dtos.CreateClaimItemRequest{
				PartCategoryID:    1,
				FaultyPartID:      faultyPartID,
				ReplacementPartID: &replacementPartID,
				IssueDescription:  "Test issue description for replacement",
				Type:              entities.ClaimItemTypeReplacement,
				Cost:              100.50,
			}
		})

		Context("when creation is successful", func() {
			It("should create and return the claim item", func() {
				createdItem := &entities.ClaimItem{
					ID:                itemID,
					ClaimID:           claimID,
					PartCategoryID:    validRequest.PartCategoryID,
					FaultyPartID:      validRequest.FaultyPartID,
					ReplacementPartID: validRequest.ReplacementPartID,
					IssueDescription:  validRequest.IssueDescription,
					Status:            entities.ClaimItemStatusPending,
					Type:              validRequest.Type,
					Cost:              validRequest.Cost,
				}

				mockService.EXPECT().Create(mockTx, claimID, mock.MatchedBy(func(cmd *services.CreateClaimItemCommand) bool {
					return cmd.PartCategoryID == validRequest.PartCategoryID &&
						cmd.FaultyPartID == validRequest.FaultyPartID &&
						uuidPtrEqual(cmd.ReplacementPartID, validRequest.ReplacementPartID) &&
						cmd.IssueDescription == validRequest.IssueDescription &&
						cmd.Status == entities.ClaimItemStatusPending &&
						cmd.Type == validRequest.Type &&
						cmd.Cost == validRequest.Cost
				})).Return(createdItem, nil).Once()

				mockTxManager.EXPECT().Do(mock.Anything, mock.Anything).
					Run(func(ctx context.Context, fn func(application.Tx) error) {
						_ = fn(mockTx)
					}).Return(nil).Once()

				r.POST("/claims/:id/items", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScStaff)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/items", w, validRequest)
				ExpectResponseNotNil(w, http.StatusCreated)
			})
		})

		Context("when user is not authorized", func() {
			It("should return forbidden error", func() {
				r.POST("/claims/:id/items", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleEvmStaff)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/items", w, validRequest)
				ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
			})
		})

		Context("when claim ID is invalid UUID", func() {
			It("should return invalid UUID error", func() {
				r.POST("/claims/:id/items", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScStaff)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/claims/invalid-uuid/items", w, validRequest)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when request body is invalid JSON", func() {
			It("should return invalid JSON error", func() {
				r.POST("/claims/:id/items", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScStaff)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/items", w, "invalid json")
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidJsonRequest)
			})
		})

		Context("when claim item type is invalid", func() {
			It("should return invalid claim item type error", func() {
				invalidRequest := validRequest
				invalidRequest.Type = "INVALID_TYPE"

				r.POST("/claims/:id/items", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScStaff)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/items", w, invalidRequest)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidClaimItemType)
			})
		})

		Context("when service returns error", func() {
			It("should return the service error", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Run(func(ctx context.Context, fn func(application.Tx) error) {
						mockService.EXPECT().Create(mockTx, claimID, mock.Anything).
							Return(nil, apperrors.NewClaimNotFound()).Once()
						_ = fn(mockTx)
					}).Return(apperrors.NewClaimNotFound()).Once()

				r.POST("/claims/:id/items", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScStaff)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/claims/"+claimID.String()+"/items", w, validRequest)
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeClaimNotFound)
			})
		})
	})

	Describe("Delete", func() {
		Context("when deletion is successful", func() {
			It("should delete the claim item and return no content", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Run(func(ctx context.Context, fn func(application.Tx) error) {
						mockService.EXPECT().HardDelete(mockTx, claimID, itemID).Return(nil).Once()
						_ = fn(mockTx)
					}).Return(nil).Once()

				r.DELETE("/claims/:id/items/:itemID", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScStaff)
					handler.Delete(c)
				})

				req, _ := http.NewRequest(http.MethodDelete, "/claims/"+claimID.String()+"/items/"+itemID.String(), nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when user is not authorized", func() {
			It("should return forbidden error", func() {
				r.DELETE("/claims/:id/items/:itemID", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleEvmStaff) // Wrong role
					handler.Delete(c)
				})

				req, _ := http.NewRequest(http.MethodDelete, "/claims/"+claimID.String()+"/items/"+itemID.String(), nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
			})
		})

		Context("when claim ID is invalid UUID", func() {
			It("should return invalid UUID error", func() {
				r.DELETE("/claims/:id/items/:itemID", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScStaff)
					handler.Delete(c)
				})

				req, _ := http.NewRequest(http.MethodDelete, "/claims/invalid-uuid/items/"+itemID.String(), nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when item ID is invalid UUID", func() {
			It("should return invalid UUID error", func() {
				r.DELETE("/claims/:id/items/:itemID", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScStaff)
					handler.Delete(c)
				})

				req, _ := http.NewRequest(http.MethodDelete, "/claims/"+claimID.String()+"/items/invalid-uuid", nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when claim item is not found", func() {
			It("should return not found error", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Run(func(ctx context.Context, fn func(application.Tx) error) {
						mockService.EXPECT().HardDelete(mockTx, claimID, itemID).
							Return(apperrors.NewClaimItemNotFound()).Once()
						_ = fn(mockTx)
					}).Return(apperrors.NewClaimItemNotFound()).Once()

				r.DELETE("/claims/:id/items/:itemID", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScStaff)
					handler.Delete(c)
				})

				req, _ := http.NewRequest(http.MethodDelete, "/claims/"+claimID.String()+"/items/"+itemID.String(), nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeClaimItemNotFound)
			})
		})
	})

	Describe("Approve", func() {
		Context("when approval is successful", func() {
			It("should approve the claim item and return no content", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Run(func(ctx context.Context, fn func(application.Tx) error) {
						mockService.EXPECT().Approve(mockTx, claimID, itemID).Return(nil).Once()
						_ = fn(mockTx)
					}).Return(nil).Once()

				r.POST("/claims/:id/items/:itemID/approve", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleEvmStaff)
					handler.Approve(c)
				})

				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/"+itemID.String()+"/approve", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when user is not authorized", func() {
			It("should return forbidden error", func() {
				r.POST("/claims/:id/items/:itemID/approve", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScStaff) // Wrong role
					handler.Approve(c)
				})

				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/"+itemID.String()+"/approve", nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
			})
		})

		Context("when claim ID is invalid UUID", func() {
			It("should return invalid UUID error", func() {
				r.POST("/claims/:id/items/:itemID/approve", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleEvmStaff)
					handler.Approve(c)
				})

				req, _ := http.NewRequest(http.MethodPost, "/claims/invalid-uuid/items/"+itemID.String()+"/approve", nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when item ID is invalid UUID", func() {
			It("should return invalid UUID error", func() {
				r.POST("/claims/:id/items/:itemID/approve", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleEvmStaff)
					handler.Approve(c)
				})

				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/invalid-uuid/approve", nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when claim item is not found", func() {
			It("should return not found error", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Run(func(ctx context.Context, fn func(application.Tx) error) {
						mockService.EXPECT().Approve(mockTx, claimID, itemID).
							Return(apperrors.NewClaimItemNotFound()).Once()
						_ = fn(mockTx)
					}).Return(apperrors.NewClaimItemNotFound()).Once()

				r.POST("/claims/:id/items/:itemID/approve", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleEvmStaff)
					handler.Approve(c)
				})

				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/"+itemID.String()+"/approve", nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeClaimItemNotFound)
			})
		})
	})

	Describe("Reject", func() {
		Context("when rejection is successful", func() {
			It("should reject the claim item and return no content", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Run(func(ctx context.Context, fn func(application.Tx) error) {
						mockService.EXPECT().Reject(mockTx, claimID, itemID).Return(nil).Once()
						_ = fn(mockTx)
					}).Return(nil).Once()

				r.POST("/claims/:id/items/:itemID/reject", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleEvmStaff)
					handler.Reject(c)
				})

				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/"+itemID.String()+"/reject", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when user is not authorized", func() {
			It("should return forbidden error", func() {
				r.POST("/claims/:id/items/:itemID/reject", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScStaff) // Wrong role
					handler.Reject(c)
				})

				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/"+itemID.String()+"/reject", nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
			})
		})

		Context("when claim ID is invalid UUID", func() {
			It("should return invalid UUID error", func() {
				r.POST("/claims/:id/items/:itemID/reject", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleEvmStaff)
					handler.Reject(c)
				})

				req, _ := http.NewRequest(http.MethodPost, "/claims/invalid-uuid/items/"+itemID.String()+"/reject", nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when item ID is invalid UUID", func() {
			It("should return invalid UUID error", func() {
				r.POST("/claims/:id/items/:itemID/reject", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleEvmStaff)
					handler.Reject(c)
				})

				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/invalid-uuid/reject", nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when claim item is not found", func() {
			It("should return not found error", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Run(func(ctx context.Context, fn func(application.Tx) error) {
						mockService.EXPECT().Reject(mockTx, claimID, itemID).
							Return(apperrors.NewClaimItemNotFound()).Once()
						_ = fn(mockTx)
					}).Return(apperrors.NewClaimItemNotFound()).Once()

				r.POST("/claims/:id/items/:itemID/reject", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleEvmStaff)
					handler.Reject(c)
				})

				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/items/"+itemID.String()+"/reject", nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeClaimItemNotFound)
			})
		})
	})
})
