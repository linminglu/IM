package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

var (
	JSON_TYPE_ERROR = errors.New("type error")
	JSON_NOT_FOUND  = errors.New("not found")
)

type Json struct {
	//	body []byte
	object map[string]interface{}
}

func Marshal(data *Json) ([]byte, error) {
	return json.Marshal(data.object)
}

func New() *Json {
	return &Json{
		object: make(map[string]interface{}),
	}
}

func NewWithMap(obj map[string]interface{}) *Json {
	return &Json{
		object: obj,
	}
}

func (this *Json) Map() map[string]interface{} {
	return this.object
}

func Unmarshal(data []byte) (*Json, error) {
	jsonObject := &Json{
		object: make(map[string]interface{}),
	}

	if err := json.Unmarshal(data, &jsonObject.object); err != nil {
		return nil, err
	}
	return jsonObject, nil
}

func (this *Json) Set(key string, value interface{}) *Json {
	this.object[key] = value
	return this
}

func (this *Json) Has(key string) bool {
	_, ok := this.object[key]
	return ok
}
func (this *Json) Delete(key string) *Json {
	delete(this.object, key)
	return this
}

func (this *Json) ToString() string {
	return fmt.Sprint(this.object)
}

func (this *Json) Bool(key string) (bool, error) {
	value, ok := this.object[key]
	if !ok {
		return false, JSON_NOT_FOUND
	}
	switch v := value.(type) {
	case bool:
		return v, nil
	default:
		return v != nil, nil
	}
	return false, nil
}

func (this *Json) DefaultBool(key string, defaultValue bool) bool {
	value, err := this.Bool(key)
	if err != nil {
		return defaultValue
	}
	return value
}

func (this *Json) Int(key string) (int64, error) {
	value, ok := this.object[key]
	if !ok {
		return 0, JSON_NOT_FOUND
	}
	switch v := value.(type) {
	case string:
		return strconv.ParseInt(v, 10, 64)
	case int8:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case int64:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case int:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case float64:
		return int64(v), nil
	default:
		return int64(0), JSON_TYPE_ERROR
	}
	return 0, nil
}

func (this *Json) DefaultInt(key string, defaultValue int64) int64 {
	value, err := this.Int(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func (this *Json) Float(key string) (float64, error) {
	value, ok := this.object[key]
	if !ok {
		return float64(0), JSON_NOT_FOUND
	}

	switch v := value.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case float64:
		return v, nil
	default:
		return float64(0), JSON_TYPE_ERROR
	}
	return float64(0), nil
}

func (this *Json) DefaultFloat(key string, defaultValue float64) float64 {
	value, err := this.Float(key)
	if err != nil {
		return defaultValue
	}
	return value
}
func (this *Json) String(key string) (string, error) {
	value, ok := this.object[key]
	if !ok {
		return "", JSON_NOT_FOUND
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case float64:
		return fmt.Sprint(v), nil
	default:
		return "", JSON_TYPE_ERROR
	}
	return "", nil

}

func (this *Json) DefaultString(key string, defaultValue string) string {
	value, err := this.String(key)
	if err != nil {
		return defaultValue
	}
	return value
}

func (this *Json) Object(key string) (*Json, error) {
	value, ok := this.object[key]
	if !ok {
		return nil, JSON_NOT_FOUND
	}

	obj, ok := value.(map[string]interface{})
	if !ok {
		return nil, JSON_TYPE_ERROR
	}
	return &Json{object: obj}, nil
}

/*
func main() {
	str := `{"username":"xiaoyao", "userId":"10000", "object":{"subObject":"test"}}`
	jsonObj, err := Unmarshal([]byte(str))
	if err != nil {
		return
	}

	log.Println(obj.String("subObject"))
}
*/
