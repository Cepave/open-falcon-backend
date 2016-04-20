package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func trimFpingResults(fpingResults []string) []string {
	var trimmedData []string
	for i, result := range fpingResults {
		if strings.HasPrefix(result, "ICMP Time Exceeded from") {
			fmt.Println("Result - ICMP Time Exceeded", i+1, ":", result)
		} else {
			trimmedData = append(trimmedData, result)
			fmt.Println("Result", i+1, ":", result)
		}
	}
	return trimmedData
}

func Probe(probingCmd []string) []string {
	cmdOutput, err := exec.Command(probingCmd[0], probingCmd[1:]...).CombinedOutput()
	if err != nil {
		log.Println("An error occured:", string(cmdOutput), err)
	}
	fpingResults := strings.Split(string(cmdOutput), "\n")
	fpingResults = fpingResults[:len(fpingResults)-1]
	rawData := trimFpingResults(fpingResults)
	for i, result := range rawData {
		fmt.Print("Result ", i+1, ": ")
		fmt.Println(result)
	}
	return rawData
}
