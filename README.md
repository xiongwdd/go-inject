# go-inject

## 引入
```
go get -u github.com/xiongwdd/go-inject
```


## 基本使用

一般配合文本框架使用,或者对map的键值对进行解析注入

自定义的接收

```GO

type student struct {
    Name        string      `inject:""`                         //表示注入参数
    Age         int         ``                                  //不填写则不会对参数进行注入
    Addr        string      `inject:"NotNull"`                  // 表示注入该参数,并且这个参数不能为空,如果为空会返回一个err
    School      string      `inject:"NotNull,def=南山外国语"`    // def表示如果该值必填,但是又没有传值的话采用的默认值
    CreateTime  time.Time   `date="2006-01-02 15:04:05"`        // 可以根据传入的时间字符串转换为时间格式
    Tel         string      `regular:"^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\\d{8}$"` //表示可以对字段正则校验
}

```

参数目前支持`url.values`格式或者 `map[string][]string`
