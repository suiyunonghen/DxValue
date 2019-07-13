package DxMsgPack


type  MsgPackCode		uint8

const(
	CodeUnkonw		MsgPackCode = 0
	PosFixedNumHigh MsgPackCode = 0x7f  //0-7f的正数最大编码
	NegFixedNumLow  MsgPackCode = 0xe0	//固定大小的负数编码

	CodeNil MsgPackCode = 0xc0

	CodeFalse MsgPackCode = 0xc2
	CodeTrue  MsgPackCode = 0xc3

	CodeFloat  MsgPackCode = 0xca
	CodeDouble MsgPackCode = 0xcb

	CodeUint8  MsgPackCode = 0xcc
	CodeUint16 MsgPackCode = 0xcd
	CodeUint32 MsgPackCode = 0xce
	CodeUint64 MsgPackCode = 0xcf

	CodeInt8  MsgPackCode = 0xd0
	CodeInt16 MsgPackCode = 0xd1
	CodeInt32 MsgPackCode = 0xd2
	CodeInt64 MsgPackCode = 0xd3

	CodeFixedStrLow  MsgPackCode = 0xa0
	CodeFixedStrHigh MsgPackCode = 0xbf
	FixedStrMask MsgPackCode = 0x1f
	CodeStr8         MsgPackCode = 0xd9
	CodeStr16        MsgPackCode = 0xda
	CodeStr32        MsgPackCode = 0xdb


	CodeBin8  MsgPackCode = 0xc4
	CodeBin16 MsgPackCode = 0xc5
	CodeBin32 MsgPackCode = 0xc6

	CodeFixedArrayLow  MsgPackCode = 0x90
	CodeFixedArrayHigh MsgPackCode = 0x9f
	FixedArrayMask MsgPackCode = 0xf
	CodeArray16        MsgPackCode = 0xdc
	CodeArray32        MsgPackCode = 0xdd

	CodeFixedMapLow  MsgPackCode = 0x80
	CodeFixedMapHigh MsgPackCode = 0x8f
	FixedMapMask MsgPackCode = 0xf
	CodeMap16        MsgPackCode = 0xde
	CodeMap32        MsgPackCode = 0xdf

	CodeFixExt1  MsgPackCode = 0xd4
	CodeFixExt2  MsgPackCode = 0xd5
	CodeFixExt4  MsgPackCode = 0xd6
	CodeFixExt8  MsgPackCode = 0xd7 //64位时间格式
	CodeFixExt16 MsgPackCode = 0xd8
	CodeExt8     MsgPackCode = 0xc7 //96位时间格式
	CodeExt16    MsgPackCode = 0xc8
	CodeExt32    MsgPackCode = 0xc9
)


func (code MsgPackCode)IsExt() bool {
	return (code >= CodeFixExt1 && code <= CodeFixExt16 && code != CodeFixExt4 && code != CodeFixExt8) || (code >= CodeExt8 && code <= CodeExt32)
}

func (code MsgPackCode)IsMap()bool  {
	return code >= CodeFixedMapLow && code <= CodeFixedMapHigh || code==CodeMap16 || code == CodeMap32
}

func (code MsgPackCode)IsFixedNum()bool  {
	return code <= PosFixedNumHigh || code >= NegFixedNumLow
}

func (code MsgPackCode)IsInt()bool  {
	return code <= PosFixedNumHigh || code >= NegFixedNumLow ||  code>=CodeUint8 && code<=CodeUint64 || code>=CodeInt8 && code<=CodeInt64
}

func (code MsgPackCode)IsStr()bool  {
	return code >= CodeFixedStrLow && code <= CodeFixedStrHigh || code>=CodeStr8 && code<=CodeStr32
}

func (code MsgPackCode)IsArray()bool  {
	return code >= CodeFixedArrayLow && code <= CodeFixedArrayHigh || code == CodeArray16 || code == CodeArray32
}

func (code MsgPackCode)IsBin()bool  {
	return code >= CodeBin8 && code <= CodeBin32
}