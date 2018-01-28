package DxValue

import (
	"unsafe"
	"math"
	"github.com/suiyunonghen/DxCommonLib"
	"io/ioutil"
	"io"
	"reflect"
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
		default:
			v.fValue.ClearValue(true)
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
		default:
			v.fValue.ClearValue(true)
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
			default:
				v.fValue.ClearValue(true)
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

func (v *DxValue)AsDateTime()(DxCommonLib.TDateTime,error)  {
	if v.fValue == nil{
		return 0,nil
	}
	return v.fValue.AsDateTime()
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
	return v.fValue.AsBytes()
}

func (v *DxValue)ClearValue()  {
	if v.fValue != nil{
		v.fValue.ClearValue(true)
	}
	v.fValue = nil
}

func (v *DxValue)NewRecord()*DxRecord  {
	var rec *DxRecord
	if v.fValue == nil || v.fValue.fValueType != DVT_Record{
		v.fValue.ClearValue(true)
		rec = &DxRecord{}
		rec.PathSplitChar = '.'
		rec.fValueType = DVT_Record
		rec.fRecords = make(map[string]*DxBaseValue,32)
		v.fValue = &rec.DxBaseValue
	}else{
		rec = (*DxRecord)(unsafe.Pointer(v.fValue))
		rec.ClearValue(false)
	}
	return rec
}

func (v *DxValue)NewIntRecord()*DxIntKeyRecord  {
	var rec *DxIntKeyRecord
	if v.fValue == nil || v.fValue.fValueType != DVT_RecordIntKey{
		v.fValue.ClearValue(true)
		rec = &DxIntKeyRecord{}
		rec.PathSplitChar = '.'
		rec.fValueType = DVT_RecordIntKey
		rec.fRecords = make(map[int64]*DxBaseValue,32)
		v.fValue = &rec.DxBaseValue
	}else{
		rec = (*DxIntKeyRecord)(unsafe.Pointer(v.fValue))
		rec.ClearValue(false)
	}
	return rec
}

func (v *DxValue)NewArray()*DxArray  {
	var arr *DxArray
	if v.fValue == nil || v.fValue.fValueType != DVT_Array{
		v.fValue.ClearValue(true)
		arr = &DxArray{}
		arr.fValueType = DVT_Array
		v.fValue = &arr.DxBaseValue
	}else{
		arr = (*DxArray)(unsafe.Pointer(v.fValue))
		arr.ClearValue(false)
	}
	return arr
}

func (v *DxValue)JsonParserFromByte(JsonByte []byte,ConvertEscape bool)(parserlen int, err error)  {
	v.ClearValue()
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

func (r *DxValue)LoadJsonReader(reader io.Reader)error  {
	return nil
}


func (v *DxValue)LoadJsonFile(fileName string,ConvertEscape bool)error  {
	databytes, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		return err
	}
	if databytes[0] == 0xEF && databytes[1] == 0xBB && databytes[2] == 0xBF{//BOM
		databytes = databytes[3:]
	}
	_,err = v.JsonParserFromByte(databytes,ConvertEscape)
	return err
}

func (v *DxValue)SaveJsonFile(fileName string,BOMFile bool)error  {
	if v.fValue != nil{
		switch v.fValue.fValueType {
		case DVT_Array:
			return (*DxArray)(unsafe.Pointer(v.fValue)).SaveJsonFile(fileName,BOMFile)
		case DVT_Record:
			return (*DxRecord)(unsafe.Pointer(v.fValue)).SaveJsonFile(fileName,BOMFile)
		case DVT_RecordIntKey:
			return (*DxIntKeyRecord)(unsafe.Pointer(v.fValue)).SaveJsonFile(fileName,BOMFile)
		}
	}
	return nil
}

func NewDxValue(avalue interface{})(result *DxValue)  {
	result = &DxValue{}
	if avalue == nil{
		return
	}
	switch value := avalue.(type) {
	case int:
		result.SetIntValue(value)
	case int32:
		result.SetInt32Value(value)
	case int64:
		result.SetInt64Value(value)
	case int8:
		result.SetIntValue(int(value))
	case uint8:
		result.SetIntValue(int(value))
	case int16:
		result.SetIntValue(int(value))
	case uint16:
		result.SetIntValue(int(value))
	case uint32:
		result.SetIntValue(int(value))
	case *int:
		result.SetIntValue(*value)
	case *int32:
		result.SetInt32Value(*value)
	case *int64:
		result.SetInt64Value(*value)
	case *int8:
		result.SetIntValue(int(*value))
	case *uint8:
		result.SetIntValue(int(*value))
	case *int16:
		result.SetIntValue(int(*value))
	case *uint16:
		result.SetIntValue(int(*value))
	case *uint32:
		result.SetIntValue(int(*value))
	case string:
		var v DxStringValue
		v.fvalue = value
		v.fValueType = DVT_String
		result.fValue = &v.DxBaseValue
	case []byte:
		var v DxBinaryValue
		v.fValueType = DVT_Binary
		v.fbinary = value
		result.fValue = &v.DxBaseValue
	case *[]byte:
		var v DxBinaryValue
		v.fValueType = DVT_Binary
		v.fbinary = *value
		result.fValue = &v.DxBaseValue
	case bool:
		var v DxBoolValue
		v.fValueType = DVT_Bool
		v.fvalue = value
		result.fValue = &v.DxBaseValue
	case *bool:
		var v DxBoolValue
		v.fValueType = DVT_Bool
		v.fvalue = *value
		result.fValue = &v.DxBaseValue
	case *string:
		var v DxStringValue
		v.fvalue = *value
		v.fValueType = DVT_String
		result.fValue = &v.DxBaseValue
	case float32:
		var v DxFloatValue
		v.fvalue = value
		v.fValueType = DVT_Float
		result.fValue = &v.DxBaseValue
	case float64:
		var v DxDoubleValue
		v.fvalue = value
		v.fValueType = DVT_Double
		result.fValue = &v.DxBaseValue
	case *float32:
		var v DxFloatValue
		v.fvalue = *value
		v.fValueType = DVT_Float
		result.fValue = &v.DxBaseValue
	case *float64:
		var v DxDoubleValue
		v.fvalue = *value
		v.fValueType = DVT_Double
		result.fValue = &v.DxBaseValue
	case *DxRecord:
		result.fValue = &value.DxBaseValue
	case DxRecord:
		result.fValue = &value.DxBaseValue
	case *DxIntKeyRecord:
		result.fValue = &value.DxBaseValue
	case DxIntKeyRecord:
		result.fValue = &value.DxBaseValue
	case DxArray:
		result.fValue = &value.DxBaseValue
	case *DxArray:
		result.fValue = &value.DxBaseValue
	case DxInt64Value:
		result.SetInt64Value(value.fvalue)
	case *DxInt64Value:
		result.SetInt64Value(value.fvalue)
	case DxInt32Value:
		result.SetInt32Value(value.fvalue)
	case *DxInt32Value:
		result.SetInt32Value(value.fvalue)
	case DxFloatValue:
		var v DxFloatValue
		v.fvalue = value.fvalue
		v.fValueType = DVT_Float
		result.fValue = &v.DxBaseValue
	case *DxFloatValue:
		var v DxFloatValue
		v.fvalue = value.fvalue
		v.fValueType = DVT_Float
		result.fValue = &v.DxBaseValue
	case DxDoubleValue:
		var v DxDoubleValue
		v.fvalue = value.fvalue
		v.fValueType = DVT_Double
		result.fValue = &v.DxBaseValue
	case *DxDoubleValue:
		var v DxDoubleValue
		v.fvalue = value.fvalue
		v.fValueType = DVT_Double
		result.fValue = &v.DxBaseValue
	case DxBoolValue:
		var v DxBoolValue
		v.fvalue = value.fvalue
		v.fValueType = DVT_Bool
		result.fValue = &v.DxBaseValue
	case *DxBoolValue:
		var v DxBoolValue
		v.fvalue = value.fvalue
		v.fValueType = DVT_Bool
		result.fValue = &v.DxBaseValue
	case DxIntValue:
		result.SetIntValue(value.fvalue)
	case *DxIntValue:
		result.SetIntValue(value.fvalue)
	case DxStringValue:
		var v DxStringValue
		v.fvalue = value.fvalue
		v.fValueType = DVT_String
		result.fValue = &v.DxBaseValue
	case *DxStringValue:
		var v DxStringValue
		v.fvalue = value.fvalue
		v.fValueType = DVT_String
		result.fValue = &v.DxBaseValue
	case DxBinaryValue:
		var v DxBinaryValue
		v.fValueType = DVT_Binary
		v.SetBinary(value.Bytes(),true)
		result.fValue = &v.DxBaseValue
	case *DxBinaryValue:
		var v DxBinaryValue
		v.fValueType = DVT_Binary
		v.SetBinary(value.Bytes(),true)
		result.fValue = &v.DxBaseValue
	default:
		reflectv := reflect.ValueOf(avalue)
		rv := getRealValue(&reflectv)
		if rv == nil{
			return
		}
		switch rv.Kind(){
		case reflect.Struct:
			rec := NewRecord()
			result.fValue = &rec.DxBaseValue
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
				rbase = &NewRecord().DxBaseValue
			case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
				rbase = &NewIntKeyRecord().DxBaseValue
			default:
				panic("Invalidate Record Key,Can Only Int or String")
			}
			rvalue := rv.MapIndex(mapkeys[0])
			result.fValue = rbase
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
			arr := NewArray()
			result.fValue = &arr.DxBaseValue
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
	return
}