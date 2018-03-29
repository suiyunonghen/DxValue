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
	"github.com/suiyunonghen/DxCommonLib"
	"strings"
	"math"
	"strconv"
	"io/ioutil"
	"io"
	"bufio"
	"os"
	"time"
)

/******************************************************
*  DxRecord
******************************************************/
type(
	DxRecord		struct{
		DxBaseValue
		fRecords		map[string]*DxBaseValue
		PathSplitChar	byte
	}
)


func (r *DxRecord)ClearValue(clearInner bool)  {
	if r.fRecords != nil{
		for _,v := range r.fRecords{
			v.ClearValue(clearInner)
			v.fParent = nil
		}
	}
	if !clearInner{
		if r.fRecords == nil || len(r.fRecords) > 0{
			r.fRecords = make(map[string]*DxBaseValue,32)
		}
	}else{
		r.fRecords = nil
	}
}

func (r *DxRecord)splitPathFields(charrune rune) bool {
	return charrune == rune(r.PathSplitChar)
}

func (r *DxRecord)getSize()int  {
	result := 0
	if r.fRecords != nil {
		for k,v := range r.fRecords{
			result += len(k)+v.Size()
		}
	}
	return result
}

func (r *DxRecord)NewRecord(keyName string)(rec *DxRecord)  {
	if value,ok := r.fRecords[keyName];ok && value != nil{
		if value.fValueType == DVT_Record{
			rec = (*DxRecord)(unsafe.Pointer(value))
			rec.ClearValue(false)
			rec.fParent = &r.DxBaseValue
			return
		}
		value.ClearValue(true)
	}
	rec = new(DxRecord)
	rec.fValueType = DVT_Record
	rec.PathSplitChar = r.PathSplitChar
	rec.fRecords = make(map[string]*DxBaseValue,32)
	r.fRecords[keyName] = &rec.DxBaseValue
	rec.fParent = &r.DxBaseValue
	return
}

func (r *DxRecord)NewIntRecord(keyName string)(rec *DxIntKeyRecord)  {
	if value,ok := r.fRecords[keyName];ok && value != nil{
		if value.fValueType == DVT_RecordIntKey{
			rec = (*DxIntKeyRecord)(unsafe.Pointer(value))
			rec.ClearValue(false)
			rec.fParent = &r.DxBaseValue
			return
		}
		value.ClearValue(true)
	}
	rec = new(DxIntKeyRecord)
	rec.fValueType = DVT_RecordIntKey
	rec.PathSplitChar = r.PathSplitChar
	rec.fRecords = make(map[int64]*DxBaseValue,32)
	r.fRecords[keyName] = &rec.DxBaseValue
	rec.fParent = &r.DxBaseValue
	return
}

func (r *DxRecord)Find(keyName string)*DxBaseValue  {
	if v,ok := r.fRecords[keyName];ok{
		return v
	}
	return nil
}



func (r *DxRecord)ForcePath(path string,v interface{}) {
	fields := strings.FieldsFunc(path,r.splitPathFields)
	vlen := len(fields)
	if vlen == 0{
		return
	}
	vbase := r.Find(fields[0])
	if vbase == nil{
		vbase = &r.NewRecord(fields[0]).DxBaseValue
	}
	for i := 1;i<vlen - 1;i++{
		if vbase != nil {
			if vbase.fValueType == DVT_Record{
				vbase = (*DxRecord)(unsafe.Pointer(vbase)).Find(fields[i])
			}else{
				if intkey,er := strconv.ParseInt(fields[i],10,64);er == nil{
					vbase = (*DxIntKeyRecord)(unsafe.Pointer(vbase)).Find(intkey)
				}else{
					vbase = vbase.Parent()
					vbase = &(*DxRecord)(unsafe.Pointer(vbase)).NewRecord(fields[i - 1]).NewRecord(fields[i]).DxBaseValue
				}
			}
		}
	}
	if vbase.fValueType == DVT_Record{
		(*DxRecord)(unsafe.Pointer(vbase)).SetValue(fields[vlen - 1],v)
	}else{
		if intkey,er := strconv.ParseInt(fields[vlen - 1],10,64);er == nil{
			(*DxIntKeyRecord)(unsafe.Pointer(vbase)).SetValue(intkey,v)
		}else{
			vbase = vbase.Parent()
			if vlen > 2{
				vbase = &(*DxRecord)(unsafe.Pointer(vbase)).NewRecord(fields[vlen - 2]).NewRecord(fields[vlen - 1]).DxBaseValue
			}else{
				vbase = &r.NewRecord(fields[vlen - 2]).DxBaseValue
			}
			(*DxRecord)(unsafe.Pointer(vbase)).SetValue(fields[vlen - 1],v)
		}
	}
}

