package mysqlPool

import (
	"database/sql"
	"fmt"
	"strings"
)

func (q *builder) Get() (*sql.Rows, error) {
	if q.table == nil {
		return nil, q.logger.Error(nil, "Table is required")
	}

	fieldNames := make([]string, len(q.selectList))
	for i, field := range q.selectList {
		switch {
		case field == "*":
			fieldNames[i] = "*"
		case strings.ContainsAny(field, ".()"):
			fieldNames[i] = field
		default:
			fieldNames[i] = fmt.Sprintf("`%s`", field)
		}
	}

	query := fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(fieldNames, ", "), *q.table)

	if len(q.joinList) > 0 {
		query += " " + strings.Join(q.joinList, " ")
	}

	if len(q.whereList) > 0 {
		query += " WHERE " + strings.Join(q.whereList, " AND ")
	}

	if q.withTotal {
		query = fmt.Sprintf("SELECT COUNT(*) OVER() AS total, data.* FROM (%s) AS data", query)
	}

	if len(q.orderList) > 0 {
		query += " ORDER BY " + strings.Join(q.orderList, ", ")
	}

	if q.limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *q.limit)
	}

	if q.offset != nil {
		query += fmt.Sprintf(" OFFSET %d", *q.offset)
	}

	return q.query(query, q.bindingList...)
}
