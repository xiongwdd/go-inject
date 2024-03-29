package inject_test

import (
	"fmt"
	"github.com/xiongwdd/go-inject"
	"reflect"
	"testing"
	"time"
)

func Test_time(t *testing.T) {
	tt := time.Now()
	v := reflect.ValueOf(tt)
	to := reflect.TypeOf(tt)
	println(v.Type())
	println(v.Kind())
	println(to.Kind())
	println(to.String())
}

func Test_tag(t *testing.T) {
	type demo struct {
		Name string `inject:"1231231"`
	}
	var d demo
	d.Name = "zhangsan"
	tt := reflect.TypeOf(d)
	aa, _ := tt.FieldByName("Name")
	println(aa.Type.Name())
	println(aa.Name)
	println(aa.PkgPath)
	println(aa.Tag)
}

type FriendSave struct {
	OrgUid     int       `json:"orgUid" inject:"NotNull"`
	Uid        int       `json:"uid"  inject:"NotNull"`
	LUid       int64     `json:"type" inject:""`
	LoUid      int64     `json:"type" inject:"NotNull"`
	Context    string    `json:"state" inject:"NotNull,MaxLen=10"`
	CreateTime time.Time `json:"remark" inject:"" date:"2006-01-02 15:04:05"`
	Tel        string    `json:"method" inject:"" regular:"^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\\d{8}$"`
	Name       string    `json:"name" inject:""`
	Addr       string    `inject:""`
}

var arr = []map[string][]string{{
	"orguid":     {"10001"},
	"uid":        {"20001"},
	"luid":       {"30001"},
	"louid":      {"40001"},
	"context":    {"小于5"},
	"createtime": {"2019-09-26 18:03:01"},
	"tel":        {"18571702046"},
	"name":       {"张一"},
	"addr":       {"哈哈jtyjfdjh"},
}, {
	"orguid":     {"10002"},
	"uid":        {"20002"},
	"LUid":       {"30002"},
	"LoUid":      {"40002"},
	"Context":    {"一段超市那"},
	"CreateTime": {"2019-09-26 18:03:02"},
	"Tel":        {"18571702047"},
	"Name":       {"张二"},
	"Addr":       {"哈哈哈hhh女为人家估计是个"},
}, {
	"orguid":     {"10003"},
	"uid":        {"20003"},
	"LUid":       {"30003"},
	"LoUid":      {"40003"},
	"Context":    {"nihaoma"},
	"CreateTime": {"2019-09-26 18:03:03"},
	"Tel":        {"18571702222"},
	"Name":       {"张三"},
}, {
	"orguid":     {"10004"},
	"uid":        {"20004"},
	"LUid":       {"30004"},
	"LoUid":      {"40004"},
	"Context":    {"小于34"},
	"CreateTime": {"2019-09-26 18:03:04"},
	"Tel":        {"18571702046"},
	"Name":       {"张五"},
	"Addr":       {"哈adshhh"},
}, {
	"orguid":     {"10005"},
	"uid":        {"20005"},
	"LUid":       {"30005"},
	"LoUid":      {"40005"},
	"Context":    {"小于asdas"},
	"CreateTime": {"2019-09-26 18:03:05"},
	"Tel":        {"18571702046"},
	"Name":       {"张四"},
	"Addr":       {"哈哈哈hhh"},
},
}

func Test_inject(t *testing.T) {
	var flist []FriendSave
	for _, a := range arr {
		var f = FriendSave{}
		err := inject.InjectionCheck(a, &f)
		if err != nil {
			fmt.Println(err)
		} else {
			flist = append(flist, f)
		}
	}
	fmt.Println(flist)
}
