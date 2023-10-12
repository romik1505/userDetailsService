package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/romik1505/userDetailsService/internal/service"
)

type Handler struct {
	PersonService service.Persons
}

func NewHandler(ps service.Persons) *Handler {
	return &Handler{
		PersonService: ps,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		user := v1.Group("/person")
		{
			user.GET("", h.GetPersons)
			user.POST("", h.CreatePerson)
			user.DELETE(":id", h.DeletePerson)
			user.PATCH("", h.UpdatePerson)
		}
	}
}
