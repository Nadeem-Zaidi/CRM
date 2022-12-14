package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/Nadeem-Zaidi/CRM/errorh"
	"github.com/Nadeem-Zaidi/CRM/pexception"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var DB *sql.DB
var databasename string

func InitDB(drivername, username, password, dbname string) {
	pexception.RecoverFromPanic()
	q := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", username, password, dbname)
	db, err := sql.Open(drivername, q)
	if err != nil {
		panic(err.Error())
	}
	DB = db
	databasename = dbname
}

// func SingleRow(withid int, columnname []string) {
// 	pexception.RecoverFromPanic()
// 	if len(columnname) == 1 {
// 		query := fmt.Sprintf("SELECT %s from %s WHERE id=%d", columnname, databasename, withid)
// 		print(query)
// 	} else {
// 		s := tupletype.FetchTuple(columnname)
// 		query := fmt.Sprintf("SELECT %s from %s where id=%d", s, databasename, withid)
// 		fmt.Println(query)
// 	}

// }

// func (a *DataBaseApp) All() {
// 	var p schema.Product
// 	query := `select * from products`
// 	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancelFunc()
// 	stmt, err := a.DB.PrepareContext(ctx, query)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer stmt.Close()
// 	row := stmt.QueryRowContext(ctx)
// 	if err := row.Scan(&p.ID, &p.Name, &p.Price); err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(p.Name, p.Price)
// }

type Number interface {
	int64 | int | float64 | float32
}

func Create(data interface{}) {

	typeData := map[string]interface{}{
		"string":  "varchar(30)",
		"int":     "int",
		"float64": "decimal(20,0)",
	}
	fieldMap := make(map[string]string)
	t := reflect.ValueOf(data)
	for i := 0; i < t.Type().NumField(); i++ {
		// fieldMap[t.Type().Field(i).Name] = t.Type().Field(i).Type.String()
		fieldMap[t.Type().Field(i).Tag.Get("json")] = fmt.Sprintf("%s %s", typeData[t.Type().Field(i).Type.String()], t.Type().Field(i).Tag.Get("gm"))

	}
	//rr := fmt.Sprintf("\"%s\"", fieldMap["Address"])
	result := make([]string, 0)

	for i, e := range fieldMap {
		r := fmt.Sprintf("%s %s", i, e)
		result = append(result, r)

	}
	tablename := strings.Split(reflect.TypeOf(data).String(), ".")
	ct := strings.Join(result, ",")
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(%s);", tablename[1], ct)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	res, err := DB.ExecContext(ctx, query)
	if err != nil {
		fmt.Println(err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("Rows affected when creating table: %d", rows)

}

func Insert(data interface{}) {
	d := reflect.ValueOf(data)
	propertyNameList := make([]string, 0)
	valueNameList := make([]string, 0)

	for i := 0; i < d.Type().NumField(); i++ {
		property := d.Type().Field(i).Name
		p := fmt.Sprintf("`%s`", property)
		propertyNameList = append(propertyNameList, p)
		propertyType := d.Type().Field(i).Type
		values := d.Field(i)
		fmt.Println(propertyType)
		if propertyType.Kind() == reflect.Int || propertyType.Kind() == reflect.Float64 {
			r := fmt.Sprintf("%v", values)
			valueNameList = append(valueNameList, r)

		} else {
			r := fmt.Sprintf("\"%s\"", values)
			valueNameList = append(valueNameList, r)

		}

	}
	tablename := strings.Split(reflect.TypeOf(data).String(), ".")
	propertyNameString := strings.Join(propertyNameList, ",")
	valueNameString := strings.Join(valueNameList, ",")

	query := fmt.Sprintf("INSERT INTO %s(%s) VALUES (%s);", tablename[1], propertyNameString, valueNameString)
	fmt.Println(query)

}
func b2s(bs []uint8) string {
	ba := make([]byte, 0, len(bs))
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

func StructScan(rows *sql.Rows, data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination") // @todo add new error message

	}
	v = reflect.Indirect(v)
	t := v.Type()
	cols, _ := rows.Columns()
	var m map[string]interface{}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnsPointer := make([]interface{}, len(cols))
		for i := range columns {
			columnsPointer[i] = &columns[i]
		}
		if err := rows.Scan(columnsPointer...); err != nil {
			return err
		}
		m = make(map[string]interface{})
		// fmt.Println(*columnsPointer[2].(*interface{}))
		for i, colName := range cols {
			val := columnsPointer[i].(*interface{})
			m[colName] = *val
		}
		list := make([]interface{}, 0)

		for i := 0; i < v.NumField(); i++ {
			field := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]

			if item, ok := m[field]; ok {
				if v.Field(i).CanSet() {
					if item != nil {
						switch v.Field(i).Kind() {
						case reflect.String:
							v.Field(i).SetString(b2s(item.([]uint8)))
						case reflect.Float32, reflect.Float64:
							v.Field(i).SetFloat(item.(float64))
						case reflect.Ptr:
							if reflect.ValueOf(item).Kind() == reflect.Bool {
								itemBool := item.(bool)
								v.Field(i).Set(reflect.ValueOf(&itemBool))
							}
						case reflect.Struct:
							v.Field(i).Set(reflect.ValueOf(item))
						default:
							fmt.Println(t.Field(i).Name, ": ", v.Field(i).Kind(), " - > - ", reflect.ValueOf(item).Kind()) // @todo remove after test out the Get methods
						}
					}
				}
			}

		}
		fmt.Println("below")
		fmt.Println(list)

	}

	return nil
}