func (r *DxRecord)NewArray(keyName string)(arr *DxArray)  {
	if value,ok := r.fRecords[keyName];ok && value != nil{
		if value.fValueType == DVT_Array{
			arr = (*DxArray)(unsafe.Pointer(value))
			arr.ClearValue(false)
			arr.fParent = &r.DxBaseValue
			return
		}
		value.ClearValue(true)
	}
	arr = new(DxArray)
	arr.fValueType = DVT_Array
	arr.fParent = &r.DxBaseValue
	r.fRecords[keyName] = &arr.DxBaseValue
	return
}

func (r *DxRecord)SetInt(KeyName string,v int)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
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
	var m DxIntValue
	m.fvalue = v
	m.fValueType = DVT_Int
	m.fParent = &r.DxBaseValue
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetInt32(KeyName string,v int32)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
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
	var m DxInt32Value
	m.fvalue = v
	m.fValueType = DVT_Int32
	m.fParent = &r.DxBaseValue
	r.fRecords[KeyName] = &m.DxBaseValue
}


func (r *DxRecord)SetInt64(KeyName string,v int64)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
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
	if v <= math.MaxInt32 && v >= math.MinInt32{
		var m DxInt32Value
		m.fvalue = int32(v)
		m.fValueType = DVT_Int32
		m.fParent = &r.DxBaseValue
		r.fRecords[KeyName] = &m.DxBaseValue
	}else{
		var m DxInt64Value
		m.fvalue = v
		m.fValueType = DVT_Int64
		m.fParent = &r.DxBaseValue
		r.fRecords[KeyName] = &m.DxBaseValue
	}
}

func (r *DxRecord)SetBool(KeyName string,v bool)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Bool {
			(*DxBoolValue)(unsafe.Pointer(value)).fvalue = v
			return
		}
		value.ClearValue(true)
	}
	var m DxBoolValue
	m.fvalue = v
	m.fParent = &r.DxBaseValue
	m.fValueType = DVT_Bool
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetNull(KeyName string)  {
	if v,ok := r.fRecords[KeyName];ok{
		v.ClearValue(true)
	}
	r.fRecords[KeyName] = nil
}



func (r *DxRecord)SetFloat(KeyName string,v float32)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Float{
			(*DxFloatValue)(unsafe.Pointer(value)).fvalue = v
			return
		}else if value.fValueType == DVT_Double || value.fValueType == DVT_DateTime{
			(*DxDoubleValue)(unsafe.Pointer(value)).fvalue = float64(v)
			return
		}
		value.ClearValue(true)
	}
	var m DxFloatValue
	m.fvalue = v
	m.fParent = &r.DxBaseValue
	m.fValueType = DVT_Float
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetDouble(KeyName string,v float64)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
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
	var m DxDoubleValue
	m.fvalue = v
	m.fParent = &r.DxBaseValue
	m.fValueType = DVT_Double
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetDateTime(KeyName string,v DxCommonLib.TDateTime)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_DateTime{
			(*DxDoubleValue)(unsafe.Pointer(value)).fvalue = float64(v)
			return
		}
		value.ClearValue(true)
	}
	var m DxDoubleValue
	m.fvalue = float64(v)
	m.fValueType = DVT_DateTime
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)AsDateTime(keyName string,defavalue DxCommonLib.TDateTime)DxCommonLib.TDateTime  {
	if value,ok := r.fRecords[keyName];ok && value != nil{
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
		case DVT_String:
			if result,err := value.AsDateTime();err == nil{
				return result
			}else{
				panic("can not convert Type to TDateTime")
			}
		default:
			panic("can not convert Type to TDateTime")
		}
	}
	return defavalue
}

