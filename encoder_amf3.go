package amf

import (
	"encoding/binary"
	"io"
	"reflect"
)

// amf3 polymorphic router

func (e *Encoder) EncodeAmf3(w io.Writer, val interface{}) (int, error) {
	if val == nil {
		return e.EncodeAmf3Null(w, true)
	}

	v := reflect.ValueOf(val)
	if !v.IsValid() {
		return e.EncodeAmf3Null(w, true)
	}

	switch v.Kind() {
	case reflect.String:
		return e.EncodeAmf3String(w, v.String(), true)
	case reflect.Bool:
		if v.Bool() {
			return e.EncodeAmf3True(w, true)
		} else {
			return e.EncodeAmf3False(w, true)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return e.EncodeAmf3Integer(w, uint32(v.Int()), true)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return e.EncodeAmf3Integer(w, uint32(v.Uint()), true)
	case reflect.Int64:
		return e.EncodeAmf3Double(w, float64(v.Int()), true)
	case reflect.Uint64:
		return e.EncodeAmf3Double(w, float64(v.Uint()), true)
	case reflect.Float32, reflect.Float64:
		return e.EncodeAmf3Double(w, float64(v.Float()), true)
	case reflect.Array, reflect.Slice:
		length := v.Len()
		arr := make(Array, length)
		for i := 0; i < length; i++ {
			arr[i] = v.Index(int(i)).Interface()
		}
		return e.EncodeAmf3Array(w, arr, true)
	case reflect.Map:
		return 0, Error("encode amf3: unsupported type object")
	}

	if _, ok := val.(TypedObject); ok {
		return 0, Error("encode amf3: unsupported type typed object")
	}

	return 0, Error("encode amf3: unsupported type %s", v.Type())
}

// marker: 1 byte 0x00
// no additional data
func (e *Encoder) EncodeAmf3Undefined(w io.Writer, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF3_UNDEFINED_MARKER); err != nil {
			return
		}
		n += 1
	}

	return
}

// marker: 1 byte 0x01
// no additional data
func (e *Encoder) EncodeAmf3Null(w io.Writer, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF3_NULL_MARKER); err != nil {
			return
		}
		n += 1
	}

	return
}

// marker: 1 byte 0x02
// no additional data
func (e *Encoder) EncodeAmf3False(w io.Writer, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF3_FALSE_MARKER); err != nil {
			return
		}
		n += 1
	}

	return
}

// marker: 1 byte 0x03
// no additional data
func (e *Encoder) EncodeAmf3True(w io.Writer, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF3_TRUE_MARKER); err != nil {
			return
		}
		n += 1
	}

	return
}

// marker: 1 byte 0x04
func (e *Encoder) EncodeAmf3Integer(w io.Writer, val uint32, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF3_INTEGER_MARKER); err != nil {
			return
		}
		n += 1
	}

	var m int

	if val <= 0x0000007F {
		err = WriteByte(w, byte(val))
		if err == nil {
			n += 1
		}
		return
	}

	if val <= 0x00003FFF {
		m, err = w.Write([]byte{byte(val>>7 | 0x80), byte(val & 0x7F)})
		n += m
		return
	}

	if val <= 0x001FFFFF {
		m, err = w.Write([]byte{byte(val>>14 | 0x80), byte(val>>7&0x7F | 0x80), byte(val & 0x7F)})
		n += m
		return
	}

	if val <= 0x1FFFFFFF {
		m, err = w.Write([]byte{byte(val>>22 | 0x80), byte(val>>15&0x7F | 0x80), byte(val>>8&0x7F | 0x80), byte(val)})
		n += m
		return
	}

	return n, Error("amf3 encode: cannot encode u29 (out of range)")
}

// marker: 1 byte 0x05
func (e *Encoder) EncodeAmf3Double(w io.Writer, val float64, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF3_DOUBLE_MARKER); err != nil {
			return
		}
		n += 1
	}

	err = binary.Write(w, binary.BigEndian, &val)
	if err != nil {
		return
	}
	n += 8

	return
}

// marker: 1 byte 0x06
// format:
// - u29 reference int. if reference, no more data. if not reference,
//   length value of bytes to read to complete string.
func (e *Encoder) EncodeAmf3String(w io.Writer, val string, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF3_STRING_MARKER); err != nil {
			return
		}
		n += 1
	}

	var m int

	for i, s := range e.stringRefs {
		if s == val {
			u29 := uint32(i<<1 | 0x01)
			m, err = e.EncodeAmf3Integer(w, u29, false)
			if err != nil {
				n += m
			}
			return
		}
	}

	length := uint32(len(val))
	u29 := uint32(length<<1) | 0x01

	m, err = e.EncodeAmf3Integer(w, u29, false)
	if err != nil {
		return n, Error("amf3 encode: cannot encode u29 for string: %s", err)
	}
	n += m

	m, err = w.Write([]byte(val))
	if err != nil {
		return n, Error("encode amf3: unable to encode string value: %s", err)
	}
	n += m

	if val != "" {
		e.stringRefs = append(e.stringRefs, val)
	}

	return
}

// marker: 1 byte 0x09
// format:
// - u29 reference int. if reference, no more data.
// - string representing associative array if present
// - n values (length of u29)
func (e *Encoder) EncodeAmf3Array(w io.Writer, val Array, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF3_ARRAY_MARKER); err != nil {
			return
		}
		n += 1
	}

	var m int
	length := uint32(len(val))
	u29 := uint32(length<<1) | 0x01

	m, err = e.EncodeAmf3Integer(w, u29, false)
	if err != nil {
		return n, Error("amf3 encode: cannot encode u29 for array: %s", err)
	}
	n += m

	m, err = e.EncodeAmf3String(w, "", false)
	if err != nil {
		return n, Error("amf3 encode: cannot encode empty string for array: %s", err)
	}
	n += m

	for _, v := range val {
		m, err := e.EncodeAmf3(w, v)
		if err != nil {
			return n, Error("amf3 encode: cannot encode array element: %s", err)
		}
		n += m
	}

	return
}
