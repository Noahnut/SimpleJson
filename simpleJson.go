package simplejson

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type structTypeField struct {
	fieldString []byte
	fieldIndex  int
	fieldType   reflect.Type
}

func Marshal(v interface{}) ([]byte, error) {
	marshal := make([]byte, 0)
	valueOf := reflect.ValueOf(v)
	vtype := valueOf.Type()

	switch vtype.Kind() {
	case reflect.Struct:
		marshal = append(marshal, structEncodeState(valueOf)...)
	}

	return marshal, nil
}

func structEncodeState(v reflect.Value) []byte {
	structType := v.Type()
	structfields := structField(structType)

	jsonString := "{"
	next := ""
	for _, structfield := range structfields {
		jsonString += next
		value := v.Field(structfield.fieldIndex)
		fieldString := string(structfield.fieldString)
		jsonfield := getJsonValue(fieldString, value, structfield.fieldType)
		jsonString += jsonfield
		next = ","
	}

	jsonString += "}"

	return []byte(jsonString)
}

func getJsonValue(fieldString string, value reflect.Value, valuetype reflect.Type) string {
	switch valuetype.Kind() {
	case reflect.String:
		return fmt.Sprintf(`%s"%s"`, fieldString, value.String())
	case reflect.Int:
		return fmt.Sprintf(`%s%d`, fieldString, value.Int())
	case reflect.Struct:
		return fmt.Sprintf(`%s%s`, fieldString, string(structEncodeState(value)))
	case reflect.Slice:
		return fmt.Sprintf(`%s%s`, fieldString, arrayEncodeState(value))
	default:
		return ""
	}
}

func arrayEncodeState(value reflect.Value) string {

	arrayString := "["
	slice := make([]interface{}, value.Len())
	for i := 0; i < value.Len(); i++ {
		slice[i] = value.Index(i).Interface()
	}

	next := ""
	switch slice[0].(type) {
	case int:
		for _, s := range slice {
			arrayString += next
			arrayString += strconv.Itoa(s.(int))
			next = ","
		}
	case string:
		for _, s := range slice {
			arrayString += next
			valueString := `"` + s.(string) + `"`
			arrayString += valueString
			next = ","
		}
	}

	arrayString += "]"

	return arrayString
}

func structField(structType reflect.Type) []structTypeField {
	typeField := make([]structTypeField, 0)
	for i := 0; i < structType.NumField(); i++ {
		fieldName := structType.Field(i).Name
		if len(structType.Field(i).Tag.Get("json")) != 0 {
			parseTagName(structType.Field(i).Tag.Get("json"), &fieldName)
		}

		next := `"` + fieldName + `"` + `:`

		f := structTypeField{
			fieldString: []byte(next),
			fieldIndex:  i,
			fieldType:   structType.Field(i).Type,
		}

		typeField = append(typeField, f)
	}

	return typeField
}

func parseTagName(tagName string, fieldName *string) {
	if !isValidTag(tagName) {
		return
	}
	*fieldName = tagName
	log.Println(tagName)
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:;<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}
	return true
}
