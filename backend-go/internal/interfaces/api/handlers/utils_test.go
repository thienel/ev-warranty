package handlers_test

import (
	"bytes"
	"encoding/json"
	"ev-warranty-go/internal/interfaces/api/dtos"
	"ev-warranty-go/pkg/mocks"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

func SetupMock(t FullGinkgoTInterface) (*mocks.Logger, *gin.Engine, *httptest.ResponseRecorder) {
	mockLogger := mocks.NewLogger(t)
	mockLogger.EXPECT().Info(mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	mockLogger.EXPECT().Info(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything).Return().Maybe()
	mockLogger.EXPECT().Error(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	responseRecorder := httptest.NewRecorder()
	return mockLogger, router, responseRecorder
}

func SetHeaderRole(c *gin.Context, role string) {
	c.Request.Header.Set("X-User-Role", role)
}

func SetHeaderID(c *gin.Context, id uuid.UUID) {
	c.Request.Header.Set("X-User-ID", id.String())
}

func SetContentTypeJSON(c *gin.Context) {
	c.Request.Header.Set("Content-Type", "application/json")
}

func SendRequest(router *gin.Engine, method, url string, w *httptest.ResponseRecorder, req any) {
	var body []byte
	if req == nil {
		body = nil
	} else {
		body, _ = json.Marshal(req)
	}
	httpReq, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	router.ServeHTTP(w, httpReq)
}

func ExpectResponseNotNil(w *httptest.ResponseRecorder, httpStatus int) {
	GinkgoHelper()
	Expect(w.Code).To(Equal(httpStatus))
	var response dtos.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	Expect(err).NotTo(HaveOccurred())
	Expect(response.Data).NotTo(BeNil())
}

func ExpectErrorCode(w *httptest.ResponseRecorder, httpStatus int, error string) {
	GinkgoHelper()
	Expect(w.Code).To(Equal(httpStatus))
	var response dtos.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	Expect(err).NotTo(HaveOccurred())
	Expect(response.Error).To(Equal(error))
}
