package DxMsgPack

import (
	"reflect"
	"github.com/suiyunonghen/DxValue"
	"errors"
	"fmt"
	"github.com/suiyunonghen/DxCommonLib"
	"github.com/suiyunonghen/DxValue/Coders"
)


var (
	valueDecoders []Coders.DecoderFunc
	ErrFiledNotSetable = errors.New("The Field not Setable")
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	ErrStructKey = errors.New("Struct Can only use String Key")
)

func init() {
	valueDecoders = []Coders.DecoderFunc{
		reflect.Bool:          decodeBoolValue,
		reflect.Int:           decodeInt64Value,
		reflect.Int8:          decodeInt64Value,
		reflect.Int16:         decodeInt64Value,
		reflect.Int32:         decodeInt64Value,
		reflect.Int64:         decodeInt64Value,
		reflect.Uint:          decodeUint64Value,
		reflect.Uint8:         decodeUint64Value,
		reflect.Uint16:        decodeUint64Value,
		reflect.Uint32:        decodeUint64Value,
		reflect.Uint64:        decodeUint64Value,
		reflect.Float32:       decodeFloat32Value,
		reflect.Float64:       decodeFloat64Value,
		reflect.Complex64:     decodeUnsupportedValue,
		reflect.Complex128:    decodeUnsupportedValue,
		reflect.Array:         decodeArrayValue,
		reflect.Chan:          decodeUnsupportedValue,
		reflect.Func:          decodeUnsupportedValue,
		reflect.Interface:     decodeInterfaceValue,
		reflect.Ptr:           decodeUnsupportedValue,
		reflect.Slice:         decodeSliceValue,
		reflect.String:        decodeStringValue,
		reflect.Struct:        decodeStructValue,
		reflect.UnsafePointer: decodeUnsupportedValue,
	}
	Coders.RegisterCoderName("msgpack")
}

func decodeBoolValue(coder Coders.Decoder, v reflect.Value) error {

	if code,err := coder.(*MsgPackDecoder).readCode();err!=nil{
		return err
	}else if code == CodeFalse{
		v.SetBool(false)
	}else if code == CodeTrue{
		v.SetBool(true)
	}else {
		return DxValue.ErrValueType
	}
	return nil
}

func decodeInt64Value(coder Coders.Decoder, v reflect.Value) error {
	if !v.CanSet() {
		return ErrFiledNotSetable
	}
	n, err := coder.(*MsgPackDecoder).DecodeInt(CodeUnkonw)
	if err != nil {
		return err
	}
	v.SetInt(n)
	return nil
}

func decodeUint64Value(coder Coders.Decoder, v reflect.Value) error {
	if !v.CanSet() {
		return ErrFiledNotSetable
	}
	n, err := coder.(*MsgPackDecoder).DecodeInt(CodeUnkonw)
	if err != nil {
		return err
	}
	v.SetUint(uint64(n))
	return nil
}

func decodeFloat32Value(coder Coders.Decoder, v reflect.Value) error {
	if !v.CanSet() {
		return ErrFiledNotSetable
	}
	f, err := coder.(*MsgPackDecoder).DecodeFloat(CodeUnkonw)
	if err != nil {
		return err
	}
	v.SetFloat(float64(f))
	return nil
}

func decodeStringValue(coder Coders.Decoder, v reflect.Value) error {
	if !v.CanSet() {
		return ErrFiledNotSetable
	}
	s, err := coder.(*MsgPackDecoder).DecodeString(CodeUnkonw)
	if err != nil {
		return err
	}
	v.SetString(DxCommonLib.FastByte2String(s))
	return nil
}


func decodeStructValue(coder Coders.Decoder, v reflect.Value)(err error) {
	maplen := 0
	decoder,_ := coder.(*MsgPackDecoder)
	if maplen,err = decoder.DecodeMapLen(CodeUnkonw);err!=nil{
		return err
	}
	strcode := CodeUnkonw
	//判断键值，是Int还是str
	if strcode,err = decoder.readCode();err!=nil{
		return err
	}
	if !strcode.IsStr(){
		return ErrStructKey
	}
	//获得Struct的结构
	fields := Coders.Fields(v.Type())
	for i := 0; i < maplen; i++ {
		name, err := decoder.DecodeString(strcode)
		if err != nil {
			return err
		}
		strcode = CodeUnkonw
		keyName := DxCommonLib.FastByte2String(name)
		if f := fields.FieldByName(coder.Name(),keyName); f != nil {
			if err := f.DecodeValue(coder, v); err != nil {
				return err
			}
		} else {
			if err := decoder.SkipByCode(strcode); err != nil {
				return err
			}
		}
	}
	return nil
}


