package DxValue

import (
	"github.com/suiyunonghen/DxValue/Coders/DxMsgPack"
	"github.com/suiyunonghen/DxCommonLib"
	"unsafe"
	"time"
	"io"
	"github.com/suiyunonghen/DxValue/Coders"
	"encoding/binary"
	"bytes"
	"errors"
)

type  DxMsgPackDecoder  struct{
	DxMsgPack.MsgPackDecoder
}


func (coder *DxMsgPackDecoder)DecodeStrMapKvRecord(strMap *DxRecord,strcode DxMsgPack.MsgPackCode)(error)  {
	keybt,err := coder.DecodeString(strcode)
	if err != nil{
		return err
	}
	if strcode,err = coder.ReadCode();err!=nil{
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
			strMap.SetBinary(keyName,bin,true,BET_Base64)
		}
	}else if strcode.IsExt(){
		if bin ,err := coder.DecodeExtValue(strcode);err!=nil{
			return err
		}else{
			strMap.SetExtValue(keyName,bin)
		}
	}else{
		switch strcode {
		case DxMsgPack.CodeTrue:	strMap.SetBool(keyName,true)
		case DxMsgPack.CodeFalse: strMap.SetBool(keyName,false)
		case DxMsgPack.CodeNil:   strMap.SetNull(keyName)
		case DxMsgPack.CodeFloat:
			if v32,err := coder.ReadBigEnd32();err!=nil{
				return err
			}else{
				strMap.SetFloat(keyName,*(*float32)(unsafe.Pointer(&v32)))
			}
		case DxMsgPack.CodeDouble:
			if v64,err := coder.ReadBigEnd64();err!=nil{
				return err
			}else{
				strMap.SetDouble(keyName,*(*float64)(unsafe.Pointer(&v64)))
			}

		case DxMsgPack.CodeFixExt4:
			if strcode,err = coder.ReadCode();err!=nil{
				return err
			}
			if int8(strcode) == -1{
				if ms,err := coder.ReadBigEnd32();err!=nil{
					return err
				}else{
					ntime := time.Now()
					ns := ntime.Unix()
					ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
					strMap.SetDateTime(keyName, DxCommonLib.Time2DelphiTime(&ntime))
				}

			}else{
				var mb [5]byte
				if err = coder.Read(mb[1:]);err!=nil{
					return err
				}
				mb[0] = byte(strcode)
				strMap.SetExtValue(keyName,mb[:])
			}
		}
	}
	return nil
}

