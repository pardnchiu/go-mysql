package mysqlPool

import (
	"fmt"
	"log"
	"strings"
)

var (
	supportFunction = []string{
		"NOW()", "CURRENT_TIMESTAMP", "UUID()", "RAND()", "CURDATE()",
		"CURTIME()", "UNIX_TIMESTAMP()", "UTC_TIMESTAMP()", "SYSDATE()",
		"LOCALTIME()", "LOCALTIMESTAMP()", "PI()", "DATABASE()", "USER()",
		"VERSION()",
	}
)

func (db *Pool) DB(dbName string) *builder {
	_, err := db.db.Exec(fmt.Sprintf("USE `%s`", dbName))
	if err != nil {
		db.Logger.Error(err, "Failed to switch to database "+dbName)
	}

	return &builder{
		db:         db.db,
		dbName:     &dbName,
		selectList: []string{"*"},
		logger:     db.Logger,
	}
}

func (q *builder) Table(tableName string) *builder {
	q.table = &tableName
	return q
}

func (q *builder) Select(fields ...string) *builder {
	if len(fields) > 0 {
		q.selectList = fields
	}
	return q
}

func (q *builder) Total() *builder {
	q.withTotal = true
	return q
}

func (q *builder) InnerJoin(table, first, operator string, second ...string) *builder {
	return q.join("INNER", table, first, operator, second...)
}

func (q *builder) LeftJoin(table, first, operator string, second ...string) *builder {
	return q.join("LEFT", table, first, operator, second...)
}

func (q *builder) RightJoin(table, first, operator string, second ...string) *builder {
	return q.join("RIGHT", table, first, operator, second...)
}

// * private method
func (q *builder) join(joinType, table, first, operator string, second ...string) *builder {
	var secondField string
	if len(second) > 0 {
		secondField = second[0]
	} else {
		secondField = operator
		operator = "="
	}

	if !strings.Contains(first, ".") {
		first = fmt.Sprintf("`%s`", first)
	}
	if !strings.Contains(secondField, ".") {
		secondField = fmt.Sprintf("`%s`", secondField)
	}

	joinClause := fmt.Sprintf("%s JOIN `%s` ON %s %s %s", joinType, table, first, operator, secondField)
	q.joinList = append(q.joinList, joinClause)
	return q
}

func (q *builder) Where(column string, operator interface{}, value ...interface{}) *builder {
	var targetValue interface{}
	var targetOperator string

	if len(value) == 0 {
		targetValue = operator
		targetOperator = "="
	} else {
		targetOperator = fmt.Sprintf("%v", operator)
		targetValue = value[0]
	}

	if targetOperator == "LIKE" {
		if str, ok := targetValue.(string); ok {
			targetValue = fmt.Sprintf("%%%s%%", str)
		}
	}

	if !strings.Contains(column, "(") && !strings.Contains(column, ".") {
		column = fmt.Sprintf("`%s`", column)
	}

	placeholder := "?"
	if targetOperator == "IN" {
		placeholder = "(?)"
	}

	whereClause := fmt.Sprintf("%s %s %s", column, targetOperator, placeholder)
	q.whereList = append(q.whereList, whereClause)
	q.bindingList = append(q.bindingList, targetValue)

	return q
}

func (q *builder) OrderBy(column string, direction ...string) *builder {
	dir := "ASC"
	if len(direction) > 0 {
		dir = strings.ToUpper(direction[0])
	}

	if dir != "ASC" && dir != "DESC" {
		log.Printf("Invalid order direction: %s", dir)
		return q
	}

	if !strings.Contains(column, ".") {
		column = fmt.Sprintf("`%s`", column)
	}

	orderClause := fmt.Sprintf("%s %s", column, dir)
	q.orderList = append(q.orderList, orderClause)
	return q
}

func (q *builder) Limit(num int) *builder {
	q.limit = &num
	return q
}

func (q *builder) Offset(num int) *builder {
	q.offset = &num
	return q
}

func (q *builder) Increase(target string, number ...int) *builder {
	num := 1
	if len(number) > 0 {
		num = number[0]
	}

	setClause := fmt.Sprintf("%s = %s + %d", target, target, num)
	q.setList = append(q.setList, setClause)
	return q
}

// * private method
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
