package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func SetStructDefaults(data interface{}) (err error) {
	dataValue := reflect.ValueOf(data).Elem()
	typeOfValue := dataValue.Type()

	for i := 0; i < dataValue.NumField(); i++ {
		fieldValue := dataValue.Field(i)
		fieldType := typeOfValue.Field(i)

		switch fieldType.Type.Kind() {
		case reflect.Struct:
			err = setStructDefaults(fieldValue, fieldType)
			if err != nil {
				return err
			}
		default:
			err = setBaseTypeDefaults(fieldValue, fieldType)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func setStructDefaults(v reflect.Value, fieldType reflect.StructField) error {
	for i := 0; i < v.NumField(); i++ {
		structValue := v.Field(i)
		structType := fieldType.Type.Field(i)

		if err := setBaseTypeDefaults(structValue, structType); err != nil {
			return err
		}
	}
	return nil
}

func setBaseTypeDefaults(v reflect.Value, fieldType reflect.StructField) (err error) {
	defaultValue := fieldType.Tag.Get("default")
	if defaultValue == "" {
		return nil
	}
	switch v.Kind() {
	case reflect.Slice:
		if v.Len() == 0 {
			return setDefaultValue(v, defaultValue)
		}
	default:
		if v.Interface() == reflect.Zero(v.Type()).Interface() {
			return setDefaultValue(v, defaultValue)
		}
	}
	return nil
}

func setDefaultValue(field reflect.Value, defaultValue string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(defaultValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(defaultValue, 10, 64)
		if err != nil {
			return fmt.Errorf("%s not a number", defaultValue)
		}
		field.SetInt(intValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(defaultValue)
		if err != nil {
			return fmt.Errorf("%s not a bool", defaultValue)
		}
		field.SetBool(boolValue)
	case reflect.Slice, reflect.Array:
		elementType := field.Type().Elem()
		if !elementType.ConvertibleTo(reflect.TypeOf("")) {
			return fmt.Errorf("unsupported slice element type: %v", elementType)
		}
		var stringSlice []string
		err := json.Unmarshal([]byte(defaultValue), &stringSlice)
		if err != nil {
			return fmt.Errorf("format error: %v", err)
		}
		sliceValue := reflect.MakeSlice(field.Type(), len(stringSlice), len(stringSlice))
		for i, v := range stringSlice {
			convertedValue := reflect.ValueOf(v).Convert(elementType)
			sliceValue.Index(i).Set(convertedValue)
		}
		field.Set(sliceValue)
	default:
		return fmt.Errorf("unsupported field type: %v", field.Kind())
	}
	return nil
}
