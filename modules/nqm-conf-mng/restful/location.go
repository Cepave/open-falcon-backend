package restful

import (
	"net/http"

	commonOwlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	"gopkg.in/gin-gonic/gin.v1"
)

func listProvinces(c *gin.Context) {
	ProvinceName := c.Param("name")
	provinces := commonOwlDb.GetProvincesByName(ProvinceName)
	c.JSON(http.StatusOK, provinces)
}

func listCities(c *gin.Context) {
	CityName := c.Param("name")
	cities := commonOwlDb.GetCitiesByName(CityName)
	c.JSON(http.StatusOK, cities)
}
