package DxMsgPack

import (
	"io"
	"github.com/suiyunonghen/DxCommonLib"
	"encoding/binary"
	"unsafe"
	"math"
	"time"
)

const(
	Max_fixmap_len		= 15
	Max_map16_len		= 1 << 16 - 1
	Max_map32_len		= 1 << 32 - 1

	Max_fixstr_len		= 32 - 1
	Max_str8_len		= 1 << 8 - 1
	Max_str16_len		= 1 << 16 - 1
	Max_str32_len		= 1 << 32 - 1
)


type  MsgPackEncoder  struct{
	w   io.Writer
	buf	[]byte
}

func (encoder *MsgPackEncoder)WriteByte(b byte)(error)  {
	encoder.buf[0] = b
	_, err := encoder.w.Write(encoder.buf[:1])
	return err
}

func (encoder *MsgPackEncoder)WriteUint16(u16 uint16,bytecode MsgPackCode)(err error)  {
	idx := 0
	if bytecode != CodeUnkonw{
		encoder.buf[0] = byte(bytecode)
		idx = 1
	}
	binary.BigEndian.PutUint16(encoder.buf[idx:],u16)
	_,err = encoder.w.Write(encoder.buf[:idx+2])
	return err
}

func (encoder *MsgPackEncoder)Name()string{
	return msgPackName
}

func (encoder *MsgPackEncoder)WriteUint32(u32 uint32,bytecode MsgPackCode)(err error)  {
	idx := 0
	if bytecode != CodeUnkonw{
		encoder.buf[0] = byte(bytecode)
		idx = 1
	}
	binary.BigEndian.PutUint32(encoder.buf[idx:],u32)
	_,err = encoder.w.Write(encoder.buf[:idx+4])
	return err
}

func (encoder *MsgPackEncoder)WriteUint64(u64 uint64,bytecode MsgPackCode)(err error)  {
	idx := 0
	if bytecode != CodeUnkonw{
		encoder.buf[0] = byte(bytecode)
		idx = 1
	}
	binary.BigEndian.PutUint64(encoder.buf[idx:],u64)
	_,err = encoder.w.Write(encoder.buf[:idx+8])
	return err
}

func (encoder *MsgPackEncoder)EncodeString(str string)(err error)  {
	strbt := DxCommonLib.FastString2Byte(str)
	strlen := len(strbt)
	switch {
	case strlen <= Max_fixstr_len:
		encoder.buf[0] = byte(CodeFixedStrLow) | byte(strlen)
		_,err = encoder.w.Write(encoder.buf[:1])
	case strlen <= Max_str8_len:
		encoder.buf[0] = byte(CodeStr8)
		encoder.buf[1] = byte(strlen)
		_,err = encoder.w.Write(encoder.buf[:2])
	case strlen <= Max_str16_len:
		err = encoder.WriteUint16(uint16(strlen),CodeStr16)
	default:
		if strlen > Max_str32_len{
			strlen = Max_str32_len
		}
		err = encoder.WriteUint32(uint32(strlen),CodeStr32)
		strbt = strbt[:strlen]
	}
	if err != nil || strlen == 0{
		return
	}
	_,err = encoder.w.Write(strbt)
	return err
}

func (encoder *MsgPackEncoder)EncodeFloat(v float32)(err error)  {
	return encoder.WriteUint32(*(*uint32)(unsafe.Pointer(&v)),CodeFloat)
}

func (encoder *MsgPackEncoder)EncodeDouble(v float64)(err error)  {
	return encoder.WriteUint64(*(*uint64)(unsafe.Pointer(&v)),CodeDouble)
}

func (encoder *MsgPackEncoder)Write(b []byte)error  {
	_,err := encoder.w.Write(b)
	return err
}

