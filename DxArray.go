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
)

/******************************************************
*  DxArray
******************************************************/
type  DxArray		struct{
	DxBaseValue
	fValues		[]*DxBaseValue
}

func (arr *DxArray)ClearValue()  {
	if arr.fValues != nil{
		arr.fValues = arr.fValues[:0]
	}
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
	if idx == -1{
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
			rec.ClearValue()
			rec.fParent = &arr.DxBaseValue
			return
		}
		arr.fValues[idx].fParent = nil
	}
	rec = new(DxRecord)
	rec.fValueType = DVT_Record
	rec.fRecords = make(map[string]*DxBaseValue,32)
	rec.fParent = &arr.DxBaseValue
	arr.fValues = append(arr.fValues,&rec.DxBaseValue)
	return
}

func (arr *DxArray)NewArray(idx int)(ararr *DxArray)  {
	if idx == -1{
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
			ararr.ClearValue()
			ararr.fParent = &arr.DxBaseValue
			return
		}
		arr.fValues[idx].fParent = nil
	}
	ararr = new(DxArray)
	ararr.fValueType = DVT_Array
	ararr.fParent = &arr.DxBaseValue
	arr.fValues = append(arr.fValues,&ararr.DxBaseValue)
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
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		arr.fValues[idx].fParent = nil
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
		case DVT_Double:return int((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return int((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
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
		case DVT_Double:return int32((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return int32((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
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
		case DVT_Double:return int64((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return int64((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		default:
			panic("can not convert Type to int")
		}
	}
	return defValue
}

func (arr *DxArray)SetInt(idx,value int)  {
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
		}
	}

	dv := new(DxIntValue)
	dv.fValueType = DVT_Int
	dv.fvalue = value
	if arr.fValues[idx] != nil{
		arr.fValues[idx].fParent = nil
	}
	dv.fParent = &arr.DxBaseValue
	arr.fValues[idx] = &dv.DxBaseValue
}

func (arr *DxArray)SetInt32(idx int,value int32)  {
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
		}
	}
	dv := new(DxInt32Value)
	dv.fValueType = DVT_Int32
	dv.fvalue = value
	if arr.fValues[idx] != nil{
		arr.fValues[idx].fParent = nil
	}
	dv.fParent = &arr.DxBaseValue
	arr.fValues[idx] = &dv.DxBaseValue
}

func (arr *DxArray)SetInt64(idx int,value int64)  {
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
		}
	}
	dv := new(DxInt64Value)
	dv.fValueType = DVT_Int64
	dv.fvalue = value
	if arr.fValues[idx] != nil{
		arr.fValues[idx].fParent = nil
	}
	dv.fParent = &arr.DxBaseValue
	arr.fValues[idx] = &dv.DxBaseValue
}

func (arr *DxArray)SetBool(idx int,value bool)  {
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_Bool{
		(*DxBoolValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
	}else{
		dv := new(DxBoolValue)
		dv.fValueType = DVT_Bool
		dv.fvalue = value
		if arr.fValues[idx] != nil{
			arr.fValues[idx].fParent = nil
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
		case DVT_Double:return float64((*DxDoubleValue)(unsafe.Pointer(value)).fvalue) != 0
		case DVT_Float:return float32((*DxFloatValue)(unsafe.Pointer(value)).fvalue) != 0
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
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_String{
		(*DxStringValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
	}else{
		dv := new(DxStringValue)
		dv.fValueType = DVT_String
		dv.fvalue = value
		if arr.fValues[idx] != nil{
			arr.fValues[idx].fParent = nil
		}
		dv.fParent = &arr.DxBaseValue
		arr.fValues[idx] = &dv.DxBaseValue
	}
}

func (arr *DxArray)SetFloat(idx int,value float32)  {
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil {
		switch arr.fValueType {
		case DVT_Float:
			(*DxFloatValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
			return
		case DVT_Double:
			(*DxDoubleValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = float64(value)
			return
		}
	}
	dv := new(DxFloatValue)
	dv.fValueType = DVT_Float
	dv.fvalue = value
	if arr.fValues[idx] != nil{
		arr.fValues[idx].fParent = nil
	}
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
		case DVT_Double:return float32((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return (*DxFloatValue)(unsafe.Pointer(value)).fvalue
		default:
			panic("can not convert Type to Float")
		}
	}
	return defValue
}

func (arr *DxArray)SetDouble(idx int,value float64)  {
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		switch arr.fValueType {
		case DVT_Float:
			if value <= math.MaxFloat32 && value >= math.MinInt32 {
				(*DxFloatValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = float32(value)
				return
			}
		case DVT_Double:
			(*DxDoubleValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
			return
		}
	}
	dv := new(DxDoubleValue)
	dv.fValueType = DVT_Double
	dv.fvalue = value
	if arr.fValues[idx] != nil{
		arr.fValues[idx].fParent = nil
	}
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
		case DVT_Double:return (*DxDoubleValue)(unsafe.Pointer(value)).fvalue
		case DVT_Float:return float64((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		default:
			panic("can not convert Type to Float")
		}
	}
	return defValue
}

func (arr *DxArray)SetArray(idx int,value *DxArray)  {
	if value != nil && value.fParent != nil {
		panic("Must Set A Single Array(no Parent)")
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		if arr.fValues[idx].fValueType == DVT_Array{
			arr := (*DxArray)(unsafe.Pointer(arr.fValues[idx]))
			arr.fParent = nil
			*arr=*value
			arr.fParent = &arr.DxBaseValue
			return
		}
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
	if value != nil && value.fParent != nil {
		panic("Must Set A Single Record(no Parent)")
	}
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil{
		if arr.fValues[idx].fValueType == DVT_Record{
			rec := (*DxRecord)(unsafe.Pointer(arr.fValues[idx]))
			rec.fParent = nil
			*rec = *value
			rec.fParent = &arr.DxBaseValue
			return
		}
		arr.fValues[idx].fParent = nil
	}
	arr.fValues[idx] = &value.DxBaseValue
	arr.fValues[idx].fParent = &arr.DxBaseValue
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

func (arr *DxArray)SetBinary(idx int,bt []byte)  {
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_Binary{
		(*DxBinaryValue)(unsafe.Pointer(arr.fValues[idx])).SetBinary(bt,true)
	}else{
		dv := new(DxBinaryValue)
		dv.fValueType = DVT_Binary
		dv.SetBinary(bt,true)
		if arr.fValues[idx] != nil{
			arr.fValues[idx].fParent = nil
		}
		dv.fParent = &arr.DxBaseValue
		arr.fValues[idx] = &dv.DxBaseValue
	}
}

func (arr *DxArray)SetValue(idx int,value interface{})  {
	arr.ifNilInitArr2idx(idx)
	if value == nil{
		arr.fValues[idx] = nil
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
			rec := arr.NewRecord(idx)
			mapkeys := rv.MapKeys()
			if len(mapkeys) == 0{
				return
			}
			kv := mapkeys[0]
			if getBaseType(kv.Type()) != reflect.String{
				panic("Invalidate Record Key")
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
						rec.SetInt(kv.String(),int(prvalue.Int()))
					case reflect.Bool:
						rec.SetBool(kv.String(),prvalue.Bool())
					case reflect.Int64:
						rec.SetInt64(kv.String(),prvalue.Int())
					case reflect.Int32,reflect.Int8,reflect.Int16,reflect.Uint8,reflect.Uint16:
						rec.SetInt32(kv.String(),int32(prvalue.Int()))
					case reflect.Float32:
						rec.SetFloat(kv.String(),float32(prvalue.Float()))
					case reflect.Float64:
						rec.SetDouble(kv.String(),prvalue.Float())
					case reflect.String:
						rec.SetString(kv.String(),prvalue.String())
					default:
						if prvalue.CanInterface(){
							rec.SetValue(kv.String(),prvalue.Interface())
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

func (arr *DxArray)parserValue(idx int, b []byte)(parserlen int, err error)  {
	i := 0;
	btlen := len(b)
	for i < btlen{
		if !IsSpace(b[i]){
			switch b[i] {
			case '[':
				narr := NewArray()
				if parserlen,err = narr.JsonParserFromByte(b[i:]);err!=nil{
					return
				}
				arr.SetArray(idx,narr)
				parserlen += 1
				return
			case '{':
				rec := NewRecord()
				if parserlen,err = rec.JsonParserFromByte(b[i:]);err != nil{
					return
				}
				arr.SetRecord(idx,rec)
				parserlen += 1
				return

			case ',',']':
				bvalue := bytes.Trim(b[:i]," \r\n\t")
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
					arr.SetString(idx,st)
					return plen + i + 2,nil
				}
				return i,ErrInvalidateJson
			}
		}
		i++
	}
	return btlen,ErrInvalidateJson
}

func (arr *DxArray)JsonParserFromByte(JsonByte []byte)(parserlen int, err error)  {
	btlen := len(JsonByte)
	i := 0
	idx := 0
	arrStart := false
	valuestart := false
	for i < btlen{
		if !arrStart && JsonByte[i] != '[' && !IsSpace(JsonByte[i]){
			return 0,ErrInvalidateJson
		}
		switch JsonByte[i]{
		case '[':
			if arrStart{
				if parserlen,err = arr.parserValue(idx,JsonByte[i:]);err!=nil{
					return parserlen + i,err
				}
				idx++
				i += parserlen
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
				if parserlen,err = arr.parserValue(idx,JsonByte[i:]);err!=nil{
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

func NewArray()*DxArray  {
	result := new(DxArray)
	result.fValueType = DVT_Array
	return result
}