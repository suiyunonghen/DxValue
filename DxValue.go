package DxValue

import "unsafe"

/******************************************************
*  DxValue
******************************************************/
type 	DxValue		struct{
	fValue			*DxBaseValue
}

func (v *DxValue)ValueType()DxValueType  {
	if v.fValue == nil{
		return DVT_Null
	}
	return v.fValue.ValueType()
}

func (v *DxValue)SetIntValue(value int)  {
	if v.fValue != nil{
		if v.fValue.ValueType() == DVT_Int{
			(*DxIntValue)(unsafe.Pointer(v.fValue)).fvalue = value
			return
		}
	}
	var m DxIntValue
	m.fValueType = DVT_Int
	m.fvalue = value
	v.fValue = &m.DxBaseValue
}

func (v *DxValue)SetVariant(value DxValue)  {
	v.fValue = value.fValue
}

func (v *DxValue)ClearValue()  {
	v.fValue = nil
}

