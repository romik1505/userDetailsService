package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/romik1505/userDetailsService/internal/common"
	"net/http"
)

// UpdatePerson	godoc
//
// @Summary update person by id
// @Tags persons
// @Accept json
// @Produce json
// @Param input body common.Person true "person for update"
// @Success 200,204
// @Failure 400,500 {string} string
// @Router /person [patch]
func (h *Handler) UpdatePerson(c *gin.Context) {
	var person common.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	err := h.PersonService.UpdatePerson(c, &person)
	if err != nil {
		ErrorWrapper(c, err)
	}

	c.JSON(200, person)
}
