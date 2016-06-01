package main

import "fmt"

func getTargetAddressList() []string {
	var targetAddressList []string
	for _, target := range GetGeneralConfig().hbsResp.Targets {
		targetAddressList = append(targetAddressList, target.Host)
	}
	return targetAddressList
}

func Task(u Utility) ([]string, error) {
	/**
	 * Only 2 possible responses come from hbs:
	 *     1. NeedPing==false (default condition)
	 *         NqmAgent, NQMTargets, Command are nil
	 *     2. NeedPing==ture
	 *         NqmAgent, NQMTargets, Command are not nil
	 */
	if !GetGeneralConfig().hbsResp.NeedPing {
		return nil, fmt.Errorf("[ " + u.UtilName() + " ] No tasks assigned.")
	}

	targetAddressList := getTargetAddressList()
	probingCmd := u.ProbingCommand(targetAddressList)
	return probingCmd, nil
}
