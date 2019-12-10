package sii

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/pkg/errors"
)

func writeSIIPlainFile(unit *Unit, w io.Writer) error {
	var err error

	// Write file header
	if _, err = fmt.Fprintln(w, "SiiNunit\n{"); err != nil {
		return errors.Wrap(err, "Unable to write header")
	}

	for _, block := range unit.Entries {
		// Write block header
		if _, err = fmt.Fprintf(w, "%s : %s {\n", block.Class(), block.Name()); err != nil {
			return errors.Wrap(err, "Unable to write block header")
		}

		// Obtain and write block content
		var raw []byte

		if reflect.TypeOf(block).Elem().Implements(reflect.TypeOf((*Marshaler)(nil)).Elem()) {
			raw, err = block.(Marshaler).MarshalSII()
		} else {
			raw, err = genericMarshal(block)
		}

		if err != nil {
			return errors.Wrap(err, "Unable to marshal block")
		}

		if len(raw) > 0 {
			raw = append(bytes.TrimRight(raw, "\n"), '\n')
		}

		if _, err = w.Write(raw); err != nil {
			return errors.Wrap(err, "Unable to write block data")
		}

		// Close block
		if _, err = fmt.Fprintf(w, "}\n\n"); err != nil {
			return errors.Wrap(err, "Unable to close block")
		}
	}

	// Write file footer
	if _, err = fmt.Fprintln(w, "}"); err != nil {
		return errors.Wrap(err, "Unable to write footer")
	}

	return nil
}
