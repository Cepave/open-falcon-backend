package nqm

import (
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

var regexpTimeInDay = regexp.MustCompile("^(\\d\\d):(\\d\\d)$")

func ValidateTimeWithUnit(sl validator.StructLevel) {
	timeUnit := sl.Current().Interface().(TimeWithUnit)

	/**
	 * Checks if the range value are provided with start/end range
	 */
	nilCheck := 0
	if timeUnit.StartTimeOfDay != nil {
		nilCheck++
	}
	if timeUnit.EndTimeOfDay != nil {
		nilCheck++
	}
	if nilCheck == 1 {
		if timeUnit.StartTimeOfDay == nil {
			sl.ReportError(timeUnit.StartTimeOfDay, "StartTimeOfDay", "", "Not nil", "")
		} else if timeUnit.EndTimeOfDay == nil {
			sl.ReportError(timeUnit.EndTimeOfDay, "EndTimeOfDay", "", "Not nil", "")
		}
		return
	}
	// :~)

	if nilCheck == 2 {
		if !regexpTimeInDay.MatchString(*timeUnit.StartTimeOfDay) {
			sl.ReportError(timeUnit.StartTimeOfDay, "StartTimeOfDay", "", "Match 00:00", "")
		}
		if !regexpTimeInDay.MatchString(*timeUnit.EndTimeOfDay) {
			sl.ReportError(timeUnit.EndTimeOfDay, "EndTimeOfDay", "", "Match 00:00", "")
		}

		switch timeUnit.Unit {
		case TimeUnitYear, TimeUnitMonth, TimeUnitWeek, TimeUnitDay:
		default:
			sl.ReportError(timeUnit.Unit, "Unit", "", "Need to be day or larger unit of time", "")
		}
	}
}
