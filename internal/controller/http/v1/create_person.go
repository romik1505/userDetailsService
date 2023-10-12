package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/romik1505/userDetailsService/internal/common"
	"net/http"
)

// CreatePerson	godoc
//
// @Summary Добавить нового человека
// @Tags persons
// @Accept json
// @Produce json
// @Param input body common.CreatePersonRequest true "person info"
// @Success 200,204
// @Failure 400,500 {string} string
// @Router /person [post]
func (h *Handler) CreatePerson(c *gin.Context) {
	var person common.CreatePersonRequest
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.PersonService.CreatePersonPushMessage(c, person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}
