/*
DxValue的Record记录集对象
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
*  DxRecord
******************************************************/
type  DxRecord		struct{
	DxBaseValue
	fRecords		map[string]*DxBaseValue
}


func (r *DxRecord)ClearValue()  {
	r.fRecords = make(map[string]*DxBaseValue,32)
}


func (r *DxRecord)NewRecord(keyName string)(rec *DxRecord)  {
	if value,ok := r.fRecords[keyName];ok && value != nil{
		if value.fValueType == DVT_Record{
			rec = (*DxRecord)(unsafe.Pointer(value))
			rec.ClearValue()
			return
		}
	}
	rec = new(DxRecord)
	rec.fValueType = DVT_Record
	rec.fRecords = make(map[string]*DxBaseValue,32)
	r.fRecords[keyName] = &rec.DxBaseValue
	return
}

func (r *DxRecord)NewArray(keyName string)(arr *DxArray)  {
	if value,ok := r.fRecords[keyName];ok && value != nil{
		if value.fValueType == DVT_Array{
			arr = (*DxArray)(unsafe.Pointer(value))
			arr.ClearValue()
			return
		}
	}
	arr = new(DxArray)
	arr.fValueType = DVT_Array
	r.fRecords[keyName] = &arr.DxBaseValue
	return
}

func (r *DxRecord)SetInt(KeyName string,v int)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Int{
			(*DxIntValue)(unsafe.Pointer(value)).fvalue = v
			return
		}
	}
	var m DxIntValue
	m.fvalue = v
	m.fValueType = DVT_Int
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetInt32(KeyName string,v int32)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Int32{
			(*DxInt32Value)(unsafe.Pointer(value)).fvalue = v
			return
		}
	}
	var m DxInt32Value
	m.fvalue = v
	m.fValueType = DVT_Int32
	r.fRecords[KeyName] = &m.DxBaseValue
}


func (r *DxRecord)SetInt64(KeyName string,v int64)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Int64{
			(*DxInt64Value)(unsafe.Pointer(value)).fvalue = v
			return
		}
	}
	var m DxInt64Value
	m.fvalue = v
	m.fValueType = DVT_Int64
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetBool(KeyName string,v bool)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Bool{
			(*DxBoolValue)(unsafe.Pointer(value)).fvalue = v
			return
		}
	}
	var m DxBoolValue
	m.fvalue = v
	m.fValueType = DVT_Bool
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetFloat(KeyName string,v float32)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Float{
			(*DxFloatValue)(unsafe.Pointer(value)).fvalue = v
			return
		}
	}
	var m DxFloatValue
	m.fvalue = v
	m.fValueType = DVT_Float
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetDouble(KeyName string,v float64)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Double{
			(*DxDoubleValue)(unsafe.Pointer(value)).fvalue = v
			return
		}
	}
	var m DxDoubleValue
	m.fvalue = v
	m.fValueType = DVT_Double
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetString(KeyName string,v string)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Double{
			(*DxStringValue)(unsafe.Pointer(value)).fvalue = v
			return
		}
	}
	var m DxStringValue
	m.fvalue = v
	m.fValueType = DVT_String
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetBinary(KeyName string,v []byte,reWrite bool)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Binary{
			bv := (*DxBinaryValue)(unsafe.Pointer(value))
			if reWrite{
				bv.SetBinary(v,false)
			}else{
				bv.Append(v)
			}
			return
		}
	}
	var m DxBinaryValue
	m.Append(v)
	m.fValueType = DVT_Binary
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)AsBytes(keyName string)[]byte  {
	if value,ok := r.fRecords[keyName];ok && value != nil{
		return value.Bytes()
	}
	return nil
}

func getBaseType(vt reflect.Type)reflect.Kind  {
	if vt.Kind() == reflect.Ptr{
		return getBaseType(vt.Elem())
	}
	return vt.Kind()
}

func getRealValue(v *reflect.Value)*reflect.Value  {
	if !v.IsValid(){
		return nil
	}
	if v.Kind() == reflect.Ptr{
		if !v.IsNil(){
			va := v.Elem()
			return getRealValue(&va)
		}else{
			return nil
		}
	}
	return v
}

func (r *DxRecord)SetRecordValue(keyName string,v *DxRecord,isbyref bool)  {
	if value,ok := r.fRecords[keyName];ok && value != nil{
		if value.fValueType == DVT_Record{
			nrec := (*DxRecord)(unsafe.Pointer(value))
			if isbyref{
				nrec.fRecords = v.fRecords
			}else if v.fRecords == nil || len(v.fRecords) == 0{
				nrec.fRecords = nil
			}else{
				nrec.ClearValue()
				for k,v := range v.fRecords{
					nrec.fRecords[k] = v
				}
			}
			return
		}
	}
	if isbyref{
		r.fRecords[keyName] = &v.DxBaseValue
	}else{
		nrec := r.NewRecord(keyName)
		nrec.ClearValue()
		for k,v := range v.fRecords{
			nrec.fRecords[k] = v
		}
	}
}

