package models

import (
	"database/sql"
	"fmt"
	"reflect"
)

func Findone(data interface{}, r *sql.Rows) {
	dType := reflect.TypeOf(data)
	fmt.Println(dType)
	dValue := reflect.ValueOf(data)
	fmt.Println(dValue)

}