func (coder *DxMsgPackDecoder)DecodeUnkown()(*DxBaseValue,error)  {
	code,err := coder.ReadCode()
	if err!=nil{
		return nil,err
	}
	if code.IsStr(){
		if stbt,err := coder.DecodeString(code);err!=nil{
			return nil,err
		}else{
			v := DxStringValue{}
			v.fvalue = DxCommonLib.FastByte2String(stbt)
			v.fValueType = DVT_String
			return &v.DxBaseValue,nil
		}
	}else if code.IsFixedNum(){
		var v DxInt32Value
		v.fvalue =  int32(int8(code))
		v.fValueType = DVT_Int32
		return &v.DxBaseValue,nil
	}else if code.IsInt(){
		if i64,err := coder.DecodeInt(code);err!=nil{
			return nil,err
		}else{
			var v DxIntValue
			v.fvalue =  int(i64)
			v.fValueType = DVT_Int
			return &v.DxBaseValue,nil
		}
	}else if code.IsMap(){
		return coder.DecodeUnknownMap(code)
	}else if code.IsArray(){
		arr := NewArray()
		if err := coder.Decode2Array(code,arr);err!=nil{
			return nil,err
		}
		return &arr.DxBaseValue,nil
	}else if code.IsBin(){
		if bin,err := coder.DecodeBinary(code);err!=nil{
			return nil,err
		}else{
			var v DxBinaryValue
			v.fValueType = DVT_Binary
			v.fbinary = bin
			return  &v.DxBaseValue,nil
		}
	}else if code.IsExt(){
		if bin,err := coder.DecodeExtValue(code);err!=nil{
			return nil,err
		}else {
			var v DxExtValue
			v.fValueType = DVT_Ext
			v.fExtType = bin[0]
			v.fdata = bin
		}
	} else{
		switch code {
		case DxMsgPack.CodeTrue:
			var v DxBoolValue
			v.fvalue = true
			v.fValueType = DVT_Bool
			return &v.DxBaseValue,nil
		case DxMsgPack.CodeFalse:
			var v DxBoolValue
			v.fvalue = false
			v.fValueType = DVT_Bool
			return &v.DxBaseValue,nil
		case DxMsgPack.CodeNil:	return nil,nil
		case DxMsgPack.CodeFloat:
			if v32,err := coder.ReadBigEnd32();err!=nil{
				return nil,err
			}else{
				var v DxFloatValue
				v.fvalue = *(*float32)(unsafe.Pointer(&v32))
				v.fValueType = DVT_Float
				return  &v.DxBaseValue,nil
			}
		case DxMsgPack.CodeDouble:
			if v64,err := coder.ReadBigEnd64();err!=nil{
				return nil,err
			}else{
				var v DxDoubleValue
				v.fvalue = *(*float64)(unsafe.Pointer(&v64))
				v.fValueType = DVT_Double
				return  &v.DxBaseValue,nil
			}

		case DxMsgPack.CodeFixExt4:
			if code,err = coder.ReadCode();err!=nil{
				return nil,err
			}
			if int8(code) == -1{
				if ms,err := coder.ReadBigEnd32();err!=nil{
					return nil,err
				}else{
					ntime := time.Now()
					ns := ntime.Unix()
					ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
					var v DxDoubleValue
					v.fvalue = float64(DxCommonLib.Time2DelphiTime(&ntime))
					v.fValueType = DVT_DateTime
					return  &v.DxBaseValue,nil
				}
			}else{
				var mb [5]byte
				if err = coder.Read(mb[1:]);err!=nil{
					return nil,err
				}
				mb[0] = byte(code)
				var v DxExtValue
				v.fValueType = DVT_Ext
				v.fExtType = mb[0]
				v.fdata = mb[:]
				return  &v.DxBaseValue,nil
			}
		}
	}
	return nil,nil
}

func (coder *DxMsgPackDecoder)DecodeIntKeyMapKvRecord(intKeyMap *DxIntKeyRecord,keycode DxMsgPack.MsgPackCode)(error)  {
	intKey,err := coder.DecodeInt(keycode)
	if err != nil{
		return err
	}
	if keycode,err = coder.ReadCode();err!=nil{
		return err
	}
	if keycode.IsStr(){
		if strbt,err := coder.DecodeString(keycode);err!=nil{
			return err
		}else{
			intKeyMap.SetString(intKey,DxCommonLib.FastByte2String(strbt))
		}
	}else if keycode.IsFixedNum(){
		intKeyMap.SetInt32(intKey,int32(int8(keycode)))
	}else if keycode.IsInt(){
		if i64,err := coder.DecodeInt(keycode);err!=nil{
			return err
		}else{
			intKeyMap.SetInt64(intKey,i64)
		}
	}else if keycode.IsMap(){
		if baseV,err := coder.DecodeUnknownMap(keycode);err!=nil{
			return err
		}else{
			intKeyMap.SetBaseValue(intKey,baseV)
		}
	}else if keycode.IsArray(){
		if arr,err := coder.DecodeArray(keycode);err!=nil{
			return err
		}else{
			intKeyMap.SetArray(intKey,arr)
		}
	}else if keycode.IsBin(){
		if bin,err := coder.DecodeBinary(keycode);err!=nil{
			return err
		} else{
			intKeyMap.SetBinary(intKey,bin,true)
		}
	}else if keycode.IsExt(){
		if bin ,err := coder.DecodeExtValue(keycode);err!=nil{
			return err
		}else{
			intKeyMap.SetExtValue(intKey,bin)
		}
	}else{
		switch keycode {
		case DxMsgPack.CodeTrue:	intKeyMap.SetBool(intKey,true)
		case DxMsgPack.CodeFalse: intKeyMap.SetBool(intKey,false)
		case DxMsgPack.CodeNil:   intKeyMap.SetNull(intKey)
		case DxMsgPack.CodeFloat:
			if v32,err := coder.ReadBigEnd32();err!=nil{
				return err
			}else{
				intKeyMap.SetFloat(intKey,*(*float32)(unsafe.Pointer(&v32)))
			}
		case DxMsgPack.CodeDouble:
			if v64,err := coder.ReadBigEnd64();err!=nil{
				return err
			}else{
				intKeyMap.SetDouble(intKey,*(*float64)(unsafe.Pointer(&v64)))
			}

		case DxMsgPack.CodeFixExt4:
			if keycode,err = coder.ReadCode();err!=nil{
				return err
			}
			if int8(keycode) == -1{
				if ms,err := coder.ReadBigEnd32();err!=nil{
					return err
				}else{
					ntime := time.Now()
					ns := ntime.Unix()
					ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
					intKeyMap.SetDateTime(intKey, DxCommonLib.Time2DelphiTime(&ntime))
				}

			}else{
				var mb [5]byte
				if err = coder.Read(mb[1:]);err!=nil{
					return err
				}
				mb[0] = byte(keycode)
				intKeyMap.SetExtValue(intKey,mb[:])
			}
		}
	}
	return nil
}


