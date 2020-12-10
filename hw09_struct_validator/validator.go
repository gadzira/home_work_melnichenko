package hw09_struct_validator //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const digit = `[0-9]+`

// const pattern = `^[A-Za-z0-9_\.]+`
// var alphaNum = regexp.MustCompile(pattern)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errorString string
	for _, err := range v {
		errorString += fmt.Sprintf("Field: %s \t Error: %s\n", err.Field, err.Err)
	}

	return errorString
}

func Validate(v interface{}) error {

	// Check object is exist
	if v == nil {
		return errors.New("Empty object")
	}

	// Check if obj is struct
	if reflect.ValueOf(v).Kind() != reflect.Struct {
		return errors.New("Object is not a Struct")
	}

	st := reflect.ValueOf(v).Type()

	var lineToValidate []string
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		t := f.Tag.Get("validate")

		if t != "" {
			// fmt.Println("Field:", reflect.ValueOf(v).Type().Field(i).Name,
			// 	"Value:", reflect.ValueOf(v).Field(i).Interface(),
			// 	"Tag:", t)
			lineToValidate = append(lineToValidate, t)
			field := reflect.ValueOf(v).Type().Field(i).Name
			value := reflect.ValueOf(v).Field(i).Interface()
			tag := t
			// listOfValidationErrors := validateElement(field, value, tag)
			_ = validateElement(field, value, tag)
			// listOfValidationErrors := validateElement(reflect.ValueOf(v).Type().Field(i).Name, reflect.ValueOf(v).Field(i).Interface(), t)
			// fmt.Println(listOfValidationErrors)
		}
	}
	return nil
}

func validateElement(f string, v interface{}, t string) ValidationErrors {
	// if reflect.ValueOf(v).Kind() == reflect.String {
	// fmt.Println("It is a", reflect.TypeOf(v).String(), reflect.ValueOf(v))
	// }

	// if reflect.ValueOf(v).Kind() == reflect.Int {
	// fmt.Println("It is a", reflect.TypeOf(v).String(), reflect.ValueOf(v))
	// }

	// Split string by |
	listOfValidators := strings.Split(t, "|")
	for _, j := range listOfValidators {
		// fmt.Printf("v:%s\n", v)
		t := getTag(j)
		switch t {
		case "len":
			fmt.Println("Call len validator for:", v)
			fmt.Println("Print tag:", j)
			//нужно получить цифру
			// var d = regexp.MustCompile(digit)
			// length := d.FindString(j)
			// lenInt, _ := strconv.Atoi(length)
			// if len((v.(string)) > length {
			// 	fmt.Println("ff")
			// }

			var i int
			fmt.Sscanf(j, "len:%5d", &i)
			fmt.Println(i)
			str := fmt.Sprintf("%v", v)
			if len(str) > i {
				fmt.Printf("Lenght of %s more than %d", f, i)
			}

		case "regexp":
			fmt.Println("Call regexp validator for:", v)
		case "in":
			fmt.Println("Call in validator for:", v)
		case "min":
			fmt.Println("Call min validator for:", v)
		case "max":
			fmt.Println("Call max validator for:", v)
		default:
			fmt.Println("Unsupported type")
		}
	}
	var rr ValidationErrors
	return rr

}

// I'm not proud of it
func getTag(s string) string {
	var v string
	if strings.Contains(s, "len") {
		v = "len"
	}
	if strings.Contains(s, "regexp") {
		v = "regexp"
	}
	if strings.Contains(s, "in") {
		v = "in"
	}
	if strings.Contains(s, "min") {
		v = "min"
	}
	if strings.Contains(s, "max") {
		v = "max"
	}
	return v
}
