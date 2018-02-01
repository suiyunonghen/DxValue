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
	"errors"
)

type(
	DxMsgPackCoder struct {}
)

const(
	max_fixmap_len		= 15
	max_map16_len		= 1 << 16 - 1
	max_map32_len		= 1 << 32 - 1

	max_fixstr_len		= 32 - 1
	max_str8_len		= 1 << 8 - 1
	max_str16_len		= 1 << 16 - 1
	max_str32_len		= 1 << 32 - 1
)

var (
	ErrInvalidateMsgPack = errors.New("Is not a Validate MsgPack format")
	msgpackDecodefuncs = map[byte]func(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(value *DxBaseValue, err error){
		0xc0:decodeNil,
		0xc2:decodefalse,
		0xc3:decodetrue,
		0xcc:decodeUint8,
		0xcd:decodeUint16,
		0xce:decodeUint32,
		0xcf:decodeUint64,
		0xd0:decodeInt8,
		0xd1:decodeInt16,
		0xd2:decodeInt32,
		0xd3:decodeInt64,
		0xca:decodeFloat32,
		0xcb:decodeFloat64,
		0xd9:decodeStr8,
		0xda:decodeStr16,
		0xdb:decodeStr32,
		0xc4:decodeBin8,
		0xc5:decodeBin16,
		0xc6:decodeBin32,
		0xd6:decodeDateTime,
		0xdc:decodeArr16,
		0xdd:decodeArr32,
		0xde:decodeRecord16,
		0xdf:decodeRecord32,
	}
)
func (coder DxMsgPackCoder)Encode(v *DxBaseValue, w io.Writer)(err error)  {
	return EncodeMsgPackBaseValue(v,w)
}

func (coder DxMsgPackCoder)DecodeResult(r io.Reader)(*DxBaseValue,error)  {
	return DecodeMsgPack(r)
}


func decode2Record(r io.Reader,rec2Fill *DxBaseValue,recLen int)(error)  {
	if recLen > 0 {
		//先读取一个，判定Map类型
		var b [1]byte
		_,err := r.Read(b[:])
		if err != nil{
			return err
		}
		//字符串类型
		if rec2Fill.fValueType == DVT_Record && (b[0] == 0xd9 || b[0] == 0xda || b[0] == 0xdb || b[0] & 0xa0 == 0xa0) ||
		   rec2Fill.fValueType == DVT_RecordIntKey && (b[0] >= 0xcc && b[0] <= 0xcf || b[0] >= 0xd0 && b[0] <= 0xd3 || b[0] <= 0x7f){
			_,err = decodeStrRecord((*DxRecord)(unsafe.Pointer(rec2Fill)),r,b[0],recLen)
			return err
		}else{
			return ErrInvalidateMsgPack
		}
	}
	return nil
}

func (coder DxMsgPackCoder)Decode(r io.Reader,v *DxBaseValue)(err error)  {
	var (
		reader *bufio.Reader
		bt byte
	)
	ok := false
	if reader,ok = r.(*bufio.Reader);!ok{
		reader = bufio.NewReader(r)
	}
	bt,err = reader.ReadByte()
	if err != nil{
		return err
	}
	switch v.fValueType {
	case DVT_Ext:
		if bt >= 0xd4 && bt == 0xd8 || bt>= 0xc7 && bt <= 0xc9{
			if mb,err := decodeExtvalue(reader,bt);err!=nil{
				return err
			}else{
				vext := (*DxExtValue)(unsafe.Pointer(v))
				vext.fdata = mb[1:]
				vext.ExtType = mb[0]
				return nil
			}
		}
		return ErrValueType
	case DVT_Record,DVT_RecordIntKey:
		rlen := -1
		if bt >= 0x80 && bt<=0x8f{
			rlen = int(bt & 0xf)
		}else if bt == 0xde{
			clen := uint16(0)
			if err = binary.Read(reader,binary.BigEndian,&clen);err != nil{
				return err
			}
			rlen = int(clen)
		}else if bt == 0xdf{
			clen := uint32(0)
			if err = binary.Read(reader,binary.BigEndian,&clen);err != nil{
				return err
			}
			rlen = int(clen)
		}
		if rlen >= 0{
			v.ClearValue(false)
			return decode2Record(reader,v,rlen)
		}
		return ErrValueType
	case DVT_String:
		if bt >= 0xd9 && bt <= 0xdb || bt >= 0xa0 && bt <= 0xbf{
			if st,err := decodeString(reader,bt);err!=nil{
				return err
			}else{
				(*DxStringValue)(unsafe.Pointer(v)).fvalue = st
			}
		}else{
			return ErrValueType
		}
	case DVT_Array:
		arlen := -1
		if bt >= 0x90 && bt <= 0x9f{
			arlen = int(bt & 0xf)
		}else if bt == 0xdc{
			clen := uint16(0)
			if err = binary.Read(r,binary.BigEndian,&clen);err != nil{
				return err
			}
			arlen = int(clen)
		}else if bt == 0xdd{
			clen := uint32(0)
			if err = binary.Read(r,binary.BigEndian,&clen);err != nil{
				return err
			}
			arlen = int(clen)
		}
		if arlen > 0 {
			arr := (*DxArray)(unsafe.Pointer(v))
			arr.TruncateArray(arlen)
			if err = decodeArray(r,arr,int(bt));err!=nil{
				return err
			}
			return nil
		}
		return ErrValueType
	case DVT_Int32:
		if vt,err := decodeInt(reader,bt);err!=nil{
			return err
		}else{
			(*DxInt32Value)(unsafe.Pointer(v)).fvalue = int32(vt)
			return nil
		}
	case DVT_Int:
		if vt,err := decodeInt(reader,bt);err!=nil{
			return err
		}else{
			(*DxIntValue)(unsafe.Pointer(v)).fvalue = int(vt)
			return nil
		}
	case DVT_Int64:
		if vt,err := decodeInt(reader,bt);err!=nil{
			return err
		}else{
			(*DxInt64Value)(unsafe.Pointer(v)).fvalue = vt
			return nil
		}
	case DVT_DateTime:
		if bt != 0xd6{
			return ErrValueType
		}
		bt,err = reader.ReadByte()
		if err != nil{
			return err
		}
		if int8(bt) != -1{
			return ErrValueType
		}
		ms := uint32(0)
		err := binary.Read(r,binary.BigEndian,&ms)
		if err!= nil{
			return err
		}
		ntime := time.Now()
		ns := ntime.Unix()
		ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
		(*DxDoubleValue)(unsafe.Pointer(v)).fvalue =float64(DxCommonLib.Time2DelphiTime(&ntime))
		return nil
	}
	return ErrValueType
}



func EncodeMsgPackRecord(r *DxRecord,w io.Writer)(err error)  {
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
			err = EncodeMsgPackBaseValue(v,writer)
		}else{
			err = writer.WriteByte(0xc0) //null
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

func decodeNil(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(value *DxBaseValue,err error)  {
	if parentValue == nil{
		return nil,nil
	}else {
		return nil,ErrInvalidateMsgPack
	}
}

func decodefalse(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(value *DxBaseValue,err error)  {
	if parentValue == nil{
		var bf DxBoolValue
		bf.fValueType = DVT_Bool
		return &bf.DxBaseValue,nil
	}else {
		return nil,ErrInvalidateMsgPack
	}
}

func decodetrue(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(value *DxBaseValue,err error)  {
	if parentValue == nil{
		var bf DxBoolValue
		bf.fValueType = DVT_Bool
		bf.fvalue = true
		return &bf.DxBaseValue,nil
	}else {
		return nil,ErrInvalidateMsgPack
	}
}

func decodeUint8(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		var b [1]byte
		_,err := r.Read(b[:])
		if err != nil{
			return nil,err
		}
		var bv DxInt32Value
		bv.fValueType = DVT_Int32
		bv.fvalue = int32(b[0])
		return &bv.DxBaseValue,nil
	}else {
		return nil,ErrInvalidateMsgPack
	}
}

func decodeUint16(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		vuint16 := uint16(0)
		if err := binary.Read(r,binary.BigEndian,&vuint16);err!= nil{
			return nil,err
		}
		var bv DxInt32Value
		bv.fValueType = DVT_Int32
		bv.fvalue = int32(vuint16)
		return &bv.DxBaseValue,nil
	}else {
		return nil,ErrInvalidateMsgPack
	}
}

func decodeUint32(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		vuint32 := uint32(0)
		if err := binary.Read(r,binary.BigEndian,&vuint32);err!= nil{
			return nil,err
		}
		var bv DxInt64Value
		bv.fValueType = DVT_Int64
		bv.fvalue = int64(vuint32)
		return &bv.DxBaseValue,nil
	}else {
		return nil,ErrInvalidateMsgPack
	}
}

func decodeUint64(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		vuint64 := uint64(0)
		if err := binary.Read(r,binary.BigEndian,&vuint64);err!= nil{
			return nil,err
		}
		var bv DxInt64Value
		bv.fValueType = DVT_Int64
		bv.fvalue = int64(vuint64)
		return &bv.DxBaseValue,nil
	}else {
		return nil,ErrInvalidateMsgPack
	}
}

func decodeInt8(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		var b [1]byte
		_,err := r.Read(b[:])
		if err != nil{
			return nil,err
		}
		var bv DxInt32Value
		bv.fValueType = DVT_Int32
		bv.fvalue = int32(int8(b[0]))
		return &bv.DxBaseValue,nil
	}else {
		return nil,ErrInvalidateMsgPack
	}
}

