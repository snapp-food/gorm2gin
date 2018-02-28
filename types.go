package gorm2gin

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
)

const queryParamTransactionUid = "_tr_uid"
const queryParamLimit = "_limit"
const queryParamOffset = "_offset"
const queryParamPageOrderField = "_page_order_field"
const queryParamPageOrderDirection = "_page_order_direction"

type CRUDerModelInterface interface {
	NewOne() interface{}
	NewSlice() interface{}
}

type CRUDer struct {
	m  CRUDerModelInterface
	db *gorm.DB
}

type Pagination struct {
	Limit, Offset int
}

type Order map[string]string

type Criterion struct {
	Field    string
	Value    interface{}
	Operator WhereOperator
}

func (Criterion *Criterion) Query() (string, interface{}) {
	return fmt.Sprintf("%s %s ?", Criterion.Field, Criterion.Operator), Criterion.Value
}

type Criteria []*Criterion

func (Criteria Criteria) Query() (string, []interface{}) {
	var (
		queries []string
		values  []interface{}
	)
	for _, c := range Criteria {
		var query, value = c.Query()
		queries = append(queries, query)
		values = append(values, value)
	}
	return strings.Join(queries, " AND "), values
}

type WhereOperator string

const (
	WhereOpEqual    WhereOperator = "="
	WhereOpGT       WhereOperator = ">"
	WhereOpGTEqual  WhereOperator = ">="
	WhereOpLT       WhereOperator = "<"
	WhereOpLTEqual  WhereOperator = "<="
	WhereOpNotEqual WhereOperator = "!="
	//WhereIn         WhereOperator = "in"
	WhereOpLike WhereOperator = "like"
	//WhereOpBetween  WhereOperator = "between"
)