func decodeFloat64Value(coder Coders.Decoder, v reflect.Value) error {
	if !v.CanSet() {
		return ErrFiledNotSetable
	}
	f, err := coder.(*MsgPackDecoder).DecodeFloat(CodeUnkonw)
	if err != nil {
		return err
	}
	v.SetFloat(f)
	return nil
}


func decodeTimeValue(coder Coders.Decoder, v reflect.Value) error {
	if !v.CanSet() {
		return ErrFiledNotSetable
	}
	if vt,err := coder.(*MsgPackDecoder).DecodeDateTime_Go(CodeUnkonw);err!=nil{
		return err
	}else{
		v.Set(reflect.ValueOf(vt))
	}
	return nil
}

func decodeUnsupportedValue(coder Coders.Decoder, v reflect.Value) error {
	return fmt.Errorf("msgpack: Decode(unsupported %s)", v.Type())
}

func growSliceValue(v reflect.Value, n int) reflect.Value {
	diff := n - v.Len()
	if diff > 256 {
		diff = 256
	}
	v = reflect.AppendSlice(v, reflect.MakeSlice(v.Type(), diff, diff))
	return v
}

func decodeSliceValue(coder Coders.Decoder, v reflect.Value) error {
	decoder,_ := coder.(*MsgPackDecoder)
	n, err := decoder.DecodeArrayLen(CodeUnkonw)
	if err != nil {
		return err
	}

	if n == -1 {
		v.Set(reflect.Zero(v.Type()))
		return nil
	}
	if n == 0 && v.IsNil() {
		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
		return nil
	}

	if v.Cap() >= n {
		v.Set(v.Slice(0, n))
	} else if v.Len() < v.Cap() {
		v.Set(v.Slice(0, v.Cap()))
	}

	for i := 0; i < n; i++ {
		if i >= v.Len() {
			v.Set(growSliceValue(v, n))
		}
		sv := v.Index(i)
		if err := decoder.DecodeValue(sv); err != nil {
			return err
		}
	}

	return nil
}

func decodeArrayValue(coder Coders.Decoder, v reflect.Value) error {
	decoder,_ := coder.(*MsgPackDecoder)
	n, err := decoder.DecodeArrayLen(CodeUnkonw)
	if err != nil {
		return err
	}

	if n == -1 {
		return nil
	}

	if n > v.Len() {
		return fmt.Errorf("%s len is %d, but msgpack has %d elements", v.Type(), v.Len(), n)
	}
	for i := 0; i < n; i++ {
		sv := v.Index(i)
		if err := decoder.DecodeValue(sv); err != nil {
			return err
		}
	}

	return nil
}


func decodeInterfaceValue(coder Coders.Decoder, v reflect.Value) error {
	if v.IsNil() {
		return coder.(*MsgPackDecoder).interfaceValue(v)
	}
	return coder.(*MsgPackDecoder).DecodeValue(v.Elem())
}

func decodeBytesValue(coder Coders.Decoder, v reflect.Value) error {
	if bt,err := coder.(*MsgPackDecoder).DecodeBinary(CodeUnkonw);err!=nil{
		return err
	}else{
		v.SetBytes(bt)
		return nil
	}
}


func decodeStringSliceValue(coder Coders.Decoder, v reflect.Value) error {
	decoder,_ := coder.(*MsgPackDecoder)
	if arrlen,err := decoder.DecodeArrayLen(CodeUnkonw);err!=nil{
		return err
	}else if arrlen ==-1{
		return nil
	}else{
		ptr := v.Addr().Convert(sliceStringPtrType).Interface().(*[]string)
		ss := setStringsCap(*ptr,arrlen)
		for i := 0; i < arrlen; i++ {
			s, err := decoder.DecodeString(CodeUnkonw)
			if err != nil {
				return err
			}
			ss = append(ss, DxCommonLib.FastByte2String(s))
		}
		*ptr = ss
		return nil
	}
}




