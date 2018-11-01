/*
DxValue的Array数组对象
可以用来序列化反序列化Json,MsgPack等，并提供一系列的操作函数
Autor: 不得闲
QQ:75492895
 */
package DxValue

import (
	"unsafe"
	"bytes"
	"reflect"
	"github.com/suiyunonghen/DxCommonLib"
	"math"
	"strconv"
	"strings"
	"io/ioutil"
	"io"
	"bufio"
	"os"
	"time"
	"github.com/suiyunonghen/DxValue/Coders/DxMsgPack"
	"github.com/suiyunonghen/DxValue/Coders"
)

/******************************************************
*  DxArray
******************************************************/
type  DxArray		struct{
	DxBaseValue
	fValues		[]*DxBaseValue
}

func (arr *DxArray)ClearValue(clearInner bool)  {
	if clearInner{
		arr.fValues = nil
		return
	}
	if arr.fValues != nil{
		arr.fValues = arr.fValues[:0]
	}
}

func (arr *DxArray)getSize() int {
	result := 0
	if arr.fValues != nil{
		for _,v := range arr.fValues{
			result += v.Size()
		}
	}
	return result
}

func (arr *DxArray)TruncateArray(ArrLen int)  {
	if arr.fValues == nil{
		caplen := ArrLen
		if caplen < 128{
			caplen = 128
		}
		arr.fValues = make([]*DxBaseValue,ArrLen,caplen)
		return
	}
	al := len(arr.fValues)
	if al < ArrLen{
		mv := make([]*DxBaseValue,ArrLen - al)
		arr.fValues = append(arr.fValues,mv...)
	}else{
		arr.fValues = arr.fValues[:ArrLen]
	}
}


func (arr *DxArray)NewRecord(idx int)(rec *DxRecord)  {
	if idx < 0{
		if arr.fValues != nil{
			idx = 0
		}else{
			idx = len(arr.fValues)
		}
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		if arr.fValues[idx].fValueType == DVT_Record{
			rec = (*DxRecord)(unsafe.Pointer(arr.fValues[idx]))
			rec.ClearValue(false)
			rec.fParent = &arr.DxBaseValue
			return
		}
		arr.fValues[idx].ClearValue(true)
	}
	root := arr.NearestRecord()
	spchar := DefaultPathSplit
	if root != nil{
		spchar = root.PathSplitChar
	}
	rec = new(DxRecord)
	rec.PathSplitChar = spchar
	rec.fValueType = DVT_Record
	rec.fRecords = make(map[string]*DxBaseValue,32)
	rec.fParent = &arr.DxBaseValue
	arr.fValues[idx] = &rec.DxBaseValue
	return
}

func (arr *DxArray)NewIntRecord(idx int)(rec *DxIntKeyRecord)   {
	if idx < 0{
		if arr.fValues != nil{
			idx = 0
		}else{
			idx = len(arr.fValues)
		}
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		if arr.fValues[idx].fValueType == DVT_RecordIntKey{
			rec = (*DxIntKeyRecord)(unsafe.Pointer(arr.fValues[idx]))
			rec.ClearValue(false)
			rec.fParent = &arr.DxBaseValue
			return
		}
		arr.fValues[idx].ClearValue(true)
	}
	root := arr.NearestRecord()
	spchar := DefaultPathSplit
	if root != nil{
		spchar = root.PathSplitChar
	}
	rec = new(DxIntKeyRecord)
	rec.PathSplitChar = spchar
	rec.fValueType = DVT_RecordIntKey
	rec.fRecords = make(map[int64]*DxBaseValue,32)
	rec.fParent = &arr.DxBaseValue
	arr.fValues[idx] = &rec.DxBaseValue
	return
}

func (arr *DxArray)NewArray(idx int)(ararr *DxArray)  {
	if idx < 0{
		if arr.fValues != nil{
			idx = 0
		}else{
			idx = len(arr.fValues)
		}
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		if arr.fValues[idx].fValueType == DVT_Array{
			ararr = (*DxArray)(unsafe.Pointer(arr.fValues[idx]))
			ararr.ClearValue(false)
			ararr.fParent = &arr.DxBaseValue
			return
		}
		arr.fValues[idx].ClearValue(true)
	}
	ararr = new(DxArray)
	ararr.fValueType = DVT_Array
	ararr.fParent = &arr.DxBaseValue
	arr.fValues[idx] = &ararr.DxBaseValue
	return
}

func (arr *DxArray)ifNilInitArr2idx(idx int)  {
	if idx < 0{
		idx = 0
	}
	vlen := 0
	if arr.fValues == nil{
		caplen := 128
		if idx > caplen - 1{
			caplen = idx+1
		}
		arr.fValues = make([]*DxBaseValue,idx+1,caplen)
		vlen = idx+1
	}else{
		vlen = len(arr.fValues)
	}
	if idx > vlen - 1{
		mv := make([]*DxBaseValue,idx + 1 - vlen)
		arr.fValues = append(arr.fValues,mv...)
	}
}

func (arr *DxArray)Length()int  {
	if arr.fValues != nil{
		return len(arr.fValues)
	}
	return 0
}

func (arr *DxArray)VaueTypeByIndex(idx int)DxValueType  {
	if arr.fValues != nil && idx >= 0 && idx < len(arr.fValues) && arr.fValues[idx] != nil{
		return arr.fValues[idx].fValueType
	}
	return DVT_Null
}

func (arr *DxArray)SetNull(idx int)  {
	if idx < 0{
		return
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		arr.fValues[idx].fParent = nil
		arr.fValues[idx].ClearValue(true)
	}
	arr.fValues[idx] = nil
}

