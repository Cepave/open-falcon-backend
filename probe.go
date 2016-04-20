package main

import (
	"log"
	"os/exec"
	"strings"
)

func trimFpingResults(fpingResults []string) []string {
	var trimmedData []string
	for _, result := range fpingResults {
		if strings.HasPrefix(result, "ICMP Time Exceeded from") {
			continue
		}
		trimmedData = append(trimmedData, result)
	}
	return trimmedData
}

func Probe(probingCmd []string) []string {
	cmdOutput, err := exec.Command(probingCmd[0], probingCmd[1:]...).CombinedOutput()
	if err != nil {
		// 'exit status 1' happens when there is at least
		// one target with 100% packet loss.
		log.Println("An error occured:", err)
	}
	fpingResults := strings.Split(string(cmdOutput), "\n")
	fpingResults = fpingResults[:len(fpingResults)-1]
	rawData := trimFpingResults(fpingResults)
	return rawData
}