func decodeInt16(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		vuint16 := int16(0)
		if err := binary.Read(r,binary.BigEndian,&vuint16);err!= nil{
			return nil,err
		}
		var bv DxInt32Value
		bv.fValueType = DVT_Int32
		bv.fvalue = int32(vuint16)
		return &bv.DxBaseValue,nil
	}else {
		return nil,ErrInvalidateMsgPack
	}
}

func decodeInt32(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		vuint32 := int32(0)
		if err := binary.Read(r,binary.BigEndian,&vuint32);err!= nil{
			return nil,err
		}
		var bv DxInt32Value
		bv.fValueType = DVT_Int32
		bv.fvalue = vuint32
		return &bv.DxBaseValue,nil
	}else {
		return nil,ErrInvalidateMsgPack
	}
}

func decodeInt64(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		vuint64 := int64(0)
		if err := binary.Read(r,binary.BigEndian,&vuint64);err!= nil{
			return nil,err
		}
		var bv DxInt64Value
		bv.fValueType = DVT_Int64
		bv.fvalue = vuint64
		return &bv.DxBaseValue,nil
	}else {
		return nil,ErrInvalidateMsgPack
	}
}

func decodeFloat32(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		v32 := int32(0)
		if err := binary.Read(r,binary.BigEndian,&v32);err!= nil{
			return nil,err
		}
		var bv DxFloatValue
		bv.fValueType = DVT_Float
		bv.fvalue = *(*float32)(unsafe.Pointer(&v32))
		return &bv.DxBaseValue,nil
	}else {
		return nil,ErrInvalidateMsgPack
	}
}

func decodeFloat64(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		v64 := int64(0)
		if err := binary.Read(r,binary.BigEndian,&v64);err!= nil{
			return nil,err
		}
		var bv DxDoubleValue
		bv.fValueType = DVT_Double
		bv.fvalue = *(*float64)(unsafe.Pointer(&v64))
		return &bv.DxBaseValue,nil
	}else {
		return nil,ErrInvalidateMsgPack
	}
}


func decodeStr8(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		var b [1]byte
		_,err := r.Read(b[:])
		if err!=nil{
			return nil,err
		}
		var bv DxStringValue
		if b[0] > 0{
			mb := make([]byte,int(b[0]))
			if _,err = r.Read(mb);err!=nil{
				return nil,err
			}
			bv.fvalue = DxCommonLib.FastByte2String(mb)
		}
		bv.fValueType = DVT_String
		return &bv.DxBaseValue,nil
	}else{
		return nil,ErrInvalidateMsgPack
	}
}

