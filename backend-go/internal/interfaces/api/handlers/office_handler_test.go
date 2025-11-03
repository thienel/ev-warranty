package handlers_test

import (
	"errors"
	"ev-warranty-go/internal/apperrors"
	"ev-warranty-go/internal/application/services"
	"ev-warranty-go/internal/domain/entities"
	"ev-warranty-go/internal/interfaces/api/dto"
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
		mockLogger   *mocks.Logger
		mockService  *mocks.OfficeService
		handler      handlers.OfficeHandler
		r            *gin.Engine
		w            *httptest.ResponseRecorder
		validReq     dto.CreateOfficeRequest
		sampleOffice *entities.Office
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
		mockService = mocks.NewOfficeService(GinkgoT())
		handler = handlers.NewOfficeHandler(mockLogger, mockService)

		validReq = dto.CreateOfficeRequest{
			OfficeName: "Test Office",
			OfficeType: entities.OfficeTypeEVM,
			Address:    "123 Test Street",
			IsActive:   true,
		}
		sampleOffice = CreateOfficeFromRequest(validReq)
	})

	Describe("Create", func() {
		Context("when authorized as ADMIN", func() {
			BeforeEach(func() {
				setupRoute("POST", "/offices", entities.UserRoleAdmin, handler.Create)
			})

			It("should create office successfully", func() {
				mockService.EXPECT().Create(mock.Anything, mock.MatchedBy(func(cmd *services.CreateOfficeCommand) bool {
					return cmd.OfficeName == validReq.OfficeName && cmd.OfficeType == validReq.OfficeType
				})).Return(sampleOffice, nil).Once()

				SendRequest(r, http.MethodPost, "/offices", w, validReq)
				ExpectResponseNotNil(w, http.StatusCreated)
			})

			DescribeTable("should handle validation errors",
				func(modifyReq func(*dto.CreateOfficeRequest), expectedError string) {
					req := validReq
					if modifyReq != nil {
						modifyReq(&req)
					}
					SendRequest(r, http.MethodPost, "/offices", w, req)
					ExpectErrorCode(w, http.StatusBadRequest, expectedError)
				},
				Entry("invalid office type",
					func(req *dto.CreateOfficeRequest) {
						req.OfficeType = "INVALID_TYPE"
					},
					apperrors.ErrorCodeInvalidOfficeType),
			)

			It("should handle invalid JSON", func() {
				SendRequest(r, http.MethodPost, "/offices", w, "invalid json")
				ExpectErrorCode(w, http.StatusBadRequest, apperrors.ErrorCodeInvalidJsonRequest)
			})

			It("should handle service errors", func() {
				mockService.EXPECT().Create(mock.Anything, mock.Anything).
					Return(nil, apperrors.NewDBOperationError(errors.New("database error"))).Once()

				SendRequest(r, http.MethodPost, "/offices", w, validReq)
				ExpectErrorCode(w, http.StatusInternalServerError, apperrors.ErrorCodeDBOperation)
			})
		})

		It("should deny access for non-admin users", func() {
			setupRoute("POST", "/offices", entities.UserRoleScStaff, handler.Create)
			SendRequest(r, http.MethodPost, "/offices", w, validReq)
			ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
		})
	})

	Describe("GetAll", func() {
		BeforeEach(func() {
			setupRoute("GET", "/offices", "", handler.GetAll)
		})

		DescribeTable("should handle different scenarios",
			func(setupMock func(), expectedStatus int, expectedError string) {
				setupMock()
				SendRequest(r, http.MethodGet, "/offices", w, nil)

				if expectedError != "" {
					ExpectErrorCode(w, expectedStatus, expectedError)
				} else {
					ExpectResponseNotNil(w, expectedStatus)
				}
			},
			Entry("successful retrieval",
				func() {
					offices := []*entities.Office{sampleOffice}
					mockService.EXPECT().GetAll(mock.Anything).Return(offices, nil).Once()
				},
				http.StatusOK, ""),
			Entry("empty results",
				func() {
					mockService.EXPECT().GetAll(mock.Anything).Return([]*entities.Office{}, nil).Once()
				},
				http.StatusOK, ""),
			Entry("service error",
				func() {
					mockService.EXPECT().GetAll(mock.Anything).
						Return(nil, errors.New("database error")).Once()
				},
				http.StatusInternalServerError, apperrors.ErrorCodeInternalServerError),
		)
	})

	Describe("GetByID", func() {
		officeID := uuid.New()

		BeforeEach(func() {
			setupRoute("GET", "/offices/:id", "", handler.GetByID)
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
					mockService.EXPECT().GetByID(mock.Anything, officeID).Return(sampleOffice, nil).Once()
				},
				"/offices/"+officeID.String(), http.StatusOK, ""),
			Entry("invalid UUID",
				nil,
				"/offices/invalid-uuid", http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID),
			Entry("office not found",
				func() {
					mockService.EXPECT().GetByID(mock.Anything, officeID).
						Return(nil, apperrors.NewOfficeNotFound()).Once()
				},
				"/offices/"+officeID.String(), http.StatusNotFound, apperrors.ErrorCodeOfficeNotFound),
		)
	})

	Describe("Update", func() {
		officeID := uuid.New()
		updateReq := dto.UpdateOfficeRequest{
			OfficeName: "Updated Office",
			Address:    "456 Updated Street",
			IsActive:   false,
		}

		Context("when authorized as ADMIN", func() {
			BeforeEach(func() {
				setupRoute("PUT", "/offices/:id", entities.UserRoleAdmin, handler.Update)
			})

			It("should update office successfully", func() {
				updatedOffice := *sampleOffice
				updatedOffice.OfficeName = updateReq.OfficeName

				mockService.EXPECT().Update(mock.Anything, officeID, mock.MatchedBy(func(cmd *services.UpdateOfficeCommand) bool {
					return cmd.OfficeName == updateReq.OfficeName && cmd.Address == updateReq.Address
				})).Return(nil).Once()

				SendRequest(r, http.MethodPut, "/offices/"+officeID.String(), w, updateReq)
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
					"/offices/invalid-uuid", updateReq, http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID),
				Entry("invalid JSON",
					nil,
					"/offices/"+officeID.String(), "invalid json", http.StatusBadRequest, apperrors.ErrorCodeInvalidJsonRequest),
				Entry("office not found",
					func() {
						mockService.EXPECT().Update(mock.Anything, officeID, mock.Anything).
							Return(apperrors.NewOfficeNotFound()).Once()
					},
					"/offices/"+officeID.String(), updateReq, http.StatusNotFound, apperrors.ErrorCodeOfficeNotFound),
			)
		})

		Context("when not authorized as ADMIN", func() {
			It("should deny access", func() {
				setupRoute("PUT", "/offices/:id", entities.UserRoleScStaff, handler.Update)
				SendRequest(r, http.MethodPut, "/offices/"+officeID.String(), w, updateReq)
				ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
			})
		})
	})

	Describe("Delete", func() {
		officeID := uuid.New()

		Context("when authorized as ADMIN", func() {
			BeforeEach(func() {
				setupRoute("DELETE", "/offices/:id", entities.UserRoleAdmin, handler.Delete)
			})

			It("should delete office successfully", func() {
				mockService.EXPECT().DeleteByID(mock.Anything, officeID).Return(nil).Once()
				SendRequest(r, http.MethodDelete, "/offices/"+officeID.String(), w, nil)
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
					"/offices/invalid-uuid", http.StatusBadRequest, apperrors.ErrorCodeInvalidUUID),
				Entry("office not found",
					func() {
						mockService.EXPECT().DeleteByID(mock.Anything, officeID).
							Return(apperrors.NewOfficeNotFound()).Once()
					},
					"/offices/"+officeID.String(), http.StatusNotFound, apperrors.ErrorCodeOfficeNotFound),
			)
		})

		Context("when not authorized as ADMIN", func() {
			It("should deny access", func() {
				setupRoute("DELETE", "/offices/:id", entities.UserRoleScStaff, handler.Delete)
				SendRequest(r, http.MethodDelete, "/offices/"+officeID.String(), w, nil)
				ExpectErrorCode(w, http.StatusForbidden, apperrors.ErrorCodeUnauthorizedRole)
			})
		})
	})
})

func CreateOfficeFromRequest(req dto.CreateOfficeRequest) *entities.Office {
	return entities.NewOffice(req.OfficeName, req.OfficeType, req.Address, req.IsActive)
}
