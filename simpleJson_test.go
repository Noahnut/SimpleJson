package simplejson

import (
	"log"
	"testing"
)

type testStructThree struct {
	G int
	K string
}

type testStructTwo struct {
	D    string
	E    string
	test testStructThree
}

type testStruct struct {
	A int `json:"test"`
	B string
	C string
	E testStructTwo
	F []string
	G []int
}

func TestSimpleJson(t *testing.T) {
	ts := testStruct{
		A: 1,
		B: "D",
		C: "E",
		E: testStructTwo{D: "F", E: "H", test: testStructThree{
			G: 10,
			K: "test",
		}},
		F: []string{"A", "B"},
		G: []int{1, 2},
	}
	result, _ := Marshal(ts)
	log.Println(string(result))
	if string(result) != `{"test":1,"B":"D","C":"E","E":{"D":"F","E":"H","test":{"G":10,"K":"test"}},"F":["A","B"],"G":[1,2]}` {
		t.Error("Marshal fail")
	}
}
