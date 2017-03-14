package restful

import (
	"net/http"
	"strconv"

	commonOwlDb "github.com/Cepave/open-falcon-backend/common/db/owl"
	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
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

func getProvinceByID(
	p *struct {
		ProvinceID int16 `mvc:"param[province_id]"`
	},
) mvc.OutputBody {
	return mvc.JsonOutputOrNotFound(
		commonOwlDb.GetProvinceById(p.ProvinceID),
	)
}

func listCities(c *gin.Context) {
	var cityName string
	if v, ok := c.GetQuery("name"); ok {
		cityName = v
	}
	cities := commonOwlDb.GetCitiesByName(cityName)
	c.JSON(http.StatusOK, cities)
}

func getCityByID(
	p *struct {
		CityID int16 `mvc:"param[city_id]"`
	},
) mvc.OutputBody {
	return mvc.JsonOutputOrNotFound(
		commonOwlDb.GetCityById(p.CityID),
	)
}

func listCitiesInProvince(c *gin.Context) {
	var cityName string
	if v, ok := c.GetQuery("name"); ok {
		cityName = v
	}

	provinceId, err := strconv.Atoi(c.Param("province_id"))
	if err != nil {
		commonGin.OutputJsonIfNotNil(c, nil)
	}

	cities := commonOwlDb.GetCitiesInProvinceByName(provinceId, cityName)
	c.JSON(http.StatusOK, cities)
}
