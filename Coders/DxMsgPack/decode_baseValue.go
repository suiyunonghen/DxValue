package DxMsgPack

import (
	"encoding/binary"
	"github.com/suiyunonghen/DxCommonLib"
	"unsafe"
	"time"
	"io"
	"bufio"
	"errors"
	"github.com/suiyunonghen/DxValue/Coders"
)


type bufReader interface{
	io.ByteScanner
	io.Reader
}

type MsgPackDecoder   struct{
	r  					bufReader
	buffer				[64]byte //内部一个缓存，主要用来做一些数据读取转换
	OnParserNormalValue	func(v interface{})
	OnStartMap			func(mapLen int,keyIsStr bool)(mapInterface interface{}) //开始一个字符串键值Map，指定返回的Map结构对象
	OnParserStrMapKv    func(mapInterface interface{},key string,v interface{})
	OnParserIntKeyMapKv func(mapInterface interface{},intKey int64,v interface{})
	OnStartArray		func(arrLen int)(arrInterface interface{}) //开始数组时候触发
	OnStartStrMapArray	func(arrLen int,Key string,mapInterface interface{})(arrInterface interface{}) //开始数组时候触发
	OnStartIntMapArray	func(arrLen int,Key int64,mapInterface interface{})(arrInterface interface{}) //开始数组时候触发
	OnParserArrElement	func(arrInterface interface{},index int,v interface{}) //解析数组元素触发
	OnParserArrObject   func(arrInterface interface{},index int)(object interface{}) //数组中是复杂对象
}

var(
	ErrInvalidateMsgPack = errors.New("Invalidate MsgPack Format")
	ErrInvalidateMap	= errors.New("Invalidate Map Format")
	ErrInvalidateMapKey	= errors.New("Invalidate Map Key,Key Can Only Int or String")
	ErrInvalidateArrLen = errors.New("Is not a Array Len Flag")
)

const msgPackName = "msgpack"

func setStringsCap(s []string, n int) []string {
	if n > 256 {
		n = 256
	}

	if s == nil {
		return make([]string, 0, n)
	}

	if cap(s) >= n {
		return s[:0]
	}

	s = s[:cap(s)]
	s = append(s, make([]string, n-len(s))...)
	return s[:0]
}

func (coder *MsgPackDecoder)Name()string  {
	return msgPackName
}

func (coder *MsgPackDecoder)Read(b []byte)error  {
	_,err := coder.r.Read(b)
	return err
}

func (coder *MsgPackDecoder)UnreadByte() error {
	return coder.r.UnreadByte()
}

func (coder *MsgPackDecoder)ReadBigEnd16()(uint16,error)  {
	if _,err := coder.r.Read(coder.buffer[:2]);err!=nil{
		return 0,err
	}
	return binary.BigEndian.Uint16(coder.buffer[:2]),nil
}

func (coder *MsgPackDecoder)ReadBigEnd32()(uint32,error)  {
	if _,err := coder.r.Read(coder.buffer[:4]);err!=nil{
		return 0,err
	}
	return binary.BigEndian.Uint32(coder.buffer[:4]),nil
}

func (coder *MsgPackDecoder)ReadBigEnd64()(uint64,error)  {
	if _,err := coder.r.Read(coder.buffer[:8]);err!=nil{
		return 0,err
	}
	return binary.BigEndian.Uint64(coder.buffer[:8]),nil
}

func (coder *MsgPackDecoder)ReadCode()(MsgPackCode,error)  {
	c, err := coder.r.ReadByte()
	if err != nil {
		return 0, err
	}
	return MsgPackCode(c), nil
}

func (coder *MsgPackDecoder) hasNilCode() bool {
	if code,err := coder.ReadCode();err==nil{
		coder.r.UnreadByte()
		return code == CodeNil
	}
	return false
}

