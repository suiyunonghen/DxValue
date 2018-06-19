package DxJson

import (
	"github.com/suiyunonghen/DxValue/Coders"
	"bufio"
	"strconv"
	"time"
	"errors"
	"reflect"
	"fmt"
	"io"
	"bytes"
)

type JsonWriter interface {
	io.Writer
	WriteString(s string) (int, error)
	WriteByte(c byte) error
}

type  JsonEncoder  struct{
	w   JsonWriter
}



const jsonName  = "json"

var(
	valueEncoders []Coders.EncoderFunc
)

func init()  {
	valueEncoders = []Coders.EncoderFunc{
		reflect.Bool:          encodeBoolValue,
		reflect.Int:           encodeInt64Value,
		reflect.Int8:          encodeInt64Value,
		reflect.Int16:         encodeInt64Value,
		reflect.Int32:         encodeInt64Value,
		reflect.Int64:         encodeInt64Value,
		reflect.Uint:          encodeInt64Value,
		reflect.Uint8:         encodeInt64Value,
		reflect.Uint16:        encodeInt64Value,
		reflect.Uint32:        encodeInt64Value,
		reflect.Uint64:        encodeInt64Value,
		reflect.Float32:       encodeFloat32Value,
		reflect.Float64:       encodeFloat64Value,
		reflect.Complex64:     encodeUnsupportedValue,
		reflect.Complex128:    encodeUnsupportedValue,
		reflect.Array:         encodeArrayValue,
		reflect.Chan:          encodeUnsupportedValue,
		reflect.Func:          encodeUnsupportedValue,
		reflect.Interface:     encodeInterfaceValue,
		reflect.Ptr:           encodeUnsupportedValue,
		reflect.Slice:         encodeSliceValue,
		reflect.String:        encodeStringValue,
		reflect.Struct:        encodeStructValue,
		reflect.UnsafePointer: encodeUnsupportedValue,
	}
	Coders.RegisterCoderName(jsonName)
}

func encodeStructValue(encoder Coders.Encoder, strct reflect.Value) error {
	structFields := Coders.Fields(strct.Type())
	mapLen := structFields.Len(jsonName)
	jsonencoder := encoder.(*JsonEncoder)
	err := jsonencoder.StartObject()
	if err != nil {
		return err
	}
	IsFist := true
	for i:=0;i<mapLen;i++{
		f := structFields.Field(i)
		MarshName := f.MarshalName(jsonName)
		if MarshName != "-" {
			if IsFist{
				IsFist = false
			}else{
				err := jsonencoder.w.WriteByte(',')
				if err != nil{
					return err
				}
			}
			err = jsonencoder.EncodeKeyName(MarshName)
			if err!=nil{
				return err
			}
			err = f.EncodeValue(encoder,strct)
			if err != nil{
				return err
			}
		}
	}

	return jsonencoder.EndObject()
}

func encodeBoolValue(encoder Coders.Encoder,v reflect.Value) error  {
	return encoder.(*JsonEncoder).EncodeBool(v.Bool(),true)
}

func encodeInt64Value(encoder Coders.Encoder,v reflect.Value) error  {
	return encoder.(*JsonEncoder).EncodeInt(v.Int(),true)
}

func encodeFloat32Value(encoder Coders.Encoder,v reflect.Value) error  {
	return encoder.(*JsonEncoder).EncodeFloat(float32(v.Float()),true)
}

func encodeFloat64Value(encoder Coders.Encoder,v reflect.Value) error  {
	return encoder.(*JsonEncoder).EncodeDouble(v.Float(),true)
}

func encodeUnsupportedValue(encoder Coders.Encoder,v reflect.Value) error  {
	return fmt.Errorf("msgpack: Encode(unsupported %s)", v.Type())
}

func encodeStringValue(encoder Coders.Encoder,v reflect.Value) error  {
	return encoder.(*JsonEncoder).EncodeString(v.String(),true)
}

