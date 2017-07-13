package digest

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

type TestStructSuite struct{}

var _ = Suite(&TestStructSuite{})

type sampleContainer struct {
	TargetValue interface{} `digest:"1"`
}

func (c *sampleContainer) String() string {
	return fmt.Sprintf("%v", c.TargetValue)
}

// Tests the digesting for struct
func (suite *TestStructSuite) TestDigestDirectVariable(c *C) {
	type stringAlias string
	type boolAlias bool
	type intAlias int
	type int64Alias int64
	type uintAlias uint
	type uint64Alias uint64
	type float64Alias float64
	type complex128Alias complex128

	testCases := []*struct {
		sampleValue  interface{}
		expectedHash string
	}{
		{"cool", "5704c6d2c94e357a8c73792db7c255aa"},
		{stringAlias("cool"), "5704c6d2c94e357a8c73792db7c255aa"},
		{true, "dcefbbcedf7c000634275e5cb68c6795"},
		{boolAlias(true), "dcefbbcedf7c000634275e5cb68c6795"},
		{int(21), "6beac09dcdfb64c4a65ffc6665a503c6"},
		{intAlias(21), "6beac09dcdfb64c4a65ffc6665a503c6"},
		{int8(22), "6a27b0cadeff550ea73118414b527628"},
		{int16(23), "ff2b42176613aa3dd15b1d4ca049a8fe"},
		{int32(24), "bf33b388ebf6dd8a60605657d1a401fc"},
		{int64(25), "9db65919b548200c9fd878810f317819"},
		{int64Alias(25), "9db65919b548200c9fd878810f317819"},
		{uint(31), "fd02fd92c571ad8f456bf792451da647"},
		{uintAlias(31), "fd02fd92c571ad8f456bf792451da647"},
		{uint8(32), "7331dd0ef933e97f6685ca658c0f9c9d"},
		{uint16(33), "1b3b8be6e59169b386220e002281c503"},
		{uint32(34), "c739141a29d0dd3e3088dc2e4b251198"},
		{uint64(35), "c1ed7a417d0398f5b46a7d85c8c1d2ea"},
		{uint64Alias(35), "c1ed7a417d0398f5b46a7d85c8c1d2ea"},
		{float32(35.71), "5d8021925e7368a317c4421f834ff775"},
		{float64(77.81), "4c867d1ceaaf64d081c2b9b11c6a7e05"},
		{float64Alias(77.81), "4c867d1ceaaf64d081c2b9b11c6a7e05"},
		{complex64(4.7651i), "07352a0e5b6ef6ecfb1436e44f6ab634"},
		{complex128(90.871i), "eacabb75c7a2c9fc2ad8c5c2fc184dbb"},
		{complex128Alias(90.871i), "eacabb75c7a2c9fc2ad8c5c2fc184dbb"},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		sampleStruct := &sampleContainer{testCase.sampleValue}
		testedHexValue := digestAndEncodeToString(sampleStruct)

		c.Logf("Source Value: [%v]. Digested Result(by object): [%s].", sampleStruct, testedHexValue)

		c.Check(testedHexValue, Equals, testCase.expectedHash, comment)
	}
}

type sha1String string

func (s sha1String) GetDigest() []byte {
	digest := sha1.Sum([]byte(s))
	return digest[:]
}

type sha512String string

func (s sha512String) GetDigest() []byte {
	digest := sha512.Sum512([]byte(s))
	return digest[:]
}

type mapWithDigest map[string]int

func (m mapWithDigest) GetDigest() []byte {
	return []byte{19, 29, 66, 51, 40, 98, 76, 17}
}

// Tests the digesting for interface{}
func (suite *TestStructSuite) TestDigestInterface(c *C) {
	sampleMap := map[string]int{
		"AB": 20, "GD": 344, "KS": 901,
	}

	testCases := []*struct {
		sampleData   interface{} `digest:"1"`
		expectedHash string
	}{
		/**
		 * Tests the effective of implementing Digestor interface
		 */
		{sha1String("GoGo2"), "7dae61955bdac2701aa9cd0138fdf97f"},
		{sha512String("GoGo2"), "4de33e1688406f03c2f33870afab4a90"},
		{mapWithDigest(sampleMap), "4dbe325ba746d85dcbdb85eefff2ef6b"},
		// :~)
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedHexValue := digestAndEncodeToString(&sampleContainer{testCase.sampleData})
		c.Logf("Source value: \"%v\". Hex: [%s]", testCase.sampleData, testedHexValue)

		c.Check(testedHexValue, Equals, testCase.expectedHash, comment)
	}
}

