package handlers_test

import (
	"errors"
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

var _ = Describe("OfficeHandler", func() {
	var (
		mockLogger  *mocks.Logger
		mockService *mocks.OfficeService
		handler     handlers.OfficeHandler
		r           *gin.Engine
		w           *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		mockLogger, r, w = SetupMock(GinkgoT())
		mockService = mocks.NewOfficeService(GinkgoT())
		handler = handlers.NewOfficeHandler(mockLogger, mockService)
	})

	Describe("Create", func() {
		var (
			validReq dtos.CreateOfficeRequest
		)

		BeforeEach(func() {
			validReq = dtos.CreateOfficeRequest{
				OfficeName: "Test Office",
				OfficeType: entities.OfficeTypeEVM,
				Address:    "123 Test Street",
				IsActive:   true,
			}
		})

		Context("when office is created successfully", func() {
			It("should return 201 with created office", func() {
				office := CreateOfficeFromRequest(validReq)

				mockService.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cmd *services.CreateOfficeCommand) bool {
					return cmd.OfficeName == validReq.OfficeName &&
						cmd.OfficeType == validReq.OfficeType &&
						cmd.Address == validReq.Address &&
						cmd.IsActive == validReq.IsActive
				})).Return(office, nil).Once()

				r.POST("/offices", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/offices", w, validReq)
				ExpectResponseNotNil(w, http.StatusCreated)
			})
		})

		Context("when user is not admin", func() {
			It("should return 403 forbidden role error", func() {
				r.POST("/offices", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScStaff)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/offices", w, validReq)
				ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
			})
		})

		Context("when request body is invalid JSON", func() {
			It("should return 400 bad request", func() {
				r.POST("/offices", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/offices", w, "invalid json")
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidJsonRequest)
			})
		})

		Context("when office type is invalid", func() {
			It("should return 400 bad request", func() {
				invalidReq := validReq
				invalidReq.OfficeType = "INVALID_TYPE"

				r.POST("/offices", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/offices", w, invalidReq)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidOfficeType)
			})
		})

		Context("when service returns error", func() {
			It("should return error from service", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, dbErr).Once()

				r.POST("/offices", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/offices", w, validReq)
				ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeDBOperation)
			})
		})

		Context("when required fields are missing", func() {
			It("should return 400 bad request for missing office_name", func() {
				invalidReq := validReq
				invalidReq.OfficeName = ""

				r.POST("/offices", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/offices", w, invalidReq)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when office type is SC", func() {
			It("should create office with SC type successfully", func() {
				scReq := validReq
				scReq.OfficeType = entities.OfficeTypeSC
				office := CreateOfficeFromRequest(scReq)

				mockService.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cmd *services.CreateOfficeCommand) bool {
					return cmd.OfficeType == entities.OfficeTypeSC
				})).Return(office, nil).Once()

				r.POST("/offices", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/offices", w, scReq)
				Expect(w.Code).To(Equal(http.StatusCreated))
			})
		})

		Context("when office is created as inactive", func() {
			It("should create inactive office successfully", func() {
				inactiveReq := validReq
				inactiveReq.IsActive = false
				office := CreateOfficeFromRequest(inactiveReq)

				mockService.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cmd *services.CreateOfficeCommand) bool {
					return cmd.IsActive == false
				})).Return(office, nil).Once()

				r.POST("/offices", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/offices", w, inactiveReq)
				Expect(w.Code).To(Equal(http.StatusCreated))
			})
		})
	})

	Describe("GetById", func() {
		var officeID uuid.UUID

		BeforeEach(func() {
			officeID = uuid.New()
		})

		Context("when office is found", func() {
			It("should return 200 with office data", func() {
				office := &entities.Office{
					ID:         officeID,
					OfficeName: "Test Office",
					OfficeType: entities.OfficeTypeEVM,
					Address:    "123 Test Street",
					IsActive:   true,
				}
				mockService.EXPECT().GetByID(mock.Anything, officeID).Return(office, nil).Once()

				r.GET("/offices/:id", handler.GetById)

				SendRequest(r, http.MethodGet, "/offices/"+officeID.String(), w, nil)
				ExpectResponseNotNil(w, http.StatusOK)
			})
		})

		Context("when office ID is invalid UUID", func() {
			It("should return 400 bad request", func() {
				r.GET("/offices/:id", handler.GetById)

				SendRequest(r, http.MethodGet, "/offices/invalid-uuid", w, nil)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when office is not found", func() {
			It("should return 404 not found", func() {
				notFoundErr := apperrors.NewOfficeNotFound()
				mockService.EXPECT().GetByID(mock.Anything, officeID).Return(nil, notFoundErr).Once()

				r.GET("/offices/:id", handler.GetById)

				SendRequest(r, http.MethodGet, "/offices/"+officeID.String(), w, nil)
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeOfficeNotFound)
			})
		})

		Context("when service returns database error", func() {
			It("should return 500 internal server error", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().GetByID(mock.Anything, officeID).Return(nil, dbErr).Once()

				r.GET("/offices/:id", handler.GetById)

				SendRequest(r, http.MethodGet, "/offices/"+officeID.String(), w, nil)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("GetAll", func() {
		Context("when offices are found", func() {
			It("should return 200 with all offices", func() {
				offices := []*entities.Office{
					{
						ID:         uuid.New(),
						OfficeName: "Office 1",
						OfficeType: entities.OfficeTypeEVM,
						Address:    "Address 1",
						IsActive:   true,
					},
					{
						ID:         uuid.New(),
						OfficeName: "Office 2",
						OfficeType: entities.OfficeTypeSC,
						Address:    "Address 2",
						IsActive:   false,
					},
				}

				mockService.EXPECT().GetAll(mock.Anything).Return(offices, nil).Once()

				r.GET("/offices", handler.GetAll)

				SendRequest(r, http.MethodGet, "/offices", w, nil)
				ExpectResponseNotNil(w, http.StatusOK)
			})
		})

		Context("when no offices exist", func() {
			It("should return 200 with empty array", func() {
				mockService.EXPECT().GetAll(mock.Anything).Return([]*entities.Office{}, nil).Once()

				r.GET("/offices", handler.GetAll)

				SendRequest(r, http.MethodGet, "/offices", w, nil)
				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})

		Context("when service returns error", func() {
			It("should return 500 internal server error", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().GetAll(mock.Anything).Return(nil, dbErr).Once()

				r.GET("/offices", handler.GetAll)

				SendRequest(r, http.MethodGet, "/offices", w, nil)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("Update", func() {
		var (
			officeID uuid.UUID
			validReq dtos.UpdateOfficeRequest
		)

		BeforeEach(func() {
			officeID = uuid.New()
			validReq = dtos.UpdateOfficeRequest{
				OfficeName: "Updated Office",
				OfficeType: entities.OfficeTypeEVM,
				Address:    "Updated Address",
				IsActive:   true,
			}
		})

		Context("when office is updated successfully", func() {
			It("should return 204 no content", func() {
				mockService.EXPECT().Update(mock.Anything, officeID, mock.MatchedBy(func(cmd *services.UpdateOfficeCommand) bool {
					return cmd.OfficeName == validReq.OfficeName &&
						cmd.OfficeType == validReq.OfficeType &&
						cmd.Address == validReq.Address &&
						cmd.IsActive == validReq.IsActive
				})).Return(nil).Once()

				r.PUT("/offices/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/offices/"+officeID.String(), w, validReq)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when user is not admin", func() {
			It("should return 403 unauthorized role error", func() {
				r.PUT("/offices/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleEvmStaff)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/offices/"+officeID.String(), w, validReq)
				Expect(w.Code).To(Equal(http.StatusForbidden))
			})
		})

		Context("when office ID is invalid UUID", func() {
			It("should return 400 bad request", func() {
				r.PUT("/offices/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/offices/invalid-uuid", w, validReq)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when request body is invalid JSON", func() {
			It("should return 400 bad request", func() {
				r.PUT("/offices/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/offices/"+officeID.String(), w, "invalid json")
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidJsonRequest)
			})
		})

		Context("when office is not found", func() {
			It("should return 404 not found", func() {
				notFoundErr := apperrors.NewOfficeNotFound()
				mockService.EXPECT().Update(mock.Anything, officeID, mock.Anything).Return(notFoundErr).Once()

				r.PUT("/offices/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/offices/"+officeID.String(), w, validReq)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when service returns invalid office type error", func() {
			It("should return 400 bad request", func() {
				invalidTypeErr := apperrors.NewInvalidOfficeType()
				mockService.EXPECT().Update(mock.Anything, officeID, mock.Anything).Return(invalidTypeErr).Once()

				r.PUT("/offices/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/offices/"+officeID.String(), w, validReq)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when updating to SC office type", func() {
			It("should update successfully", func() {
				scReq := validReq
				scReq.OfficeType = entities.OfficeTypeSC

				mockService.EXPECT().Update(mock.Anything, officeID, mock.MatchedBy(func(cmd *services.UpdateOfficeCommand) bool {
					return cmd.OfficeType == entities.OfficeTypeSC
				})).Return(nil).Once()

				r.PUT("/offices/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/offices/"+officeID.String(), w, scReq)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when updating to inactive status", func() {
			It("should update successfully", func() {
				inactiveReq := validReq
				inactiveReq.IsActive = false

				mockService.EXPECT().Update(mock.Anything, officeID, mock.MatchedBy(func(cmd *services.UpdateOfficeCommand) bool {
					return cmd.IsActive == false
				})).Return(nil).Once()

				r.PUT("/offices/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/offices/"+officeID.String(), w, inactiveReq)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when service returns database error", func() {
			It("should return 500 internal server error", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().Update(mock.Anything, officeID, mock.Anything).Return(dbErr).Once()

				r.PUT("/offices/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/offices/"+officeID.String(), w, validReq)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("Delete", func() {
		var officeID uuid.UUID

		BeforeEach(func() {
			officeID = uuid.New()
		})

		Context("when office is deleted successfully", func() {
			It("should return 204 no content", func() {
				mockService.EXPECT().DeleteByID(mock.Anything, officeID).Return(nil).Once()

				r.DELETE("/offices/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Delete(c)
				})

				SendRequest(r, http.MethodDelete, "/offices/"+officeID.String(), w, nil)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when user is not admin", func() {
			It("should return 403 unauthorized role error", func() {
				r.DELETE("/offices/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScTechnician)
					SetContentTypeJSON(c)
					handler.Delete(c)
				})

				SendRequest(r, http.MethodDelete, "/offices/"+officeID.String(), w, nil)
				ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
			})
		})

		Context("when office ID is invalid UUID", func() {
			It("should return 400 bad request", func() {
				r.DELETE("/offices/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Delete(c)
				})

				SendRequest(r, http.MethodDelete, "/offices/invalid-uuid", w, nil)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when office is not found", func() {
			It("should return 404 not found", func() {
				notFoundErr := apperrors.NewOfficeNotFound()
				mockService.EXPECT().DeleteByID(mock.Anything, officeID).Return(notFoundErr).Once()

				r.DELETE("/offices/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Delete(c)
				})

				SendRequest(r, http.MethodDelete, "/offices/"+officeID.String(), w, nil)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when service returns database error", func() {
			It("should return 500 internal server error", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().DeleteByID(mock.Anything, officeID).Return(dbErr).Once()

				r.DELETE("/offices/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Delete(c)
				})

				SendRequest(r, http.MethodDelete, "/offices/"+officeID.String(), w, nil)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})

func CreateOfficeFromRequest(req dtos.CreateOfficeRequest) *entities.Office {
	return entities.NewOffice(req.OfficeName, req.OfficeType, req.Address, req.IsActive)
}
