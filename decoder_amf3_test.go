package amf

import (
	"bytes"
	"testing"
)

type u29TestCase struct {
	value  uint32
	expect []byte
}

var u29TestCases = []u29TestCase{
	{1, []byte{0x01}},
	{2, []byte{0x02}},
	{127, []byte{0x7F}},
	{128, []byte{0x81, 0x00}},
	{255, []byte{0x81, 0x7F}},
	{256, []byte{0x82, 0x00}},
	{0x3FFF, []byte{0xFF, 0x7F}},
	{0x4000, []byte{0x81, 0x80, 0x00}},
	{0x7FFF, []byte{0x81, 0xFF, 0x7F}},
	{0x8000, []byte{0x82, 0x80, 0x00}},
	{0x1FFFFF, []byte{0xFF, 0xFF, 0x7F}},
	{0x200000, []byte{0x80, 0xC0, 0x80, 0x00}},
	{0x3FFFFF, []byte{0x80, 0xFF, 0xFF, 0xFF}},
	{0x400000, []byte{0x81, 0x80, 0x80, 0x00}},
	{0x0FFFFFFF, []byte{0xBF, 0xFF, 0xFF, 0xFF}},
}

func TestDecodeAmf3Undefined(t *testing.T) {
	buf := bytes.NewReader([]byte{0x00})

	dec := new(Decoder)

	got, err := dec.DecodeAmf3(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if got != nil {
		t.Errorf("expect nil got %v", got)
	}
}

func TestDecodeAmf3Null(t *testing.T) {
	buf := bytes.NewReader([]byte{0x01})

	dec := new(Decoder)

	got, err := dec.DecodeAmf3(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if got != nil {
		t.Errorf("expect nil got %v", got)
	}
}

func TestDecodeAmf3False(t *testing.T) {
	buf := bytes.NewReader([]byte{0x02})
	expect := false

	dec := new(Decoder)

	got, err := dec.DecodeAmf3(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}
}

func TestDecodeAmf3True(t *testing.T) {
	buf := bytes.NewReader([]byte{0x03})
	expect := true

	dec := new(Decoder)

	got, err := dec.DecodeAmf3(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}
}

func TestDecodeAmf3Integer(t *testing.T) {
	dec := new(Decoder)

	for _, tc := range u29TestCases {
		buf := bytes.NewBuffer(tc.expect)
		n, err := dec.DecodeAmf3Integer(buf, false)
		if err != nil {
			t.Errorf("DecodeAmf3Integer error: %s", err)
		}
		if n != tc.value {
			t.Errorf("DecodeAmf3Integer expect n %x got %x", tc.value, n)
		}
	}

	buf := bytes.NewReader([]byte{0x04, 0xFF, 0xFF, 0x7F})
	expect := uint32(2097151)

	got, err := dec.DecodeAmf3(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}

	buf.Seek(0, 0)
	got, err = dec.DecodeAmf3Integer(buf, true)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}

	buf.Seek(1, 0)
	got, err = dec.DecodeAmf3Integer(buf, false)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}
}

func TestDecodeAmf3Double(t *testing.T) {
	buf := bytes.NewReader([]byte{0x05, 0x3f, 0xf3, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33})
	expect := float64(1.2)

	dec := new(Decoder)

	got, err := dec.DecodeAmf3(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}
}

func TestDecodeAmf3String(t *testing.T) {
	buf := bytes.NewReader([]byte{0x06, 0x07, 'f', 'o', 'o'})
	expect := "foo"

	dec := new(Decoder)

	got, err := dec.DecodeAmf3(buf)
	if err != nil {
		t.Errorf("%s", err)
	}
	if expect != got {
		t.Errorf("expect %v got %v", expect, got)
	}
}