// Tests the multi-level-pointer of digest
func (suite *TestStructSuite) TestDigestPointers(c *C) {
	var intValue int32 = 8071
	var pointerValue1 *int32 = &intValue
	var pointerValue2 **int32 = &pointerValue1
	var pointerValue3 ***int32 = &pointerValue2
	var pointerNil1 *int32 = nil
	var pointerNil2 **int32 = nil

	intValueHash := digestAndEncodeToString(&sampleContainer{intValue})
	c.Logf("Original Value Hash: [%s]", intValueHash)

	pointerValue1Hash := digestAndEncodeToString(&sampleContainer{pointerValue1})
	c.Logf("Level 1 pointer Hash: [%s]", pointerValue1Hash)
	c.Check(intValueHash, Equals, pointerValue1Hash)

	pointerValue2Hash := digestAndEncodeToString(&sampleContainer{pointerValue2})
	c.Logf("Level 2 pointer Hash: [%s]", pointerValue2Hash)
	c.Check(intValueHash, Equals, pointerValue2Hash)

	pointerValue3Hash := digestAndEncodeToString(&sampleContainer{pointerValue3})
	c.Logf("Level 3 pointer Hash: [%s]", pointerValue3Hash)
	c.Check(intValueHash, Equals, pointerValue3Hash)

	pointerValueNil1Hash := digestAndEncodeToString(&sampleContainer{pointerNil1})
	c.Logf("Nil pointer Hash: [%s]", pointerValueNil1Hash)
	c.Check(pointerValueNil1Hash, Equals, "d41d8cd98f00b204e9800998ecf8427e")

	pointerValueNil2Hash := digestAndEncodeToString(&sampleContainer{pointerNil2})
	c.Logf("Nil pointer Hash: [%s]", pointerValueNil2Hash)
	c.Check(pointerValueNil2Hash, Equals, "d41d8cd98f00b204e9800998ecf8427e")
}

type byteAlias byte

func (a byteAlias) GetDigest() []byte {
	return []byte{10, byte(a), 98}
}

type twiceMd5String string

func (s twiceMd5String) GetDigest() []byte {
	originalBytes := md5.Sum([]byte(s))

	finalBytes := make([]byte, 0, md5.Size+md5.Size)
	finalBytes = append(finalBytes, originalBytes[:]...)
	finalBytes = append(finalBytes, originalBytes[:]...)

	return finalBytes
}

// Test the digesting for arrayg
func (suite *TestStructSuite) TestDigestArray(c *C) {
	type sampleWheel struct {
		Age   int64  `digest:"1"`
		Model string `digest:"2"`
	}

	var iv1, iv2, iv3 int16 = 98, 51, 1045
	var s1, s2, s3 = "cg-1", "bk-2", "zc-3"
	var w1, w2, w3 = sampleWheel{2, "ZK-1"}, sampleWheel{38, "PO-99"}, sampleWheel{76801, "II-78"}

	testCases := []*struct {
		sampleArray  interface{}
		expectedHash string
	}{
		{[]int16{iv1, iv2, iv3}, "272df57244703dce2d5ff8be59bb9833"},
		{[]*int16{&iv1, &iv2, &iv3}, "272df57244703dce2d5ff8be59bb9833"},
		{[]string{"cg-1", "bk-2", "zc-3"}, "5e14a5230f79a13c3c61c64c1a756f44"},
		{[]*string{&s1, &s2, &s3}, "5e14a5230f79a13c3c61c64c1a756f44"},
		{[]byteAlias{9, 13, 19}, "7ca5b6391adb1627348231480c0622d1"},
		{[]*sampleWheel{&w1, &w2, &w3}, "76b33e9ef42a4fd8b9369696bae35251"},
		{[]sampleWheel{w1, w2, w3}, "76b33e9ef42a4fd8b9369696bae35251"},
		{[]byte{81, 21, 9, 9, 76}, "c58d516a810d80f71747735130580df3"},
		{[]twiceMd5String{"Good!", "Nice!"}, "d132d63f1583b30f00f917584d866d7a"},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		testedHexValue := digestAndEncodeToString(&sampleContainer{testCase.sampleArray})
		c.Logf("Array: %v. Hash: [%s]", testCase.sampleArray, testedHexValue)

		c.Check(testedHexValue, Equals, testCase.expectedHash, comment)
	}
}

