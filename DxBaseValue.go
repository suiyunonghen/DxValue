package DxValue

import (
	"strconv"
	"encoding/base64"
	"github.com/suiyunonghen/DxCommonLib"
	"unsafe"
	"errors"
	"strings"
	"time"
	"io"
)


type(
	DxValueType				uint8
	DxBinaryEncodeType		uint8
	IDxValue			interface{
		ValueType()DxValueType
		CanParent()bool
		AsInt()(int,error)
		AsBool()(bool,error)
		AsInt32()(int32,error)
		AsInt64()(int64,error)
		AsArray()(*DxArray,error)
		AsRecord()(*DxRecord,error)
		AsString()string
		AsFloat()(float32,error)
		AsDouble()(float64,error)
		AsBytes()([]byte,error)
		AsDateTime()(DxCommonLib.TDateTime,error)
	}

	IDxValueCoder		interface{
		Encode(v *DxBaseValue,w io.Writer)
		Decode(r io.Reader)*DxBaseValue
	}

	DxBaseValue			struct{
		fValueType		DxValueType
		fParent			*DxBaseValue
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
	DVT_DateTime
	DVT_String
	DVT_Binary
	DVT_Array
	DVT_Record
	DVT_RecordIntKey
)

const(
	BET_Hex DxBinaryEncodeType = iota
	BET_Base64
)

var (
	ErrValueType = errors.New("Value Data Type not Match")
	ErrInvalidateJson = errors.New("Is not a Validate Json format")
)
func (v DxBaseValue)ValueType()DxValueType  {
	return v.fValueType
}

func (v *DxBaseValue)AsInt()(int,error){
	switch v.fValueType {
	case DVT_Int:
		return (*DxIntValue)(unsafe.Pointer(v)).fvalue,nil
	case DVT_Int32:
		return int((*DxInt32Value)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Int64:
		return int((*DxInt64Value)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Float:
		return int((*DxFloatValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Double,DVT_DateTime:
		return int((*DxDoubleValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Bool:
		if (*DxBoolValue)(unsafe.Pointer(v)).fvalue{
			return 1,nil
		}
		return 0,nil
	case DVT_Null:
		return 0,nil
	case DVT_String:
		return  strconv.Atoi((*DxStringValue)(unsafe.Pointer(v)).fvalue)
	default:
		return 0,ErrValueType
	}
}

func (v *DxBaseValue)AsDateTime()(DxCommonLib.TDateTime,error){
	switch v.fValueType {
	case DVT_Int:
		return (DxCommonLib.TDateTime)((*DxIntValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Int32:
		return (DxCommonLib.TDateTime)((*DxInt32Value)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Int64:
		return (DxCommonLib.TDateTime)((*DxInt64Value)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Float:
		return (DxCommonLib.TDateTime)((*DxFloatValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Double,DVT_DateTime:
		return (DxCommonLib.TDateTime)((*DxDoubleValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_String:
		t,err := time.Parse("2006-01-02T15:04:05Z",(*DxStringValue)(unsafe.Pointer(v)).fvalue)
		if err == nil{
			return DxCommonLib.Time2DelphiTime(&t),err
		}
		t,err = time.Parse("2006-01-02 15:04:05",(*DxStringValue)(unsafe.Pointer(v)).fvalue)
		if err == nil{
			return DxCommonLib.Time2DelphiTime(&t),err
		}
		t,err = time.Parse("2006/01/02 15:04:05",(*DxStringValue)(unsafe.Pointer(v)).fvalue)
		if err == nil{
			return DxCommonLib.Time2DelphiTime(&t),err
		}
		return -1,err
	default:
		return -1,ErrValueType
	}
}

func (v *DxBaseValue)AsBool()(bool,error){
	switch v.fValueType {
	case DVT_Int:
		return (*DxIntValue)(unsafe.Pointer(v)).fvalue != 0,nil
	case DVT_Int32:
		return int((*DxInt32Value)(unsafe.Pointer(v)).fvalue) != 0,nil
	case DVT_Int64:
		return int((*DxInt64Value)(unsafe.Pointer(v)).fvalue) != 0,nil
	case DVT_Float:
		return int((*DxFloatValue)(unsafe.Pointer(v)).fvalue) != 0,nil
	case DVT_Double,DVT_DateTime:
		return int((*DxDoubleValue)(unsafe.Pointer(v)).fvalue) != 0,nil
	case DVT_Bool: return (*DxBoolValue)(unsafe.Pointer(v)).fvalue,nil
	case DVT_Null: return false,nil
	case DVT_String:
		return strings.ToUpper((*DxStringValue)(unsafe.Pointer(v)).fvalue)== "TRUE",nil
	default:
		return false,ErrValueType
	}
}

func (v *DxBaseValue)AsInt32()(int32,error){
	switch v.fValueType {
	case DVT_Int:
		return int32((*DxIntValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Int32:
		return (*DxInt32Value)(unsafe.Pointer(v)).fvalue,nil
	case DVT_Int64:
		return int32((*DxInt64Value)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Float:
		return int32((*DxFloatValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Double,DVT_DateTime:
		return int32((*DxDoubleValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Bool:
		if (*DxBoolValue)(unsafe.Pointer(v)).fvalue{
			return 1,nil
		}
		return 0,nil
	case DVT_Null:
		return 0,nil
	case DVT_String:
		rv,err := strconv.Atoi((*DxStringValue)(unsafe.Pointer(v)).fvalue)
		return int32(rv),err
	default:
		return 0,ErrValueType
	}
}

func (v *DxBaseValue)AsInt64()(int64,error){
	switch v.fValueType {
	case DVT_Int:
		return int64((*DxIntValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Int32:
		return int64((*DxInt32Value)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Int64:
		return (*DxInt64Value)(unsafe.Pointer(v)).fvalue,nil
	case DVT_Float:
		return int64((*DxFloatValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Double,DVT_DateTime:
		return int64((*DxDoubleValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Bool:
		if (*DxBoolValue)(unsafe.Pointer(v)).fvalue{
			return 1,nil
		}
		return 0,nil
	case DVT_Null:
		return 0,nil
	case DVT_String:
		rv,err := strconv.Atoi((*DxStringValue)(unsafe.Pointer(v)).fvalue)
		return int64(rv),err
	default:
		return 0,ErrValueType
	}
}

func (v *DxBaseValue)AsArray()(*DxArray,error){
	if v.fValueType == DVT_Array{
		return (*DxArray)(unsafe.Pointer(v)),nil
	}
	return nil,ErrValueType
}

func (v *DxBaseValue)AsRecord()(*DxRecord,error){
	if v.fValueType == DVT_Record{
		return (*DxRecord)(unsafe.Pointer(v)),nil
	}
	return nil,ErrValueType
}

func (v *DxBaseValue)AsIntRecord()(*DxIntKeyRecord,error){
	if v.fValueType == DVT_RecordIntKey{
		return (*DxIntKeyRecord)(unsafe.Pointer(v)),nil
	}
	return nil,ErrValueType
}

func (v *DxBaseValue)AsString()string{
	return v.ToString()
}

func (v *DxBaseValue)AsFloat()(float32,error){
	switch v.fValueType {
	case DVT_Int:
		return float32((*DxIntValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Int32:
		return float32((*DxInt32Value)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Int64:
		return float32((*DxInt64Value)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Float:
		return (*DxFloatValue)(unsafe.Pointer(v)).fvalue,nil
	case DVT_Double,DVT_DateTime:
		return float32((*DxDoubleValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Bool:
		if (*DxBoolValue)(unsafe.Pointer(v)).fvalue{
			return 1,nil
		}
		return 0,nil
	case DVT_Null:
		return 0,nil
	case DVT_String:
		rv,err := strconv.ParseFloat((*DxStringValue)(unsafe.Pointer(v)).fvalue,32)
		return float32(rv),err
	default:
		return 0,ErrValueType
	}
}

func (v *DxBaseValue)AsDouble()(float64,error){
	switch v.fValueType {
	case DVT_Int:
		return float64((*DxIntValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Int32:
		return float64((*DxInt32Value)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Int64:
		return float64((*DxInt64Value)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Float:
		return float64((*DxFloatValue)(unsafe.Pointer(v)).fvalue),nil
	case DVT_Double,DVT_DateTime:
		return (*DxDoubleValue)(unsafe.Pointer(v)).fvalue,nil
	case DVT_Bool:
		if (*DxBoolValue)(unsafe.Pointer(v)).fvalue{
			return 1,nil
		}
		return 0,nil
	case DVT_Null:
		return 0,nil
	case DVT_String:
		rv,err := strconv.ParseFloat((*DxStringValue)(unsafe.Pointer(v)).fvalue,64)
		return rv,err
	default:
		return 0,ErrValueType
	}
}


func (v *DxBaseValue)CanParent()bool  {
	return v.fValueType >= DVT_Array
}

func (v *DxBaseValue)Parent()*DxBaseValue  {
	return v.fParent
}

func (v *DxBaseValue)ClearValue(clearInner bool)  {
	switch v.fValueType {
	case DVT_Record: (*DxRecord)(unsafe.Pointer(v)).ClearValue(clearInner)
	case DVT_RecordIntKey: (*DxIntKeyRecord)(unsafe.Pointer(v)).ClearValue(clearInner)
	case DVT_Binary: (*DxBinaryValue)(unsafe.Pointer(v)).ClearValue(clearInner)
	case DVT_Array:  (*DxArray)(unsafe.Pointer(v)).ClearValue(clearInner)
	default:
		v.fParent = nil
	}
}

func (v *DxBaseValue)ToString()string  {
	switch v.fValueType {
	case DVT_String:
		return (*DxStringValue)(unsafe.Pointer(v)).fvalue
	case DVT_Binary:
		return (*DxBinaryValue)(unsafe.Pointer(v)).ToString()
	case DVT_Float:
		return (*DxFloatValue)(unsafe.Pointer(v)).ToString()
	case DVT_Double,DVT_DateTime:
		return (*DxDoubleValue)(unsafe.Pointer(v)).ToString()
	case DVT_Bool:
		return (*DxBoolValue)(unsafe.Pointer(v)).ToString()
	case DVT_Record:
		return (*DxRecord)(unsafe.Pointer(v)).ToString()
	case DVT_RecordIntKey:
		return (*DxIntKeyRecord)(unsafe.Pointer(v)).ToString()
	case DVT_Int32:
		return (*DxInt32Value)(unsafe.Pointer(v)).ToString()
	case DVT_Int:
		return (*DxIntValue)(unsafe.Pointer(v)).ToString()
	case DVT_Int64:
		return (*DxInt64Value)(unsafe.Pointer(v)).ToString()
	case DVT_Array:
		return (*DxArray)(unsafe.Pointer(v)).ToString()
	default:
		return ""
	}
}

func (v *DxBaseValue)AsBytes()([]byte,error)  {
	switch v.fValueType {
	case DVT_String:
		return (*DxStringValue)(unsafe.Pointer(v)).Bytes(),nil
	case DVT_Binary:
		return (*DxBinaryValue)(unsafe.Pointer(v)).Bytes(),nil
	case DVT_Float:
		return (*DxFloatValue)(unsafe.Pointer(v)).Bytes(),nil
	case DVT_Double,DVT_DateTime:
		return (*DxDoubleValue)(unsafe.Pointer(v)).Bytes(),nil
	case DVT_Bool:
		return (*DxBoolValue)(unsafe.Pointer(v)).Bytes(),nil
	case DVT_Record:
		return (*DxRecord)(unsafe.Pointer(v)).Bytes(),nil
	case DVT_RecordIntKey:
		return (*DxIntKeyRecord)(unsafe.Pointer(v)).Bytes(),nil
	case DVT_Int32:
		return (*DxInt32Value)(unsafe.Pointer(v)).Bytes(),nil
	case DVT_Int:
		return (*DxIntValue)(unsafe.Pointer(v)).Bytes(),nil
	case DVT_Int64:
		return (*DxInt64Value)(unsafe.Pointer(v)).Bytes(),nil
	case DVT_Array:
		return (*DxArray)(unsafe.Pointer(v)).Bytes(),nil
	default:
		return nil,nil
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
	if v.fValueType == DVT_DateTime{
		return DxCommonLib.TDateTime(v.fvalue).ToTime().Format("2006-01-02 15:04:05")
	}
	return strconv.FormatFloat(float64(v.fvalue),'f','e',64)
}

func (v *DxDoubleValue)Bytes()[]byte  {
	mb := make([]byte,0,unsafe.Sizeof(v.fvalue))
	*(*float64)(unsafe.Pointer(&mb[0])) = v.fvalue
	return mb
}

func (v *DxDoubleValue)Time()time.Time  {
	return DxCommonLib.TDateTime(v.fvalue).ToTime()
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
	switch v.EncodeType {
	case BET_Base64: return base64.StdEncoding.EncodeToString(v.fbinary)
	case BET_Hex: return DxCommonLib.Binary2Hex(v.fbinary)
	}
	return ""
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

func (v *DxBinaryValue)ClearValue(clearInner bool)  {
	if clearInner{
		v.fbinary = nil
	}
	if v.fbinary != nil{
		v.fbinary = v.fbinary[:0]
	}
}

func (v *DxBinaryValue)Bytes()[]byte  {
	return v.fbinary
}

func IsSpace(b byte)bool  {
	return b == ' ' || b == '\r' || b == '\n' || b == '\t'
}