/*
DxValue的Record记录集对象，整数键
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
	"strings"
	"math"
	"strconv"
	"io/ioutil"
	"io"
	"bufio"
	"os"
	"time"
	"github.com/suiyunonghen/DxValue/Coders/DxMsgPack"
	"github.com/suiyunonghen/DxValue/Coders"
)

/******************************************************
*  DxIntKeyRecord
******************************************************/
type(
	DxIntKeyRecord		struct{
		DxBaseRecord
		fRecords		map[int64]*DxBaseValue
	}
)

func (r *DxIntKeyRecord)ClearValue(clearInner bool)  {
	if r.fRecords != nil{
		for _,v := range r.fRecords{
			if v != nil{
				v.ClearValue(true)
				v.fParent = nil
			}
		}
	}
	if clearInner{
		r.fRecords = nil
		return
	}
	if r.fRecords == nil || len(r.fRecords) > 0{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
}

func (r *DxIntKeyRecord)getSize() int {
	result := 0
	if r.fRecords != nil {
		for _,v := range r.fRecords{
			result += v.Size() + 8
		}
	}
	return result
}

func (r *DxIntKeyRecord)splitPathFields(charrune rune) bool {
	if r.PathSplitChar == 0{
		r.PathSplitChar = DefaultPathSplit
	}
	return charrune == rune(r.PathSplitChar)
}

func (r *DxIntKeyRecord)NewIntRecord(key int64)(rec *DxIntKeyRecord)  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			if value.fValueType == DVT_RecordIntKey{
				rec = (*DxIntKeyRecord)(unsafe.Pointer(value))
				rec.ClearValue(false)
				rec.fParent = &r.DxBaseValue
				return
			}
			value.ClearValue(true)
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	rec = new(DxIntKeyRecord)
	rec.fValueType = DVT_RecordIntKey
	rec.PathSplitChar = r.PathSplitChar
	rec.fRecords = make(map[int64]*DxBaseValue,32)
	r.fRecords[key] = &rec.DxBaseValue
	rec.fParent = &r.DxBaseValue
	return
}

func (r *DxIntKeyRecord)NewRecord(key int64,ExistsReset bool)(rec *DxRecord)  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			if value.fValueType == DVT_Record{
				rec = (*DxRecord)(unsafe.Pointer(value))
				if ExistsReset{
					rec.ClearValue(false)
					rec.fParent = &r.DxBaseValue
				}
				return
			}
			value.ClearValue(true)
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	rec = new(DxRecord)
	rec.fValueType = DVT_Record
	rec.PathSplitChar = r.PathSplitChar
	rec.fRecords = make(map[string]*DxBaseValue,32)
	r.fRecords[key] = &rec.DxBaseValue
	rec.fParent = &r.DxBaseValue
	return
}

func (r *DxIntKeyRecord)Find(key int64)*DxBaseValue  {
	if r.fRecords != nil{
		if v,ok := r.fRecords[key];ok{
			return v
		}
	}
	return nil
}


func (r *DxIntKeyRecord)NewArray(key int64)(arr *DxArray)  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			if value.fValueType == DVT_Array{
				arr = (*DxArray)(unsafe.Pointer(value))
				arr.ClearValue(false)
				arr.fParent = &r.DxBaseValue
				return
			}
			value.ClearValue(true)
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	arr = new(DxArray)
	arr.fValueType = DVT_Array
	arr.fParent = &r.DxBaseValue
	r.fRecords[key] = &arr.DxBaseValue
	return
}

func (r *DxIntKeyRecord)SetInt(key int64,v int)  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			switch value.fValueType {
			case DVT_Int:
				(*DxIntValue)(unsafe.Pointer(value)).fvalue = v
				return
			case DVT_Int32:
				if v <= math.MaxInt32 && v >= math.MinInt32{
					(*DxInt32Value)(unsafe.Pointer(value)).fvalue = int32(v)
					return
				}
			case DVT_Int64:
				(*DxInt64Value)(unsafe.Pointer(value)).fvalue = int64(v)
				return
			}
			value.ClearValue(true)
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	var m DxIntValue
	m.fvalue = v
	m.fValueType = DVT_Int
	m.fParent = &r.DxBaseValue
	r.fRecords[key] = &m.DxBaseValue
}

func (r *DxIntKeyRecord)SetInt32(key int64,v int32)  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			switch value.fValueType {
			case DVT_Int:
				(*DxIntValue)(unsafe.Pointer(value)).fvalue = int(v)
				return
			case DVT_Int32:
				(*DxInt32Value)(unsafe.Pointer(value)).fvalue = v
				return
			case DVT_Int64:
				(*DxInt64Value)(unsafe.Pointer(value)).fvalue = int64(v)
				return
			}
			value.ClearValue(true)
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	var m DxInt32Value
	m.fvalue = v
	m.fValueType = DVT_Int32
	m.fParent = &r.DxBaseValue
	r.fRecords[key] = &m.DxBaseValue
}


func (r *DxIntKeyRecord)SetInt64(key int64,v int64)  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			switch value.fValueType {
			case DVT_Int:
				if DxCommonLib.IsAmd64 || v <= math.MaxInt32 && v >= math.MinInt32{
					(*DxIntValue)(unsafe.Pointer(value)).fvalue = int(v)
					return
				}
			case DVT_Int32:
				if v <= math.MaxInt32 && v >= math.MinInt32{
					(*DxInt32Value)(unsafe.Pointer(value)).fvalue = int32(v)
					return
				}
				return
			case DVT_Int64:
				(*DxInt64Value)(unsafe.Pointer(value)).fvalue = v
				return
			}
			value.ClearValue(true)
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	var m DxInt64Value
	m.fvalue = v
	m.fValueType = DVT_Int64
	m.fParent = &r.DxBaseValue
	r.fRecords[key] = &m.DxBaseValue
}

