package sii

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"

	"github.com/pkg/errors"
)

func float2sii(f float32) ([]byte, error) {
	var (
		buf = new(bytes.Buffer)
		err error
	)

	if math.Floor(float64(f)) == float64(f) && f < 1000 {
		return []byte(fmt.Sprintf("%.0f", f)), nil
	}

	err = binary.Write(buf, binary.BigEndian, f)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to encode float")
	}

	dst := make([]byte, hex.EncodedLen(buf.Len()))
	hex.Encode(dst, buf.Bytes())

	return append([]byte("&"), dst...), nil
}

func sii2float(f []byte) (float32, error) {
	if f[0] != '&' {
		out, err := strconv.ParseFloat(string(f), 32)
		return float32(out), err
	}

	// Strip leading '&'
	f = f[1:]

	var (
		err error
		out float32
	)

	dst := make([]byte, hex.DecodedLen(len(f)))
	_, err = hex.Decode(dst, f)
	if err != nil {
		return 0, errors.Wrap(err, "Unable to read hex format")
	}

	err = binary.Read(bytes.NewReader(dst), binary.BigEndian, &out)
	return out, errors.Wrap(err, "Unable to decode hex notation")
}