// func ScanAgainv(v interface{}) {
// 	vType := reflect.TypeOf(v)
// 	if k := vType.Kind(); k != reflect.Ptr {
// 		fmt.Println("must be a pointer ")
// 	}
// 	vType = vType.Elem()
// 	vVal := reflect.ValueOf(v).Elem()
// 	if vType.Kind() == reflect.Slice {
// 		fmt.Println("error for slice")
// 	}
// 	sl := reflect.New(reflect.SliceOf(vType))
// 	rows(sl.Interface())

// 	sl = sl.Elem()
// 	if sl.Len() == 0 {
// 		fmt.Println("0 rows")
// 	}
// 	vVal.Set(sl.Index(0))
// }

// func rows(v interface{}, r *sql.Rows) {
// 	vType := reflect.TypeOf(v)
// 	if k := vType.Kind(); k != reflect.Ptr {
// 		fmt.Println("must be a pointer ")
// 	}
// 	sliceType := vType.Elem()
// 	if reflect.Slice != sliceType.Kind() {
// 		fmt.Println(",must be a slice")
// 	}
// 	sliceVal := reflect.Indirect(reflect.ValueOf(v))
// 	itemType := sliceType.Elem()
// 	cols, err := r.Columns()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	isPrimitive:=itemType.Kind()!=reflect.Struct

// 	for r.Next(){
// 		sliceItem:=reflect.New(itemType).Elem()
// 		 var pointers []interface{}
// 		 if primitive{
// 			if len(cols)>1{

// 			}
// 		 }
// 	}

// }

var (
	// ErrTooManyColumns indicates that a select query returned multiple columns and
	// attempted to bind to a slice of a primitive type. For example, trying to bind
	// `select col1, col2 from mutable` to []string
	ErrTooManyColumns = errors.New("too many columns returned for primitive slice")

	// ErrSliceForRow occurs when trying to use Row on a slice
	ErrSliceForRow = errors.New("cannot scan Row into slice")

	// AutoClose is true when scan should automatically close Scanner when the scan
	// is complete. If you set it to false, then you must defer rows.Close() manually
	AutoClose = true

	// OnAutoCloseError can be used to log errors which are returned from rows.Close()
	// By default this is a NOOP function
	OnAutoCloseError = func(error) {}
)

// Row scans a single row into a single variable. It requires that you use
// db.Query and not db.QueryRow, because QueryRow does not return column names.
// There is no performance impact in using one over the other. QueryRow only
// defers returning err until Scan is called, which is an unnecessary
// optimization for this library.
func Row(v interface{}, r *sql.Rows) error {
	return row(v, r, false)
}

// RowStrict scans a single row into a single variable. It is identical to
// Row, but it ignores fields that do not have a db tag
func RowStrict(v interface{}, r *sql.Rows) error {
	return row(v, r, true)
}

func row(v interface{}, r *sql.Rows, strict bool) error {
	vType := reflect.TypeOf(v)
	if k := vType.Kind(); k != reflect.Ptr {
		return fmt.Errorf("%q must be a pointer: %w", k.String())
	}

	vType = vType.Elem()
	vVal := reflect.ValueOf(v).Elem()
	if vType.Kind() == reflect.Slice {
		return ErrSliceForRow
	}

	sl := reflect.New(reflect.SliceOf(vType))
	err := rows(sl.Interface(), r, strict)
	if err != nil {
		return err
	}

	sl = sl.Elem()

	if sl.Len() == 0 {
		return sql.ErrNoRows
	}

	vVal.Set(sl.Index(0))

	return nil
}

// Rows scans sql rows into a slice (v)
func Rows(v interface{}, r *sql.Rows) (outerr error) {
	return rows(v, r, false)
}

// RowsStrict scans sql rows into a slice (v) only using db tags
func RowsStrict(v interface{}, r *sql.Rows) (outerr error) {
	return rows(v, r, true)
}