func (r *DxRecord)SetString(KeyName string,v string)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_String:
			(*DxStringValue)(unsafe.Pointer(value)).fvalue = v
			return
		case DVT_DateTime:
			if jt := DxCommonLib.ParserJsonTime(v);jt >= 0{
				(*DxDoubleValue)(unsafe.Pointer(value)).fvalue = float64(jt)
				return
			}
			t,err := time.Parse("2006-01-02T15:04:05Z",v)
			if err == nil{
				(*DxDoubleValue)(unsafe.Pointer(value)).fvalue = float64(DxCommonLib.Time2DelphiTime(&t))
				return
			}
			t,err = time.Parse("2006-01-02 15:04:05",v)
			if err == nil{
				(*DxDoubleValue)(unsafe.Pointer(value)).fvalue = float64(DxCommonLib.Time2DelphiTime(&t))
				return
			}
			t,err = time.Parse("2006/01/02 15:04:05",v)
			if err == nil{
				(*DxDoubleValue)(unsafe.Pointer(value)).fvalue = float64(DxCommonLib.Time2DelphiTime(&t))
				return
			}
		case DVT_Bool:
			if v == "true" || strings.ToLower(v) == "true"{
				(*DxBoolValue)(unsafe.Pointer(value)).fvalue = true
				return
			}else if v == "false" || strings.ToLower(v) == "false" {
				(*DxBoolValue)(unsafe.Pointer(value)).fvalue = false
				return
			}
		case DVT_Int:
			if iv,err := strconv.Atoi(v);err == nil{
				(*DxIntValue)(unsafe.Pointer(value)).fvalue = iv
				return
			}
		case DVT_Int32:
			if iv,err := strconv.Atoi(v);err == nil{
				(*DxInt32Value)(unsafe.Pointer(value)).fvalue = int32(iv)
				return
			}
		case DVT_Int64:
			if iv,err := strconv.ParseInt(v,10,64);err == nil{
				(*DxInt64Value)(unsafe.Pointer(value)).fvalue = iv
				return
			}
		case DVT_Float:
			if iv,err := strconv.ParseFloat(v,32);err == nil{
				(*DxFloatValue)(unsafe.Pointer(value)).fvalue = float32(iv)
				return
			}
		case DVT_Double:
			if iv,err := strconv.ParseFloat(v,64);err == nil{
				(*DxDoubleValue)(unsafe.Pointer(value)).fvalue = iv
				return
			}
		}
		value.ClearValue(true)
	}
	var m DxStringValue
	m.fvalue = v
	m.fValueType = DVT_String
	r.fRecords[KeyName] = &m.DxBaseValue
}


func (r *DxRecord)SetBinary(KeyName string,v []byte,reWrite bool,encodeType DxBinaryEncodeType)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Binary{
			bv := (*DxBinaryValue)(unsafe.Pointer(value))
			bv.EncodeType = encodeType
			if reWrite{
				bv.SetBinary(v,true)
			}else{
				bv.Append(v)
			}
			return
		}
		value.ClearValue(true)
	}
	var m DxBinaryValue
	m.fbinary = v
	m.fParent = &r.DxBaseValue
	m.fValueType = DVT_Binary
	m.EncodeType = encodeType
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetExtValue(keyName string,extbt []byte)  {
	if value,ok := r.fRecords[keyName];ok && value != nil{
		if value.fValueType == DVT_Ext{
			value.SetExtValue(extbt)
			return
		}
		value.ClearValue(true)
	}
	var m DxExtValue
	m.fdata = extbt
	m.fParent = &r.DxBaseValue
	if extbt != nil && len(extbt) > 0{
		m.fExtType = extbt[0]
	}
	m.fValueType = DVT_Ext
	r.fRecords[keyName] = &m.DxBaseValue
}

func (r *DxRecord)AsExtValue(keyName string)(*DxExtValue)  {
	if value,ok := r.fRecords[keyName];ok && value != nil && value.fValueType == DVT_Ext{
		return (*DxExtValue)(unsafe.Pointer(value))
	}
	return nil
}

func (r *DxRecord)AsBytes(keyName string)[]byte  {
	if value,ok := r.fRecords[keyName];ok && value != nil{
		bt,_ := value.AsBytes()
		return bt
	}
	return nil
}

type targetBuffer interface {
	WriteString(string)(int, error)
	WriteByte(c byte) error
}

func (r *DxRecord)EncodeJson2Writer(w io.Writer)  {
	var buffer targetBuffer
	if buf,ok := w.(targetBuffer);ok{
		buffer = buf
	} else {
		buffer = bufio.NewWriter(w)
	}
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
			if vt == DVT_String || vt == DVT_Binary || vt == DVT_DateTime || vt == DVT_Ext{
				buffer.WriteByte('"')
			}
			if vt == DVT_DateTime{
				buffer.WriteString("/Date(")
				buffer.WriteString(strconv.Itoa(int(DxCommonLib.TDateTime((*DxDoubleValue)(unsafe.Pointer(v)).fvalue).ToTime().Unix())*1000))
				buffer.WriteString(")/")
			}else{
				buffer.WriteString(v.ToString())
			}

			if vt == DVT_String || vt == DVT_Binary || vt == DVT_DateTime || vt == DVT_Ext{
				buffer.WriteByte('"')
			}
		}else{
			buffer.WriteString("null")
		}
	}
	buffer.WriteByte('}')
}