func encodeInterfaceValue(encoder Coders.Encoder, v reflect.Value) error {
	if v.IsNil() {
		return errors.New("nil Can't Encode")
	}
	v = v.Elem()
	if encodeFunc := encoder.(*JsonEncoder).GetEncoderFunc(v.Type());encodeFunc !=nil{
		return encodeFunc(encoder,v)
	}else if v.CanInterface(){
		return encoder.EncodeStand(v.Interface())
	}else {
		return  fmt.Errorf("json: Encode(unsupported %s)", v.Type())
	}
}

func encodeArrayValue(encoder Coders.Encoder,v reflect.Value)(err error)  {
	arlen := uint(v.Len())
	jsonEncoder := encoder.(*JsonEncoder)
	IsFirst := true
	jsonEncoder.StartArray()
	for i := uint(0);i< arlen;i++ {
		if IsFirst{
			IsFirst = false
		}else{
			jsonEncoder.w.WriteByte(',')
		}
		av := v.Index(int(i))
		arrvalue := Coders.GetRealValue(&av)
		if arrvalue == nil {
			err = jsonEncoder.EncodeNil(i==arlen - 1)
		}else{
			if encodeFunc := jsonEncoder.GetEncoderFunc(arrvalue.Type());encodeFunc !=nil{
				err = encodeFunc(encoder,*arrvalue)
			}else if arrvalue.CanInterface(){
				err = encoder.EncodeStand(arrvalue.Interface())
			}
		}
		if err!=nil{
			return
		}
	}
	return jsonEncoder.EndObject()
}

func encodeSliceValue(encoder Coders.Encoder,v reflect.Value)(err error)  {
	if v.IsNil() {
		return encoder.(*JsonEncoder).EncodeNil(true)
	}else{
		return encodeArrayValue(encoder,v)
	}
}

func (encoder *JsonEncoder)EncodeCustom()(error){
	return nil
}

func (encoder *JsonEncoder)Name()string{
	return jsonName
}

func encodeMapStringStringValue(e Coders.Encoder, v reflect.Value) error {
	encoder := e.(*JsonEncoder)
	if v.IsNil() {
		return encoder.EncodeNil(true)
	}

	m := v.Convert(Coders.MapStringStringType).Interface().(map[string]string)
	err := encoder.StartObject()
	if err != nil{
		return err
	}
	IsFirst := true
	for mk, mv := range m {
		if IsFirst{
			IsFirst = false
		}else{
			err := encoder.w.WriteByte(',')
			if err != nil{
				return err
			}
		}
		if err := encoder.EncodeKeyName(mk); err != nil {
			return err
		}
		if err := encoder.EncodeString(mv,true); err != nil {
			return err
		}
	}
	return encoder.EndObject()
}

func encodeMapStringInterfaceValue(e Coders.Encoder, v reflect.Value) error {
	encoder := e.(*JsonEncoder)
	if v.IsNil() {
		return encoder.EncodeNil(true)
	}

	if err := encoder.StartObject(); err != nil {
		return err
	}

	IsFirst := true
	m := v.Convert(Coders.MapStringInterfaceType).Interface().(map[string]interface{})
	for mk, mv := range m {
		if IsFirst{
			IsFirst = false
		}else{
			err := encoder.w.WriteByte(',')
			if err != nil{
				return err
			}
		}
		if err := encoder.EncodeKeyName(mk); err != nil {
			return err
		}
		if err := encoder.EncodeStand(mv); err != nil {
			return err
		}
	}

	return encoder.EndObject()
}

