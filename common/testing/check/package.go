//
// This package provindes extensions to "gopkg.in/check.v1".
//
// The build-in checkers are defined as public variables.
//
// ChannelEquals
//
// This checker checks two channels that if they are same(for both len() and content).
//
// 	c.Assert(checkedChannel, check.JsonEquals, expectedChannel)
//
// JsonEquals
//
// This checker checks JSON content(supports various types).
//
// 	c.Assert("[3, 4, 5]", check.JsonEquals, "[3, 4, 5")
// 	c.Assert("[3, 4, 5]", check.JsonEquals, []byte("[3, 4, 5]"))
// 	c.Assert(jsonObject, check.JsonEquals, "[3, 4, 5]")
// 	c.Assert(anyObject, check.JsonEquals, "[3, 4, 5]")
//
// Number checkers
//
// These number checkers supports comparison between numbers without enforcement on types.
//
// For example: compares "int16(34)" and "int32(77)"
//
// 	LargerThan - obtained > expected
// 	c.Assert(obtained, check.LargerThan, expected)
//
// 	LargerThanOrequalTo - obtained >= expected
// 	c.Assert(obtained, check.LargerThanOrEqualTo, expected)
//
// 	SmallerThan - obtained < expected
// 	c.Assert(obtained, check.SmallerThan, expected)
//
// 	SmallerThanOrequalTo - obtained <= expected
// 	c.Assert(obtained, check.SmallerThanOrEqualTo, expected)
//
// StringContains
//
// Checks if a sub-string is contained in checked string.
//
// 	c.Assert("Hello World", check.StringContains, "llo")
//
// Time checkers
//
// The time checkers only checks value of UNIX time, the time zone is ignored.
//
//	TimeEquals - Equivalence of two values of time
// 	c.Assert(obtainedTime, TimeEquals, expectedTime)
//
//	TimeBefore - The obtained one must be earlier the expected one
// 	c.Assert(obtainedTime, TimeEquals, expectedTime)
//
//	TimeAfter - The obtained one must be latter the expected one
// 	c.Assert(obtainedTime, TimeEquals, expectedTime)
//
// ViableValue
//
// This checker checks if a value is viable.
//
// See "utils.ValueExt.IsViable()" for detail.
package check
