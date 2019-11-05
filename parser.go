package sii

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var blockStartRegex = regexp.MustCompile(`^([^\s:]+)\s?:\s?([^\s]+)(?:\s?\{)?$`)

func parseSIIPlainFile(r io.Reader) (*Unit, error) {
	var (
		blockContent []byte
		blockName    string
		blockClass   string
		inBlock      = false
		inComment    = false
		inUnit       = false
		scanner      = bufio.NewScanner(r)
		unit         = &Unit{}
	)

	for scanner.Scan() {
		var line = strings.TrimSpace(scanner.Text())

		switch {

		case line == "{":
			if !inUnit {
				inUnit = true
				continue
			}

			if !inBlock {
				if blockClass == "" || blockName == "" {
					return nil, errors.New("Unpexpected block open without unit class / name")
				}

				inBlock = true
				continue
			}

			return nil, errors.New("Unexpected opening braces")

		case line == "}":
			if inBlock {
				if err := processBlock(unit, blockClass, blockName, blockContent); err != nil {
					return nil, errors.Wrap(err, "Unable to process block")
				}

				inBlock = false
				blockClass = ""
				blockName = ""
				continue
			}

			if inUnit {
				inUnit = false
				continue
			}

			return nil, errors.New("Unexpected closing braces")

		case blockStartRegex.MatchString(line) && !inBlock:
			if !inUnit {
				return nil, errors.New("Unexpected block start outside unit")
			}

			groups := blockStartRegex.FindStringSubmatch(line)

			blockClass = groups[1]
			blockName = groups[2]

			if strings.HasSuffix(line, `{`) {
				inBlock = true
			}

		case (strings.HasPrefix(line, `\*`) && strings.HasSuffix(line, `*\`)) || strings.HasPrefix(line, `#`) || strings.HasPrefix(line, `//`):
			// one-line-comment, just drop

		case strings.HasPrefix(line, `/*`):
			inComment = true

		case strings.HasSuffix(line, `*/`):
			inComment = false

		default:
			if inComment {
				// Inside multi-line-comment, just drop
				continue
			}

			if !inBlock {
				// Outside block, drop line
				continue
			}

			// Append line to block content
			blockContent = bytes.Join([][]byte{
				blockContent,
				scanner.Bytes(),
			}, []byte{'\n'})

		}

	}

	if scanner.Err() != nil {
		return nil, errors.Wrap(scanner.Err(), "Unable to scan file")
	}

	return unit, nil
}

func processBlock(unit *Unit, blockClass, blockName string, blockContent []byte) error {
	block := getBlockInstance(blockClass)
	block.Init(blockClass, blockName)

	var err error
	if reflect.TypeOf(block).Implements(reflect.TypeOf((*Unmarshaler)(nil)).Elem()) {
		err = block.(Unmarshaler).UnmarshalSII(blockContent)
	} else {
		// TODO: Add generic unmarshal
	}

	if err != nil {
		return errors.Wrap(err, "Unable to unmarshal block content")
	}

	unit.Entries = append(unit.Entries, block)

	return nil
}
