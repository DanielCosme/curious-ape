package sqlite

import (
	"fmt"
	"strings"
	"time"

	"github.com/danielcosme/curious-ape/internal/core/entity"
	_ "github.com/mattn/go-sqlite3"
)

const DriverName = "sqlite3"

type sqlQueryBuilder struct {
	TableName string
	Data      []filterData
}

type filterData struct {
	columnName string
	values     []interface{}
}

func newBuilder(tableName string) *sqlQueryBuilder {
	return &sqlQueryBuilder{TableName: tableName, Data: []filterData{}}
}

func (qb *sqlQueryBuilder) AddFilter(columnName string, values []interface{}) {
	qb.Data = append(qb.Data, filterData{columnName, values})
}

func (qb *sqlQueryBuilder) generate() (string, []interface{}) {
	var args []interface{}
	q := fmt.Sprintf("SELECT * FROM %s ", qb.TableName)

	if len(qb.Data) > 0 {
		q += "WHERE"

		for idx, data := range qb.Data {
			if idx > 0 {
				q = fmt.Sprintf("%s AND", q)
			}
			q = fmt.Sprintf("%s %s IN (%s)", q, data.columnName, getArgs(len(data.values)))
			args = append(args, data.values...)
		}
	}

	return q, args
}

func getArgs(length int) string {
	ss := []string{}
	for i := 0; i < length; i++ {
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
		iSlice[i] = sanitizeDate(v)
	}
	return iSlice
}

func habitTypeToInterface(hts []entity.HabitType) []interface{} {
	iSlice := make([]interface{}, len(hts))
	for i, v := range hts {
		iSlice[i] = v
	}
	return iSlice
}
