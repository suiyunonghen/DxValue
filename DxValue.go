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


func (v *DxValue)ClearValue()  {
	v.fValue = nil
}

