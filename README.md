# DxValue
一个万能值复合变量，其目的是将Json,Msgpack这类数据协议格式综合，提供一个复合变量，在变量内部使用Go的常规简单数据类型进行操作管理（将各类数据类型直接序列化到本
本变量内部存储，不采用反射加大性能，内部采用KV结构模型。）
整体数据结构如下：      
![image](https://github.com/suiyunonghen/DxValue/blob/master/DxValueStruct.png)

整体有效常用对象主要有3个，DxRecord,DxArray,DxValue
1. DxRecord对象
    - 记录集对象  
    - 支持针对Json的Object模式格式的编码解码，并生成记录对象内容
    - 支持针对MsgPack的Map模式格式的编码解码，并生成记录对象内容
    - 记录值可以任意嵌套
    - 可以包含任意Json,MsgPack格式支持的数据类型，使用SetValue(K,v)赋，或者使用SetInt,SetInt32,SetArray等强制类型函数赋值    
    - KV存储，采用Key查找对应的值，具备有AsInt(Key),AsBool(key),AsString(Key)等获取相关值的函数
    - 支持路径模式查找获取和路径模式创建，比如AsStringByPath，AsRecordByPath等
    - 使用ForcePath来创建路径并赋值
    - JsonParserFromByte用来将Json串解码，其中参数2主要用来设定是否针对字符串做自动转义检查转码并解码转义字符
    > 路径分隔采用对象的PathSplitChar来设置分隔符，比如，JSon格式如下
    ```json
    {"BoolValue":  true  ,"object":{"objBool":  false  }}
    ```
    >设置PathSplitChar='.'(PathSplitChar默认值是.) ，那么可以使用object.objBool来获取objBool的值，如下
    ```go
    rc := NewRecord()
	rc.JsonParserFromByte([]byte(`{"BoolValue":  true  ,"object":{"objBool":  false  }}`),false)
	fmt.Println("BoolValue=",rc.AsBool("BoolValue",false))
	fmt.Println("object.objBool=",rc.AsBoolByPath("object.objBool",true))
    ```
    > 使用ForcePath来创建路径并赋值，本函数在路径存在的时候，直接赋值，如果不存在那么会创建路径然后赋值比如：
    ```go
    rc := NewRecord()
    rc.SetBool("BoolValue",true)
    rc.ForcePath("object.objBool",false)
    fmt.Println(rc.ToString()
    ```
    > 使用本功能则可以获得以上的Json格式字符串

    > 使用LoadJsonFile加载Json文件如下
    ```go
    rec := NewRecord()
    rec.LoadJsonFile("DataProxy.config.json",true) //参数2指定是否自动解析转义符
    ```
    > 使用SaveJsonFile保存内容到Json格式文件
    ```go
    rec := NewRecord()
    rec.SetInt("Age",12)
    rec.SetString("Name","suiyunonghen")
    rec.SetValue("Home",map[string]interface{}{
        "Addres": "湖北武汉",
        "code":"430000",
        "Peoples":4,
    })
    rec.SaveJsonFile("d:\\testJson.json",true)
    ```
    
2. DxArray对象
    - 数组对象  
    - 支持针对Json的数组对象格式的编码解码，并生成数组对象内容
    - 支持针对MsgPack的数组对象格式的编码解码，并生成数组对象内容
    - 数组值可以任意嵌套
    - 可以包含任意Json,MsgPack格式支持的数据类型，使用SetValue(idx,v)赋，或者使用SetInt,SetInt32,SetArray等强制类型函数赋值  
    - JsonParserFromByte用来将Json串解码，其中参数2主要用来设定是否针对字符串做自动转义检查转码并解码转义字符
    - 具备有AsInt(idx),AsBool(idx),AsString(idx)等获取相关值的函数
    - 使用LoadJsonFile加载Json文件
    - 使用SaveJsonFile保存内容到Json格式文件
    >用法如下：
    ```go
    arr := NewArray()
	  _,err := arr.JsonParserFromByte([]byte(`[  32  ,  "2342"  ,[ 2 , true , false  ,{ "Name" : "DxSoft" , "Age"  :  32 } ] ]`),false)
    if err == nil {
		fmt.Println(arr.ToString())
	}else{
		fmt.Println("Paser Error")
	}
    ```
2. DxValue对象
    - 万能值对象  
    - 支持任意Json,MsgPack对象格式的编码解码，并且生成对象内容
    - 用法使用DxRecore和DxArray结合使用
    - 使用LoadJsonFile加载Json文件（参考Record)
    - 使用SaveJsonFile保存内容到Json格式文件（参考Record)
    - JsonParserFromByte用来将Json串解码，其中参数2主要用来设定是否针对字符串做自动转义检查转码并解码转义字符，并且自动识别JSON格式
    ```go
    var v DxValue
    v.LoadJsonFile("DataProxy.config.json",true)
    fmt.Println(v.AsString())
    if rec,_ := v.AsRecord();rec != nil{
        if arr := rec.AsArray("list");arr!=nil{
            arr.SaveJsonFile("d:\\1.json",true)
            fmt.Println("SaveJsonOK")
        }
    }
    ```
