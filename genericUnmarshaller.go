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

func genericUnmarshal(in []byte, out interface{}, unit *Unit) error {
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
			v := Ptr{unit: unit}
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

		case reflect.TypeOf(RawValue{}):
			data := getSingleValue(in, attributeName)
			v := RawValue{}
			if err := v.UnmarshalSII(data); err != nil {
				return errors.Wrapf(err, "Unable to parse RawValue for attribute %q", attributeName)
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
			bv := getSingleValue(in, attributeName)
			if isNilValue(bv) || len(bv) == 0 {
				continue
			}
			v, err := strconv.ParseInt(string(bv), 10, 64)
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

		case reflect.Array:

			switch typeField.Type.Elem().Kind() {

			case reflect.Float32:
				switch typeField.Type.Len() {

				case 3:
					grps := regexp.MustCompile(`^\(([0-9.-]+|&[0-9a-f]+), ([0-9.-]+|&[0-9a-f]+), ([0-9.-]+|&[0-9a-f]+)\)$`).
						FindSubmatch(getSingleValue(in, attributeName))
					var v [3]float32

					for i := range v {
						val, err := sii2float(grps[i+1][:])
						if err != nil {
							return errors.Wrapf(err, "Unable to parse float32 for attribute %q", attributeName)
						}
						v[i] = val
					}
					valField.Set(reflect.ValueOf(v))

				default:
					return errors.Errorf("Unsupported type: [%d]%s", typeField.Type.Len(), typeField.Type.Elem().Kind())

				}

			case reflect.Int64:
				switch typeField.Type.Len() {

				case 3:
					grps := regexp.MustCompile(`^\(([0-9.-]+), ([0-9.-]+), ([0-9.-]+)\)$`).
						FindSubmatch(getSingleValue(in, attributeName))
					var v [3]int64

					for i := range v {
						val, err := strconv.ParseInt(string(grps[i+1][:]), 10, 64)
						if err != nil {
							return errors.Wrapf(err, "Unable to parse int64 for attribute %q", attributeName)
						}
						v[i] = val
					}
					valField.Set(reflect.ValueOf(v))

				default:
					return errors.Errorf("Unsupported type: [%d]%s", typeField.Type.Len(), typeField.Type.Elem().Kind())

				}

			default:
				return errors.Errorf("Unsupported type: [%d]%s", typeField.Type.Len(), typeField.Type.Elem().Kind())

			}

		case reflect.Ptr:

			switch typeField.Type.Elem().Kind() {

			case reflect.Int64:
				bv := getSingleValue(in, attributeName)
				if !isNilValue(bv) && len(bv) > 0 {
					v, err := strconv.ParseInt(string(bv), 10, 64)
					if err != nil {
						return errors.Wrapf(err, "Unable to parse int for attribute %q", attributeName)
					}
					valField.Set(reflect.ValueOf(&v))
				}

			default:
				return errors.Errorf("Unsupported type: *%s", typeField.Type.Elem().Kind())

			}

		case reflect.Slice:
			ba, err := getArrayValues(in, attributeName)
			if err != nil {
				return errors.Wrapf(err, "Unable to fetch array values for attribute %q", attributeName)
			}

			switch typeField.Type.Elem() {

			case reflect.TypeOf(Ptr{}):
				var v []Ptr
				for _, bv := range ba {
					e := Ptr{unit: unit}
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
						return errors.Wrapf(err, "Unable to parse Placement for attribute %q", attributeName)
					}
					v = append(v, e)
				}
				valField.Set(reflect.ValueOf(v))
				continue

			case reflect.TypeOf(RawValue{}):
				var v []RawValue
				for _, bv := range ba {
					e := RawValue{}
					if err := e.UnmarshalSII(bv); err != nil {
						return errors.Wrapf(err, "Unable to parse RawValue for attribute %q", attributeName)
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

			case reflect.Float32:
				var v []float32
				for _, bv := range ba {
					pbv, err := sii2float(bv)
					if err != nil {
						return errors.Wrapf(err, "Unable to parse float32 for attribute %q", attributeName)
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
					if len(bv) == 0 {
						v = append(v, 0)
						continue
					}
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

			case reflect.Array:

				switch typeField.Type.Elem().Elem().Kind() {

				case reflect.Float32:
					switch typeField.Type.Elem().Len() {

					case 3:
						var v [][3]float32
						for _, bv := range ba {
							grps := regexp.MustCompile(`^\(([0-9.-]+|&[0-9a-f]+), ([0-9.-]+|&[0-9a-f]+), ([0-9.-]+|&[0-9a-f]+)\)$`).
								FindSubmatch(bv)
							var sv [3]float32

							for i := range sv {
								val, err := sii2float(grps[i+1][:])
								if err != nil {
									return errors.Wrapf(err, "Unable to parse float32 for attribute %q", attributeName)
								}
								sv[i] = val
							}
							v = append(v, sv)
						}
						valField.Set(reflect.ValueOf(v))

					case 4:
						var v [][4]float32
						for _, bv := range ba {
							grps := regexp.MustCompile(`^\(([0-9.-]+|&[0-9a-f]+); ([0-9.-]+|&[0-9a-f]+), ([0-9.-]+|&[0-9a-f]+), ([0-9.-]+|&[0-9a-f]+)\)$`).
								FindSubmatch(bv)
							var sv [4]float32

							for i := range sv {
								val, err := sii2float(grps[i+1][:])
								if err != nil {
									return errors.Wrapf(err, "Unable to parse float32 for attribute %q", attributeName)
								}
								sv[i] = val
							}
							v = append(v, sv)
						}
						valField.Set(reflect.ValueOf(v))

					default:
						return errors.Errorf("Unsupported len of type: [][%d]%s", typeField.Type.Elem().Len(), typeField.Type.Elem().Elem().Kind())

					}

				default:
					return errors.Errorf("Unsupported type: [][%d]%s", typeField.Type.Elem().Len(), typeField.Type.Elem().Elem().Kind())

				}

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

			out[idx] = make([]byte, len(grp[3]))
			for i, b := range grp[3] {
				out[idx][i] = b
			}
		}
	}

	return out, errors.Wrap(scanner.Err(), "Unable to parse array lines")
}

func isNilValue(in []byte) bool {
	return reflect.DeepEqual(in, []byte("nil"))
}