func decodeStr16(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		strlen := uint16(0)
		err := binary.Read(r,binary.BigEndian,&strlen)
		if err!= nil{
			return nil,err
		}
		var bv DxStringValue
		if strlen > 0{
			mb := make([]byte,int(strlen))
			if _,err = r.Read(mb);err!=nil{
				return nil,err
			}
			bv.fvalue = DxCommonLib.FastByte2String(mb)
		}
		bv.fValueType = DVT_String
		return &bv.DxBaseValue,nil
	}else{
		return nil,ErrInvalidateMsgPack
	}
}

func decodeStr32(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		strlen := uint32(0)
		err := binary.Read(r,binary.BigEndian,&strlen)
		if err!= nil{
			return nil,err
		}
		var bv DxStringValue
		if strlen > 0{
			mb := make([]byte,int(strlen))
			if _,err = r.Read(mb);err!=nil{
				return nil,err
			}
			bv.fvalue = DxCommonLib.FastByte2String(mb)
		}
		bv.fValueType = DVT_String
		return &bv.DxBaseValue,nil
	}else{
		return nil,ErrInvalidateMsgPack
	}
}

func decodeBin8(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		var b [1]byte
		_,err := r.Read(b[:])
		if err!=nil{
			return nil,err
		}
		var bv DxBinaryValue
		if b[0] > 0{
			mb := make([]byte,int(b[0]))
			if _,err = r.Read(mb);err!=nil{
				return nil,err
			}
			bv.fbinary = mb
		}
		bv.fValueType = DVT_Binary
		return &bv.DxBaseValue,nil
	}else{
		return nil,ErrInvalidateMsgPack
	}
}

func decodeBin16(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		strlen := uint16(0)
		err := binary.Read(r,binary.BigEndian,&strlen)
		if err!= nil{
			return nil,err
		}
		var bv DxBinaryValue
		if strlen > 0{
			mb := make([]byte,int(strlen))
			if _,err = r.Read(mb);err!=nil{
				return nil,err
			}
			bv.fbinary = mb
		}
		bv.fValueType = DVT_Binary
		return &bv.DxBaseValue,nil
	}else{
		return nil,ErrInvalidateMsgPack
	}
}

func decodeBin32(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		strlen := uint32(0)
		err := binary.Read(r,binary.BigEndian,&strlen)
		if err!= nil{
			return nil,err
		}
		var bv DxBinaryValue
		if strlen > 0{
			mb := make([]byte,int(strlen))
			if _,err = r.Read(mb);err!=nil{
				return nil,err
			}
			bv.fbinary = mb
		}
		bv.fValueType = DVT_Binary
		return &bv.DxBaseValue,nil
	}else{
		return nil,ErrInvalidateMsgPack
	}
}

func decodeDateTime(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	if parentValue == nil{
		var b [1]byte
		_,err := r.Read(b[:])
		if err!=nil{
			return nil,err
		}
		if int8(b[0]) == -1{
			ms := uint32(0)
			err := binary.Read(r,binary.BigEndian,&ms)
			if err!= nil{
				return nil,err
			}
			var bv DxDoubleValue
			ntime := time.Now()
			ns := ntime.Unix()
			ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
			bv.fvalue = float64(DxCommonLib.Time2DelphiTime(&ntime))
			bv.fValueType = DVT_DateTime
			return &bv.DxBaseValue,nil
		}else{
			var mb [4]byte
			if _,err = r.Read(mb[:]);err!=nil{
				return nil,err
			}
			var bv DxExtValue
			bv.fValueType = DVT_Ext
			bv.fdata = mb[1:]
			bv.ExtType = mb[0]
			return &bv.DxBaseValue,nil
		}

	}
	return nil,ErrInvalidateMsgPack
}