func (coder *DxMsgPackDecoder)DecodeArrayElement(arr *DxArray,eleIndex int)(error)  {
	code,err := coder.ReadCode()
	if err!=nil{
		return err
	}
	if code.IsStr(){
		if stbt,err := coder.DecodeString(code);err!=nil{
			return err
		}else{
			arr.SetString(eleIndex,DxCommonLib.FastByte2String(stbt))
		}
	}else if code.IsFixedNum(){
		arr.SetInt32(eleIndex,int32(int8(code)))
	}else if code.IsInt(){
		if i64,err := coder.DecodeInt(code);err!=nil{
			return err
		}else{
			arr.SetInt64(eleIndex,i64)
		}
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
		case DxMsgPack.CodeTrue:	arr.SetBool(eleIndex,true)
		case DxMsgPack.CodeFalse: arr.SetBool(eleIndex,false)
		case DxMsgPack.CodeNil:	arr.SetNull(eleIndex)
		case DxMsgPack.CodeFloat:
			if v32,err := coder.ReadBigEnd32();err!=nil{
				return err
			}else{
				arr.SetFloat(eleIndex,*(*float32)(unsafe.Pointer(&v32)))
			}
		case DxMsgPack.CodeDouble:
			if v64,err := coder.ReadBigEnd64();err!=nil{
				return err
			}else{
				arr.SetDouble(eleIndex,*(*float64)(unsafe.Pointer(&v64)))
			}

		case DxMsgPack.CodeFixExt4:
			if code,err = coder.ReadCode();err!=nil{
				return err
			}
			if int8(code) == -1{
				if ms,err := coder.ReadBigEnd32();err!=nil{
					return err
				}else{
					ntime := time.Now()
					ns := ntime.Unix()
					ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
					arr.SetDateTime(eleIndex, DxCommonLib.Time2DelphiTime(&ntime))
				}
			}else{
				var mb [5]byte
				if err = coder.Read(mb[1:]);err!=nil{
					return err
				}
				mb[0] = byte(code)
				arr.SetExtValue(eleIndex,mb[:])
			}
		}
	}
	return nil
}

func (coder *DxMsgPackDecoder)DecodeArray(code DxMsgPack.MsgPackCode)(*DxArray,error)  {
	var (
		err error
		arrlen int
	)
	if code == DxMsgPack.CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
			return nil,err
		}
	}
	if arrlen,err = coder.DecodeArrayLen(code);err!=nil{
		return nil,err
	}
	arr := NewArray()
	for i := 0;i<arrlen;i++{
		if err = coder.DecodeArrayElement(arr,i);err!=nil{
			return nil,err
		}
	}
	return arr,nil
}