func (r *DxRecord)Bytes()[]byte  {
	buffer := bytes.NewBuffer(make([]byte,0,512))
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
			if vt == DVT_String || vt == DVT_Binary || vt == DVT_DateTime || vt == DVT_Ext{
				buffer.WriteByte('"')
			}
			if vt == DVT_DateTime{
				buffer.WriteString("/Date(")
				buffer.WriteString(strconv.Itoa(int(DxCommonLib.TDateTime((*DxDoubleValue)(unsafe.Pointer(v)).fvalue).ToTime().Unix())*1000))
				buffer.WriteString(")/")
			}else{
				buffer.WriteString(v.ToString())
			}

			if vt == DVT_String || vt == DVT_Binary || vt == DVT_DateTime || vt == DVT_Ext{
				buffer.WriteByte('"')
			}
		}else{
			buffer.WriteString("null")
		}
	}
	buffer.WriteByte('}')
	return buffer.Bytes()
}

func (r *DxRecord)findPathNode(path string)(rec *DxBaseValue,keyName string)  {
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
	}
	return rParent,fields[vlen - 1]
}

func (r *DxRecord)AsBytesByPath(Path string)[]byte  {
	parentBase,keyName := r.findPathNode(Path)
	if parentBase != nil {
		if keyName != ""{
			switch parentBase.fValueType {
			case DVT_Array:
				if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
					panic("can not Convert to IntKey")
				}else{
					values := (*DxArray)(unsafe.Pointer(parentBase)).fValues
					if values != nil && intkey >= 0 && intkey < int64(len(values)){
						bt,_ :=  values[intkey].AsBytes()
						return bt
					}else{
						panic("KeyIndex Out bands of Array")
					}
				}
			case DVT_RecordIntKey:
				if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
					panic("can not Convert to IntKey")
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


func (r *DxRecord)SetRecordValue(keyName string,v *DxRecord) {
	if v != nil && v.fParent != nil {
		panic("Must Set A Single Record(no Parent)")
	}
	if value, ok := r.fRecords[keyName]; ok && value != nil {
		value.ClearValue(true)
		if  value.fValueType == DVT_Record {
			nrec := (*DxRecord)(unsafe.Pointer(value))
			nrec.fParent = nil
			*nrec = *v
			nrec.fParent = &r.DxBaseValue
			return
		}
	}
	if v != nil {
		r.fRecords[keyName] = &v.DxBaseValue
		v.fParent = &r.DxBaseValue
	}
}

func (r *DxRecord)SetBaseValue(keyName string,v *DxBaseValue)  {
	if v != nil{
		switch v.fValueType {
		case DVT_Record:
			r.SetRecordValue(keyName,(*DxRecord)(unsafe.Pointer(v)))
		case DVT_RecordIntKey:
			r.SetIntRecordValue(keyName,(*DxIntKeyRecord)(unsafe.Pointer(v)))
		case DVT_Array:
			r.SetArray(keyName,(*DxArray)(unsafe.Pointer(v)))
		}
	}else {
		r.SetNull(keyName)
	}
}

func (r *DxRecord)SetIntRecordValue(keyName string,v *DxIntKeyRecord) {
	if v != nil && v.fParent != nil {
		panic("Must Set A Single Record(no Parent)")
	}
	if value, ok := r.fRecords[keyName]; ok && value != nil {
		value.ClearValue(true)
		if  value.fValueType == DVT_RecordIntKey {
			nrec := (*DxIntKeyRecord)(unsafe.Pointer(value))
			nrec.fParent = nil
			*nrec = *v
			nrec.fParent = &r.DxBaseValue
			return
		}
	}
	if v != nil {
		r.fRecords[keyName] = &v.DxBaseValue
		v.fParent = &r.DxBaseValue
	}
}

func (r *DxRecord)SetArray(KeyName string,v *DxArray)  {
	if v != nil && v.fParent != nil {
		panic("Must Set A Single Array(no Parent)")
	}
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		value.ClearValue(true)
		if value.fValueType == DVT_Array{
			arr := (*DxArray)(unsafe.Pointer(value))
			arr.fParent = nil
			*arr = *v
			arr.fParent = &r.DxBaseValue
			return
		}
	}
	if v!=nil{
		r.fRecords[KeyName] = &v.DxBaseValue
		v.fParent = &r.DxBaseValue
	}
}

func (r *DxRecord)AsBaseValue(keyName string)*DxBaseValue{
	if r.fRecords != nil{
		return r.fRecords[keyName]
	}
	return nil
}

func (r *DxRecord)SetValue(keyName string,v interface{})  {
	if v == nil{
		r.SetNull(keyName)
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
	case []byte: r.SetBinary(keyName,value,true,BET_Base64)
	case *[]byte: r.SetBinary(keyName,*value,true,BET_Base64)
	case bool: r.SetBool(keyName,value)
	case *bool: r.SetBool(keyName,*value)
	case *string: r.SetString(keyName,*value)
	case float32: r.SetFloat(keyName,value)
	case float64: r.SetDouble(keyName,value)
	case *float32: r.SetFloat(keyName,*value)
	case *float64: r.SetDouble(keyName,*value)
	case time.Time: r.SetDateTime(keyName,DxCommonLib.Time2DelphiTime(&value))
	case *time.Time: r.SetDateTime(keyName,DxCommonLib.Time2DelphiTime(value))
	case *DxRecord: r.SetRecordValue(keyName,value)
	case DxRecord: r.SetRecordValue(keyName,&value)
	case *DxIntKeyRecord: r.SetIntRecordValue(keyName,value)
	case DxIntKeyRecord: r.SetIntRecordValue(keyName,&value)
	case DxArray:  r.SetArray(keyName,&value)
	case *DxArray: r.SetArray(keyName,value)
	case DxInt64Value: r.SetInt64(keyName,value.fvalue)
	case *DxInt64Value: r.SetInt64(keyName,value.fvalue)
	case DxInt32Value: r.SetInt32(keyName,value.fvalue)
	case *DxInt32Value: r.SetInt32(keyName,value.fvalue)
	case DxFloatValue: r.SetFloat(keyName,value.fvalue)
	case *DxFloatValue: r.SetFloat(keyName,value.fvalue)
	case DxDoubleValue: r.SetDouble(keyName,value.fvalue)
	case *DxDoubleValue: r.SetDouble(keyName,value.fvalue)
	case DxBoolValue: r.SetBool(keyName,value.fvalue)
	case *DxBoolValue: r.SetBool(keyName,value.fvalue)
	case DxIntValue: r.SetInt(keyName,value.fvalue)
	case *DxIntValue: r.SetInt(keyName,value.fvalue)
	case DxStringValue: r.SetString(keyName,value.fvalue)
	case *DxStringValue: r.SetString(keyName,value.fvalue)
	case DxBinaryValue:  r.SetBinary(keyName,value.Bytes(),true,BET_Base64)
	case *DxBinaryValue:  r.SetBinary(keyName,value.Bytes(),true,BET_Base64)
	default:
		reflectv := reflect.ValueOf(v)
		rv := getRealValue(&reflectv)
		if rv == nil{
			if _,ok := r.fRecords[keyName];!ok{
				r.fRecords[keyName] = nil
			}
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
			mapkeys := rv.MapKeys()
			if len(mapkeys) == 0{
				return
			}
			keytype := rv.Type().Key()
			var rbase *DxBaseValue
			switch getBaseType(keytype) {
			case reflect.String:
				rbase = &r.NewRecord(keyName).DxBaseValue
			case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
				rbase = &r.NewIntRecord(keyName).DxBaseValue
			default:
				panic("Invalidate Record Key,Can Only Int or String")
			}
			rvalue := rv.MapIndex(mapkeys[0])
			//获得Value类型
			valueKind := getBaseType(rvalue.Type())
			for _,kv := range mapkeys{
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
			arr := r.NewArray(keyName)
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


func (r *DxRecord)KeyValueType(KeyName string)DxValueType  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		return value.fValueType
	}
	return DVT_Null
}

func (r *DxRecord)AsInt32(KeyName string,defavalue int32)int32  {
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
		case DVT_Double,DVT_DateTime:return int32((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return int32((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		default:
			panic("can not convert Type to int32")
		}
	}
	return defavalue
}

func (r *DxRecord)AsInt32ByPath(path string,defavalue int32)int32  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsInt32(keyName,defavalue)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsInt32(intkey,defavalue)
			}else{
				panic(err)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						v,err := values[intkey].AsInt32()
						if err != nil{
							panic(err)
						}
						return v
					}else{
						panic("array Index out of bounds")
					}
				}
			}else{
				panic(err)
			}
		default:
			panic("Path not A Parent Node")
		}
	}
	return defavalue
}

