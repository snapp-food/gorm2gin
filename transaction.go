package gorm2gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

const TrLifeTime = time.Second * 2

type Transaction struct {
	DB *gorm.DB
}

type transactions map[uint64]*Transaction

var Transactions = transactions{}

func (CRUDer *CRUDer) GetOrNewTransaction(id uint64) *Transaction {
	if tr, ok := Transactions[id]; ok {
		return tr
	} else { // Begin:
		Transactions[id] = new(Transaction)
		Transactions[id].DB = CRUDer.db.Begin()
		go CRUDer.DestroyDeadTr(time.After(TrLifeTime), id)
		return Transactions[id]
	}
}

func (CRUDer *CRUDer) DestroyDeadTr(c <-chan time.Time, id uint64) {
	<-c
	if tr, ok := Transactions[id]; ok {
		tr.DB.Rollback()
		delete(Transactions, id)
		fmt.Println("Dead transaction:", id)
	}
}

func (CRUDer *CRUDer) CommitTransaction(context *gin.Context) {
	var trUid, _ = context.GetQuery(queryParamTransactionUid)
	var id, err = strconv.ParseUint(trUid, 10, 64)
	if tr, ok := Transactions[id]; err != nil || !ok {
		context.Status(400)
	} else if err := tr.DB.Commit().Error; err != nil {
		tr.DB.Rollback()
		context.Status(500)
	} else {
		delete(Transactions, id)
		context.Status(200)
	}
}
func (CRUDer *CRUDer) RollbackTransaction(context *gin.Context) {
	var trUid, _ = context.GetQuery(queryParamTransactionUid)
	var id, err = strconv.ParseUint(trUid, 10, 64)
	if tr, ok := Transactions[id]; err != nil || !ok {
		context.Status(400)
	} else {
		tr.DB.Rollback()
		delete(Transactions, id)
		context.Status(200)
	}
}
