package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	cModel "github.com/Cepave/open-falcon-backend/common/model"
	cModelConfig "github.com/Cepave/open-falcon-backend/common/model/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test agentHeartbeatCall() of AgentHeartbeat service", func() {
	var (
		dataNum int = 3
		agents  []*cModel.FalconAgentHeartbeat
	)

	BeforeEach(func() {
		agents = make([]*cModel.FalconAgentHeartbeat, 0)
		for i := 0; i < dataNum; i++ {
			agents = append(agents, requestToHeartbeat(generateRandomRequest(), time.Now().Unix()))
		}
	})

	Context("when the call succeed", func() {
		var ts *httptest.Server

		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()
				decorder := json.NewDecoder(r.Body)
				var rAgents []*cModel.FalconAgentHeartbeat
				err := decorder.Decode(&rAgents)
				Expect(err).To(BeNil())

				rowsAffectedCnt := int64(len(rAgents))
				res := cModel.FalconAgentHeartbeatResult{rowsAffectedCnt}
				resp, err := json.Marshal(res)
				Expect(err).To(BeNil())
				w.Write(resp)
			}))
		})

		AfterEach(func() {
			ts.Close()
		})

		It("should return correct affected amount", func() {
			InitPackage(&cModelConfig.MysqlApiConfig{Host: ts.URL}, "")
			rowsAffectedCnt, agentsDroppedCnt := agentHeartbeatCall(agents)

			Expect(rowsAffectedCnt).To(Equal(int64(dataNum)))
			Expect(agentsDroppedCnt).To(BeZero())
		})
	})

	Context("when the call fail", func() {
		It("should return correct dropped amount", func() {
			InitPackage(&cModelConfig.MysqlApiConfig{Host: "dummyHost"}, "")
			rowsAffectedCnt, agentsDroppedCnt := agentHeartbeatCall(agents)

			Expect(rowsAffectedCnt).To(BeZero())
			Expect(agentsDroppedCnt).To(Equal(int64(dataNum)))
		})
	})
})