func encodeMapInt64StringValue(e Coders.Encoder, v reflect.Value) error {
	encoder := e.(*JsonEncoder)
	if v.IsNil() {
		return encoder.EncodeNil(true)
	}

	if err := encoder.StartObject(); err != nil {
		return err
	}
	IsFirst := true
	m := v.Convert(Coders.MapIntStringType).Interface().(map[int64]string)
	for mk, mv := range m {
		if IsFirst{
			IsFirst = false
		}else{
			err := encoder.w.WriteByte(',')
			if err != nil{
				return err
			}
		}
		if err := encoder.EncodeKeyName(strconv.Itoa(int(mk))); err != nil {
			return err
		}
		if err := encoder.EncodeString(mv,true); err != nil {
			return err
		}
	}

	return encoder.EndObject()
}

func encodeMapInt64InterfaceValue(e Coders.Encoder, v reflect.Value) error {
	encoder := e.(*JsonEncoder)
	if v.IsNil() {
		return encoder.EncodeNil(true)
	}

	if err := encoder.StartObject(); err != nil {
		return err
	}

	IsFirst := true
	m := v.Convert(Coders.MapStringInterfaceType).Interface().(map[int64]interface{})
	for mk, mv := range m {
		if IsFirst{
			IsFirst = false
		}else{
			err := encoder.w.WriteByte(',')
			if err != nil{
				return err
			}
		}
		if err := encoder.EncodeKeyName(strconv.Itoa(int(mk))); err != nil {
			return err
		}
		if err := encoder.EncodeStand(mv); err != nil {
			return err
		}
	}

	return encoder.EndObject()
}

func encodeMapIntStringValue(e Coders.Encoder, v reflect.Value) error {
	encoder := e.(*JsonEncoder)
	if v.IsNil() {
		return encoder.EncodeNil(true)
	}

	if err := encoder.StartObject(); err != nil {
		return err
	}
	IsFirst := true
	m := v.Convert(Coders.MapIntStringType).Interface().(map[int]string)
	for mk, mv := range m {
		if IsFirst{
			IsFirst = false
		}else{
			err := encoder.w.WriteByte(',')
			if err != nil{
				return err
			}
		}
		if err := encoder.EncodeKeyName(strconv.Itoa(mk)); err != nil {
			return err
		}
		if err := encoder.EncodeString(mv,true); err != nil {
			return err
		}
	}
	return encoder.EndObject()
}

func encodeMapIntInterfaceValue(e Coders.Encoder, v reflect.Value) error {
	encoder := e.(*JsonEncoder)
	if v.IsNil() {
		return encoder.EncodeNil(true)
	}

	if err := encoder.StartObject(); err != nil {
		return err
	}

	IsFirst := true
	m := v.Convert(Coders.MapStringInterfaceType).Interface().(map[int]interface{})
	for mk, mv := range m {
		if IsFirst{
			IsFirst = false
		}else{
			err := encoder.w.WriteByte(',')
			if err != nil{
				return err
			}
		}
		if err := encoder.EncodeKeyName(strconv.Itoa(int(mk))); err != nil {
			return err
		}
		if err := encoder.EncodeStand(mv); err != nil {
			return err
		}
	}

	return nil
}

func (encoder *JsonEncoder)GetEncoderFunc(typ reflect.Type)Coders.EncoderFunc{
	kind := typ.Kind()

	if typ == Coders.ErrorType {
		return nil
	}

	switch kind {
	case reflect.Ptr:
		return func(e Coders.Encoder, v reflect.Value) error {
			if v.IsNil() {
				return e.(*JsonEncoder).EncodeNil(true)
			}
			encoderFunc := encoder.GetEncoderFunc(typ.Elem())
			return encoderFunc(e, v.Elem())
		}
	case reflect.Map:
		switch typ.Key() {
		case Coders.StringType:
			switch typ.Elem(){
			case Coders.StringType:
				return encodeMapStringStringValue
			case Coders.InterfaceType:
				return encodeMapStringInterfaceValue
			}
		case Coders.Int64Type:
			switch typ.Elem(){
			case Coders.StringType:
				return encodeMapInt64StringValue
			case Coders.InterfaceType:
				return encodeMapInt64InterfaceValue
			}
		case Coders.IntType:
			switch typ.Elem(){
			case Coders.StringType:
				return encodeMapIntStringValue
			case Coders.InterfaceType:
				return encodeMapIntInterfaceValue
			}
		}

	}
	result := valueEncoders[kind]
	if result == nil{
		result = encodeUnsupportedValue
	}
	return result
}