func (coder *DxMsgPackDecoder)Decode2Array(code DxMsgPack.MsgPackCode,arr *DxArray)(error)  {
	var (
		err error
		arrlen int
	)
	if code == DxMsgPack.CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
			return err
		}
	}
	if arrlen,err = coder.DecodeArrayLen(code);err!=nil{
		return err
	}
	for i := 0;i<arrlen;i++{
		if err = coder.DecodeArrayElement(arr,i);err!=nil{
			return err
		}
	}
	return nil
}

func (coder *DxMsgPackDecoder)DecodeUnknownMap(code DxMsgPack.MsgPackCode)(*DxBaseValue,error)  {
	if maplen,err := coder.DecodeMapLen(code);err!=nil{
		return nil, err
	}else if maplen > 0{
		//判断键值，是Int还是str
		var baseV *DxBaseValue
		if code,err = coder.ReadCode();err!=nil{
			return nil,err
		}
		if code.IsInt(){
			iMap := NewIntKeyRecord()
			baseV = &iMap.DxBaseValue
			if err = coder.DecodeIntKeyMapKvRecord(iMap,code);err!=nil{
				return nil,err
			}
			for j := 1;j<maplen;j++{
				if err = coder.DecodeIntKeyMapKvRecord(iMap,DxMsgPack.CodeUnkonw);err!=nil{
					return nil,err
				}
			}
			return baseV,nil
		}else if code.IsStr(){
			iMap := NewRecord()
			baseV = &iMap.DxBaseValue
			if err = coder.DecodeStrMapKvRecord(iMap,code);err!=nil{
				return nil,err
			}
			for j := 1;j<maplen;j++{
				if err = coder.DecodeStrMapKvRecord(iMap,DxMsgPack.CodeUnkonw);err!=nil{
					return nil,err
				}
			}
			return baseV,nil
		}
		return nil,DxMsgPack.ErrInvalidateMapKey
	}else{
		return nil, nil
	}
}

func (coder *DxMsgPackDecoder)DecodeStrMap(code DxMsgPack.MsgPackCode,rec *DxRecord)error  {
	var err error
	if code == DxMsgPack.CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
			return err
		}
	}
	maplen := 0
	switch code {
	case DxMsgPack.CodeMap16:
		if v16,err := coder.ReadBigEnd16();err!=nil{
			return err
		}else {
			maplen = int(v16)
		}
	case DxMsgPack.CodeMap32:
		if v32,err := coder.ReadBigEnd32();err!=nil{
			return err
		}else {
			maplen = int(v32)
		}
	default:
		if code >= DxMsgPack.CodeFixedMapLow && code<= DxMsgPack.CodeFixedMapHigh{
			maplen = int(code & 0xf)
		}
	}
	for i := 0;i<maplen;i++{
		if err = coder.DecodeStrMapKvRecord(rec,DxMsgPack.CodeUnkonw);err!=nil{
			return err
		}
	}
	return nil
}


func (coder *DxMsgPackDecoder)DecodeIntKeyMap(code DxMsgPack.MsgPackCode,rec *DxIntKeyRecord)error  {
	var err error
	if code == DxMsgPack.CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
			return err
		}
	}
	maplen := 0
	switch code {
	case DxMsgPack.CodeMap16:
		if v16,err := coder.ReadBigEnd16();err!=nil{
			return err
		}else {
			maplen = int(v16)
		}
	case DxMsgPack.CodeMap32:
		if v32,err := coder.ReadBigEnd32();err!=nil{
			return err
		}else {
			maplen = int(v32)
		}
	default:
		if code >= DxMsgPack.CodeFixedMapLow && code<= DxMsgPack.CodeFixedMapHigh{
			maplen = int(code & 0xf)
		}
	}
	for i := 0;i<maplen;i++{
		if err = coder.DecodeIntKeyMapKvRecord(rec,DxMsgPack.CodeUnkonw);err!=nil{
			return err
		}
	}
	return nil
}

