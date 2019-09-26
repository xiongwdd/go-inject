package inject

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	Binary      = 2
	Octal       = 8
	Decimal     = 10
	Hexadecimal = 16
)

const (
	Int8  = 8
	Int16 = 16
	Int32 = 32
	Int64 = 64
)

const (
	// 表示需要注入的值
	Inject  = "inject"
	Regular = "regular"
	Date    = "date"
	NotNull = "NotNull"
	MaxLen  = "MaxLen"
	Def     = "def"
)

const (
	Space     = " "
	comma     = ","
	Slash     = "/"
	Backslash = `\`
	equal     = "="
)

func InjectionCheck(param url.Values, object interface{}) (err error) {

	if len(param) == 0 || object == nil {
		return errors.New(fmt.Sprintf("The Parameter is empty"))
	}

	var (
		name      string
		values    []string
		value     string
		ok        bool
		tag       string
		regu      string
		d         string
		save      interface{}
		t         = reflect.TypeOf(object)
		v         = reflect.ValueOf(object)
		p         = getFieldName(object)
		maxlen    int
		isNotNull = false
		isMaxLen  = false
	)

	for _, name = range p {
		tt, _ := t.Elem().FieldByName(name)
		// 用户自定义的标签
		tag, ok = tt.Tag.Lookup(Inject)

		if len(tag) == 0 && !ok {
			continue
		}

		vt := v.Elem().FieldByName(name)
		if strings.Contains(tag, Def) {
			// 如果可以有默认值
			if v, o := param[name]; !o || len(v) == 0 || len(v[0]) == 0 {
				param[name] = []string{getDef(tag)}
			}
		}

		isNotNull = strings.Contains(tag, NotNull)
		isMaxLen = strings.Contains(tag, MaxLen)

		if values, ok = param[name]; ok {
			if isNotNull {
				if len(values) == 0 || len(values[0]) == 0 {
					return errors.New("bad request")
				}
			}

			value = values[0]
			// 获取到值,其他类型的
			save, err = getFieldValue(tt.Type.Kind(), value)
			if err != nil {
				if isNotNull {
					return errors.New("param error")
				} else {
					//  如果没有标明是非空的,则如果解析失败也不影响程序继续执行
					continue
				}
			}

			if vt.CanSet() {
				// 长度校验
				if isMaxLen {
					arr := strings.Split(tag, comma)
					for _, l := range arr {
						if strings.Contains(l, MaxLen) {
							maxlen, _ = strconv.Atoi(strings.Split(l, equal)[1])
							if s, o := save.(string); !o || utf8.RuneCountInString(s) > maxlen {
								return errors.New("String length is too large ")
							}
							break
						}
					}
				}

				// 正则校验
				if regu, ok = tt.Tag.Lookup(Regular); ok && len(regu) > 0 {
					if tt.Type.Kind() == reflect.String {
						if o, err := regexp.MatchString(regu, save.(string)); !o || err != nil {
							return errors.New("Regular match failed ")
						}
					}
				}

				// 时间格式校验
				if d, ok = tt.Tag.Lookup(Date); ok && strings.Contains(tt.Type.Name(), "Time") {
					if s, o := save.(string); len(d) > 0 && o {
						save, err = time.Parse(strings.Trim(d, Space), s)
						if err != nil {
							return errors.New("Time format error ")
						}
					} else {
						// 不能为空,但是目标格式为时间格式,无法转换为未知的时间格式
						return errors.New("Time format error ")
					}
				}

				vt.Set(reflect.ValueOf(save))
			} else {
				if isNotNull {
					return errors.New("param error")
				}
			}
		} else {
			if isNotNull {
				return errors.New("bad request")
			}

		}
	}
	return
}

func getDef(s string) string {
	arr := strings.Split(s, comma)
	for _, l := range arr {
		if strings.Contains(l, Def) {
			return strings.Split(l, equal)[1]
		}
	}
	return ""
}

func getFieldValue(kind reflect.Kind, value string) (save interface{}, err error) {
	switch kind {
	case reflect.String:
		save = value
	case reflect.Bool:
		save, err = strconv.ParseBool(value)
	case reflect.Int:
		save, err = strconv.Atoi(value)
	case reflect.Int8:
		save, err = strconv.ParseInt(value, Decimal, Int8)
	case reflect.Int16:
		save, err = strconv.ParseInt(value, Decimal, Int16)
	case reflect.Int32:
		save, err = strconv.ParseInt(value, Decimal, Int32)
	case reflect.Int64:
		save, err = strconv.ParseInt(value, Decimal, Int64)
	case reflect.Float32:
		save, err = strconv.ParseFloat(value, 32)
	case reflect.Float64:
		save, err = strconv.ParseFloat(value, 64)
	case reflect.Struct:
		save = value
	default:
		save = value
	}
	return
}

// 获取结构体中字段的名称
func getFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, strings.ToLower(t.Field(i).Name))
	}
	return result
}