func (encoder *JsonEncoder)EncodeStand(v interface{})(error){
	switch value := v.(type) {
	case *map[string]interface{}:
	case map[string]interface{}:
	case *map[string]string:
	case map[string]string:
	case []interface{}:
	case []string:
	case []*string:
	case string:
		encoder.EncodeString(value,true)
	case *string:
		encoder.EncodeString(*value,true)
	case time.Time:
		encoder.EncodeTime(value,true)
	case *time.Time:
		encoder.EncodeTime(*value,true)
	case *int8:
		return encoder.EncodeInt(int64(*value),true)
	case int8:
		return encoder.EncodeInt(int64(value),true)
	case *int16:
		return encoder.EncodeInt(int64(*value),true)
	case int16:
		return encoder.EncodeInt(int64(value),true)
	case *int32:
		return encoder.EncodeInt(int64(*value),true)
	case int32:
		return encoder.EncodeInt(int64(value),true)
	case *int64:
		return encoder.EncodeInt(*value,true)
	case int64:
		return encoder.EncodeInt(value,true)
	case *uint8:
		return encoder.EncodeInt(int64(*value),true)
	case uint8:
		return encoder.EncodeInt(int64(value),true)
	case *uint16:
		return encoder.EncodeInt(int64(*value),true)
	case uint16:
		return encoder.EncodeInt(int64(value),true)
	case *uint32:
		return encoder.EncodeInt(int64(*value),true)
	case uint32:
		return encoder.EncodeInt(int64(value),true)
	case *uint64:
		return encoder.EncodeInt(int64(*value),true)
	case uint64:
		return encoder.EncodeInt(int64(value),true)
	case *float32:
		return encoder.EncodeFloat(*value,true)
	case *float64:
		return encoder.EncodeDouble(*value,true)
	case float32:
		return encoder.EncodeFloat(value,true)
	case float64:
		return encoder.EncodeDouble(value,true)
	case *bool:
		return encoder.EncodeBool(*value,true)
	case bool:
		return encoder.EncodeBool(value,true)
	case []byte,*[]byte:
		return errors.New("DataType Can't Encode")
	default:
		v := reflect.ValueOf(value)
		rv := Coders.GetRealValue(&v)
		if rv == nil{
			return errors.New("nil Can't Encode")
		}
		//首先判断这个V是否具备值编码器接口
		if v.Type().Implements(Coders.ValueCoderType){
			switch rv.Kind() {
			case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
				if rv.IsNil() {
					return errors.New("nil Can't Encode")
				}
			default:
				vcoder := v.Interface().(Coders.ValueCoder)
				return vcoder.Encode(encoder)
			}
		}
		switch rv.Kind(){
		case reflect.Struct:
			if rv.Type() == Coders.TimeType{
				return errors.New("DataType Can't Encode")
			}
			structFields := Coders.Fields(rv.Type())
			err := encoder.StartObject()
			if err != nil {
				return err
			}
			isFirst := true
			for i:=0;i<structFields.Len(jsonName);i++{
				f := structFields.Field(i)
				MarshName := f.MarshalName(jsonName)
				if MarshName != "-" {
					if isFirst{
						isFirst = false
					}else{
						encoder.w.WriteByte(',')
					}
					if err = encoder.EncodeKeyName(MarshName);err!=nil{
						return err
					}
					if err = f.EncodeValue(encoder,*rv);err!=nil{
						return err
					}
				}
			}
			return encoder.EndObject()
		default:
			return encoder.GetEncoderFunc(rv.Type())(encoder,*rv)
		}
	}
	return nil
}



func (encoder *JsonEncoder)StartObject()error  {
	return encoder.w.WriteByte('{')
}

