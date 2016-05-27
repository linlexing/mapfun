package mapfun

import (
	"reflect"
	"testing"
)

func TestChanges(t *testing.T) {
	v1 := map[string]interface{}{
		"流水号": "123",
		"ID":  int64(11),
	}
	v2 := map[string]interface{}{
		"流水号": "123",
		"ID":  int64(11),
	}
	pre, post := Changes(v1, v2)
	if len(pre) != 0 || len(post) != 0 {
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
	pre, post = Changes(v1, v2)
	if !reflect.DeepEqual(pre, map[string]interface{}{
		"ID": 12,
	}) || !reflect.DeepEqual(post, map[string]interface{}{
		"ID": 11,
	}) {
		t.Error("not equ", pre, post)
	}
}