func decodeArray(r io.Reader,arr *DxArray,arrlen int)error  {
	var (
		bt [1]byte
		err error
	)
	for i := 0;i<arrlen;i++{
		_,err = r.Read(bt[:])
		if err != nil{
			return err
		}
		switch {
		case bt[0] <= 0x7f:
			arr.SetInt32(i,int32(bt[0]))
		case bt[0] >= 0x80 && bt[0]<=0x8f:
			rlen := int(bt[0] & 0xf)
			if rec,err := decodeRecord(r,nil,rlen);err!=nil{
				return err
			}else{
				arr.SetBaseValue(i,rec)
			}
		case bt[0] >= 0x90 && bt[0] <= 0x9f:
			bt[0] = bt[0] & 0xf
			carr := NewArray()
			carr.TruncateArray(int(bt[0]))
			if err = decodeArray(r,carr,int(bt[0]));err!=nil{
				return err
			}
			arr.SetArray(i,carr)
		case bt[0] >= 0xd9 && bt[0] <= 0xdb || bt[0] >= 0xa0 && bt[0] <= 0xbf:
			if st,err := decodeString(r,bt[0]);err!=nil{
				return err
			}else{
				arr.SetString(i,st)
			}
		case bt[0] >= 0xe0 && bt[0]<= 0xff:
			if vint,err := decodeInt(r,bt[0]);err!=nil{
				return err
			}else{
				arr.SetInt32(i,int32(vint))
			}
		case bt[0] == 0xc0:
			arr.SetNull(i)
		case bt[0] == 0xc1: //un used
		case bt[0] == 0xc2:
			arr.SetBool(i,false)
		case bt[0] == 0xc3:
			arr.SetBool(i,true)
		case bt[0] == 0xcc || bt[0] == 0xcd || bt[0] == 0xd0 || bt[0] == 0xd1 || bt[0] == 0xd2:
			if v64,err := decodeInt(r,bt[0]);err!=nil{
				return err
			}else{
				arr.SetInt32(i,int32(v64))
			}
		case bt[0] == 0xce || bt[0] == 0xcf || bt[0] == 0xd3:
			if v64,err := decodeInt(r,bt[0]);err!=nil{
				return err
			}else{
				arr.SetInt64(i,v64)
			}
		case bt[0] == 0xca:
			v32 := int32(0)
			if err := binary.Read(r,binary.BigEndian,&v32);err!= nil{
				return err
			}
			arr.SetFloat(i,*(*float32)(unsafe.Pointer(&v32)))
		case bt[0] == 0xcb:
			v64 := int64(0)
			if err := binary.Read(r,binary.BigEndian,&v64);err!= nil{
				return err
			}
			arr.SetDouble(i,*(*float64)(unsafe.Pointer(&v64)))
		case bt[0] >= 0xc4 && bt[0] <= 0xc6:
			if b,err :=decodeBinary(r,bt[0]);err!=nil{
				return err
			}else{
				arr.SetBinary(i,b)
			}
		case bt[0] == 0xd6:
			if _,err := r.Read(bt[:]);err!=nil{
				return err
			}
			if int8(bt[0]) == -1{
				ms := uint32(0)
				err := binary.Read(r,binary.BigEndian,&ms)
				if err!= nil{
					return err
				}
				ntime := time.Now()
				ns := ntime.Unix()
				ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
				arr.SetDateTime(i,DxCommonLib.Time2DelphiTime(&ntime))
			}else{
				var mb [4]byte
				if _,err = r.Read(mb[:]);err!=nil{
					return err
				}
				arr.SetExtValue(i,bt[0],mb[:]) //扩展类型
			}
		case bt[0] >= 0xd4 && bt[0] == 0xd8 || bt[0]>= 0xc7 && bt[0] <= 0xc9://扩展类型
			if mb,err := decodeExtvalue(r,bt[0]);err!=nil{
				return err
			}else{
				arr.SetExtValue(i,mb[0],mb[1:])
			}
		case bt[0] == 0xdc:
			clen := uint16(0)
			if err = binary.Read(r,binary.BigEndian,&clen);err != nil{
				return err
			}
			//解析Array
			carr := NewArray()
			carr.TruncateArray(int(clen))
			if err = decodeArray(r,carr,int(clen));err!=nil{
				return err
			}
			arr.SetArray(i,carr)
		case bt[0] == 0xdd:
			crecLen := uint32(0)
			if err = binary.Read(r,binary.BigEndian,&crecLen);err != nil{
				return err
			}
			//解析Array
			carr := NewArray()
			carr.TruncateArray(int(crecLen))
			if err = decodeArray(r,carr,int(crecLen));err!=nil{
				return err
			}
			arr.SetArray(i,carr)
		case bt[0] >= 0x90 && bt[0] <= 0x9f:
			bt[0] = bt[0] & 0xf
			carr := NewArray()
			carr.TruncateArray(int(bt[0]))
			if err = decodeArray(r,carr,int(bt[0]));err!=nil{
				return err
			}
			arr.SetArray(i,carr)
		default:
			return ErrInvalidateMsgPack
		}
	}
	return nil
}

func decodeArr16(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	arrlen := uint16(0)
	err := binary.Read(r,binary.BigEndian,&arrlen)
	if err!= nil{
		return nil,err
	}
	arr := NewArray()
	if err = decodeArray(r,arr,int(arrlen));err!=nil{
		return nil,err
	}
	return &arr.DxBaseValue,nil
}

func decodeArr32(r io.Reader,parentValue *DxBaseValue,strkey string,intInfo int64)(*DxBaseValue,error)  {
	arrlen := uint32(0)
	err := binary.Read(r,binary.BigEndian,&arrlen)
	if err!= nil{
		return nil,err
	}
	arr := NewArray()
	if err = decodeArray(r,arr,int(arrlen));err!=nil{
		return nil,err
	}
	return &arr.DxBaseValue,nil
}

func decodeBinary(r io.Reader,bflag byte)([]byte,error)  {
	btlen := 0
	if bflag == 0{
		var b [1]byte
		_,err := r.Read(b[:])
		if err!=nil{
			return nil,err
		}
		bflag = b[0]
	}
	switch bflag {
	case 0xc4:
		var b [1]byte
		_,err := r.Read(b[:])
		if err!=nil{
			return nil,err
		}
		btlen = int(b[0])
	case 0xc5:
		alen := uint16(0)
		err := binary.Read(r,binary.BigEndian,&alen)
		if err!= nil{
			return nil,err
		}
		btlen = int(alen)
	case 0xc6:
		alen := uint32(0)
		err := binary.Read(r,binary.BigEndian,&alen)
		if err!= nil{
			return nil,err
		}
		btlen = int(alen)
	default:
		return nil,ErrInvalidateMsgPack
	}
	if btlen > 0{
		mb := make([]byte,btlen)
		if _,err := r.Read(mb);err!=nil{
			return nil,err
		}
		return mb,nil
	}
	return nil,nil
}

func decodeString(r io.Reader,bflag byte)(string, error)  {
	strlen := 0
	if bflag == 0{
		var b [1]byte
		_,err := r.Read(b[:])
		if err!=nil{
			return "",err
		}
		bflag = b[0]
	}
	switch bflag {
	case 0xd9:
		var b [1]byte
		_,err := r.Read(b[:])
		if err!=nil{
			return "",err
		}
		strlen = int(b[0])
	case 0xda:
		alen := uint16(0)
		err := binary.Read(r,binary.BigEndian,&alen)
		if err!= nil{
			return "",err
		}
		strlen = int(alen)
	case 0xdb:
		alen := uint32(0)
		err := binary.Read(r,binary.BigEndian,&alen)
		if err!= nil{
			return "",err
		}
		strlen = int(alen)
	default:
		if bflag < 0xa0 || bflag> 0xbf {
			return "",ErrInvalidateMsgPack
		}
		strlen = int(bflag & 0x1f)
	}
	if strlen > 0{
		mb := make([]byte,strlen)
		if _,err := r.Read(mb);err!=nil{
			return "",err
		}
		return DxCommonLib.FastByte2String(mb),nil
	}
	return "",nil
}

