package DxMsgPack

import (
	"github.com/suiyunonghen/DxValue"
	"bufio"
	"io"
	"time"
	"github.com/suiyunonghen/DxCommonLib"
)

func (coder *MsgPackDecoder)DecodeStand(r io.Reader, v interface{})(error)  {
	if bytebf,ok := r.(bufReader);ok {
		coder.r = bytebf
	}else {
		if bf,ok := r.(*bufio.Reader);ok{
			bf.Reset(r)
			coder.r = bf
		}else{
			coder.r = bufio.NewReader(r)
		}
	}
	switch value := v.(type) {
	case *string:
		if strbt,err := coder.DecodeString(CodeUnkonw);err!=nil{
			return err
		}else{
			*value = DxCommonLib.FastByte2String(strbt)
		}
	case *[]interface{}:
	case *time.Time:
		if dt,err := coder.DecodeDateTime(CodeUnkonw);err !=nil{
			return err
		}else{
			*value = dt.ToTime()
		}
	case *int8:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = int8(i64)
		}
	case *int16:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = int16(i64)
		}
	case *int32:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = int32(i64)
		}
	case *int64:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = i64
		}
	case *uint8:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = uint8(i64)
		}
	case *uint16:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = uint16(i64)
		}
	case *uint32:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = uint32(i64)
		}
	case *uint64:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = uint64(i64)
		}
	case *float32:
		if vf,err := coder.DecodeFloat(CodeUnkonw);err!=nil{
			return err
		}else{
			*value = float32(vf)
		}
	case *float64:
		if vf,err := coder.DecodeFloat(CodeUnkonw);err!=nil{
			return err
		}else{
			*value = vf
		}
	case *bool:
		if code,err := coder.readCode();err!=nil{
			return err
		}else if code == CodeFalse{
			*value = false
		}else if code == CodeTrue{
			*value = true
		}
	case *[]byte:
		if bt,err := coder.DecodeBinary(CodeUnkonw);err!=nil{
			return err
		}else{
			*value = bt
		}
	case *map[string]interface{}:
	case *map[int]interface{}:
	case *time.Duration:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			*value = time.Duration(i64)
		}
	}
	coder.r.UnreadByte()
	return DxValue.ErrValueType
}
