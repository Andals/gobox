package mysql

import (
	//     "fmt"
	"strings"
)

const (
	COND_EQUAL         = "="
	COND_NOT_EQUAL     = "!="
	COND_LESS          = "<"
	COND_LESS_EQUAL    = "<="
	COND_GREATER       = ">"
	COND_GREATER_EQUAL = ">="
	COND_IN            = "in"
	COND_NOT_IN        = "not in"
	COND_LIKE          = "like"
	COND_BETWEEN       = "between"
)

type ColQueryItem struct {
	Name      string
	Condition string
	Values    []interface{}
}

func NewColQueryItem(name, condition string, values ...interface{}) *ColQueryItem {
	return &ColQueryItem{
		Name:      name,
		Condition: condition,
		Values:    values,
	}
}

type SimpleQueryBuilder struct {
	query string
	args  []interface{}
}

func (this *SimpleQueryBuilder) Query() string {
	return this.query
}

func (this *SimpleQueryBuilder) Args() []interface{} {
	return this.args
}

func (this *SimpleQueryBuilder) Insert(tableName string, colNames ...string) *SimpleQueryBuilder {
	this.args = nil

	this.query = "INSERT INTO " + tableName + " ("
	this.query += strings.Join(colNames, ", ") + ")"

	return this
}

func (this *SimpleQueryBuilder) Values(colsValues ...[]interface{}) *SimpleQueryBuilder {
	l := len(colsValues) - 1
	if l == -1 {
		return this
	}

	this.query += " VALUES "
	for i := 0; i < l; i++ {
		this.buildColValues(colsValues[i])
		this.query += ", "
	}
	this.buildColValues(colsValues[l])

	return this
}

func (this *SimpleQueryBuilder) Delete(tableName string) *SimpleQueryBuilder {
	this.args = nil

	this.query = "DELETE FROM " + tableName

	return this
}

func (this *SimpleQueryBuilder) Update(tableName string) *SimpleQueryBuilder {
	this.args = nil

	this.query = "UPDATE " + tableName

	return this
}

func (this *SimpleQueryBuilder) Set(setItems ...*ColQueryItem) *SimpleQueryBuilder {
	l := len(setItems) - 1
	if l == -1 {
		return this
	}

	this.query += " SET "
	for i := 0; i < l; i++ {
		this.query += setItems[i].Name + " = ?, "
		this.args = append(this.args, setItems[i].Values...)
	}
	this.query += setItems[l].Name + " = ? "
	this.args = append(this.args, setItems[l].Values...)

	return this
}

func (this *SimpleQueryBuilder) Select(what, tableName string) *SimpleQueryBuilder {
	this.args = nil

	this.query = "SELECT " + what + " FROM " + tableName

	return this
}

func (this *SimpleQueryBuilder) WhereConditionAnd(condItems ...*ColQueryItem) *SimpleQueryBuilder {
	this.query += " WHERE "

	this.buildWhereCondition("AND", condItems...)

	return this
}

func (this *SimpleQueryBuilder) WhereConditionOr(condItems ...*ColQueryItem) *SimpleQueryBuilder {
	this.query += " WHERE "

	this.buildWhereCondition("OR", condItems...)

	return this
}

func (this *SimpleQueryBuilder) OrderBy(orderBy string) *SimpleQueryBuilder {
	this.query += " ORDER BY " + orderBy

	return this
}

func (this *SimpleQueryBuilder) GroupBy(groupBy string) *SimpleQueryBuilder {
	this.query += " GROUP BY " + groupBy

	return this
}

func (this *SimpleQueryBuilder) HavingConditionAnd(condItems ...*ColQueryItem) *SimpleQueryBuilder {
	this.query += " HAVING "

	this.buildWhereCondition("AND", condItems...)

	return this
}

func (this *SimpleQueryBuilder) HavingConditionOr(condItems ...*ColQueryItem) *SimpleQueryBuilder {
	this.query += " HAVING "

	this.buildWhereCondition("OR", condItems...)

	return this
}

func (this *SimpleQueryBuilder) Limit(offset, rowCnt int) *SimpleQueryBuilder {
	this.query += " LIMIT ?, ?"
	this.args = append(this.args, offset, rowCnt)

	return this
}

func (this *SimpleQueryBuilder) buildColValues(colValues []interface{}) {
	l := len(colValues) - 1
	if l == -1 {
		return
	}

	this.query += "("

	for i := 0; i < l; i++ {
		this.query += "?, "
		this.args = append(this.args, colValues[i])
	}

	this.query += "?)"
	this.args = append(this.args, colValues[l])
}

func (this *SimpleQueryBuilder) buildWhereCondition(andOr string, condItems ...*ColQueryItem) {
	l := len(condItems) - 1
	if l == -1 {
		return
	}

	for i := 0; i < l; i++ {
		this.buildCondition(condItems[i])
		this.query += " " + andOr + " "
	}
	this.buildCondition(condItems[l])
}

func (this *SimpleQueryBuilder) buildCondition(condItem *ColQueryItem) {
	switch condItem.Condition {
	case COND_EQUAL, COND_NOT_EQUAL, COND_LESS, COND_LESS_EQUAL, COND_GREATER, COND_GREATER_EQUAL:
		this.query += condItem.Name + " " + condItem.Condition + " ?"
		this.args = append(this.args, condItem.Values...)
	case COND_IN:
		this.buildConditionInOrNotIn(condItem, "IN")
	case COND_NOT_IN:
		this.buildConditionInOrNotIn(condItem, "NOT IN")
	case COND_LIKE:
		this.query += condItem.Name + " LIKE ?"
		this.args = append(this.args, condItem.Values...)
	case COND_BETWEEN:
		if len(condItem.Values) != 2 {
			return
		}

		this.query += condItem.Name + " BETWEEN ? AND ?"
		this.args = append(this.args, condItem.Values...)
	}
}

func (this *SimpleQueryBuilder) buildConditionInOrNotIn(condItem *ColQueryItem, inOrNotIn string) {
	l := len(condItem.Values) - 1
	if l == -1 {
		return
	}

	this.query += condItem.Name + " " + inOrNotIn + " ("
	for i := 0; i < l; i++ {
		this.query += "?, "
	}
	this.query += "?)"
	this.args = append(this.args, condItem.Values...)
}
