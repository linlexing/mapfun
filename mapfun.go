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

//ValueEqu 判断值相等
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
		}

	case string:
		if tv2, ok := v2.(string); ok {
			return tv1 == tv2
		}

	case []byte:
		if tv2, ok := v2.([]byte); ok {
			return bytes.Equal(tv1, tv2)
		}

	case int64:
		if tv2, ok := v2.(int64); ok {
			return tv1 == tv2
		}
	case uint64:
		if tv2, ok := v2.(uint64); ok {
			return tv1 == tv2
		}

	case int32:
		if tv2, ok := v2.(int32); ok {
			return tv1 == tv2
		}
	case uint32:
		if tv2, ok := v2.(uint32); ok {
			return tv1 == tv2
		}
	case float64:
		if tv2, ok := v2.(float64); ok {
			return tv1 == tv2
		}
	case float32:
		if tv2, ok := v2.(float32); ok {
			return tv1 == tv2
		}
	default:
		return reflect.DeepEqual(v1, v2)
	}
	return false
}

//SubOf v2是否包含在v1中
func SubOf(v1, v2 map[string]interface{}) bool {
	if v1 == nil && v2 != nil {
		return false
	}
	for k, v := range v2 {
		if sv, ok := v1[k]; ok && !ValueEqu(v, sv) || !ok && v != nil {
			return false
		}
	}
	return true
}

//ValuesByKeys 根据指定的key顺序，返回值的数组
func ValuesByKeys(val map[string]interface{}, keys ...string) []interface{} {
	rev := []interface{}{}
	for _, v := range keys {
		rev = append(rev, val[v])
	}
	return rev
}

//Trans 翻译字段名，如果maps为nil，则原样返回，否则仅返回有对照的属性
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

//Pack 去除value是nil的属性
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

//Without 从一个数组中删除指定的值
func Without(val []interface{}, byRemove ...interface{}) []interface{} {
	if len(byRemove) == 0 {
		return val
	}
	rev := []interface{}{}
	for _, one := range val {
		found := false
		for _, findVal := range byRemove {
			if reflect.DeepEqual(one, findVal) {
				found = true
				break
			}
		}
		if !found {
			rev = append(rev, one)
		}
	}
	return rev
}

//WithoutStr 从一个字符串数组中删除指定的字符串
func WithoutStr(val []string, byRemove ...string) []string {
	if len(byRemove) == 0 {
		return val
	}
	rev := []string{}
	for _, one := range val {
		found := false
		for _, findVal := range byRemove {
			if one == findVal {
				found = true
				break
			}
		}
		if !found {
			rev = append(rev, one)
		}
	}
	return rev
}

//MakeJSON 转换int64和time，使其成为string，以便于用json序列化
func MakeJSON(val map[string]interface{}) {
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

//FromBytes 反序列化map
func FromBytes(data []byte) (rev map[string]interface{}, err error) {
	in := bytes.NewBuffer(data)
	rev = map[string]interface{}{}
	err = gob.NewDecoder(in).Decode(&rev)
	return
}

//Bytes 序列化map
func Bytes(v map[string]interface{}) (rev []byte, err error) {
	buf := bytes.NewBuffer(nil)
	err = gob.NewEncoder(buf).Encode(v)
	if err != nil {
		return
	}
	rev = buf.Bytes()
	return
}

//Changes 返回差异部分
func Changes(v1, v2 map[string]interface{}) (post map[string]interface{}) {
	if v1 == nil || v2 == nil {
		return nil
	}
	post = Clone(v2)
	for k, v := range v2 {
		if sv, ok := v1[k]; ok && ValueEqu(sv, v) || !ok && v == nil {
			delete(post, k)
		}
	}
	return
}

//Object 根据名称和值列表，组装map
func Object(keys []string, values []interface{}) map[string]interface{} {
	r := map[string]interface{}{}
	for i, v := range keys {
		r[v] = values[i]
	}
	return r
}

//Has 判断指定的key是不是在map中存在
func Has(row map[string]interface{}, keys ...string) bool {
	for _, k := range keys {
		if _, ok := row[k]; !ok {
			return false
		}
	}
	return true
}

//Group 分类汇总一个数组，重复值只返回一个
func Group(vals []interface{}) (result []interface{}) {
	if vals == nil {
		return
	}
	result = []interface{}{}
	for _, v := range vals {
		bfound := false
		for _, rv := range result {
			if reflect.DeepEqual(v, rv) {
				bfound = true
				break
			}
		}
		if !bfound {
			result = append(result, v)
		}
	}
	return
}

//GroupStr 排重一个字符串数组
func GroupStr(vals []string) (result []string) {
	if vals == nil {
		return
	}
	result = []string{}
	for _, v := range vals {
		bfound := false
		for _, rv := range result {
			if v == rv {
				bfound = true
				break
			}
		}
		if !bfound {
			result = append(result, v)
		}
	}
	return
}

//Pick 属性白名单
func Pick(row map[string]interface{}, keys ...string) map[string]interface{} {
	r := map[string]interface{}{}
	for _, k := range keys {
		if _, ok := row[k]; ok {
			r[k] = row[k]
		}
	}
	return r
}

//Values 返回所有的值
func Values(row map[string]interface{}) []interface{} {
	r := []interface{}{}
	for _, v := range row {
		r = append(r, v)
	}
	return r
}

//Keys 返回所有的key
func Keys(data map[string]interface{}) []string {
	if data == nil || len(data) == 0 {
		return nil
	}
	result := []string{}
	for k := range data {
		result = append(result, k)
	}
	return result
}

//UpperKeys 将所有的Key转换成大写
func UpperKeys(data map[string]interface{}) map[string]interface{} {
	r := map[string]interface{}{}
	for k, v := range data {
		r[strings.ToUpper(k)] = v
	}
	return r
}

//PickType 计算出指定key的value的数据类型，返回 STR、INT、FLOAT、DATE之一
//警告：要取消
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

//Extend 扩展一个map
func Extend(row1, row2 map[string]interface{}) map[string]interface{} {
	for k, v := range row2 {
		row1[k] = v
	}
	return row1
}

//Intersection 交集，用指定的keys来判断是否相等，返回两个集合交集部分
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

//Difference 找出src在dest中不存在的记录,src、dest可以为nil，注意不是差集，是减集
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

//Omit 黑名单
func Omit(data map[string]interface{}, keys ...string) map[string]interface{} {
	if len(keys) == 0 {
		return data
	}
	rev := map[string]interface{}{}
	keymaps := map[string]bool{}
	for _, k := range keys {
		keymaps[k] = true
	}
	for k, v := range data {
		if _, ok := keymaps[k]; !ok {
			rev[k] = v
		}
	}
	return rev
}

//Pluck 从数组map中找出某个属性的值
func Pluck(data []map[string]interface{}, name string) []interface{} {
	result := []interface{}{}
	for _, one := range data {
		if v, ok := one[name]; ok {
			result = append(result, v)
		}
	}
	return result
}

//Clone 复制一个map
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

//FindWhere 根据一个where去搜索
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

//String 可视化一个map
func String(v map[string]interface{}) string {
	bys, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(bys)
}
