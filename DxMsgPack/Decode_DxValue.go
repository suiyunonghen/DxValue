package DxMsgPack

import (
	"bufio"
	"github.com/suiyunonghen/DxValue"
	"io"
	"github.com/suiyunonghen/DxCommonLib"
	"encoding/binary"
	"time"
	"unsafe"
	"errors"
)

var(
	ErrInvalidateMap	= errors.New("Invalidate Map Format")
	ErrInvalidateMapKey	= errors.New("Invalidate Map Key,Key Can Only Int or String")
	ErrInvalidateArrLen = errors.New("Is not a Array Len Flag")
)

type bufReader interface{
	io.ByteScanner
	io.Reader
}

type MsgPackDecoder   struct{
	r  bufReader
}


func (coder *MsgPackDecoder)DecodeStrMapKvRecord(strMap *DxValue.DxRecord,strcode MsgPackCode)(error)  {
	keybt,err := coder.DecodeString(strcode)
	if err != nil{
		return err
	}
	if strcode,err = coder.readCode();err!=nil{
		return err
	}
	keyName := DxCommonLib.FastByte2String(keybt)
	if strcode.IsStr(){
		if keybt,err = coder.DecodeString(strcode);err!=nil{
			return err
		}
		strMap.SetString(keyName,DxCommonLib.FastByte2String(keybt))
	}else if strcode.IsFixedNum(){
		strMap.SetInt32(keyName,int32(int8(strcode)))
	}else if strcode.IsInt(){
		if i64,err := coder.DecodeInt(strcode);err!=nil{
			return err
		}else{
			strMap.SetInt64(keyName,i64)
		}
	}else if strcode.IsMap(){
		if baseV,err := coder.DecodeUnknownMap(strcode);err!=nil{
			return err
		}else{
			strMap.SetBaseValue(keyName,baseV)
		}
	}else if strcode.IsArray(){
		if arr,err := coder.DecodeArray(strcode);err!=nil{
			return err
		}else{
			strMap.SetArray(keyName,arr)
		}
	}else if strcode.IsBin(){
		if bin,err := coder.DecodeBinary(strcode);err!=nil{
			return err
		} else{
			strMap.SetBinary(keyName,bin,true)
		}
	}else if strcode.IsExt(){
		if bin ,err := coder.DecodeExtValue(strcode);err!=nil{
			return err
		}else{
			strMap.SetExtValue(keyName,bin)
		}
	}else{
		switch strcode {
		case CodeTrue:	strMap.SetBool(keyName,true)
		case CodeFalse: strMap.SetBool(keyName,false)
		case CodeNil:   strMap.SetNull(keyName)
		case CodeFloat:
			var b [4]byte
			if _,err := coder.r.Read(b[:]);err!=nil{
				return err
			}
			v32 := binary.BigEndian.Uint32(b[:])
			strMap.SetFloat(keyName,*(*float32)(unsafe.Pointer(&v32)))
		case CodeDouble:
			var b [8]byte
			if _,err := coder.r.Read(b[:]);err!=nil{
				return err
			}
			v64 := binary.BigEndian.Uint64(b[:])
			strMap.SetDouble(keyName,*(*float64)(unsafe.Pointer(&v64)))
		case CodeFixExt4:
			if strcode,err = coder.readCode();err!=nil{
				return err
			}
			if int8(strcode) == -1{
				var b [4]byte
				if _,err := coder.r.Read(b[:]);err!=nil{
					return err
				}
				ms := binary.BigEndian.Uint32(b[:])
				ntime := time.Now()
				ns := ntime.Unix()
				ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
				strMap.SetDateTime(keyName, DxCommonLib.Time2DelphiTime(&ntime))
			}else{
				var mb [5]byte
				if _,err = coder.r.Read(mb[1:]);err!=nil{
					return err
				}
				mb[0] = byte(strcode)
				strMap.SetExtValue(keyName,mb[:])
			}
		}
	}
	return nil
}


