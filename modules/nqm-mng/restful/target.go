package restful

import (
	"net/http"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"

	commonNqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	commonNqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
)

func addNewTarget(c *gin.Context) {
	/**
	 * Builds data from body of request
	 */
	targetForAdding := commonNqmModel.NewTargetForAdding()

	commonGin.BindJson(c, targetForAdding)
	commonGin.ConformAndValidateStruct(targetForAdding, commonNqmModel.Validator)
	targetForAdding.UniqueGroupTags()
	// :~)

	newTarget, err := commonNqmDb.AddTarget(targetForAdding)
	if err != nil {
		switch err.(type) {
		case commonNqmDb.ErrDuplicatedNqmTarget:
			commonGin.JsonConflictHandler(
				c,
				commonGin.DataConflictError{
					ErrorCode:    1,
					ErrorMessage: err.Error(),
				},
			)
		default:
			panic(err)
		}

		return
	}

	c.JSON(http.StatusOK, newTarget)
}
func modifyTarget(c *gin.Context) {
	/**
	 * Loads target from database
	 */
	targetId, targetIdErr := strconv.Atoi(c.Param("target_id"))
	if targetIdErr != nil {
		panic(targetIdErr)
	}

	originalTarget := commonNqmDb.GetTargetById(int32(targetId))
	if originalTarget == nil {
		commonGin.JsonNoMethodHandler(c)
		return
	}
	// :~)

	/**
	 * Binding JSON body to modified target
	 */
	modifiedTarget := originalTarget.ToTargetForAdding()
	commonGin.BindJson(c, modifiedTarget)
	commonGin.ConformAndValidateStruct(modifiedTarget, commonNqmModel.Validator)
	modifiedTarget.UniqueGroupTags()
	// :~)

	updatedTarget, err := commonNqmDb.UpdateTarget(originalTarget, modifiedTarget)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, updatedTarget)
}
func getTargetById(c *gin.Context) {
	targetId, err := strconv.Atoi(c.Param("target_id"))
	if err != nil {
		commonGin.OutputJsonIfNotNil(c, nil)
	}

	target := commonNqmDb.GetTargetById(int32(targetId))

	commonGin.OutputJsonIfNotNil(c, target)
}
func listTargets(c *gin.Context) {
	/**
	 * Set-up paging
	 */
	paging := commonGin.PagingByHeader(
		c,
		&commonModel.Paging{
			Size:     50,
			Position: 1,
		},
	)
	// :~)

	query := buildQueryForListTargets(c)
	targets, resultPaging := commonNqmDb.ListTargets(query, *paging)

	commonGin.HeaderWithPaging(c, resultPaging)
	c.JSON(http.StatusOK, targets)
}

func buildQueryForListTargets(c *gin.Context) *commonNqmModel.TargetQuery {
	query := &commonNqmModel.TargetQuery{}
	queryWrapper := commonGin.NewQueryWrapper(c)

	/**
	 * Set-up query parameters
	 */
	if v, ok := c.GetQuery("name"); ok {
		query.Name = v
	}
	if v, ok := c.GetQuery("host"); ok {
		query.Host = v
	}
	if param := queryWrapper.GetInt64("isp_id"); param.Viable {
		query.HasIspId = true
		query.IspId = int16(param.Value.(int64))
	}
	if param := queryWrapper.GetBool("status"); param.Viable {
		query.HasStatusCondition = true
		query.Status = param.Value.(bool)
	}
	// :~)

	return query
}
