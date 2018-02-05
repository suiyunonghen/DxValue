package DxMsgPack

import (
	"time"
	"reflect"
	"errors"
	"github.com/suiyunonghen/DxValue"
	"github.com/suiyunonghen/DxValue/Coders"
)


func (encoder *MsgPackEncoder)EncodeStand(v interface{})(error)  {
	switch value := v.(type) {
	case *string:
		return encoder.EncodeString(*value)
	case string:
		return encoder.EncodeString(value)
	case *[]interface{}:
	case []interface{}:
		//return coder.DecodeArray2StdSlice(CodeUnkonw,value)
	case *DxValue.DxBaseValue:
		return encoder.Encode(value)
	case DxValue.DxBaseValue:
		return encoder.Encode(&value)
	case *time.Time:
		return encoder.EncodeTime(*value)
	case time.Time:
		return encoder.EncodeTime(value)
	case *int8:
		return encoder.EncodeInt(int64(*value))
	case int8:
		return encoder.EncodeInt(int64(value))
	case *int16:
		return encoder.EncodeInt(int64(*value))
	case int16:
		return encoder.EncodeInt(int64(value))
	case *int32:
		return encoder.EncodeInt(int64(*value))
	case int32:
		return encoder.EncodeInt(int64(value))
	case *int64:
		return encoder.EncodeInt(*value)
	case int64:
		return encoder.EncodeInt(value)
	case *uint8:
		return encoder.EncodeInt(int64(*value))
	case uint8:
		return encoder.EncodeInt(int64(value))
	case *uint16:
		return encoder.EncodeInt(int64(*value))
	case uint16:
		return encoder.EncodeInt(int64(value))
	case *uint32:
		return encoder.EncodeInt(int64(*value))
	case uint32:
		return encoder.EncodeInt(int64(value))
	case *uint64:
		return encoder.EncodeInt(int64(*value))
	case uint64:
		return encoder.EncodeInt(int64(value))
	case *float32:
		return encoder.EncodeFloat(*value)
	case *float64:
		return encoder.EncodeDouble(*value)
	case float32:
		return encoder.EncodeFloat(value)
	case float64:
		return encoder.EncodeDouble(value)
	case *bool:
		return encoder.EncodeBool(*value)
	case bool:
		return encoder.EncodeBool(value)
	case *[]byte:
		return encoder.EncodeBinary(*value)
	case *map[string]interface{}:
		//return coder.decodeStrMapFunc(value)
	case *map[int]interface{}:
		//return coder.decodeIntKeyMapFunc(value)
	case *map[int64]interface{}:
		//return coder.decodeIntKeyMapFunc64(value)
	case *map[string]string:
		//return coder.decodeStrValueMapFunc(value)
	case *time.Duration:
		return encoder.EncodeInt(int64(*value))
	default:
		v := reflect.ValueOf(value)
		if !v.IsValid() {
			return errors.New("msgpack: Decode(nil)")
		}
		rv := Coders.GetRealValue(&v)
		if rv == nil{
			encoder.WriteByte(0xc0) //null
		}
		v = v.Elem()
		switch rv.Kind(){
		case reflect.Struct:
		case reflect.Map:
		case reflect.Slice,reflect.Array:

		}
	}
	return nil
}
