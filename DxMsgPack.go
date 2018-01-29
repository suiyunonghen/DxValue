/*
MsgPack的编码解码库
Autor: 不得闲
QQ:75492895
 */
package DxValue

import (
	"io"
	"unsafe"
	"github.com/suiyunonghen/DxCommonLib"
	"encoding/binary"
	"math"
	"time"
	"bufio"
)

type(
	DxMsgPackCoder struct {}
)

const(
	max_fixmap_len		= 16
	max_map16_len		= 1 << 16
	max_map32_len		= 1 << 32

	max_fixstr_len		= 32
	max_str8_len		= 1 << 8
	max_str16_len		= 1 << 16
	max_str32_len		= 1 << 32
)

func (coder DxMsgPackCoder)Encode(v *DxBaseValue, w io.Writer)(err error)  {
	switch v.fValueType {
	case DVT_Record:
		return coder.EncodeRecord((*DxRecord)(unsafe.Pointer(v)),w)
	}
	return nil
}

func (coder DxMsgPackCoder)EncodeRecord(r *DxRecord,w io.Writer)(err error)  {
	var writer *bufio.Writer
	ok := false
	if writer,ok = w.(*bufio.Writer);!ok{
		writer = bufio.NewWriter(w)
	}
	maplen := len(r.fRecords)
	if maplen <= max_fixmap_len{   //fixmap
		err = writer.WriteByte(0x80 | byte(maplen))
	}else if maplen <= max_map16_len{
		//写入长度
		mb := [3]byte{}
		mb[0] = 0xDE
		binary.BigEndian.PutUint16(mb[1:],uint16(maplen))
		_,err = writer.Write(mb[:])
	}else{
		if maplen > max_map32_len{
			maplen = max_map32_len
		}
		mb := [5]byte{}
		mb[0] = 0xDF
		binary.BigEndian.PutUint32(mb[1:],uint32(maplen))
		_,err = writer.Write(mb[:])
	}
	if err != nil{
		return
	}
	//写入对象信息,Kv对
	for k,v := range r.fRecords{
		if err = EncodeMsgPackString(k,writer);err != nil{
			return
		}
		//写入v
		if v != nil{
			switch v.fValueType {
			case DVT_Record:
				err = coder.EncodeRecord((*DxRecord)(unsafe.Pointer(v)),writer)
			case DVT_Int:
				err = EncodeMsgPackInt(int64((*DxIntValue)(unsafe.Pointer(v)).fvalue),writer)
			case DVT_Int32:
				err = EncodeMsgPackInt(int64((*DxInt32Value)(unsafe.Pointer(v)).fvalue),writer)
			case DVT_Int64:
				err = EncodeMsgPackInt((*DxInt64Value)(unsafe.Pointer(v)).fvalue,writer)
			case DVT_Bool:
				err = EncodeMsgPackBool((*DxBoolValue)(unsafe.Pointer(v)).fvalue,writer)
			case DVT_String:
				err = EncodeMsgPackString((*DxStringValue)(unsafe.Pointer(v)).fvalue,writer)
			case DVT_Float:
				err = EncodeMsgPackFloat((*DxFloatValue)(unsafe.Pointer(v)).fvalue,writer)
			case DVT_Double:
				err = EncodeMsgPackDouble((*DxDoubleValue)(unsafe.Pointer(v)).fvalue,writer)
			case DVT_Binary:
				if (*DxBinaryValue)(unsafe.Pointer(v)).fbinary != nil{
					err = EncodeMsgPackBinary((*DxBinaryValue)(unsafe.Pointer(v)).fbinary,writer)
				}else{
					writer.WriteByte(0xc0) //null
				}
			case DVT_Array:
			case DVT_DateTime:
				err = EncodeMsgPackDateTime(DxCommonLib.TDateTime((*DxDoubleValue)(unsafe.Pointer(v)).fvalue),writer)
			default:
				writer.WriteByte(0xc0) //null
			}
		}else{
			writer.WriteByte(0xc0) //null
		}

		if err != nil{
			return
		}
		maplen--
		if maplen == 0{
			break
		}
	}
	err = writer.Flush()
	return err
}

func EncodeMsgPackString(str string,w io.Writer)(err error)  {
	strbt := DxCommonLib.FastString2Byte(str)
	strlen := len(strbt)
	switch {
	case strlen <= max_fixstr_len:
		_,err = w.Write([]byte{0xA0 | byte(strlen)})
	case strlen <= max_str8_len:
		_,err = w.Write([]byte{0xD9, byte(strlen)})
	case strlen <= max_str16_len:
		var bt [3]byte
		bt[0] = 0xDA
		binary.BigEndian.PutUint16(bt[1:],uint16(strlen))
		_,err = w.Write(bt[:])
	default:
		if strlen > max_str32_len{
			strlen = max_str32_len
		}
		var bt [5]byte
		bt[0] = 0xDB
		binary.BigEndian.PutUint32(bt[1:],uint32(strlen))
		_,err = w.Write(bt[:])
		strbt = strbt[:strlen]
	}
	if err != nil || strlen == 0{
		return
	}
	_,err = w.Write(strbt)
	return
}

