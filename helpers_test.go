package sii

import "testing"

func TestSii2FloatConversion(t *testing.T) {
	var (
		err error
		f   float32
	)

	for b, exp := range map[string]float32{
		"0.00250711967": 0.00250711967,
		"1.0":           1.0,
		"&3b244e7d":     0.00250711967,
		"&3f46da61":     0.7767697,
		"&47135818":     37720.0938,
	} {
		f, err = sii2float([]byte(b))
		if err != nil {
			t.Errorf("Conversion of %q failed: %s", b, err)
			continue
		}

		if f != exp {
			t.Errorf("Conversion of %q has unxpected result: %f != %f", b, f, exp)
		}
	}
}

func TestFloat2SiiConversion(t *testing.T) {
	var (
		err error
		f   []byte
	)

	for b, exp := range map[float32]string{
		0.00250711967: "&3b244e7d",
		0.7767697:     "&3f46da61",
		37720.0938:    "&47135818",
	} {
		f, err = float2sii(b)
		if err != nil {
			t.Errorf("Conversion of %f failed: %s", b, err)
			continue
		}

		if string(f) != exp {
			t.Errorf("Conversion of %f has unxpected result: %s != %s", b, f, exp)
		}
	}
}
