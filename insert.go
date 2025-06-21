package mysqlPool

import (
	"fmt"
	"strings"
)

func (q *builder) Insert(data map[string]interface{}) (int64, error) {
	if q.table == nil {
		return 0, q.logger.Error(nil, "Table is required")
	}

	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	placeholders := make([]string, 0, len(data))

	for column, value := range data {
		columns = append(columns, fmt.Sprintf("`%s`", column))
		values = append(values, value)
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)",
		*q.table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	result, err := q.exec(query, values...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
