package hw09_struct_validator //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrLen      = errors.New("length is invalid")
	ErrMatched  = errors.New("not matched with regexp")
	ErrIncluded = errors.New("not included in the set")
	ErrMin      = errors.New("incoming value is low than a minimum limit")
	ErrMax      = errors.New("incoming value is more than a maximum limit")
	ErrDefault  = errors.New("unsupported type of validator")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errorString string
	for _, err := range v {
		errorString += fmt.Sprintf("Field: %s Error: %s", err.Field, err.Err)
	}

	return errorString
}

func Validate(v interface{}) error {
	var listOfErrors ValidationErrors
	if v == nil {
		return errors.New("empty object")
	}
	if reflect.ValueOf(v).Kind() != reflect.Struct {
		return errors.New("object is not a Struct")
	}

	st := reflect.ValueOf(v).Type()
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		t := f.Tag.Get("validate")
		// nolint:nestif
		if t != "" {
			field := reflect.ValueOf(v).Type().Field(i).Name
			value := reflect.ValueOf(v).Field(i).Interface()
			tag := t
			if reflect.TypeOf(value).Kind() != reflect.Slice {
				resultOfValidation, ok := validateElement(field, value, tag)
				if ok {
					listOfErrors = append(listOfErrors, resultOfValidation)
				}
			}
			if reflect.TypeOf(value).Kind() == reflect.Slice {
				s := reflect.ValueOf(value)
				for i := 0; i < s.Len(); i++ {
					ss := s.Index(i).Interface()
					resultOfValidation, ok := validateElement(field, ss, tag)
					if ok {
						listOfErrors = append(listOfErrors, resultOfValidation)
					}
				}
			}
		}
	}
	if len(listOfErrors) > 0 {
		err := fmt.Errorf(listOfErrors.Error())

		return err
	}

	return nil
}

// nolint:funlen
func validateElement(f string, v interface{}, t string) (ValidationError, bool) {
	var resultOfValidation ValidationError
	listOfValidators := strings.Split(t, "|")
	for _, j := range listOfValidators {
		st := getTag(j)
		switch st {
		case "len":
			var i int
			fmt.Sscanf(j, "len:%5d", &i)
			str := fmt.Sprintf("%v", v)
			if len(str) != i {
				resultOfValidation.Field = f
				resultOfValidation.Err = ErrLen
			}
		case "regexp":
			r := strings.Split(j, ":")[1]
			re := regexp.MustCompile(r)
			if !re.MatchString(v.(string)) {
				resultOfValidation.Field = f
				resultOfValidation.Err = ErrMatched
			}
		case "in":
			fs := strings.Split(j, ":")[1]
			ss := strings.Split(fs, ",")
			t := reflect.TypeOf(v).String()
			switch t {
			case "int":
				vt := v.(int)
				vs := strconv.Itoa(vt)
				if !stringInSlice(vs, ss) {
					resultOfValidation.Field = f
					resultOfValidation.Err = ErrIncluded
				}
			default:
				vs := fmt.Sprintf("%s", v)
				if !stringInSlice(vs, ss) {
					resultOfValidation.Field = f
					resultOfValidation.Err = ErrIncluded
				}
			}
		case "min":
			var i int
			fmt.Sscanf(j, "min:%5d", &i)
			if v.(int) < i {
				resultOfValidation.Field = f
				resultOfValidation.Err = ErrMin
			}
		case "max":
			var i int
			fmt.Sscanf(j, "max:%5d", &i)
			if v.(int) > i {
				resultOfValidation.Field = f
				resultOfValidation.Err = ErrMax
			}
		default:
			resultOfValidation.Field = f
			resultOfValidation.Err = ErrDefault
		}
	}
	if resultOfValidation.Err != nil {
		return resultOfValidation, true
	}

	return resultOfValidation, false
}

func getTag(s string) string {
	var v []string
	v = strings.Split(s, ":")

	return v[0]
}

func stringInSlice(a interface{}, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}