func (coder *MsgPackDecoder)DecodeDateTime(code MsgPackCode)(DxCommonLib.TDateTime,error)  {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
			return -1,err
		}
	}
	switch code {
	case CodeFloat:
		if v32,err := coder.ReadBigEnd32();err!=nil{
			return -1,err
		}else{
			return DxCommonLib.TDateTime(*(*float32)(unsafe.Pointer(&v32))),nil
		}
	case CodeDouble:
		if v64,err := coder.ReadBigEnd64();err!=nil{
			return -1,err
		}else{
			return DxCommonLib.TDateTime(*(*float32)(unsafe.Pointer(&v64))),nil
		}
	case CodeFixExt4:
		var b byte
		if b,err = coder.r.ReadByte();err!=nil{
			return -1,err
		}
		if int8(b) == -1{
			if ms,err := coder.ReadBigEnd32();err!=nil{
				return -1,err
			}else{
				ntime := time.Now()
				ns := ntime.Unix()
				ntime = ntime.Add((time.Duration(int64(ms) - ns)*time.Second))
				return DxCommonLib.Time2DelphiTime(&ntime),nil
			}
		}
	}
	return -1,Coders.ErrValueType
}

func (coder *MsgPackDecoder)DecodeDateTime_Go(code MsgPackCode)(time.Time,error)  {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
			return time.Time{},err
		}
	}
	switch code {
	case CodeFloat:
		if v32,err := coder.ReadBigEnd32();err!=nil{
			return time.Time{},err
		}else{
			return DxCommonLib.TDateTime(*(*float32)(unsafe.Pointer(&v32))).ToTime(),nil
		}
	case CodeDouble:
		if v64,err := coder.ReadBigEnd64();err!=nil{
			return time.Time{},err
		}else{
			return DxCommonLib.TDateTime(*(*float32)(unsafe.Pointer(&v64))).ToTime(),nil
		}
	case CodeFixExt4:
		var b byte
		if b,err = coder.r.ReadByte();err!=nil{
			return time.Time{},err
		}
		if int8(b) == -1{
			if ms,err := coder.ReadBigEnd32();err!=nil{
				return time.Time{},err
			}else{
				ntime := time.Now()
				ns := ntime.Unix()
				return ntime.Add((time.Duration(int64(ms) - ns)*time.Second)),nil
			}
		}
	}
	return time.Time{},Coders.ErrValueType
}

func (coder *MsgPackDecoder)DecodeBool(code MsgPackCode)(bool,error)  {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
			return false,err
		}
	}
	switch code {
	case CodeTrue:
		return true,nil
	case CodeFalse:
		return false,nil
	default:
		return false,errors.New("invalidate Bool type")
	}
}

func (coder *MsgPackDecoder)DecodeInt(code MsgPackCode)(int64,error)  {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
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
		if vwrod,err := coder.ReadBigEnd16();err!=nil{
			return 0,err
		}else{
			if code == CodeInt16{
				return int64(int16(vwrod)),nil
			}else{
				return int64(vwrod),nil
			}
		}
	case CodeInt32,CodeUint32:
		if vuint32,err := coder.ReadBigEnd32();err!=nil{
			return 0,err
		} else{
			if code == CodeInt32{
				return int64(int32(vuint32)),nil
			}else{
				return int64(vuint32),nil
			}
		}
	case CodeInt64,CodeUint64:
		if vuint64,err := coder.ReadBigEnd64();err!=nil{
			return 0,err
		}else{
			return int64(vuint64),nil
		}
	}
	return 0,Coders.ErrValueType
}

func (coder *MsgPackDecoder)DecodeFloat(code MsgPackCode)(float64,error) {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
			return 0,err
		}
	}
	switch code {
	case CodeFloat:
		if v32,err := coder.ReadBigEnd32();err!=nil{
			return 0,err
		}else{
			return float64(*(*float32)(unsafe.Pointer(&v32))),nil
		}

	case CodeDouble:
		if v64,err := coder.ReadBigEnd64();err!=nil{
			return 0,err
		}else{
			return *(*float64)(unsafe.Pointer(&v64)),nil
		}

	default:
		if iv,err := coder.DecodeInt(code);err!=nil{
			return 0,err
		}else{
			return float64(iv),nil
		}
	}
	return 0,Coders.ErrValueType
}

