package ghost

import "strconv"

type Map map[string]interface{}

func NewEmptyMap() Map{
	return Map{}
}

func NewMapFromData(data map[string]interface{}) Map{
	return Map(data)
}

func (m Map) Get(key string) interface{}{
	if v, ok := m[key]; ok{
		return v
	}else{
		return nil
	}
}

func (m Map) GetString(key string, args ...string) string{
	defaultVal := ""
	switch len(args) {
	case 1:
		defaultVal = args[0]
	}
	v := m.Get(key)
	if v != nil{
		return v.(string)
	}else{
		return defaultVal
	}
}

func (m Map) GetInt(key string, args ...int) int{
	var defaultVal = 0
	switch len(args) {
	case 1:
		defaultVal = args[0]
	}
	v := m.Get(key)
	if v != nil{
		switch v.(type) {
		case float64:
			return int(v.(float64))
		case string:
			dv, _ := strconv.Atoi(v.(string))
			return dv
		default:
			return v.(int)
		}
	}else{
		return defaultVal
	}
}

func (m Map) GetFloat(key string, args ...float64) float64{
	var defaultVal = 0.00
	switch len(args) {
	case 1:
		defaultVal = args[0]
	}
	v := m.Get(key)
	if v != nil{
		switch v.(type) {
		case string:
			dv, _ :=  strconv.ParseFloat(v.(string), 64)
			return dv
		default:
			return v.(float64)
		}
	}else{
		return defaultVal
	}
}

func (m Map) GetBool(key string, args ...bool) bool{
	var defaultVal = false
	switch len(args) {
	case 1:
		defaultVal = args[0]
	}
	v := m.Get(key)
	if v != nil{
		return v.(bool)
	}else{
		return defaultVal
	}
}