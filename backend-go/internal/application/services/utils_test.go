package services_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"ev-warranty-go/internal/apperrors"
)

func ExpectAppError(err error, expectedCode string) {
	GinkgoHelper()
	Expect(err).To(HaveOccurred())
	var appErr *apperrors.AppError
	Expect(errors.As(err, &appErr)).To(BeTrue(), "error should be an AppError")
	Expect(appErr.ErrorCode).To(Equal(expectedCode), "error code should match")
}