func (r *DxIntKeyRecord)SetBool(key int64,v bool)  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			if value.fValueType == DVT_Bool {
				(*DxBoolValue)(unsafe.Pointer(value)).fvalue = v
				return
			}
			value.ClearValue(true)
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	var m DxBoolValue
	m.fvalue = v
	m.fParent = &r.DxBaseValue
	m.fValueType = DVT_Bool
	r.fRecords[key] = &m.DxBaseValue
}

func (r *DxIntKeyRecord)SetNull(key int64)  {
	if r.fRecords != nil{
		if v,ok := r.fRecords[key];ok{
			v.ClearValue(true)
			v.fParent = nil
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	r.fRecords[key] = nil
}



func (r *DxIntKeyRecord)SetFloat(key int64,v float32)  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			if value.fValueType == DVT_Float{
				(*DxFloatValue)(unsafe.Pointer(value)).fvalue = v
				return
			}else if value.fValueType == DVT_Double || value.fValueType == DVT_DateTime{
				(*DxDoubleValue)(unsafe.Pointer(value)).fvalue = float64(v)
				return
			}
			value.ClearValue(true)
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	var m DxFloatValue
	m.fvalue = v
	m.fParent = &r.DxBaseValue
	m.fValueType = DVT_Float
	r.fRecords[key] = &m.DxBaseValue
}

func (r *DxIntKeyRecord)SetDouble(key int64,v float64)  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			if value.fValueType == DVT_Double || value.fValueType == DVT_DateTime{
				(*DxDoubleValue)(unsafe.Pointer(value)).fvalue = v
				return
			}else if value.fValueType == DVT_Float{
				if v <= math.MaxFloat32 && v >= math.MinInt32{
					(*DxFloatValue)(unsafe.Pointer(value)).fvalue = float32(v)
					return
				}
			}
			value.ClearValue(true)
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	var m DxDoubleValue
	m.fvalue = v
	m.fParent = &r.DxBaseValue
	m.fValueType = DVT_Double
	r.fRecords[key] = &m.DxBaseValue
}

func (r *DxIntKeyRecord)SetDateTime(key int64,v DxCommonLib.TDateTime)  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			if value.fValueType == DVT_Double || value.fValueType == DVT_DateTime{
				(*DxDoubleValue)(unsafe.Pointer(value)).fvalue = float64(v)
				(*DxDoubleValue)(unsafe.Pointer(value)).fValueType = DVT_DateTime
				return
			}else if value.fValueType == DVT_Float{
				if v <= math.MaxFloat32 && v >= math.MinInt32{
					(*DxFloatValue)(unsafe.Pointer(value)).fvalue = float32(v)
					return
				}
			}
			value.ClearValue(true)
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	var m DxDoubleValue
	m.fvalue = float64(v)
	m.fParent = &r.DxBaseValue
	m.fValueType = DVT_DateTime
	r.fRecords[key] = &m.DxBaseValue
}

func (r *DxIntKeyRecord)SetGoTime(key int64,v time.Time)  {
	r.SetDateTime(key,DxCommonLib.Time2DelphiTime(v))
}

func (r *DxIntKeyRecord)SetString(key int64,v string)  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			if value.fValueType == DVT_String{
				(*DxStringValue)(unsafe.Pointer(value)).fvalue = v
				return
			}
			value.ClearValue(true)
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	var m DxStringValue
	m.fvalue = v
	m.fValueType = DVT_String
	r.fRecords[key] = &m.DxBaseValue
}

func (r *DxIntKeyRecord)SetBinary(key int64,v []byte,reWrite bool)  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			if value.fValueType == DVT_Binary{
				bv := (*DxBinaryValue)(unsafe.Pointer(value))
				if reWrite{
					bv.SetBinary(v,false)
				}else{
					bv.Append(v)
				}
				return
			}
			value.ClearValue(true)
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	var m DxBinaryValue
	m.Append(v)
	m.fParent = &r.DxBaseValue
	m.fValueType = DVT_Binary
	r.fRecords[key] = &m.DxBaseValue
}

func (r *DxIntKeyRecord)AsBytes(key int64)[]byte  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			bt,_ := value.AsBytes()
			return bt
		}
	}
	return nil
}

func (r *DxIntKeyRecord)BytesWithSort(escapJsonStr bool)[]byte  {
	buffer := bytes.NewBuffer(make([]byte,0,512))
	buffer.WriteByte('{')
	if r.fRecords != nil{
		keys := make([]int64,len(r.fRecords))
		idx := 0
		for k,_ := range r.fRecords{
			keys[idx] = k
			idx++
		}
		for i := 0;i<idx;i++{
			if i == 0{
				buffer.WriteByte('"')
			}else{
				buffer.WriteString(`,"`)
			}
			buffer.WriteString(strconv.Itoa(int(keys[i])))
			buffer.WriteString(`":`)
			v := r.fRecords[keys[i]]
			if v != nil{
				vt := v.fValueType
				if vt == DVT_String || vt == DVT_Binary || vt == DVT_DateTime || vt == DVT_Ext{
					buffer.WriteByte('"')
				}
				switch vt {
				case DVT_DateTime:
					buffer.WriteString("/Date(")
					buffer.WriteString(strconv.Itoa(int(DxCommonLib.TDateTime((*DxDoubleValue)(unsafe.Pointer(v)).fvalue).ToTime().Unix())*1000))
					buffer.WriteString(")/")
				case DVT_RecordIntKey:
					buffer.Write((*DxIntKeyRecord)(unsafe.Pointer(v)).BytesWithSort(escapJsonStr))
				case DVT_Record:
					buffer.Write((*DxRecord)(unsafe.Pointer(v)).BytesWithSort(escapJsonStr))
				case DVT_Array:
					buffer.Write((*DxArray)(unsafe.Pointer(v)).BytesWithSort(escapJsonStr))
				case DVT_String:
					if escapJsonStr{
						buffer.WriteString(DxCommonLib.EscapeJsonStr((*DxStringValue)(unsafe.Pointer(v)).fvalue))
					}else{
						buffer.WriteString((*DxStringValue)(unsafe.Pointer(v)).fvalue)
					}
				default:
					buffer.WriteString(v.ToString())
				}
				if vt == DVT_String || vt == DVT_Binary{
					buffer.WriteByte('"')
				}
			}else{
				buffer.WriteString("null")
			}
		}
	}
	buffer.WriteByte('}')
	return buffer.Bytes()
}

