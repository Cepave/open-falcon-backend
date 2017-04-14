package restful

import (
	"net/http"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"

	commonNqmDb "github.com/Cepave/open-falcon-backend/common/db/nqm"
	commonGin "github.com/Cepave/open-falcon-backend/common/gin"
	mvc "github.com/Cepave/open-falcon-backend/common/gin/mvc"
	commonModel "github.com/Cepave/open-falcon-backend/common/model"
	commonNqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
)

func addNewAgent(c *gin.Context) {
	/**
	 * Builds data from body of request
	 */
	agentForAdding := commonNqmModel.NewAgentForAdding()

	commonGin.BindJson(c, agentForAdding)
	commonGin.ConformAndValidateStruct(agentForAdding, commonNqmModel.Validator)
	agentForAdding.UniqueGroupTags()
	// :~)

	newAgent, err := commonNqmDb.AddAgent(agentForAdding)
	if err != nil {
		switch err.(type) {
		case commonNqmDb.ErrDuplicatedNqmAgent:
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

	c.JSON(http.StatusOK, newAgent)
}

func modifyAgent(c *gin.Context) {
	/**
	 * Loads agent from database
	 */
	agentId, agentIdErr := strconv.Atoi(c.Param("agent_id"))
	if agentIdErr != nil {
		panic(agentIdErr)
	}

	originalAgent := commonNqmDb.GetAgentById(int32(agentId))
	if originalAgent == nil {
		commonGin.JsonNoMethodHandler(c)
		return
	}
	// :~)

	/**
	 * Binding JSON body to modified agent
	 */
	modifiedAgent := originalAgent.ToAgentForAdding()
	commonGin.BindJson(c, modifiedAgent)
	commonGin.ConformAndValidateStruct(modifiedAgent, commonNqmModel.Validator)
	modifiedAgent.UniqueGroupTags()
	// :~)

	updatedAgent, err := commonNqmDb.UpdateAgent(originalAgent, modifiedAgent)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, updatedAgent)
}

func getAgentById(c *gin.Context) {
	agentId, err := strconv.Atoi(c.Param("agent_id"))
	if err != nil {
		commonGin.OutputJsonIfNotNil(c, nil)
	}

	agent := commonNqmDb.GetAgentById(int32(agentId))

	commonGin.OutputJsonIfNotNil(c, agent)
}

func listAgents(
	query *commonNqmModel.AgentQuery,
	paging *struct {
		Page *commonModel.Paging `mvc:"pageSize[50] pageOrderBy[status#desc:connection_id#asc]"`
	},
) (*commonModel.Paging, mvc.OutputBody) {
	agents, resultPaging := commonNqmDb.ListAgents(query, *paging.Page)

	return resultPaging, mvc.JsonOutputBody(agents)
}

func listTargetsOfAgentById(
	c *gin.Context,
	q *commonNqmModel.TargetsOfAgentQuery,
	p *struct {
		Paging *commonModel.Paging `mvc:"pageSize[50] pageOrderBy[status#desc:name#asc:host#asc]"`
	},
) (*commonModel.Paging, mvc.OutputBody) {
	p.Paging = commonGin.PagingByHeader(c, p.Paging)
	targets, resultPaging := commonNqmDb.ListTargetsOfAgentById(q, *p.Paging)

	return resultPaging, mvc.JsonOutputOrNotFound(targets)
}
