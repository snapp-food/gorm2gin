package gorm2gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"reflect"
)

type CRUDerModelInterface interface {
	NewOne() interface{}
	NewSlice() interface{}
}

type CRUDer struct {
	m  CRUDerModelInterface
	db *gorm.DB
}

func InitCRUDer(db *gorm.DB, model CRUDerModelInterface) *CRUDer {
	var cruder = new(CRUDer)
	cruder.m = model
	cruder.db = db
	return cruder
}

func (CRUDer *CRUDer) List(context *gin.Context) {
	var (
		res                                       = CRUDer.m.NewSlice()
		pageIndex, pageSize                       int
		pageIndexStr, hasPageIndex                = context.GetQuery("page_index")
		pageSizeStr, hasPageSize                  = context.GetQuery("page_size")
		pageOrderField, hasPageOrderField         = context.GetQuery("page_order_field")
		pageOrderDirection, hasPageOrderDirection = context.GetQuery("page_order_direction")
	)

	if !hasPageIndex {
		pageIndex = 1
	} else {
		pageIndex, _ = strconv.Atoi(pageIndexStr)
	}
	if !hasPageSize {
		pageSize = 10
	} else {
		pageSize, _ = strconv.Atoi(pageSizeStr)
	}
	if !hasPageOrderField || pageOrderField == "" {
		pageOrderField = "id"
	}
	if !hasPageOrderDirection {
		pageOrderDirection = "asc"
	}
	CRUDer.db.Offset(pageIndex - 1*pageSize).Limit(pageSize).Order(pageOrderField + " " + pageOrderDirection).Find(res)

	context.JSON(200, res)
}

func (CRUDer *CRUDer) Read(context *gin.Context) {
	var (
		id  = context.Param("rid")
		res = CRUDer.m.NewOne()
	)
	CRUDer.db.Find(res, id)
	context.JSON(200, res)
}

func (CRUDer *CRUDer) Update(context *gin.Context) {
	var (
		res      = CRUDer.m.NewOne()
		id       = context.Param("rid")
		intId, _ = strconv.Atoi(id)
	)
	var query = CRUDer.db.Find(res, id)
	context.BindJSON(res)
	reflect.Indirect(reflect.ValueOf(res)).FieldByName("ID").SetInt(int64(intId))
	query.Update(res)
	context.JSON(200, res)
}
