package handlers_test

import (
	"bytes"
	"encoding/json"
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
		router      *gin.Engine
		w           *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		mockLogger = mocks.NewLogger(GinkgoT())
		mockService = mocks.NewOfficeService(GinkgoT())
		handler = handlers.NewOfficeHandler(mockLogger, mockService)
		router = gin.New()
		w = httptest.NewRecorder()

		mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
		mockLogger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
		mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
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
				office := &entities.Office{
					ID:         uuid.New(),
					OfficeName: validReq.OfficeName,
					OfficeType: validReq.OfficeType,
					Address:    validReq.Address,
					IsActive:   validReq.IsActive,
				}

				mockService.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cmd *services.CreateOfficeCommand) bool {
					return cmd.OfficeName == validReq.OfficeName &&
						cmd.OfficeType == validReq.OfficeType &&
						cmd.Address == validReq.Address &&
						cmd.IsActive == validReq.IsActive
				})).Return(office, nil).Once()

				router.POST("/offices", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Create(c)
				})

				body, _ := json.Marshal(validReq)
				req, _ := http.NewRequest(http.MethodPost, "/offices", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusCreated))
				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Data).NotTo(BeNil())
			})
		})

		Context("when user is not admin", func() {
			It("should return 403 unauthorized role error", func() {
				router.POST("/offices", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleScStaff)
					handler.Create(c)
				})

				body, _ := json.Marshal(validReq)
				req, _ := http.NewRequest(http.MethodPost, "/offices", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleScStaff)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusForbidden))
				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Error).To(Equal(apperrors.ErrorCodeUnauthorizedRole))
			})
		})

		Context("when request body is invalid JSON", func() {
			It("should return 400 bad request", func() {
				router.POST("/offices", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Create(c)
				})

				req, _ := http.NewRequest(http.MethodPost, "/offices", bytes.NewBuffer([]byte("invalid json")))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Error).To(Equal(apperrors.ErrorCodeInvalidJsonRequest))
			})
		})

		Context("when office type is invalid", func() {
			It("should return 400 bad request", func() {
				invalidReq := validReq
				invalidReq.OfficeType = "INVALID_TYPE"

				router.POST("/offices", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Create(c)
				})

				body, _ := json.Marshal(invalidReq)
				req, _ := http.NewRequest(http.MethodPost, "/offices", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Error).To(Equal(apperrors.ErrorCodeInvalidOfficeType))
			})
		})

		Context("when service returns error", func() {
			It("should return error from service", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, dbErr).Once()

				router.POST("/offices", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Create(c)
				})

				body, _ := json.Marshal(validReq)
				req, _ := http.NewRequest(http.MethodPost, "/offices", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Error).To(Equal(apperrors.ErrorCodeDBOperation))
			})
		})

		Context("when required fields are missing", func() {
			It("should return 400 bad request for missing office_name", func() {
				invalidReq := validReq
				invalidReq.OfficeName = ""

				router.POST("/offices", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Create(c)
				})

				body, _ := json.Marshal(invalidReq)
				req, _ := http.NewRequest(http.MethodPost, "/offices", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when office type is SC", func() {
			It("should create office with SC type successfully", func() {
				scReq := validReq
				scReq.OfficeType = entities.OfficeTypeSC

				office := &entities.Office{
					ID:         uuid.New(),
					OfficeName: scReq.OfficeName,
					OfficeType: scReq.OfficeType,
					Address:    scReq.Address,
					IsActive:   scReq.IsActive,
				}

				mockService.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cmd *services.CreateOfficeCommand) bool {
					return cmd.OfficeType == entities.OfficeTypeSC
				})).Return(office, nil).Once()

				router.POST("/offices", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Create(c)
				})

				body, _ := json.Marshal(scReq)
				req, _ := http.NewRequest(http.MethodPost, "/offices", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusCreated))
			})
		})

		Context("when office is created as inactive", func() {
			It("should create inactive office successfully", func() {
				inactiveReq := validReq
				inactiveReq.IsActive = false

				office := &entities.Office{
					ID:         uuid.New(),
					OfficeName: inactiveReq.OfficeName,
					OfficeType: inactiveReq.OfficeType,
					Address:    inactiveReq.Address,
					IsActive:   false,
				}

				mockService.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cmd *services.CreateOfficeCommand) bool {
					return cmd.IsActive == false
				})).Return(office, nil).Once()

				router.POST("/offices", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Create(c)
				})

				body, _ := json.Marshal(inactiveReq)
				req, _ := http.NewRequest(http.MethodPost, "/offices", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

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

				router.GET("/offices/:id", handler.GetById)

				req, _ := http.NewRequest(http.MethodGet, "/offices/"+officeID.String(), nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Data).NotTo(BeNil())
			})
		})

		Context("when office ID is invalid UUID", func() {
			It("should return 400 bad request", func() {
				router.GET("/offices/:id", handler.GetById)

				req, _ := http.NewRequest(http.MethodGet, "/offices/invalid-uuid", nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Error).To(Equal(apperrors.ErrorCodeInvalidUUID))
			})
		})

		Context("when office is not found", func() {
			It("should return 404 not found", func() {
				notFoundErr := apperrors.NewOfficeNotFound()
				mockService.EXPECT().GetByID(mock.Anything, officeID).Return(nil, notFoundErr).Once()

				router.GET("/offices/:id", handler.GetById)

				req, _ := http.NewRequest(http.MethodGet, "/offices/"+officeID.String(), nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Error).To(Equal(apperrors.ErrorCodeOfficeNotFound))
			})
		})

		Context("when service returns database error", func() {
			It("should return 500 internal server error", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().GetByID(mock.Anything, officeID).Return(nil, dbErr).Once()

				router.GET("/offices/:id", handler.GetById)

				req, _ := http.NewRequest(http.MethodGet, "/offices/"+officeID.String(), nil)
				router.ServeHTTP(w, req)

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

				router.GET("/offices", handler.GetAll)

				req, _ := http.NewRequest(http.MethodGet, "/offices", nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Data).NotTo(BeNil())
			})
		})

		Context("when no offices exist", func() {
			It("should return 200 with empty array", func() {
				mockService.EXPECT().GetAll(mock.Anything).Return([]*entities.Office{}, nil).Once()

				router.GET("/offices", handler.GetAll)

				req, _ := http.NewRequest(http.MethodGet, "/offices", nil)
				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})

		Context("when service returns error", func() {
			It("should return 500 internal server error", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().GetAll(mock.Anything).Return(nil, dbErr).Once()

				router.GET("/offices", handler.GetAll)

				req, _ := http.NewRequest(http.MethodGet, "/offices", nil)
				router.ServeHTTP(w, req)

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

				router.PUT("/offices/:id", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Update(c)
				})

				body, _ := json.Marshal(validReq)
				req, _ := http.NewRequest(http.MethodPut, "/offices/"+officeID.String(), bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when user is not admin", func() {
			It("should return 403 unauthorized role error", func() {
				router.PUT("/offices/:id", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleEvmStaff)
					handler.Update(c)
				})

				body, _ := json.Marshal(validReq)
				req, _ := http.NewRequest(http.MethodPut, "/offices/"+officeID.String(), bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleEvmStaff)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusForbidden))
			})
		})

		Context("when office ID is invalid UUID", func() {
			It("should return 400 bad request", func() {
				router.PUT("/offices/:id", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Update(c)
				})

				body, _ := json.Marshal(validReq)
				req, _ := http.NewRequest(http.MethodPut, "/offices/invalid-uuid", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Error).To(Equal(apperrors.ErrorCodeInvalidUUID))
			})
		})

		Context("when request body is invalid JSON", func() {
			It("should return 400 bad request", func() {
				router.PUT("/offices/:id", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Update(c)
				})

				req, _ := http.NewRequest(http.MethodPut, "/offices/"+officeID.String(), bytes.NewBuffer([]byte("invalid json")))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Error).To(Equal(apperrors.ErrorCodeInvalidJsonRequest))
			})
		})

		Context("when office is not found", func() {
			It("should return 404 not found", func() {
				notFoundErr := apperrors.NewOfficeNotFound()
				mockService.EXPECT().Update(mock.Anything, officeID, mock.Anything).Return(notFoundErr).Once()

				router.PUT("/offices/:id", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Update(c)
				})

				body, _ := json.Marshal(validReq)
				req, _ := http.NewRequest(http.MethodPut, "/offices/"+officeID.String(), bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when service returns invalid office type error", func() {
			It("should return 400 bad request", func() {
				invalidTypeErr := apperrors.NewInvalidOfficeType()
				mockService.EXPECT().Update(mock.Anything, officeID, mock.Anything).Return(invalidTypeErr).Once()

				router.PUT("/offices/:id", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Update(c)
				})

				body, _ := json.Marshal(validReq)
				req, _ := http.NewRequest(http.MethodPut, "/offices/"+officeID.String(), bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

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

				router.PUT("/offices/:id", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Update(c)
				})

				body, _ := json.Marshal(scReq)
				req, _ := http.NewRequest(http.MethodPut, "/offices/"+officeID.String(), bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

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

				router.PUT("/offices/:id", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Update(c)
				})

				body, _ := json.Marshal(inactiveReq)
				req, _ := http.NewRequest(http.MethodPut, "/offices/"+officeID.String(), bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when service returns database error", func() {
			It("should return 500 internal server error", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().Update(mock.Anything, officeID, mock.Anything).Return(dbErr).Once()

				router.PUT("/offices/:id", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Update(c)
				})

				body, _ := json.Marshal(validReq)
				req, _ := http.NewRequest(http.MethodPut, "/offices/"+officeID.String(), bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

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

				router.DELETE("/offices/:id", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Delete(c)
				})

				req, _ := http.NewRequest(http.MethodDelete, "/offices/"+officeID.String(), nil)
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when user is not admin", func() {
			It("should return 403 unauthorized role error", func() {
				router.DELETE("/offices/:id", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleScTechnician)
					handler.Delete(c)
				})

				req, _ := http.NewRequest(http.MethodDelete, "/offices/"+officeID.String(), nil)
				req.Header.Set("X-User-Role", entities.UserRoleScTechnician)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusForbidden))
				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Error).To(Equal(apperrors.ErrorCodeUnauthorizedRole))
			})
		})

		Context("when office ID is invalid UUID", func() {
			It("should return 400 bad request", func() {
				router.DELETE("/offices/:id", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Delete(c)
				})

				req, _ := http.NewRequest(http.MethodDelete, "/offices/invalid-uuid", nil)
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				var response dtos.APIResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				Expect(err).NotTo(HaveOccurred())
				Expect(response.Error).To(Equal(apperrors.ErrorCodeInvalidUUID))
			})
		})

		Context("when office is not found", func() {
			It("should return 404 not found", func() {
				notFoundErr := apperrors.NewOfficeNotFound()
				mockService.EXPECT().DeleteByID(mock.Anything, officeID).Return(notFoundErr).Once()

				router.DELETE("/offices/:id", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Delete(c)
				})

				req, _ := http.NewRequest(http.MethodDelete, "/offices/"+officeID.String(), nil)
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when service returns database error", func() {
			It("should return 500 internal server error", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().DeleteByID(mock.Anything, officeID).Return(dbErr).Once()

				router.DELETE("/offices/:id", func(c *gin.Context) {
					c.Request.Header.Set("X-User-Role", entities.UserRoleAdmin)
					handler.Delete(c)
				})

				req, _ := http.NewRequest(http.MethodDelete, "/offices/"+officeID.String(), nil)
				req.Header.Set("X-User-Role", entities.UserRoleAdmin)

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