func (r *DxRecord)AsIntByPath(path string,defavalue int)int  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsInt(keyName,defavalue)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsInt(intkey,defavalue)
			}else{
				panic(err)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						v,err := values[intkey].AsInt()
						if err != nil{
							panic(err)
						}
						return v
					}else{
						panic("array Index out of bounds")
					}
				}
			}else{
				panic(err)
			}
		default:
			panic("Path not A Parent Node")
		}
	}
	return defavalue
}

func (r *DxRecord)AsInt(KeyName string,defavalue int)int  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
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
		default:
			panic("can not convert Type to int")
		}
	}
	return defavalue
}

func (r *DxRecord)AsInt64ByPath(path string,defavalue int64)int64  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsInt64(keyName,defavalue)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsInt64(intkey,defavalue)
			}else{
				panic(err)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						v,err := values[intkey].AsInt64()
						if err != nil{
							panic(err)
						}
						return v
					}else{
						panic("array Index out of bounds")
					}
				}
			}else{
				panic(err)
			}
		default:
			panic("Path not A Parent Node")
		}
	}
	return defavalue
}

func (r *DxRecord)AsInt64(KeyName string,defavalue int64)int64  {
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
		case DVT_Double,DVT_DateTime:return int64((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return int64((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		default:
			panic("can not convert Type to int64")
		}
	}
	return defavalue
}

func (r *DxRecord)AsBoolByPath(path string,defavalue bool)bool  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsBool(keyName,defavalue)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsBool(intkey,defavalue)
			}else{
				panic(err)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						v,err := values[intkey].AsBool()
						if err != nil{
							panic(err)
						}
						return v
					}else{
						panic("array Index out of bounds")
					}
				}
			}else{
				panic(err)
			}
		default:
			panic("Path not A Parent Node")
		}
	}
	return defavalue
}

