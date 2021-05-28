package dao

import (
	"errors"
	"net/http"
	"time"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PGUrlDAO", func() {
	var pgUrlDAO *PGUrlDAO
	var logger *loglib.Logger

	BeforeEach(func() {
		logger = loglib.NewNopLogger()
		pgUrlDAO = &PGUrlDAO{logger, testPGClient}
	})

	var _ = Describe("Create", func() {
		var (
			expectURL *URL
			createErr *business.Error
		)

		actualOriginalURL := "http://example.com"
		actualKey := Key{
			ID:        "random",
			CreatedAt: time.Now(),
		}
		actualURL := &URL{
			ID:        actualKey.ID,
			Original:  actualOriginalURL,
			CreatedAt: time.Now(),
			ExpiredAt: time.Now(),
		}

		JustBeforeEach(func() {
			expectURL, createErr = pgUrlDAO.Create(actualOriginalURL)
		})

		Context("success", func() {
			BeforeEach(func() {
				_, err := testPGClient.Model(&actualKey).Insert()
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				_, err := testPGClient.Model((*URL)(nil)).Where("id = ?", actualKey.ID).Delete()
				Expect(err).To(BeNil())
			})

			It("result", func() {
				Expect(createErr).To(BeNil())
				Expect(expectURL.ID).To(Equal(actualURL.ID))
				Expect(expectURL.Original).To(Equal(actualURL.Original))
				Ω(testPGClient.Model(&Key{}).Where("id = ?", actualKey.ID).Count()).To(Equal(0))
				Ω(testPGClient.Model(&URL{}).Where("id = ?", actualKey.ID).Count()).To(Equal(1))
			})
		})

		Context("key not found", func() {
			It("result", func() {
				Expect(createErr).To(Equal(business.NewError(business.NotFound, http.StatusNotFound, "record not found", errors.New(PGErrMsgNoRowsFound))))
			})
		})
	})

	var _ = Describe("Get", func() {
		var (
			expectURL *URL
			getErr    *business.Error
		)

		now := time.Now().UTC().Truncate(time.Millisecond)
		actualURL := &URL{
			ID:        "random",
			Original:  "http://example.com",
			CreatedAt: now,
			ExpiredAt: now.Add(time.Minute),
		}

		JustBeforeEach(func() {
			expectURL, getErr = pgUrlDAO.Get(actualURL.ID)
		})

		Context("success", func() {
			BeforeEach(func() {
				_, err := testPGClient.Model(actualURL).Insert()
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				_, err := testPGClient.Model(actualURL).WherePK().Delete()
				Expect(err).To(BeNil())
			})

			It("result", func() {
				Expect(getErr).To(BeNil())
				Expect(expectURL).To(Equal(actualURL))
			})
		})

		Context("url not found when expired_at <= now", func() {
			BeforeEach(func() {
				actualURL.ExpiredAt = actualURL.CreatedAt
				_, err := testPGClient.Model(actualURL).Insert()
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				_, err := testPGClient.Model(actualURL).WherePK().Delete()
				Expect(err).To(BeNil())
			})

			It("result", func() {
				Expect(getErr).To(Equal(business.NewError(business.NotFound, http.StatusNotFound, "record not found", getErr.Reason)))
				Expect(expectURL).To(BeNil())
			})
		})

		Context("url not found when id not correct", func() {
			It("result", func() {
				Expect(getErr).To(Equal(business.NewError(business.NotFound, http.StatusNotFound, "record not found", getErr.Reason)))
				Expect(expectURL).To(BeNil())
			})
		})
	})

	var _ = Describe("Delete", func() {
		var (
			deleteErr *business.Error
		)

		now := time.Now().UTC()
		actualURL := &URL{
			ID:        "random",
			Original:  "http://example.com",
			CreatedAt: now,
			ExpiredAt: now.Add(1 * time.Second),
		}

		JustBeforeEach(func() {
			deleteErr = pgUrlDAO.Delete(actualURL.ID)
		})

		Context("success", func() {
			BeforeEach(func() {
				_, err := testPGClient.Model(actualURL).Insert()
				Expect(err).To(BeNil())
			})

			It("result", func() {
				Expect(deleteErr).To(BeNil())
				Expect(testPGClient.Model((*URL)(nil)).Where("id = ?", actualURL.ID).Select()).To(HaveOccurred())
			})
		})
	})

	var _ = Describe("Expire", func() {
		var (
			expectExpireNum int
			expireErr       *business.Error
		)

		now := time.Now().UTC()
		actualURLs := []URL{
			{
				ID:        "000000",
				Original:  "http://example.com",
				CreatedAt: now,
				ExpiredAt: now.Add(-100 * time.Second),
			},
			{
				ID:        "111111",
				Original:  "http://example.com",
				CreatedAt: now,
				ExpiredAt: now.Add(-100 * time.Second),
			},
			{
				ID:        "222222",
				Original:  "http://example.com",
				CreatedAt: now,
				ExpiredAt: now.Add(-100 * time.Second),
			},
		}
		actualLimitNum := len(actualURLs) - 1

		JustBeforeEach(func() {
			expectExpireNum, expireErr = pgUrlDAO.Expire(actualLimitNum)
		})

		Context("success", func() {
			BeforeEach(func() {
				_, err := testPGClient.Model(&actualURLs).Insert()
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				deleteIDs := make([]string, 0)
				for _, url := range actualURLs[actualLimitNum:] {
					deleteIDs = append(deleteIDs, url.ID)
				}
				_, err := testPGClient.Model((*URL)(nil)).WhereIn("id in (?)", deleteIDs).Delete()
				Expect(err).To(BeNil())
			})

			It("result", func() {
				Expect(expireErr).To(BeNil())
				Expect(expectExpireNum).To(Equal(actualLimitNum))
			})
		})
	})

})
