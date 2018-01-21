package gorm2gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"reflect"
)

func InitCRUDer(db *gorm.DB, model CRUDerModelInterface) *CRUDer {
	var cruder = new(CRUDer)
	cruder.m = model
	cruder.db = db
	return cruder
}

//todo:: normal queries like page_index,.. get a prefix like `_` ke beshe filter ha ro fieldName=value kar kard, hata fieldName>=value baye parser middleware!  regex :)
/*
limit
offset
order: map[string]string (field:direction)
criteria: map[string]string (field:operator:value)
 */
func (CRUDer *CRUDer) List(context *gin.Context) {
	var (
		count                                     int64
		criteria                                  = Criteria{}
		res                                       = CRUDer.m.NewSlice()
		pageIndex, pageSize                       int
		pageIndexStr, hasPageIndex                = context.GetQuery("_page_index")
		pageSizeStr, hasPageSize                  = context.GetQuery("_page_size")
		pageOrderField, hasPageOrderField         = context.GetQuery("_page_order_field")
		pageOrderDirection, hasPageOrderDirection = context.GetQuery("_page_order_direction")
	)

	var qStr = context.Request.URL.Query()
	for key, value := range qStr {
		if key[0] == '_' {
			continue
		}
		for _, v := range value {
			criteria = append(criteria, &Criterion{
				Field:    key,
				Value:    v,
				Operator: WhereOpEqual,
			})
		}
	}
	var queries, values = criteria.Query()

	//for _,cr := range criteria {
	//	fmt.Println(cr.Query())
	//}
	//fmt.Println(criteria.Query())

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

	CRUDer.db.
		Where(queries, values...).
		Model(res).Count(&count)
	CRUDer.db.
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Order(pageOrderField + " " + pageOrderDirection).
		Where(queries, values...).
		Find(res)

	context.JSON(200, map[string]interface{}{
		"results": res,
		"pagination": map[string]interface{}{
			"total":      count,
			"page_index": pageIndex,
			"page_size":  pageSize,
		},
	})
}

func (CRUDer *CRUDer) Create(context *gin.Context) {
	var res = CRUDer.m.NewOne()
	/*var err = */ context.ShouldBindJSON(res)
	/*if err != nil {
		context.Status(400)
		panic(err)
	}*/
	var rId = reflect.Indirect(reflect.ValueOf(res)).FieldByName("ID")
	rId.Set(reflect.Zero(rId.Type()))
	if err := CRUDer.db.Create(res).Error; err != nil {
		context.JSON(400, err)
	} else {
		context.JSON(201, res)
	}
}

func (CRUDer *CRUDer) Read(context *gin.Context) {
	var (
		id  = context.Param("rid")
		res = CRUDer.m.NewOne()
	)
	if CRUDer.db.Find(res, id).RecordNotFound() {
		context.Status(404)
	} else {
		context.JSON(200, res)
	}
}

func (CRUDer *CRUDer) Update(context *gin.Context) {
	var (
		res      = CRUDer.m.NewOne()
		id       = context.Param("rid")
		intId, _ = strconv.Atoi(id)
	)
	var query = CRUDer.db.Find(res, id)
	/*var err = */ context.ShouldBindJSON(res)
	//if err != nil {
	//	context.Status(400)
	//	panic(err)
	//}
	var i = int64(intId)
	reflect.Indirect(reflect.ValueOf(res)).FieldByName("ID").Set(reflect.ValueOf(&i))
	//reflect.Indirect(reflect.ValueOf(res)).FieldByName("ID").SetInt(int64(intId))
	if err := query.Update(res).Error; err != nil {
		context.JSON(400, err)
	} else {
		context.JSON(200, res)
	}
}

func (CRUDer *CRUDer) Delete(context *gin.Context) {
	var (
		id  = context.Param("rid")
		res = CRUDer.m.NewOne()
	)
	CRUDer.db.Delete(res, id)
	context.Status(200)
}

func CRUDerMiddleware(context *gin.Context) {
	//context.Request.URL.Query()

}
