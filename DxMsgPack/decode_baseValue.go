package DxMsgPack

import (
	"encoding/binary"
	"github.com/suiyunonghen/DxCommonLib"
	"unsafe"
	"time"
	"github.com/suiyunonghen/DxValue"
)

func (coder *MsgPackDecoder)readCode()(MsgPackCode,error)  {
	c, err := coder.r.ReadByte()
	if err != nil {
		return 0, err
	}
	return MsgPackCode(c), nil
}


func (coder *MsgPackDecoder)DecodeDateTime(code MsgPackCode)(DxCommonLib.TDateTime,error)  {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.readCode();err!=nil{
			return -1,err
		}
	}
	switch code {
	case CodeFloat:
		var b [4]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return -1,err
		}
		v32 := binary.BigEndian.Uint32(b[:])
		return DxCommonLib.TDateTime(*(*float32)(unsafe.Pointer(&v32))),nil
	case CodeDouble:
		var b [8]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return -1,err
		}
		v64 := binary.BigEndian.Uint64(b[:])
		return DxCommonLib.TDateTime(*(*float32)(unsafe.Pointer(&v64))),nil
	case CodeFixExt4:
		var b byte
		if b,err = coder.r.ReadByte();err!=nil{
			return -1,err
		}
		if int8(b) == -1{
			var b [4]byte
			if _,err := coder.r.Read(b[:]);err!=nil{
				return -1,err
			}
			ms := binary.BigEndian.Uint32(b[:])
			ntime := time.Now()
			ns := ntime.Unix()
			ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
			return DxCommonLib.Time2DelphiTime(&ntime),nil
		}
	}
	return -1,DxValue.ErrValueType
}

func (coder *MsgPackDecoder)DecodeInt(code MsgPackCode)(int64,error)  {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.readCode();err!=nil{
			return 0,err
		}
	}
	if code <= PosFixedNumHigh{
		return int64(code),nil
	}else if code >= 0xe0{
		return int64(int8(code)),nil
	}
	switch code {
	case CodeInt8,CodeUint8:
		if bt,err := coder.r.ReadByte();err!=nil{
			return 0,err
		}else if code == CodeInt8{
			return int64(int8(bt)),nil
		}else{
			return int64(bt),nil
		}
	case CodeInt16,CodeUint16:
		var b [2]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return 0,err
		}else{
			vuint16 := binary.BigEndian.Uint16(b[:])
			if code == CodeInt16{
				return int64(int16(vuint16)),nil
			}else{
				return int64(vuint16),nil
			}
		}
	case CodeInt32,CodeUint32:
		var b [4]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return 0,err
		}else{
			vuint32 := binary.BigEndian.Uint32(b[:])
			if code == CodeInt32{
				return int64(int32(vuint32)),nil
			}else{
				return int64(vuint32),nil
			}
		}
	case CodeInt64,CodeUint64:
		var b [8]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return 0,err
		}else{
			vuint64 := binary.BigEndian.Uint64(b[:])
			return int64(vuint64),nil
		}
	}
	return 0,DxValue.ErrValueType
}

func (coder *MsgPackDecoder)DecodeFloat(code MsgPackCode)(float64,error) {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.readCode();err!=nil{
			return 0,err
		}
	}
	switch code {
	case CodeFloat:
		var b [4]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return 0,err
		}
		v32 := binary.BigEndian.Uint32(b[:])
		return float64(*(*float32)(unsafe.Pointer(&v32))),nil
	case CodeDouble:
		var b [8]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return 0,err
		}
		v64 := binary.BigEndian.Uint64(b[:])
		return *(*float64)(unsafe.Pointer(&v64)),nil
	default:
		if iv,err := coder.DecodeInt(code);err!=nil{
			return 0,err
		}else{
			return float64(iv),nil
		}
	}
	return 0,DxValue.ErrValueType
}