func (r *DxIntKeyRecord)Bytes(escapJsonStr bool)[]byte  {
	buffer := bytes.NewBuffer(make([]byte,0,512))
	buffer.WriteByte('{')
	if r.fRecords != nil{
		isFirst := true
		for k,v := range r.fRecords{
			if !isFirst{
				buffer.WriteString(`,"`)
			}else{
				isFirst = false
				buffer.WriteByte('"')
			}
			buffer.WriteString(strconv.Itoa(int(k)))
			buffer.WriteString(`":`)
			if v != nil{
				vt := v.fValueType
				if vt == DVT_String || vt == DVT_Binary || vt == DVT_DateTime || vt == DVT_Ext{
					buffer.WriteByte('"')
				}
				switch vt {
				case DVT_DateTime:
					buffer.WriteString("/Date(")
					buffer.WriteString(strconv.Itoa(int(DxCommonLib.TDateTime((*DxDoubleValue)(unsafe.Pointer(v)).fvalue).ToTime().Unix())*1000))
					buffer.WriteString(")/")
				case DVT_String:
					if escapJsonStr{
						buffer.WriteString(DxCommonLib.EscapeJsonStr((*DxStringValue)(unsafe.Pointer(v)).fvalue))
					}else{
						buffer.WriteString((*DxStringValue)(unsafe.Pointer(v)).fvalue)
					}
				case DVT_Array:
					buffer.Write((*DxArray)(unsafe.Pointer(v)).Bytes(escapJsonStr))
				case DVT_Record:
					buffer.Write((*DxRecord)(unsafe.Pointer(v)).Bytes(escapJsonStr))
				case DVT_RecordIntKey:
					buffer.Write((*DxIntKeyRecord)(unsafe.Pointer(v)).Bytes(escapJsonStr))
				default:
					buffer.Write(DxCommonLib.FastString2Byte(v.ToString()))
				}
				if vt == DVT_String || vt == DVT_Binary || vt == DVT_DateTime || vt == DVT_Ext{
					buffer.WriteByte('"')
				}
			}else{
				buffer.WriteString("null")
			}
		}
	}
	buffer.WriteByte('}')
	return buffer.Bytes()
}

func (r *DxIntKeyRecord)findPathNode(path string)(rec *DxBaseValue,keyName string)  {
	fields := strings.FieldsFunc(path,r.splitPathFields)
	vlen := len(fields)
	if vlen == 0{
		return nil,""
	}
	rParent := &r.DxBaseValue
	for i := 0;i < vlen - 1;i++{
		switch rParent.fValueType {
		case DVT_Record:
			rParent = (*DxRecord)(unsafe.Pointer(rParent)).Find(fields[i])
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(fields[i],10,64);err == nil{
				rParent = (*DxIntKeyRecord)(unsafe.Pointer(rParent)).Find(intkey)
			}else{
				panic("can not Convert to IntKey")
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(fields[i],10,64);err == nil{
				values := (*DxArray)(unsafe.Pointer(rParent)).fValues
				if values != nil && intkey >= 0 && intkey < int64(len(values)){
					rParent = values[intkey]
				}else{
					panic("KeyIndex Out bands of Array")
				}
			}else{
				panic("can not Convert to Int Index")
			}
		default:
			return nil,""
		}
		if rParent==nil{
			return nil,""
		}
	}
	return rParent,fields[vlen - 1]
}

func (r *DxIntKeyRecord)AsBytesByPath(Path string)[]byte  {
	parentBase,keyName := r.findPathNode(Path)
	if parentBase != nil {
		if keyName != ""{
			switch parentBase.fValueType {
			case DVT_Array:
				if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
					//panic("can not Convert to IntKey")
				}else{
					values := (*DxArray)(unsafe.Pointer(parentBase)).fValues
					if values != nil && intkey >= 0 && intkey < int64(len(values)){
						bt,_ :=  values[intkey].AsBytes()
						return bt
					}
				}
			case DVT_RecordIntKey:
				if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
					//panic("can not Convert to IntKey")
				}else{
					return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsBytes(intkey)
				}
			case DVT_Record:
				return (*DxRecord)(unsafe.Pointer(parentBase)).AsBytes(keyName)
			}
		}
	}
	return nil
}


