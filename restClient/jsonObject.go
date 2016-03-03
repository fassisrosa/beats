// author: Francisco Assis Rosa
// date: 2016-02-23
package restClient

type JsonObject map[string]interface{}

func NewJsonObject(data []byte) (JsonObject, error) {
	jsonData, err := decodeJson(data)
	if err != nil {
		return nil, err
	}
	return jsonData.(JsonObject), nil
}

func (json JsonObject) HasKey(key string) bool {
	return json[key] != nil
}

func (json JsonObject) GetObject(key string) JsonObject {
	return json[key].(JsonObject)
}

func (json JsonObject) GetArray(key string) JsonArray {
	return json[key].(JsonArray)
}

func (json JsonObject) GetString(key string) string {
	return json[key].(string)
}

func (json JsonObject) GetInteger(key string) int64 {
	return json[key].(int64)
}

func (json JsonObject) GetBoolean(key string) bool {
	return json[key].(bool)
}
