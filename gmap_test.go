package ghost

import (
	"testing"
)

func TestGMap(t *testing.T) {

	data := map[string]interface{}{
		"int":    12,
		"float":  12.56,
		"string": "abc",
		"list":   []int{1, 2, 3},
		"dict": map[string]string{
			"a": "aaa",
			"b": "bbb",
			"c": "ccc",
		},
	}
	m := NewGMapFromData(data)

	t.Log(m.Get("list").([]int))
}
