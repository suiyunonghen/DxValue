package Coders

import (
	"github.com/suiyunonghen/DxValue"
	"reflect"
)

type DecoderFunc func(Decoder, reflect.Value) error
//type EncoderFunc func(Encoder)

type  Decoder  interface{
	DecodeStand(v interface{})(error)
	DecodeCustom()(error)
	Decode(v *DxValue.DxBaseValue)(error)
	GetDecoderFunc(typ reflect.Type) DecoderFunc
	Skip()(error)
	Name()string
}

type  Encoder   interface{
	EncodeStand(v interface{})(error)
	EncodeCustom()(error)
	Encode(v *DxValue.DxBaseValue)(error)
	Name()string
}



var(
	coders	[]string
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
