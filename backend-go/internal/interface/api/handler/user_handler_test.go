package handler_test

import (
	"errors"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/internal/interface/api/dto"
	"ev-warranty-go/internal/interface/api/handler"
	apperrors2 "ev-warranty-go/pkg/apperror"
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
		userHandler handler.UserHandler
		r           *gin.Engine
		w           *httptest.ResponseRecorder
		validReq    dto.CreateUserRequest
		sampleUser  *entity.User
	)

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
		mockService = mocks.NewUserService(GinkgoT())
		userHandler = handler.NewUserHandler(mockLogger, mockService)

		officeID := uuid.New()
		validReq = dto.CreateUserRequest{
			Name:     "John Doe",
			Email:    "john.doe@example.com",
			Role:     entity.UserRoleAdmin,
			Password: "Password123!",
			IsActive: true,
			OfficeID: officeID,
		}
		sampleUser = CreateUserFromRequest(validReq)
	})

	Describe("Create", func() {
		Context("when authorized as ADMIN", func() {
			BeforeEach(func() {
				setupRoute("POST", "/users", entity.UserRoleAdmin, userHandler.Create)
			})

			It("should create user successfully", func() {
				mockService.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cmd *services.UserCreateCommand) bool {
					return cmd.Name == validReq.Name && cmd.Email == validReq.Email && cmd.Role == validReq.Role
				})).Return(sampleUser, nil).Once()

				SendRequest(r, http.MethodPost, "/users", w, validReq)
				ExpectResponseNotNil(w, http.StatusCreated)
			})

			DescribeTable("should handle validation errors",
				func(modifyReq func(*dto.CreateUserRequest), expectedError string) {
					req := validReq
					modifyReq(&req)
					SendRequest(r, http.MethodPost, "/users", w, req)
					ExpectErrorCode(w, http.StatusBadRequest, expectedError)
				},
				Entry("invalid role", func(req *dto.CreateUserRequest) {
					req.Role = "INVALID_ROLE"
					mockService.EXPECT().Create(mock.Anything, mock.Anything).
						Return(nil, apperrors2.NewInvalidUserInput()).Once()
				}, apperrors2.ErrorCodeInvalidUserInput),
			)

			It("should handle invalid JSON", func() {
				SendRequest(r, http.MethodPost, "/users", w, "invalid json")
				ExpectErrorCode(w, http.StatusBadRequest, apperrors2.ErrorCodeInvalidJsonRequest)
			})

			It("should handle validation errors from service", func() {
				mockService.EXPECT().Create(mock.Anything, mock.Anything).
					Return(nil, apperrors2.NewInvalidUserInput()).Once()

				SendRequest(r, http.MethodPost, "/users", w, validReq)
				ExpectErrorCode(w, http.StatusBadRequest, apperrors2.ErrorCodeInvalidUserInput)
			})

			It("should handle email already exists error", func() {
				mockService.EXPECT().Create(mock.Anything, mock.Anything).
					Return(nil, apperrors2.NewDBDuplicateKeyError("email")).Once()

				SendRequest(r, http.MethodPost, "/users", w, validReq)
				ExpectErrorCode(w, http.StatusConflict, apperrors2.ErrorCodeDuplicateKey)
			})

			It("should handle service errors", func() {
				mockService.EXPECT().Create(mock.Anything, mock.Anything).
					Return(nil, apperrors2.NewDBOperationError(errors.New("database error"))).Once()

				SendRequest(r, http.MethodPost, "/users", w, validReq)
				ExpectErrorCode(w, http.StatusInternalServerError, apperrors2.ErrorCodeDBOperation)
			})
		})

		It("should deny access for non-admin users", func() {
			setupRoute("POST", "/users", entity.UserRoleScStaff, userHandler.Create)
			SendRequest(r, http.MethodPost, "/users", w, validReq)
			ExpectErrorCode(w, http.StatusForbidden, apperrors2.ErrorCodeUnauthorizedRole)
		})
	})

	Describe("GetAll", func() {
		BeforeEach(func() {
			setupRoute("GET", "/users", entity.UserRoleAdmin, userHandler.GetAll)
		})

		DescribeTable("should handle different scenarios",
			func(setupMock func(), expectedStatus int, expectedError string) {
				setupMock()
				SendRequest(r, http.MethodGet, "/users", w, nil)

				if expectedError != "" {
					ExpectErrorCode(w, expectedStatus, expectedError)
				} else {
					ExpectResponseNotNil(w, expectedStatus)
				}
			},
			Entry("successful retrieval",
				func() {
					users := []*entity.User{sampleUser}
					mockService.EXPECT().GetAll(mock.Anything).Return(users, nil).Once()
				},
				http.StatusOK, ""),
			Entry("empty results",
				func() {
					mockService.EXPECT().GetAll(mock.Anything).Return([]*entity.User{}, nil).Once()
				},
				http.StatusOK, ""),
			Entry("service error",
				func() {
					mockService.EXPECT().GetAll(mock.Anything).
						Return(nil, errors.New("database error")).Once()
				},
				http.StatusInternalServerError, apperrors2.ErrorCodeInternalServerError),
		)
	})

	Describe("GetByID", func() {
		userID := uuid.New()

		BeforeEach(func() {
			setupRoute("GET", "/users/:id", entity.UserRoleAdmin, userHandler.GetByID)
		})

		DescribeTable("should handle different scenarios",
			func(setupMock func(), url string, expectedStatus int, expectedError string) {
				if setupMock != nil {
					setupMock()
				}
				SendRequest(r, http.MethodGet, url, w, nil)

				if expectedError != "" {
					ExpectErrorCode(w, expectedStatus, expectedError)
				} else {
					ExpectResponseNotNil(w, expectedStatus)
				}
			},
			Entry("successful retrieval",
				func() {
					mockService.EXPECT().GetByID(mock.Anything, userID).Return(sampleUser, nil).Once()
				},
				"/users/"+userID.String(), http.StatusOK, ""),
			Entry("invalid UUID",
				nil,
				"/users/invalid-uuid", http.StatusBadRequest, apperrors2.ErrorCodeInvalidUUID),
			Entry("user not found",
				func() {
					mockService.EXPECT().GetByID(mock.Anything, userID).
						Return(nil, apperrors2.NewUserNotFound()).Once()
				},
				"/users/"+userID.String(), http.StatusNotFound, apperrors2.ErrorCodeUserNotFound),
		)
	})

	Describe("Update", func() {
		userID := uuid.New()
		updateReq := dto.UpdateUserRequest{
			Name:     "Updated Name",
			Role:     entity.UserRoleEvmStaff,
			IsActive: true,
		}

		Context("when authorized as ADMIN", func() {
			BeforeEach(func() {
				setupRoute("PUT", "/users/:id", entity.UserRoleAdmin, userHandler.Update)
			})

			It("should update user successfully", func() {
				updatedUser := *sampleUser
				updatedUser.Name = updateReq.Name
				updatedUser.Role = updateReq.Role

				mockService.EXPECT().Update(mock.Anything, userID, mock.MatchedBy(func(cmd *services.UserUpdateCommand) bool {
					return cmd.Name == updateReq.Name && cmd.Role == updateReq.Role
				})).Return(nil).Once()

				SendRequest(r, http.MethodPut, "/users/"+userID.String(), w, updateReq)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})

			DescribeTable("should handle error scenarios",
				func(setupMock func(), url string, reqBody interface{}, expectedStatus int, expectedError string) {
					if setupMock != nil {
						setupMock()
					}
					SendRequest(r, http.MethodPut, url, w, reqBody)
					ExpectErrorCode(w, expectedStatus, expectedError)
				},
				Entry("invalid UUID",
					nil,
					"/users/invalid-uuid", updateReq, http.StatusBadRequest, apperrors2.ErrorCodeInvalidUUID),
				Entry("invalid JSON",
					nil,
					"/users/"+userID.String(), "invalid json", http.StatusBadRequest, apperrors2.ErrorCodeInvalidJsonRequest),
				Entry("validation error from service",
					func() {
						mockService.EXPECT().Update(mock.Anything, userID, mock.Anything).
							Return(apperrors2.NewInvalidUserInput()).Once()
					},
					"/users/"+userID.String(), updateReq, http.StatusBadRequest, apperrors2.ErrorCodeInvalidUserInput),
				Entry("user not found",
					func() {
						mockService.EXPECT().Update(mock.Anything, userID, mock.Anything).
							Return(apperrors2.NewUserNotFound()).Once()
					},
					"/users/"+userID.String(), updateReq, http.StatusNotFound, apperrors2.ErrorCodeUserNotFound),
			)
		})

		Context("when not authorized", func() {
			It("should deny access for non-admin users", func() {
				setupRoute("PUT", "/users/:id", entity.UserRoleScStaff, userHandler.Update)
				SendRequest(r, http.MethodPut, "/users/"+userID.String(), w, updateReq)
				ExpectErrorCode(w, http.StatusForbidden, apperrors2.ErrorCodeUnauthorizedRole)
			})
		})
	})

	Describe("Delete", func() {
		userID := uuid.New()

		Context("when performed by ADMIN", func() {
			BeforeEach(func() {
				setupRoute("DELETE", "/users/:id", entity.UserRoleAdmin, userHandler.Delete)
			})

			It("should delete user successfully", func() {
				mockService.EXPECT().Delete(mock.Anything, userID).Return(nil).Once()
				SendRequest(r, http.MethodDelete, "/users/"+userID.String(), w, nil)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})

			DescribeTable("should handle error scenarios",
				func(setupMock func(), url string, expectedStatus int, expectedError string) {
					if setupMock != nil {
						setupMock()
					}
					SendRequest(r, http.MethodDelete, url, w, nil)
					ExpectErrorCode(w, expectedStatus, expectedError)
				},
				Entry("invalid UUID",
					nil,
					"/users/invalid-uuid", http.StatusBadRequest, apperrors2.ErrorCodeInvalidUUID),
				Entry("service error",
					func() {
						mockService.EXPECT().Delete(mock.Anything, userID).
							Return(apperrors2.NewDBOperationError(errors.New("database error"))).Once()
					},
					"/users/"+userID.String(), http.StatusInternalServerError, apperrors2.ErrorCodeDBOperation),
				Entry("user not found",
					func() {
						mockService.EXPECT().Delete(mock.Anything, userID).
							Return(apperrors2.NewUserNotFound()).Once()
					},
					"/users/"+userID.String(), http.StatusNotFound, apperrors2.ErrorCodeUserNotFound),
			)
		})

		Context("when not authorized", func() {
			It("should deny access for non-admin users", func() {
				setupRoute("DELETE", "/users/:id", entity.UserRoleScStaff, userHandler.Delete)
				SendRequest(r, http.MethodDelete, "/users/"+userID.String(), w, nil)
				ExpectErrorCode(w, http.StatusForbidden, apperrors2.ErrorCodeUnauthorizedRole)
			})
		})
	})
})

func CreateUserFromRequest(req dto.CreateUserRequest) *entity.User {
	return entity.NewUser(req.Name, req.Email, req.Role, "hashed_password", req.IsActive, req.OfficeID)
}
