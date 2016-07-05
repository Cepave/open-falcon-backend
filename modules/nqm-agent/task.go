package main

import (
	"fmt"
	"time"

	"github.com/Cepave/open-falcon-backend/common/model"
)

func getTargetAddressList(targets []model.NqmTarget) []string {
	var targetAddressList []string
	for _, target := range targets {
		targetAddressList = append(targetAddressList, target.Host)
	}
	return targetAddressList
}

func Task(u Utility) ([]string, []model.NqmTarget, model.NqmAgent, time.Duration, error) {
	/**
	 * Only 2 possible responses come from hbs:
	 *     1. NeedPing==false (default condition)
	 *         NqmAgent, NQMTargets, Measurements are nil
	 *     2. NeedPing==ture
	 *         NqmAgent, NQMTargets, Measurements are not nil
	 */
	hbsResp := GetGeneralConfig().hbsResp.Load().(model.NqmTaskResponse)

	if !hbsResp.NeedPing {
		return nil, nil, model.NqmAgent{}, GetGeneralConfig().Hbs.Interval, fmt.Errorf("[ " + u.UtilName() + " ] No tasks assigned.")
	}
	if !hbsResp.Measurements[u.UtilName()].Enabled {
		return nil, nil, model.NqmAgent{}, GetGeneralConfig().Hbs.Interval, fmt.Errorf("[ " + u.UtilName() + " ] Not enabled.")

	}
	targets := make([]model.NqmTarget, len(hbsResp.Targets))
	copy(targets, hbsResp.Targets)

	agent := *hbsResp.Agent

	command := make([]string, len(hbsResp.Measurements[u.UtilName()].Command))
	copy(command, hbsResp.Measurements[u.UtilName()].Command)

	targetAddressList := getTargetAddressList(targets)
	probingCmd := u.ProbingCommand(command, targetAddressList)
	return probingCmd, targets, agent, hbsResp.Measurements[u.UtilName()].Interval, nil
}
