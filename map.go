package ghost

type Map map[string]interface{}

func NewEmptyMap() Map{
	return Map{}
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
		return v.(int)
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
		return v.(float64)
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

func (m Map) GetMap(key string) Map{
	v := m.Get(key)
	if v != nil {
		if vt, ok := v.(Map); ok{
			return vt
		}
	}
	return Map{}
}

func (m Map) GetArray(key string) []interface{}{
	df := make([]interface{}, 0)

	data := m.Get(key)
	if data == nil{
		return df
	}

	if v, ok := data.([]interface{}); ok{
		return v
	}else{
		return df
	}
}