package tupletype

import (
	"fmt"
	"strings"
)

func tupleformysql(data []string) {
	temp := make([]string, 0)
	for _, e := range data {
		r := fmt.Sprintf("`%s`", e)
		temp = append(temp, r)

	}
	fmt.Println(strings.Join(temp, ","))

}

func listtostring(data []string) string {

	r := strings.Join(data, ",")
	return r

}

type Tuple struct {
	Data []string
}

func ToTuple(data []string) {
	tupleformysql(data)

}

func FetchTuple(data []string) string {
	return listtostring(data)

}
