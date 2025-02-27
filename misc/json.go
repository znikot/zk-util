package misc

import (
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-json"
)

// marshal obj to json and print to stdout
func PrintJSON(obj any, beauty bool) {
	var str []byte
	var err error

	if beauty {
		str, err = json.MarshalIndent(obj, "", " ")
	} else {
		str, err = json.Marshal(obj)
	}

	if err != nil {
		fmt.Printf("marshal json failed: %v\n", err)
		return
	}
	fmt.Printf("%s\n", str)
}

// marshal obj to json and print to stdout with pretty
func PrintJSONPretty(obj any) {
	PrintJSON(obj, true)
}

// read data from reader and unmarshal to obj
func ReadJSON(reader io.Reader, obj any) (err error) {
	err = json.NewDecoder(reader).Decode(obj)

	return
}

// read data from json file and unmarshal to v
//
// filePath will resolve width function ResolveFilePath(string)
func ReadJSONFile(filePath string, v any) error {
	jsonFile, err := os.Open(ResolveFilePath(filePath))
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	return ReadJSON(jsonFile, v)
}

// marshal anything to json string
// avoid err return value of json.Marshal
//
//	json.Marshal:
//
//	str,_ := json.Marshal(obj)
//	doSomethingWithStr(str)
//
// ToJSON:
//
//	doSomethingWithStr(ToJSON(obj))
func ToJSON(obj any) string {
	str, _ := json.Marshal(obj)
	return string(str)
}