func decodeInt(r io.Reader,bflag byte)(int64, error)  {
	var b [1]byte
	if bflag == 0{
		_,err := r.Read(b[:])
		if err!=nil{
			return 0,err
		}
		bflag = b[0]
	}
	switch bflag {
	case 0xcc,0xd0:
		if _,err := r.Read(b[:]);err!=nil{
			return 0,err
		}
		if bflag == 0xcc{
			return int64(b[0]),nil
		}
		return int64(int8(b[0])),nil
	case 0xcd,0xd1:
		vuint16 := uint16(0)
		if err := binary.Read(r,binary.BigEndian,&vuint16);err!= nil{
			return 0,err
		}
		if bflag == 0xcd{
			return int64(vuint16),nil
		}
		return int64(int16(vuint16)),nil
	case 0xce,0xd2:
		vuint32 := uint32(0)
		if err := binary.Read(r,binary.BigEndian,&vuint32);err!= nil{
			return 0,err
		}
		if bflag == 0xcd{
			return int64(vuint32),nil
		}
		return int64(int32(vuint32)),nil
	case 0xcf,0xd3:
		vuint64 := uint64(0)
		if err := binary.Read(r,binary.BigEndian,&vuint64);err!= nil{
			return 0,err
		}
		return int64(vuint64),nil
	default:
		if bflag <= 0x7f{
			return int64(bflag),nil

		}else if bflag > 0xe0 && bflag < 0xff{
			return -int64(bflag & 0x1f),nil
		} else if bflag & 0x80 == 0x80{
			return int64(bflag & 0xF),nil
		}else{
			return 0,ErrInvalidateMsgPack
		}
	}
}

func decodeExtvalue(r io.Reader,bflag byte)(mb []byte,err error){
	btlen := -1
	if bflag == 0{
		var b [1]byte
		_,err := r.Read(b[:])
		if err!=nil{
			return nil,err
		}
		bflag = b[0]
	}
	switch bflag {
	case 0xd4:
		btlen = 2
	case 0xd5:
		btlen = 3
	case 0xd7:
		btlen = 9
	case 0xd8:
		btlen = 17
	case 0xc7:
		var b [1]byte
		_,err := r.Read(b[:])
		if err!=nil{
			return nil,err
		}
		btlen = int(b[0])
	case 0xc8:
		u16 := uint16(0)
		if err = binary.Read(r,binary.BigEndian,&u16);err!=nil{
			return nil,err
		}
		btlen = int(u16)
	case 0xc9:
		u32 := uint32(0)
		if err = binary.Read(r,binary.BigEndian,&u32);err!=nil{
			return nil,err
		}
		btlen = int(u32)
	}
	mb = make([]byte,btlen)
	if _,err = r.Read(mb);err!=nil{
		return nil,err
	}
	return mb,nil
}


func decodeStrRecord(rec *DxRecord,r io.Reader,keyflag byte,recLen int)(*DxBaseValue,error)  {
	var b [1]byte
	keyName,err := decodeString(r,keyflag)
	if err != nil{
		return nil,err
	}
	_,err = r.Read(b[:])
	if err != nil{
		return nil,err
	}
	switch  {
	case b[0]>=0xd9 && b[0] <= 0xdb ||  b[0] >= 0xa0 && b[0] <= 0xbf:
		str,aerr := decodeString(r,b[0])
		if aerr != nil{
			return nil,aerr
		}
		rec.SetString(keyName,str)
	case b[0] >= 0xe0 && b[0] <= 0xff:
		rec.SetInt32(keyName, -int32((b[0] & 0x1f)))
	case b[0] >= 0xcc && b[0] <= 0xcf || b[0] >= 0xd0 && b[0] <= 0xd3 || b[0] <= 0x7f:
		if i64,aerr := decodeInt(r,b[0]);aerr != nil{
			return nil,aerr
		}else{
			rec.SetInt64(keyName,i64)
		}
	case b[0] == 0xca:
		v32 := int32(0)
		if err := binary.Read(r,binary.BigEndian,&v32);err!= nil{
			return nil,err
		}
		rec.SetFloat(keyName,*(*float32)(unsafe.Pointer(&v32)))
	case b[0] == 0xcb:
		v64 := int64(0)
		if err := binary.Read(r,binary.BigEndian,&v64);err!= nil{
			return nil,err
		}
		rec.SetDouble(keyName,*(*float64)(unsafe.Pointer(&v64)))
	case b[0] >= 0xc4 && b[0] <= 0xc6:
		mb,aerr := decodeBinary(r,b[0])
		if aerr != nil{
			return nil,aerr
		}
		rec.SetBinary(keyName,mb,true)
	case b[0] == 0xde:
		clen := uint16(0)
		if err = binary.Read(r,binary.BigEndian,&clen);err != nil{
			return nil,err
		}
		if recbase, err := decodeRecord(r,nil,int(clen));err!=nil{
			return nil,err
		}else{
			switch recbase.fValueType {
			case DVT_Record: rec.SetRecordValue(keyName,(*DxRecord)(unsafe.Pointer(recbase)))
			case DVT_RecordIntKey: rec.SetIntRecordValue(keyName,(*DxIntKeyRecord)(unsafe.Pointer(recbase)))
			}
		}
	case b[0] == 0xdf:
		crecLen := uint32(0)
		if err = binary.Read(r,binary.BigEndian,&crecLen);err != nil{
			return nil,err
		}
		if recbase, err := decodeRecord(r,nil,int(crecLen));err!=nil{
			return nil,err
		}else{
			switch recbase.fValueType {
			case DVT_Record: rec.SetRecordValue(keyName,(*DxRecord)(unsafe.Pointer(recbase)))
			case DVT_RecordIntKey: rec.SetIntRecordValue(keyName,(*DxIntKeyRecord)(unsafe.Pointer(recbase)))
			}
		}
	case b[0] >= 0x80 && b[0] <= 0x8f:
		rlen := int(b[0] & 0xf)
		if recbase, err := decodeRecord(r,nil,int(rlen));err!=nil{
			return nil,err
		}else{
			switch recbase.fValueType {
			case DVT_Record: rec.SetRecordValue(keyName,(*DxRecord)(unsafe.Pointer(recbase)))
			case DVT_RecordIntKey: rec.SetIntRecordValue(keyName,(*DxIntKeyRecord)(unsafe.Pointer(recbase)))
			}
		}
	case b[0] == 0xd6:
		if _,err = r.Read(b[:]);err!=nil{
			return nil,err
		}
		if int8(b[0]) == -1{
			ms := uint32(0)
			err := binary.Read(r,binary.BigEndian,&ms)
			if err!= nil{
				return nil,err
			}
			ntime := time.Now()
			ns := ntime.Unix()
			ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
			rec.SetDateTime(keyName,DxCommonLib.Time2DelphiTime(&ntime))
		}else {
			var mb [4]byte
			if _,err = r.Read(mb[:]);err!=nil{
				return nil,err
			}
			rec.SetExtValue(keyName,b[0],mb[:]) //扩展类型
		}
	case b[0] >= 0xd4 && b[0] == 0xd8 || b[0]>= 0xc7 && b[0] <= 0xc9://扩展类型
		if mb,err := decodeExtvalue(r,b[0]);err!=nil{
			return nil,err
		}else{
			rec.SetExtValue(keyName,mb[0],mb[1:])
		}
	case b[0] == 0xdc:
		clen := uint16(0)
		if err = binary.Read(r,binary.BigEndian,&clen);err != nil{
			return nil,err
		}
		//解析Array
		arr := NewArray()
		arr.TruncateArray(int(clen))
		if err = decodeArray(r,arr,int(clen));err!=nil{
			return nil,err
		}
		rec.SetArray(keyName,arr)
	case b[0] == 0xdd:
		crecLen := uint32(0)
		if err = binary.Read(r,binary.BigEndian,&crecLen);err != nil{
			return nil,err
		}
		//解析Array
		arr := NewArray()
		arr.TruncateArray(int(crecLen))
		if err = decodeArray(r,arr,int(crecLen));err!=nil{
			return nil,err
		}
		rec.SetArray(keyName,arr)
	case b[0] >= 0x90 && b[0] <= 0x9f:
		b[0] = b[0] & 0xf
		arr := NewArray()
		arr.TruncateArray(int(b[0]))
		if err = decodeArray(r,arr,int(b[0]));err!=nil{
			return nil,err
		}
		rec.SetArray(keyName,arr)
	}

	for i := 0;i<recLen-1;i++{
		if _,err := decodeRecord(r,&rec.DxBaseValue,1);err != nil{
			return nil,err
		}
	}
	return &rec.DxBaseValue,nil
}


