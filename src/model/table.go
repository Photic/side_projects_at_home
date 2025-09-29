package model

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type TableData struct {
	Headers []string   `db:"headers"`
	Rows    [][]string `db:"rows"`
}

func BuildTableData(data interface{}) TableData {
	table_data := TableData{}

	if data == nil {
		return table_data
	}

	rv := reflect.ValueOf(data)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return table_data
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		if rv.Len() == 0 {
			return table_data
		}
		first := rv.Index(0)
		for first.Kind() == reflect.Ptr {
			first = first.Elem()
		}
		switch first.Kind() {
		case reflect.Struct:
			t := first.Type()
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				tag := strings.Split(f.Tag.Get("json"), ",")[0]
				if tag == "" {
					tag = f.Name
				}
				if tag == "-" {
					continue
				}
				table_data.Headers = append(table_data.Headers, tag)
			}
			for i := 0; i < rv.Len(); i++ {
				elem := rv.Index(i)
				for elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				var row []string
				for j := 0; j < elem.NumField(); j++ {
					row = append(row, fmt.Sprint(elem.Field(j).Interface()))
				}
				table_data.Rows = append(table_data.Rows, row)
			}
			return table_data
		case reflect.Map:
			keysSet := map[string]struct{}{}
			var keys []string
			for i := 0; i < rv.Len(); i++ {
				elem := rv.Index(i)
				for elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				for _, k := range elem.MapKeys() {
					ks := fmt.Sprint(k.Interface())
					if _, ok := keysSet[ks]; !ok {
						keysSet[ks] = struct{}{}
						keys = append(keys, ks)
					}
				}
			}
			sort.Strings(keys)
			table_data.Headers = keys
			for i := 0; i < rv.Len(); i++ {
				elem := rv.Index(i)
				for elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				var row []string
				for _, k := range keys {
					v := elem.MapIndex(reflect.ValueOf(k))
					if !v.IsValid() {
						row = append(row, "")
					} else {
						row = append(row, fmt.Sprint(v.Interface()))
					}
				}
				table_data.Rows = append(table_data.Rows, row)
			}
			return table_data
		default:
			table_data.Headers = []string{"Value"}
			for i := 0; i < rv.Len(); i++ {
				table_data.Rows = append(table_data.Rows, []string{fmt.Sprint(rv.Index(i).Interface())})
			}
			return table_data
		}
	case reflect.Struct:
		t := rv.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			tag := strings.Split(f.Tag.Get("json"), ",")[0]
			if tag == "" {
				tag = f.Name
			}
			if tag == "-" {
				continue
			}
			table_data.Headers = append(table_data.Headers, tag)
		}
		var row []string
		for i := 0; i < rv.NumField(); i++ {
			row = append(row, fmt.Sprint(rv.Field(i).Interface()))
		}
		table_data.Rows = append(table_data.Rows, row)
		return table_data
	case reflect.Map:
		var keys []string
		for _, k := range rv.MapKeys() {
			keys = append(keys, fmt.Sprint(k.Interface()))
		}
		sort.Strings(keys)
		table_data.Headers = keys
		var row []string
		for _, k := range keys {
			row = append(row, fmt.Sprint(rv.MapIndex(reflect.ValueOf(k)).Interface()))
		}
		table_data.Rows = append(table_data.Rows, row)
		return table_data
	default:
		table_data.Headers = []string{"Value"}
		table_data.Rows = append(table_data.Rows, []string{fmt.Sprint(rv.Interface())})
		return table_data
	}
}
