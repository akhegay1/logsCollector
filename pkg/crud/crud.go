package crud

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func CreateInsert(q interface{}) string {
	log.Println("CreateInsert started")
	var fnames, fvalues string
	if reflect.ValueOf(q).Kind() == reflect.Struct {
		t := reflect.TypeOf(q).Name()
		nam := reflect.Indirect(reflect.ValueOf(q))

		v := reflect.ValueOf(q)
		for i := 0; i < v.NumField(); i++ {
			if reflect.TypeOf(q).Field(i).Tag.Get("keys") != "pk" {
				fnames = fnames + ", " + nam.Type().Field(i).Name
			}
		}
		fnames = "(" + fnames[2:] + ")"
		//log.Println("fnames", fnames)

		query := fmt.Sprintf("insert into monit_sch.%s %s values(", t, fnames)

		for i := 0; i < v.NumField(); i++ {
			switch v.Field(i).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if reflect.TypeOf(q).Field(i).Tag.Get("keys") != "pk" {
					fvalues = fmt.Sprintf("%s, %d", fvalues, v.Field(i).Int())
				}
			case reflect.String:
				if v.Field(i).String() != "" {
					fvalues = fmt.Sprintf("%s, '%s'", fvalues, strings.Replace(v.Field(i).String(), "'", "''", -1))
				} else {
					fvalues = fmt.Sprintf("%s, %s", fvalues, "null")
				}

			case reflect.Bool:
				fvalues = fmt.Sprintf("%s, %v", fvalues, v.Field(i).Interface().(bool))
			case reflect.Float32, reflect.Float64:
				fvalues = fmt.Sprintf("%s, %.2f", fvalues, v.Field(i).Interface().(float64))
			default:
				return "Unsupported type"
			}
		}

		fvalues = fvalues[2:]
		query = fmt.Sprintf("%s %s) returning id", query, fvalues)
		//fmt.Println(query)
		return query

	}
	return "unsupported type"
}
