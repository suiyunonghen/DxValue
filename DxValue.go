package DxValue

import (
	"unsafe"
	"math"
	"github.com/suiyunonghen/DxCommonLib"
)

/******************************************************
*  DxValue
******************************************************/
type 	DxValue		struct{
	fValue			*DxBaseValue
}

func (v *DxValue)ValueType()DxValueType  {
	if v.fValue == nil{
		return DVT_Null
	}
	return v.fValue.ValueType()
}

func (v *DxValue)SetIntValue(value int)  {
	if v.fValue != nil{
		switch v.fValue.fValueType {
		case DVT_Int:
			(*DxIntValue)(unsafe.Pointer(v.fValue)).fvalue = value
			return
		case DVT_Int64:
			(*DxInt64Value)(unsafe.Pointer(v.fValue)).fvalue = int64(value)
			return
		case DVT_Int32:
			if value <= math.MaxInt32 && value >= math.MinInt32{
				(*DxInt32Value)(unsafe.Pointer(v.fValue)).fvalue = int32(value)
				return
			}
		case DVT_Record:
			(*DxRecord)(unsafe.Pointer(v.fValue)).ClearValue()
		case DVT_Array:
			(*DxArray)(unsafe.Pointer(v.fValue)).ClearValue()
		}
	}
	var m DxIntValue
	m.fValueType = DVT_Int
	m.fvalue = value
	v.fValue = &m.DxBaseValue
}

func (v *DxValue)SetInt32Value(value int32)  {
	if v.fValue != nil{
		switch v.fValue.fValueType {
		case DVT_Int:
			(*DxIntValue)(unsafe.Pointer(v.fValue)).fvalue = int(value)
			return
		case DVT_Int32:
			(*DxInt32Value)(unsafe.Pointer(v.fValue)).fvalue = value
			return
		case DVT_Int64:
			(*DxInt64Value)(unsafe.Pointer(v.fValue)).fvalue = int64(value)
			return
		case DVT_Record:
			(*DxRecord)(unsafe.Pointer(v.fValue)).ClearValue()
		case DVT_Array:
			(*DxArray)(unsafe.Pointer(v.fValue)).ClearValue()
		}

	}
	var m DxInt32Value
	m.fValueType = DVT_Int32
	m.fvalue = value
	v.fValue = &m.DxBaseValue
}

func (v *DxValue)SetInt64Value(value int64)  {
	if v.fValue != nil{
		if v.fValue != nil{
			switch v.fValue.fValueType {
			case DVT_Int64:
				(*DxInt64Value)(unsafe.Pointer(v.fValue)).fvalue = value
				return
			case DVT_Int:
				if DxCommonLib.IsAmd64 || value <= math.MaxInt32 && value >= math.MinInt32{
					(*DxIntValue)(unsafe.Pointer(v.fValue)).fvalue = int(value)
					return
				}
			case DVT_Int32:
				if value <= math.MaxInt32 && value >= math.MinInt32{
					(*DxInt32Value)(unsafe.Pointer(v.fValue)).fvalue = int32(value)
					return
				}
			case DVT_Record:
				(*DxRecord)(unsafe.Pointer(v.fValue)).ClearValue()
			case DVT_Array:
				(*DxArray)(unsafe.Pointer(v.fValue)).ClearValue()
			}
		}
	}
	var m DxInt64Value
	m.fValueType = DVT_Int64
	m.fvalue = value
	v.fValue = &m.DxBaseValue
}

func (v *DxValue)CanParent()bool  {
	if v.fValue == nil{
		return false
	}
	return v.fValue.CanParent()
}

func (v *DxValue)AsInt()(int,error){
	if v.fValue == nil{
		return 0,nil
	}
	return v.fValue.AsInt()
}

func (v *DxValue)AsBool()(bool,error){
	if v.fValue == nil{
		return false,nil
	}
	return v.fValue.AsBool()
}

func (v *DxValue)AsInt32()(int32,error){
	if v.fValue == nil{
		return 0,nil
	}
	return v.fValue.AsInt32()
}

func (v *DxValue)AsInt64()(int64,error){
	if v.fValue == nil{
		return 0,nil
	}
	return v.fValue.AsInt64()
}

func (v *DxValue)AsArray()(*DxArray,error){
	if v.fValue == nil{
		return nil,nil
	}
	return v.fValue.AsArray()
}

func (v *DxValue)AsRecord()(*DxRecord,error){
	if v.fValue == nil{
		return nil,nil
	}
	return v.fValue.AsRecord()
}

func (v *DxValue)AsString()string{
	if v.fValue == nil{
		return ""
	}
	return v.fValue.AsString()
}

func (v *DxValue)AsFloat()(float32,error){
	if v.fValue == nil{
		return 0,nil
	}
	return v.fValue.AsFloat()
}

func (v *DxValue)AsDouble()(float64,error){
	if v.fValue == nil{
		return 0,nil
	}
	return v.fValue.AsDouble()
}

func (v *DxValue)AsBytes()([]byte,error){
	if v.fValue == nil{
		return nil,nil
	}
	return v.fValue.AsBytes(),nil
}

func (v *DxValue)ClearValue()  {
	if v.fValue != nil{
		v.fValue.ClearValue()
	}
	v.fValue = nil
}

func (v *DxValue)JsonParserFromByte(JsonByte []byte,ConvertEscape bool)(parserlen int, err error)  {
	for i := 0; i < len(JsonByte) ; i++  {
		if !IsSpace(JsonByte[i]){
			switch JsonByte[i] {
			case '{':
				var rec *DxRecord
				if v.fValue == nil || v.fValue.fValueType != DVT_Record{
					rec = &DxRecord{}
					rec.PathSplitChar = '.'
					rec.fValueType = DVT_Record
					rec.fRecords = make(map[string]*DxBaseValue,32)
				}else{
					rec = (*DxRecord)(unsafe.Pointer(v.fValue))
				}
				parserlen, err = rec.JsonParserFromByte(JsonByte[i:],ConvertEscape)
				if err == nil {
					v.fValue = &rec.DxBaseValue
				}
				return
			case '[':
				var arr *DxArray
				if v.fValue == nil || v.fValue.fValueType != DVT_Array{
					arr = &DxArray{}
					arr.fValueType = DVT_Array
				}else{
					arr = (*DxArray)(unsafe.Pointer(v.fValue))
				}
				parserlen, err = arr.JsonParserFromByte(JsonByte[i:],ConvertEscape)
				if err == nil {
					v.fValue = &arr.DxBaseValue
				}
				return
			default:
				return i,ErrInvalidateJson
			}
		}
	}
	return 0,ErrInvalidateJson
}

