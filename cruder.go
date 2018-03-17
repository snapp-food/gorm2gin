package gorm2gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"reflect"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin/binding"
)

//fixme: handle connection is closed errors

//todo:: sync and async mod! (bekhosus tu tr mod)

func InitCRUDer(db *gorm.DB, model CRUDerModelInterface) *CRUDer {
	var cruder = new(CRUDer)
	cruder.m = model
	cruder.db = db
	return cruder
}

func (CRUDer *CRUDer) GetDB(context *gin.Context) *gorm.DB {
	if trUid, trMode := context.GetQuery(queryParamTransactionUid); trMode {
		var id, err = strconv.ParseUint(trUid, 10, 64)
		if err != nil {
			panic(err)
		}
		return CRUDer.GetOrNewTransaction(id).DB
	} else {
		return CRUDer.db
	}
}

func (CRUDer *CRUDer) List(context *gin.Context) {
	var (
		count                                     int64
		pageOrder                                 []string
		criteria                                  = Criteria{}
		res                                       = CRUDer.m.NewSlice()
		qStr                                      = context.Request.URL.Query()
		limitStr, wantLimit                       = context.GetQuery(queryParamLimit)
		offsetStr, _                              = context.GetQuery(queryParamOffset)
		pageOrderField, _                         = context.GetQuery(queryParamPageOrderField)
		pageOrderDirection, hasPageOrderDirection = context.GetQuery(queryParamPageOrderDirection)
	)

	if !wantLimit {
		limitStr = "100"
	}

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

	if pageOrderField == "" {
		pageOrder = []string{"id"}
	} else {
		pageOrder = strings.Split(pageOrderField, ",")
	}
	if !hasPageOrderDirection {
		pageOrder[0] += " asc"
	} else {
		for i, dir := range strings.Split(pageOrderDirection, ",") {
			pageOrder[i] += " " + dir
		}
	}

	CRUDer.GetDB(context).
		Where(queries, values...).
		Model(res).Count(&count)
	CRUDer.GetDB(context).
		Offset(offsetStr).
		Limit(limitStr).
		Order(strings.Join(pageOrder, ",")).
		Where(queries, values...).
		Find(res)

	context.JSON(200, map[string]interface{}{
		"results": res,
		"pagination": map[string]interface{}{
			"total":  count,
			"offset": offsetStr,
			"limit":  limitStr,
		},
	})
}

func (CRUDer *CRUDer) Create(context *gin.Context) {
	var res = CRUDer.m.NewOne()
	//err :=context.ShouldBindJSON(res)
	err :=context.ShouldBindWith(res,binding.JSON)

	if err != nil {
		context.String(400, err.Error())
		panic(err)
	}
	var rId = reflect.Indirect(reflect.ValueOf(res)).FieldByName("ID")
	rId.Set(reflect.Zero(rId.Type()))
	if err := CRUDer.GetDB(context).Create(res).Error; err != nil {
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
	if CRUDer.GetDB(context).Find(res, id).RecordNotFound() {
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
	var query = CRUDer.GetDB(context).Find(res, id)
	//err :=context.ShouldBindJSON(res)
	err :=context.ShouldBindWith(res,binding.JSON)
	if err != nil {
		context.String(400, err.Error())
		panic(err)
	}
	var i = int64(intId)
	reflect.Indirect(reflect.ValueOf(res)).FieldByName("ID").Set(reflect.ValueOf(&i))
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
	CRUDer.GetDB(context).Delete(res, id)
	context.Status(200)
}

func CRUDerMiddleware(context *gin.Context) {
	//context.Request.URL.Query()

}
