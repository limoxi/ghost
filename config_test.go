package ghost

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefault(t *testing.T){
	c := new(config)
	t.Log(Map(nil))
	assert.Equal(t, Map(nil), c.Map)
}

func TestParseFn(t *testing.T){
	str := `{
		"mode": "${MODE}||prod",
		"engine_status": "${STATUS}||off",
		"started": true,
		"power": 3.14,
		"nicer": "${nicer}||a",
		"timeout": "${timeout}||10.1",
		"abc": {
			"a": 1,
			"b": "B",
			"c": true
		},
		"array": [1,2,3]
	}`
	var data Map
	err := json.Unmarshal([]byte(str), &data)
	if err != nil{
		assert.Fail(t, err.Error())
	}
	s := parseEnvArgs(data)
	c := new(config)
	c.Map = s
	assert.Equal(t, "prod", c.GetString("mode"))
	assert.Equal(t, 3.14, c.GetFloat("power"))
	assert.Equal(t, "a", c.GetString("nicer", "a"))
	assert.Equal(t, "B", c.GetString("abc.b"))
	assert.Equal(t, true, c.GetBool("abc.c"))
	assert.Equal(t, 1.0, c.GetFloat("abc.a"))
	assert.Equal(t, 10.1, c.GetFloat("timeout"))
	assert.Equal(t, Map{
		"a": 1.0,
		"b": "B",
		"c": true,
	}, c.GetMap("abc"))
	assert.Equal(t, []interface{}{1.0,2.0,3.0}, c.GetArray("array"))
}

func TestA(t *testing.T){
	nums := []int{4,1,2,2,1}
	for i, _ := range nums{
		for j:=i+1;j<len(nums);j++{
			if nums[i]>nums[j]{
				nums[i], nums[j] = nums[j], nums[i]
			}
		}

	}
	t.Log(nums)
	a := 0

out:	for i, num := range nums{
		if i+1 >=len(nums){
			a = num
			break
		}
		if i&1==1{
			continue
		}
		for j:=i+1; j<len(nums);j++{
			if num == nums[j]{
				break
			}
			if num != nums[j]{
				a = num
				break out

			}
		}
	}
	t.Log(a)
}