func (r *DxIntKeyRecord)SetIntRecordValue(key int64,v *DxIntKeyRecord) {
	if v != nil && v.fParent != nil {
		panic("Must Set A Single Record(no Parent)")
	}
	if r.fRecords != nil{
		if value, ok := r.fRecords[key]; ok && value != nil {
			value.ClearValue(true)
			value.fParent = nil
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	if v != nil {
		r.fRecords[key] = &v.DxBaseValue
		v.fParent = &r.DxBaseValue
	}else {
		r.fRecords[key] = nil
	}
}

func (r *DxIntKeyRecord)SetRecordValue(key int64,v *DxRecord) {
	if v != nil && v.fParent != nil {
		panic("Must Set A Single Record(no Parent)")
	}
	if r.fRecords != nil{
		if value, ok := r.fRecords[key]; ok && value != nil {
			value.ClearValue(true)
			value.fParent = nil
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	if v != nil {
		r.fRecords[key] = &v.DxBaseValue
		v.PathSplitChar = r.PathSplitChar
		v.fParent = &r.DxBaseValue
	}else{
		r.fRecords[key] = nil
	}
}

func (r *DxIntKeyRecord)SetArray(key int64,v *DxArray)  {
	if v != nil && v.fParent != nil {
		panic("Must Set A Single Array(no Parent)")
	}
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			value.ClearValue(true)
			value.fParent = nil
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	if v!=nil{
		r.fRecords[key] = &v.DxBaseValue
		v.fParent = &r.DxBaseValue
	}else{
		r.fRecords[key] = nil
	}
}

func (r *DxIntKeyRecord)AsBaseValue(key int64)*DxBaseValue{
	if r.fRecords != nil{
		return r.fRecords[key]
	}
	return nil
}

func (r *DxIntKeyRecord)SetValue(key int64,v interface{})  {
	if v == nil{
		r.SetNull(key)
		return
	}
	switch value := v.(type) {
	case int: r.SetInt(key,value)
	case int32: r.SetInt32(key,value)
	case int64: r.SetInt64(key,value)
	case int8: r.SetInt(key,int(value))
	case uint8: r.SetInt(key,int(value))
	case int16: r.SetInt(key,int(value))
	case uint16: r.SetInt(key,int(value))
	case uint32: r.SetInt(key,int(value))
	case *int: r.SetInt(key,*value)
	case *int32: r.SetInt32(key,*value)
	case *int64: r.SetInt64(key,*value)
	case *int8: r.SetInt(key,int(*value))
	case *uint8: r.SetInt(key,int(*value))
	case *int16: r.SetInt(key,int(*value))
	case *uint16: r.SetInt(key,int(*value))
	case *uint32: r.SetInt(key,int(*value))
	case string: r.SetString(key,value)
	case []byte: r.SetBinary(key,value,true)
	case *[]byte: r.SetBinary(key,*value,true)
	case bool: r.SetBool(key,value)
	case *bool: r.SetBool(key,*value)
	case *string: r.SetString(key,*value)
	case float32: r.SetFloat(key,value)
	case float64: r.SetDouble(key,value)
	case *float32: r.SetFloat(key,*value)
	case *float64: r.SetDouble(key,*value)
	case *time.Time: r.SetDateTime(key,DxCommonLib.Time2DelphiTime(*value))
	case time.Time: r.SetDateTime(key,DxCommonLib.Time2DelphiTime(value))
	case *DxRecord: r.SetRecordValue(key,value)
	case DxRecord: r.SetRecordValue(key,&value)
	case *DxIntKeyRecord: r.SetIntRecordValue(key,value)
	case DxIntKeyRecord: r.SetIntRecordValue(key ,&value)
	case DxArray:  r.SetArray(key,&value)
	case *DxArray: r.SetArray(key,value)
	case DxInt64Value: r.SetInt64(key,value.fvalue)
	case *DxInt64Value: r.SetInt64(key,value.fvalue)
	case DxInt32Value: r.SetInt32(key,value.fvalue)
	case *DxInt32Value: r.SetInt32(key,value.fvalue)
	case DxFloatValue: r.SetFloat(key,value.fvalue)
	case *DxFloatValue: r.SetFloat(key,value.fvalue)
	case DxDoubleValue: r.SetDouble(key,value.fvalue)
	case *DxDoubleValue: r.SetDouble(key,value.fvalue)
	case DxBoolValue: r.SetBool(key,value.fvalue)
	case *DxBoolValue: r.SetBool(key,value.fvalue)
	case DxIntValue: r.SetInt(key,value.fvalue)
	case *DxIntValue: r.SetInt(key,value.fvalue)
	case DxStringValue: r.SetString(key,value.fvalue)
	case *DxStringValue: r.SetString(key,value.fvalue)
	case DxBinaryValue:  r.SetBinary(key,value.Bytes(),true)
	case *DxBinaryValue:  r.SetBinary(key,value.Bytes(),true)
	default:
		reflectv := reflect.ValueOf(v)
		rv := getRealValue(&reflectv)
		if rv == nil{
			if r.fRecords != nil{
				if _,ok := r.fRecords[key];!ok{
					r.fRecords[key] = nil
				}
			}
			return
		}
		switch rv.Kind(){
		case reflect.Struct:
			rec := r.NewRecord(key,true)
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
				rbase = &r.NewRecord(key,true).DxBaseValue
			case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
				rbase = &r.NewIntRecord(key).DxBaseValue
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
		case reflect.Slice,reflect.Array:
			arr := r.NewArray(key)
			vlen := rv.Len()
			for i := 0;i< vlen;i++{
				av := rv.Index(i)
				arrvalue := getRealValue(&av)
				switch arrvalue.Kind() {
				case reflect.Int,reflect.Uint32:
					arr.SetInt(i,int(arrvalue.Int()))
				case reflect.Bool:
					arr.SetBool(i,arrvalue.Bool())
				case reflect.Int64:
					arr.SetInt64(i,arrvalue.Int())
				case reflect.Int32,reflect.Int8,reflect.Int16,reflect.Uint8,reflect.Uint16:
					arr.SetInt32(i,int32(arrvalue.Int()))
				case reflect.Float32:
					arr.SetFloat(i,float32(arrvalue.Float()))
				case reflect.Float64:
					arr.SetDouble(i,arrvalue.Float())
				case reflect.String:
					arr.SetString(i,arrvalue.String())
				default:
					if arrvalue.CanInterface(){
						arr.SetValue(i,arrvalue.Interface())
					}
				}
			}
		}
	}
}


func (r *DxIntKeyRecord)KeyValueType(key int64)DxValueType  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			return value.fValueType
		}
	}
	return DVT_Null
}

func (r *DxIntKeyRecord)AsInt32(key int64,defavalue int32)int32  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
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
			case DVT_Double,DVT_DateTime:return int32((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
			case DVT_Float:return int32((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
			case DVT_String:
				v,err := strconv.ParseInt((*DxStringValue)(unsafe.Pointer(value)).fvalue,0,0)
				if err == nil{
					return int32(v)
				}
			default:
				//panic("can not convert Type to int32")
			}
		}
	}
	return defavalue
}

func (r *DxIntKeyRecord)AsInt32ByPath(path string,defavalue int32)int32  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsInt32(keyName,defavalue)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsInt32(intkey,defavalue)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						v,err := values[intkey].AsInt32()
						if err == nil{
							return v
						}
					}
				}
			}
		default:
			//panic("Path not A Parent Node")
		}
	}
	return defavalue
}