func (arr *DxArray)AsInt(idx int,defValue int)int  {
	if arr.fValues != nil && idx >= 0 && idx < len(arr.fValues) && arr.fValues[idx] != nil{
		value := arr.fValues[idx]
		switch value.fValueType {
		case DVT_Int: return (*DxIntValue)(unsafe.Pointer(value)).fvalue
		case DVT_Int32: return int((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Int64: return int((*DxInt64Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Bool:
			if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
				return 1
			}else{
				return 0
			}
		case DVT_Double,DVT_DateTime:return int((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return int((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		case DVT_String:
			v,err := strconv.ParseInt((*DxStringValue)(unsafe.Pointer(value)).fvalue,0,0)
			if err != nil{
				panic(err)
			}else{
				return int(v)
			}
		default:
			panic("can not convert Type to int")
		}
	}
	return defValue
}

func (arr *DxArray)AsInt32(idx int,defValue int32)int32  {
	if arr.fValues != nil && idx >= 0 && idx < len(arr.fValues) && arr.fValues[idx] != nil{
		value := arr.fValues[idx]
		switch value.fValueType {
		case DVT_Int: return int32((*DxIntValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Int32: return (*DxInt32Value)(unsafe.Pointer(value)).fvalue
		case DVT_Int64: return int32((*DxInt64Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Bool:
			if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
				return 1
			}else{
				return 0
			}
		case DVT_Double,DVT_DateTime:return int32((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return int32((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		case DVT_String:
			v,err := strconv.ParseInt((*DxStringValue)(unsafe.Pointer(value)).fvalue,0,0)
			if err != nil{
				panic(err)
			}else{
				return int32(v)
			}
		default:
			panic("can not convert Type to int")
		}
	}
	return defValue
}

func (arr *DxArray)AsInt64(idx int,defValue int64)int64  {
	if arr.fValues != nil && idx >= 0 && idx < len(arr.fValues) && arr.fValues[idx] != nil{
		value := arr.fValues[idx]
		switch value.fValueType {
		case DVT_Int: return int64((*DxIntValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Int32: return int64((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Int64: return (*DxInt64Value)(unsafe.Pointer(value)).fvalue
		case DVT_Bool:
			if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
				return 1
			}else{
				return 0
			}
		case DVT_Double,DVT_DateTime:return int64((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return int64((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		case DVT_String:
			v,err := strconv.ParseInt((*DxStringValue)(unsafe.Pointer(value)).fvalue,0,0)
			if err != nil{
				panic(err)
			}else{
				return v
			}
		default:
			panic("can not convert Type to int")
		}
	}
	return defValue
}


func (arr *DxArray)SetInt(idx,value int)  {
	if idx < 0{
		return
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil {
		switch arr.fValues[idx].fValueType {
		case DVT_Int:
			(*DxIntValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
			return
		case DVT_Int32:
			if value <= math.MaxInt32 && value >= math.MinInt32{
				(*DxInt32Value)(unsafe.Pointer(arr.fValues[idx])).fvalue = int32(value)
				return
			}
		case DVT_Int64:
			(*DxInt64Value)(unsafe.Pointer(arr.fValues[idx])).fvalue = int64(value)
			return
		default:
			arr.fValues[idx].ClearValue(true)
		}
	}

	dv := new(DxIntValue)
	dv.fValueType = DVT_Int
	dv.fvalue = value
	dv.fParent = &arr.DxBaseValue
	arr.fValues[idx] = &dv.DxBaseValue
}

func (arr *DxArray)SetInt32(idx int,value int32)  {
	if idx < 0{
		return
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil {
		switch arr.fValues[idx].fValueType {
		case DVT_Int:
			(*DxIntValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = int(value)
			return
		case DVT_Int32:
			(*DxInt32Value)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
			return
		case DVT_Int64:
			(*DxInt64Value)(unsafe.Pointer(arr.fValues[idx])).fvalue = int64(value)
			return
		default:
			arr.fValues[idx].ClearValue(true)
		}
	}
	dv := new(DxInt32Value)
	dv.fValueType = DVT_Int32
	dv.fvalue = value
	dv.fParent = &arr.DxBaseValue
	arr.fValues[idx] = &dv.DxBaseValue
}

func (arr *DxArray)SetInt64(idx int,value int64)  {
	if idx < 0{
		return
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		switch arr.fValues[idx].fValueType {
		case DVT_Int:
			if DxCommonLib.IsAmd64 || value <= math.MaxInt32 && value >= math.MinInt32{
				(*DxIntValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = int(value)
				return
			}
		case DVT_Int32:
			if value <= math.MaxInt32 && value >= math.MinInt32{
				(*DxInt32Value)(unsafe.Pointer(arr.fValues[idx])).fvalue = int32(value)
				return
			}
		case DVT_Int64:
			(*DxInt64Value)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
			return
		default:
			arr.fValues[idx].ClearValue(true)
		}
	}
	if value <= math.MaxInt32 && value >= math.MinInt32{
		dv := new(DxInt32Value)
		dv.fValueType = DVT_Int32
		dv.fvalue = int32(value)
		dv.fParent = &arr.DxBaseValue
		arr.fValues[idx] = &dv.DxBaseValue
	}else{
		dv := new(DxInt64Value)
		dv.fValueType = DVT_Int64
		dv.fvalue = value
		dv.fParent = &arr.DxBaseValue
		arr.fValues[idx] = &dv.DxBaseValue
	}
}

func (arr *DxArray)SetBool(idx int,value bool)  {
	if idx < 0{
		return
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_Bool{
		(*DxBoolValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
	}else{
		dv := new(DxBoolValue)
		dv.fValueType = DVT_Bool
		dv.fvalue = value
		if arr.fValues[idx] != nil{
			arr.fValues[idx].ClearValue(true)
		}
		dv.fParent = &arr.DxBaseValue
		arr.fValues[idx] = &dv.DxBaseValue
	}
}

func (arr *DxArray)AsBool(idx int,defValue bool)bool  {
	if arr.fValues != nil && idx >= 0 && idx < len(arr.fValues) && arr.fValues[idx] != nil{
		value := arr.fValues[idx]
		switch value.fValueType {
		case DVT_Int: return (*DxIntValue)(unsafe.Pointer(value)).fvalue != 0
		case DVT_Int32: return (*DxInt32Value)(unsafe.Pointer(value)).fvalue != 0
		case DVT_Int64: return (*DxInt64Value)(unsafe.Pointer(value)).fvalue != 0
		case DVT_Bool: return bool((*DxBoolValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Double,DVT_DateTime:return float64((*DxDoubleValue)(unsafe.Pointer(value)).fvalue) != 0
		case DVT_Float:return float32((*DxFloatValue)(unsafe.Pointer(value)).fvalue) != 0
		case DVT_String:
			return strings.ToUpper((*DxStringValue)(unsafe.Pointer(value)).fvalue) == "TRUE"
		default:
			panic("can not convert Type to Bool")
		}
	}
	return defValue
}

func (arr *DxArray)AsString(idx int,defValue string)string  {
	if arr.fValues != nil && idx >= 0 && idx < len(arr.fValues) && arr.fValues[idx] != nil{
		return arr.fValues[idx].ToString()
	}
	return defValue
}

func (arr *DxArray)SetString(idx int,value string)  {
	if idx < 0{
		return
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		switch arr.fValues[idx].fValueType {
		case DVT_String:
			(*DxStringValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
		case DVT_DateTime:
			if jt := DxCommonLib.ParserJsonTime(value);jt >= 0{
				(*DxDoubleValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = float64(jt)
				return
			}
			t,err := time.Parse("2006-01-02T15:04:05Z",value)
			if err == nil{
				(*DxDoubleValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = float64(DxCommonLib.Time2DelphiTime(&t))
				return
			}
			t,err = time.Parse("2006-01-02 15:04:05",value)
			if err == nil{
				(*DxDoubleValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = float64(DxCommonLib.Time2DelphiTime(&t))
				return
			}
			t,err = time.Parse("2006/01/02 15:04:05",value)
			if err == nil{
				(*DxDoubleValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = float64(DxCommonLib.Time2DelphiTime(&t))
				return
			}
		case DVT_Int:
			if iv,err := strconv.Atoi(value);err == nil{
				(*DxIntValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = iv
				return
			}
		case DVT_Int32:
			if iv,err := strconv.Atoi(value);err == nil{
				(*DxInt32Value)(unsafe.Pointer(arr.fValues[idx])).fvalue = int32(iv)
				return
			}
		case DVT_Int64:
			if iv,err := strconv.ParseInt(value,10,64);err == nil{
				(*DxInt64Value)(unsafe.Pointer(arr.fValues[idx])).fvalue = iv
				return
			}
		case DVT_Float:
			if iv,err := strconv.ParseFloat(value,32);err == nil{
				(*DxFloatValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = float32(iv)
				return
			}
		case DVT_Double:
			if iv,err := strconv.ParseFloat(value,64);err == nil{
				(*DxDoubleValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = iv
				return
			}
		}
		arr.fValues[idx].ClearValue(true)
	}
	dv := new(DxStringValue)
	dv.fValueType = DVT_String
	dv.fvalue = value
	if arr.fValues[idx] != nil{
		arr.fValues[idx].ClearValue(true)
	}
	dv.fParent = &arr.DxBaseValue
	arr.fValues[idx] = &dv.DxBaseValue
}


func (arr *DxArray)SetFloat(idx int,value float32)  {
	if idx < 0{
		return
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil {
		switch arr.fValueType {
		case DVT_Float:
			(*DxFloatValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
			return
		case DVT_Double,DVT_DateTime:
			(*DxDoubleValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = float64(value)
			return
		default:
			arr.fValues[idx].ClearValue(true)
		}
	}
	dv := new(DxFloatValue)
	dv.fValueType = DVT_Float
	dv.fvalue = value
	dv.fParent = &arr.DxBaseValue
	arr.fValues[idx] = &dv.DxBaseValue
}

func (arr *DxArray)AsFloat(idx int,defValue float32)float32  {
	if arr.fValues != nil && idx >= 0 && idx < len(arr.fValues) && arr.fValues[idx] != nil{
		value := arr.fValues[idx]
		switch value.fValueType {
		case DVT_Int: return float32((*DxIntValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Int32: return float32((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Int64: return float32((*DxInt64Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Bool:
			if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
				return 1
			}
			return 0
		case DVT_Double,DVT_DateTime:return float32((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return (*DxFloatValue)(unsafe.Pointer(value)).fvalue
		default:
			panic("can not convert Type to Float")
		}
	}
	return defValue
}

func (arr *DxArray)SetDouble(idx int,value float64)  {
	if idx < 0{
		return
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		switch arr.fValueType {
		case DVT_Float:
			if value <= math.MaxFloat32 && value >= math.MinInt32 {
				(*DxFloatValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = float32(value)
				return
			}
		case DVT_Double,DVT_DateTime:
			(*DxDoubleValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
			return
		default:
			arr.fValues[idx].ClearValue(true)
		}
	}
	dv := new(DxDoubleValue)
	dv.fValueType = DVT_Double
	dv.fvalue = value
	dv.fParent = &arr.DxBaseValue
	arr.fValues[idx] = &dv.DxBaseValue
}

func (arr *DxArray)SetDateTime(idx int,t DxCommonLib.TDateTime)  {
	if idx < 0{
		return
	}
	arr.ifNilInitArr2idx(idx)
	value := float64(t)
	if arr.fValues[idx] != nil{
		switch arr.fValueType {
		case DVT_Float:
			if value <= math.MaxFloat32 && t >= math.MinInt32 {
				(*DxFloatValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = float32(value)
				return
			}
		case DVT_Double,DVT_DateTime:
			(*DxDoubleValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
			return
		default:
			arr.fValues[idx].ClearValue(true)
		}
	}
	dv := new(DxDoubleValue)
	dv.fValueType = DVT_DateTime
	dv.fvalue = value
	dv.fParent = &arr.DxBaseValue
	arr.fValues[idx] = &dv.DxBaseValue
}

func (arr *DxArray)AsDouble(idx int,defValue float64)float64  {
	if arr.fValues != nil && idx >= 0 && idx < len(arr.fValues) && arr.fValues[idx] != nil{
		value := arr.fValues[idx]
		switch value.fValueType {
		case DVT_Int: return float64((*DxIntValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Int32: return float64((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Int64: return float64((*DxInt64Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Bool:
			if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
				return 1
			}
			return 0
		case DVT_Double,DVT_DateTime:return (*DxDoubleValue)(unsafe.Pointer(value)).fvalue
		case DVT_Float:return float64((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		default:
			panic("can not convert Type to Float")
		}
	}
	return defValue
}

func (arr *DxArray)AsBaseValue(idx int)*DxBaseValue  {
	if arr.fValues != nil && idx >= 0 && idx < len(arr.fValues){
		return arr.fValues[idx]
	}
	return nil
}

func (arr *DxArray)AsDateTime(idx int,defValue DxCommonLib.TDateTime)DxCommonLib.TDateTime  {
	if arr.fValues != nil && idx >= 0 && idx < len(arr.fValues) && arr.fValues[idx] != nil{
		value := arr.fValues[idx]
		switch value.fValueType {
		case DVT_Int: return DxCommonLib.TDateTime((*DxIntValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Int32: return DxCommonLib.TDateTime((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Int64: return DxCommonLib.TDateTime((*DxInt64Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Bool:
			if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
				return 1
			}
			return 0
		case DVT_Double,DVT_DateTime:return DxCommonLib.TDateTime((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return DxCommonLib.TDateTime((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		default:
			panic("can not convert Type to Float")
		}
	}
	return defValue
}


func (arr *DxArray)SetArray(idx int,value *DxArray)  {
	if idx < 0{
		return
	}
	if value != nil && value.fParent != nil {
		panic("Must Set A Single Array(no Parent)")
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		arr.fValues[idx].ClearValue(true)
		arr.fValues[idx].fParent = nil
	}
	arr.fValues[idx] = &value.DxBaseValue
	arr.fValues[idx].fParent = &arr.DxBaseValue
}

func (arr *DxArray)AsArray(idx int)(*DxArray)  {
	if arr.fValues != nil && idx >= 0 && idx < len(arr.fValues) && arr.fValues[idx] != nil{
		if arr.fValues[idx].fValueType == DVT_Array{
			return (*DxArray)(unsafe.Pointer(arr.fValues[idx]))
		}
		panic("not Array Value")
	}
	return nil
}


func (arr *DxArray)SetRecord(idx int,value *DxRecord)  {
	if idx < 0{
		return
	}
	if value != nil && value.fParent != nil {
		panic("Must Set A Single Record(no Parent)")
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		arr.fValues[idx].ClearValue(true)
		arr.fValues[idx].fParent = nil
	}
	if value!=nil{
		root := arr.NearestRecord()
		spchar := DefaultPathSplit
		if root != nil{
			spchar = root.PathSplitChar
		}
		value.PathSplitChar = spchar
		arr.fValues[idx] = &value.DxBaseValue
		arr.fValues[idx].fParent = &arr.DxBaseValue
	}else{
		arr.fValues[idx] = nil
	}
}

func (arr *DxArray)SetIntRecord(idx int,value *DxIntKeyRecord)  {
	if idx < 0{
		return
	}
	if value != nil && value.fParent != nil {
		panic("Must Set A Single Record(no Parent)")
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		arr.fValues[idx].ClearValue(true)
		arr.fValues[idx].fParent = nil
	}
	if value!=nil{
		root := arr.NearestRecord()
		spchar := DefaultPathSplit
		if root != nil{
			spchar = root.PathSplitChar
		}
		value.PathSplitChar = spchar
		arr.fValues[idx] = &value.DxBaseValue
		arr.fValues[idx].fParent = &arr.DxBaseValue
	}else{
		arr.fValues[idx] = nil
	}
}

func (arr *DxArray)AsRecord(idx int)(*DxRecord)  {
	if arr.fValues != nil && idx >= 0 && idx < len(arr.fValues) && arr.fValues[idx] != nil{
		if arr.fValues[idx].fValueType == DVT_Record{
			return (*DxRecord)(unsafe.Pointer(arr.fValues[idx]))
		}
		panic("not Record Value")
	}
	return nil
}

func (arr *DxArray)AsIntRecord(idx int)(*DxIntKeyRecord)  {
	if arr.fValues != nil && idx >= 0 && idx < len(arr.fValues) && arr.fValues[idx] != nil{
		if arr.fValues[idx].fValueType == DVT_RecordIntKey{
			return (*DxIntKeyRecord)(unsafe.Pointer(arr.fValues[idx]))
		}
		panic("not IntKeyRecord Value")
	}
	return nil
}

func (arr *DxArray)SetBinary(idx int,bt []byte)  {
	if idx < 0{
		return
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_Binary{
		(*DxBinaryValue)(unsafe.Pointer(arr.fValues[idx])).SetBinary(bt,true)
	}else{
		dv := new(DxBinaryValue)
		dv.fValueType = DVT_Binary
		dv.SetBinary(bt,true)
		if arr.fValues[idx] != nil{
			arr.fValues[idx].ClearValue(true)
		}
		dv.fParent = &arr.DxBaseValue
		arr.fValues[idx] = &dv.DxBaseValue
	}
}

func (arr *DxArray)SetBaseValue(idx int,v *DxBaseValue)  {
	if v.fValueType == DVT_Unknown{
		panic("UnKnown Value")
	}
	arr.ifNilInitArr2idx(idx)
	if v == nil{
		arr.SetNull(idx)
		return
	}
	if arr.fValues[idx] != nil{
		arr.fValues[idx].ClearValue(true)
		arr.fValues[idx].fParent = nil
	}
	arr.fValues[idx] = v
	v.fParent = &arr.DxBaseValue
}

func (arr *DxArray)Append(value ...interface{})  {
	arrlen := arr.Length()
	for i := 0;i<len(value);i++{
		arr.SetValue(i+arrlen,value[i])
	}
}

func (arr *DxArray)SetValue(idx int,value interface{})  {
	if idx < 0{
		return
	}
	arr.ifNilInitArr2idx(idx)
	if value == nil{
		arr.SetNull(idx)
		return 
	}
	switch value := value.(type) {
	case int:
		arr.SetInt(idx, value)
	case int32:
		arr.SetInt32(idx, value)
	case int64:
		arr.SetInt64(idx, value)
	case int8:
		arr.SetInt(idx, int(value))
	case uint8:
		arr.SetInt(idx, int(value))
	case int16:
		arr.SetInt(idx, int(value))
	case uint16:
		arr.SetInt(idx, int(value))
	case uint32:
		arr.SetInt(idx, int(value))
	case *int:
		arr.SetInt(idx, *value)
	case *int32:
		arr.SetInt32(idx, *value)
	case *int64:
		arr.SetInt64(idx, *value)
	case *int8:
		arr.SetInt(idx, int(*value))
	case *uint8:
		arr.SetInt(idx, int(*value))
	case *int16:
		arr.SetInt(idx, int(*value))
	case *uint16:
		arr.SetInt(idx, int(*value))
	case *uint32:
		arr.SetInt(idx, int(*value))
	case string:
		arr.SetString(idx, value)
	case []byte:
		arr.SetBinary(idx, value)
	case *[]byte:
		arr.SetBinary(idx, *value)
	case bool:
		arr.SetBool(idx, value)
	case *bool:
		arr.SetBool(idx, *value)
	case *string:
		arr.SetString(idx, *value)
	case float32:
		arr.SetFloat(idx, value)
	case float64:
		arr.SetDouble(idx, value)
	case *float32:
		arr.SetFloat(idx, *value)
	case *float64:
		arr.SetDouble(idx, *value)
	case *DxRecord:
		arr.SetRecord(idx, value)
	case DxRecord:
		arr.SetRecord(idx, &value)
	case DxArray:
		arr.SetArray(idx,&value)
	case *DxArray:
		arr.SetArray(idx,value)
	case *DxBaseValue: arr.SetBaseValue(idx,value)
	case DxInt64Value: arr.SetInt64(idx,value.fvalue)
	case *DxInt64Value: arr.SetInt64(idx,value.fvalue)
	case DxInt32Value: arr.SetInt32(idx,value.fvalue)
	case *DxInt32Value: arr.SetInt32(idx,value.fvalue)
	case DxFloatValue: arr.SetFloat(idx,value.fvalue)
	case *DxFloatValue: arr.SetFloat(idx,value.fvalue)
	case DxDoubleValue: arr.SetDouble(idx,value.fvalue)
	case *DxDoubleValue: arr.SetDouble(idx,value.fvalue)
	case DxBoolValue: arr.SetBool(idx,value.fvalue)
	case *DxBoolValue: arr.SetBool(idx,value.fvalue)
	case DxIntValue: arr.SetInt(idx,value.fvalue)
	case *DxIntValue: arr.SetInt(idx,value.fvalue)
	case DxStringValue: arr.SetString(idx,value.fvalue)
	case *DxStringValue: arr.SetString(idx,value.fvalue)
	case DxBinaryValue:  arr.SetBinary(idx,value.Bytes())
	case *DxBinaryValue:  arr.SetBinary(idx,value.Bytes())
	case time.Time: arr.SetDateTime(idx,DxCommonLib.Time2DelphiTime(&value))
	case *time.Time: arr.SetDateTime(idx,DxCommonLib.Time2DelphiTime(value))
	default:
		reflectv := reflect.ValueOf(value)
		rv := getRealValue(&reflectv)
		if rv == nil{
			arr.SetNull(idx)
			return
		}
		switch rv.Kind(){
		case reflect.Struct:
			rec := arr.NewRecord(idx)
			rtype := rv.Type()
			for i := 0;i < rtype.NumField();i++{
				sfield := rtype.Field(i)
				fv := rv.Field(i)
				fieldvalue := getRealValue(&fv)
				if fieldvalue != nil{
					switch fieldvalue.Kind() {
					case reflect.Int,reflect.Uint32:
						rec.SetInt(sfield.Name,int(fieldvalue.Int()))
					case reflect.Bool:
						rec.SetBool(sfield.Name,fieldvalue.Bool())
					case reflect.Int64:
						rec.SetInt64(sfield.Name,fieldvalue.Int())
					case reflect.Int32,reflect.Int8,reflect.Int16,reflect.Uint8,reflect.Uint16:
						rec.SetInt32(sfield.Name,int32(fieldvalue.Int()))
					case reflect.Float32:
						rec.SetFloat(sfield.Name,float32(fieldvalue.Float()))
					case reflect.Float64:
						rec.SetDouble(sfield.Name,fieldvalue.Float())
					case reflect.String:
						rec.SetString(sfield.Name,fieldvalue.String())
					default:
						if fieldvalue.CanInterface(){
							rec.SetValue(sfield.Name,fieldvalue.Interface())
						}
					}
				}
			}
		case reflect.Map:
			mapkeys := rv.MapKeys()
			if len(mapkeys) == 0{
				return
			}
			kv := mapkeys[0]
			var rbase *DxBaseValue
			switch getBaseType(kv.Type()) {
			case reflect.String:
				rbase = &arr.NewRecord(idx).DxBaseValue
			case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
				rbase = &arr.NewIntRecord(idx).DxBaseValue
			default:
				panic("Invalidate Record Key,Can Only Int or String")
			}
			rvalue := rv.MapIndex(mapkeys[0])
			//获得Value类型
			valueKind := getBaseType(rvalue.Type())
			for _,kv = range mapkeys{
				rvalue = rv.MapIndex(kv)
				prvalue := getRealValue(&rvalue)
				if prvalue != nil{
					switch valueKind {
					case reflect.Int,reflect.Uint32:
						if rbase.fValueType == DVT_Record{
							(*DxRecord)(unsafe.Pointer(rbase)).SetInt(kv.String(),int(prvalue.Int()))
						}else{
							(*DxIntKeyRecord)(unsafe.Pointer(rbase)).SetInt(kv.Int(),int(prvalue.Int()))
						}
					case reflect.Bool:
						if rbase.fValueType == DVT_Record{
							(*DxRecord)(unsafe.Pointer(rbase)).SetBool(kv.String(),prvalue.Bool())
						}else {
							(*DxIntKeyRecord)(unsafe.Pointer(rbase)).SetBool(kv.Int(),prvalue.Bool())
						}
					case reflect.Int64:
						if rbase.fValueType == DVT_Record{
							(*DxRecord)(unsafe.Pointer(rbase)).SetInt64(kv.String(),prvalue.Int())
						}else{
							(*DxIntKeyRecord)(unsafe.Pointer(rbase)).SetInt64(kv.Int(),prvalue.Int())
						}
					case reflect.Int32,reflect.Int8,reflect.Int16,reflect.Uint8,reflect.Uint16:
						if rbase.fValueType == DVT_Record{
							(*DxRecord)(unsafe.Pointer(rbase)).SetInt32(kv.String(),int32(prvalue.Int()))
						}else{
							(*DxIntKeyRecord)(unsafe.Pointer(rbase)).SetInt32(kv.Int(),int32(prvalue.Int()))
						}
					case reflect.Float32:
						if rbase.fValueType == DVT_Record{
							(*DxRecord)(unsafe.Pointer(rbase)).SetFloat(kv.String(),float32(prvalue.Float()))
						}else{
							(*DxIntKeyRecord)(unsafe.Pointer(rbase)).SetFloat(kv.Int(),float32(prvalue.Float()))
						}
					case reflect.Float64:
						if rbase.fValueType == DVT_Record {
							(*DxRecord)(unsafe.Pointer(rbase)).SetDouble(kv.String(), prvalue.Float())
						}else{
							(*DxIntKeyRecord)(unsafe.Pointer(rbase)).SetDouble(kv.Int(), prvalue.Float())
						}
					case reflect.String:
						if rbase.fValueType == DVT_Record {
							(*DxRecord)(unsafe.Pointer(rbase)).SetString(kv.String(), prvalue.String())
						}else{
							(*DxIntKeyRecord)(unsafe.Pointer(rbase)).SetString(kv.Int(), prvalue.String())
						}
					default:
						if prvalue.CanInterface(){
							if rbase.fValueType == DVT_Record {
								(*DxRecord)(unsafe.Pointer(rbase)).SetValue(kv.String(), prvalue.Interface())
							}else{
								(*DxIntKeyRecord)(unsafe.Pointer(rbase)).SetValue(kv.Int(), prvalue.Interface())
							}
						}
					}
				}
			}
		case reflect.Array,reflect.Slice:
			carr := arr.NewArray(idx)
			vlen := rv.Len()
			for i := 0;i< vlen;i++{
				av := rv.Index(i)
				arrvalue := getRealValue(&av)
				switch arrvalue.Kind() {
				case reflect.Int,reflect.Uint32:
					carr.SetInt(i,int(arrvalue.Int()))
				case reflect.Bool:
					carr.SetBool(i,arrvalue.Bool())
				case reflect.Int64:
					carr.SetInt64(i,arrvalue.Int())
				case reflect.Int32,reflect.Int8,reflect.Int16,reflect.Uint8,reflect.Uint16:
					carr.SetInt32(i,int32(arrvalue.Int()))
				case reflect.Float32:
					carr.SetFloat(i,float32(arrvalue.Float()))
				case reflect.Float64:
					carr.SetDouble(i,arrvalue.Float())
				case reflect.String:
					carr.SetString(i,arrvalue.String())
				default:
					if arrvalue.CanInterface(){
						carr.SetValue(i,arrvalue.Interface())
					}
				}
			}
		}
	}
}


func (arr *DxArray)SetExtValue(idx int,extbt []byte)  {
	if idx < 0{
		return
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		if arr.fValues[idx].fValueType == DVT_Ext{
			arr.fValues[idx].SetExtValue(extbt)
			return
		}
		arr.fValues[idx].ClearValue(true)
	}
	var bv DxExtValue
	bv.fdata = extbt
	if extbt != nil && len(extbt) > 0{
		bv.fExtType = extbt[0]
	}
	bv.fValueType = DVT_Ext
	arr.fValues[idx] = &bv.DxBaseValue
	arr.fValues[idx].fParent = &arr.DxBaseValue
}

func (arr *DxArray)AsExtValue(idx int)(*DxExtValue)  {
	if arr.fValues != nil && idx >= 0 && idx < len(arr.fValues) && arr.fValues[idx] != nil{
		if arr.fValues[idx].fValueType == DVT_Ext{
			return (*DxExtValue)(unsafe.Pointer(arr.fValues[idx]))
		}
		panic("not ExtType Value")
	}
	return nil
}

func (arr *DxArray)Bytes()[]byte  {
	var buf bytes.Buffer
	buf.WriteByte('[')
	if arr.fValues != nil{
		isFirst := true
		for i := 0;i<len(arr.fValues);i++{
			av := arr.fValues[i]
			if !isFirst{
				buf.WriteByte(',')
			}else{
				isFirst = false
			}
			if av == nil{
				buf.WriteString("null")
			}else{
				if av.fValueType == DVT_String || av.fValueType == DVT_Binary{
					buf.WriteByte('"')
				}
				buf.WriteString(av.ToString())
				if av.fValueType == DVT_String || av.fValueType == DVT_Binary{
					buf.WriteByte('"')
				}
			}
		}
	}
	buf.WriteByte(']')
	return buf.Bytes()
}

func (arr *DxArray)BytesWithSort()[]byte  {
	var buf bytes.Buffer
	buf.WriteByte('[')
	if arr.fValues != nil{
		isFirst := true
		for i := 0;i<len(arr.fValues);i++{
			av := arr.fValues[i]
			if !isFirst{
				buf.WriteByte(',')
			}else{
				isFirst = false
			}
			if av == nil{
				buf.WriteString("null")
			}else{
				switch av.fValueType {
				case DVT_String,DVT_Binary:
					buf.WriteByte('"')
					buf.WriteString(av.ToString())
					buf.WriteByte('"')
				case DVT_Record:
					buf.Write((*DxRecord)(unsafe.Pointer(av)).BytesWithSort())
				case DVT_RecordIntKey:
					buf.Write((*DxIntKeyRecord)(unsafe.Pointer(av)).BytesWithSort())
				case DVT_Array:
					buf.Write((*DxArray)(unsafe.Pointer(av)).BytesWithSort())
				default:
					buf.WriteString(av.ToString())
				}
			}
		}
	}
	buf.WriteByte(']')
	return buf.Bytes()
}

func (arr *DxArray)Delete(idx int)  {
	if arr.fValues != nil{
		if idx >= 0 && idx < len(arr.fValues){
			if arr.fValues[idx] != nil{
				arr.fValues[idx].ClearValue(true)
				arr.fValues[idx].fParent = nil
			}
			arr.fValues = append(arr.fValues[:idx],arr.fValues[idx+1:]...)
		}
	}
}

func (arr *DxArray)parserValue(idx int, b []byte,ConvertEscape,structRest bool)(parserlen int, err error)  {
	record := arr.NearestRecord()
	spchar := DefaultPathSplit
	if record != nil{
		spchar = record.PathSplitChar
	}
	i := 0
	btlen := len(b)
	validCharIndex := -1
	for i < btlen{
		if !IsSpace(b[i]){
			switch b[i] {
			case '[':
				narr := NewArray()
				if parserlen,err = narr.JsonParserFromByte(b[i:],ConvertEscape,structRest);err!=nil{
					return
				}
				arr.SetArray(idx,narr)
				parserlen += 1+i
				return
			case '{':
				rec := NewRecord()
				rec.PathSplitChar = spchar
				if parserlen,err = rec.JsonParserFromByte(b[i:],ConvertEscape,structRest);err != nil{
					return
				}
				arr.SetRecord(idx,rec)
				parserlen += 1+i
				return
			case ',',']':
				//bvalue := bytes.Trim(b[:i]," \r\n\t")
				//获取有效的字符开始的位置
				bvalue := b[:validCharIndex+1]
				if len(bvalue) == 0{
					return i,ErrInvalidateJson
				}
				if bytes.IndexByte(bvalue,'.') > -1{
					if vf,err := strconv.ParseFloat(DxCommonLib.FastByte2String(bvalue),64);err!=nil{
						return i,ErrInvalidateJson
					}else{
						arr.SetDouble(idx,vf)
					}
				}else{
					st := DxCommonLib.FastByte2String(bvalue)
					if st == "true" || strings.ToUpper(st) == "TRUE"{
						arr.SetBool(idx,true)
					}else if st == "false" || strings.ToUpper(st) == "FALSE"{
						arr.SetBool(idx,false)
					}else if st == "null" || strings.ToUpper(st) == "NULL"{
						arr.SetNull(idx)
					}else{
						if vf,err := strconv.Atoi(st);err!=nil{
							return i,ErrInvalidateJson
						}else{
							arr.SetInt(idx,vf)
						}
					}
				}
				return i,nil
			case '"':
				plen := bytes.IndexByte(b[i+1:btlen],'"')
				if plen > -1{
					st := DxCommonLib.FastByte2String(b[i+1:plen+i+1])
					if ConvertEscape{
						if jt := DxCommonLib.ParserJsonTime(st);jt>=0{
							arr.SetDateTime(idx,jt)
							return i,ErrInvalidateJson
						}
					}
					arr.SetString(idx,st)
					return plen + i + 2,nil
				}
				return i,ErrInvalidateJson
			default:
				validCharIndex = i
			}
		}
		i++
	}
	return btlen,ErrInvalidateJson
}

func (arr *DxArray)JsonParserFromByte(JsonByte []byte,ConvertEscape,structRest bool)(parserlen int, err error)  {
	btlen := len(JsonByte)
	i := 0
	idx := 0
	arrStart := false
	valuestart := false
	for i < btlen{
		if IsSpace(JsonByte[i]){
			i++
			continue
		}
		if !arrStart && JsonByte[i] != '['{
			return 0,ErrInvalidateJson
		}
		switch JsonByte[i]{
		case '[':
			if arrStart{
				if parserlen,err = arr.parserValue(idx,JsonByte[i:],ConvertEscape,structRest);err!=nil{
					return parserlen + i,err
				}
				idx++
				i += parserlen
				continue
			}
			arrStart = true
			valuestart = true
		case ']':
			return i,nil
		case ',','}':
			valuestart = true
		default:
			if valuestart {
				valuestart = false
				if parserlen,err = arr.parserValue(idx,JsonByte[i:],ConvertEscape,structRest);err!=nil{
					return parserlen + i,err
				}
				idx++
				i += parserlen
				continue
			}else{
				return i,ErrInvalidateJson
			}
		}
		i++
	}
	return btlen,ErrInvalidateJson
}

func (arr *DxArray)ToString()string  {
	return DxCommonLib.FastByte2String(arr.Bytes())
}

func (arr *DxArray)LoadJsonFile(fileName string,ConvertEscape,structRest bool)error  {
	databytes, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		return err
	}
	if databytes[0] == 0xEF && databytes[1] == 0xBB && databytes[2] == 0xBF{//BOM
		databytes = databytes[3:]
	}
	_,err = arr.JsonParserFromByte(databytes,ConvertEscape,structRest)
	return err
}

func (arr *DxArray)SaveJsonFile(fileName string,BOMFile bool)(err error){
	if file,err := os.OpenFile(fileName,os.O_CREATE | os.O_TRUNC,0644);err == nil{
		defer file.Close()
		if BOMFile{
			file.Write([]byte{0xEF,0xBB,0xBF})
		}
		return arr.SaveJsonWriter(file)
	}else{
		return err
	}
}

func (arr *DxArray)SaveJsonWriter(w io.Writer)(err error) {
	writer := bufio.NewWriter(w)
	err = writer.WriteByte('[')
	if err != nil{
		return
	}
	if arr.fValues != nil{
		isFirst := true
		for i := 0;i<len(arr.fValues);i++{
			av := arr.fValues[i]
			if !isFirst{
				err = writer.WriteByte(',')
				if err != nil{
					return
				}
			}else{
				isFirst = false
			}
			if av == nil{
				_,err = writer.WriteString("null")
			}else{
				if av.fValueType == DVT_String || av.fValueType == DVT_Binary{
					if err = writer.WriteByte('"');err!=nil{
						return
					}
				}
				_,err = writer.WriteString(av.ToString())
				if err == nil && (av.fValueType == DVT_String || av.fValueType == DVT_Binary){
					err = writer.WriteByte('"')
				}
			}
			if err != nil{
				return
			}
		}
	}
	writer.WriteByte(']')
	return writer.Flush()
}

func (arr *DxArray)LoadJsonReader(reader io.Reader)error  {
	return nil
}


func (arr *DxArray)LoadMsgPackReader(reader io.Reader)error  {
	return NewDecoder(reader).Decode(&arr.DxBaseValue)
}

func (arr *DxArray)LoadMsgPackFile(fileName string)error  {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return NewDecoder(f).Decode(&arr.DxBaseValue)
}

func (arr *DxArray)SaveMsgPackFile(fileName string)error  {
	if file,err := os.OpenFile(fileName,os.O_CREATE | os.O_TRUNC,0644);err == nil{
		defer file.Close()
		return NewEncoder(file).EncodeArray(arr)
	}else{
		return err
	}
}

func (arr *DxArray)Encode(valuecoder Coders.Encoder) error{
	var err error
	switch valuecoder.Name() {
	case "msgpack":
		if msgpacker,ok := valuecoder.(*DxMsgPackEncoder);ok{
			return msgpacker.EncodeArray(arr)
		}
		encoder := valuecoder.(*DxMsgPack.MsgPackEncoder)
		arlen := uint(arr.Length())
		switch {
		case arlen < 16: //1001XXXX|    N objects
			err = encoder.WriteByte(byte(DxMsgPack.CodeFixedArrayLow) | byte(arlen))
		case arlen <= DxMsgPack.Max_map16_len:  //0xdc  |YYYYYYYY|YYYYYYYY|    N objects
			encoder.WriteUint16(uint16(arlen),DxMsgPack.CodeArray16)
		default:
			if arlen > DxMsgPack.Max_map32_len{
				arlen = DxMsgPack.Max_map32_len
			}
			encoder.WriteUint32(uint32(arlen),DxMsgPack.CodeArray32)
		}

		for i := uint(0);i <= arlen - 1;i++{
			vbase := arr.AsBaseValue(int(i))
			if vbase == nil{
				err = encoder.WriteByte(0xc0) //null
			}else{
				err = vbase.Encode(encoder)
			}
			if err != nil{
				return err
			}
		}
		return err
	}
	return nil
}

func NewArray()*DxArray  {
	result := new(DxArray)
	result.fValueType = DVT_Array
	return result
}