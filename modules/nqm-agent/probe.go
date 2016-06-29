package main

import (
	"log"
	"os/exec"
	"strings"
)

func trimResults(results []string) []string {
	var trimmedData []string
	for _, result := range results {
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
	results := strings.Split(string(cmdOutput), "\n")
	results = results[:len(results)-1]
	rawData := trimResults(results)
	return rawData
}
