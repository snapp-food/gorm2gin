package gorm2gin

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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

func (CRUDer *CRUDer) List(c *gin.Context) {
	var res = CRUDer.m.NewSlice()
	CRUDer.db.Limit(100).Find(res)
	c.JSON(200, res)
}

func (CRUDer *CRUDer) Read(c *gin.Context) {
	var id = c.Param("rid")
	var res = CRUDer.m.NewOne()
	CRUDer.db.Find(res, id)
	c.JSON(200, res)
}
