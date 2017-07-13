package utils

import (
	. "gopkg.in/check.v1"
)

type TestGroupingSuite struct{}

var _ = Suite(&TestGroupingSuite{})

type sampleKey string

func (k sampleKey) GetKey() interface{} {
	return k
}

// Tests the putting and retrieving of grouping data
func (suite *TestGroupingSuite) TestGroupingProcessor(c *C) {
	testedProcessor := NewGroupingProcessorOfTargetType(int(0))

	testedProcessor.Put(sampleKey("GD-1"), 20)
	testedProcessor.Put(sampleKey("GD-1"), 30)
	testedProcessor.Put(sampleKey("GD-1"), 40)
	testedProcessor.Put(sampleKey("GD-2"), 70)
	testedProcessor.Put(sampleKey("GD-2"), 80)

	c.Assert(testedProcessor.Keys(), HasLen, 2)
	c.Assert(testedProcessor.KeyObject(sampleKey("GD-1")), Equals, sampleKey("GD-1"))
	c.Assert(testedProcessor.KeyObject(sampleKey("GD-2")), Equals, sampleKey("GD-2"))
	c.Assert(testedProcessor.Children(sampleKey("GD-1")), HasLen, 3)
	c.Assert(testedProcessor.Children(sampleKey("GD-2")), HasLen, 2)
}

// Tests the type conversion of grouping
func (suite *TestGroupingSuite) TestTypedValues(c *C) {
	testedProcessor := NewGroupingProcessorOfTargetType(int(0))

	testedProcessor.Put(sampleKey("GD-1"), 20)
	testedProcessor.Put(sampleKey("GD-1"), 30)

	typedKey := testedProcessor.KeyObject(sampleKey("GD-1")).(sampleKey)
	c.Assert(typedKey, Equals, sampleKey("GD-1"))

	intValues := testedProcessor.Children(sampleKey("GD-1")).([]int)
	c.Assert(intValues, HasLen, 2)
}
