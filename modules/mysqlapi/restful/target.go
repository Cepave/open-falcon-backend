package restful

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	commonNqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	"github.com/Cepave/open-falcon-backend/common/gin/mvc"
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

func listTargets(
	q *commonNqmModel.TargetQuery,
	p *struct {
		Paging *commonModel.Paging `mvc:"pageSize[50] pageOrderBy[status#desc:name#asc:host#asc]"`
	},
) (*commonModel.Paging, mvc.OutputBody) {
	pingtasks, resultPaging := commonNqmDb.ListTargets(q, *p.Paging)

	return resultPaging, mvc.JsonOutputBody(pingtasks)
}
