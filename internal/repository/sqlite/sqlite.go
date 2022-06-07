package sqlite

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"time"
)

const DriverName = "sqlite3"

type sqlBuilder struct {
	TableName string
	Data      []filterData
}

type filterData struct {
	columnName string
	values     []interface{}
}

func newBuilder(tableName string) *sqlBuilder {
	return &sqlBuilder{TableName: tableName, Data: []filterData{}}
}

func (qb *sqlBuilder) AddFilter(columnName string, values []interface{}) {
	qb.Data = append(qb.Data, filterData{columnName, values})
}

func (qb *sqlBuilder) generate() (string, []interface{}) {
	var args []interface{}
	q := fmt.Sprintf("SELECT * FROM %s ", qb.TableName)

	if len(qb.Data) > 0 {
		q = q + "WHERE"

		for idx, data := range qb.Data {
			if idx > 0 {
				q = fmt.Sprintf("%s AND", q)
			}
			q = fmt.Sprintf("%s %s IN (%s)", q, data.columnName, getArgs(data.values))
			args = append(args, data.values...)
		}
	}

	return q, args
}

func getArgs(args []interface{}) string {
	ss := []string{}
	for i := 0; i < len(args); i++ {
		ss = append(ss, "?")
	}
	return strings.Join(ss, ",")
}

func intToInterface(ints []int) []interface{} {
	iSlice := make([]interface{}, len(ints))
	for i, v := range ints {
		iSlice[i] = v
	}
	return iSlice
}

func dateToInterface(ds []time.Time) []interface{} {
	iSlice := make([]interface{}, len(ds))
	for i, v := range ds {
		iSlice[i] = v
	}
	return iSlice
}
