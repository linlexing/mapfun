package mapfun

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"
)

func valueEqu(v1, v2 interface{}) bool {

	if v1 == nil && v2 == nil {
		return true
	}
	if v1 == nil && v2 != nil ||
		v1 != nil && v2 == nil {
		return false
	}
	switch tv := v1.(type) {
	case time.Time:
		return tv.Equal(v2.(time.Time))
	default:
		return reflect.DeepEqual(v1, v2)
	}
}

//根据指定的key顺序，返回值的数组
func ValuesByKeys(val map[string]interface{}, keys ...string) []interface{} {
	rev := []interface{}{}
	for _, v := range keys {
		rev = append(rev, val[v])
	}
	return rev
}

//翻译字段名，如果maps为nil，则原样返回，否则仅返回有对照的属性
func Trans(v map[string]interface{}, maps map[string]string) map[string]interface{} {
	if len(maps) == 0 {
		return Clone(v)
	}
	rev := map[string]interface{}{}
	for fieldName, propertyName := range maps {
		if fv, ok := v[fieldName]; ok {
			rev[propertyName] = fv
		}
	}
	return rev
}
func Pack(val map[string]interface{}) {
	list := []string{}
	for k, v := range val {
		if v == nil {
			list = append(list, k)
		}
	}
	for _, k := range list {
		delete(val, k)
	}
}

//返回差异部分
func Changes(v1, v2 map[string]interface{}) (pre, post map[string]interface{}) {
	if v1 == nil || v2 == nil {
		return v1, v2
	}
	pre = Clone(v1)
	post = Clone(v2)
	removeList := []string{}
	//删除pre、post相同的值
	for k, v := range pre {
		if sv, ok := post[k]; ok && valueEqu(sv, v) {
			removeList = append(removeList, k)
		}
	}
	for _, str := range removeList {
		delete(pre, str)
		delete(post, str)
	}
	return
}
func Object(list []string, values []interface{}) map[string]interface{} {
	r := map[string]interface{}{}
	for i, v := range list {
		r[v] = values[i]
	}
	return r
}
func Has(row map[string]interface{}, keys ...string) bool {
	for _, k := range keys {
		if _, ok := row[k]; !ok {
			return false
		}
	}
	return true
}
func Pick(row map[string]interface{}, keys ...string) map[string]interface{} {
	r := map[string]interface{}{}
	for _, k := range keys {
		r[k] = row[k]
	}
	return r
}
func Values(row map[string]interface{}) []interface{} {
	r := []interface{}{}
	for _, v := range row {
		r = append(r, v)
	}
	return r
}
func Keys(data map[string]interface{}) []string {
	if data == nil || len(data) == 0 {
		return nil
	}
	result := []string{}
	for k, _ := range data {
		result = append(result, k)
	}
	return result
}
func UpperKeys(data map[string]interface{}) map[string]interface{} {
	r := map[string]interface{}{}
	for k, v := range data {
		r[strings.ToUpper(k)] = v
	}
	return r
}

//计算出制定关键字的数据类型
func PickType(rows []map[string]interface{}, keys ...string) map[string]string {
	r := map[string]string{}
	for _, k := range keys {
		r[k] = ""
	}
	for _, row := range rows {
		for _, k := range keys {
			if r[k] == "" {
				switch row[k].(type) {
				case string, []byte:
					r[k] = "STR"
				case int32, int64, uint, uint64:
					r[k] = "INT"
				case float32, float64:
					r[k] = "FLOAT"
				case time.Time, *time.Time:
					r[k] = "DATE"
				case nil:
				default:
					panic("not impl")
				}
			}
		}
	}
	return r
}

//交集，用指定的keys来判断是否相等，返回两个集合交集部分
func Intersection(rows1, rows2 []map[string]interface{}, keys []string) (result1, result2 []map[string]interface{}) {
	if len(rows1) == 0 || len(rows2) == 0 {
		return []map[string]interface{}{}, []map[string]interface{}{}
	}
	result1 = []map[string]interface{}{}
	result2 = []map[string]interface{}{}
	for _, row1 := range rows1 {
		keyValue1 := Pick(row1, keys...)
		for _, row2 := range rows2 {
			if reflect.DeepEqual(keyValue1, Pick(row2, keys...)) {
				result1 = append(result1, row1)
				result2 = append(result2, row2)
				break
			}
		}
	}
	return
}

//找出src在dest中不存在的记录
func Difference(src, dest []map[string]interface{}, keys []string) (result []map[string]interface{}) {
	result = []map[string]interface{}{}
	if len(src) == 0 {
		return
	}
	if len(dest) == 0 {
		result = src
		return
	}
	for _, row1 := range src {
		keyValue1 := Pick(row1, keys...)
		found := false
		for _, row2 := range dest {
			if reflect.DeepEqual(keyValue1, Pick(row2, keys...)) {
				found = true
				break
			}
		}
		if !found {
			result = append(result, row1)
		}
	}
	return
}

//找出某个属性的值
func Pluck(data []map[string]interface{}, name string) []interface{} {
	result := []interface{}{}
	for _, one := range data {
		if v, ok := one[name]; ok {
			result = append(result, v)
		}
	}
	return result
}
func Clone(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return nil
	}
	r := map[string]interface{}{}
	for k, v := range src {
		r[k] = v
	}
	return r
}
func FindWhere(src []map[string]interface{}, where map[string]interface{}) map[string]interface{} {
	for _, row := range src {
		bHas := true
		for k, v := range where {
			if fv, ok := row[k]; !ok || fv != v {
				bHas = false
				break
			}
		}
		if bHas {
			return row
		}
	}
	return nil
}
func String(v map[string]interface{}) string {
	bys, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(bys)
}
