// author: Francisco Assis Rosa
// date: 2016-02-23
package restClient

type JsonArray []interface{}

func NewJsonArray(data []byte) (JsonArray, error) {
	jsonData, err := decodeJson(data)
	if err != nil {
		return nil, err
	}
	return jsonData.(JsonArray), nil
}

func (json JsonArray) GetObject(key int) JsonObject {
	return json[key].(JsonObject)
}

func (json JsonArray) GetArray(key int) JsonArray {
	return json[key].(JsonArray)
}

func (json JsonArray) GetString(key int) string {
	return json[key].(string)
}

func (json JsonArray) GetInteger(key int) int64 {
	return json[key].(int64)
}

func (json JsonArray) GetBoolean(key int) bool {
	return json[key].(bool)
}