func (r *DxIntKeyRecord)Clone()*DxIntKeyRecord  {
	result := NewIntKeyRecord()
	result.PathSplitChar = r.PathSplitChar
	if r.fRecords != nil{
		for k,v := range r.fRecords{
			vbase := v.Clone()
			vbase.fParent = &result.DxBaseValue
			result.fRecords[k] = vbase
		}
	}
	return  result
}

func (r *DxIntKeyRecord)AsIntByPath(path string,defavalue int)int  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsInt(keyName,defavalue)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsInt(intkey,defavalue)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						v,err := values[intkey].AsInt()
						if err == nil{
							return v
						}
					}
				}
			}
		default:
			//panic("Path not A Parent Node")
		}
	}
	return defavalue
}

func (r *DxIntKeyRecord)AsInt(key int64,defavalue int)int  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
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
				if v,err := strconv.ParseInt((*DxStringValue)(unsafe.Pointer(value)).fvalue,0,0);err == nil{
					return int(v)
				}
			default:
				//panic("can not convert Type to int")
			}
		}
	}
	return defavalue
}

func (r *DxIntKeyRecord)AsInt64ByPath(path string,defavalue int64)int64  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsInt64(keyName,defavalue)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsInt64(intkey,defavalue)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						if v,err := values[intkey].AsInt64();err == nil{
							return v
						}
					}
				}
			}
		default:
			//panic("Path not A Parent Node")
		}
	}
	return defavalue
}

func (r *DxIntKeyRecord)AsInt64(key int64,defavalue int64)int64  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
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
			case DVT_Double,DVT_DateTime:return int64((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
			case DVT_Float:return int64((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
			case DVT_String:
				if v,err := strconv.ParseInt((*DxStringValue)(unsafe.Pointer(value)).fvalue,0,0);err == nil{
					return v
				}
			default:
				//panic("can not convert Type to int64")
			}
		}
	}
	return defavalue
}

func (r *DxIntKeyRecord)AsBoolByPath(path string,defavalue bool)bool  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsBool(keyName,defavalue)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsBool(intkey,defavalue)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						if v,err := values[intkey].AsBool();err == nil{
							return v
						}
					}
				}
			}
		default:
			//panic("Path not A Parent Node")
		}
	}
	return defavalue
}

func (r *DxIntKeyRecord)AsBool(key int64,defavalue bool)bool  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
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
				//panic("can not convert Type to Bool")
			}
		}
	}
	return defavalue
}


func (r *DxIntKeyRecord)AsFloatByPath(path string,defavalue float32)float32  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsFloat(keyName,defavalue)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsFloat(intkey,defavalue)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						if v,err := values[intkey].AsFloat();err == nil{
							return v
						}
					}
				}
			}
		default:
			//panic("Path not A Parent Node")
		}
	}
	return defavalue
}

func (r *DxIntKeyRecord)AsFloat(key int64,defavalue float32)float32  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			switch value.fValueType {
			case DVT_Int: return float32((*DxIntValue)(unsafe.Pointer(value)).fvalue)
			case DVT_Int32: return float32((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
			case DVT_Int64: return float32((*DxInt64Value)(unsafe.Pointer(value)).fvalue)
			case DVT_Bool:
				if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
					return 1
				}
				return 0
			case DVT_String: {
				v,e := strconv.ParseFloat((*DxStringValue)(unsafe.Pointer(value)).fvalue,2)
				if e == nil{
					return float32(v)
				}
			}
			case DVT_Double,DVT_DateTime:return float32((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
			case DVT_Float:return (*DxFloatValue)(unsafe.Pointer(value)).fvalue
			default:
			}
		}
	}
	return defavalue
}

func (r *DxIntKeyRecord)SetExtValue(intKey int64,extbt []byte)  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[intKey];ok && value != nil{
			if value.fValueType == DVT_Ext{
				value.SetExtValue(extbt)
				return
			}
			value.ClearValue(true)
		}
	}else{
		r.fRecords = make(map[int64]*DxBaseValue,32)
	}
	var m DxExtValue
	m.fdata = extbt
	m.fParent = &r.DxBaseValue
	m.fValueType = DVT_Ext
	if extbt != nil && len(extbt) > 0{
		m.fExtType = extbt[0]
	}
	r.fRecords[intKey] = &m.DxBaseValue
}

func (r *DxIntKeyRecord)SetBaseValue(intKey int64,v *DxBaseValue)  {
	if v != nil{
		switch v.fValueType {
		case DVT_Record,DVT_RecordIntKey,DVT_Array:
			if v != nil && v.fParent != nil {
				panic("Must Set A Single Record(no Parent)")
			}
			if r.fRecords != nil{
				if value, ok := r.fRecords[intKey]; ok && value != nil {
					value.ClearValue(true)
					value.fParent = nil
				}
			}else{
				r.fRecords = make(map[int64]*DxBaseValue,32)
			}
			if v != nil {
				r.fRecords[intKey] = v
				v.fParent = &r.DxBaseValue
			}else{
				r.fRecords[intKey] = nil
			}
		case DVT_Int:
			r.SetInt(intKey,(*DxIntValue)(unsafe.Pointer(v)).fvalue)
		case DVT_Int32:
			r.SetInt32(intKey,(*DxInt32Value)(unsafe.Pointer(v)).fvalue)
		case DVT_Int64:
			r.SetInt64(intKey,(*DxInt64Value)(unsafe.Pointer(v)).fvalue)
		case DVT_Binary:
			r.SetBinary(intKey,(*DxBinaryValue)(unsafe.Pointer(v)).fbinary,true)
		case DVT_Ext:
			r.SetExtValue(intKey,(*DxExtValue)(unsafe.Pointer(v)).ExtData())
		}
	}else {
		r.SetNull(intKey)
	}
}

