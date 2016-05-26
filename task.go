package main

import (
	"fmt"

	"github.com/Cepave/common/model"
)

func Task(u Utility) ([]string, []model.NqmTarget, *model.NqmAgent, error) {
	/**
	 * Only 2 possible responses come from hbs:
	 *     1. NeedPing==false (default condition)
	 *         NqmAgent, NQMTargets, Command are nil
	 *     2. NeedPing==ture
	 *         NqmAgent, NQMTargets, Command are not nil
	 */
	if !GetGeneralConfig().hbsResp.NeedPing {
		return nil, nil, nil, fmt.Errorf("[ " + u.utilName() + " ] No tasks assigned.")
	}

	targetAddressList := getTargetAddressList()
	probingCmd := u.ProbingCommand(targetAddressList)
	return probingCmd, GetGeneralConfig().hbsResp.Targets, GetGeneralConfig().hbsResp.Agent, nil
}
