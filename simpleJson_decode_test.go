package simplejson

import "testing"

func Test_checkValid(t *testing.T) {
	testJson := `{"A":"B"}`
	result := checkValid([]byte(testJson))

	if result != nil {
		t.Fatal(result)
	}

	testJson = `{"A": true}`
	result = checkValid([]byte(testJson))

	if result != nil {
		t.Fatal(result)
	}

	testJson = `{"A": {"testOne":"testTwo"}, "B":"C", "num":10}`
	result = checkValid([]byte(testJson))

	if result != nil {
		t.Fatal(result)
	}
}
