package DxMsgPack

import (
	"encoding/binary"
	"github.com/suiyunonghen/DxValue"
	"unsafe"
	"github.com/suiyunonghen/DxCommonLib"
)

func (encoder *MsgPackEncoder)EncodeExtValue(v *DxValue.DxExtValue)(err error)  {
	btlen := 0
	bt := v.ExtData()
	btlen = len(bt)
	encoder.buf[1] = v.ExtType()
	switch {
	case btlen == 1:
		encoder.buf[0] = byte(CodeFixExt1)
		_,err = encoder.w.Write(encoder.buf[:1])
	case btlen == 2:
		encoder.buf[0] = byte(CodeFixExt2)
		_,err = encoder.w.Write(encoder.buf[:1])
	case btlen == 4:
		encoder.buf[0] = byte(CodeFixExt4)
		_,err = encoder.w.Write(encoder.buf[:1])
	case btlen == 8:
		encoder.buf[0] = byte(CodeFixExt8)
		_,err = encoder.w.Write(encoder.buf[:1])
	case btlen <= 16:
		encoder.buf[0] = byte(CodeFixExt16)
		_,err = encoder.w.Write(encoder.buf[:1])
	case btlen <= max_str8_len:
		encoder.buf[0] = byte(CodeExt8)
		encoder.buf[1] = byte(btlen)
		encoder.buf[2] = v.ExtType()
		_,err = encoder.w.Write(encoder.buf[:3])
	case btlen <= max_str16_len:
		encoder.buf[0] = byte(CodeExt16)
		binary.BigEndian.PutUint16(encoder.buf[1:3],uint16(btlen))
		encoder.buf[3] = v.ExtType()
		_,err = encoder.w.Write(encoder.buf[:4])
	default:
		if btlen > max_str32_len{
			btlen = max_str32_len
		}
		encoder.buf[0] = 0xc6
		binary.BigEndian.PutUint32(encoder.buf[1:5],uint32(btlen))
		encoder.buf[5] = v.ExtType()
		_,err = encoder.w.Write(encoder.buf[:6])
	}
	if err == nil && btlen > 0{
		_,err = encoder.w.Write(bt[:btlen])
	}
	return
}

func (encoder *MsgPackEncoder)rangeDxRecord(keyName string,value *DxValue.DxBaseValue)bool  {
	if encoder.curErr = encoder.EncodeString(keyName);encoder.curErr != nil{
		return false
	}
	//写入v
	if value != nil{
		encoder.curErr = encoder.Encode(value)
	}else{
		encoder.curErr = encoder.WriteByte(0xc0) //null
	}

	return encoder.curErr == nil
}

func (encoder *MsgPackEncoder)rangeDxIntRecord(key int64,value *DxValue.DxBaseValue)bool  {
	if encoder.curErr = encoder.EncodeInt(key);encoder.curErr != nil{
		return false
	}
	//写入v
	if value != nil{
		encoder.curErr = encoder.Encode(value)
	}else{
		encoder.curErr = encoder.WriteByte(0xc0) //null
	}
	return encoder.curErr == nil
}

func (encoder *MsgPackEncoder)EncodeRecord(r *DxValue.DxRecord)(err error)  {
	maplen := r.Length()
	if maplen <= max_fixmap_len{   //fixmap
		err = encoder.WriteByte(0x80 | byte(maplen))
	}else if maplen <= max_map16_len{
		//写入长度
		err = encoder.writeUint16(uint16(maplen),CodeMap16)
	}else{
		if maplen > max_map32_len{
			maplen = max_map32_len
		}
		err = encoder.writeUint32(uint32(maplen),CodeMap32)
	}
	if err != nil{
		return
	}
	//写入对象信息,Kv对
	encoder.curErr = nil
	r.Range(encoder.rangeDxRecord)
	err = encoder.curErr
	encoder.curErr = nil
	return err
}