func (coder *DxMsgPackDecoder)Decode(v *DxBaseValue)(error)  {
	switch v.ValueType() {
	case DVT_Array:
		arrlen,err  := coder.DecodeArrayLen(DxMsgPack.CodeUnkonw)
		if err !=nil{
			return err
		}
		arr := (*DxArray)(unsafe.Pointer(v))
		arr.TruncateArray(arrlen)
		for i := 0;i<arrlen;i++{
			if err = coder.DecodeArrayElement(arr,i);err!=nil{
				return err
			}

		}
	case DVT_DateTime:
		if dt,err := coder.DecodeDateTime(DxMsgPack.CodeUnkonw);err !=nil{
			return err
		}else{
			v.SetDateTime(dt)
		}
	case DVT_Int,DVT_Int32,DVT_Int64:
		if i64,err := coder.DecodeInt(DxMsgPack.CodeUnkonw);err!=nil{
			return err
		}else {
			v.SetInt64(i64)
		}
	case DVT_Float,DVT_Double:
		if vf,err := coder.DecodeFloat(DxMsgPack.CodeUnkonw);err!=nil{
			return err
		}else{
			v.SetDouble(vf)
		}
	case DVT_Ext:
		if vb,err := coder.DecodeExtValue(DxMsgPack.CodeUnkonw);err!=nil{
			return err
		}else{
			v.SetExtValue(vb)
		}
	case DVT_Bool:
		if code,err := coder.ReadCode();err!=nil{
			return err
		}else if code == DxMsgPack.CodeFalse{
			v.SetBool(false)
		}else if code == DxMsgPack.CodeTrue{
			v.SetBool(true)
		}
	case DVT_Binary:
		if bt,err := coder.DecodeBinary(DxMsgPack.CodeUnkonw);err!=nil{
			return err
		}else{
			v.SetBinary(bt)
		}
	case DVT_Record:
		rec,_ := v.AsRecord()
		return coder.DecodeStrMap(DxMsgPack.CodeUnkonw,rec)
	case DVT_RecordIntKey:
		rec,_ := v.AsIntRecord()
		return coder.DecodeIntKeyMap(DxMsgPack.CodeUnkonw,rec)
	}
	coder.UnreadByte()
	return Coders.ErrValueType
}

func (dcoder *DxMsgPackDecoder)DecodeStand(v interface{})(error)  {
	switch value := v.(type) {
	case *DxRecord:
		return dcoder.DecodeStrMap(DxMsgPack.CodeUnkonw,value)
	case DxRecord:
		return errors.New("must pointer value")
	case *DxArray:
		arrlen,err  := dcoder.DecodeArrayLen(DxMsgPack.CodeUnkonw)
		if err !=nil{
			return err
		}
		value.TruncateArray(arrlen)
		for i := 0;i<arrlen;i++{
			if err = dcoder.DecodeArrayElement(value,i);err!=nil{
				return err
			}
		}
		return nil
	case DxArray:
		return errors.New("must pointer value")
	case *DxIntKeyRecord:
		return dcoder.DecodeIntKeyMap(DxMsgPack.CodeUnkonw,value)
	case DxIntKeyRecord:
		return errors.New("must pointer value")
	default:
		return dcoder.MsgPackDecoder.DecodeStand(v)
	}
}

func NewDecoder(r io.Reader)*DxMsgPackDecoder  {
	var result DxMsgPackDecoder
	result.ReSetReader(r)
	return &result
}



type  DxMsgPackEncoder struct{
	DxMsgPack.MsgPackEncoder
}