func decodeByteArrayValue(coder Coders.Decoder, v reflect.Value) error {
	msgdecoder,_ := coder.(*MsgPackDecoder)
	c, err := msgdecoder.readCode()
	if err != nil {
		return err
	}

	btlen,err := msgdecoder.BinaryLen(c)
	if err != nil {
		return err
	}
	if btlen == -1 {
		return nil
	}
	if btlen > v.Len() {
		return fmt.Errorf("%s len is %d, but msgpack has %d elements", v.Type(), v.Len(), btlen)
	}

	b := v.Slice(0, btlen).Bytes()
	_,err = msgdecoder.r.Read(b)
	return err
}

func decodeMapStringStringValue(coder Coders.Decoder, v reflect.Value) error {
	mptr := v.Addr().Convert(Coders.MapStringStringPtrType).Interface().(*map[string]string)
	return coder.(*MsgPackDecoder).decodeStrValueMapFunc(mptr)
}

func decodeMapStringInterfaceValue(coder Coders.Decoder, v reflect.Value) error {
	ptr := v.Addr().Convert(Coders.MapStringInterfacePtrType).Interface().(*map[string]interface{})
	return coder.(*MsgPackDecoder).decodeStrMapFunc(ptr)
}

func decodeMapIntInterfaceValue(coder Coders.Decoder, v reflect.Value) error {
	ptr := v.Addr().Convert(Coders.MapIntInterfacePtrType).Interface().(*map[int]interface{})
	return coder.(*MsgPackDecoder).decodeIntKeyMapFunc(ptr)
}

func decodeMapInt64InterfaceValue(coder Coders.Decoder, v reflect.Value) error {
	ptr := v.Addr().Convert(Coders.MapIntInterfacePtrType).Interface().(*map[int64]interface{})
	return coder.(*MsgPackDecoder).decodeIntKeyMapFunc64(ptr)
}

func (coder *MsgPackDecoder)GetDecoderFunc(typ reflect.Type) Coders.DecoderFunc {
	kind := typ.Kind()
	switch kind {
	case reflect.Ptr:
		decoder := coder.GetDecoderFunc(typ.Elem())
		return func(coder Coders.Decoder, v reflect.Value) error {
			msgdecoder, _ := coder.(*MsgPackDecoder)
			if msgdecoder.hasNilCode() {
				v.Set(reflect.Zero(v.Type()))
				msgdecoder.readCode()
			}
			if v.IsNil() {
				if !v.CanSet(){
					return ErrFiledNotSetable
				}
				v.Set(reflect.New(v.Type().Elem()))
			}
			return decoder(msgdecoder, v.Elem())
		}
	case reflect.Slice:
		elem := typ.Elem()
		switch elem.Kind() {
		case reflect.Uint8:
			return decodeBytesValue
		}
		switch elem {
		case Coders.StringType:
			return decodeStringSliceValue
		}
	case reflect.Array:
		if typ.Elem().Kind() == reflect.Uint8 {
			return decodeByteArrayValue
		}
	case reflect.Struct:
		if typ == Coders.TimeType{
			return decodeTimeValue
		}
	case reflect.Map:
		switch typ.Key() {
		case Coders.StringType:
			switch typ.Elem() {
			case Coders.StringType:
				return decodeMapStringStringValue
			case Coders.InterfaceType:
				return decodeMapStringInterfaceValue
			}
		case Coders.IntType:
			switch typ.Elem() {
			case Coders.StringType:
				return decodeMapStringStringValue
			case Coders.InterfaceType:
				return decodeMapIntInterfaceValue
			}
		case Coders.Int64Type:
			switch typ.Elem() {
			case Coders.StringType:
			case Coders.InterfaceType:
				return decodeMapInt64InterfaceValue
			}
		}
	}
	return valueDecoders[kind]
}


func (coder *MsgPackDecoder) interfaceValue(v reflect.Value) error {
	vv, err := coder.Decode2Interface()
	if err != nil {
		return err
	}
	if vv != nil {
		if v.Type() == errorType {
			if vv, ok := vv.(string); ok {
				v.Set(reflect.ValueOf(errors.New(vv)))
				return nil
			}
		}

		v.Set(reflect.ValueOf(vv))
	}
	return nil
}

func (coder *MsgPackDecoder)DecodeValue(v reflect.Value)(error)  {
	if !v.CanSet(){
		return ErrCannotSet
	}
	typ := v.Type()
	return coder.GetDecoderFunc(typ)(coder,v)
}
