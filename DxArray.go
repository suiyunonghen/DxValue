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
	"fmt"
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
	rec = new(DxRecord)
	rec.fValueType = DVT_Record
	rec.fRecords = make(map[string]*DxBaseValue,20)
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
	arr = new(DxArray)
	ararr.fValueType = DVT_Array
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

func (arr *DxArray)SetNull(idx int)  {
	arr.ifNilInitArr2idx(idx)
	arr.fValues[idx] = nil
}

func (arr *DxArray)SetInt(idx,value int)  {
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_Int{
		(*DxIntValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
	}else{
		dv := new(DxIntValue)
		dv.fValueType = DVT_Int
		dv.fvalue = value
		arr.fValues[idx] = &dv.DxBaseValue
	}
}

func (arr *DxArray)SetInt32(idx int,value int32)  {
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_Int32{
		(*DxInt32Value)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
	}else{
		dv := new(DxInt32Value)
		dv.fValueType = DVT_Int32
		dv.fvalue = value
		arr.fValues[idx] = &dv.DxBaseValue
	}
}

func (arr *DxArray)SetInt64(idx int,value int64)  {
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_Int64{
		(*DxInt64Value)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
	}else{
		dv := new(DxInt64Value)
		dv.fValueType = DVT_Int64
		dv.fvalue = value
		arr.fValues[idx] = &dv.DxBaseValue
	}
}

func (arr *DxArray)SetBool(idx int,value bool)  {
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_Bool{
		(*DxBoolValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
	}else{
		dv := new(DxBoolValue)
		dv.fValueType = DVT_Bool
		dv.fvalue = value
		arr.fValues[idx] = &dv.DxBaseValue
	}
}

func (arr *DxArray)SetString(idx int,value string)  {
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_String{
		(*DxStringValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
	}else{
		dv := new(DxStringValue)
		dv.fValueType = DVT_String
		dv.fvalue = value
		arr.fValues[idx] = &dv.DxBaseValue
	}
}

func (arr *DxArray)SetFloat(idx int,value float32)  {
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_Float{
		(*DxFloatValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
	}else{
		dv := new(DxFloatValue)
		dv.fValueType = DVT_Float
		dv.fvalue = value
		arr.fValues[idx] = &dv.DxBaseValue
	}
}

func (arr *DxArray)SetDouble(idx int,value float64)  {
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_Double{
		(*DxDoubleValue)(unsafe.Pointer(arr.fValues[idx])).fvalue = value
	}else{
		dv := new(DxDoubleValue)
		dv.fValueType = DVT_Double
		dv.fvalue = value
		arr.fValues[idx] = &dv.DxBaseValue
	}
}

func (arr *DxArray)SetArray(idx int,value *DxArray,isbyref bool)  {
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_Array{
		arr := (*DxArray)(unsafe.Pointer(arr.fValues[idx]))
		if !isbyref{
			if value.fValues == nil ||  len(value.fValues) == 0{
				arr.ClearValue()
			}else{
				if arr.fValues == nil{
					arr.fValues = make([]*DxBaseValue,len(value.fValues))
				}
				copy(arr.fValues,value.fValues)
			}
		}else{
			*arr=*value
		}
	}else if isbyref{
		arr.fValues[idx] = &value.DxBaseValue
	} else{
		dv := new(DxArray)
		dv.fValueType = DVT_Array
		if value.fValues!= nil && len(value.fValues)!=0{
			arr.fValues = make([]*DxBaseValue,len(value.fValues))
			copy(arr.fValues,value.fValues)
		}
	}
}


func (arr *DxArray)SetRecord(idx int,value *DxRecord,isbyref bool)  {
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_Record{
		if isbyref{
			*(*DxRecord)(unsafe.Pointer(arr.fValues[idx])) = *value
		}else{
			rec := (*DxRecord)(unsafe.Pointer(arr.fValues[idx]))
			rec.ClearValue()
			if value.fRecords != nil && len(value.fRecords) > 0{
				for k,v := range value.fRecords{
					rec.fRecords[k] = v
				}
			}
		}
	}else if isbyref{
		arr.fValues[idx] = &value.DxBaseValue
	}else {
		dv := new(DxRecord)
		dv.fValueType = DVT_Record
		dv.fRecords = make(map[string]*DxBaseValue,32)
		if value.fRecords != nil && len(value.fRecords) > 0{
			for k,v := range value.fRecords{
				dv.fRecords[k] = v
			}
		}
		arr.fValues[idx] = &dv.DxBaseValue
	}
}

func (arr *DxArray)SetBinary(idx int,bt []byte)  {
	arr.ifNilInitArr2idx(idx)
	if arr.fValues[idx] != nil && arr.fValues[idx].fValueType == DVT_Binary{
		(*DxBinaryValue)(unsafe.Pointer(arr.fValues[idx])).SetBinary(bt,true)
	}else{
		dv := new(DxBinaryValue)
		dv.fValueType = DVT_Binary
		dv.SetBinary(bt,true)
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
		arr.SetRecord(idx, value,true)
	case DxRecord:
		arr.SetRecord(idx, &value,true)
	default:
		reflectv := reflect.ValueOf(value)
		rv := getRealValue(&reflectv)
		if rv == nil{
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
			fmt.Println(valueKind)
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
		case reflect.Array:
		}
	}
}

func (arr *DxArray)ToString()string  {
	var buf bytes.Buffer
	buf.WriteByte('[')
	if arr.fValues != nil{
		for i := 0;i<len(arr.fValues);i++{

		}
	}
	buf.WriteByte(']')
	return string(buf.Bytes())
}

func NewArray()*DxArray  {
	result := new(DxArray)
	result.fValueType = DVT_Array
	return result
}