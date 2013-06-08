package amf

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func EncodeAndDecode(val interface{}, ver Version) (result interface{}, err error) {
	enc := new(Encoder)
	dec := new(Decoder)

	buf := new(bytes.Buffer)

	_, err = enc.Encode(buf, val, ver)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in encode: %s", err))
	}

	result, err = dec.Decode(buf, ver)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in decode: %s", err))
	}

	return
}

func Compare(val interface{}, ver Version, name string, t *testing.T) {
	result, err := EncodeAndDecode(val, ver)
	if err != nil {
		t.Errorf("%s: %s", name, err)
	}
	if val != result {
		t.Errorf("%s: comparison failed between %+v and %+v", name, val, result)
	}
}

func TestAmf0Number(t *testing.T) {
	Compare(float64(3.14159), 0, "amf0 number float", t)
	Compare(float64(124567890), 0, "amf0 number high", t)
	Compare(float64(-34.2), 0, "amf0 number negative", t)
}

func TestAmf0String(t *testing.T) {
	Compare("a pup!", 0, "amf0 string simple", t)
	Compare("日本語", 0, "amf0 string utf8", t)
}

func TestAmf0Boolean(t *testing.T) {
	Compare(true, 0, "amf0 boolean true", t)
	Compare(false, 0, "amf0 boolean false", t)
}

func TestAmf0Null(t *testing.T) {
	Compare(nil, 0, "amf0 boolean nil", t)
}

func TestAmf0Object(t *testing.T) {
	obj := make(Object)
	obj["dog"] = "alfie"
	obj["coffee"] = true
	obj["drugs"] = false
	obj["pi"] = 3.14159

	res, err := EncodeAndDecode(obj, 0)
	if err != nil {
		t.Errorf("amf0 object: %s", err)
	}

	result, ok := res.(Object)
	if ok != true {
		t.Errorf("amf0 object conversion failed")
	}

	if result["dog"] != "alfie" {
		t.Errorf("amf0 object string: comparison failed")
	}

	if result["coffee"] != true {
		t.Errorf("amf0 object true: comparison failed")
	}

	if result["drugs"] != false {
		t.Errorf("amf0 object false: comparison failed")
	}

	if result["pi"] != float64(3.14159) {
		t.Errorf("amf0 object float: comparison failed")
	}
}

func TestAmf0Array(t *testing.T) {
	arr := [5]float64{1, 2, 3, 4, 5}

	res, err := EncodeAndDecode(arr, 0)
	if err != nil {
		t.Error("amf0 object: %s", err)
	}

	result, ok := res.(Array)
	if ok != true {
		t.Errorf("amf0 array conversion failed")
	}

	for i := 0; i < len(arr); i++ {
		if arr[i] != result[i] {
			t.Errorf("amf0 array %d comparison failed: %v / %v", i, arr[i], result[i])
		}
	}
}

func TestAmf3Integer(t *testing.T) {
	Compare(uint32(0), 3, "amf3 integer zero", t)
	Compare(uint32(1245), 3, "amf3 integer low", t)
	Compare(uint32(123456), 3, "amf3 integer high", t)
}

func TestAmf3Double(t *testing.T) {
	Compare(float64(3.14159), 3, "amf3 double float", t)
	Compare(float64(1234567890), 3, "amf3 double high", t)
	Compare(float64(-12345), 3, "amf3 double negative", t)
}

func TestAmf3String(t *testing.T) {
	Compare("a pup!", 0, "amf0 string simple", t)
	Compare("日本語", 0, "amf0 string utf8", t)
}

func TestAmf3Boolean(t *testing.T) {
	Compare(true, 3, "amf3 boolean true", t)
	Compare(false, 3, "amf3 boolean false", t)
}

func TestAmf3Null(t *testing.T) {
	Compare(nil, 3, "amf3 boolean nil", t)
}
