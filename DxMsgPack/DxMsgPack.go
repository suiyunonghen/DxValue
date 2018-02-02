package DxMsgPack

import (
	"bufio"
	"github.com/suiyunonghen/DxValue"
	"io"
	"github.com/suiyunonghen/DxCommonLib"
	"encoding/binary"
	"time"
	"unsafe"
)

type MsgPackDecoder   struct{
	r  *bufio.Reader
}

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
		v32 := int32(0)
		if err := binary.Read(coder.r,binary.BigEndian,&v32);err!= nil{
			return -1,err
		}
		return DxCommonLib.TDateTime(*(*float32)(unsafe.Pointer(&v32))),nil
	case CodeDouble:
		v64 := int64(0)
		if err := binary.Read(coder.r,binary.BigEndian,&v64);err!= nil{
			return -1,err
		}
		return DxCommonLib.TDateTime(*(*float32)(unsafe.Pointer(&v64))),nil
	case CodeFixExt4:
		var b byte
		if b,err = coder.r.ReadByte();err!=nil{
			return -1,err
		}
		if int8(b) == -1{
			ms := uint32(0)
			err := binary.Read(coder.r,binary.BigEndian,&ms)
			if err!= nil{
				return -1,err
			}
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
		vuint16 := uint16(0)
		if err := binary.Read(coder.r,binary.BigEndian,&vuint16);err!= nil{
			return 0,err
		}else if code == CodeInt16{
			return int64(int16(vuint16)),nil
		}else{
			return int64(vuint16),nil
		}
	case CodeInt32,CodeUint32:
		vuint32 := uint32(0)
		if err := binary.Read(coder.r,binary.BigEndian,&vuint32);err!= nil{
			return 0,err
		}else if code == CodeInt32{
			return int64(int32(vuint32)),nil
		}else{
			return int64(vuint32),nil
		}
	case CodeInt64,CodeUint64:
		vuint64 := uint64(0)
		if err := binary.Read(coder.r,binary.BigEndian,&vuint64);err!= nil{
			return 0,err
		}else{
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
		v32 := int32(0)
		if err := binary.Read(coder.r,binary.BigEndian,&v32);err!= nil{
			return 0,err
		}
		return float64(*(*float32)(unsafe.Pointer(&v32))),nil
	case CodeDouble:
		v64 := int64(0)
		if err := binary.Read(coder.r,binary.BigEndian,&v64);err!= nil{
			return 0,err
		}
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
		alen := uint16(0)
		if err := binary.Read(coder.r,binary.BigEndian,&alen);err!=nil{
			return nil,err
		}else{
			btlen = int(alen)
		}
	case CodeBin32:
		alen := uint32(0)
		err := binary.Read(coder.r,binary.BigEndian,&alen)
		if err!= nil{
			return nil,err
		}
		btlen = int(alen)
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
		u16 := uint16(0)
		if err = binary.Read(coder.r,binary.BigEndian,&u16);err!=nil{
			return nil,err
		}
		btlen = int(u16)+1
	case CodeExt32:
		u32 := uint32(0)
		if err = binary.Read(coder.r,binary.BigEndian,&u32);err!=nil{
			return nil,err
		}
		btlen = int(u32)+1
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

func (coder *MsgPackDecoder)DecodeStrMap(code MsgPackCode,rec *DxValue.DxRecord)error  {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.readCode();err!=nil{
			return err
		}
	}
	maplen := 0
	switch code {
	case CodeMap16:
		recLen := uint16(0)
		if err = binary.Read(coder.r,binary.BigEndian,&recLen);err!=nil{
			return err
		}
		maplen = int(recLen)
	case CodeMap32:
		recLen := uint32(0)
		if err = binary.Read(coder.r,binary.BigEndian,&recLen);err!= nil{
			return err
		}
		maplen = int(recLen)
	default:
		if code >= CodeFixedMapLow && code<= CodeFixedMapHigh{
			maplen = int(code & 0xf)
		}
	}
	//先判断第一个元素是否为字符串
	var keybt []byte
	for i := 0;i<maplen;i++{
		if keybt,err = coder.DecodeString(CodeUnkonw);err!=nil{
			return err
		}

	}
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
		alen := uint16(0)
		if err = binary.Read(coder.r,binary.BigEndian,&alen);err!=nil{
			return nil,err
		}
		strlen = int(alen)
	case CodeStr32:
		alen := uint32(0)
		if err = binary.Read(coder.r,binary.BigEndian,&alen);err!= nil{
			return nil,err
		}
		strlen = int(alen)
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

func (coder *MsgPackDecoder)Decode(r io.Reader, v *DxValue.DxBaseValue)(error)  {
	if bf,ok := r.(*bufio.Reader);ok{
		coder.r = bf
	}else{
		coder.r = bufio.NewReader(r)
	}
	coder.r.Reset(r)
	switch v.ValueType() {
	case DxValue.DVT_Array:
	case DxValue.DVT_DateTime:
		if dt,err := coder.DecodeDateTime(CodeUnkonw);err !=nil{
			return err
		}else{
			v.SetDateTime(dt)
		}
	case DxValue.DVT_Int,DxValue.DVT_Int32,DxValue.DVT_Int64:
		if i64,err := coder.DecodeInt(CodeUnkonw);err!=nil{
			return err
		}else {
			v.SetInt64(i64)
		}
	case DxValue.DVT_Float,DxValue.DVT_Double:
		if vf,err := coder.DecodeFloat(CodeUnkonw);err!=nil{
			return err
		}else{
			v.SetDouble(vf)
		}
	case DxValue.DVT_Ext:
		if vb,err := coder.DecodeExtValue(CodeUnkonw);err!=nil{
			return err
		}else{
			v.SetExtValue(vb)
		}
	case DxValue.DVT_Bool:
		if code,err := coder.readCode();err!=nil{
			return err
		}else if code == CodeFalse{
			v.SetBool(false)
		}else if code == CodeTrue{
			v.SetBool(true)
		}
	case DxValue.DVT_Binary:
		if bt,err := coder.DecodeBinary(CodeUnkonw);err!=nil{
			return err
		}else{
			v.SetBinary(bt)
		}
	case DxValue.DVT_Record:
		rec,_ := v.AsRecord()
		return coder.DecodeStrMap(CodeUnkonw,rec)
	case DxValue.DVT_RecordIntKey:

	}
	coder.r.UnreadByte()
	return DxValue.ErrValueType
}