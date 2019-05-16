package fileoperations

import (
	"os"

	"github.com/gocarina/gocsv"
)

//ClientData stores data fetched from csv file - OBTAIN PROPER ADDRESS FROM IT

type ClientData struct {
	CompanyID      int    `csv:"company_id"`
	CompanyName    string `csv:"company_name"`
	ParentCompany  string `csv:"parent_company"`
	IsDealership   bool   `csv:"is_dealership"`
	IsStore        bool   `csv:"is_store"`
	CompanyType    string `csv:"company_type"`
	Address        string `csv:"address"`
	City           string `csv:"city"`
	State          string `csv:"state"`
	Zipcode        string `csv:"zipcode"`
	Stage          string `csv:"stage"`
	Description    string `csv:"description"`
	StatusField    string
	ZipPlus4       string
	Latitude       string
	Longitude      string
	ServiceAddress string
	JSONResponse   string
}

// ClientDatasExport Refer type of client Data
var ClientDatasExport = []*ClientData{}


//GetClientDataFromFilePath takes filename as input and produce client data struct as output
func GetClientDataFromFilePath(filepath string) ([]*ClientData, error) {

	// clientDatas Get client data from csv file
	clientDatas := []*ClientData{}

	clientsFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {

		return clientDatas, err

	}

	if err := gocsv.UnmarshalFile(clientsFile, &clientDatas); err != nil { // Load clients from file

		return clientDatas, err
	}

	defer func() {

		clientsFile.Close()

	}()


	return clientDatas, nil

	
}


//WriteDatatoFile takes filepath and data as input and write it to the output file
func WriteDatatoFile(filePath string,clientDatas []*ClientData)error{


	var OutputPath = filePath

	// making directory if not exist

	if _, err := os.Stat(OutputPath); os.IsNotExist(err) {
		err := os.Mkdir(OutputPath, 0700)

		if err != nil {
			//logger.Fatalf("Failed to create log file: %v", err)

			return err

		}

	}

	// Making output file
	_, err := os.Create(OutputPath + "/output_store-data.csv")

	outputFile, err := os.OpenFile(OutputPath+"/output_store-data.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	
	err = gocsv.MarshalFile(&clientDatas, outputFile)
	if err != nil {

		return err

	}

	defer func() {
		outputFile.Close()

	}()

       return nil

}