func (r *DxIntKeyRecord)AsExtValue(intKey int64)(*DxExtValue)  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[intKey];ok && value != nil && value.fValueType == DVT_Ext{
			return (*DxExtValue)(unsafe.Pointer(value))
		}
	}
	return nil
}

func (r *DxIntKeyRecord)AsDoubleByPath(path string,defavalue float64)float64  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsDouble(keyName,defavalue)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsDouble(intkey,defavalue)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						if v,err := values[intkey].AsDouble();err == nil{
							return v
						}
					}
				}
			}
		default:
			//panic("Path not A Parent Node")
		}
	}
	return defavalue
}

func (r *DxIntKeyRecord)AsDouble(key int64,defavalue float64)float64  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			switch value.fValueType {
			case DVT_Int: return float64((*DxIntValue)(unsafe.Pointer(value)).fvalue)
			case DVT_Int32: return float64((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
			case DVT_Int64: return float64((*DxInt64Value)(unsafe.Pointer(value)).fvalue)
			case DVT_Bool:
				if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
					return 1
				}
				return 0
			case DVT_String: {
				v,e := strconv.ParseFloat((*DxStringValue)(unsafe.Pointer(value)).fvalue,2)
				if e == nil{
					return v
				}
			}
			case DVT_Double,DVT_DateTime:return (*DxDoubleValue)(unsafe.Pointer(value)).fvalue
			case DVT_Float:return float64((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
			default:
				//panic("can not convert Type to Double")
			}
		}
	}
	return defavalue
}

func (r *DxIntKeyRecord)AsDateTime(key int64,defavalue DxCommonLib.TDateTime)DxCommonLib.TDateTime  {
	if r.fRecords!= nil{
		if value,ok := r.fRecords[key];ok && value != nil{
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
				//panic("can not convert Type to Double")
			}
		}
	}
	return defavalue
}

func (r *DxIntKeyRecord)AsStringByPath(path string,defavalue string)string  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsString(keyName,defavalue)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsString(intkey,defavalue)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						return values[intkey].AsString()
					}
				}
			}
		default:
			//panic("Path not A Parent Node")
		}
	}
	return defavalue
}

func (r *DxIntKeyRecord)AsString(key int64,defavalue string)string  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			return value.ToString()
		}
	}
	return defavalue
}

func (r *DxIntKeyRecord)AsRecordByPath(path string)*DxRecord  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsRecord(keyName)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsRecord(intkey)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						if v,err := values[intkey].AsRecord();err == nil{
							return v
						}
					}
				}
			}
		default:
			//panic("Path not A Parent Node")
		}
	}
	return nil
}

func (r *DxIntKeyRecord)AsRecord(key int64)*DxRecord  {
	if value,ok := r.fRecords[key];ok && value != nil{
		if value.fValueType == DVT_Record{
			return (*DxRecord)(unsafe.Pointer(value))
		}
	}
	return nil
}

func (r *DxIntKeyRecord)AsIntRecord(key int64)*DxIntKeyRecord  {
	if r.fRecords != nil{
		if value,ok := r.fRecords[key];ok && value != nil{
			if value.fValueType == DVT_Record{
				return (*DxIntKeyRecord)(unsafe.Pointer(value))
			}
		}
	}
	return nil
}

func (r *DxIntKeyRecord)AsIntRecordByPath(path string)*DxIntKeyRecord  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsIntRecord(keyName)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsIntRecord(intkey)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						if v,err := values[intkey].AsIntRecord();err == nil{
							return v
						}
					}
				}
			}
		default:
			//panic("Path not A Parent Node")
		}
	}
	return nil
}

func (r *DxIntKeyRecord)AsArray(key int64)*DxArray  {
	if r.fRecords != nil {
		if value, ok := r.fRecords[key]; ok && value != nil {
			if value.fValueType == DVT_Array {
				return (*DxArray)(unsafe.Pointer(value))
			}
			//panic("not Array Value")
		}
	}
	return nil
}

func (r *DxIntKeyRecord)AsArrayByPath(path string)*DxArray  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsArray(keyName)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsArray(intkey)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						if v,err := values[intkey].AsArray();err == nil{
							return v
						}
					}
				}
			}
		default:
			//panic("Path not A Parent Node")
		}
	}
	return nil
}

func (r *DxIntKeyRecord)Length()int  {
	if r.fRecords != nil{
		return len(r.fRecords)
	}
	return 0
}

func (r *DxIntKeyRecord)Contains(keyName string)bool  {
	if r.fRecords != nil{
		if vr,vk := r.findPathNode(keyName);vr!=nil{
			switch vr.fValueType {
			case DVT_Record:
				_,ok := (*DxRecord)(unsafe.Pointer(vr)).fRecords[vk]
				return ok
			case DVT_RecordIntKey:
				if intkey,err := strconv.ParseInt(vk,10,64);err!=nil{
					//panic(err)
				}else{
					_,ok := (*DxIntKeyRecord)(unsafe.Pointer(vr)).fRecords[intkey]
					return ok
				}
			case DVT_Array:
				if intkey,err := strconv.ParseInt(vk,10,64);err!=nil{
					//panic(err)
				}else{
					arr := (*DxArray)(unsafe.Pointer(vr))
					return arr.fValues != nil && intkey >= 0 && intkey < int64(len(arr.fValues))
				}
			}
		}
	}
	return false
}

