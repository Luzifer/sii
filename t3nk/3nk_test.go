package t3nk

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestEncodeToDecode(t *testing.T) {
	expect := "Ohai!"
	f := strings.NewReader(expect)

	r, err := Encode(f)
	if err != nil {
		t.Fatalf("Unable to encode test string: %s", err)
	}

	dr, err := Decode(r)
	if err != nil {
		t.Fatalf("Unable to decode test string: %s", err)
	}

	raw, err := ioutil.ReadAll(dr)
	if err != nil {
		t.Fatalf("Unable to read decoded test string: %s", err)
	}

	if s := string(raw); s != expect {
		t.Errorf("Did not receive expected string: exp=%q got=%q", expect, s)
	}
}
