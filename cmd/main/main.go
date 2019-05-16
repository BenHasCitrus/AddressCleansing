//Testing git
package main

import (
	"fmt"

	getconfig "github.com/hasanul.benna/AddressCleansing/cmd/config"
	"github.com/hasanul.benna/AddressCleansing/cmd/fileoperations"
	"github.com/hasanul.benna/AddressCleansing/cmd/webservice"
)

func main() {

	// get configuration
	configFilePath := "config.json"
	configFile, err := getconfig.GetConfiguration(configFilePath)

	if err != nil {

		fmt.Println(err.Error())
		return

	}

	var sourceFile = configFile.SourceFile
	var targetLocation = configFile.OutputLocation

	// Get data from file
	clientDatas, err := fileoperations.GetClientDataFromFilePath(sourceFile)

	if err != nil {

		fmt.Println(err.Error())
		return

	}

	//Update data using service
	clientDatasWeb, err := webservice.GetDataFromService(clientDatas, configFile)

	if err != nil {

		fmt.Printf(err.Error())

	}

	//Update data to filepath

	err = fileoperations.WriteDatatoFile(targetLocation, clientDatasWeb)

	if err != nil {

		fmt.Println(err.Error())

	}

}
