package smartystreet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	extract "github.com/smartystreets/smartystreets-go-sdk/us-extract-api"
	zipcode "github.com/smartystreets/smartystreets-go-sdk/us-zipcode-api"
	"github.com/smartystreets/smartystreets-go-sdk/wireup"

	config "github.com/hasanul.benna/AddressCleansing/cmd/config"
	fileOperations "github.com/hasanul.benna/AddressCleansing/cmd/fileoperations"
)

//GetDataFromAddress - Get details from given file details
func GetDataFromFields(clientDatasInput []*fileOperations.ClientData, config config.ConfigFileStruct) {
	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildUSZIPCodeAPIClient(
		wireup.SecretKeyCredential(config.AuthID, config.AuthToken),
		// wireup.DebugHTTPOutput(), // uncomment this line to see detailed HTTP request/response information.
	)

	// Documentation for input fields can be found at:
	// https://smartystreets.com/docs/us-zipcode-api#input-fields

	lookup1 := &zipcode.Lookup{
		InputID: "dfc33cb6-829e-4fea-aa1b-b6d6580f0817", // Optional ID from your system
		City:    "PROVO",
		State:   "UT",
		ZIPCode: "84604",
	}

	lookup2 := &zipcode.Lookup{
		InputID: "01189998819991197253",
		ZIPCode: "90210",
	}

	batch := zipcode.NewBatch()
	batch.Append(lookup1)
	batch.Append(lookup2)

	fmt.Println("\nBatch full, preparing to send inputs:", batch.Length())

	if err := client.SendBatch(batch); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	for i, input := range batch.Records() {
		fmt.Println("Input:", i, input.City, input.State, input.ZIPCode)
		fmt.Printf("%#v\n", input.Result)
		fmt.Println()
	}

	log.Println("OK")
}

//getDataFromAddress - Get details from given street
func getDataFromAddress(clientDatasInput []*fileOperations.ClientData, config config.ConfigFileStruct) {

	log.SetFlags(log.Ltime | log.Llongfile)

	client := wireup.BuildUSExtractAPIClient(
		wireup.SecretKeyCredential(config.AuthID, config.AuthToken),
		// wireup.DebugHTTPOutput(), // uncomment this line to see detailed HTTP request/response information.
	)

	// Documentation for input fields can be found at:
	// https://smartystreets.com/docs/cloud/us-extract-api#http-request-input-fields

	lookup := &extract.Lookup{
		Text:                    "Meet me at 3214 N University Ave Provo UT 84604 just after 3pm.",
		Aggressive:              true,
		AddressesWithLineBreaks: false,
		AddressesPerLine:        1,
	}

	if err := client.SendLookup(lookup); err != nil {
		log.Fatal("Error sending batch:", err)
	}

	fmt.Println(DumpJSON(lookup))

	log.Println("OK")

}

//getDataFromDetails - provide all the given details
func getDataFromDetails(clientDatasInput []*fileOperations.ClientData, config config.ConfigFileStruct) {

}

func DumpJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}

	var indent bytes.Buffer
	err = json.Indent(&indent, b, "", "  ")
	if err != nil {
		return err.Error()
	}
	return indent.String()
}