// Tests the nested struct
func (suite *TestStructSuite) TestNestedStruct(c *C) {
	type engine struct {
		Age  int    `digest:"1"`
		Name string `digest:"2"`
	}
	type smallCar struct {
		WheelLeft  *engine `digest:"1"`
		WheelRight *engine `digest:"2"`
		WheelNil   *engine `digest:"3"`
	}
	type superCar struct {
		CarOfLeft  *smallCar `digest:"1"`
		CarOfRight *smallCar `digest:"2"`
	}

	sampleStruct := &superCar{
		&smallCar{
			&engine{38, "ok"},
			&engine{44, "no"},
			nil, // Tested nil value of struct
		},
		&smallCar{
			&engine{7, "red"},
			&engine{9, "green"},
			nil, // Tested nil value of struct
		},
	}

	testedHexValue := digestAndEncodeToString(sampleStruct)
	c.Logf("Nested struct: %s", testedHexValue)
	c.Check(testedHexValue, Equals, "a5ac2890cbb114a482e3bc1fd809fe75")
}

// Tests the various form which make sure the digest are same
func (suite *TestStructSuite) TestSameDigestWithDifferentForms(c *C) {
	type elf struct {
		Name string `digest:"1"`
		Age  int32  `digest:"2"`

		Color byte

		bowRange int
	}
	type elf2 struct {
		Age  int32  `digest:"2"`
		Name string `digest:"1"`
	}
	type human struct {
		Name string `digest:"1"`
		Age  int32  `digest:"2"`

		Exp uint16

		weight int
	}
	type human2 struct {
		Name string `digest:"1"`
		Age  int32  `digest:"2"`

		Value1 *int32   `digest:"3"`
		Value2 []int32  `digest:"4"`
		Value3 []int32  `digest:"5"`
		Value4 [0]int32 `digest:"6"`
	}

	type stringAlias string

	arrayOfValues := []int32{10, 20, 45}
	arrayOfPointers := []*int32{
		&arrayOfValues[0],
		&arrayOfValues[1],
		&arrayOfValues[2],
	}

	testCases := []*struct {
		leftData  interface{}
		rightData interface{}
	}{
		// Same fields
		{&elf{Name: "John", Age: 48}, &human{Name: "John", Age: 48}},
		// Checks the effective of sequence on tag value
		{&elf{Name: "Lif", Age: 19}, &elf2{Name: "Lif", Age: 19}},
		// Pointer/Object
		{&elf{Name: "Alice", Age: 32}, human{Name: "Alice", Age: 32}},
		// Array of same type
		{&sampleContainer{[]string{"TB-01", "TB-02"}}, &sampleContainer{[]stringAlias{"TB-01", "TB-02"}}},
		// Array of object/pointers
		{&sampleContainer{arrayOfValues}, &sampleContainer{arrayOfPointers}},
		// Nil/empty array values
		{&human{Name: "Steve", Age: 76},
			&human2{
				Name: "Steve", Age: 76,
				Value1: nil, Value2: nil, Value3: []int32{},
			},
		},
	}

	for i, testCase := range testCases {
		comment := Commentf("Test Case: %d", i+1)

		leftHex := digestAndEncodeToString(testCase.leftData)
		rightHex := digestAndEncodeToString(testCase.rightData)

		c.Logf("Left hash: [%s]. Right hash: [%s]", leftHex, rightHex)

		c.Check(leftHex, Equals, rightHex, comment)
	}
}

type listOfInt32 []int32

func (l listOfInt32) GetDigest() []byte {
	return GetBytesGetter([]int32(l), Md5SumFunc)()
}

// Tests the usage of GetBytesGetter to build customized digesting
func (suite *TestStructSuite) TestGetBytesGetter(c *C) {
	var sampleData listOfInt32 = []int32{9081, 765142, 33045, -176, -9878}

	testedHash := hex.EncodeToString(sampleData.GetDigest())
	c.Logf("Digestor: [%s]", testedHash)

	c.Assert(testedHash, Equals, "00002379000bacd600008115ffffff50ffffd96a")
}

func digestAndEncodeToString(v interface{}) string {
	hashValue := DigestStruct(v, Md5SumFunc)
	return hex.EncodeToString(hashValue)
}

func (s *TestStructSuite) SetUpSuite(c *C) {
	Logger.Level = log.DebugLevel
}
func (s *TestStructSuite) TearDownSuite(c *C) {
	Logger.Level = log.WarnLevel
}