func (r *DxRecord)AsBool(KeyName string,defavalue bool)bool  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_Int: return (*DxIntValue)(unsafe.Pointer(value)).fvalue != 0
		case DVT_Int32: return (*DxInt32Value)(unsafe.Pointer(value)).fvalue != 0
		case DVT_Int64: return (*DxInt64Value)(unsafe.Pointer(value)).fvalue != 0
		case DVT_Bool: return bool((*DxBoolValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Double,DVT_DateTime:return float64((*DxDoubleValue)(unsafe.Pointer(value)).fvalue) != 0
		case DVT_Float:return float32((*DxFloatValue)(unsafe.Pointer(value)).fvalue) != 0
		default:
			panic("can not convert Type to Bool")
		}
	}
	return defavalue
}


func (r *DxRecord)AsFloatByPath(path string,defavalue float32)float32  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsFloat(keyName,defavalue)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsFloat(intkey,defavalue)
			}else{
				panic(err)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						v,err := values[intkey].AsFloat()
						if err != nil{
							panic(err)
						}
						return v
					}else{
						panic("array Index out of bounds")
					}
				}
			}else{
				panic(err)
			}
		default:
			panic("Path not A Parent Node")
		}
	}
	return defavalue
}

func (r *DxRecord)AsFloat(KeyName string,defavalue float32)float32  {
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
		case DVT_Double,DVT_DateTime:return float32((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return (*DxFloatValue)(unsafe.Pointer(value)).fvalue
		default:
			panic("can not convert Type to Float")
		}
	}
	return defavalue
}


func (r *DxRecord)AsDoubleByPath(path string,defavalue float64)float64  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsDouble(keyName,defavalue)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsDouble(intkey,defavalue)
			}else{
				panic(err)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						v,err := values[intkey].AsDouble()
						if err != nil{
							panic(err)
						}
						return v
					}else{
						panic("array Index out of bounds")
					}
				}
			}else{
				panic(err)
			}
		default:
			panic("Path not A Parent Node")
		}
	}
	return defavalue
}

func (r *DxRecord)AsDouble(KeyName string,defavalue float64)float64  {
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
		case DVT_Float,DVT_DateTime:return float64((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		default:
			panic("can not convert Type to Double")
		}
	}
	return defavalue
}

func (r *DxRecord)AsStringByPath(path string,defavalue string)string  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsString(keyName,defavalue)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsString(intkey,defavalue)
			}else{
				panic(err)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						return values[intkey].AsString()
					}else{
						panic("array Index out of bounds")
					}
				}
			}else{
				panic(err)
			}
		default:
			panic("Path not A Parent Node")
		}
	}
	return defavalue
}

func (r *DxRecord)AsString(KeyName string,defavalue string)string  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		return value.ToString()
	}
	return defavalue
}

func (r *DxRecord)AsRecordByPath(path string)*DxRecord  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsRecord(keyName)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsRecord(intkey)
			}else{
				panic(err)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						v,err := values[intkey].AsRecord()
						if err != nil{
							panic(err)
						}
						return v
					}else{
						panic("array Index out of bounds")
					}
				}
			}else{
				panic(err)
			}
		default:
			panic("Path not A Parent Node")
		}
	}
	return nil
}

func (r *DxRecord)AsRecord(KeyName string)*DxRecord  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Record{
			return (*DxRecord)(unsafe.Pointer(value))
		}
		panic("not Record Value")
	}
	return nil
}

func (r *DxRecord)AsIntRecord(KeyName string)*DxIntKeyRecord  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_RecordIntKey{
			return (*DxIntKeyRecord)(unsafe.Pointer(value))
		}
		panic("not Record Value")
	}
	return nil
}

