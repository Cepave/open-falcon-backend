package main

import (
	"log"
	"os/exec"
	"strings"
)

func trimResults(fpingResults []string) []string {
	var trimmedData []string
	for _, result := range fpingResults {
		if strings.HasPrefix(result, "ICMP Time Exceeded from") {
			continue
		}
		trimmedData = append(trimmedData, result)
	}
	return trimmedData
}

func Probe(probingCmd []string, util string) []string {
	cmdOutput, err := exec.Command(probingCmd[0], probingCmd[1:]...).CombinedOutput()
	if err != nil {
		// fping output 'exit status 1' when there is at least
		// one target with 100% packet loss.
		log.Println("[", util, "] An error occured:", err)
	}
	fpingResults := strings.Split(string(cmdOutput), "\n")
	fpingResults = fpingResults[:len(fpingResults)-1]
	rawData := trimResults(fpingResults)
	return rawData
}
