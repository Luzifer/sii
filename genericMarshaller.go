package sii

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	singleLineValue = `^\s*%s\s?:\s?(.+)$`
	arrayLineValue  = `^\s*%s(\[([0-9]*)\])?\s?:\s?(.+)$`
)

func genericMarshal(in interface{}) ([]byte, error) {
	return nil, errors.New("Not implemented")
}

func genericUnmarshal(in []byte, out interface{}) error {
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return errors.New("Calling parser with non-pointer")
	}

	if reflect.ValueOf(out).Elem().Kind() != reflect.Struct {
		return errors.New("Calling parser with pointer to non-struct")
	}

	st := reflect.ValueOf(out).Elem()
	for i := 0; i < st.NumField(); i++ {
		valField := st.Field(i)
		typeField := st.Type().Field(i)

		attributeName := typeField.Tag.Get("sii")
		if attributeName == "" {
			// Names must be explicitly defined, library does not support name guessing
			continue
		}

		switch typeField.Type {

		case reflect.TypeOf(Ptr{}):
			data := getSingleValue(in, attributeName)
			v := Ptr{}
			if err := v.UnmarshalSII(data); err != nil {
				return errors.Wrapf(err, "Unable to parse Ptr for attribute %q", attributeName)
			}
			valField.Set(reflect.ValueOf(v))
			continue

		case reflect.TypeOf(Placement{}):
			data := getSingleValue(in, attributeName)
			v := Placement{}
			if err := v.UnmarshalSII(data); err != nil {
				return errors.Wrapf(err, "Unable to parse Placement for attribute %q", attributeName)
			}
			valField.Set(reflect.ValueOf(v))
			continue

		}

		switch typeField.Type.Kind() {
		case reflect.Bool:
			v, err := strconv.ParseBool(string(getSingleValue(in, attributeName)))
			if err != nil {
				return errors.Wrapf(err, "Unable to parse boolean for attribute %q", attributeName)
			}
			valField.SetBool(v)

		case reflect.Float32:
			v, err := sii2float(getSingleValue(in, attributeName))
			if err != nil {
				return errors.Wrapf(err, "Unable to parse float for attribute %q", attributeName)
			}
			valField.Set(reflect.ValueOf(v))

		case reflect.Int, reflect.Int64:
			v, err := strconv.ParseInt(string(getSingleValue(in, attributeName)), 10, 64)
			if err != nil {
				return errors.Wrapf(err, "Unable to parse int for attribute %q", attributeName)
			}
			valField.SetInt(v)

		case reflect.String:
			v := strings.Trim(string(getSingleValue(in, attributeName)), `"`)
			valField.SetString(v)

		case reflect.Uint64:
			v, err := strconv.ParseUint(string(getSingleValue(in, attributeName)), 10, 64)
			if err != nil {
				return errors.Wrapf(err, "Unable to parse uint for attribute %q", attributeName)
			}
			valField.SetUint(v)

		case reflect.Slice:
			ba, err := getArrayValues(in, attributeName)
			if err != nil {
				return errors.Wrapf(err, "Unable to fetch array values for attribute %q", attributeName)
			}

			switch typeField.Type.Elem() {

			case reflect.TypeOf(Ptr{}):
				var v []Ptr
				for _, bv := range ba {
					e := Ptr{}
					if err := e.UnmarshalSII(bv); err != nil {
						return errors.Wrapf(err, "Unable to parse Ptr for attribute %q", attributeName)
					}
					v = append(v, e)
				}
				valField.Set(reflect.ValueOf(v))
				continue

			case reflect.TypeOf(Placement{}):
				var v []Placement
				for _, bv := range ba {
					e := Placement{}
					if err := e.UnmarshalSII(bv); err != nil {
						return errors.Wrapf(err, "Unable to parse Ptr for attribute %q", attributeName)
					}
					v = append(v, e)
				}
				valField.Set(reflect.ValueOf(v))
				continue

			}

			switch typeField.Type.Elem().Kind() {
			case reflect.Bool:
				var v []bool
				for _, bv := range ba {
					pbv, err := strconv.ParseBool(string(bv))
					if err != nil {
						return errors.Wrapf(err, "Unable to parse boolean for attribute %q", attributeName)
					}
					v = append(v, pbv)
				}
				valField.Set(reflect.ValueOf(v))

			case reflect.Int:
				var v []int
				for _, bv := range ba {
					pbv, err := strconv.Atoi(string(bv))
					if err != nil {
						return errors.Wrapf(err, "Unable to parse int for attribute %q", attributeName)
					}
					v = append(v, pbv)
				}
				valField.Set(reflect.ValueOf(v))

			case reflect.Int64:
				var v []int64
				for _, bv := range ba {
					pbv, err := strconv.ParseInt(string(bv), 10, 64)
					if err != nil {
						return errors.Wrapf(err, "Unable to parse int for attribute %q", attributeName)
					}
					v = append(v, pbv)
				}
				valField.Set(reflect.ValueOf(v))

			case reflect.String:
				var v []string
				for _, bv := range ba {
					v = append(v, strings.Trim(string(bv), `"`))
				}
				valField.Set(reflect.ValueOf(v))

			default:
				return errors.Errorf("Unsupported type: []%s", typeField.Type.Elem().Kind())
			}

		default:
			return errors.Errorf("Unsupported type: %s", typeField.Type.Kind())
		}
	}

	return nil
}

func getSingleValue(in []byte, name string) []byte {
	rex := regexp.MustCompile(fmt.Sprintf(singleLineValue, name))

	var scanner = bufio.NewScanner(bytes.NewReader(in))
	for scanner.Scan() {
		if rex.Match(scanner.Bytes()) {
			grp := rex.FindSubmatch(scanner.Bytes())
			return grp[1]
		}
	}
	return nil
}

func getArrayValues(in []byte, name string) ([][]byte, error) {
	rex := regexp.MustCompile(fmt.Sprintf(arrayLineValue, name))
	var out [][]byte

	var scanner = bufio.NewScanner(bytes.NewReader(in))
	for scanner.Scan() {
		if rex.Match(scanner.Bytes()) {
			grp := rex.FindSubmatch(scanner.Bytes())
			if len(grp[1]) == 0 {
				arrayLen, err := strconv.Atoi(string(grp[3]))
				if err != nil {
					return nil, errors.Wrap(err, "Unable to parse array capacity")
				}
				out = make([][]byte, arrayLen)
				continue
			}

			if len(grp[2]) == 0 {
				out = append(out, grp[3])
				continue
			}

			idx, err := strconv.Atoi(string(grp[2]))
			if err != nil {
				return nil, errors.Wrap(err, "Unable to parse array index")
			}

			out[idx] = grp[3]
		}
	}

	return out, nil
}
