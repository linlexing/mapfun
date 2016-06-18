package mapfun

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func init() {
	gob.Register(map[string]interface{}{})
}
func ValueEqu(v1, v2 interface{}) bool {

	if v1 == nil && v2 == nil {
		return true
	}
	if v1 == nil && v2 != nil ||
		v1 != nil && v2 == nil {
		return false
	}
	switch tv1 := v1.(type) {
	case time.Time:
		if tv2, ok := v2.(time.Time); ok {
			return tv1.Equal(tv2)
		} else {
			return false
		}
	case string:
		if tv2, ok := v2.(string); ok {
			return tv1 == tv2
		} else {
			return false
		}
	case []byte:
		if tv2, ok := v2.([]byte); ok {
			return bytes.Equal(tv1, tv2)
		} else {
			return false
		}
	case int64:
		if tv2, ok := v2.(int64); ok {
			return tv1 == tv2
		} else {
			return false
		}
	case uint64:
		if tv2, ok := v2.(uint64); ok {
			return tv1 == tv2
		} else {
			return false
		}
	case int32:
		if tv2, ok := v2.(int32); ok {
			return tv1 == tv2
		} else {
			return false
		}
	case uint32:
		if tv2, ok := v2.(uint32); ok {
			return tv1 == tv2
		} else {
			return false
		}
	case float64:
		if tv2, ok := v2.(float64); ok {
			return tv1 == tv2
		} else {
			return false
		}
	case float32:
		if tv2, ok := v2.(float32); ok {
			return tv1 == tv2
		} else {
			return false
		}
	default:
		return reflect.DeepEqual(v1, v2)
	}
}

//v2是否包含在v1中
func SubOf(v1, v2 map[string]interface{}) bool {
	for k, v := range v2 {
		if sv, ok := v1[k]; !ok || !ValueEqu(v, sv) {
			return false
		}
	}
	return true
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
	if len(v) == 0 || len(maps) == 0 {
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
func Pack(val map[string]interface{}) map[string]interface{} {
	list := []string{}
	for k, v := range val {
		if v == nil {
			list = append(list, k)
		}
	}
	for _, k := range list {
		delete(val, k)
	}
	return val
}

//转换int64和time，使其成为string，以便于用json序列化
func MakeJson(val map[string]interface{}) {
	for k, v := range val {
		switch tv := v.(type) {
		case int64, uint64:
			val[k] = fmt.Sprintf("%d", tv)
		case time.Time:
			val[k] = tv.Format(time.RFC3339)
		}
	}
	return
}
func FromBytes(data []byte) (rev map[string]interface{}, err error) {
	in := bytes.NewBuffer(data)
	rev = map[string]interface{}{}
	err = gob.NewDecoder(in).Decode(&rev)
	return
}

func Bytes(v map[string]interface{}) (rev []byte, err error) {
	buf := bytes.NewBuffer(nil)
	err = gob.NewEncoder(buf).Encode(v)
	if err != nil {
		return
	}
	rev = buf.Bytes()
	return
}

//返回差异部分
func Changes(v1, v2 map[string]interface{}) (pre, post map[string]interface{}) {
	if v1 == nil || v2 == nil {
		return v1, v2
	}
	//一定要Pack，否则nil值会出问题
	pre = Pack(Clone(v1))
	post = Pack(Clone(v2))
	removeList := []string{}
	//删除pre、post相同的值
	for k, v := range pre {
		if sv, ok := post[k]; ok && ValueEqu(sv, v) {
			removeList = append(removeList, k)
		}
	}
	for _, str := range removeList {
		delete(pre, str)
		delete(post, str)
	}
	return
}
func Object(keys []string, values []interface{}) map[string]interface{} {
	r := map[string]interface{}{}
	for i, v := range keys {
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
func Extend(row1, row2 map[string]interface{}) map[string]interface{} {
	for k, v := range row2 {
		row1[k] = v
	}
	return row1
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
