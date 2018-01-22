/*
DxValue的Array数组对象
可以用来序列化反序列化Json,MsgPack等，并提供一系列的操作函数
Autor: 不得闲
QQ:75492895
 */
package DxValue


/******************************************************
*  DxArray
******************************************************/
type  DxArray		struct{
	DxBaseValue
	fRecords		[]*DxBaseValue
}