func (coder *MsgPackDecoder)DecodeBinary(code MsgPackCode)([]byte,error)  {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.readCode();err!=nil{
			return nil,err
		}
	}
	btlen := 0
	switch code {
	case CodeBin8:
		if b,err := coder.r.ReadByte();err!=nil{
			return nil,err
		}else{
			btlen = int(b)
		}
	case CodeBin16:
		var b [2]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return nil,err
		}
		btlen =  int(binary.BigEndian.Uint16(b[:]))
	case CodeBin32:
		var b [4]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return nil,err
		}
		btlen =  int(binary.BigEndian.Uint32(b[:]))
	default:
		return nil,DxValue.ErrValueType
	}
	if btlen > 0{
		mb := make([]byte,btlen)
		if _,err := coder.r.Read(mb);err!=nil{
			return nil,err
		}
		return mb,nil
	}
	return nil,nil
}

func (coder *MsgPackDecoder)DecodeExtValue(code MsgPackCode)([]byte,error) {
	btlen := -1
	var err error
	if code == CodeUnkonw{
		if code,err = coder.readCode();err!=nil{
			return nil,err
		}
	}
	switch code {
	case CodeFixExt1: btlen = 2
	case CodeFixExt2: btlen = 3
	case CodeFixExt4: btlen = 5
	case CodeFixExt8: btlen = 9
	case CodeFixExt16: btlen = 17
	case CodeExt8:
		if blen,err := coder.r.ReadByte();err!=nil {
			return nil,err
		}else{
			btlen = int(blen)+1
		}
	case CodeExt16:
		var b [2]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return nil,err
		}
		btlen =  int(binary.BigEndian.Uint16(b[:]))+1
	case CodeExt32:
		var b [4]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return nil,err
		}
		btlen =  int(binary.BigEndian.Uint32(b[:]))+1
	}
	if btlen <= 0{
		return nil,DxValue.ErrInvalidateMsgPack
	}
	mb := make([]byte,btlen)
	if _,err = coder.r.Read(mb);err!=nil{
		return nil,err
	}
	return mb,nil
}

func (coder *MsgPackDecoder)DecodeMapLen(mapcode MsgPackCode)(int,error)  {
	var err error
	if mapcode == CodeUnkonw{
		if mapcode,err = coder.readCode();err!=nil{
			return 0,err
		}
	}
	switch mapcode {
	case CodeMap16:
		var b [2]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return 0,err
		}
		return int(binary.BigEndian.Uint16(b[:])),nil
	case CodeMap32:
		var b [4]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return 0,err
		}
		return int(binary.BigEndian.Uint32(b[:])),nil
	default:
		if mapcode >= CodeFixedMapLow && mapcode<=CodeFixedMapHigh{
			return int(mapcode & FixedMapMask),nil
		}
	}
	return 0,ErrInvalidateMap
}

func (coder *MsgPackDecoder)DecodeArrayLen(code MsgPackCode)(int,error)  {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.readCode();err!=nil{
			return 0,err
		}
	}
	switch code {
	case CodeArray16:
		var b [2]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return 0,err
		}
		return int(binary.BigEndian.Uint16(b[:])),nil
	case CodeArray32:
		var b [4]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return 0,err
		}
		return int(binary.BigEndian.Uint32(b[:])),nil
	default:
		if code >= CodeFixedArrayLow && code <= CodeFixedArrayHigh{
			return int(code & FixedArrayMask),nil
		}
	}
	return 0,ErrInvalidateArrLen
}

func (coder *MsgPackDecoder)DecodeString(code MsgPackCode)([]byte,error) {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.readCode();err!=nil{
			return nil,err
		}
	}
	strlen := 0
	switch code {
	case CodeStr8:
		if bl,err := coder.r.ReadByte();err!=nil{
			return nil,err
		}else{
			strlen = int(bl)
		}
	case CodeStr16:
		var b [2]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return nil,err
		}
		strlen = int(binary.BigEndian.Uint16(b[:]))
	case CodeStr32:
		var b [4]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return nil,err
		}
		strlen = int(binary.BigEndian.Uint32(b[:]))
	default:
		if code < 0xa0 || code> 0xbf {
			return nil,DxValue.ErrValueType
		}
		strlen = int(code & FixedStrMask)
	}
	if strlen > 0{
		mb := make([]byte,strlen)
		if _,err = coder.r.Read(mb);err!=nil{
			return nil,err
		}

		return mb,err
	}
	return nil,nil
}
