package util

import (
	"reflect"
)

func New(v interface{}) interface{} {
	type_ := reflect.TypeOf(v)
	switch type_.Kind() {
	case reflect.Ptr:
		u := (reflect.New(type_.Elem())).Elem().Interface()
		return u
	case reflect.Struct:
		u := (reflect.New(type_)).Elem().Interface()
		return u
	}
	return nil

}

func NewByType(type_ reflect.Type) interface{} {
	u := (reflect.New(type_)).Elem().Interface()
	return u
}
