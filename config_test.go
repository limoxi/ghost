package ghost

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseFn(t *testing.T){
	str := `{
		"mode": "${MODE}||prod",
		"engine_status": "${STATUS}||off",
		"started": true,
		"power": 3.14,
		"nicer": "${nicer}",
		"abc": {
			"a": 1,
			"b": "B",
			"c": true
		}
	}`
	var data map[string]interface{}
	err := json.Unmarshal([]byte(str), &data)
	if err != nil{
		assert.Fail(t, err.Error())
	}
	s := parseEnvArgs(data)
	c := new(config)
	c.data = s
	assert.Equal(t, "prod", c.GetString("mode"))
	assert.Equal(t, 3.14, c.GetFloat("power"))
	assert.Equal(t, "a", c.GetString("nicer", "a"))
	assert.Equal(t, "B", c.GetString("abc.b"))
	assert.Equal(t, true, c.GetBool("abc.c"))
}