func (r *DxRecord)AsIntRecordByPath(path string)*DxIntKeyRecord  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsIntRecord(keyName)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsIntRecord(intkey)
			}else{
				panic(err)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						v,err := values[intkey].AsIntRecord()
						if err != nil{
							panic(err)
						}
						return v
					}else{
						panic("array Index out of bounds")
					}
				}
			}else{
				panic(err)
			}
		default:
			panic("Path not A Parent Node")
		}
	}
	return nil
}

func (r *DxRecord)AsArray(KeyName string)*DxArray  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Array{
			return (*DxArray)(unsafe.Pointer(value))
		}
		panic("not Array Value")
	}
	return nil
}

func (r *DxRecord)AsArrayByPath(path string)*DxArray  {
	parentBase,keyName := r.findPathNode(path)
	if parentBase != nil && keyName != ""{
		switch parentBase.fValueType {
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(parentBase)).AsArray(keyName)
		case DVT_RecordIntKey:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				return (*DxIntKeyRecord)(unsafe.Pointer(parentBase)).AsArray(intkey)
			}else{
				panic(err)
			}
		case DVT_Array:
			if intkey,err := strconv.ParseInt(keyName,10,64);err == nil{
				if values := (*DxArray)(unsafe.Pointer(parentBase)).fValues;values != nil{
					if intkey >= 0 && intkey < int64(len(values)){
						v,err := values[intkey].AsArray()
						if err != nil{
							panic(err)
						}
						return v
					}else{
						panic("array Index out of bounds")
					}
				}
			}else{
				panic(err)
			}
		default:
			panic("Path not A Parent Node")
		}
	}
	return nil
}

func (r *DxRecord)Length()int  {
	if r.fRecords != nil{
		return len(r.fRecords)
	}
	return 0
}

func (r *DxRecord)Contains(keyName string)bool  {
	if r.fRecords != nil{
		if vr,vk := r.findPathNode(keyName);vr!=nil{
			switch vr.fValueType {
			case DVT_Record:
				_,ok := (*DxRecord)(unsafe.Pointer(vr)).fRecords[vk]
				return ok
			case DVT_RecordIntKey:
				if intkey,err := strconv.ParseInt(vk,10,64);err!=nil{
					return false
				}else{
					_,ok := (*DxIntKeyRecord)(unsafe.Pointer(vr)).fRecords[intkey]
					return ok
				}
			case DVT_Array:
				if intkey,err := strconv.ParseInt(vk,10,64);err!=nil{
					return false
				}else{
					arr := (*DxArray)(unsafe.Pointer(vr))
					return arr.fValues != nil && intkey >= 0 && intkey < int64(len(arr.fValues))
				}
			}
		}
	}
	return false
}

func (r *DxRecord)Remove(keyOrPath string)  {
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
					panic(err)
				}else if v,ok := (*DxIntKeyRecord)(unsafe.Pointer(vr)).fRecords[intkey];ok{
					v.ClearValue(true)
					delete((*DxIntKeyRecord)(unsafe.Pointer(vr)).fRecords,intkey)
				}
			case DVT_Array:
				if intkey,err := strconv.ParseInt(vk,10,64);err!=nil{
					panic(err)
				}else{
					arr := (*DxArray)(unsafe.Pointer(vr))
					arr.Delete(int(intkey))
				}
			}
		}
	}
}

func (r *DxRecord)Delete(key string)  {
	if v,ok := r.fRecords[key];ok{
		v.ClearValue(true)
		delete(r.fRecords,key)
	}
}

func (r *DxRecord)Range(iteafunc func(keyName string,value *DxBaseValue)bool){
	if r.fRecords != nil && iteafunc!=nil{
		for k,v := range r.fRecords{
			if !iteafunc(k,v){
				return
			}
		}
	}
}

func (r *DxRecord)ToString()string  {
	return DxCommonLib.FastByte2String(r.Bytes())
}

