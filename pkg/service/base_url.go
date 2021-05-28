package service

import (
	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *BaseService) CreateShorteningURL(c *gin.Context) {
	var request struct {
		URL string `json:"url" binding:"required,min=1,max=2048"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		sendErrorResponse(c, s, business.NewError(business.Validation, http.StatusBadRequest, "invalid url field", err))
		return
	}

	url, err := s.urlDAO.Create(request.URL)
	if err != nil {
		sendErrorResponse(c, s, err)
		return
	}
	sendSuccessResponse(c, http.StatusCreated, gin.H{"id": url.ID, "shortUrl": combineFQDNAndShorteningURLID(s.config.FQDN, url.ID)})
}

func (s *BaseService) GetShorteningURL(c *gin.Context) {
	var request struct {
		ID string `json:"id" uri:"id" binding:"len=6"`
	}
	if err := c.ShouldBindUri(&request); err != nil {
		sendErrorResponse(c, s, business.NewError(business.Validation, http.StatusBadRequest, "invalid id field", err))
		return
	}

	id := c.Param("id")
	url, err := s.urlDAO.Get(id)
	if err != nil {
		sendErrorResponse(c, s, err)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url.Original)
}

func (s *BaseService) DeleteShorteningURL(c *gin.Context) {
	var request struct {
		ID string `json:"id" uri:"id" binding:"len=6"`
	}
	if err := c.ShouldBindUri(&request); err != nil {
		sendErrorResponse(c, s, business.NewError(business.Validation, http.StatusBadRequest, "invalid id field", err))
		return
	}

	err := s.urlDAO.Delete(request.ID)
	if err != nil {
		sendErrorResponse(c, s, err)
		return
	}
	sendSuccessResponse(c, http.StatusNoContent, nil)
}
