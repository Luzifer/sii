package sii

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

func genericMarshal(in interface{}) ([]byte, error) {
	if reflect.TypeOf(in).Kind() == reflect.Ptr {
		in = reflect.ValueOf(in).Elem().Interface()
	}

	if reflect.TypeOf(in).Kind() != reflect.Struct {
		return nil, errors.New("Calling marshaller with non-struct")
	}

	var buf = new(bytes.Buffer)

	st := reflect.ValueOf(in)
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
			v := valField.Interface().(Ptr).MarshalSII()
			buf.WriteString(fmt.Sprintf(" %s: %s\n", attributeName, v))
			continue

		case reflect.TypeOf(Placement{}):
			v, err := valField.Interface().(Placement).MarshalSII()
			if err != nil {
				return nil, errors.Wrap(err, "Unable to encode Placement")
			}
			buf.WriteString(fmt.Sprintf(" %s: %s\n", attributeName, v))
			continue

		case reflect.TypeOf(RawValue{}):
			v, err := valField.Interface().(RawValue).MarshalSII()
			if err != nil {
				return nil, errors.Wrap(err, "Unable to encode RawValue")
			}
			buf.WriteString(fmt.Sprintf(" %s: %s\n", attributeName, v))
			continue

		}

		switch typeField.Type.Kind() {

		case reflect.Bool:
			v := valField.Bool()
			buf.WriteString(fmt.Sprintf(" %s: %s\n", attributeName, strconv.FormatBool(v)))

		case reflect.Float32:
			v, err := float2sii(float32(valField.Float()))
			if err != nil {
				return nil, errors.Wrap(err, "Unable to encode float32")
			}
			buf.WriteString(fmt.Sprintf(" %s: %s\n", attributeName, v))

		case reflect.Int, reflect.Int64:
			buf.WriteString(fmt.Sprintf(" %s: %d\n", attributeName, valField.Int()))

		case reflect.String:
			buf.WriteString(fmt.Sprintf(" %s: %q\n", attributeName, valField.String()))

		case reflect.Uint64:
			v := strconv.FormatUint(valField.Uint(), 10)
			buf.WriteString(fmt.Sprintf(" %s: %s\n", attributeName, v))

		case reflect.Array:

			switch typeField.Type.Elem().Kind() {

			case reflect.Float32:
				var vals [][]byte

				switch typeField.Type.Len() {

				case 3:
					for _, v := range valField.Interface().([3]float32) {
						bv, err := float2sii(v)
						if err != nil {
							return nil, errors.Wrap(err, "Unable to encode float32")
						}
						vals = append(vals, bv)
					}

				default:
					return nil, errors.Errorf("Unsupported type: [%d]%s", typeField.Type.Len(), typeField.Type.Elem().Kind())

				}

				buf.WriteString(fmt.Sprintf(" %s: (%s)\n", attributeName, bytes.Join(vals, []byte(", "))))

			default:
				return nil, errors.Errorf("Unsupported type: [%d]%s", typeField.Type.Len(), typeField.Type.Elem().Kind())

			}

		case reflect.Ptr:

			switch typeField.Type.Elem().Kind() {

			case reflect.Int64:
				var v string
				if valField.IsNil() {
					v = "nil"
				} else {
					v = strconv.FormatInt(valField.Elem().Int(), 10)
				}
				buf.WriteString(fmt.Sprintf(" %s: %s\n", attributeName, v))

			default:
				return nil, errors.Errorf("Unsupported type: *%s", typeField.Type.Elem().Kind())

			}

		case reflect.Slice:
			var values []string

			switch typeField.Type.Elem() {

			case reflect.TypeOf(Ptr{}):
				for _, val := range valField.Interface().([]Ptr) {
					values = append(values, string(val.MarshalSII()))
				}
				buf.Write(encodeSliceValue(attributeName, values))
				continue

			case reflect.TypeOf(Placement{}):
				for _, val := range valField.Interface().([]Placement) {
					ev, err := val.MarshalSII()
					if err != nil {
						return nil, errors.Wrap(err, "Unable to encode Placement")
					}
					values = append(values, string(ev))
				}
				buf.Write(encodeSliceValue(attributeName, values))
				continue

			case reflect.TypeOf(RawValue{}):
				for _, val := range valField.Interface().([]RawValue) {
					ev, err := val.MarshalSII()
					if err != nil {
						return nil, errors.Wrap(err, "Unable to encode RawValue")
					}
					values = append(values, string(ev))
				}
				buf.Write(encodeSliceValue(attributeName, values))
				continue

			}

			switch typeField.Type.Elem().Kind() {

			case reflect.Bool:
				for _, val := range valField.Interface().([]bool) {
					values = append(values, strconv.FormatBool(val))
				}
				buf.Write(encodeSliceValue(attributeName, values))

			case reflect.Float32:
				for _, val := range valField.Interface().([]float32) {
					v, err := float2sii(val)
					if err != nil {
						return nil, errors.Wrap(err, "Unable to encode float32")
					}
					values = append(values, string(v))
				}
				buf.Write(encodeSliceValue(attributeName, values))

			case reflect.Int:
				for _, val := range valField.Interface().([]int) {
					values = append(values, strconv.FormatInt(int64(val), 10))
				}
				buf.Write(encodeSliceValue(attributeName, values))

			case reflect.Int64:
				for _, val := range valField.Interface().([]int64) {
					values = append(values, strconv.FormatInt(val, 10))
				}
				buf.Write(encodeSliceValue(attributeName, values))

			case reflect.String:
				for _, val := range valField.Interface().([]string) {
					values = append(values, fmt.Sprintf("%q", val))
				}
				buf.Write(encodeSliceValue(attributeName, values))

			default:
				return nil, errors.Errorf("Unsupported type: []%s", typeField.Type.Elem().Kind())

			}

		default:
			return nil, errors.Errorf("Unsupported type: %s", typeField.Type.Kind())

		}

	}

	return buf.Bytes(), nil
}

func encodeSliceValue(attributeName string, values []string) []byte {
	var buf = new(bytes.Buffer)

	fmt.Fprintf(buf, " %s: %d\n", attributeName, len(values))

	for i, v := range values {
		fmt.Fprintf(buf, " %s[%d]: %s\n", attributeName, i, v)
	}

	return buf.Bytes()
}