func (coder *MsgPackDecoder)BinaryLen(code MsgPackCode)(int,error)  {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
			return -1,err
		}
	}
	btlen := 0
	switch code {
	case CodeBin8:
		if b,err := coder.r.ReadByte();err!=nil{
			return -1,err
		}else{
			btlen = int(b)
		}
	case CodeBin16:
		if v16,err := coder.ReadBigEnd16();err!=nil{
			return -1,err
		}else{
			btlen = int(v16)
		}
	case CodeBin32:
		if v32,err := coder.ReadBigEnd32();err!=nil{
			return -1,err
		}else{
			btlen = int(v32)
		}
	default:
		return -1,Coders.ErrValueType
	}
	return btlen,nil
}

func (coder *MsgPackDecoder)DecodeBinary(code MsgPackCode)([]byte,error)  {
	btlen,err := coder.BinaryLen(code)
	if btlen > 0{
		mb := make([]byte,btlen)
		if _,err = coder.r.Read(mb);err!=nil{
			return nil,err
		}
		return mb,nil
	}
	return nil,err
}

func (coder *MsgPackDecoder)DecodeExtValue(code MsgPackCode)([]byte,error) {
	btlen := -1
	var err error
	if code == CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
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
		if v16,err := coder.ReadBigEnd16();err!=nil{
			return nil,err
		}else{
			btlen = int(v16)+1
		}
	case CodeExt32:
		if v32,err := coder.ReadBigEnd32();err!=nil{
			return nil,err
		}else{
			btlen = int(v32) + 1
		}
	}
	if btlen <= 0{
		return nil,ErrInvalidateMsgPack
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
		if mapcode,err = coder.ReadCode();err!=nil{
			return 0,err
		}
	}
	switch mapcode {
	case CodeMap16:
		if v16,err := coder.ReadBigEnd16();err!=nil{
			return 0,err
		} else{
			return int(v16),nil
		}
	case CodeMap32:
		if v32,err := coder.ReadBigEnd32();err!=nil{
			return 0,err
		} else{
			return int(v32),nil
		}
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
		if code,err = coder.ReadCode();err!=nil{
			return 0,err
		}
	}
	switch code {
	case CodeArray16:
		if v16,err := coder.ReadBigEnd16();err!=nil{
			return 0,err
		} else{
			return int(v16),nil
		}
	case CodeArray32:
		if v32,err := coder.ReadBigEnd32();err!=nil{
			return 0,err
		} else{
			return int(v32),nil
		}
	default:
		if code >= CodeFixedArrayLow && code <= CodeFixedArrayHigh{
			return int(code & FixedArrayMask),nil
		}
	}
	return 0,ErrInvalidateArrLen
}

func (coder *MsgPackDecoder)ReSetReader(r io.Reader)  {
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
}

func (coder *MsgPackDecoder)skipN(nbyte int)(err error)  {
	if seeker,ok := coder.r.(io.Seeker);ok{
		_,err = seeker.Seek(int64(nbyte),io.SeekCurrent)
		return err
	}
	var buf []byte
	if nbyte <= len(coder.buffer){
		buf = coder.buffer[:nbyte]
	}else{
		buf = make([]byte,nbyte)
	}
	_,err = coder.r.Read(buf)
	return err
}

func (coder *MsgPackDecoder)Skip()(error)  {
	return coder.SkipByCode(CodeUnkonw)
}

