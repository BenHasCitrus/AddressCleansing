package getconfig

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//ConfigFileStruct stores configurable fields
type ConfigFileStruct struct {
	SourceFile     string
	OutputLocation string
	LogLocation    string
	AuthID         string
	AuthToken      string
	Candidates     string
}

// GetConfiguration to provide configuration fields extracted from given path
func GetConfiguration(filepath string) (ConfigFileStruct, error) {

	configFileData := ConfigFileStruct{}

	// opening file
	configfile, err := os.OpenFile(filepath, os.O_RDWR, os.ModePerm)
	if err != nil {

		return configFileData, err

	}

	byteValue, _ := ioutil.ReadAll(configfile)

	//map to json
	//var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &configFileData)

	if err != nil {

		return configFileData, err

	}

	// configVariable.SourceFile = fmt.Sprint(result["SourceFile"])
	// configVariable.OutputLocation= fmt.Sprint(result["OutputFile"])
	// configVariable.AuthID = fmt.Sprint(result["Auth_id"])
	// configVariable.AuthToken = fmt.Sprint(result["Auth_token"])
	// configVariable.Candidates = fmt.Sprint(result["candidates"])

	return configFileData, nil

}
