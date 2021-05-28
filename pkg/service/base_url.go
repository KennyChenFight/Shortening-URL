package service

import (
	"net/http"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"

	"github.com/gin-gonic/gin"
)

func (s *BaseService) CreateShorteningURL(c *gin.Context) {
	var request struct {
		URL string `json:"url" binding:"required,min=1,max=2048"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		s.responseWithError(c, business.NewError(business.Validation, http.StatusBadRequest, "invalid url field", err))
		return
	}

	url, err := s.urlRepository.CreateShorteningURL(request.URL)
	if err != nil {
		s.responseWithError(c, err)
		return
	}

	s.responseWithSuccess(c, business.NewSuccess(http.StatusCreated, gin.H{"id": url.ID, "shortUrl": combineFQDNWithShorteningURLID(s.config.FQDN, url.ID)}))
}

func (s *BaseService) GetOriginalURL(c *gin.Context) {
	var request struct {
		ID string `json:"id" uri:"id" binding:"len=6"`
	}
	if err := c.ShouldBindUri(&request); err != nil {
		s.responseWithError(c, business.NewError(business.Validation, http.StatusBadRequest, "invalid id field", err))
		return
	}

	id := c.Param("id")
	originalURL, err := s.urlRepository.GetOriginalURL(id)
	if err != nil {
		s.responseWithError(c, err)
		return
	}
	s.responseWithSuccess(c, business.NewSuccess(http.StatusTemporaryRedirect, originalURL))
}

func (s *BaseService) DeleteShorteningURL(c *gin.Context) {
	var request struct {
		ID string `json:"id" uri:"id" binding:"len=6"`
	}
	if err := c.ShouldBindUri(&request); err != nil {
		s.responseWithError(c, business.NewError(business.Validation, http.StatusBadRequest, "invalid id field", err))
		return
	}

	err := s.urlRepository.DeleteShorteningURL(request.ID)
	if err != nil {
		s.responseWithError(c, err)
		return
	}
	s.responseWithSuccess(c, business.NewSuccess(http.StatusNoContent, nil))
}