func (coder *MsgPackDecoder)SkipByCode(code MsgPackCode)(err error)  {
	if code == CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
			return err
		}
	}
	if code.IsStr(){
		return coder.skipString(code)
	}else if code.IsArray(){
		return coder.skipArray(code)
	}else if code.IsBin(){
		return coder.skipBinary(code)
	}else if code.IsExt(){
		return coder.skipExtValue(code)
	}else if code.IsInt(){
		_,err = coder.DecodeInt(code)
		return err
	}else if code.IsMap(){
		return coder.skipMap(code)
	}else{
		switch code {
		case CodeTrue:	return nil
		case CodeFalse: return nil
		case CodeNil:   return nil
		case CodeFloat:
			_,err = coder.ReadBigEnd32()
			return err
		case CodeDouble:
			_,err = coder.ReadBigEnd64()
			return err
		case CodeFixExt4:
			if code,err = coder.ReadCode();err!=nil{
				return err
			}
			_,err = coder.ReadBigEnd32()
			return err
		}
	}
	return nil
}

func (coder *MsgPackDecoder)skipStrMapKvRecord(strcode MsgPackCode)(error)  {
	err := coder.skipString(strcode)
	if err != nil{
		return err
	}
	if strcode,err = coder.ReadCode();err!=nil{
		return  err
	}
	if strcode.IsStr(){
		return coder.skipString(strcode)
	}else if strcode.IsFixedNum(){
		return nil
	}else if strcode.IsInt(){
		_,err = coder.DecodeInt(strcode)
		return err
	}else if strcode.IsMap(){
		return coder.skipMap(strcode)
	}else if strcode.IsArray(){
		return coder.skipArray(strcode)
	}else if strcode.IsBin(){
		return coder.skipBinary(strcode)
	}else if strcode.IsExt(){
		return coder.skipExtValue(strcode)
	}else{
		switch strcode {
		case CodeTrue:	return nil
		case CodeFalse: return nil
		case CodeNil:   return nil
		case CodeFloat:
			_,err = coder.ReadBigEnd32()
			return err
		case CodeDouble:
			_,err = coder.ReadBigEnd64()
			return err
		case CodeFixExt4:
			if strcode,err = coder.ReadCode();err!=nil{
				return err
			}
			_,err = coder.ReadBigEnd32()
			return err
		}
	}
	return err
}

func (coder *MsgPackDecoder)skipIntKeyMapKvRecord(intkeyCode MsgPackCode)(error)  {
	_,err := coder.DecodeInt(intkeyCode)
	if err != nil{
		return err
	}
	if intkeyCode,err = coder.ReadCode();err!=nil{
		return err
	}

	if intkeyCode.IsStr(){
		return coder.skipString(intkeyCode)
	}else if intkeyCode.IsFixedNum(){
		return nil
	}else if intkeyCode.IsInt(){
		_,err = coder.DecodeInt(intkeyCode)
		return err
	}else if intkeyCode.IsMap(){
		return coder.skipMap(intkeyCode)
	}else if intkeyCode.IsArray(){
		return coder.skipArray(intkeyCode)
	}else if intkeyCode.IsBin(){
		return coder.skipBinary(intkeyCode)
	}else if intkeyCode.IsExt(){
		return coder.skipExtValue(intkeyCode)
	}else{
		switch intkeyCode {
		case CodeTrue:	return nil
		case CodeFalse: return nil
		case CodeNil:   return nil
		case CodeFloat:
			_,err = coder.ReadBigEnd32()
			return err
		case CodeDouble:
			_,err = coder.ReadBigEnd64()
			return err
		case CodeFixExt4:
			if intkeyCode,err = coder.ReadCode();err!=nil{
				return  err
			}
			_,err = coder.ReadBigEnd32()
			return err
		}
	}
	return err
}

func (coder *MsgPackDecoder)skipMap(strcode MsgPackCode)(error)  {
	if maplen,err := coder.DecodeMapLen(strcode);err!=nil{
		return err
	}else{
		//判断键值，是Int还是str
		if strcode,err = coder.ReadCode();err!=nil{
			return err
		}
		if strcode.IsInt(){
			if err = coder.skipIntKeyMapKvRecord(strcode);err!=nil{
				return err
			}
			for j := 1;j<maplen;j++{
				if err = coder.skipIntKeyMapKvRecord(CodeUnkonw);err!=nil{
					return err
				}
			}
		}else if strcode.IsStr(){
			if err = coder.skipStrMapKvRecord(strcode);err!=nil{
				return err
			}
			for j := 1;j<maplen;j++{
				if err = coder.skipStrMapKvRecord(CodeUnkonw);err!=nil{
					return err
				}
			}
		}
		return nil
	}
}

