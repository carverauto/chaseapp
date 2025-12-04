package transform

import (
	"reflect"
	"strings"
	"time"
)

const (
	tagName      = "firestore"
	tagOmitEmpty = "omitempty"
	delimiter    = ","
)

type FirestoreMap map[string]interface{}

func ToFirestoreMap(value interface{}) FirestoreMap {
	var result = parseData(value)
	return result.(FirestoreMap)
}

func isZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func parseData(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	var firestoreMap = FirestoreMap{}
	var tag string
	//var value interface{}
	var fieldCount int
	var val = reflect.ValueOf(value)

	switch value.(type) {
	case time.Time, *time.Time:
		return value
	}

	switch val.Kind() {
	case reflect.Map:
		for _, key := range val.MapKeys() {
			val := val.MapIndex(key)
			firestoreMap[key.String()] = parseData(val.Interface())
		}
		return firestoreMap
	case reflect.Ptr:
		if val.IsNil() {
			return nil
		}
		fieldCount = val.Elem().NumField()
		for i := 0; i < fieldCount; i++ {
			tag = val.Elem().Type().Field(i).Tag.Get(tagName)
			value = val.Elem().Field(i).Interface()
			setValue(firestoreMap, tag, value)
		}
		return firestoreMap
	case reflect.Struct, reflect.Interface:
		fieldCount = val.NumField()
		for i := 0; i < fieldCount; i++ {
			tag = val.Type().Field(i).Tag.Get(tagName)
			value = val.Field(i).Interface()
			setValue(firestoreMap, tag, value)
		}
		return firestoreMap
	}
	return value
}

func setValue(firestoreMap FirestoreMap, tag string, value interface{}) {
	if tag == "" || tag == "-" || value == nil {
		return
	}

	tagValues := strings.Split(tag, delimiter)
	if strings.Contains(tag, tagOmitEmpty) {
		if isZeroOfUnderlyingType(value) {
			return
		}
	}
	firestoreMap[tagValues[0]] = parseData(value)
}