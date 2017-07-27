package nqm

import (
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/model"
	nqmModel "github.com/Cepave/open-falcon-backend/common/model/nqm"
	"github.com/icrowley/fake"
	"sync"
	"time"

	dbTest "github.com/Cepave/open-falcon-backend/common/testing/db"

	//ocheck "github.com/Cepave/open-falcon-backend/common/testing/check"
	. "gopkg.in/check.v1"
)

type TestNqmStressSuite struct{}

var _ = Suite(&TestNqmStressSuite{})

func (suite *TestNqmStressSuite) TestGpa(c *C) {
	const numberOfRoutines = 16
	const numberOfNqmAgents = 0

	if numberOfNqmAgents == 0 {
		c.Skip("Number of NQM agents is 0. See numberOfNqmAgents in code")
		return
	}

	pool := make(chan bool, numberOfRoutines)
	for i := 0; i < numberOfRoutines; i++ {
		pool <- true
	}

	now := time.Now()
	for i := 0; i < numberOfNqmAgents; i++ {
		<-pool

		go func() {
			defer func() {
				pool <- true
			}()

			taskRequest := generateTaskRequest()

			c.Logf("%+v", taskRequest)
			RefreshAgentInfo(
				nqmModel.NewNqmAgent(taskRequest),
				now,
			)
		}()
	}

	for i := 0; i < numberOfRoutines; i++ {
		<-pool
	}
}

var lockForTaskRequest = &sync.Mutex{}

func generateTaskRequest() *model.NqmTaskRequest {
	lockForTaskRequest.Lock()
	defer lockForTaskRequest.Unlock()

	hostname := fake.DomainName()
	ipAddress := fake.IPv4()

	return &model.NqmTaskRequest{
		ConnectionId: fmt.Sprintf("%s@%s", ipAddress, hostname),
		Hostname:     hostname,
		IpAddress:    ipAddress,
	}
}

func (s *TestNqmStressSuite) SetUpSuite(c *C) {
	DbFacade = dbTest.InitDbFacade(c)
}
func (s *TestNqmStressSuite) TearDownSuite(c *C) {
	dbTest.ReleaseDbFacade(c, DbFacade)
}
