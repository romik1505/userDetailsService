package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// DeletePerson	godoc
//
// @Summary delete person by id
// @Tags persons
// @Accept json
// @Produce json
// @Param id query int64 true "person for delete"
// @Success 200,204
// @Failure 400,500 {string} string
// @Router /person [delete]
func (h *Handler) DeletePerson(c *gin.Context) {
	id := c.Query("id")
	numID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	if err := h.PersonService.DeletePerson(c, numID); err != nil {
		ErrorWrapper(c, err)
	}

	c.Status(200)
}
