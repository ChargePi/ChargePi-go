package settings

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"regexp"
	"strings"
)

const (
	YamlFile = "yaml"
	YmlFile  = "yml"
	JSON     = "json"
)

var ErrUnsupportedFileFormat = errors.New("unsupported file format")

func WriteToFile(filename string, structure interface{}) (err error) {
	var (
		encodingType          string
		marshal               []byte
		splitFile             = strings.Split(filename, ".")
		isValidFile, matchErr = regexp.MatchString("^.*\\.(json|yaml|yml)$", filename)
	)

	if matchErr != nil {
		return matchErr
	}

	// Check if the file format is supported
	if len(splitFile) > 0 && isValidFile {
		encodingType = splitFile[len(splitFile)-1]
	}

	switch encodingType {
	case YamlFile, YmlFile:
		marshal, err = yaml.Marshal(&structure)
		break
	case JSON:
		marshal, err = json.MarshalIndent(&structure, "", "\t")
		break
	default:
		return ErrUnsupportedFileFormat
	}

	return ioutil.WriteFile(filename, marshal, 0644)
}
