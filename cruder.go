package gorm2gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"reflect"
	"strings"
)

//todo: transactional requests: ye flag e transactionMode migire ba ye UniqueID (rand) ke baraye majmueye trans estefade mikone
// va akhare kar Commit ya Rollback call mishe!
// mitune Sync va ASync bashe juri ke poshte har request query bezane ya ASync query bezane va faghat tu commit result ro bege!
// sare Commit mibine age hameye query ha anjam shode ke result o mide magarna sabr mikone hamashun anjam beshan bad result mide!
// todo INJURIIIIIIIIIIII: ye service e singleton darim ke handler e transactionas:
// vaghti ye tr e jadid miad ba hamin uid:::  tx := db.Begin() ro zakhire mikone!! va bad ru un query mizane dige!
// akharesham ru hamun tx.commit ya tx.rollback HOOOOOOOOOOOOOORAAAAy
// har kodumam ye expire e moshakhas dare ke age vaziatesh moshakhas nashod khodesh rollback ya commit kone!

func InitCRUDer(db *gorm.DB, model CRUDerModelInterface) *CRUDer {
	var cruder = new(CRUDer)
	cruder.m = model
	cruder.db = db
	return cruder
}

func (CRUDer *CRUDer) List(context *gin.Context) {
	var (
		count                                     int64
		criteria                                  = Criteria{}
		res                                       = CRUDer.m.NewSlice()
		qStr                                      = context.Request.URL.Query()
		offsetStr, _                              = context.GetQuery("_offset")
		limitStr, _                               = context.GetQuery("_limit")
		pageOrderField, _                         = context.GetQuery("_page_order_field")
		pageOrderDirection, hasPageOrderDirection = context.GetQuery("_page_order_direction")
		pageOrder                                 []string
	)

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

	CRUDer.db.
		Where(queries, values...).
		Model(res).Count(&count)
	CRUDer.db.
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
