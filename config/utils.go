package config

import (
	"fmt"
	"github.com/incubator/common/maybe"
)

func CheckInt32GT0(key string, val int32) (err maybe.MaybeError) {
	if val > 0 {
		err.Error(nil)
		return
	}
	err.Error(fmt.Errorf("illegal int32(expecting >0) %s: %d", key, val))
	return
}

func CheckInt64GT0(key string, val int64) (err maybe.MaybeError) {
	if val > 0 {
		err.Error(nil)
		return
	}
	err.Error(fmt.Errorf("illegal int64(expecting >0) %s: %d", key, val))
	return
}

func CheckIntGT0(key string, val int) (err maybe.MaybeError) {
	if val > 0 {
		err.Error(nil)
		return
	}
	err.Error(fmt.Errorf("illegal int(expecting >0) %s: %d", key, val))
	return
}

func CheckInt32GET0(key string, val int32) (err maybe.MaybeError) {
	if val >= 0 {
		err.Error(nil)
		return
	}
	err.Error(fmt.Errorf("illegal int32(expecting >=0) %s: %d", key, val))
	return
}

func CheckInt64GET0(key string, val int64) (err maybe.MaybeError) {
	if val >= 0 {
		err.Error(nil)
		return
	}
	err.Error(fmt.Errorf("illegal int64(expecting >=0) %s: %d", key, val))
	return
}

func CheckIntGET0(key string, val int) (err maybe.MaybeError) {
	if val >= 0 {
		err.Error(nil)
		return
	}
	err.Error(fmt.Errorf("illegal int(expecting >=0) %s: %d", key, val))
	return
}

func CheckStringNotEmpty(key string, val string) (err maybe.MaybeError) {
	if val != "" {
		err.Error(nil)
		return
	}
	err.Error(fmt.Errorf("illegal string(expecting !=\"\") %s: %s", key, val))
	return
}

func GetAttrsMap(attrs interface{}, key string) (ret maybe.MaybeEface) {
	if attrs == nil {
		ret.Error(fmt.Errorf("attrs is nil when gettring attr: %s", key))
		return
	}

	val, ok := attrs.(map[string]interface{})
	if !ok {
		ret.Error(fmt.Errorf("illegal attrs(expecting map[string]interface{} when getting attr: %s", key))
		return
	}
	ret.Value(val)
	return
}

func GetAttrMapEface(attrs interface{}, key string) (ret maybe.MaybeEface) {
	attrsMap := GetAttrsMap(attrs, key).Right().(map[string]interface{})

	attr, ok := attrsMap[key]
	if !ok {
		ret.Error(fmt.Errorf("attr not exists: %s", key))
		return
	}
	val, ok := attr.(map[string]interface{})
	if !ok {
		ret.Error(fmt.Errorf("attr type error(expecting map[string]interface{}) %s: %v", key, attr))
	}

	ret.Value(val)
	return
}

func GetAttrMapEfaceArray(attrs interface{}, key string) (ret maybe.MaybeEface) {
	attrsMap := GetAttrsMap(attrs, key).Right().(map[string]interface{})

	attr, ok := attrsMap[key]
	if !ok {
		ret.Error(fmt.Errorf("attr not exists: %s", key))
		return
	}
	val, ok := attr.([]map[string]interface{})
	if !ok {
		ret.Error(fmt.Errorf("attr type error(expecting []map[string]interface{}) %s: %v", key, attr))
	}

	ret.Value(val)
	return
}

func GetAttrInt32(attrs interface{}, key string, check func(string, int32) maybe.MaybeError) (ret maybe.MaybeInt32) {
	attrsMap := GetAttrsMap(attrs, key).Right().(map[string]interface{})

	attr, ok := attrsMap[key]
	if !ok {
		ret.Error(fmt.Errorf("attr not exists: %s", key))
		return
	}
	val, ok := attr.(int32)
	if !ok {
		ret.Error(fmt.Errorf("attr type error(expecting int32) %s: %v", key, attr))
	}

	if check != nil {
		check(key, val).Test()
	}

	ret.Value(val)
	return
}

func GetAttrInt64(attrs interface{}, key string, check func(string, int64) maybe.MaybeError) (ret maybe.MaybeInt64) {
	attrsMap := GetAttrsMap(attrs, key).Right().(map[string]interface{})

	attr, ok := attrsMap[key]
	if !ok {
		ret.Error(fmt.Errorf("attr not exists: %s", key))
		return
	}
	val, ok := attr.(int64)
	if !ok {
		ret.Error(fmt.Errorf("attr type error(expecting int64) %s: %v", key, attr))
	}

	if check != nil {
		check(key, val).Test()
	}

	ret.Value(val)
	return
}

func GetAttrInt(attrs interface{}, key string, check func(string, int) maybe.MaybeError) (ret maybe.MaybeInt) {
	attrsMap := GetAttrsMap(attrs, key).Right().(map[string]interface{})

	attr, ok := attrsMap[key]
	if !ok {
		ret.Error(fmt.Errorf("attr not exists: %s", key))
		return
	}
	val, ok := attr.(int)
	if !ok {
		ret.Error(fmt.Errorf("attr type error(expecting int) %s: %v", key, attr))
	}

	if check != nil {
		check(key, val).Test()
	}

	ret.Value(val)
	return
}

func GetAttrString(attrs interface{}, key string, check func(string, string) maybe.MaybeError) (ret maybe.MaybeString) {
	attrsMap := GetAttrsMap(attrs, key).Right().(map[string]interface{})

	attr, ok := attrsMap[key]
	if !ok {
		ret.Error(fmt.Errorf("attr not exists: %s", key))
		return
	}
	val, ok := attr.(string)
	if !ok {
		ret.Error(fmt.Errorf("attr type error(expecting string) %s: %v", key, attr))
	}

	if check != nil {
		check(key, val).Test()
	}

	ret.Value(val)
	return
}

func GetAttrBool(attrs interface{}, key string) (ret maybe.MaybeBool) {
	attrsMap := GetAttrsMap(attrs, key).Right().(map[string]interface{})

	attr, ok := attrsMap[key]
	if !ok {
		ret.Error(fmt.Errorf("attr not exists: %s", key))
		return
	}
	val, ok := attr.(bool)
	if !ok {
		ret.Error(fmt.Errorf("attr type error(expecting bool) %s: %v", key, attr))
	}

	ret.Value(val)
	return
}
