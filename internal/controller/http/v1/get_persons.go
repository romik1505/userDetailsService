package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/romik1505/userDetailsService/internal/domain"
	"github.com/romik1505/userDetailsService/internal/util"
	"net/http"
)

// GetPersons	godoc
//
// @Summary Get list persons with pagination
// @Tags persons
// @Accept json
// @Produce json
// @Param input query domain.ListPersonsFilter true "filter for list persons"
// @Success 200 {object} common.Person
// @Failure 400,500 {string} string
// @Router /person [get]
func (h *Handler) GetPersons(c *gin.Context) {
	filter, err := BindListPersonFilter(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	pagination, err := h.PersonService.ListPersons(c, filter)
	if err != nil {
		ErrorWrapper(c, err)
	}

	c.JSON(200, pagination)
}

func BindListPersonFilter(c *gin.Context) (domain.ListPersonsFilter, error) {
	filter := domain.ListPersonsFilter{}
	if err := c.ShouldBindQuery(&filter); err != nil {
		return domain.ListPersonsFilter{}, err
	}

	rawIds, _ := c.GetQuery("ids")
	ids := util.Values(rawIds, util.QueryListToInts64)

	rawGenders, _ := c.GetQuery("gender_in")
	genders := util.Values(rawGenders, util.QueryListToStrings)

	rawNationalities, _ := c.GetQuery("nationality_in")
	nationalities := util.Values(rawNationalities, util.QueryListToStrings)

	filter.IDs = ids
	filter.GenderIn = genders
	filter.NationalityIn = nationalities

	return filter, nil
}
