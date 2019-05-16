package webservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	config "github.com/hasanul.benna/AddressCleansing/cmd/config"
	fileOperations "github.com/hasanul.benna/AddressCleansing/cmd/fileoperations"
)

//structClientAddressFields stores data necessary for service call

type structClientAddressFields struct {
	Street    string
	City      string
	State     string
	Zipcode   string
	Component string
}

//streetdata stores data obtained from service call

type streetdata struct {
	InputIndex           int
	CandidateIndex       int    `json:"candidate_index"`
	DeliveryLine1        string `json:"delivery_line_1"`
	LastLine             string `json:"last_line"`
	DeliveryPointBarcode string `json:"delivery_point_barcode"`
	Components           struct {
		PrimaryNumber           string `json:"primary_number"`
		StreetPredirection      string `json:"street_predirection"`
		StreetName              string `json:"street_name"`
		StreetSuffix            string `json:"street_suffix"`
		CityName                string `json:"city_name"`
		DefaultCityName         string `json:"default_city_name"`
		StateAbbreviation       string `json:"state_abbreviation"`
		Zipcode                 string `json:"zipcode"`
		Plus4Code               string `json:"plus4_code"`
		DeliveryPoint           string `json:"delivery_point"`
		DeliveryPointCheckDigit string `json:"delivery_point_check_digit"`
	} `json:"components"`
	Metadata struct {
		RecordType            string  `json:"record_type"`
		ZipType               string  `json:"zip_type"`
		CountyFips            string  `json:"county_fips"`
		CountyName            string  `json:"county_name"`
		CarrierRoute          string  `json:"carrier_route"`
		CongressionalDistrict string  `json:"congressional_district"`
		Rdi                   string  `json:"rdi"`
		ElotSequence          string  `json:"elot_sequence"`
		ElotSort              string  `json:"elot_sort"`
		Latitude              float64 `json:"latitude"`
		Longitude             float64 `json:"longitude"`
		Precision             string  `json:"precision"`
		TimeZone              string  `json:"time_zone"`
		UtcOffset             int     `json:"utc_offset"`
	} `json:"metadata"`
	Analysis struct {
		DpvMatchCode string `json:"dpv_match_code"`
		DpvFootnotes string `json:"dpv_footnotes"`
		DpvCmra      string `json:"dpv_cmra"`
		DpvVacant    string `json:"dpv_vacant"`
		Active       string `json:"active"`
		Footnotes    string `json:"footnotes"`
	} `json:"analysis"`
}

//clientAddressFields Get individual address fields
var clientAddressFields = []structClientAddressFields{}