func decodeIntKeyRecord(rec *DxIntKeyRecord,r io.Reader,keyflag byte,recLen int)(*DxBaseValue,error)  {
	var b [1]byte
	intKey,err := decodeInt(r,keyflag)
	if err != nil{
		return nil,err
	}
	_,err = r.Read(b[:])
	if err != nil{
		return nil,err
	}
	switch  {
	case b[0]>=0xd9 && b[0] <= 0xdb ||  b[0] >= 0xa0 && b[0] <= 0xbf:
		str,aerr := decodeString(r,b[0])
		if aerr != nil{
			return nil,aerr
		}
		rec.SetString(intKey,str)
	case b[0] >= 0xe0 && b[0] <= 0xff:
		rec.SetInt32(intKey, -int32((b[0] & 0x1f)))
	case b[0] >= 0xcc && b[0] <= 0xcf || b[0] >= 0xd0 && b[0] <= 0xd3 || b[0] <= 0x7f:
		if i64,aerr := decodeInt(r,b[0]);aerr != nil{
			return nil,aerr
		}else{
			rec.SetInt64(intKey,i64)
		}
	case b[0] == 0xca:
		v32 := int32(0)
		if err := binary.Read(r,binary.BigEndian,&v32);err!= nil{
			return nil,err
		}
		rec.SetFloat(intKey,*(*float32)(unsafe.Pointer(&v32)))
	case b[0] == 0xcb:
		v64 := int64(0)
		if err := binary.Read(r,binary.BigEndian,&v64);err!= nil{
			return nil,err
		}
		rec.SetDouble(intKey,*(*float64)(unsafe.Pointer(&v64)))
	case b[0] >= 0xc4 && b[0] <= 0xc6:
		mb,aerr := decodeBinary(r,b[0])
		if aerr != nil{
			return nil,aerr
		}
		rec.SetBinary(intKey,mb,true)
	case b[0] == 0xde:
		clen := uint16(0)
		if err = binary.Read(r,binary.BigEndian,&clen);err != nil{
			return nil,err
		}
		if recbase, err := decodeRecord(r,nil,int(clen));err!=nil{
			return nil,err
		}else{
			switch recbase.fValueType {
			case DVT_Record: rec.SetRecordValue(intKey,(*DxRecord)(unsafe.Pointer(recbase)))
			case DVT_RecordIntKey: rec.SetIntRecordValue(intKey,(*DxIntKeyRecord)(unsafe.Pointer(recbase)))
			}
		}
	case b[0] == 0xdf:
		crecLen := uint32(0)
		if err = binary.Read(r,binary.BigEndian,&crecLen);err != nil{
			return nil,err
		}
		if recbase, err := decodeRecord(r,nil,int(crecLen));err!=nil{
			return nil,err
		}else{
			switch recbase.fValueType {
			case DVT_Record: rec.SetRecordValue(intKey,(*DxRecord)(unsafe.Pointer(recbase)))
			case DVT_RecordIntKey: rec.SetIntRecordValue(intKey,(*DxIntKeyRecord)(unsafe.Pointer(recbase)))
			}
		}
	case b[0] >= 0x80 && b[0] <= 0x8f:
		rlen := int(b[0] & 0xf)
		if recbase, err := decodeRecord(r,nil,int(rlen));err!=nil{
			return nil,err
		}else{
			switch recbase.fValueType {
			case DVT_Record: rec.SetRecordValue(intKey,(*DxRecord)(unsafe.Pointer(recbase)))
			case DVT_RecordIntKey: rec.SetIntRecordValue(intKey,(*DxIntKeyRecord)(unsafe.Pointer(recbase)))
			}
		}
	case b[0] == 0xd6:
		if _,err = r.Read(b[:]);err!=nil{
			return nil,err
		}
		if int8(b[0]) == -1{
			ms := uint32(0)
			err := binary.Read(r,binary.BigEndian,&ms)
			if err!= nil{
				return nil,err
			}
			ntime := time.Now()
			ns := ntime.Unix()
			ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
			rec.SetDateTime(intKey,DxCommonLib.Time2DelphiTime(&ntime))
		}else {
			var mb [4]byte
			if _,err = r.Read(mb[:]);err!=nil{
				return nil,err
			}
			rec.SetExtValue(intKey,b[0],mb[:]) //扩展类型
		}
	case b[0] >= 0xd4 && b[0] == 0xd8 || b[0]>= 0xc7 && b[0] <= 0xc9://扩展类型
		if mb,err := decodeExtvalue(r,b[0]);err!=nil{
			return nil,err
		}else{
			rec.SetExtValue(intKey,mb[0],mb[1:])
		}
	case b[0] == 0xdc:
		clen := uint16(0)
		if err = binary.Read(r,binary.BigEndian,&clen);err != nil{
			return nil,err
		}
		//解析Array
		arr := NewArray()
		arr.TruncateArray(int(clen))
		if err = decodeArray(r,arr,int(clen));err!=nil{
			return nil,err
		}
		rec.SetArray(intKey,arr)
	case b[0] == 0xdd:
		crecLen := uint32(0)
		if err = binary.Read(r,binary.BigEndian,&crecLen);err != nil{
			return nil,err
		}
		//解析Array
		arr := NewArray()
		arr.TruncateArray(int(crecLen))
		if err = decodeArray(r,arr,int(crecLen));err!=nil{
			return nil,err
		}
		rec.SetArray(intKey,arr)
	case b[0] >= 0x90 && b[0] <= 0x9f:
		b[0] = b[0] & 0xf
		arr := NewArray()
		arr.TruncateArray(int(b[0]))
		if err = decodeArray(r,arr,int(b[0]));err!=nil{
			return nil,err
		}
		rec.SetArray(intKey,arr)
	}

	for i := 0;i<recLen-1;i++{
		if _,err := decodeRecord(r,&rec.DxBaseValue,1);err != nil{
			return nil,err
		}
	}
	return &rec.DxBaseValue,nil
}