func (coder *MsgPackDecoder)DecodeArrayElement(arr *DxValue.DxArray,eleIndex int)(error)  {
	code,err := coder.readCode()
	if err!=nil{
		return err
	}
	if code.IsStr(){
		if stbt,err := coder.DecodeString(code);err!=nil{
			return err
		}else{
			arr.SetString(eleIndex,DxCommonLib.FastByte2String(stbt))
		}
	}else if code.IsInt(){
		if i64,err := coder.DecodeInt(code);err!=nil{
			return err
		}else{
			arr.SetInt64(eleIndex,i64)
		}
	}else if code.IsFixedNum(){
		arr.SetInt32(eleIndex,int32(int8(code)))
	}else if code.IsMap(){
		if mpbv,err := coder.DecodeUnknownMap(code);err!=nil{
			return err
		}else{
			arr.SetBaseValue(eleIndex,mpbv)
		}
	}else if code.IsArray(){
		if carr,err := coder.DecodeArray(code);err!=nil{
			return err
		}else{
			arr.SetArray(eleIndex,carr)
		}
	}else if code.IsBin(){
		if bin,err := coder.DecodeBinary(code);err!=nil{
			return err
		}else{
			arr.SetBinary(eleIndex,bin)
		}
	}else if code.IsExt(){
		if bin,err := coder.DecodeExtValue(code);err!=nil{
			return err
		}else {
			arr.SetExtValue(eleIndex,bin)
		}
	} else{
		switch code {
		case CodeTrue:	arr.SetBool(eleIndex,true)
		case CodeFalse: arr.SetBool(eleIndex,false)
		case CodeNil:	arr.SetNull(eleIndex)
		case CodeFloat:
			var b [4]byte
			if _,err := coder.r.Read(b[:]);err!=nil{
				return err
			}
			v32 := binary.BigEndian.Uint32(b[:])
			arr.SetFloat(eleIndex,*(*float32)(unsafe.Pointer(&v32)))
		case CodeDouble:
			var b [8]byte
			if _,err := coder.r.Read(b[:]);err!=nil{
				return err
			}
			v64 := binary.BigEndian.Uint64(b[:])
			arr.SetDouble(eleIndex,*(*float64)(unsafe.Pointer(&v64)))
		case CodeFixExt4:
			if code,err = coder.readCode();err!=nil{
				return err
			}
			if int8(code) == -1{
				var b [4]byte
				if _,err := coder.r.Read(b[:]);err!=nil{
					return err
				}
				ms := binary.BigEndian.Uint32(b[:])
				ntime := time.Now()
				ns := ntime.Unix()
				ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
				arr.SetDateTime(eleIndex, DxCommonLib.Time2DelphiTime(&ntime))
			}else{
				var mb [5]byte
				if _,err = coder.r.Read(mb[1:]);err!=nil{
					return err
				}
				mb[0] = byte(code)
				arr.SetExtValue(eleIndex,mb[:])
			}
		}
	}
	return nil
}

func (coder *MsgPackDecoder)DecodeArray(code MsgPackCode)(*DxValue.DxArray,error)  {
	var (
		err error
		arrlen int
	)
	if code == CodeUnkonw{
		if code,err = coder.readCode();err!=nil{
			return nil,err
		}
	}
	if arrlen,err = coder.DecodeArrayLen(code);err!=nil{
		return nil,err
	}
	arr := DxValue.NewArray()
	for i := 0;i<arrlen;i++{
		if err = coder.DecodeArrayElement(arr,i);err!=nil{
			return nil,err
		}
	}
	return arr,nil
}

func (coder *MsgPackDecoder)DecodeUnknownMap(code MsgPackCode)(*DxValue.DxBaseValue,error)  {
	if maplen,err := coder.DecodeMapLen(code);err!=nil{
		return nil, err
	}else{
		//判断键值，是Int还是str
		var baseV *DxValue.DxBaseValue
		if code,err = coder.readCode();err!=nil{
			return nil,err
		}
		if code.IsInt(){
			iMap := DxValue.NewIntKeyRecord()
			baseV = &iMap.DxBaseValue
		}else if code.IsStr(){
			iMap := DxValue.NewRecord()
			baseV = &iMap.DxBaseValue
			if err = coder.DecodeStrMapKvRecord(iMap,code);err!=nil{
				return nil,err
			}
			for j := 1;j<maplen;j++{
				if err = coder.DecodeStrMapKvRecord(iMap,CodeUnkonw);err!=nil{
					return nil,err
				}
			}
		}else{
			return nil,ErrInvalidateMapKey
		}
		return baseV,nil
	}
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
		var b [2]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return err
		}
		maplen =  int(binary.BigEndian.Uint16(b[:]))
	case CodeMap32:
		var b [4]byte
		if _,err := coder.r.Read(b[:]);err!=nil{
			return err
		}
		maplen =  int(binary.BigEndian.Uint32(b[:]))
	default:
		if code >= CodeFixedMapLow && code<= CodeFixedMapHigh{
			maplen = int(code & 0xf)
		}
	}
	for i := 0;i<maplen;i++{
		if err = coder.DecodeStrMapKvRecord(rec,CodeUnkonw);err!=nil{
			return err
		}
	}
	return nil
}

func (coder *MsgPackDecoder)Decode(r io.Reader, v *DxValue.DxBaseValue)(error)  {
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