package main

import (
	"fmt"

	"github.com/Cepave/common/model"
)

func getTargetAddressList(targets []model.NqmTarget) []string {
	var targetAddressList []string
	for _, target := range targets {
		targetAddressList = append(targetAddressList, target.Host)
	}
	return targetAddressList
}

func Task(u Utility) ([]string, []model.NqmTarget, model.NqmAgent, error) {
	/**
	 * Only 2 possible responses come from hbs:
	 *     1. NeedPing==false (default condition)
	 *         NqmAgent, NQMTargets, Command are nil
	 *     2. NeedPing==ture
	 *         NqmAgent, NQMTargets, Command are not nil
	 */
	if !GetGeneralConfig().hbsResp.NeedPing {
		return nil, nil, model.NqmAgent{}, fmt.Errorf("[ " + u.UtilName() + " ] No tasks assigned.")
	}
	if !GetGeneralConfig().Measurements[u.UtilName()].enabled {
		return nil, nil, model.NqmAgent{}, fmt.Errorf("[ " + u.UtilName() + " ] Not enabled.")

	}
	targets := make([]model.NqmTarget, len(GetGeneralConfig().hbsResp.Targets))
	copy(targets, GetGeneralConfig().hbsResp.Targets)

	agent := *GetGeneralConfig().hbsResp.Agent

	command := make([]string, len(GetGeneralConfig().hbsResp.Command))
	copy(command, GetGeneralConfig().hbsResp.Command)

	targetAddressList := getTargetAddressList(targets)
	probingCmd := u.ProbingCommand(targetAddressList)
	return probingCmd, targets, agent, nil
}
