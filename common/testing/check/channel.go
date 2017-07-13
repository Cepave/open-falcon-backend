package check

import (
	"fmt"
	"reflect"

	"gopkg.in/check.v1"
)

// This checking would check:
//
// 	1. The len() for both of the channels and
// 	2. Every element(be put back) in both of the channels
var ChannelEquals = channelEquals(true)

type channelEquals bool

func (c channelEquals) Info() *check.CheckerInfo {
	return &check.CheckerInfo{
		Name:   "ChannelEquals",
		Params: []string{"obtained", "expected"},
	}
}
func (c channelEquals) Check(params []interface{}, names []string) (bool, string) {
	checkedChannel := reflect.ValueOf(params[0])
	expectedChannel := reflect.ValueOf(params[1])

	expectedLen := expectedChannel.Len()
	if checkedChannel.Len() != expectedLen {
		return false, fmt.Sprintf(
			"Size of channels are not same. Obtained: [%d]. Expected: [%d]",
			checkedChannel.Len(), expectedLen,
		)
	}

	for i := 0; i < expectedLen; i++ {
		expected, _ := expectedChannel.TryRecv()
		checked, _ := checkedChannel.TryRecv()

		if !reflect.DeepEqual(expected.Interface(), checked.Interface()) {
			return false, fmt.Sprintf(
				"The [%d] element of channels are not same. Obtained: [%#v]. Expected: [%#v]",
				i, checked.Interface(), expected.Interface(),
			)
		}

		/**
		 * Put the element back to the channel
		 */
		expectedChannel.Send(expected)
		checkedChannel.Send(checked)
		// :~)
	}

	return true, ""
}