func (encoder *MsgPackEncoder)EncodeRecordIntKey(r *DxValue.DxIntKeyRecord)(err error)  {
	maplen := r.Length()
	if maplen <= max_fixmap_len{   //fixmap
		err = encoder.WriteByte(0x80 | byte(maplen))
	}else if maplen <= max_map16_len{
		//写入长度
		err = encoder.writeUint16(uint16(maplen),CodeMap16)
	}else{
		if maplen > max_map32_len{
			maplen = max_map32_len
		}
		err = encoder.writeUint32(uint32(maplen),CodeMap32)
	}
	if err != nil{
		return
	}
	//写入对象信息,Kv对
	//写入对象信息,Kv对
	encoder.curErr = nil
	r.Range(encoder.rangeDxIntRecord)
	err = encoder.curErr
	encoder.curErr = nil
	return err
	return err
}

func (encoder *MsgPackEncoder)Encode(v *DxValue.DxBaseValue)(err error)  {
	switch v.ValueType() {
	case DxValue.DVT_Record:
		err = encoder.EncodeRecord((*DxValue.DxRecord)(unsafe.Pointer(v)))
	case DxValue.DVT_RecordIntKey:
		err = encoder.EncodeRecordIntKey((*DxValue.DxIntKeyRecord)(unsafe.Pointer(v)))
	case DxValue.DVT_Int:
		return encoder.EncodeInt(int64((*DxValue.DxIntValue)(unsafe.Pointer(v)).Int()))
	case DxValue.DVT_Int32:
		return encoder.EncodeInt(int64((*DxValue.DxInt32Value)(unsafe.Pointer(v)).Int32()))
	case DxValue.DVT_Int64:
		return encoder.EncodeInt((*DxValue.DxInt64Value)(unsafe.Pointer(v)).Int64())
	case DxValue.DVT_Bool:
		return encoder.EncodeBool((*DxValue.DxBoolValue)(unsafe.Pointer(v)).Bool())
	case DxValue.DVT_String:
		return encoder.EncodeString((*DxValue.DxStringValue)(unsafe.Pointer(v)).String())
	case DxValue.DVT_Float:
		return encoder.EncodeFloat((*DxValue.DxFloatValue)(unsafe.Pointer(v)).Float())
	case DxValue.DVT_Double:
		return encoder.EncodeDouble((*DxValue.DxDoubleValue)(unsafe.Pointer(v)).Double())
	case DxValue.DVT_Binary:
		bt := (*DxValue.DxBinaryValue)(unsafe.Pointer(v)).Bytes()
		if bt != nil{
			return encoder.EncodeBinary(bt)
		}else{
			return encoder.WriteByte(0xc0)
		}
	case DxValue.DVT_Ext:
		return encoder.EncodeExtValue((*DxValue.DxExtValue)(unsafe.Pointer(v)))
	case DxValue.DVT_Array:
		return encoder.EncodeArray((*DxValue.DxArray)(unsafe.Pointer(v)))
	case DxValue.DVT_DateTime:
		return encoder.EncodeDateTime(DxCommonLib.TDateTime((*DxValue.DxDoubleValue)(unsafe.Pointer(v)).Double()))
	default:
		return encoder.WriteByte(0xc0) //null
	}
	return nil
}


func (encoder *MsgPackEncoder)EncodeArray(arr *DxValue.DxArray)(err error)  {
	arlen := arr.Length()
	switch {
	case arlen < 16: //1001XXXX|    N objects
		err = encoder.WriteByte(byte(CodeFixedArrayLow) | byte(arlen))
	case arlen <= max_map16_len:  //0xdc  |YYYYYYYY|YYYYYYYY|    N objects
		encoder.writeUint16(uint16(arlen),CodeArray16)
	default:
		if arlen > max_map32_len{
			arlen = max_map32_len
		}
		encoder.writeUint32(uint32(arlen),CodeArray32)
	}

	for i := 0;i <= arlen - 1;i++{
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