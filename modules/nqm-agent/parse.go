package main

import "strings"

func parseRow(row string) []string {
	return strings.FieldsFunc(row, func(r rune) bool {
		switch r {
		case ' ', '\n', ':', '/', '%', '=', ',':
			return true
		}
		return false
	})
}

func Parse(rawData []string) [][]string {
	var parsedRows [][]string
	for _, row := range rawData {
		parsedRow := parseRow(row)
		parsedRows = append(parsedRows, parsedRow)
	}
	return parsedRows
}
