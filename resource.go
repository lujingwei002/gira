package gira

import (
	"fmt"
	"reflect"
)

type ResourceLoader interface {
	// 加载资源
	LoadResource(dir string) error
	// 重载资源
	ReloadResource(dir string) error
	// 配置版本
	GetResVersion() string
	// loader版本
	GetLoaderVersion() string
}

type ResourceHandler interface {
	OnResourcePreLoad()
	OnResourcePostLoad()
}

type ResourceSource interface {
	// 返回资源加载器
	GetResourceLoader() ResourceLoader
}

func Make1Key_int[T any](arr []*T, dict map[int]*T, key string) error {
	for _, v := range arr {
		vv := reflect.ValueOf(*v).FieldByName(key)
		if !vv.IsValid() {
			return fmt.Errorf("Make1Key_int error, field %s not found in %#v", key, v)
		}
		id := (int)(vv.Int())
		dict[id] = v
	}
	return nil
}

func Make1Key_int32[T any](arr []*T, dict map[int32]*T, key string) error {
	for _, v := range arr {
		vv := reflect.ValueOf(*v).FieldByName(key)
		if !vv.IsValid() {
			return fmt.Errorf("Make1Key_int32 error, field %s not found in %#v", key, v)
		}
		id := (int32)(vv.Int())
		dict[id] = v
	}
	return nil
}

func Make1Key_int64[T any](arr []*T, dict map[int64]*T, key string) error {
	for _, v := range arr {
		vv := reflect.ValueOf(*v).FieldByName(key)
		if !vv.IsValid() {
			return fmt.Errorf("Make1Key_int64 error, field %s not found in %#v", key, v)
		}
		id := (int64)(vv.Int())
		dict[id] = v
	}
	return nil
}
func Make1Key_string[T any](arr []*T, dict map[string]*T, key string) error {
	for _, v := range arr {
		vv := reflect.ValueOf(*v).FieldByName(key)
		if !vv.IsValid() {
			return fmt.Errorf("Make1Key_string error, field %s not found in %#v", key, v)
		}
		id := (string)(vv.String())
		dict[id] = v
	}
	return nil
}

func Make2Key_int64_string[T any](arr []*T, dict map[int64]map[string]*T, key1 string, key2 string) error {
	for _, v := range arr {
		v1 := reflect.ValueOf(*v).FieldByName(key1)
		v2 := reflect.ValueOf(*v).FieldByName(key2)
		if !v1.IsValid() {
			return fmt.Errorf("Make2Key_int64_string error, field %s not found in %#v", key1, v)
		}
		if !v2.IsValid() {
			return fmt.Errorf("Make2Key_int64_string error, field %s not found in %#v", key2, v)
		}
		id1 := (int64)(v1.Int())
		id2 := (string)(v2.String())
		if _, ok := dict[id1]; !ok {
			dict[id1] = make(map[string]*T, 0)
		}
		if dict1, ok := dict[id1]; !ok {
			return fmt.Errorf("Make2Key_int64_string error, create map %s fail %#v", key1, v)
		} else {
			dict1[id2] = v
		}
	}
	return nil
}

func Make2Key_int_int[T any](arr []*T, dict map[int]map[int]*T, key1 string, key2 string) error {
	for _, v := range arr {
		v1 := reflect.ValueOf(*v).FieldByName(key1)
		v2 := reflect.ValueOf(*v).FieldByName(key2)
		if !v1.IsValid() {
			return fmt.Errorf("Make2Key_int_int error, field %s not found in %#v", key1, v)
		}
		if !v2.IsValid() {
			return fmt.Errorf("Make2Key_int_int error, field %s not found in %#v", key2, v)
		}
		id1 := (int)(v1.Int())
		id2 := (int)(v2.Int())
		if _, ok := dict[id1]; !ok {
			dict[id1] = make(map[int]*T, 0)
		}
		if dict1, ok := dict[id1]; !ok {
			return fmt.Errorf("Make2Key_int_int error, create map %s fail %#v", key1, v)
		} else {
			dict1[id2] = v
		}
	}
	return nil
}

func Make2Key_int64_int64[T any](arr []*T, dict map[int64]map[int64]*T, key1 string, key2 string) error {
	for _, v := range arr {
		v1 := reflect.ValueOf(*v).FieldByName(key1)
		v2 := reflect.ValueOf(*v).FieldByName(key2)
		if !v1.IsValid() {
			return fmt.Errorf("Make2Key_int64_int64 error, field %s not found in %#v", key1, v)
		}
		if !v2.IsValid() {
			return fmt.Errorf("Make2Key_int64_int64 error, field %s not found in %#v", key2, v)
		}
		id1 := (int64)(v1.Int())
		id2 := (int64)(v2.Int())
		if _, ok := dict[id1]; !ok {
			dict[id1] = make(map[int64]*T, 0)
		}
		if dict1, ok := dict[id1]; !ok {
			return fmt.Errorf("Make2Key_int64_int64 error, create map %s fail %#v", key1, v)
		} else {
			dict1[id2] = v
		}
	}
	return nil
}

func Make2Key_string_int64[T any](arr []*T, dict map[string]map[int64]*T, key1 string, key2 string) error {
	for _, v := range arr {
		v1 := reflect.ValueOf(*v).FieldByName(key1)
		v2 := reflect.ValueOf(*v).FieldByName(key2)
		if !v1.IsValid() {
			return fmt.Errorf("Make2Key_string_int64 error, field %s not found in %#v", key1, v)
		}
		if !v2.IsValid() {
			return fmt.Errorf("Make2Key_string_int64 error, field %s not found in %#v", key2, v)
		}
		id1 := (string)(v1.String())
		id2 := (int64)(v2.Int())
		if _, ok := dict[id1]; !ok {
			dict[id1] = make(map[int64]*T, 0)
		}
		if dict1, ok := dict[id1]; !ok {
			return fmt.Errorf("Make2Key_string_int64 error, create map %s fail %#v", key1, v)
		} else {
			dict1[id2] = v
		}
	}
	return nil
}
