package cron

import (
	"fmt"
	"github.com/open-falcon/alarm/g"
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/common/utils"
)

func BuildCommonSMSContent(event *model.Event) string {
	return fmt.Sprintf(
		"[P%d][%s][%s][][%s %s %s %s %s%s%s][O%d %s]",
		event.Priority(),
		event.Status,
		event.Endpoint,
		event.Note(),
		event.Func(),
		event.Metric(),
		utils.SortedTags(event.PushedTags),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.CurrentStep,
		event.FormattedTime(),
	)
}

func BuildCommonMailContent(event *model.Event) string {
	link := g.Link(event)
	return fmt.Sprintf(
		`<html><body><table border="1" cellpadding="10" cellspacing="0" style="border: 1px solid #ccc;">
 			<tr style="border: 1px solid #ccc;">
    				<td colspan="row">%s</td>
    				<td colspan="row">%d</tr>
  			<tr style="border: 1px solid #ccc;">
    				<td scope="row">Endpoint:</td>
    				<td scope="row">%s</td>
			</tr>
  			<tr style="border: 1px solid #ccc;">
    				<td scope="row">Metric:</td>
    				<td align="left">%s</td>
  			</tr>
  			<tr style="border: 1px solid #ccc;">
    				<td scope="row">Tags:</td>
    				<td align="left">%s</td>
  			</tr>
  			<tr style="border: 1px solid #ccc;">
    				<td scope="row">%s</td>
    				<td align="left">%s%s%s</td>
  			</tr>
  			<tr style="border: 1px solid #ccc;">
    				<td scope="row">Note:</td>
    				<td align="left">%s</td>
  			</tr>
  			<tr style="border: 1px solid #ccc;">
    				<td scope="row">Max:</td>
    				<td align="left">%d</td>
  			</tr>
  			<tr style="border: 1px solid #ccc;">
    				<td scope="row">Current:</td>
    				<td align="left">%d</td>
  			</tr>
  			<tr style="border: 1px solid #ccc;">
    				<td scope="row">Timesramp:</td>
    				<td align="left">%s</td>
  			</tr>
		</table>
		<br>
		<a href="%s">%s</a></body></html>`,
		event.Status,
		event.Priority(),
		event.Endpoint,
		event.Metric(),
		utils.SortedTags(event.PushedTags),
		event.Func(),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.Note(),
		event.MaxStep(),
		event.CurrentStep,
		event.FormattedTime(),
		link,
		link,
	)
}

func GenerateSmsContent(event *model.Event) string {
	return BuildCommonSMSContent(event)
}

func GenerateMailContent(event *model.Event) string {
	return BuildCommonMailContent(event)
}
