package rdb

import (
	"github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("[Intg] Test GetAgentConfig", ginkgoDb.NeedDb(func() {
	BeforeEach(func() {
		inTx(
			"INSERT INTO common_config(`key`, `value`)" +
				"VALUES('TestGetAgentConfig', 'https://example.com/Cepave/TestGetAgentConfig.git')",
		)
	})

	AfterEach(func() {
		inTx(
			"DELETE FROM common_config WHERE `key` = 'TestGetAgentConfig'",
		)
	})

	DescribeTable("when key is",
		func(key string, expectedResult *model.AgentConfigResult) {
			res := GetAgentConfig(key)
			Expect(res).To(Equal(expectedResult))
		},
		Entry("Nonexistent", "Non-existent-key", nil),
		Entry("Existent", "TestGetAgentConfig",
			&model.AgentConfigResult{
				"TestGetAgentConfig",
				"https://example.com/Cepave/TestGetAgentConfig.git",
			},
		),
	)
}))
