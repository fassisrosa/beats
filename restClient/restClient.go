// author: Francisco Assis Rosa
// date: 2016-02-23
package restClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

type RestClient struct {
	client     *http.Client
	headerInfo map[string]string
}

func NewRestClient() (restClient RestClient, err error) {

	restClient = RestClient{}

	cookieJar, err := cookiejar.New(nil)
	restClient.client = &http.Client{
		Jar: cookieJar,
	}
	if err != nil {
		return
	}
	restClient.headerInfo = map[string]string{"Content-Type": "application/json"}

	return
}

func (restClient RestClient) AddHeader(headerKey string, headerValue string) {
	restClient.headerInfo[headerKey] = headerValue
}

func (restClient RestClient) GetObject(method string, url string, payload []byte) (jsonObject JsonObject, err error) {
	body, err := restClient.doGetString(method, url, payload)
	if err != nil {
		return nil, err
	}
	return NewJsonObject(body)
}

func (restClient RestClient) GetArray(method string, url string, payload []byte) (jsonObject JsonArray, err error) {
	body, err := restClient.doGetString("GET", url, payload)
	if err != nil {
		return nil, err
	}
	return NewJsonArray(body)
}

func (restClient RestClient) doGetString(method string, url string, payload []byte) (body []byte, err error) {
	var req *http.Request
	if payload != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return
	}

	// set all required headers
	for headerKey, headerValue := range restClient.headerInfo {
		req.Header.Add(headerKey, headerValue)
	}
	resp, err := restClient.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

// DecodeJson converts a byte sequence into maps, arrays, native types.
// It does away with needing to specify structures to decode
// the byte sequence by storing JSON in generic structures.
// Users of DecodeJson can then retrieve required information
// knowing the path access to this information.
func decodeJson(data []byte) (interface{}, error) {
	// try converting to string first
	var dataString string = ""
	err := json.Unmarshal(data, &dataString)
	if err == nil {
		return dataString, nil
	}

	// try decoding an int number next
	var dataInt int64 = 0
	err = json.Unmarshal(data, &dataInt)
	if err == nil {
		return dataInt, nil
	}

	// try decoding a float number next
	var dataFloat float64 = 0.0
	err = json.Unmarshal(data, &dataFloat)
	if err == nil {
		return dataFloat, nil
	}

	// try decoding a boolean next
	var dataBool bool = false
	err = json.Unmarshal(data, &dataBool)
	if err == nil {
		return dataBool, nil
	}

	// try decoding a map next
	var dataMap map[string]*json.RawMessage
	err = json.Unmarshal(data, &dataMap)
	if err == nil {
		resultMap := JsonObject{}
		for key, oneRawMessage := range dataMap {
			resultMap[key] = nil
			if oneRawMessage != nil {
				resultMap[key], err = decodeJson(*oneRawMessage)
				if err != nil {
					return nil, err
				}
			}
		}
		return resultMap, nil
	}

	// try decoding an array next
	var dataArray []*json.RawMessage
	err = json.Unmarshal(data, &dataArray)
	if err == nil {
		resultArray := make(JsonArray, len(dataArray), len(dataArray))
		for idx, oneRawMessage := range dataArray {
			arrayEntry, err := decodeJson(*oneRawMessage)
			if err != nil {
				return nil, err
			}
			resultArray[idx] = arrayEntry
		}
		return resultArray, nil
	}

	return nil, fmt.Errorf("Could not convert data, could not find type")
}