func (r *DxIntKeyRecord)Remove(keyOrPath string)  {
	if r.fRecords != nil{
		if vr,vk := r.findPathNode(keyOrPath);vr!=nil{
			switch vr.fValueType {
			case DVT_Record:
				if v,ok := (*DxRecord)(unsafe.Pointer(vr)).fRecords[vk];ok{
					v.ClearValue(true)
					delete((*DxRecord)(unsafe.Pointer(vr)).fRecords,vk)
				}
			case DVT_RecordIntKey:
				if intkey,err := strconv.ParseInt(vk,10,64);err!=nil{
					//panic(err)
				}else if v,ok := (*DxIntKeyRecord)(unsafe.Pointer(vr)).fRecords[intkey];ok{
					v.ClearValue(true)
					delete((*DxIntKeyRecord)(unsafe.Pointer(vr)).fRecords,intkey)
				}
			case DVT_Array:
				if intkey,err := strconv.ParseInt(vk,10,64);err!=nil{
					//panic(err)
				}else{
					arr := (*DxArray)(unsafe.Pointer(vr))
					arr.Delete(int(intkey))
				}
			}
		}
	}
}

func (r *DxIntKeyRecord)Delete(key int64)  {
	if r.fRecords!=nil {
		if v, ok := r.fRecords[key]; ok {
			v.ClearValue(true)
			delete(r.fRecords, key)
		}
	}
}

func (r *DxIntKeyRecord)Range(iteafunc func(key int64,value *DxBaseValue,params ...interface{})bool,params ...interface{}){
	if r.fRecords != nil && iteafunc!=nil{
		for k,v := range r.fRecords{
			if !iteafunc(k,v,params...){
				return
			}
		}
	}
}

func (r *DxIntKeyRecord)ToString()string  {
	return DxCommonLib.FastByte2String(r.Bytes(false))
}

func (r *DxIntKeyRecord)parserValue(key int64, b []byte,ConvertEscape,structRest bool)(parserlen int, err error)  {
	blen := len(b)
	i := 0
	valuestart := -1
	validCharIndex := -1
	startValue := false
	for i<blen {
		if !IsSpace(b[i]){
			switch b[i] {
			case ':':
				startValue = true
				//valuestart = i //自己记录有效的开始位置，和有效的结束位置，省去一个trim
			case '{':
				var rec DxRecord
				rec.PathSplitChar = r.PathSplitChar
				rec.fValueType = DVT_Record
				rec.fRecords = make(map[string]*DxBaseValue,32)
				if parserlen,err = rec.JsonParserFromByte(b[i:blen],ConvertEscape,structRest);err == nil{
					r.SetRecordValue(key,&rec)
				}
				parserlen+=i+1
				return
			case '[':
				var arr DxArray
				arr.fValueType = DVT_Array
				if parserlen,err = arr.JsonParserFromByte(b[i:],ConvertEscape,structRest);err == nil{
					r.SetArray(key,&arr)
				}
				parserlen+=i+1
				return
			case ',','}':
				//bvalue := bytes.Trim(b[valuestart + 1:i]," \r\n\t")
				bvalue := b[valuestart: validCharIndex+1]
				if len(bvalue) == 0{
					return i,ErrInvalidateJson
				}
				if bytes.IndexByte(bvalue,'.') > -1{
					vf := DxCommonLib.StrToFloatDef(DxCommonLib.FastByte2String(bvalue),0)
					r.SetDouble(key,vf)
				}else {
					st := DxCommonLib.FastByte2String(bvalue)
					if st == "true" || strings.ToUpper(st) == "TRUE"{
						r.SetBool(key,true)
					}else if st == "false" || strings.ToUpper(st) == "FALSE"{
						r.SetBool(key,false)
					}else if st == "null" || strings.ToUpper(st) == "NULL"{
						r.SetNull(key)
					}else{
						vf := DxCommonLib.StrToIntDef(st,0)
						if vf <= math.MaxInt32 && vf>=math.MinInt32{
							r.SetInt(key,int(vf))
						}else{
							r.SetInt64(key,vf)
						}
					}
				}
				return i,nil
			case '"': //string
				i++
				isInEscape := false
				for j :=i; j<blen;j++{
					if IsSpace(b[j]){
						continue
					}
					switch b[j] {
					case '"':
						if isInEscape{
							isInEscape = false
							continue
						}
						//字符串完毕了！
						st := ""
						bvalue := b[i:j]
						if ConvertEscape{
							st = DxCommonLib.ParserEscapeStr(bvalue)
						}else{
							st = DxCommonLib.FastByte2String(bvalue)
						}
						jt := DxCommonLib.ParserJsonTime(st)
						if jt >= 0{
							r.SetDateTime(key,jt)
						}else{
							r.SetString(key,st)
						}
						return j+1,nil
					case '\\':
						isInEscape = !isInEscape
					default:
						if isInEscape{
							//判断是否是有效的转义
							if b[j] == 't'|| b[j] == 'b'|| b[j] == 'f'|| b[j] == 'n'|| b[j] == 'r'|| b[j] == '\\'|| b[j] == '"' || b[j]=='u'|| b[j]=='U' || b[j] == '/'{
								//有效的转义
								isInEscape = false
							}else{
								return j,ErrInvalidateJson
							}
						}
					}
				}
			default:
				if !startValue && valuestart == -1{
					return i,ErrInvalidateJson
				}
				if valuestart == -1{
					valuestart = i
					startValue = false
				}else{
					validCharIndex = i
				}
			}
		}
		i += 1
	}
	return blen,ErrInvalidateJson
}

func (r *DxIntKeyRecord)LoadJsonFile(fileName string,ConvertEscape,structRest bool)error  {
	databytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	if len(databytes) > 2 && databytes[0] == 0xEF && databytes[1] == 0xBB && databytes[2] == 0xBF{//BOM
		databytes = databytes[3:]
	}
	_,err = r.JsonParserFromByte(databytes,ConvertEscape,structRest)
	return err
}

