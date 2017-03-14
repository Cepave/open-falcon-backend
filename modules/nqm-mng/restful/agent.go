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
func listAgents(c *gin.Context) {
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

	query := buildQueryForListAgents(c)
	agents, resultPaging := commonNqmDb.ListAgents(query, *paging)

	commonGin.HeaderWithPaging(c, resultPaging)
	c.JSON(http.StatusOK, agents)
}

func buildQueryForListAgents(c *gin.Context) *commonNqmModel.AgentQuery {
	query := &commonNqmModel.AgentQuery{}
	queryWrapper := commonGin.NewQueryWrapper(c)

	/**
	 * Set-up query parameters
	 */
	if v, ok := c.GetQuery("name"); ok {
		query.Name = v
	}
	if v, ok := c.GetQuery("connection_id"); ok {
		query.ConnectionId = v
	}
	if v, ok := c.GetQuery("hostname"); ok {
		query.Hostname = v
	}
	if v, ok := c.GetQuery("ip_address"); ok {
		query.IpAddress = v
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
