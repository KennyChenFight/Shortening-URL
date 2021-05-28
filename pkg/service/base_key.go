package service

import (
	"net/http"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/gin-gonic/gin"
)

func (s *BaseService) BatchCreateKeys(c *gin.Context) {
	var request struct {
		Number int `json:"number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		s.responseWithError(c, business.NewError(business.Validation, http.StatusBadRequest, "invalid number field", err))
		return
	}

	length, err := s.urlRepository.BatchCreateKeys(request.Number)
	if err != nil {
		s.responseWithError(c, err)
		return
	}

	s.responseWithSuccess(c, business.NewSuccess(http.StatusCreated, gin.H{"length": length}))
}