func (r *DxIntKeyRecord)SaveJsonWriter(w io.Writer)error  {
	writer := bufio.NewWriter(w)
	err := writer.WriteByte('{')
	if err != nil{
		return err
	}
	if r.fRecords != nil {
		isFirst := true
		for k, v := range r.fRecords {
			if !isFirst {
				_, err = writer.WriteString(`,"`)
			} else {
				isFirst = false
				err = writer.WriteByte('"')
			}
			if err != nil {
				return err
			}
			_, err = writer.WriteString(strconv.Itoa(int(k)))
			if err != nil {
				return err
			}
			_, err = writer.WriteString(`":`)
			if err != nil {
				return err
			}
			if v != nil {
				vt := v.fValueType
				if vt == DVT_String || vt == DVT_Binary {
					err = writer.WriteByte('"')
				}
				if err != nil {
					return err
				}
				_, err = writer.WriteString(v.ToString())
				if err == nil && (vt == DVT_String || vt == DVT_Binary) {
					err = writer.WriteByte('"')
				}
			} else {
				_, err = writer.WriteString("null")
			}
			if err != nil {
				return err
			}
		}
	}
	writer.WriteByte('}')
	err = writer.Flush()
	return err
}

func (r *DxIntKeyRecord)SaveJsonFile(fileName string,BOMFile bool)error  {
	if file,err := os.OpenFile(fileName,os.O_CREATE | os.O_TRUNC,0644);err == nil{
		defer file.Close()
		if BOMFile{
			file.Write([]byte{0xEF,0xBB,0xBF})
		}
		return r.SaveJsonWriter(file)
	}else{
		return err
	}
}

func (r *DxIntKeyRecord)LoadJsonReader(reader io.Reader)error  {
	return nil
}

func (r *DxIntKeyRecord)LoadMsgPackReader(reader io.Reader)error  {
	return NewDecoder(reader).Decode(&r.DxBaseValue)
}

func (r *DxIntKeyRecord)LoadMsgPackFile(fileName string)error  {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return NewDecoder(f).Decode(&r.DxBaseValue)
}

func (r *DxIntKeyRecord)SaveMsgPackFile(fileName string)error  {
	if file,err := os.OpenFile(fileName,os.O_CREATE | os.O_TRUNC,0644);err == nil{
		defer file.Close()
		return NewEncoder(file).EncodeRecordIntKey(r)
	}else{
		return err
	}
}

func (r *DxIntKeyRecord)JsonParserFromByte(JsonByte []byte,ConvertEscape,structRest bool)(parserlen int, err error)  {
	i := 0
	if structRest{
		r.ClearValue(false)
	}
	objStart := false
	keyStart := false
	btlen := len(JsonByte)
	plen := -1
	keyName := ""
	for i < btlen{
		if IsSpace(JsonByte[i]){
			i++
			continue
		}
		if !objStart && JsonByte[i] != '{' {
			return 0,ErrInvalidateJson
		}
		switch JsonByte[i]{
		case '{':
			objStart = true
			keyStart = true
		case '}':
			if keyStart{
				return i,ErrInvalidateJson
			}
			objStart = false
			return i,nil
		case '"': //keyName
			if keyStart{
				//获取string
				plen = bytes.IndexByte(JsonByte[i+1:btlen],'"')
				if plen > -1{
					keyName = DxCommonLib.FastByte2String(JsonByte[i+1:i+1+plen])
				}
				i += plen+2
				keyStart = false
				//解析Value
				intkey,err := strconv.ParseInt(keyName,10,64)
				if err != nil{
					return i,err
				}
				if ilen,err := r.parserValue(intkey,JsonByte[i:btlen],ConvertEscape,structRest);err!=nil{
					return ilen + i,err
				}else{
					i += ilen
					continue
				}
			}
		case ',': //next key
			if keyStart{
				return i,ErrInvalidateJson
			}
			keyStart = true
		case ':': //value
			if keyStart || objStart{
				return i,ErrInvalidateJson
			}
		case '[':
			if objStart || keyStart{
				return i,ErrInvalidateJson
			}
		case ']':
			if objStart || keyStart {
				return i,ErrInvalidateJson
			}
		default:

		}
		i+=1
	}
	return btlen,ErrInvalidateJson
}

//增加值编码器
func (r *DxIntKeyRecord) Encode(valuecoder Coders.Encoder) error{
	//NewEncoder(file).EncodeRecord(r)
	var err error
	switch valuecoder.Name() {
	case "msgpack":
		if msgpacker,ok := valuecoder.(*DxMsgPackEncoder);ok{
			return msgpacker.EncodeRecordIntKey(r)
		}
		encoder := valuecoder.(*DxMsgPack.MsgPackEncoder)
		maplen := uint(r.Length())
		if maplen <= DxMsgPack.Max_fixmap_len{   //fixmap
			err = encoder.WriteByte(0x80 | byte(maplen))
		}else if maplen <= DxMsgPack.Max_map16_len{
			//写入长度
			err = encoder.WriteUint16(uint16(maplen),DxMsgPack.CodeMap16)
		}else{
			if maplen > DxMsgPack.Max_map32_len{
				maplen = DxMsgPack.Max_map32_len
			}
			err = encoder.WriteUint32(uint32(maplen),DxMsgPack.CodeMap32)
		}
		if err != nil{
			return err
		}
		//写入对象信息,Kv对
		for k,v := range r.fRecords{
			if err = encoder.EncodeInt(k);err!=nil{
				return err
			}
			if v != nil{
				err = v.Encode(encoder)
			}else{
				err = encoder.WriteByte(0xc0) //null
			}
			if err!=nil{
				return err
			}
		}

		return nil
	case "json":
	}
	return nil
}

func NewIntKeyRecord()*DxIntKeyRecord  {
	result := new(DxIntKeyRecord)
	result.PathSplitChar = DefaultPathSplit
	result.fValueType = DVT_RecordIntKey
	result.fRecords = make(map[int64]*DxBaseValue,32)
	return result
}