func (r *DxRecord)parserValue(keyName string, b []byte,ConvertEscape,structRest bool)(parserlen int, err error)  {
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
					r.SetRecordValue(keyName,&rec)
				}
				parserlen+=2 //会多解析一个{
				return
			case '[':
				var arr DxArray
				arr.fValueType = DVT_Array
				if parserlen,err = arr.JsonParserFromByte(b[i:],ConvertEscape,structRest);err == nil{
					r.SetArray(keyName,&arr)
				}
				parserlen+=2
				return
			case ',','}':
				var bvalue []byte
				if validCharIndex == -1 {
					bvalue = bytes.Trim(b[valuestart:i]," \r\n\t")
				}else{
					bvalue = b[valuestart: validCharIndex+1]
				}
				if len(bvalue) == 0{
					return i,ErrInvalidateJson
				}
				if bytes.IndexByte(bvalue,'.') > -1{
					if vf,err := strconv.ParseFloat(DxCommonLib.FastByte2String(bvalue),64);err!=nil{
						return i,ErrInvalidateJson
					}else{
						r.SetDouble(keyName,vf)
					}
				}else {
					st := DxCommonLib.FastByte2String(bvalue)
					if st == "true" || strings.ToUpper(st) == "TRUE"{
						r.SetBool(keyName,true)
					}else if st == "false" || strings.ToUpper(st) == "FALSE"{
						r.SetBool(keyName,false)
					}else if st == "null" || strings.ToUpper(st) == "NULL"{
						r.SetNull(keyName)
					}else{
						if vf,err := strconv.Atoi(st);err!=nil{
							return i,ErrInvalidateJson
						}else{
							r.SetInt(keyName,vf)
						}
					}
				}
				return i,nil
			case '"': //string
				plen := bytes.IndexByte(b[i+1:blen],'"')
				if plen > -1{
					bvalue := b[i+1:plen+i+1]
					st := ""
					if ConvertEscape{
						st = DxCommonLib.ParserEscapeStr(bvalue)
						jt := DxCommonLib.ParserJsonTime(st)
						if jt >= 0{
							r.SetDateTime(keyName,jt)
							return plen + i + 2,nil
						}
					}else{
						st = DxCommonLib.FastByte2String(bvalue)
					}
					r.SetString(keyName,st)
					return plen + i + 2,nil
				}
				return i,ErrInvalidateJson
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

func (r *DxRecord)LoadJsonFile(fileName string,ConvertEscape,structRest bool)error  {
	databytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	if databytes[0] == 0xEF && databytes[1] == 0xBB && databytes[2] == 0xBF{//BOM
		databytes = databytes[3:]
	}
	_,err = r.JsonParserFromByte(databytes,ConvertEscape,structRest)
	return err
}

func (r *DxRecord)SaveJsonWriter(w io.Writer)error  {
	writer := bufio.NewWriter(w)
	err := writer.WriteByte('{')
	if err != nil{
		return err
	}
	isFirst := true
	for k,v := range r.fRecords{
		if !isFirst{
			_,err = writer.WriteString(`,"`)
		}else{
			isFirst = false
			err = writer.WriteByte('"')
		}
		if err != nil{
			return err
		}
		_,err = writer.WriteString(k)
		if err!=nil{
			return err
		}
		_, err = writer.WriteString(`":`)
		if err!=nil{
			return err
		}
		if v != nil{
			vt := v.fValueType
			if vt == DVT_String || vt == DVT_Binary || vt == DVT_DateTime{
				err = writer.WriteByte('"')
			}
			if err != nil{
				return err
			}
			if vt == DVT_DateTime{
				_,err = writer.WriteString("/Date(")
				_,err = writer.WriteString(strconv.Itoa(int(DxCommonLib.TDateTime((*DxDoubleValue)(unsafe.Pointer(v)).fvalue).ToTime().Unix())*1000))
				_,err = writer.WriteString(")/")
			}else{
				_,err = writer.WriteString(v.ToString())
			}

			if err == nil && (vt == DVT_String || vt == DVT_Binary || vt == DVT_DateTime) {
				err = writer.WriteByte('"')
			}
		}else{
			_,err = writer.WriteString("null")
		}
		if err != nil{
			return err
		}
	}
	writer.WriteByte('}')
	err = writer.Flush()
	return err
}

func (r *DxRecord)SaveJsonFile(fileName string,BOMFile bool)error  {
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


func (r *DxRecord)LoadJsonReader(reader io.Reader)error  {
	return nil
}

func (r *DxRecord)SaveMsgPackFile(fileName string)error  {
	if file,err := os.OpenFile(fileName,os.O_CREATE | os.O_TRUNC,0644);err == nil{
		defer file.Close()
		return NewEncoder(file).EncodeRecord(r)
	}else{
		return err
	}
}

func (r *DxRecord)LoadMsgPackReader(reader io.Reader)error  {
	return NewDecoder (reader).Decode(&r.DxBaseValue)
}

func (r *DxRecord)LoadMsgPackFile(fileName string)error  {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return NewDecoder (f).Decode(&r.DxBaseValue)
}

func (r *DxRecord)JsonParserFromByte(JsonByte []byte,ConvertEscape,structRest bool)(parserlen int, err error)  {
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
				if ilen,err := r.parserValue(keyName,JsonByte[i:btlen],ConvertEscape,structRest);err!=nil{
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
			if keyStart || keyStart{
				return i,ErrInvalidateJson
			}
		default:

		}
		i+=1
	}
	return btlen,ErrInvalidateJson
}

func NewRecord()*DxRecord  {
	result := new(DxRecord)
	result.PathSplitChar = '.'
	result.fValueType = DVT_Record
	result.fRecords = make(map[string]*DxBaseValue,32)
	return result
}