func (r *DxRecord)SetArray(KeyName string,v *DxArray,copyarr bool)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Array{
			arr := (*DxArray)(unsafe.Pointer(value))
			if copyarr{
				if v.fValues == nil ||  len(v.fValues) == 0{
					arr.ClearValue()
				}else{
					if arr.fValues == nil{
						arr.fValues = make([]*DxBaseValue,len(v.fValues))
					}
					copy(arr.fValues,v.fValues)
				}
			}else{
				*arr=*v
			}
			return
		}
	}
	if !copyarr{
		r.fRecords[KeyName] = &v.DxBaseValue
	}else{
		arr := r.NewArray(KeyName)
		if v.fValues == nil ||  len(v.fValues) == 0{
			arr.ClearValue()
		}else if v.fValues!= nil && len(v.fValues)!=0{
			arr.fValues = make([]*DxBaseValue,len(v.fValues))
			copy(arr.fValues,v.fValues)
		}
	}
}

func (r *DxRecord)SetValue(keyName string,v interface{})  {
	if v == nil{
		r.fRecords[keyName] = nil
		return
	}
	switch value := v.(type) {
	case int: r.SetInt(keyName,value)
	case int32: r.SetInt32(keyName,value)
	case int64: r.SetInt64(keyName,value)
	case int8: r.SetInt(keyName,int(value))
	case uint8: r.SetInt(keyName,int(value))
	case int16: r.SetInt(keyName,int(value))
	case uint16: r.SetInt(keyName,int(value))
	case uint32: r.SetInt(keyName,int(value))
	case *int: r.SetInt(keyName,*value)
	case *int32: r.SetInt32(keyName,*value)
	case *int64: r.SetInt64(keyName,*value)
	case *int8: r.SetInt(keyName,int(*value))
	case *uint8: r.SetInt(keyName,int(*value))
	case *int16: r.SetInt(keyName,int(*value))
	case *uint16: r.SetInt(keyName,int(*value))
	case *uint32: r.SetInt(keyName,int(*value))
	case string: r.SetString(keyName,value)
	case []byte: r.SetBinary(keyName,value,true)
	case *[]byte: r.SetBinary(keyName,*value,true)
	case bool: r.SetBool(keyName,value)
	case *bool: r.SetBool(keyName,*value)
	case *string: r.SetString(keyName,*value)
	case float32: r.SetFloat(keyName,value)
	case float64: r.SetDouble(keyName,value)
	case *float32: r.SetFloat(keyName,*value)
	case *float64: r.SetDouble(keyName,*value)
	case *DxRecord: r.SetRecordValue(keyName,value,true)
	case DxRecord: r.SetRecordValue(keyName,&value,true)
	default:
		reflectv := reflect.ValueOf(v)
		rv := getRealValue(&reflectv)
		if rv == nil{
			return
		}
		switch rv.Kind(){
		case reflect.Struct:
			rec := r.NewRecord(keyName)
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
			rec := r.NewRecord(keyName)
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


func (r *DxRecord)KeyValueType(KeyName string)DxValueType  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		return value.fValueType
	}
	return DVT_Null
}

func (r *DxRecord)AsInt32(KeyName string)int32  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_Int: return int32((*DxIntValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Int32: return int32((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
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
			panic("can not convert Type to int32")
		}
	}
	return 0
}

func (r *DxRecord)AsInt(KeyName string)int  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_Int: return int((*DxIntValue)(unsafe.Pointer(value)).fvalue)
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
	return 0
}

func (r *DxRecord)AsInt64(KeyName string)int64  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_Int: return int64((*DxIntValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Int32: return int64((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Int64: return int64((*DxInt64Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Bool:
			if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
				return 1
			}else{
				return 0
			}
		case DVT_Double:return int64((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return int64((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		default:
			panic("can not convert Type to int64")
		}
	}
	return 0
}

func (r *DxRecord)AsBool(KeyName string)bool  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
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
	return false
}


func (r *DxRecord)AsFloat(KeyName string)float32  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
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
		case DVT_Float:return float32((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		default:
			panic("can not convert Type to Float")
		}
	}
	return 0
}


func (r *DxRecord)AsDouble(KeyName string)float64  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_Int: return float64((*DxIntValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Int32: return float64((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Int64: return float64((*DxInt64Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Bool:
			if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
				return 1
			}
			return 0
		case DVT_Double:return float64((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return float64((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		default:
			panic("can not convert Type to Double")
		}
	}
	return 0
}

func (r *DxRecord)AsString(KeyName string)string  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		return value.ToString()
	}
	return ""
}

func (r *DxRecord)AsRecord(KeyName string)*DxRecord  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Record{
			return (*DxRecord)(unsafe.Pointer(value))
		}
	}
	return nil
}

func (r *DxRecord)ToString()string  {
	var buffer bytes.Buffer
	buffer.WriteByte('{')
	isFirst := true
	for k,v := range r.fRecords{
		if !isFirst{
			buffer.WriteString(`,"`)
		}else{
			isFirst = false
			buffer.WriteByte('"')
		}
		buffer.WriteString(k)
		buffer.WriteString(`":`)
		if v != nil{
			vt := v.fValueType
			if vt == DVT_String{
				buffer.WriteByte('"')
			}
			buffer.WriteString(v.ToString())
			if vt == DVT_String{
				buffer.WriteByte('"')
			}
		}else{
			buffer.WriteString("NULL")
		}
	}
	buffer.WriteByte('}')
	return string(buffer.Bytes())
}

func NewRecord()*DxRecord  {
	result := new(DxRecord)
	result.fValueType = DVT_Record
	result.fRecords = make(map[string]*DxBaseValue,32)
	return result
}
