/*
针对常规对象的Struct结构的类型的注册
主要用来记录缓存结构的字段类型
Autor: 不得闲
QQ:75492895
 */
package Coders

import (
	"sync"
	"reflect"
	"time"
)

type structField struct{
	name      string
	index     []int
	fieldtype	reflect.Type
	mashNames	  map[string]string		//序列化的CoderName对应的建值，比如msgpack:"", json:""
}

var(
	InterfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
	StringType = reflect.TypeOf((*string)(nil)).Elem()

	IntType = reflect.TypeOf(int(0))
	Int64Type= reflect.TypeOf(int64(0))

	MapStringStringPtrType = reflect.TypeOf((*map[string]string)(nil))
	MapStringStringType = MapStringStringPtrType.Elem()
	MapStringInterfacePtrType = reflect.TypeOf((*map[string]interface{})(nil))
	MapStringInterfaceType = MapStringInterfacePtrType.Elem()
	MapIntStringPtrType = reflect.TypeOf((*map[int]string)(nil))
	MapIntStringType = MapIntStringPtrType.Elem()

	MapIntInterfacePtrType = reflect.TypeOf((*map[int]interface{})(nil))
	MapIntInterfaceType = MapIntInterfacePtrType.Elem()

	TimePtrType = reflect.TypeOf((*time.Time)(nil))
	TimeType = TimePtrType.Elem()
	ValueCoderType = reflect.TypeOf((*ValueCoder)(nil)).Elem()
	ErrorType = reflect.TypeOf((*error)(nil)).Elem()
)




func (f *structField) value(v reflect.Value) reflect.Value {
	if len(f.index) == 1 {
		return v.Field(f.index[0])
	}
	for i, x := range f.index {
		if i > 0 {
			var ok bool
			v, ok = indirectNew(v)
			if !ok {
				return v
			}
		}
		v = v.Field(x)
	}
	return v
}

func (f *structField) DecodeValue(coder Decoder, strct reflect.Value) error {
	return coder.GetDecoderFunc(f.fieldtype)(coder,f.value(strct))
}

func (f *structField)EncodeValue(coder Encoder,strct reflect.Value)error  {
	return coder.GetEncoderFunc(f.fieldtype)(coder,f.value(strct))
}

func (f *structField)MarshalName(coderName string)string  {
	if coderName != ""{
		return f.mashNames[coderName]
	}
	return f.name
}

func indirectNew(v reflect.Value) (reflect.Value, bool) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			if !v.CanSet() {
				return v, false
			}
			elemType := v.Type().Elem()
			if elemType.Kind() != reflect.Struct {
				return v, false
			}
			v.Set(reflect.New(elemType))
		}
		v = v.Elem()
	}
	return v, true
}

func GetRealValue(v *reflect.Value)*reflect.Value  {
	if !v.IsValid(){
		return nil
	}
	if v.Kind() == reflect.Ptr{
		if !v.IsNil(){
			va := v.Elem()
			return GetRealValue(&va)
		}else{
			return nil
		}
	}
	return v
}

type structMap struct {
	mu sync.RWMutex
	m  map[reflect.Type]*structFields
}

func newStructMap() *structMap {
	return &structMap{
		m: make(map[reflect.Type]*structFields),
	}
}

var structTypeMap = newStructMap()

func Fields(typ reflect.Type) *structFields {
	structTypeMap.mu.RLock()
	fs, ok := structTypeMap.m[typ]
	structTypeMap.mu.RUnlock()
	if ok {
		return fs
	}

	structTypeMap.mu.Lock()
	fs, ok = structTypeMap.m[typ]
	if !ok {
		fs = getFields(typ)
		structTypeMap.m[typ] = fs
	}
	structTypeMap.mu.Unlock()

	return fs
}

type structFields struct {
	List  []*structField
}

func (fs *structFields) Len(coderName string) int {
	if coderName == ""{
		return len(fs.List)
	}
	mapLen := 0
	for _,field := range fs.List{
		if mname,ok := field.mashNames[coderName];ok && mname != "-"{
			mapLen++
		}
	}
	return mapLen
}

func (fs *structFields)Field(idx int)*structField  {
	if idx>=0 && idx < len(fs.List){
		return fs.List[idx]
	}
	return nil
}

func (fs *structFields)Range(coderName string,fieldIteaFunc func(field *structField)bool)  {
	if fieldIteaFunc == nil{
		return
	}
	if coderName == ""{
		for _,field := range fs.List{
			if !fieldIteaFunc(field){
				return
			}
		}
	}else{
		for _,field := range fs.List{
			if mname,ok := field.mashNames[coderName];ok && mname != "-"{
				if !fieldIteaFunc(field){
					return
				}
			}
		}
	}
}

func (fs *structFields) Add(field *structField) {
	fs.List = append(fs.List, field)
}

func (fs *structFields)FieldByName(coderName,fieldName string)*structField  {
	if coderName!=""{
		for _,field := range fs.List{
			if mname,ok := field.mashNames[coderName];ok && mname == fieldName{
				return field
			}
		}
	}else{
		for _,field := range fs.List{
			if field.name == fieldName{
				return field
			}
		}
	}
	return nil
}

func inlineFields(fs *structFields, typ reflect.Type, f *structField) bool {
	inlinedFields := getFields(typ).List
	for _, field := range inlinedFields {
		if fld := fs.FieldByName("",field.name);fld==nil{
			field.index = append(f.index, field.index...)
			fs.Add(field)
		}
	}
	return true
}

func getFields(typ reflect.Type) *structFields {
	numField := typ.NumField()
	fs := &structFields{
		List:  make([]*structField, 0, numField),
	}
	for i := 0; i < numField; i++ {
		f := typ.Field(i)
		if f.PkgPath != "" && !f.Anonymous {
			continue
		}
		field := &structField{
			index:     f.Index,
			fieldtype:	f.Type,
			name:		f.Name,
		}
		field.mashNames = make(map[string]string,len(coders))
		for _,coderName := range coders{
			MashName := f.Tag.Get(coderName)
			if MashName == "" {
				MashName = f.Name
			}
			field.mashNames[coderName] = MashName
		}

		if f.Anonymous && inlineFields(fs, f.Type, field) {
			continue
		}

		fs.Add(field)
	}
	return fs
}



func RegisterType(typ reflect.Type)  {
	structTypeMap.mu.RLock()
	_, ok := structTypeMap.m[typ]
	structTypeMap.mu.RUnlock()
	if !ok {
		structTypeMap.mu.Lock()
		fs, ok := structTypeMap.m[typ]
		if !ok {
			fs = getFields(typ)
			structTypeMap.m[typ] = fs
		}
		structTypeMap.mu.Unlock()
	}
}