func rows(v interface{}, r *sql.Rows, strict bool) (outerr error) {
	if AutoClose {
		defer closeRows(r)
	}

	vType := reflect.TypeOf(v)
	if k := vType.Kind(); k != reflect.Ptr {
		return fmt.Errorf("%q must be a pointer: %w", k.String())
	}
	sliceType := vType.Elem()
	if reflect.Slice != sliceType.Kind() {
		return fmt.Errorf("%q must be a slice: %w", sliceType.String())
	}

	sliceVal := reflect.Indirect(reflect.ValueOf(v))
	itemType := sliceType.Elem()

	cols, err := r.Columns()
	fmt.Println("cols")
	fmt.Println(cols)
	if err != nil {
		return err
	}

	// isPrimitive := itemType.Kind() != reflect.Struct

	for r.Next() {
		sliceItem := reflect.New(itemType).Elem()

		var pointers []interface{}
		// if isPrimitive {
		// 	if len(cols) > 1 {
		// 		return ErrTooManyColumns
		// 	}
		// 	pointers = []interface{}{sliceItem.Addr().Interface()}
		// } else {
		// 	pointers = structPointers(sliceItem, cols, strict)
		// }
		pointers = structPointers(sliceItem, cols)

		if len(pointers) == 0 {
			return nil
		}

		err := r.Scan(pointers...)
		if err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, sliceItem))
	}
	return r.Err()
}

// Initialization the tags from struct.
func initFieldTag(sliceItem reflect.Value, fieldTagMap *map[string]reflect.Value) {
	typ := sliceItem.Type()
	for i := 0; i < sliceItem.NumField(); i++ {
		if typ.Field(i).Anonymous || typ.Field(i).Type.Kind() == reflect.Struct {
			// found an embedded struct
			sliceItemOfAnonymous := sliceItem.Field(i)
			initFieldTag(sliceItemOfAnonymous, fieldTagMap)
		}
		tag, ok := typ.Field(i).Tag.Lookup("db")
		if ok && tag != "" {
			(*fieldTagMap)[tag] = sliceItem.Field(i)
		}
	}
	fmt.Println("from init tag")
	fmt.Println(fieldTagMap)
}

func structPointers(sliceItem reflect.Value, cols []string) []interface{} {
	pointers := make([]interface{}, 0, len(cols))
	// fieldTag := make(map[string]reflect.Value, len(cols))
	// initFieldTag(sliceItem, &fieldTag)

	for _, colName := range cols {
		var fieldVal reflect.Value
		// if v, ok := fieldTag[colName]; ok {
		// 	fieldVal = v
		// } else {
		// 	if strict {
		// 		fieldVal = reflect.ValueOf(nil)
		// 	} else {
		// 		fieldVal = sliceItem.FieldByName(cases.Title(language.English).String(colName))
		// 	}
		// }
		fieldVal = sliceItem.FieldByName(cases.Title(language.English).String(colName))
		if !fieldVal.IsValid() || !fieldVal.CanSet() {
			// have to add if we found a column because Scan() requires
			// len(cols) arguments or it will error. This way we can scan to
			// a useless pointer
			var nothing interface{}
			pointers = append(pointers, &nothing)
			continue
		}

		pointers = append(pointers, fieldVal.Addr().Interface())
	}
	return pointers
}

func closeRows(c io.Closer) {
	if err := c.Close(); err != nil {
		if OnAutoCloseError != nil {
			OnAutoCloseError(err)
		}
	}
}

func FindOne(data interface{}, r *sql.Rows) error {

	dType := reflect.TypeOf(data)

	fmt.Println(dType.Kind())

	if dType.Kind() != reflect.Ptr {
		return &errorh.RaiseError{ErrorMessage: "must be a pointer", ErrorCode: 3}

	}

	dValue := reflect.ValueOf(data).Elem()
	if dValue.Kind() != reflect.Struct {
		return &errorh.RaiseError{ErrorMessage: "must be of type struct", ErrorCode: 3}

	}

	c, err := r.Columns()
	if err != nil {
		fmt.Println("inside error")
		fmt.Println(err)
	}

	pointers := make([]interface{}, 0, len(c))
	for _, colName := range c {

		fieldVal := dValue.FieldByName(cases.Title(language.English).String(colName))

		pointers = append(pointers, fieldVal.Addr().Interface())

	}

	for r.Next() {

		err := r.Scan(pointers...)

		if err != nil {
			fmt.Println("inside error")
			fmt.Println(err)
		}

	}
	return nil

}

func FindA(data interface{}, r *sql.Rows) error {
	dType := reflect.TypeOf(data)
	if dType.Kind() != reflect.Ptr {
		return &errorh.RaiseError{ErrorMessage: "must be a pointer", ErrorCode: 3}
	}
	sliceType := dType.Elem()

	if reflect.Slice != sliceType.Kind() {
		return &errorh.RaiseError{ErrorMessage: "must be a slice", ErrorCode: 4}
	}
	dv := reflect.ValueOf(data)
	sliceVal := reflect.Indirect(dv)
	itemType := dType.Elem().Elem()

	cols, err := r.Columns()
	fmt.Println("cols")
	fmt.Println(cols)
	if err != nil {
		fmt.Println(err)
	}

	for r.Next() {
		sliceItem := reflect.New(itemType).Elem()
		pointers := structPointers(sliceItem, cols)

		if len(pointers) == 0 {
			fmt.Println("can not process")
		}

		err := r.Scan(pointers...)
		if err != nil {
			fmt.Println(err)
		}
		sliceVal.Set(reflect.Append(sliceVal, sliceItem))

	}
	return nil

}
