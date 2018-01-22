package DxValue

import (
	"unsafe"
	"bytes"
)

/******************************************************
*  DxRecord
******************************************************/
type  DxRecord		struct{
	DxBaseValue
	fRecords		map[string]*DxBaseValue
}


func (r *DxRecord)ClearValue()  {
	r.fRecords = make(map[string]*DxBaseValue,20)
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
	rec.fRecords = make(map[string]*DxBaseValue,20)
	r.fRecords[keyName] = &rec.DxBaseValue
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

func (r *DxRecord)SetValue(keyName string,v interface{})  {
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
	result.fRecords = make(map[string]*DxBaseValue,20)
	return result
}
