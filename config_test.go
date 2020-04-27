package ghost

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefault(t *testing.T){
	c := new(config)
	t.Log(GMap(nil))
	assert.Equal(t, GMap(nil), c.GMap)
}

func TestParseFn(t *testing.T){
	c := LoadConfigFromFile("config_test.yaml")
	assert.Equal(t, "prod", c.GetString("mode"))
	assert.Equal(t, 3.14, c.GetFloat("power"))
	assert.Equal(t, "a", c.GetString("nicer", "a"))
	assert.Equal(t, "B", c.GetString("abc.b"))
	assert.Equal(t, true, c.GetBool("abc.c"))
	assert.Equal(t, 1.0, c.GetFloat("abc.a"))
	assert.Equal(t, 1, c.GetInt("abc.a"))
	assert.Equal(t, 10.1, c.GetFloat("timeout"))
	assert.Equal(t, "*", c.GetString("spec"))
	assert.Equal(t, GMap{
		"a": 1,
		"b": "B",
		"c": true,
	}, c.GetMap("abc"))
	assert.Equal(t, []interface{}{1,2,3.3}, c.GetArray("array"))
}