func (encoder *JsonEncoder)EndObject()error  {
	return encoder.w.WriteByte('}')
}

func (encoder *JsonEncoder)StartArray()error  {
	return encoder.w.WriteByte('[')
}

func (encoder *JsonEncoder)EndArray()error  {
	return encoder.w.WriteByte(']')
}

func (encoder *JsonEncoder)EncodeKeyName(key string) error {
	err := encoder.w.WriteByte('"')
	if err != nil{
		return err
	}
	_,err = encoder.w.WriteString(key)
	if err != nil{
		return err
	}
	_,err = encoder.w.WriteString(`":`)
	return err
}

func (encoder *JsonEncoder)EncodeString(v string,isEndValue bool) error {
	err := encoder.w.WriteByte('"')
	if err != nil{
		return err
	}
	_,err = encoder.w.WriteString(v)
	if err != nil{
		return err
	}
	if !isEndValue{
		_,err = encoder.w.WriteString(`",`)
	}else{
		err = encoder.w.WriteByte('"')
	}
	return err
}

func (encoder *JsonEncoder)EncodeInt(v int64,isEndValue bool) error {
	_, err := encoder.w.WriteString(strconv.FormatInt(v,10))
	if err!=nil{
		return err
	}
	if !isEndValue{
		err = encoder.w.WriteByte(',')
	}
	return err
}

func (encoder *JsonEncoder)EncodeFloat(f float32,isEndValue bool) error {
	_,err := encoder.w.WriteString(strconv.FormatFloat(float64(f),'f','e',32))
	if err != nil{
		return err
	}
	if !isEndValue{
		err = encoder.w.WriteByte(',')
	}
	return err
}

func (encoder *JsonEncoder)EncodeDouble(d float64,isEndValue bool) error {
	_,err := encoder.w.WriteString(strconv.FormatFloat(d,'f','e',64))
	if err != nil{
		return err
	}
	if !isEndValue{
		err = encoder.w.WriteByte(',')
	}
	return err
}

func (encoder *JsonEncoder)EncodeBool(b bool,isEndValue bool) error {
	var err error
	if b{
		_, err = encoder.w.WriteString("true")
	}else{
		_, err = encoder.w.WriteString("false")
	}
	if err != nil{
		return err
	}
	if !isEndValue{
		err = encoder.w.WriteByte(',')
	}
	return err
}

func (encoder *JsonEncoder)EncodeNil(isEndValue bool) error {
	_,err := encoder.w.WriteString("null")
	if err != nil{
		return err
	}
	if !isEndValue{
		err = encoder.w.WriteByte(',')
	}
	return err
}

func (encoder *JsonEncoder)EncodeTime(dt time.Time,isEndValue bool) error {
	_,err := encoder.w.WriteString("/Date(")
	if err != nil{
		return err
	}
	_,err = encoder.w.WriteString(strconv.Itoa(int(dt.Unix())*1000))
	if err != nil{
		return err
	}
	_,err = encoder.w.WriteString(")/")
	return err
}


func (encoder *JsonEncoder)ReSet(w io.Writer)  {
	if buffer, ok := w.(*bytes.Buffer);ok{
		encoder.w = buffer
	}else{
		if encoder.w == nil{
			encoder.w = bufio.NewWriter(w)
		}else{
			encoder.w.(*bufio.Writer).Reset(w)
		}
	}
}

func NewEncoder(w io.Writer) *JsonEncoder {
	result := new(JsonEncoder)
	if buffer, ok := w.(*bytes.Buffer);ok{
		result.w = buffer
	}else{
		result.w = bufio.NewWriter(w)
	}
	return result
}

func Marshal(v interface{})([]byte,error) {
	var buffer bytes.Buffer
	coder := NewEncoder(&buffer)
	if err := coder.EncodeStand(v);err!=nil{
		return nil,err
	}
	return buffer.Bytes(),nil
}