package dao

import (
	"errors"
	"net/http"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pgErrorHandle", func() {
	var originalError error
	var businessError *business.Error
	JustBeforeEach(func() {
		businessError = pgErrorHandle(loglib.NewNopLogger(), originalError)
	})

	Context("when err == PGErrMsgNoRowsFound", func() {
		BeforeEach(func() {
			originalError = errors.New(PGErrMsgNoRowsFound)
		})
		AfterEach(func() {
			originalError = nil
		})
		It("result", func() {
			Expect(businessError).To(Equal(business.NewError(business.NotFound, http.StatusNotFound, "record not found", originalError)))
		})
	})

	Context("when err == PGErrMsgNoMultiRowsFound", func() {
		BeforeEach(func() {
			originalError = errors.New(PGErrMsgNoMultiRowsFound)
		})
		AfterEach(func() {
			originalError = nil
		})
		It("result", func() {
			Expect(businessError).To(Equal(business.NewError(business.NotFound, http.StatusNotFound, "multi record not found", originalError)))
		})
	})

	Context("internal error", func() {
		internalError := errors.New("internal error")
		BeforeEach(func() {
			originalError = internalError
		})
		AfterEach(func() {
			originalError = nil
		})
		It("result", func() {
			Expect(businessError).To(Equal(business.NewError(business.PostgresInternalError, http.StatusInternalServerError, "internal error", internalError)))
		})
	})
})
