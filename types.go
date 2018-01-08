package gorm2gin

type Pagination struct {
	Limit, Offset int
}

type Order map[string]string

type Criterion struct {
	Field    string
	Value    interface{}
	Operator WhereOperator
}

type Criteria []Criterion

type WhereOperator string

const (
	WhereOpEqual    WhereOperator = "="
	WhereOpGT       WhereOperator = ">"
	WhereOpGTEqual  WhereOperator = ">="
	WhereOpLT       WhereOperator = "<"
	WhereOpLTEqual  WhereOperator = "<="
	WhereOpNotEqual WhereOperator = "<>"
	//WhereOpLike     WhereOperator = "like"
	//WhereOpBetween  WhereOperator = "between"
)