func (encoder *DxMsgPackEncoder)EncodeExtValue(v *DxExtValue)(err error)  {
	btlen := uint(0)
	bt := v.ExtData()
	btlen = uint(len(bt))
	buf := encoder.Buffer()
	buf[1] = v.ExtType()
	switch {
	case btlen == 1:
		buf[0] = byte(DxMsgPack.CodeFixExt1)
		err = encoder.Write(buf[:1])
	case btlen == 2:
		buf[0] = byte(DxMsgPack.CodeFixExt2)
		err = encoder.Write(buf[:1])
	case btlen == 4:
		buf[0] = byte(DxMsgPack.CodeFixExt4)
		err = encoder.Write(buf[:1])
	case btlen == 8:
		buf[0] = byte(DxMsgPack.CodeFixExt8)
		err = encoder.Write(buf[:1])
	case btlen <= 16:
		buf[0] = byte(DxMsgPack.CodeFixExt16)
		err = encoder.Write(buf[:1])
	case btlen <= DxMsgPack.Max_str8_len:
		buf[0] = byte(DxMsgPack.CodeExt8)
		buf[1] = byte(btlen)
		buf[2] = v.ExtType()
		err = encoder.Write(buf[:3])
	case btlen <= DxMsgPack.Max_str16_len:
		buf[0] = byte(DxMsgPack.CodeExt16)
		binary.BigEndian.PutUint16(buf[1:3],uint16(btlen))
		buf[3] = v.ExtType()
		err = encoder.Write(buf[:4])
	default:
		if btlen > DxMsgPack.Max_str32_len{
			btlen = DxMsgPack.Max_str32_len
		}
		buf[0] = 0xc6
		binary.BigEndian.PutUint32(buf[1:5],uint32(btlen))
		buf[5] = v.ExtType()
		err = encoder.Write(buf[:6])
	}
	if err == nil && btlen > 0{
		err = encoder.Write(bt[:btlen])
	}
	return
}


func (encoder *DxMsgPackEncoder)EncodeRecord(r *DxRecord)(err error)  {
	maplen := uint(r.Length())
	if maplen <= DxMsgPack.Max_fixmap_len{   //fixmap
		err = encoder.WriteByte(0x80 | byte(maplen))
	}else if maplen <= DxMsgPack.Max_map16_len{
		//写入长度
		err = encoder.WriteUint16(uint16(maplen),DxMsgPack.CodeMap16)
	}else{
		if maplen > DxMsgPack.Max_map32_len{
			maplen = DxMsgPack.Max_map32_len
		}
		err = encoder.WriteUint32(uint32(maplen),DxMsgPack.CodeMap32)
	}
	if err != nil{
		return
	}
	//写入对象信息,Kv对

	for k,v := range r.fRecords{
		if err = encoder.EncodeString(k);err!=nil{
			return err
		}
		if v != nil{
			err = encoder.Encode(v)
		}else{
			err = encoder.WriteByte(0xc0) //null
		}
		if err!=nil{
			return err
		}
	}

	return nil
}


func (encoder *DxMsgPackEncoder)EncodeRecordIntKey(r *DxIntKeyRecord)(err error)  {
	maplen := uint(r.Length())
	if maplen <= DxMsgPack.Max_fixmap_len{   //fixmap
		err = encoder.WriteByte(0x80 | byte(maplen))
	}else if maplen <= DxMsgPack.Max_map16_len{
		//写入长度
		err = encoder.WriteUint16(uint16(maplen),DxMsgPack.CodeMap16)
	}else{
		if maplen > DxMsgPack.Max_map32_len{
			maplen = DxMsgPack.Max_map32_len
		}
		err = encoder.WriteUint32(uint32(maplen),DxMsgPack.CodeMap32)
	}
	if err != nil{
		return
	}
	//写入对象信息,Kv对
	for k,v := range r.fRecords{
		if err = encoder.EncodeInt(k);err!=nil{
			return err
		}
		if v != nil{
			err = encoder.Encode(v)
		}else{
			err = encoder.WriteByte(0xc0) //null
		}
		if err!=nil{
			return err
		}
	}

	return nil
}

