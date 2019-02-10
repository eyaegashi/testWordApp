package util

import (
	"io/ioutil"
	"os"
)

// LoadjsonFile is to load Jsonfile and prepare to convert json data.
func LoadjsonFile(file string) (jsonData []byte, err error) {
	// open and close json file
	jsonFile, err := os.Open(file)
	if err != nil {
		return []byte{}, err
	}

	defer jsonFile.Close()

	// read json file
	jsonData, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		return []byte{}, err
	}
	return jsonData, err
}
