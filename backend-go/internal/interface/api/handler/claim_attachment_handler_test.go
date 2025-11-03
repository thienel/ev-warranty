package handler_test

import (
	"bytes"
	"context"
	"errors"
	"ev-warranty-go/internal/application"
	"ev-warranty-go/internal/domain/entity"
	"ev-warranty-go/internal/interface/api/handler"
	apperrors2 "ev-warranty-go/pkg/apperror"
	"ev-warranty-go/pkg/mocks"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("ClaimAttachmentHandler", func() {
	var (
		mockLogger        *mocks.Logger
		mockTxManager     *mocks.TxManager
		mockService       *mocks.ClaimAttachmentService
		mockTx            *mocks.Tx
		attachmentHandler handler.ClaimAttachmentHandler
		r                 *gin.Engine
		w                 *httptest.ResponseRecorder
		claimID           uuid.UUID
		attachmentID      uuid.UUID
		sampleAttachment  *entity.ClaimAttachment
	)

	createMultipartForm := func(filenames ...string) (*bytes.Buffer, string) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		if len(filenames) == 0 {
			filenames = []string{"test-image.jpg"}
		}

		for i, filename := range filenames {
			part, _ := writer.CreateFormFile("files", filename)
			content := []byte("fake image content " + string(rune('1'+i)))
			_, _ = part.Write(content)
		}
		_ = writer.Close()
		return body, writer.FormDataContentType()
	}

	BeforeEach(func() {
		mockLogger, r, w = SetupMock(GinkgoT())
		mockTxManager = mocks.NewTxManager(GinkgoT())
		mockService = mocks.NewClaimAttachmentService(GinkgoT())
		mockTx = mocks.NewTx(GinkgoT())
		attachmentHandler = handler.NewClaimAttachmentHandler(mockLogger, mockTxManager, mockService)

		claimID = uuid.New()
		attachmentID = uuid.New()
		sampleAttachment = &entity.ClaimAttachment{
			ID:      attachmentID,
			ClaimID: claimID,
			Type:    entity.AttachmentTypeImage,
			URL:     "https://example.com/image.jpg",
		}
	})

	Describe("GetByID", func() {
		DescribeTable("should handle different scenarios",
			func(setupMock func(), expectedStatus int, expectedError string) {
				setupMock()

				r.GET("/claims/:id/attachments/:attachmentID", attachmentHandler.GetByID)
				req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/attachments/"+attachmentID.String(), nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(expectedStatus))
				if expectedError != "" {
					ExpectErrorCode(w, expectedStatus, expectedError)
				} else {
					ExpectResponseNotNil(w, expectedStatus)
				}
			},
			Entry("successful retrieval",
				func() {
					mockService.EXPECT().GetByID(mock.Anything, attachmentID).Return(sampleAttachment, nil).Once()
				},
				http.StatusOK, ""),
			Entry("attachment not found",
				func() {
					mockService.EXPECT().GetByID(mock.Anything, attachmentID).
						Return(nil, apperrors2.NewClaimAttachmentNotFound()).Once()
				},
				http.StatusNotFound, apperrors2.ErrorCodeClaimAttachmentNotFound),
			Entry("service error",
				func() {
					mockService.EXPECT().GetByID(mock.Anything, attachmentID).
						Return(nil, errors.New("database error")).Once()
				},
				http.StatusInternalServerError, apperrors2.ErrorCodeInternalServerError),
		)

		It("should return error for invalid attachment ID", func() {
			r.GET("/claims/:id/attachments/:attachmentID", attachmentHandler.GetByID)
			req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/attachments/invalid-uuid", nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusBadRequest, apperrors2.ErrorCodeInvalidUUID)
		})
	})

	Describe("GetByClaimID", func() {
		DescribeTable("should handle different scenarios",
			func(setupMock func(), expectedStatus int, expectedError string) {
				setupMock()

				r.GET("/claims/:id/attachments", attachmentHandler.GetByClaimID)
				req, _ := http.NewRequest(http.MethodGet, "/claims/"+claimID.String()+"/attachments", nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(expectedStatus))
				if expectedError != "" {
					ExpectErrorCode(w, expectedStatus, expectedError)
				} else {
					ExpectResponseNotNil(w, expectedStatus)
				}
			},
			Entry("successful retrieval with data",
				func() {
					attachments := []*entity.ClaimAttachment{sampleAttachment}
					mockService.EXPECT().GetByClaimID(mock.Anything, claimID).Return(attachments, nil).Once()
				},
				http.StatusOK, ""),
			Entry("empty results",
				func() {
					mockService.EXPECT().GetByClaimID(mock.Anything, claimID).Return([]*entity.ClaimAttachment{}, nil).Once()
				},
				http.StatusOK, ""),
			Entry("service error",
				func() {
					mockService.EXPECT().GetByClaimID(mock.Anything, claimID).
						Return(nil, errors.New("database error")).Once()
				},
				http.StatusInternalServerError, apperrors2.ErrorCodeInternalServerError),
		)

		It("should return error for invalid claim ID", func() {
			r.GET("/claims/:id/attachments", attachmentHandler.GetByClaimID)
			req, _ := http.NewRequest(http.MethodGet, "/claims/invalid-uuid/attachments", nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusBadRequest, apperrors2.ErrorCodeInvalidUUID)
		})
	})

	Describe("Create", func() {
		Context("when authorized as SC_TECHNICIAN", func() {
			BeforeEach(func() {
				r.POST("/claims/:id/attachments", func(c *gin.Context) {
					SetHeaderRole(c, entity.UserRoleScTechnician)
					attachmentHandler.Create(c)
				})
			})

			It("should create single attachment successfully", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Run(func(ctx context.Context, fn func(application.Tx) error) {
						mockService.EXPECT().Create(mockTx, claimID, mock.Anything).
							Return(sampleAttachment, nil).Once()
						_ = fn(mockTx)
					}).Return(nil).Once()

				body, contentType := createMultipartForm()
				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/attachments", body)
				req.Header.Set("Content-Type", contentType)
				r.ServeHTTP(w, req)

				ExpectResponseNotNil(w, http.StatusCreated)
			})

			It("should create multiple attachments successfully", func() {
				attachment2 := &entity.ClaimAttachment{
					ID: uuid.New(), ClaimID: claimID, Type: entity.AttachmentTypeImage, URL: "https://example.com/image2.jpg"}

				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Run(func(ctx context.Context, fn func(application.Tx) error) {
						mockService.EXPECT().Create(mockTx, claimID, mock.Anything).Return(sampleAttachment, nil).Once()
						mockService.EXPECT().Create(mockTx, claimID, mock.Anything).Return(attachment2, nil).Once()
						_ = fn(mockTx)
					}).Return(nil).Once()

				body, contentType := createMultipartForm("image1.jpg", "image2.jpg")
				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/attachments", body)
				req.Header.Set("Content-Type", contentType)
				r.ServeHTTP(w, req)

				ExpectResponseNotNil(w, http.StatusCreated)
			})

			DescribeTable("should handle error scenarios",
				func(setupRequest func() *http.Request, expectedError string) {
					req := setupRequest()
					r.ServeHTTP(w, req)
					ExpectErrorCode(w, http.StatusBadRequest, expectedError)
				},
				Entry("invalid claim ID",
					func() *http.Request {
						body, contentType := createMultipartForm()
						req, _ := http.NewRequest(http.MethodPost, "/claims/invalid-uuid/attachments", body)
						req.Header.Set("Content-Type", contentType)
						return req
					},
					apperrors2.ErrorCodeInvalidUUID),
				Entry("non-multipart request",
					func() *http.Request {
						req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/attachments", strings.NewReader("not multipart"))
						req.Header.Set("Content-Type", "application/json")
						return req
					},
					apperrors2.ErrorCodeInvalidMultipartFormRequest),
				Entry("no files provided",
					func() *http.Request {
						body := &bytes.Buffer{}
						writer := multipart.NewWriter(body)
						_ = writer.WriteField("other_field", "value")
						_ = writer.Close()
						req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/attachments", body)
						req.Header.Set("Content-Type", writer.FormDataContentType())
						return req
					},
					apperrors2.ErrorCodeInvalidMultipartFormRequest),
			)

			It("should handle service error during creation", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Run(func(ctx context.Context, fn func(application.Tx) error) {
						mockService.EXPECT().Create(mockTx, claimID, mock.Anything).
							Return(nil, apperrors2.NewClaimNotFound()).Once()
						_ = fn(mockTx)
					}).Return(apperrors2.NewClaimNotFound()).Once()

				body, contentType := createMultipartForm()
				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/attachments", body)
				req.Header.Set("Content-Type", contentType)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusNotFound, apperrors2.ErrorCodeClaimNotFound)
			})

			It("should handle transaction failure", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Return(errors.New("transaction failed")).Once()

				body, contentType := createMultipartForm()
				req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/attachments", body)
				req.Header.Set("Content-Type", contentType)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusInternalServerError, apperrors2.ErrorCodeInternalServerError)
			})
		})

		It("should deny access for unauthorized roles", func() {
			r.POST("/claims/:id/attachments", func(c *gin.Context) {
				SetHeaderRole(c, entity.UserRoleEvmStaff)
				attachmentHandler.Create(c)
			})

			body, contentType := createMultipartForm()
			req, _ := http.NewRequest(http.MethodPost, "/claims/"+claimID.String()+"/attachments", body)
			req.Header.Set("Content-Type", contentType)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusForbidden, apperrors2.ErrorCodeUnauthorizedRole)
		})
	})

	Describe("Delete", func() {
		Context("when authorized as SC_TECHNICIAN", func() {
			BeforeEach(func() {
				r.DELETE("/claims/:id/attachments/:attachmentID", func(c *gin.Context) {
					SetHeaderRole(c, entity.UserRoleScTechnician)
					attachmentHandler.Delete(c)
				})
			})

			It("should delete attachment successfully", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Run(func(ctx context.Context, fn func(application.Tx) error) {
						mockService.EXPECT().HardDelete(mockTx, claimID, attachmentID).Return(nil).Once()
						_ = fn(mockTx)
					}).Return(nil).Once()

				req, _ := http.NewRequest(http.MethodDelete, "/claims/"+claimID.String()+"/attachments/"+attachmentID.String(), nil)
				r.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNoContent))
			})

			It("should return error for invalid claim ID", func() {
				req, _ := http.NewRequest(http.MethodDelete, "/claims/invalid-uuid/attachments/"+attachmentID.String(), nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusBadRequest, apperrors2.ErrorCodeInvalidUUID)
			})

			It("should return error for invalid attachment ID", func() {
				req, _ := http.NewRequest(http.MethodDelete, "/claims/"+claimID.String()+"/attachments/invalid-uuid", nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusBadRequest, apperrors2.ErrorCodeInvalidUUID)
			})

			It("should handle attachment not found error", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Run(func(ctx context.Context, fn func(application.Tx) error) {
						mockService.EXPECT().HardDelete(mockTx, claimID, attachmentID).
							Return(apperrors2.NewClaimAttachmentNotFound()).Once()
						_ = fn(mockTx)
					}).Return(apperrors2.NewClaimAttachmentNotFound()).Once()

				req, _ := http.NewRequest(http.MethodDelete, "/claims/"+claimID.String()+"/attachments/"+attachmentID.String(), nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusNotFound, apperrors2.ErrorCodeClaimAttachmentNotFound)
			})

			It("should handle service error", func() {
				mockTxManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(application.Tx) error")).
					Run(func(ctx context.Context, fn func(application.Tx) error) {
						mockService.EXPECT().HardDelete(mockTx, claimID, attachmentID).
							Return(errors.New("database error")).Once()
						_ = fn(mockTx)
					}).Return(errors.New("database error")).Once()

				req, _ := http.NewRequest(http.MethodDelete, "/claims/"+claimID.String()+"/attachments/"+attachmentID.String(), nil)
				r.ServeHTTP(w, req)

				ExpectErrorCode(w, http.StatusInternalServerError, apperrors2.ErrorCodeInternalServerError)
			})
		})

		It("should deny access for unauthorized roles", func() {
			r.DELETE("/claims/:id/attachments/:attachmentID", func(c *gin.Context) {
				SetHeaderRole(c, entity.UserRoleEvmStaff)
				attachmentHandler.Delete(c)
			})

			req, _ := http.NewRequest(http.MethodDelete, "/claims/"+claimID.String()+"/attachments/"+attachmentID.String(), nil)
			r.ServeHTTP(w, req)

			ExpectErrorCode(w, http.StatusForbidden, apperrors2.ErrorCodeUnauthorizedRole)
		})
	})
})
