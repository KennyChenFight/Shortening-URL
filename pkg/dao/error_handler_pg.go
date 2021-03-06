package dao

import (
	"net/http"

	"github.com/KennyChenFight/Shortening-URL/pkg/business"
	"github.com/KennyChenFight/golib/loglib"
	"go.uber.org/zap"
)

const (
	PGErrMsgNoRowsFound      = "pg: no rows in result set"
	PGErrMsgNoMultiRowsFound = "pg: multiple rows in result set"
)

func pgErrorHandle(logger *loglib.Logger, err error) *business.Error {
	switch err.Error() {
	case PGErrMsgNoRowsFound:
		return business.NewError(business.NotFound, http.StatusNotFound, "record not found", err)
	case PGErrMsgNoMultiRowsFound:
		return business.NewError(business.NotFound, http.StatusNotFound, "multi record not found", err)
	default:
		logger.Error("postgres internal error", zap.Error(err))
		return business.NewError(business.PostgresInternalError, http.StatusInternalServerError, "internal error", err)
	}
}
