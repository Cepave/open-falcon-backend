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
	hbsResp := GetGeneralConfig().hbsResp.Load().(model.NqmPingTaskResponse)
	if !hbsResp.NeedPing {
		return nil, nil, model.NqmAgent{}, fmt.Errorf("[ " + u.UtilName() + " ] No tasks assigned.")
	}
	if !GetGeneralConfig().Measurements[u.UtilName()].enabled {
		return nil, nil, model.NqmAgent{}, fmt.Errorf("[ " + u.UtilName() + " ] Not enabled.")

	}
	targets := make([]model.NqmTarget, len(hbsResp.Targets))
	copy(targets, hbsResp.Targets)

	agent := *hbsResp.Agent

	command := make([]string, len(hbsResp.Command))
	copy(command, hbsResp.Command)

	targetAddressList := getTargetAddressList(targets)
	probingCmd := u.ProbingCommand(command, targetAddressList)
	return probingCmd, targets, agent, nil
}
