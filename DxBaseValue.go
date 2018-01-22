package DxValue

import (
	"strconv"
	"encoding/base64"
	"github.com/suiyunonghen/DxCommonLib"
	"unsafe"
)

/*import (
	""
)*/



type(
	DxValueType				uint8
	DxBinaryEncodeType		uint8
	IDxValue			interface{
		ValueType()DxValueType
	}
	IDxIntValue			interface{
		AsInt32()int32
		AsInt()int
		AsInt64()int64
	}
	IDxBoolValue		interface{
		AsBoolean()bool
	}
	IDxArrayValue		interface{
		AsArray()[]IDxValue
		Count()int
	}

	DxBaseValue			struct{
		fValueType		DxValueType
	}


	DxInt32Value struct {
		DxBaseValue
		fvalue		int32
	}

	DxIntValue struct {
		DxBaseValue
		fvalue		int
	}

	DxInt64Value struct {
		DxBaseValue
		fvalue		int64
	}

	DxBoolValue struct {
		DxBaseValue
		fvalue		bool
	}

	DxFloatValue struct {
		DxBaseValue
		fvalue		float32
	}

	DxDoubleValue struct {
		DxBaseValue
		fvalue		float64
	}

	DxStringValue struct {
		DxBaseValue
		fvalue		string
	}

	DxBinaryValue struct {
		DxBaseValue
		EncodeType	DxBinaryEncodeType
		fbinary 	[]byte
	}
)


const(
	DVT_Unknown DxValueType = iota
	DVT_Null
	DVT_Int
	DVT_Int32
	DVT_Int64
	DVT_Bool
	DVT_Float
	DVT_Double
	DVT_String
	DVT_Binary
	DVT_Record
	DVT_Array
)

const(
	BET_Hex DxBinaryEncodeType = iota
	BET_Base64
)

func (v DxBaseValue)ValueType()DxValueType  {
	return v.fValueType
}

func (v *DxBaseValue)ToString()string  {
	switch v.fValueType {
	case DVT_String:
		return (*DxStringValue)(unsafe.Pointer(v)).fvalue
	case DVT_Binary:
		return (*DxBinaryValue)(unsafe.Pointer(v)).ToString()
	case DVT_Float:
		return (*DxFloatValue)(unsafe.Pointer(v)).ToString()
	case DVT_Double:
		return (*DxDoubleValue)(unsafe.Pointer(v)).ToString()
	case DVT_Bool:
		return (*DxBoolValue)(unsafe.Pointer(v)).ToString()
	case DVT_Record:
		return (*DxRecord)(unsafe.Pointer(v)).ToString()
	case DVT_Int32:
		return (*DxInt32Value)(unsafe.Pointer(v)).ToString()
	case DVT_Int:
		return (*DxIntValue)(unsafe.Pointer(v)).ToString()
	case DVT_Int64:
		return (*DxInt64Value)(unsafe.Pointer(v)).ToString()
	case DVT_Array:
		return ""
	default:
		return ""
	}
}

func (v *DxBaseValue)Bytes()[]byte  {
	switch v.fValueType {
	case DVT_String:
		return (*DxStringValue)(unsafe.Pointer(v)).Bytes()
	case DVT_Binary:
		return (*DxBinaryValue)(unsafe.Pointer(v)).Bytes()
	case DVT_Float:
		return (*DxFloatValue)(unsafe.Pointer(v)).Bytes()
	case DVT_Double:
		return (*DxDoubleValue)(unsafe.Pointer(v)).Bytes()
	case DVT_Bool:
		return (*DxBoolValue)(unsafe.Pointer(v)).Bytes()
	case DVT_Record:
		return (*DxRecord)(unsafe.Pointer(v)).Bytes()
	case DVT_Int32:
		return (*DxInt32Value)(unsafe.Pointer(v)).Bytes()
	case DVT_Int:
		return (*DxIntValue)(unsafe.Pointer(v)).Bytes()
	case DVT_Int64:
		return (*DxInt64Value)(unsafe.Pointer(v)).Bytes()
	case DVT_Array:
		return nil
	default:
		return nil
	}
}

