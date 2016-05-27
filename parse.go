package main

import "strings"

func parseOneRow(row string) []string {
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
		parsedRow := parseOneRow(row)
		parsedRows = append(parsedRows, parsedRow)
	}
	return parsedRows
}
