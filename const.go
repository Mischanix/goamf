package amf

const (
	AMF0 = 0x00
	AMF3 = 0x03
)

const (
	AMF0_NUMBER_MARKER         = 0x00
	AMF0_BOOLEAN_MARKER        = 0x01
	AMF0_STRING_MARKER         = 0x02
	AMF0_OBJECT_MARKER         = 0x03
	AMF0_MOVIECLIP_MARKER      = 0x04
	AMF0_NULL_MARKER           = 0x05
	AMF0_UNDEFINED_MARKER      = 0x06
	AMF0_REFERENCE_MARKER      = 0x07
	AMF0_ECMA_ARRAY_MARKER     = 0x08
	AMF0_OBJECT_END_MARKER     = 0x09
	AMF0_STRICT_ARRAY_MARKER   = 0x0a
	AMF0_DATE_MARKER           = 0x0b
	AMF0_LONG_STRING_MARKER    = 0x0c
	AMF0_UNSUPPORTED_MARKER    = 0x0d
	AMF0_RECORDSET_MARKER      = 0x0e
	AMF0_XML_DOCUMENT_MARKER   = 0x0f
	AMF0_TYPED_OBJECT_MARKER   = 0x10
	AMF0_ACMPLUS_OBJECT_MARKER = 0x11
)

const (
	AMF0_BOOLEAN_FALSE  = 0x00
	AMF0_BOOLEAN_TRUE   = 0x01
	AMF0_STRING_MAX_LEN = 65535
)

const (
	AMF3_UNDEFINED_MARKER = 0x00
	AMF3_NULL_MARKER      = 0x01
	AMF3_FALSE_MARKER     = 0x02
	AMF3_TRUE_MARKER      = 0x03
	AMF3_INTEGER_MARKER   = 0x04
	AMF3_DOUBLE_MARKER    = 0x05
	AMF3_STRING_MARKER    = 0x06
	AMF3_XMLDOC_MARKER    = 0x07
	AMF3_DATE_MARKER      = 0x08
	AMF3_ARRAY_MARKER     = 0x09
	AMF3_OBJECT_MARKER    = 0x0a
	AMF3_XML_MARKER       = 0x0b
	AMF3_BYTEARRAY_MARKER = 0x0c
)

type Version uint8

type Object map[string]interface{}

type TypedObject struct {
	Type   string
	Object Object
}

func MakeObject() *Object {
	return &Object{}
}

func MakeTypedObject() *TypedObject {
	return &TypedObject{
		Type:   "",
		Object: *MakeObject(),
	}
}