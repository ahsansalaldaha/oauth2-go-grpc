package utils

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

// PrintStructFields - prints the struct
func PrintStructFields(s interface{}) {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		key := t.Field(i).Name
		value := field.Interface()

		logrus.Printf("%s: %v\n", key, value)
	}
}
