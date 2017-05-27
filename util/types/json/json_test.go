// Copyright 2017 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package json

import (
	"fmt"
	"testing"

	. "github.com/pingcap/check"
)

var _ = Suite(&testJSONSuite{})

type testJSONSuite struct{}

func TestT(t *testing.T) {
	TestingT(t)
}

func parseFromStringPanic(s string) JSON {
	j, err := ParseFromString(s)
	if err != nil {
		msg := fmt.Sprintf("ParseFromString(%s) fail", s)
		panic(msg)
	}
	return j
}

func (s *testJSONSuite) TestParseFromString(c *C) {
	jstr1 := `{"a": [1, "2", {"aa": "bb"}, 4, null], "b": true, "c": null}`
	jstr2 := parseFromStringPanic(jstr1).String()
	c.Assert(jstr2, Equals, `{"a":[1,"2",{"aa":"bb"},4,null],"b":true,"c":null}`)
}

func (s *testJSONSuite) TestJSONSerde(c *C) {
	var jsonNilValue = CreateJSON(nil)
	var jsonBoolValue = CreateJSON(true)
	var jsonDoubleValue = CreateJSON(3.24)
	var jsonStringValue = CreateJSON("hello, 世界")
	j1 := parseFromStringPanic(`{"aaaaaaaaaaa": [1, "2", {"aa": "bb"}, 4.0], "bbbbbbbbbb": true, "ccccccccc": "d"}`)
	j2 := parseFromStringPanic(`[{"a": 1, "b": true}, 3, 3.5, "hello, world", null, true]`)

	var testcses = []struct {
		In  JSON
		Out JSON
	}{
		{In: jsonNilValue, Out: jsonNilValue},
		{In: jsonBoolValue, Out: jsonBoolValue},
		{In: jsonDoubleValue, Out: jsonDoubleValue},
		{In: jsonStringValue, Out: jsonStringValue},
		{In: j1, Out: j1},
		{In: j2, Out: j2},
	}

	for _, s := range testcses {
		data := Serialize(s.In)
		t, err := Deserialize(data)
		c.Assert(err, IsNil)

		v1 := t.String()
		v2 := s.Out.String()
		c.Assert(v1, Equals, v2)
	}
}

func (s *testJSONSuite) TestCompareJSON(c *C) {
	jNull := parseFromStringPanic(`null`)
	jBoolTrue := parseFromStringPanic(`true`)
	jBoolFalse := parseFromStringPanic(`false`)
	jIntegerLarge := parseFromStringPanic(`5`)
	jIntegerSmall := parseFromStringPanic(`3`)
	jStringLarge := parseFromStringPanic(`"hello, world"`)
	jStringSmall := parseFromStringPanic(`"hello"`)
	jArrayLarge := parseFromStringPanic(`["a", "c"]`)
	jArraySmall := parseFromStringPanic(`["a", "b"]`)
	jObject := parseFromStringPanic(`{"a": "b"}`)

	var caseList = []struct {
		left  JSON
		right JSON
	}{
		{jNull, jIntegerSmall},
		{jIntegerSmall, jIntegerLarge},
		{jIntegerLarge, jStringSmall},
		{jStringSmall, jStringLarge},
		{jStringLarge, jObject},
		{jObject, jArraySmall},
		{jArraySmall, jArrayLarge},
		{jArrayLarge, jBoolFalse},
		{jBoolFalse, jBoolTrue},
	}
	for _, cmpCase := range caseList {
		cmp, err := CompareJSON(cmpCase.left, cmpCase.right)
		c.Assert(err, IsNil)
		c.Assert(cmp < 0, IsTrue)
	}
}