package mapfun

import (
	"reflect"
	"testing"
)

func TestSubOf(t *testing.T) {
	var v1 map[string]interface{}

	if SubOf(v1, map[string]interface{}{"1": "aa"}) {
		t.Error("error")
	}
	v1 = map[string]interface{}{
		"流水号": "123",
		"ID":  int64(11),
	}
	v2 := map[string]interface{}{
		"流水号": "123",
		"ID":  int64(11),
	}
	if !SubOf(v1, v2) {
		t.Error("error")
	}
	v1 = map[string]interface{}{
		"流水号": "123",
		"ID":  12,
	}
	v2 = map[string]interface{}{
		"流水号": "123",
		"ID":  11,
	}
	if SubOf(v1, v2) {
		t.Error("error")
	}
	v1 = map[string]interface{}{
		"流水号": "123",
		"ID":  nil,
		"T":   nil,
		"B":   1,
	}
	v2 = map[string]interface{}{
		"流水号": "123",
		"ID":  nil,
		"A":   nil,
		"B":   1,
	}
	if !SubOf(v1, v2) {
		t.Error("error")
	}
}
func TestChanges(t *testing.T) {
	v1 := map[string]interface{}{
		"流水号": "123",
		"ID":  int64(11),
	}
	v2 := map[string]interface{}{
		"流水号": "123",
		"ID":  int64(11),
	}
	post := Changes(v1, v2)
	if len(post) != 0 {
		t.Error("not empty")
	}

	v1 = map[string]interface{}{
		"流水号": "123",
		"ID":  12,
	}
	v2 = map[string]interface{}{
		"流水号": "123",
		"ID":  11,
	}
	post = Changes(v1, v2)
	if !reflect.DeepEqual(post, map[string]interface{}{
		"ID": 11,
	}) {
		t.Error("not equ", post)
	}
	v1 = map[string]interface{}{
		"流水号": "123",
		"ID":  nil,
		"T":   nil,
		"B":   1,
	}
	v2 = map[string]interface{}{
		"流水号": "123",
		"ID":  nil,
		"A":   nil,
		"B":   nil,
	}
	post = Changes(v1, v2)
	if !reflect.DeepEqual(post, map[string]interface{}{
		"B": nil,
	}) {
		t.Error("not empty")
	}
}
