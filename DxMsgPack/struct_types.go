/*
针对常规对象的Struct结构的类型的注册
主要用来记录缓存结构的字段类型
Autor: 不得闲
QQ:75492895
 */
package DxMsgPack

import (
	"sync"
	"reflect"
	"strings"
)

type structField struct{
	name      string
	index     []int
	//encoder encoderFunc
	decoder 	decoderFunc
}

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

func (f *structField) DecodeValue(coderd *MsgPackDecoder, strct reflect.Value) error {
	return f.decoder(coderd, f.value(strct))
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

type structFields struct {
	List  []*structField
}

func (fs *structFields) Len() int {
	return len(fs.List)
}

func (fs *structFields) Add(field *structField) {
	fs.List = append(fs.List, field)
}

func (fs *structFields)FieldByName(fieldName string)*structField  {
	for _,field := range fs.List{
		if field.name == fieldName{
			return field
		}
	}
	return nil
}

func inlineFields(fs *structFields, typ reflect.Type, f *structField) bool {
	//var encoder encoderFunc
	var decoder decoderFunc

	if typ.Kind() == reflect.Struct {
		//encoder = f.encoder
		decoder = f.decoder
	} else {
		for typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
			//encoder = getEncoder(typ)
			decoder = getDecoder(typ)
		}
		if typ.Kind() != reflect.Struct {
			return false
		}
	}

	if reflect.ValueOf(encoder).Pointer() != encodeStructValuePtr {
		return false
	}
	if reflect.ValueOf(decoder).Pointer() != decodeStructValuePtr {
		return false
	}

	inlinedFields := getFields(typ,"msgpack").List
	for _, field := range inlinedFields {
		if fld := fs.FieldByName(field.name);fld==nil{
			field.index = append(f.index, field.index...)
			fs.Add(field)
		}
	}
	return true
}

func getFields(typ reflect.Type,coderName string) *structFields {
	numField := typ.NumField()
	fs := &structFields{
		List:  make([]*structField, 0, numField),
	}
	for i := 0; i < numField; i++ {
		f := typ.Field(i)
		if f.PkgPath != "" && !f.Anonymous {
			continue
		}
		MashName := f.Tag.Get(coderName)
		if MashName == "-"{ //不序列化
			continue
		}

		if MashName == "" {
			MashName = f.Name
		}
		field := &structField{
			name:      MashName,
			index:     f.Index,
			decoder:   getDecoder(f.Type),
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