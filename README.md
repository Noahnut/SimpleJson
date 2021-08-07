# Simple JSON Encode and Decode

Simple encode and decode JSON operation, encode the struct object to the string or decode the string to the struct object also provide the tag 


## Usage 
### Encode
```golang
type testStructTwo struct {
	G int
	K string
    A []int
}


type testStruct struct {
	D    string `json:"test"`
	E    string
	test testStructTwo
}

func main() {
    ts := testStruct{
        D: "testOne",
        E: "testTwo", 
        test: testStructTwo{
            G: "testThree",
            K: "testFour",
            A: []int{1, 2},
        }
	}

    result, err := simpleJson.Marshal(ts)

    if err != nil {
        log.Println(err)
    }

    log.Println(result)
}
```