package restful

import (
	"net/http"

	commonOwlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	"gopkg.in/gin-gonic/gin.v1"
)

func listProvinces(c *gin.Context) {
	var provinceName string
	if v, ok := c.GetQuery("name"); ok {
		provinceName = v
	}
	provinces := commonOwlDb.GetProvincesByName(provinceName)
	c.JSON(http.StatusOK, provinces)
}

func listCities(c *gin.Context) {
	var cityName string
	if v, ok := c.GetQuery("name"); ok {
		cityName = v
	}
	cities := commonOwlDb.GetCitiesByName(cityName)
	c.JSON(http.StatusOK, cities)
}