func decodeRecord(r io.Reader,parentValue *DxBaseValue,recLen int)(*DxBaseValue, error)  {
	if recLen > 0 || parentValue != nil{
		//先读取一个，判定Map类型
		var b [1]byte
		_,err := r.Read(b[:])
		if err != nil{
			return nil,err
		}
		//字符串类型
		if b[0] == 0xd9 || b[0] == 0xda || b[0] == 0xdb || b[0] & 0xa0 == 0xa0{
			var rec *DxRecord
			if parentValue != nil{
				if parentValue.fValueType == DVT_Record{
					rec = (*DxRecord)(unsafe.Pointer(parentValue))
				}
			}
			if rec == nil{
				rec = NewRecord()
			}
			return decodeStrRecord(rec,r,b[0],recLen)
		}else if b[0] >= 0xcc && b[0] <= 0xcf || b[0] >= 0xd0 && b[0] <= 0xd3 || b[0] <= 0x7f{
			var rec *DxIntKeyRecord
			if parentValue != nil{
				if parentValue.fValueType == DVT_RecordIntKey{
					rec = (*DxIntKeyRecord)(unsafe.Pointer(parentValue))
				}
			}
			if rec == nil{
				rec = NewIntKeyRecord()
			}

			return decodeIntKeyRecord(rec,r,b[0],recLen)
		}else{
			return nil,ErrInvalidateMsgPack
		}
	}
	return nil,nil
}

func decodeRecord32(r io.Reader,parentValue *DxBaseValue,strkey string,intkey int64)(*DxBaseValue,error)  {
	recLen := uint32(0)
	err := binary.Read(r,binary.BigEndian,&recLen)
	if err!= nil{
		return nil,err
	}
	return decodeRecord(r,parentValue,int(recLen))
}

func decodeRecord16(r io.Reader,parentValue *DxBaseValue,strkey string,intkey int64)(*DxBaseValue,error)  {
	recLen := uint16(0)
	err := binary.Read(r,binary.BigEndian,&recLen)
	if err!= nil{
		return nil,err
	}
	return decodeRecord(r,parentValue,int(recLen))
}



func DecodeMsgPack(r io.Reader)(value *DxBaseValue, err error)  {
	var (
		reader *bufio.Reader
		bt byte
	)
	ok := false
	if reader,ok = r.(*bufio.Reader);!ok{
		reader = bufio.NewReader(r)
	}
	bt,err = reader.ReadByte()
	if err != nil{
		return nil,err
	}
	if decodefunc,ok := msgpackDecodefuncs[bt];ok{
		return decodefunc(reader,value,"",-1)
	}else{
		switch {
		case bt <= 0x7f:
			var v DxInt32Value
			v.fValueType = DVT_Int32
			v.fvalue = int32(bt)
			return &v.DxBaseValue,nil
		case bt >= 0x80 && bt<=0x8f:
			rlen := int(bt & 0xf)
			if rec,err := decodeRecord(reader,nil,rlen);err!=nil{
				return nil,err
			}else{
				return rec,nil
			}
		case bt >= 0x90 && bt <= 0x9f:
			bt = bt & 0xf
			arr := NewArray()
			arr.TruncateArray(int(bt))
			if err = decodeArray(r,arr,int(bt));err!=nil{
				return nil,err
			}
			return &arr.DxBaseValue,nil
		case bt >= 0xa0 && bt <= 0xbf:
			if st,err := decodeString(reader,bt);err!=nil{
				return nil,err
			}else{
				var v DxStringValue
				v.fvalue = st
				v.fValueType = DVT_String
				return &v.DxBaseValue,nil
			}
		case bt >= 0xd4 && bt == 0xd8 || bt>= 0xc7 && bt <= 0xc9://扩展类型
			if mb,err := decodeExtvalue(reader,bt);err!=nil{
				return nil,err
			}else{
				var v DxExtValue
				v.fValueType = DVT_Ext
				v.fdata = mb[1:]
				v.ExtType = mb[0]
				return &v.DxBaseValue,nil
			}
		case bt >= 0xe0 && bt<= 0xff:
			if vint,err := decodeInt(reader,bt);err!=nil{
				return nil,err
			}else{
				var v DxInt32Value
				v.fValueType = DVT_Int32
				v.fvalue = int32(vint)
				return &v.DxBaseValue,nil
			}
		}
	}
	return nil,ErrInvalidateMsgPack
}