func EncodeMsgPackFloat(v float32,w io.Writer)(err error)  {
	var b [5]byte
	b[0] = 0xca
	binary.BigEndian.PutUint32(b[1:], *(*uint32)(unsafe.Pointer(&v)))
	_,err = w.Write(b[:])
	return
}

func EncodeMsgPackDouble(v float64,w io.Writer)(err error)  {
	var b [9]byte
	b[0] = 0xcb
	binary.BigEndian.PutUint64(b[1:], *(*uint64)(unsafe.Pointer(&v)))
	_,err = w.Write(b[:])
	return
}

func EncodeMsgPackBinary(bt []byte,w io.Writer)(err error)  {
	btlen := len(bt)
	switch {
	case btlen <= max_str8_len:
		_,err = w.Write([]byte{0xc4,byte(btlen)})
	case btlen <= max_str16_len:
		var mb [3]byte
		mb[0] = 0xc5
		binary.BigEndian.PutUint16(mb[1:],uint16(btlen))
		_,err = w.Write(mb[:])
	default:
		if btlen > max_str32_len{
			btlen = max_str32_len
		}
		var mb [5]byte
		mb[0] = 0xc6
		binary.BigEndian.PutUint32(mb[1:],uint32(btlen))
		_,err = w.Write(mb[:])
	}
	if err == nil && btlen > 0{
		_,err = w.Write(bt[:btlen])
	}
	return
}

func EncodeMsgPackInt(vint int64,w io.Writer)(err error)  {
	switch {
	case vint >= 0 && vint <= 0x7f:  //0XXXXXXX is 8-bit unsigned integer
		_,err = w.Write([]byte{byte(vint)})
	case vint >= -0x1f && vint < 0:  // 111YYYYY is 8-bit 5-bit negative integer
		_,err = w.Write([]byte{0xe0 | byte(-vint)})
	case vint >= 0 && vint <= 0xff:
		_,err = w.Write([]byte{0xcc,uint8(vint)})
	case vint >= 0 && vint <= 0xffff: //0xcd  |ZZZZZZZZ|ZZZZZZZZ
		var mb [3]byte
		mb[0] = 0xcd
		binary.BigEndian.PutUint16(mb[1:],uint16(vint))
		_,err = w.Write(mb[:])
	case vint >= 0 && vint <= 0xffffffff: //0xce  |ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ
		var mb [5]byte
		mb[0] = 0xce
		binary.BigEndian.PutUint32(mb[1:],uint32(vint))
		_,err = w.Write(mb[:])
	case uint64(vint) <= 0xffffffffffffffff: // 0xcf  |ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|
		var mb [9]byte
		mb[0] = 0xcf
		binary.BigEndian.PutUint64(mb[1:],uint64(vint))
		_,err = w.Write(mb[:])
	case vint >= math.MinInt8 && vint <= math.MaxInt8: //0xd0  |ZZZZZZZZ|
		_,err = w.Write([]byte{0xd0,uint8(vint)})
	case vint >= math.MinInt16 && vint <= math.MaxInt16: //0xd1  |ZZZZZZZZ|ZZZZZZZZ|
		var mb [3]byte
		mb[0] = 0xd1
		binary.BigEndian.PutUint16(mb[1:],uint16(vint))
		_,err = w.Write(mb[:])
	case vint >=  math.MinInt32 && vint <= math.MaxInt32: //0xd2  |ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ
		var mb [5]byte
		mb[0] = 0xd2
		binary.BigEndian.PutUint32(mb[1:],uint32(vint))
		_,err = w.Write(mb[:])
	default: //0xd3  |ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ|ZZZZZZZZ
		var mb [9]byte
		mb[0] = 0xd3
		binary.BigEndian.PutUint64(mb[1:],uint64(vint))
		_,err = w.Write(mb[:])
	}
	return
}

func EncodeMsgPackBool(v bool,w io.Writer)(err error)  {
	if v{
		_,err = w.Write([]byte{0xc3})
	}else{
		_,err = w.Write([]byte{0xc2})
	}
	return
}

func EncodeMsgPackTime(t time.Time,w io.Writer)(err error)  {
	var b [6]byte
	b[0] = 0xd6
	b[1] = 0xff
	binary.BigEndian.PutUint32(b[2:],uint32(t.Unix()))
	_,err = w.Write(b[:])
	return
}

func EncodeMsgPackDateTime(dt DxCommonLib.TDateTime,w io.Writer)(err error)  {
	var b [6]byte
	b[0] = 0xd6
	b[1] = 0xff
	binary.BigEndian.PutUint32(b[2:],uint32(dt.ToTime().Unix()))
	_,err = w.Write(b[:])
	return
}