func (v *DxInt32Value)ToString()string  {
	return strconv.Itoa(int(v.fvalue))
}

func (v *DxInt32Value)Bytes()[]byte  {
	mb := make([]byte,0,unsafe.Sizeof(v.fvalue))
	*(*int32)(unsafe.Pointer(&mb[0])) = v.fvalue
	return mb
}



func (v *DxIntValue)ToString()string  {
	return strconv.Itoa(v.fvalue)
}

func (v *DxIntValue)Bytes()[]byte  {
	mb := make([]byte,0,unsafe.Sizeof(v.fvalue))
	*(*int)(unsafe.Pointer(&mb[0])) = v.fvalue
	return mb
}

func (v *DxInt64Value)ToString()string  {
	return strconv.FormatInt(int64(v.fvalue),10)
}

func (v *DxInt64Value)Bytes()[]byte  {
	mb := make([]byte,0,unsafe.Sizeof(v.fvalue))
	*(*int64)(unsafe.Pointer(&mb[0])) = v.fvalue
	return mb
}


func (v DxBoolValue)ToString()string  {
	if v.fvalue {
		return "true"
	}
	return "false"
}

func (v *DxBoolValue)Bytes()[]byte  {
	if v.fvalue{
		return []byte{1}
	}
	return []byte{0}
}


func (v *DxFloatValue)ToString()string  {
	return strconv.FormatFloat(float64(v.fvalue),'f','e',32)
}

func (v *DxFloatValue)Bytes()[]byte  {
	mb := make([]byte,0,unsafe.Sizeof(v.fvalue))
	*(*float32)(unsafe.Pointer(&mb[0])) = v.fvalue
	return mb
}


func (v *DxDoubleValue)ToString()string  {
	return strconv.FormatFloat(float64(v.fvalue),'f','e',64)
}

func (v *DxDoubleValue)Bytes()[]byte  {
	mb := make([]byte,0,unsafe.Sizeof(v.fvalue))
	*(*float64)(unsafe.Pointer(&mb[0])) = v.fvalue
	return mb
}


func (v *DxStringValue)ToString()string  {
	return v.fvalue
}

func (v *DxStringValue)Bytes()[]byte  {
	return ([]byte)(v.fvalue)
}

func (v *DxBinaryValue)ValueType()DxValueType  {
	return DVT_Binary
}

func (v *DxBinaryValue)ToString()string  {
	if v.fbinary == nil || len(v.fbinary) == 0{
		return ""
	}
	if v.EncodeType == BET_Base64{
		return base64.StdEncoding.EncodeToString(v.fbinary)
	}
	return DxCommonLib.Binary2Hex(v.fbinary)
}

func (v *DxBinaryValue)Append(b []byte)  {
	l := len(b)
	if l == 0{
		return
	}
	if v.fbinary == nil{
		v.fbinary = make([]byte,0,l)
		copy(v.fbinary,b)
		return
	}
	v.fbinary = append(v.fbinary,b...)
}

func (v *DxBinaryValue)SetBinary(b []byte,reSet bool)  {
	if reSet{
		v.fbinary = b
		return
	}
	l := len(b)
	if v.fbinary != nil{
		lb := len(v.fbinary)
		if lb >= l{
			v.fbinary = v.fbinary[:l]
			if l > 0{
				copy(v.fbinary,b)
			}
			return
		}
		copy(v.fbinary,b[:lb])
		v.fbinary = append(v.fbinary,b[lb:]...)
		return
	}
	v.fbinary = make([]byte,0,l)
	copy(v.fbinary,b)
}

func (v *DxBinaryValue)Bytes()[]byte  {
	return v.fbinary
}