func (encoder *MsgPackEncoder)encodeIntMapFunc(vmap *map[int]interface{})(err error)  {
	if vmap == nil{
		return encoder.WriteByte(byte(CodeNil))
	}
	maplen := len(*vmap)
	if maplen <= Max_fixmap_len{   //fixmap
		err = encoder.WriteByte(0x80 | byte(maplen))
	}else if maplen <= Max_map16_len{
		//写入长度
		err = encoder.WriteUint16(uint16(maplen),CodeMap16)
	}else{
		if maplen > Max_map32_len{
			maplen = Max_map32_len
		}
		err = encoder.WriteUint32(uint32(maplen),CodeMap32)
	}
	if err != nil{
		return
	}
	//写入对象信息,Kv对
	for k,v := range *vmap{
		if err = encoder.EncodeInt(int64(k));err != nil{
			return err
		}
		//写入v
		if v != nil{
			err = encoder.EncodeStand(v)
		}else{
			err = encoder.WriteByte(0xc0) //null
		}
		if err!=nil{
			return
		}
	}
	return nil
}

func (encoder *MsgPackEncoder)encodeInt64MapFunc(vmap *map[int64]interface{})(err error)  {
	if vmap == nil{
		return encoder.WriteByte(byte(CodeNil))
	}
	maplen := len(*vmap)
	if maplen <= Max_fixmap_len{   //fixmap
		err = encoder.WriteByte(0x80 | byte(maplen))
	}else if maplen <= Max_map16_len{
		//写入长度
		err = encoder.WriteUint16(uint16(maplen),CodeMap16)
	}else{
		if maplen > Max_map32_len{
			maplen = Max_map32_len
		}
		err = encoder.WriteUint32(uint32(maplen),CodeMap32)
	}
	if err != nil{
		return
	}
	//写入对象信息,Kv对
	for k,v := range *vmap{
		if err = encoder.EncodeInt(k);err != nil{
			return err
		}
		//写入v
		if v != nil{
			err = encoder.EncodeStand(v)
		}else{
			err = encoder.WriteByte(0xc0) //null
		}
		if err!=nil{
			return
		}
	}
	return nil
}

func (encoder *MsgPackEncoder)encodeStrMapFunc(vmap *map[string]interface{})(err error)  {
	if vmap == nil{
		return encoder.WriteByte(byte(CodeNil))
	}
	maplen := len(*vmap)
	if maplen <= Max_fixmap_len{   //fixmap
		err = encoder.WriteByte(0x80 | byte(maplen))
	}else if maplen <= Max_map16_len{
		//写入长度
		err = encoder.WriteUint16(uint16(maplen),CodeMap16)
	}else{
		if maplen > Max_map32_len{
			maplen = Max_map32_len
		}
		err = encoder.WriteUint32(uint32(maplen),CodeMap32)
	}
	if err != nil{
		return
	}
	//写入对象信息,Kv对
	for k,v := range *vmap{
		if err = encoder.EncodeString(k);err != nil{
			return err
		}
		//写入v
		if v != nil{
			err = encoder.EncodeStand(v)
		}else{
			err = encoder.WriteByte(0xc0) //null
		}
		if err!=nil{
			return
		}
	}
	return nil
}

func (encoder *MsgPackEncoder)encodeStrStrMapFunc(vmap *map[string]string)(err error)  {
	if vmap == nil{
		return encoder.WriteByte(byte(CodeNil))
	}
	maplen := len(*vmap)
	if maplen <= Max_fixmap_len{   //fixmap
		err = encoder.WriteByte(0x80 | byte(maplen))
	}else if maplen <= Max_map16_len{
		//写入长度
		err = encoder.WriteUint16(uint16(maplen),CodeMap16)
	}else{
		if maplen > Max_map32_len{
			maplen = Max_map32_len
		}
		err = encoder.WriteUint32(uint32(maplen),CodeMap32)
	}
	if err != nil{
		return
	}
	//写入对象信息,Kv对
	for k,v := range *vmap{
		if err = encoder.EncodeString(k);err != nil{
			return err
		}
		//写入v
		err = encoder.EncodeString(v)
		if err!=nil{
			return
		}
	}
	return nil
}