func EncodeMsgPackRecordIntKey(r *DxIntKeyRecord,w io.Writer)(err error)  {
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
		if err = EncodeMsgPackInt(k,writer);err != nil{
			return
		}
		//写入v
		if v != nil{
			err = EncodeMsgPackBaseValue(v,writer)
		}else{
			err = writer.WriteByte(0xc0) //null
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

func EncodeMsgPackBaseValue(v *DxBaseValue,w io.Writer)(err error)  {
	switch v.fValueType {
	case DVT_Record:
		err = EncodeMsgPackRecord((*DxRecord)(unsafe.Pointer(v)),w)
	case DVT_RecordIntKey:
		err = EncodeMsgPackRecordIntKey((*DxIntKeyRecord)(unsafe.Pointer(v)),w)
	case DVT_Int:
		err = EncodeMsgPackInt(int64((*DxIntValue)(unsafe.Pointer(v)).fvalue),w)
	case DVT_Int32:
		err = EncodeMsgPackInt(int64((*DxInt32Value)(unsafe.Pointer(v)).fvalue),w)
	case DVT_Int64:
		err = EncodeMsgPackInt((*DxInt64Value)(unsafe.Pointer(v)).fvalue,w)
	case DVT_Bool:
		err = EncodeMsgPackBool((*DxBoolValue)(unsafe.Pointer(v)).fvalue,w)
	case DVT_String:
		err = EncodeMsgPackString((*DxStringValue)(unsafe.Pointer(v)).fvalue,w)
	case DVT_Float:
		err = EncodeMsgPackFloat((*DxFloatValue)(unsafe.Pointer(v)).fvalue,w)
	case DVT_Double:
		err = EncodeMsgPackDouble((*DxDoubleValue)(unsafe.Pointer(v)).fvalue,w)
	case DVT_Binary:
		if (*DxBinaryValue)(unsafe.Pointer(v)).fbinary != nil{
			err = EncodeMsgPackBinary((*DxBinaryValue)(unsafe.Pointer(v)).fbinary,w)
		}else{
			_,err = w.Write([]byte{0xc0}) //null
		}
	case DVT_Ext:
		err = EncodeMsgPackExtValue((*DxExtValue)(unsafe.Pointer(v)),w)
	case DVT_Array:
		err = EncodeMsgPackArray((*DxArray)(unsafe.Pointer(v)),w)
	case DVT_DateTime:
		err = EncodeMsgPackDateTime(DxCommonLib.TDateTime((*DxDoubleValue)(unsafe.Pointer(v)).fvalue),w)
	default:
		_,err =  w.Write([]byte{0xc0}) //null
	}
	return err
}

func EncodeMsgPackArray(arr *DxArray,w io.Writer)(err error)  {
	var writer *bufio.Writer
	ok := false
	if writer,ok = w.(*bufio.Writer);!ok{
		writer = bufio.NewWriter(w)
	}
	arlen := arr.Length()
	switch {
	case arlen < 16: //1001XXXX|    N objects
		_,err = writer.Write([]byte{0x90 | byte(arlen)})
	case arlen <= max_map16_len:  //0xdc  |YYYYYYYY|YYYYYYYY|    N objects
		var mb [3]byte
		mb[0] = 0xdc
		binary.BigEndian.PutUint16(mb[1:],uint16(arlen))
		_,err = writer.Write(mb[:])
	default:
		if arlen > max_map32_len{
			arlen = max_map32_len
		}
		var mb [5]byte
		mb[0] = 0xdd
		binary.BigEndian.PutUint32(mb[1:],uint32(arlen))
		_,err = writer.Write(mb[:])
	}
	for i := 0;i <= arlen - 1;i++{
		if arr.fValues[i] == nil{
			err = writer.WriteByte(0xc0) //null
		}else{
			err = EncodeMsgPackBaseValue(arr.fValues[i],writer)
		}
		if err != nil{
			return
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
	btlen := 0
	if bt != nil{
		btlen = len(bt)
	}
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


func EncodeMsgPackExtValue(v *DxExtValue,w io.Writer)(err error)  {
	btlen := 0
	if v.fdata == nil || len(v.fdata) == 0{
		return nil
	}

	btlen = len(v.fdata)
	bt := v.fdata
	switch {
	case btlen == 1:
		_,err = w.Write([]byte{0xd4,v.ExtType})
	case btlen == 2:
		_,err = w.Write([]byte{0xd5,v.ExtType})
	case btlen <= 4:
		_,err = w.Write([]byte{0xd6,v.ExtType})
		if btlen < 4{
			bt = append(bt,0)
			btlen = 4
		}
	case btlen <= 8:
		_,err = w.Write([]byte{0xd7,v.ExtType})
		if btlen < 8{
			bt = append(bt,make([]byte,8-btlen)...)
			btlen = 8
		}
	case btlen <= 16:
		_,err = w.Write([]byte{0xd8,v.ExtType})
		if btlen < 16{
			bt = append(bt,make([]byte,16-btlen)...)
			btlen = 16
		}
	case btlen <= max_str8_len:
		_,err = w.Write([]byte{0xc7,byte(btlen),v.ExtType})
	case btlen <= max_str16_len:
		var mb [4]byte
		mb[0] = 0xc8
		binary.BigEndian.PutUint16(mb[1:],uint16(btlen))
		mb[3] = v.ExtType
		_,err = w.Write(mb[:])
	default:
		if btlen > max_str32_len{
			btlen = max_str32_len
		}
		var mb [6]byte
		mb[0] = 0xc6
		binary.BigEndian.PutUint32(mb[1:],uint32(btlen))
		mb[5] = v.ExtType
		_,err = w.Write(mb[:])
	}
	if err == nil && btlen > 0{
		_,err = w.Write(bt[:btlen])
	}
	return
}
