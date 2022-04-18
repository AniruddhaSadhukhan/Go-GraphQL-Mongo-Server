package common

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/microcosm-cc/bluemonday"
)

const (
	MAX_PERMISSIBLE_INPUT_STRING_LENGTH = 10000
	MAX_PERMISSIBLE_INPUT_NUMBER        = 5000000
)

var initializeHTMLSanitizer sync.Once
var htmlSanitizer *bluemonday.Policy

func getHtmlSanitizer() *bluemonday.Policy {
	initializeHTMLSanitizer.Do(func() {
		if htmlSanitizer == nil {
			htmlSanitizer = bluemonday.UGCPolicy()
		}
	})
	return htmlSanitizer
}

func Sanitize(rawInput interface{}) (interface{}, error) {

	input := reflect.ValueOf(rawInput)

	switch input.Kind() {

	// Recursively sanitize all values of the map
	case reflect.Map:
		for _, key := range input.MapKeys() {
			newValue, err := Sanitize(input.MapIndex(key).Interface())
			if err != nil {
				return rawInput, err
			}
			input.SetMapIndex(key, reflect.ValueOf(newValue))
		}

	// Recursively sanitize all values of the slice
	case reflect.Slice:
		for i := 0; i < input.Len(); i++ {
			newValue, err := Sanitize(input.Index(i).Interface())
			if err != nil {
				return rawInput, err
			}
			input.Index(i).Set(reflect.ValueOf(newValue))
		}

	//Check for string length and sanitize the string using Bluemonday
	case reflect.String:
		rawInputString := rawInput.(string)
		if len(rawInputString) > MAX_PERMISSIBLE_INPUT_STRING_LENGTH {
			return rawInput, fmt.Errorf("string length exceeded")
		}
		rawInput = getHtmlSanitizer().Sanitize(rawInputString)

	// For int and float, if number exceeds 5 million, return error
	case reflect.Int:
		if rawInput.(int) > MAX_PERMISSIBLE_INPUT_NUMBER {
			return rawInput, fmt.Errorf("number exceeded")
		}

	case reflect.Float64:
		if rawInput.(float64) > MAX_PERMISSIBLE_INPUT_NUMBER {
			return rawInput, fmt.Errorf("floating point number exceeded")
		}

	}

	return rawInput, nil

}
