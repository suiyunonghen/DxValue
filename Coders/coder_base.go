package Coders

import (
	"reflect"
	"errors"
)

type DecoderFunc func(Decoder, reflect.Value) error
type EncoderFunc func(Encoder,reflect.Value) error

type  Decoder  interface{
	DecodeStand(v interface{})(error)
	DecodeCustom()(error)
	//Decode(v *DxValue.DxBaseValue)(error)
	GetDecoderFunc(typ reflect.Type) DecoderFunc
	Skip()(error)
	Name()string
}

type  Encoder   interface{
	EncodeStand(v interface{})(error)
	EncodeCustom()(error)
	//Encode(v *DxValue.DxBaseValue)(error)
	Name()string
	GetEncoderFunc(typ reflect.Type)EncoderFunc
}

//值编码器
type  ValueCoder interface {
	Encode(encoder Encoder) error
	Decode(decoder Decoder) error
}

var(
	coders	[]string
	ErrValueType = errors.New("Value Data Type not Match")
)

func init()  {
	coders = make([]string,0,8)
}


func RegisterCoderName(coderName string)int  {
	for i := 0;i<len(coders);i++{
		if coders[i] == coderName{
			return i
		}
	}
	coders = append(coders,coderName)
	return len(coders)-1
}
