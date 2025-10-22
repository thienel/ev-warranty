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

var _ = Describe("UserHandler", func() {
	var (
		mockLogger  *mocks.Logger
		mockService *mocks.UserService
		handler     handlers.UserHandler
		r           *gin.Engine
		w           *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		mockLogger, r, w = SetupMock(GinkgoT())
		mockService = mocks.NewUserService(GinkgoT())
		handler = handlers.NewUserHandler(mockLogger, mockService)
	})

	Describe("Create", func() {
		var (
			validReq dtos.CreateUserRequest
			officeID uuid.UUID
		)

		BeforeEach(func() {
			officeID = uuid.New()
			validReq = dtos.CreateUserRequest{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Role:     entities.UserRoleAdmin,
				Password: "Password123!",
				IsActive: true,
				OfficeID: officeID,
			}
		})

		Context("when user is created successfully", func() {
			It("should return 201 with created user", func() {
				user := CreateUserFromRequest(validReq)

				mockService.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cmd *services.UserCreateCommand) bool {
					return cmd.Name == validReq.Name &&
						cmd.Email == validReq.Email &&
						cmd.Role == validReq.Role &&
						cmd.Password == validReq.Password &&
						cmd.IsActive == validReq.IsActive &&
						cmd.OfficeID == validReq.OfficeID
				})).Return(user, nil).Once()

				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, validReq)
				ExpectResponseNotNil(w, http.StatusCreated)
			})
		})

		Context("when user is not admin", func() {
			It("should return 403 forbidden role error", func() {
				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScStaff)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, validReq)
				ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
			})
		})

		Context("when request body is invalid JSON", func() {
			It("should return 400 bad request", func() {
				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, "invalid json")
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidJsonRequest)
			})
		})

		Context("when user role is invalid", func() {
			It("should return 400 bad request", func() {
				invalidReq := validReq
				invalidReq.Role = "INVALID_ROLE"

				mockService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, apperrors.NewInvalidUserInput()).Once()
				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, invalidReq)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUserInput)
			})
		})

		Context("when service returns error", func() {
			It("should return error from service", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, dbErr).Once()

				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, validReq)
				ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeDBOperation)
			})
		})

		Context("when required fields are missing", func() {
			It("should return 400 bad request for missing name", func() {
				invalidReq := validReq
				invalidReq.Name = ""

				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, invalidReq)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 bad request for missing email", func() {
				invalidReq := validReq
				invalidReq.Email = ""

				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, invalidReq)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 400 bad request for missing password", func() {
				invalidReq := validReq
				invalidReq.Password = ""

				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, invalidReq)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when email is invalid format", func() {
			It("should return 400 bad request", func() {
				invalidReq := validReq
				invalidReq.Email = "invalid-email"

				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, invalidReq)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidJsonRequest)
			})
		})

		Context("when password is weak", func() {
			It("should return 400 bad request", func() {
				invalidReq := validReq
				invalidReq.Password = "weak"

				invalidInputErr := apperrors.NewInvalidUserInput()
				mockService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, invalidInputErr).Once()

				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, invalidReq)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUserInput)
			})
		})

		Context("when user is created with EVM staff role", func() {
			It("should create user successfully", func() {
				evmReq := validReq
				evmReq.Role = entities.UserRoleEvmStaff
				user := CreateUserFromRequest(evmReq)

				mockService.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cmd *services.UserCreateCommand) bool {
					return cmd.Role == entities.UserRoleEvmStaff
				})).Return(user, nil).Once()

				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, evmReq)
				Expect(w.Code).To(Equal(http.StatusCreated))
			})
		})

		Context("when user is created with SC staff role", func() {
			It("should create user successfully", func() {
				scReq := validReq
				scReq.Role = entities.UserRoleScStaff
				user := CreateUserFromRequest(scReq)

				mockService.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cmd *services.UserCreateCommand) bool {
					return cmd.Role == entities.UserRoleScStaff
				})).Return(user, nil).Once()

				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, scReq)
				Expect(w.Code).To(Equal(http.StatusCreated))
			})
		})

		Context("when user is created as inactive", func() {
			It("should create inactive user successfully", func() {
				inactiveReq := validReq
				inactiveReq.IsActive = false
				user := CreateUserFromRequest(inactiveReq)

				mockService.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cmd *services.UserCreateCommand) bool {
					return cmd.IsActive == false
				})).Return(user, nil).Once()

				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, inactiveReq)
				Expect(w.Code).To(Equal(http.StatusCreated))
			})
		})

		Context("when office is not found", func() {
			It("should return 404 not found", func() {
				officeNotFoundErr := apperrors.NewOfficeNotFound()
				mockService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, officeNotFoundErr).Once()

				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, validReq)
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeOfficeNotFound)
			})
		})

		Context("when role does not match office type", func() {
			It("should return 400 bad request", func() {
				invalidOfficeTypeErr := apperrors.NewInvalidOfficeType()
				mockService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, invalidOfficeTypeErr).Once()

				r.POST("/users", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Create(c)
				})

				SendRequest(r, http.MethodPost, "/users", w, validReq)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidOfficeType)
			})
		})
	})

	Describe("GetByID", func() {
		var userID uuid.UUID

		BeforeEach(func() {
			userID = uuid.New()
		})

		Context("when user is found", func() {
			It("should return 200 with user data", func() {
				user := &entities.User{
					ID:       userID,
					Name:     "John Doe",
					Email:    "john.doe@example.com",
					Role:     entities.UserRoleAdmin,
					IsActive: true,
					OfficeID: uuid.New(),
				}
				mockService.EXPECT().GetByID(mock.Anything, userID).Return(user, nil).Once()

				r.GET("/users/:id", handler.GetByID)

				SendRequest(r, http.MethodGet, "/users/"+userID.String(), w, nil)
				ExpectResponseNotNil(w, http.StatusOK)
			})
		})

		Context("when user ID is invalid UUID", func() {
			It("should return 400 bad request", func() {
				r.GET("/users/:id", handler.GetByID)

				SendRequest(r, http.MethodGet, "/users/invalid-uuid", w, nil)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when user is not found", func() {
			It("should return 404 not found", func() {
				notFoundErr := apperrors.NewUserNotFound()
				mockService.EXPECT().GetByID(mock.Anything, userID).Return(nil, notFoundErr).Once()

				r.GET("/users/:id", handler.GetByID)

				SendRequest(r, http.MethodGet, "/users/"+userID.String(), w, nil)
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeUserNotFound)
			})
		})

		Context("when service returns database error", func() {
			It("should return 500 internal server error", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().GetByID(mock.Anything, userID).Return(nil, dbErr).Once()

				r.GET("/users/:id", handler.GetByID)

				SendRequest(r, http.MethodGet, "/users/"+userID.String(), w, nil)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("GetAll", func() {
		Context("when users are found", func() {
			It("should return 200 with all users", func() {
				users := []*entities.User{
					{
						ID:       uuid.New(),
						Name:     "User 1",
						Email:    "user1@example.com",
						Role:     entities.UserRoleAdmin,
						IsActive: true,
						OfficeID: uuid.New(),
					},
					{
						ID:       uuid.New(),
						Name:     "User 2",
						Email:    "user2@example.com",
						Role:     entities.UserRoleEvmStaff,
						IsActive: false,
						OfficeID: uuid.New(),
					},
				}

				mockService.EXPECT().GetAll(mock.Anything).Return(users, nil).Once()

				r.GET("/users", handler.GetAll)

				SendRequest(r, http.MethodGet, "/users", w, nil)
				ExpectResponseNotNil(w, http.StatusOK)
			})
		})

		Context("when no users exist", func() {
			It("should return 200 with empty array", func() {
				mockService.EXPECT().GetAll(mock.Anything).Return([]*entities.User{}, nil).Once()

				r.GET("/users", handler.GetAll)

				SendRequest(r, http.MethodGet, "/users", w, nil)
				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})

		Context("when service returns error", func() {
			It("should return 500 internal server error", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().GetAll(mock.Anything).Return(nil, dbErr).Once()

				r.GET("/users", handler.GetAll)

				SendRequest(r, http.MethodGet, "/users", w, nil)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("Update", func() {
		var (
			userID   uuid.UUID
			validReq dtos.UpdateUserRequest
			officeID uuid.UUID
		)

		BeforeEach(func() {
			userID = uuid.New()
			officeID = uuid.New()
			validReq = dtos.UpdateUserRequest{
				Name:     "Updated Name",
				Role:     entities.UserRoleAdmin,
				IsActive: true,
				OfficeID: officeID,
			}
		})

		Context("when user is updated successfully", func() {
			It("should return 204 no content", func() {
				mockService.EXPECT().Update(mock.Anything, userID, mock.MatchedBy(func(cmd *services.UserUpdateCommand) bool {
					return cmd.Name == validReq.Name &&
						cmd.Role == validReq.Role &&
						cmd.IsActive == validReq.IsActive &&
						cmd.OfficeID == validReq.OfficeID
				})).Return(nil).Once()

				r.PUT("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/users/"+userID.String(), w, validReq)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when user is not admin", func() {
			It("should return 403 unauthorized role error", func() {
				r.PUT("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleEvmStaff)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/users/"+userID.String(), w, validReq)
				Expect(w.Code).To(Equal(http.StatusForbidden))
			})
		})

		Context("when user ID is invalid UUID", func() {
			It("should return 400 bad request", func() {
				r.PUT("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/users/invalid-uuid", w, validReq)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when request body is invalid JSON", func() {
			It("should return 400 bad request", func() {
				r.PUT("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/users/"+userID.String(), w, "invalid json")
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidJsonRequest)
			})
		})

		Context("when user is not found", func() {
			It("should return 404 not found", func() {
				notFoundErr := apperrors.NewUserNotFound()
				mockService.EXPECT().Update(mock.Anything, userID, mock.Anything).Return(notFoundErr).Once()

				r.PUT("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/users/"+userID.String(), w, validReq)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when service returns invalid user input error", func() {
			It("should return 400 bad request", func() {
				invalidInputErr := apperrors.NewInvalidUserInput()
				mockService.EXPECT().Update(mock.Anything, userID, mock.Anything).Return(invalidInputErr).Once()

				r.PUT("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/users/"+userID.String(), w, validReq)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when updating to SC staff role", func() {
			It("should update successfully", func() {
				scReq := validReq
				scReq.Role = entities.UserRoleScStaff

				mockService.EXPECT().Update(mock.Anything, userID, mock.MatchedBy(func(cmd *services.UserUpdateCommand) bool {
					return cmd.Role == entities.UserRoleScStaff
				})).Return(nil).Once()

				r.PUT("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/users/"+userID.String(), w, scReq)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when updating to inactive status", func() {
			It("should update successfully", func() {
				inactiveReq := validReq
				inactiveReq.IsActive = false

				mockService.EXPECT().Update(mock.Anything, userID, mock.MatchedBy(func(cmd *services.UserUpdateCommand) bool {
					return cmd.IsActive == false
				})).Return(nil).Once()

				r.PUT("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/users/"+userID.String(), w, inactiveReq)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when service returns database error", func() {
			It("should return 500 internal server error", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().Update(mock.Anything, userID, mock.Anything).Return(dbErr).Once()

				r.PUT("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/users/"+userID.String(), w, validReq)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("when office is not found", func() {
			It("should return 404 not found", func() {
				officeNotFoundErr := apperrors.NewOfficeNotFound()
				mockService.EXPECT().Update(mock.Anything, userID, mock.Anything).Return(officeNotFoundErr).Once()

				r.PUT("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/users/"+userID.String(), w, validReq)
				ExpectErrorCode(w, http.StatusNotFound, apperrors.ErrorCodeOfficeNotFound)
			})
		})

		Context("when role does not match office type", func() {
			It("should return 400 bad request", func() {
				invalidOfficeTypeErr := apperrors.NewInvalidOfficeType()
				mockService.EXPECT().Update(mock.Anything, userID, mock.Anything).Return(invalidOfficeTypeErr).Once()

				r.PUT("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Update(c)
				})

				SendRequest(r, http.MethodPut, "/users/"+userID.String(), w, validReq)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidOfficeType)
			})
		})
	})

	Describe("Delete", func() {
		var userID uuid.UUID

		BeforeEach(func() {
			userID = uuid.New()
		})

		Context("when user is deleted successfully", func() {
			It("should return 204 no content", func() {
				mockService.EXPECT().Delete(mock.Anything, userID).Return(nil).Once()

				r.DELETE("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Delete(c)
				})

				SendRequest(r, http.MethodDelete, "/users/"+userID.String(), w, nil)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})

		Context("when user is not admin", func() {
			It("should return 403 unauthorized role error", func() {
				r.DELETE("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleScTechnician)
					SetContentTypeJSON(c)
					handler.Delete(c)
				})

				SendRequest(r, http.MethodDelete, "/users/"+userID.String(), w, nil)
				ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
			})
		})

		Context("when user ID is invalid UUID", func() {
			It("should return 400 bad request", func() {
				r.DELETE("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Delete(c)
				})

				SendRequest(r, http.MethodDelete, "/users/invalid-uuid", w, nil)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID)
			})
		})

		Context("when user is not found", func() {
			It("should return 404 not found", func() {
				notFoundErr := apperrors.NewUserNotFound()
				mockService.EXPECT().Delete(mock.Anything, userID).Return(notFoundErr).Once()

				r.DELETE("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Delete(c)
				})

				SendRequest(r, http.MethodDelete, "/users/"+userID.String(), w, nil)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when service returns database error", func() {
			It("should return 500 internal server error", func() {
				dbErr := apperrors.NewDBOperationError(errors.New("database error"))
				mockService.EXPECT().Delete(mock.Anything, userID).Return(dbErr).Once()

				r.DELETE("/users/:id", func(c *gin.Context) {
					SetHeaderRole(c, entities.UserRoleAdmin)
					SetContentTypeJSON(c)
					handler.Delete(c)
				})

				SendRequest(r, http.MethodDelete, "/users/"+userID.String(), w, nil)
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})

func CreateUserFromRequest(req dtos.CreateUserRequest) *entities.User {
	return entities.NewUser(req.Name, req.Email, req.Role, "hashed_password", req.IsActive, req.OfficeID)
}
