package docopt_unmarshal

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Hook func(reflect.Value, string)error

type Unmarshaller struct {
	hook map[string]Hook
}

func New() *Unmarshaller {
	um := &Unmarshaller{hook: map[string]Hook{}}
	return um
}

var DefaultUnmarshaller *Unmarshaller = New()

func DocoptUnmarshal(arguments map[string]interface{}, options interface{}) error {
	return DefaultUnmarshaller.Unmarshal(arguments, options)
}

func (um *Unmarshaller) AddHook(type_string string, new_hook Hook) {
	um.hook[type_string] = new_hook
}

func (um *Unmarshaller) AddHooks(new_hooks map[string]Hook) {
	for k,v := range(new_hooks) {
		um.AddHook(k,v)
	}
}

func (um *Unmarshaller) Unmarshal(arguments map[string]interface{}, options interface{}) error {
	var seen []string
	seen, err := um.docopt_unmarshal(arguments, options, seen)
	if err != nil { return err }
	for _, a := range seen {
		delete(arguments, a)
	}
	for leftover, _ := range arguments {
		return errors.New(fmt.Sprintf("%s is missing from options struct", leftover))
	}
	return nil
}

func (um *Unmarshaller) docopt_unmarshal(arguments map[string]interface{}, options interface{}, seen []string) ([]string, error) {
	val := reflect.ValueOf(options).Elem()
	typ := val.Type()
	for i:=0; i<val.NumField(); i++ {
		f_val := val.Field(i)
		f_typ := typ.Field(i)
		flag := f_typ.Tag.Get("docopt")
		if flag != "" {
			a, exist := arguments[flag]
			if !exist {
				return seen, errors.New(fmt.Sprintf("Struct member %s has no corresponding option %s in docopt", f_typ.Name, flag))
			} else if a != nil {
				a_typ := reflect.TypeOf(a)
				if a_typ.Kind() == reflect.String {
				    if hook, exists := um.hook[f_typ.Type.String()]; exists {
						if err := hook(f_val, a.(string)); err != nil {
							return seen, errors.New(fmt.Sprintf("%s: %s", flag, err))
						}
				    } else {
					switch f_typ.Type.Kind() {
					case reflect.Bool:
						f_val.SetBool(a != nil)
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						if f_typ.Type.String() == "time.Duration" {
							dv, err := time.ParseDuration(a.(string))
							if err != nil {
								return seen, errors.New(fmt.Sprintf("%s: %s", flag, err))
							}
							f_val.Set(reflect.ValueOf(dv))
						} else {
							iv, err := strconv.ParseInt(a.(string), 10, 64)
							if err != nil {
								return seen, errors.New(fmt.Sprintf("%s: %s", flag, err))
							}
							f_val.SetInt(iv)
						}
					case reflect.Float32, reflect.Float64:
						fv, err := strconv.ParseFloat(a.(string), 64)
						if err != nil {
							return seen, errors.New(fmt.Sprintf("%s: %s", flag, err))
						}
						f_val.SetFloat(fv)
					default:
						f_val.Set(reflect.ValueOf(a))
					}
				    }
				} else {
					f_val.Set(reflect.ValueOf(a))
				}
			}
			seen = append(seen, flag)
		}
		if f_val.Type().Kind() == reflect.Struct {
			var err error
			if seen, err = um.docopt_unmarshal(arguments, f_val.Addr().Interface(), seen); err != nil {
				return seen, err
			}
		}
	}
	return seen, nil
}