func (coder *MsgPackDecoder)skipExtValue(code MsgPackCode)(error) {
	btlen := -1
	var err error
	if code == CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
			return err
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
			return err
		}else{
			btlen = int(blen)+1
		}
	case CodeExt16:
		if v16,err := coder.ReadBigEnd16();err!=nil{
			return err
		}else{
			btlen = int(v16)+1
		}
	case CodeExt32:
		if v32,err := coder.ReadBigEnd32();err!=nil{
			return err
		}else{
			btlen = int(v32) + 1
		}
	}
	if btlen <= 0{
		return nil
	}
	return coder.skipN(btlen)
}

func (coder *MsgPackDecoder)skipBinary(code MsgPackCode)(error)  {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
			return err
		}
	}
	btlen := 0
	switch code {
	case CodeBin8:
		if b,err := coder.r.ReadByte();err!=nil{
			return err
		}else{
			btlen = int(b)
		}
	case CodeBin16:
		if v16,err := coder.ReadBigEnd16();err!=nil{
			return err
		}else{
			btlen = int(v16)
		}
	case CodeBin32:
		if v32,err := coder.ReadBigEnd32();err!=nil{
			return err
		}else{
			btlen = int(v32)
		}
	default:
		return Coders.ErrValueType
	}
	if btlen > 0{
		return coder.skipN(btlen)
	}
	return nil
}


func (coder *MsgPackDecoder)skipArray(code MsgPackCode)(error)  {
	var (
		err error
		arrlen int
	)
	if code == CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
			return err
		}
	}
	if arrlen,err = coder.DecodeArrayLen(code);err!=nil{
		return err
	}
	for i := 0;i<arrlen;i++{
		if err = coder.SkipByCode(CodeUnkonw);err!=nil{
			return err
		}
	}
	return nil
}


func (coder *MsgPackDecoder)skipString(code MsgPackCode)(error) {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
			return err
		}
	}
	strlen := 0
	switch code {
	case CodeStr8:
		if bl,err := coder.r.ReadByte();err!=nil{
			return err
		}else{
			strlen = int(bl)
		}
	case CodeStr16:
		if v16,err := coder.ReadBigEnd16();err!=nil{
			return err
		} else{
			strlen = int(v16)
		}
	case CodeStr32:
		if v32,err := coder.ReadBigEnd32();err!=nil{
			return err
		} else{
			strlen = int(v32)
		}
	default:
		if code < 0xa0 || code> 0xbf {
			return Coders.ErrValueType
		}
		strlen = int(code & FixedStrMask)
	}
	if strlen > 0{
		return coder.skipN(strlen)
	}
	return nil
}

func (coder *MsgPackDecoder)DecodeString(code MsgPackCode)([]byte,error) {
	var err error
	if code == CodeUnkonw{
		if code,err = coder.ReadCode();err!=nil{
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
		if v16,err := coder.ReadBigEnd16();err!=nil{
			return nil, err
		} else{
			strlen = int(v16)
		}
	case CodeStr32:
		if v32,err := coder.ReadBigEnd32();err!=nil{
			return nil, err
		} else{
			strlen = int(v32)
		}
	default:
		if code < 0xa0 || code> 0xbf {
			return nil,Coders.ErrValueType
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

func NewDecoder(r io.Reader)*MsgPackDecoder  {
	var result MsgPackDecoder
	if bytebf,ok := r.(bufReader);ok {
		result.r = bytebf
	}else {
		if bf,ok := r.(*bufio.Reader);ok{
			bf.Reset(r)
			result.r = bf
		}else{
			result.r = bufio.NewReader(r)
		}
	}
	return &result
}