func (encoder *DxMsgPackEncoder)Encode(v *DxBaseValue)(err error)  {
	switch v.ValueType() {
	case DVT_Record:
		err = encoder.EncodeRecord((*DxRecord)(unsafe.Pointer(v)))
	case DVT_RecordIntKey:
		err = encoder.EncodeRecordIntKey((*DxIntKeyRecord)(unsafe.Pointer(v)))
	case DVT_Int:
		return encoder.EncodeInt(int64((*DxIntValue)(unsafe.Pointer(v)).fvalue))
	case DVT_Int32:
		return encoder.EncodeInt(int64((*DxInt32Value)(unsafe.Pointer(v)).fvalue))
	case DVT_Int64:
		return encoder.EncodeInt((*DxInt64Value)(unsafe.Pointer(v)).fvalue)
	case DVT_Bool:
		return encoder.EncodeBool((*DxBoolValue)(unsafe.Pointer(v)).fvalue)
	case DVT_String:
		return encoder.EncodeString((*DxStringValue)(unsafe.Pointer(v)).fvalue)
	case DVT_Float:
		return encoder.EncodeFloat((*DxFloatValue)(unsafe.Pointer(v)).fvalue)
	case DVT_Double:
		return encoder.EncodeDouble((*DxDoubleValue)(unsafe.Pointer(v)).fvalue)
	case DVT_Binary:
		bt := (*DxBinaryValue)(unsafe.Pointer(v)).fbinary
		if bt != nil{
			return encoder.EncodeBinary(bt)
		}else{
			return encoder.WriteByte(0xc0)
		}
	case DVT_Ext:
		return encoder.EncodeExtValue((*DxExtValue)(unsafe.Pointer(v)))
	case DVT_Array:
		return encoder.EncodeArray((*DxArray)(unsafe.Pointer(v)))
	case DVT_DateTime:
		return encoder.EncodeDateTime(DxCommonLib.TDateTime((*DxDoubleValue)(unsafe.Pointer(v)).fvalue))
	default:
		return encoder.WriteByte(0xc0) //null
	}
	return nil
}


func (encoder *DxMsgPackEncoder)EncodeArray(arr *DxArray)(err error)  {
	arlen := uint(arr.Length())
	switch {
	case arlen < 16: //1001XXXX|    N objects
		err = encoder.WriteByte(byte(DxMsgPack.CodeFixedArrayLow) | byte(arlen))
	case arlen <= DxMsgPack.Max_map16_len:  //0xdc  |YYYYYYYY|YYYYYYYY|    N objects
		encoder.WriteUint16(uint16(arlen),DxMsgPack.CodeArray16)
	default:
		if arlen > DxMsgPack.Max_map32_len{
			arlen = DxMsgPack.Max_map32_len
		}
		encoder.WriteUint32(uint32(arlen),DxMsgPack.CodeArray32)
	}

	for i := 0;i <= int(arlen - 1);i++{
		vbase := arr.AsBaseValue(i)
		if vbase == nil{
			err = encoder.WriteByte(0xc0) //null
		}else{
			err = encoder.Encode(vbase)
		}
		if err != nil{
			return
		}
	}
	return err
}

func (encoder *DxMsgPackEncoder)EncodeStand(v interface{})(error)  {
	switch value := v.(type) {
	case *DxRecord:
		return encoder.EncodeRecord(value)
	case DxRecord:
		return encoder.EncodeRecord(&value)
	case *DxArray:
		return encoder.EncodeArray(value)
	case DxArray:
		return encoder.EncodeArray(&value)
	case *DxIntKeyRecord:
		return encoder.EncodeRecordIntKey(value)
	case DxIntKeyRecord:
		return encoder.EncodeRecordIntKey(&value)
	default:
		return encoder.MsgPackEncoder.EncodeStand(v)
	}
}


func NewEncoder(w io.Writer) *DxMsgPackEncoder {
	var result DxMsgPackEncoder
	result.Buffer()
	result.ReSet(w)
	return &result
}


func Marshal(v...interface{})([]byte,error) {
	var buf bytes.Buffer
	coder := NewEncoder(&buf)
	for _,value := range v{

		if err := coder.EncodeStand(value);err!=nil{
			return nil,err
		}
	}
	return buf.Bytes(),nil
}

func Unmarshal(data []byte, v...interface{}) error {
	coder := NewDecoder(bytes.NewReader(data))
	for _,vdst := range v{
		if err := coder.DecodeStand(vdst);err!=nil{
			return err
		}
	}
	return nil
}