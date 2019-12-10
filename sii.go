package sii

import (
	"bytes"
	"compress/flate"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"io"
	"os"
	"reflect"

	"github.com/pkg/errors"
)

var encryptionKey []byte

var (
	sigBinary    = []byte{0x42, 0x53, 0x49, 0x49} // BSII (unencrypted binary format)
	sigEncrypted = []byte{0x53, 0x63, 0x73, 0x43} // ScsC (compressed and encrypted SiiN with headers)
	sigPlain     = []byte{0x53, 0x69, 0x69, 0x4e} // SiiN (plain-text unit file)
)

var (
	ErrNoEncryptionKeySet = errors.New("No encryption key set")
)

type scscHeader struct {
	Signature  [4]byte
	HMAC       [32]byte
	InitVector [16]byte
	DataSize   uint32
}

// DecryptRaw takes a reader of a ScsC file and extracts the SiiN raw content
func DecryptRaw(file io.ReaderAt, size int64) (io.Reader, error) {
	if encryptionKey == nil {
		return nil, ErrNoEncryptionKeySet
	}

	ftHeader, err := readFTHeader(file)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read file header")
	}

	if !reflect.DeepEqual(ftHeader, sigEncrypted) {
		return nil, errors.New("Input file does not contain ScsC header")
	}

	h := scscHeader{}
	if err := binary.Read(io.NewSectionReader(file, 0, int64(binary.Size(h))), binary.LittleEndian, &h); err != nil {
		return nil, errors.Wrap(err, "Unable to read header from file")
	}

	var content = make([]byte, size-int64(binary.Size(h)))
	if _, err := file.ReadAt(content, int64(binary.Size(h))); err != nil {
		return nil, errors.Wrap(err, "Unable to read encrypted content")
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid encryption key")
	}
	decrypter := cipher.NewCBCDecrypter(block, h.InitVector[:])
	decrypter.CryptBlocks(content, content)

	return flate.NewReader(bytes.NewReader(content[2:])), nil
}

// ReadUnitFile reads the file, decrypts it if required and parses it into the Unit struct
func ReadUnitFile(filename string) (*Unit, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read info of save-file")
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to open encrypted file")
	}
	defer f.Close()

	// Check what type of file we do have here
	ftHeader, err := readFTHeader(f)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read file header")
	}

	var r io.Reader
	switch {
	case reflect.DeepEqual(ftHeader, sigBinary):
		return nil, errors.New("File has unsupported Binary-SII format")

	case reflect.DeepEqual(ftHeader, sigEncrypted):
		r, err = DecryptRaw(f, stat.Size())
		if err != nil {
			return nil, errors.New("Unable to decrypt file")
		}

	case reflect.DeepEqual(ftHeader, sigPlain):
		// We already got the plain file: We can just read it
		r = f

	default:
		return nil, errors.New("Invalid / unknown file type header found")
	}

	return parseSIIPlainFile(r)
}

// SetEncryptionKey sets the 32-byte key to encrypt / decrypt ScsC files.
// The key is not included for legal reasons and you need to obtain it from other sources.
func SetEncryptionKey(key []byte) { encryptionKey = key }

func WriteUnitFile(filename string, encrypt bool, data *Unit) error {
	if encrypt && encryptionKey == nil {
		return ErrNoEncryptionKeySet
	}

	var buf = new(bytes.Buffer)

	if err := writeSIIPlainFile(data, buf); err != nil {
		return errors.Wrap(err, "Unable to create SII file content")
	}

	// FIXME: Implement encryption

	// Create output file
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "Unable to open output file")
	}
	defer f.Close()

	_, err = buf.WriteTo(f)
	return errors.Wrap(err, "Unable to write buffer")
}

func readFTHeader(f io.ReaderAt) ([]byte, error) {
	var ftHeader = make([]byte, 4)
	if n, err := f.ReadAt(ftHeader, 0); err != nil || n != 4 {
		if err != nil {
			err = errors.Errorf("Received %d / 4 byte header", n)
		}
		return nil, errors.Wrap(err, "Unable to read 4-byte file header")
	}

	return ftHeader, nil
}
