package main

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
)

func convToFloat(samples []string) []float64 {
	var floatData []float64

	for _, sample := range samples {
		if sample != "-" {
			rtt, err := strconv.ParseFloat(sample, 64)
			if err != nil {
				log.Println("error occured:", err)
			} else {
				floatData = append(floatData, rtt)
			}
		}
	}
	return floatData
}

func calcRow(parsedRow []string, u Utility) map[string]string {
	/*
		    assume fping command looks like:
		        fping -p 20 -i 10 -C 5 -a www.google.com www.yahoo.com
		    input argument row looks like:
				www.yahoo.com  6.72 29.08 8.55 7.40 - 6.26
				0                1   2     3     4  5   6   ....  n
	*/
	samples := parsedRow[1:]
	row := convToFloat(samples)
	return u.CalcStats(row, len(samples))
}

func Calc(parsedData [][]string, u Utility) []map[string]string {
	var statsData []map[string]string
	for _, parsedRow := range parsedData {
		statsRow := calcRow(parsedRow, u)
		statsData = append(statsData, statsRow)
	}
	return statsData
}
