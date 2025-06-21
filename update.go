package mysqlPool

import (
	"database/sql"
	"fmt"
	"strings"
)

func (q *builder) Update(data ...map[string]interface{}) (sql.Result, error) {
	if q.table == nil {
		return nil, q.logger.Error(nil, "Table is required")
	}

	values := []interface{}{}

	if len(data) > 0 {
		for column, value := range data[0] {
			columnName := column
			if !strings.Contains(column, ".") {
				columnName = fmt.Sprintf("`%s`", column)
			}

			if str, ok := value.(string); ok && contains(supportFunction, strings.ToUpper(str)) {
				q.setList = append(q.setList, fmt.Sprintf("%s = %s", columnName, str))
			} else {
				q.setList = append(q.setList, fmt.Sprintf("%s = ?", columnName))
				values = append(values, value)
			}
		}
	}

	query := fmt.Sprintf("UPDATE `%s` SET %s", *q.table, strings.Join(q.setList, ", "))

	if len(q.whereList) > 0 {
		query += " WHERE " + strings.Join(q.whereList, " AND ")
	}

	allValues := append(values, q.bindingList...)
	return q.exec(query, allValues...)
}
