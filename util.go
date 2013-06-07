package amf

import (
	"errors"
	"fmt"
	"io"
)

func ReadByte(r io.Reader) (byte, error) {
	bytes, err := ReadBytes(r, 1)
	if err != nil {
		return 0x00, err
	}

	return bytes[0], nil
}

func ReadBytes(r io.Reader, n int) ([]byte, error) {
	bytes := make([]byte, n)

	m, err := r.Read(bytes)
	if err != nil {
		return bytes, err
	}

	if m != n {
		return bytes, errors.New(fmt.Sprintf("decode read bytes failed: expected %d got %d", m, n))
	}

	return bytes, nil
}

func ReadMarker(r io.Reader) (byte, error) {
	return ReadByte(r)
}

func AssertMarker(r io.Reader, x bool, m byte) error {
	if x == false {
		return nil
	}

	marker, err := ReadMarker(r)
	if err != nil {
		return err
	}

	if marker != m {
		return errors.New(fmt.Sprintf("decode assert marker failed: expected %v got %v", m, marker))
	}

	return nil
}