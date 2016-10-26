package nqm

import (
	commonDb "github.com/Cepave/open-falcon-backend/common/db"
)

// The query conditions of agent
type AgentQuery struct {
	Name string
	ConnectionId string
	Hostname string
	IpAddress string

	HasIspId bool
	IspId int16

	HasStatusCondition bool
	Status bool
}

// Gets the []byte used to perform like in MySql
func (query *AgentQuery) GetIpForLikeCondition() []byte {
	bytes, err := commonDb.IpV4ToBytesForLike(query.IpAddress)
	if err != nil {
		panic(err)
	}

	return bytes
}
