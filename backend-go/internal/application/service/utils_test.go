package service_test

import (
	"errors"
	"ev-warranty-go/pkg/apperror"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func ExpectAppError(err error, expectedCode string) {
	GinkgoHelper()
	Expect(err).To(HaveOccurred())
	var appErr *apperror.AppError
	Expect(errors.As(err, &appErr)).To(BeTrue(), "error should be an AppError")
	Expect(appErr.ErrorCode).To(Equal(expectedCode), "error code should match")
}
