package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	commonModelConfig "github.com/Cepave/open-falcon-backend/common/model/config"
	"github.com/Cepave/open-falcon-backend/modules/nqm-mng/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test agentHeartbeatCall() of AgentHeartbeat service", func() {
	var (
		dataNum int = 3
		agents  []*model.FalconAgentHeartbeat
	)

	BeforeEach(func() {
		agents = make([]*model.FalconAgentHeartbeat, 0)
		for i := 0; i < dataNum; i++ {
			agents = append(agents, requestToHeartbeat(generateRandomRequest(), time.Now().Unix()))
		}
	})

	Context("when the call succeed", func() {
		var ts *httptest.Server

		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				decorder := json.NewDecoder(r.Body)
				var rAgents []*model.FalconAgentHeartbeat
				err := decorder.Decode(&rAgents)
				if err != nil {
					Fail(err.Error())
				}
				defer r.Body.Close()
				rowsAffectedCnt := int64(len(rAgents))
				res := model.FalconAgentHeartbeatResult{rowsAffectedCnt}
				if resp, err := json.Marshal(res); err != nil {
					Fail(err.Error())
				} else {
					w.Write(resp)
				}
			}))
		})

		AfterEach(func() {
			ts.Close()
		})

		It("should return correct affected amount", func() {
			InitPackage(&commonModelConfig.MysqlApiConfig{Host: ts.URL}, "")
			rowsAffectedCnt, agentsDroppedCnt := agentHeartbeatCall(agents)

			Expect(rowsAffectedCnt).To(Equal(int64(dataNum)))
			Expect(agentsDroppedCnt).To(BeZero())
		})
	})

	Context("when the call fail", func() {
		It("should return correct dropped amount", func() {
			InitPackage(&commonModelConfig.MysqlApiConfig{Host: "dummyHost"}, "")
			rowsAffectedCnt, agentsDroppedCnt := agentHeartbeatCall(agents)

			Expect(rowsAffectedCnt).To(BeZero())
			Expect(agentsDroppedCnt).To(Equal(int64(dataNum)))
		})
	})
})