func (encoder *MsgPackEncoder)EncodeBinary(bt []byte)(err error) {
	btlen := 0
	if bt != nil{
		btlen = len(bt)
	}
	switch {
	case btlen <= Max_str8_len:
		encoder.buf[0] = byte(0xc4)
		encoder.buf[1] = byte(btlen)
		_,err = encoder.w.Write(encoder.buf[:2])
	case btlen <= Max_str16_len:
		err = encoder.WriteUint16(uint16(btlen),CodeBin16)
	default:
		if btlen > Max_str32_len{
			btlen = Max_str32_len
		}
		err = encoder.WriteUint32(uint32(btlen),CodeBin32)
	}
	if err == nil && btlen > 0{
		_,err = encoder.w.Write(bt[:btlen])
	}
	return err
}

func (encoder *MsgPackEncoder)EncodeInt(vint int64)(err error)  {
	switch {
	case vint >= 0 && vint <= 0x7f:  //0XXXXXXX is 8-bit unsigned integer
		encoder.buf[0] = byte(vint)
		_,err = encoder.w.Write(encoder.buf[:1])
	case vint >= -32 && vint < 0:  // 111YYYYY is 8-bit 5-bit negative integer
		encoder.buf[0] = byte(NegFixedNumLow)
		encoder.buf[1] = byte(vint)
		_,err = encoder.w.Write(encoder.buf[:2])
	case vint >= 0 && vint <= 0xff:
		encoder.buf[0] = byte(CodeUint8)
		encoder.buf[1] = byte(vint)
		_,err = encoder.w.Write(encoder.buf[:2])
	case vint >= 0 && vint <= 0xffff: //0xcd  |ZZZZZZZZ|ZZZZZZZZ
		return encoder.WriteUint16(uint16(vint),CodeUint16)
	case vint >= 0 && vint <= 0xffffffff: //0xce  |ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ
		return encoder.WriteUint32(uint32(vint),CodeUint32)
	case uint64(vint) <= 0xffffffffffffffff: // 0xcf  |ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|
		return encoder.WriteUint64(uint64(vint),CodeUint64)
	case vint >= math.MinInt8 && vint <= math.MaxInt8: //0xd0  |ZZZZZZZZ|
		encoder.buf[0] = byte(CodeInt8)
		encoder.buf[1] = byte(vint)
		_,err = encoder.w.Write(encoder.buf[:2])
	case vint >= math.MinInt16 && vint <= math.MaxInt16: //0xd1  |ZZZZZZZZ|ZZZZZZZZ|
		return encoder.WriteUint16(uint16(vint),CodeInt16)
	case vint >=  math.MinInt32 && vint <= math.MaxInt32: //0xd2  |ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ
		return encoder.WriteUint32(uint32(vint),CodeInt32)
	default: //0xd3  |ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ
		return encoder.WriteUint64(uint64(vint),CodeInt64)
	}
	return
}

func (encoder *MsgPackEncoder)EncodeBool(v bool)(err error)  {
	if v{
		encoder.buf[0] = byte(CodeTrue)
	}else{
		encoder.buf[0] = byte(CodeFalse)
	}
	_,err = encoder.w.Write(encoder.buf[:1])
	return
}


func (encoder *MsgPackEncoder)EncodeTime(t time.Time)(err error)  {
	encoder.buf[0] = 0xd6
	encoder.buf[1] = 0xff
	binary.BigEndian.PutUint32(encoder.buf[2:],uint32(t.Unix()))
	_,err = encoder.w.Write(encoder.buf[:6])
	return
}

func (encoder *MsgPackEncoder)EncodeDateTime(dt DxCommonLib.TDateTime)(err error)  {
	encoder.buf[0] = 0xd6
	encoder.buf[1] = 0xff
	binary.BigEndian.PutUint32(encoder.buf[2:],uint32(dt.ToTime().Unix()))
	_,err = encoder.w.Write(encoder.buf[:6])
	return
}

func (encoder *MsgPackEncoder)Buffer()[]byte {
	if encoder.buf == nil{
		encoder.buf = make([]byte,9)
	}
	return encoder.buf
}

func (encoder *MsgPackEncoder)ReSet(w io.Writer)  {
	encoder.w = w
}

func NewEncoder(w io.Writer) *MsgPackEncoder {
	return &MsgPackEncoder{
		w:   w,
		buf: make([]byte, 9),
	}
}