//GetDataFromService used to call data service and store in clientData field
func GetDataFromService(clientDatasInput []*fileOperations.ClientData, config config.ConfigFileStruct) ([]*fileOperations.ClientData, error) {

	for _, clientDataItem := range clientDatasInput {

		itemStreet := clientDataItem.Address
		//itemCity := arrayStreetCityItems[1]
		itemCity := clientDataItem.City
		itemState := clientDataItem.State
		itemZipcode := clientDataItem.Zipcode

		obtainedItemStruct := structClientAddressFields{Street: itemStreet, City: itemCity, State: itemState, Zipcode: itemZipcode}

		clientAddressFields = append(clientAddressFields, obtainedItemStruct)

	}

	for itemRow, urlFields := range clientAddressFields {

		urlFields.City = strings.TrimSpace(urlFields.City)
		urlFields.State = strings.TrimSpace(urlFields.State)
		urlFields.Street = strings.TrimSpace(urlFields.Street)
		urlFields.Zipcode = strings.TrimSpace(urlFields.Zipcode)

		// fmt.Println("city   ----->", urlFields.City)
		// fmt.Println("state ------>", urlFields.State)
		// fmt.Println("street ------->", urlFields.Street)
		// fmt.Println("zipcode ------->", urlFields.Zipcode)

		//Fetrch from Url

		URL1 := fmt.Sprintf("https://us-street.api.smartystreets.com/street-address?auth-id=%s&auth-token=%s&candidates=%s", config.AuthID, config.AuthToken, config.Candidates)

		var itemValues = ""

		if urlFields.Street != "" {

			itemValues = fmt.Sprintf("%s&street=%s", itemValues, urlFields.Street)

			//URL1 = fmt.Sprintf("%s&street=%s", URL1, urlFields.Street)

		}

		itemValues = url.PathEscape(itemValues)

		URL1 = fmt.Sprintf("%s%s", URL1, itemValues)
		//URL1 := fmt.Sprintf("https://us-street.api.smartystreets.com/street-address?auth-id=%s&auth-token=%s&candidates=%s&street=3301%20South%20Greenfield%20Rd", authID, authToken, candidates)

		//fmt.Println(URL1)
		res, err := http.Get(URL1)
		if err != nil {

			// Implement loger
			// insertLog := fmt.Sprintf("Failed to load from URL  %s  %s Row - %d", URL1, err.Error(), itemRow)
			// logger.Info(insertLog)
			// panic(err.Error())

			return clientDatasInput, err

		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {

			return clientDatasInput, err

			// clientDatas[itemRow].StatusField = "F"
			// continue
		}

		// convert and store as struct data
		var streetData []streetdata
		err = json.Unmarshal(body, &streetData)

		if err != nil {

			return clientDatasInput, err

		} else { // check values or next service call

			//Implement Updated value details

			var statusStrng = ""

			if len(streetData) == 0 {

				// clientDatas[itemRow].StatusField = "F"
				// continue

				// Got some error :(    has to proceed with remaining address fields
				// FAILED CASE - 1

				if urlFields.City != "" {

					itemValues = fmt.Sprintf("%s&city=%s", itemValues, urlFields.City)

					//URL1 = fmt.Sprintf("%s&city=%s", URL1, urlFields.City)

					//params.Add("city", urlFields.City)

				}
				if urlFields.State != "" {

					itemValues = fmt.Sprintf("%s&state=%s", itemValues, urlFields.State)

					//URL1 = fmt.Sprintf("%s&state=%s", URL1, urlFields.State)

					//params.Add("state", urlFields.State)

				}

				if urlFields.Zipcode != "" {

					itemValues = fmt.Sprintf("%s&state=%s", itemValues, urlFields.Zipcode)
					//URL1 = fmt.Sprintf("%s&zipcode=%s", URL1, urlFields.Zipcode)

					//params.Add("zipcode", urlFields.Zipcode)

				}

				itemValues = url.PathEscape(itemValues)

				URL1 = fmt.Sprintf("%s%s", URL1, itemValues)

				res, err := http.Get(URL1)
				if err != nil {

					// Implement loger
					// insertLog := fmt.Sprintf("Failed to load from URL  %s  %s Row - %d", URL1, err.Error(), itemRow)
					// logger.Info(insertLog)
					// panic(err.Error())
					// fmt.Println(err.Error())
					//break

					return clientDatasInput, err

				}

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					// Implement loger

					return clientDatasInput, err
					//continue
				}

				// convert and store as struct data
				var streetData []streetdata
				err = json.Unmarshal(body, &streetData)

				if err != nil {

					return clientDatasInput, err

					//

				} else { // values obtained

					//Implement Updated value details

					if len(streetData) == 0 {

						fmt.Println("no update in given data")
						clientDatasInput[itemRow].StatusField = "Service Failed"

					} else {

						resJSON, _ := json.Marshal(streetData)
						//fmt.Println(string(res2B))
						clientDatasInput[itemRow].JSONResponse = string(resJSON)

						var statusStrng = ""
						var obtainedField = streetData[0].Components
						obtainedFullItem := streetData[0]

						if urlFields.State != obtainedField.StateAbbreviation {

							if urlFields.State == "" {
								statusStrng += fmt.Sprintf("Empty state --> %s,", obtainedField.StateAbbreviation)

							} else {
								statusStrng += fmt.Sprintf("%s --> %s,", urlFields.State, obtainedField.StateAbbreviation)
							}

							clientDatasInput[itemRow].State = obtainedField.StateAbbreviation

						}
						if urlFields.City != obtainedField.CityName {

							if urlFields.City == "" {
								statusStrng += fmt.Sprintf("Empty City --> %s,", obtainedField.CityName)

							} else {
								statusStrng += fmt.Sprintf("%s --> %s,", urlFields.City, obtainedField.CityName)

							}

							clientDatasInput[itemRow].City = obtainedField.CityName

						}
						if urlFields.Zipcode != obtainedField.Zipcode {

							if urlFields.Zipcode == "" {

								statusStrng += fmt.Sprintf("No zip --> %s", obtainedField.Zipcode)
							} else {
								statusStrng += fmt.Sprintf("%s --> %s", urlFields.Zipcode, obtainedField.Zipcode)

							}

							clientDatasInput[itemRow].Zipcode = obtainedField.Zipcode

						}

						if statusStrng == "" {

							statusStrng = "NO UPDATES"
						}
						clientDatasInput[itemRow].StatusField = statusStrng
						clientDatasInput[itemRow].ZipPlus4 = fmt.Sprintf("%s-%s", obtainedField.Zipcode, obtainedField.Plus4Code)
						clientDatasInput[itemRow].Latitude = fmt.Sprintf("%f", obtainedFullItem.Metadata.Latitude)
						clientDatasInput[itemRow].Longitude = fmt.Sprintf("%f", obtainedFullItem.Metadata.Longitude)

						CombinedserviceAddress := fmt.Sprintf("%s,%s,%s", obtainedFullItem.DeliveryLine1, obtainedFullItem.LastLine, obtainedFullItem.DeliveryPointBarcode)

						clientDatasInput[itemRow].ServiceAddress = CombinedserviceAddress

					}
				}

				// obatained values from first service itself - Update if any
			} else {

				resJSON, _ := json.Marshal(streetData)
				//fmt.Println(string(res2B))
				clientDatasInput[itemRow].JSONResponse = string(resJSON)

				var obtainedField = streetData[0].Components
				obtainedFullItem := streetData[0]

				if urlFields.State != obtainedField.StateAbbreviation {

					if urlFields.State == "" {
						statusStrng += fmt.Sprintf("Empty state --> %s,", obtainedField.StateAbbreviation)

					} else {
						statusStrng += fmt.Sprintf("%s --> %s,", urlFields.State, obtainedField.StateAbbreviation)
					}

					clientDatasInput[itemRow].State = obtainedField.StateAbbreviation

				}
				if urlFields.City != obtainedField.CityName {

					if urlFields.City == "" {
						statusStrng += fmt.Sprintf("Empty City --> %s,", obtainedField.CityName)

					} else {
						statusStrng += fmt.Sprintf("%s --> %s,", urlFields.City, obtainedField.CityName)

					}

					clientDatasInput[itemRow].City = obtainedField.CityName

				}
				if urlFields.Zipcode != obtainedField.Zipcode {

					if urlFields.Zipcode == "" {

						statusStrng += fmt.Sprintf("No zip --> %s", obtainedField.Zipcode)
					} else {
						statusStrng += fmt.Sprintf("%s --> %s", urlFields.Zipcode, obtainedField.Zipcode)

					}

					clientDatasInput[itemRow].Zipcode = obtainedField.Zipcode

				}

				if statusStrng == "" {

					statusStrng = "NO UPDATES"
				}
				clientDatasInput[itemRow].StatusField = statusStrng
				clientDatasInput[itemRow].ZipPlus4 = fmt.Sprintf("%s-%s", obtainedField.Zipcode, obtainedField.Plus4Code)
				clientDatasInput[itemRow].Latitude = fmt.Sprintf("%f", obtainedFullItem.Metadata.Latitude)
				clientDatasInput[itemRow].Longitude = fmt.Sprintf("%f", obtainedFullItem.Metadata.Longitude)

				CombinedserviceAddress := fmt.Sprintf("%s,%s,%s", obtainedFullItem.DeliveryLine1, obtainedFullItem.LastLine, obtainedFullItem.DeliveryPointBarcode)

				clientDatasInput[itemRow].ServiceAddress = CombinedserviceAddress
			}

		}

	}

	return clientDatasInput, nil

}
