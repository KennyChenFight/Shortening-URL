package dao

import (
	"github.com/KennyChenFight/Shortening-URL/internal/randomstrgeneratormock"
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/KennyChenFight/randstr"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PGKeyDAO", func() {
	var pgKeyDAO *PGKeyDAO
	var logger *loglib.Logger
	var randomStrGeneratorMock *randomstrgeneratormock.MockRandomStrGenerator
	var mockCtrl *gomock.Controller

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		randomStrGeneratorMock = randomstrgeneratormock.NewMockRandomStrGenerator(mockCtrl)
		logger = loglib.NewNopLogger()
		pgKeyDAO = &PGKeyDAO{logger, testPGClient, randomStrGeneratorMock}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	var _ = Describe("BatchCreate", func() {
		var (
			expectCreateNumber int
			createErr          *business.Error
		)

		num := 2

		JustBeforeEach(func() {
			expectCreateNumber, createErr = pgKeyDAO.BatchCreate(num)
		})

		Context("success", func() {
			actualRandStrGenerator := randstr.NewFastGenerator(randstr.CharSetEnglishAlphabet)
			ids := []string{actualRandStrGenerator.GenerateRandomStr(randStrLength), actualRandStrGenerator.GenerateRandomStr(randStrLength)}
			BeforeEach(func() {
				for _, id := range ids {
					randomStrGeneratorMock.EXPECT().GenerateRandomStr(randStrLength).Return(id)
				}
			})

			AfterEach(func() {
				res, err := testPGClient.Model((*Key)(nil)).WhereIn("id in (?)", ids).Delete()
				Expect(err).To(BeNil())
				Expect(res.RowsAffected()).To(Equal(len(ids)))
			})

			It("result", func() {
				Expect(createErr).To(BeNil())
				Expect(expectCreateNumber).To(Equal(len(ids)))
			})
		})

		Context("success but insert on conflict", func() {
			ids := []string{"random", "random"}
			BeforeEach(func() {
				for _, id := range ids {
					randomStrGeneratorMock.EXPECT().GenerateRandomStr(randStrLength).Return(id)
				}
			})

			AfterEach(func() {
				res, err := testPGClient.Model((*Key)(nil)).WhereIn("id in (?)", ids).Delete()
				Expect(err).To(BeNil())
				Expect(res.RowsAffected()).To(Equal(len(ids) - 1))
			})

			It("result", func() {
				Expect(createErr).To(BeNil())
				Expect(expectCreateNumber).To(Equal(len(ids) - 1))
			})
		})
	})
})
