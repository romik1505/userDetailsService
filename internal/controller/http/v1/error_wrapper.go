package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romik1505/userDetailsService/internal/service"
	"net/http"
)

func ErrorWrapper(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrBadRequest):
		c.AbortWithError(http.StatusBadRequest, err)
	case errors.Is(err, service.ErrNotFound), errors.Is(err, service.ErrEmptyResponse):
		c.AbortWithError(http.StatusNoContent, err)
	case errors.Is(err, service.ErrInternalError):
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}
