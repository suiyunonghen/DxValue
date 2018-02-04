package DxMsgPack

import (
	"reflect"
	"github.com/suiyunonghen/DxValue"
	"errors"
	"fmt"
	"github.com/suiyunonghen/DxCommonLib"
)

type decoderFunc func(*MsgPackDecoder, reflect.Value) error
var (
	interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
	stringType = reflect.TypeOf((*string)(nil)).Elem()
    mapStringStringPtrType = reflect.TypeOf((*map[string]string)(nil))
	mapStringStringType = mapStringStringPtrType.Elem()
	mapStringInterfacePtrType = reflect.TypeOf((*map[string]interface{})(nil))
	mapStringInterfaceType = mapStringInterfacePtrType.Elem()
	mapIntStringPtrType = reflect.TypeOf((*map[int]string)(nil))
	mapIntStringType = mapIntStringPtrType.Elem()

	mapIntInterfacePtrType = reflect.TypeOf((*map[int]interface{})(nil))
	mapIntInterfaceType = mapIntInterfacePtrType.Elem()


	valueDecoders []decoderFunc
	ErrFiledNotSetable = errors.New("The Field not Setable")
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	ErrStructKey = errors.New("Struct Can only use String Key")
)

func init() {
	valueDecoders = []decoderFunc{
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
}

func decodeBoolValue(coder *MsgPackDecoder, v reflect.Value) error {
	if code,err := coder.readCode();err!=nil{
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

func decodeInt64Value(coder *MsgPackDecoder, v reflect.Value) error {
	if !v.CanSet() {
		return ErrFiledNotSetable
	}
	n, err := coder.DecodeInt(CodeUnkonw)
	if err != nil {
		return err
	}
	v.SetInt(n)
	return nil
}

func decodeUint64Value(coder *MsgPackDecoder, v reflect.Value) error {
	if !v.CanSet() {
		return ErrFiledNotSetable
	}
	n, err := coder.DecodeInt(CodeUnkonw)
	if err != nil {
		return err
	}
	v.SetUint(uint64(n))
	return nil
}

func decodeFloat32Value(coder *MsgPackDecoder, v reflect.Value) error {
	if !v.CanSet() {
		return ErrFiledNotSetable
	}
	f, err := coder.DecodeFloat(CodeUnkonw)
	if err != nil {
		return err
	}
	v.SetFloat(float64(f))
	return nil
}

func decodeStringValue(coder *MsgPackDecoder, v reflect.Value) error {
	if !v.CanSet() {
		return ErrFiledNotSetable
	}
	s, err := coder.DecodeString(CodeUnkonw)
	if err != nil {
		return err
	}
	v.SetString(DxCommonLib.FastByte2String(s))
	return nil
}


func decodeStructValue(coder *MsgPackDecoder, v reflect.Value)(err error) {
	maplen := 0
	if maplen,err = coder.DecodeMapLen(CodeUnkonw);err!=nil{
		return err
	}
	strcode := CodeUnkonw
	//判断键值，是Int还是str
	if strcode,err = coder.readCode();err!=nil{
		return err
	}
	if !strcode.IsStr(){
		return ErrStructKey
	}
	return nil
	//获得Struct的结构
	/*if k,v,err := coder.decodeStrMapKvRecord(strcode);err!=nil{
		return nil
	}else{
		structs
		//
		for j := 1;j<maplen;j++{
			if k,v,err = coder.decodeStrMapKvRecord(CodeUnkonw);err!=nil{
				return err
			}
			//strMap[k] = v
		}
	}*/
}


func decodeFloat64Value(coder *MsgPackDecoder, v reflect.Value) error {
	if !v.CanSet() {
		return ErrFiledNotSetable
	}
	f, err := coder.DecodeFloat(CodeUnkonw)
	if err != nil {
		return err
	}
	v.SetFloat(f)
	return nil
}

func decodeUnsupportedValue(coder *MsgPackDecoder, v reflect.Value) error {
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

func decodeSliceValue(coder *MsgPackDecoder, v reflect.Value) error {
	n, err := coder.DecodeArrayLen(CodeUnkonw)
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
		if err := coder.DecodeValue(sv); err != nil {
			return err
		}
	}

	return nil
}

func decodeArrayValue(coder *MsgPackDecoder, v reflect.Value) error {
	n, err := coder.DecodeArrayLen(CodeUnkonw)
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
		if err := coder.DecodeValue(sv); err != nil {
			return err
		}
	}

	return nil
}


func decodeInterfaceValue(coder *MsgPackDecoder, v reflect.Value) error {
	if v.IsNil() {
		return coder.interfaceValue(v)
	}
	return coder.DecodeValue(v.Elem())
}

func decodeBytesValue(coder *MsgPackDecoder, v reflect.Value) error {
	if bt,err := coder.DecodeBinary(CodeUnkonw);err!=nil{
		return err
	}else{
		v.SetBytes(bt)
		return nil
	}
}


func decodeStringSliceValue(coder *MsgPackDecoder, v reflect.Value) error {
	if arrlen,err := coder.DecodeArrayLen(CodeUnkonw);err!=nil{
		return err
	}else if arrlen ==-1{
		return nil
	}else{
		ptr := v.Addr().Convert(sliceStringPtrType).Interface().(*[]string)
		ss := setStringsCap(*ptr,arrlen)
		for i := 0; i < arrlen; i++ {
			s, err := coder.DecodeString(CodeUnkonw)
			if err != nil {
				return err
			}
			ss = append(ss, DxCommonLib.FastByte2String(s))
		}
		*ptr = ss
		return nil
	}
}



func ptrDecoderFunc(typ reflect.Type) decoderFunc {
	decoder := getDecoder(typ.Elem())
	return func(coder *MsgPackDecoder, v reflect.Value) error {
		if coder.hasNilCode() {
			v.Set(reflect.Zero(v.Type()))
			coder.readCode()
		}
		if v.IsNil() {
			if !v.CanSet(){
				return ErrFiledNotSetable
			}
			v.Set(reflect.New(v.Type().Elem()))
		}
		return decoder(coder, v.Elem())
	}
}



func decodeByteArrayValue(coder *MsgPackDecoder, v reflect.Value) error {
	c, err := coder.readCode()
	if err != nil {
		return err
	}

	btlen,err := coder.BinaryLen(c)
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
	_,err = coder.r.Read(b)
	return err
}

func decodeMapStringStringValue(coder *MsgPackDecoder, v reflect.Value) error {
	mptr := v.Addr().Convert(mapStringStringPtrType).Interface().(*map[string]string)
	return coder.decodeStrValueMapFunc(mptr)
}

func decodeMapStringInterfaceValue(coder *MsgPackDecoder, v reflect.Value) error {
	ptr := v.Addr().Convert(mapStringInterfacePtrType).Interface().(*map[string]interface{})
	return coder.decodeStrMapFunc(ptr)
}

func decodeMapIntInterfaceValue(coder *MsgPackDecoder, v reflect.Value) error {
	ptr := v.Addr().Convert(mapIntInterfacePtrType).Interface().(*map[int]interface{})
	return coder.decodeIntKeyMapFunc(ptr)
}

func getDecoder(typ reflect.Type) decoderFunc {
	kind := typ.Kind()
	switch kind {
	case reflect.Ptr:
		return ptrDecoderFunc(typ)
	case reflect.Slice:
		elem := typ.Elem()
		switch elem.Kind() {
		case reflect.Uint8:
			return decodeBytesValue
		}
		switch elem {
		case stringType:
			return decodeStringSliceValue
		}
	case reflect.Array:
		if typ.Elem().Kind() == reflect.Uint8 {
			return decodeByteArrayValue
		}
	case reflect.Map:
		if typ.Key() == stringType {
			switch typ.Elem() {
			case stringType:
				return decodeMapStringStringValue
			case interfaceType:
				return decodeMapStringInterfaceValue
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
	typ := v.Type()
	if !v.CanSet(){
		return ErrCannotSet
	}
	switch typ.Kind() {
	case reflect.Bool:
		if code,err := coder.readCode();err!=nil{
			return err
		}else if code == CodeFalse{
			v.SetBool(false)
		}else if code == CodeTrue{
			v.SetBool(true)
		}
	case reflect.Struct:
	case reflect.Map:
	case reflect.Slice:
	default:
		if v.CanInterface(){
			vt := v.Interface()
			coder.DecodeStand(vt)
		}
	}
	return DxValue.ErrValueType
}
