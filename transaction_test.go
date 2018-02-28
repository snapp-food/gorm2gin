package gorm2gin

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
	"sync"
	"testing"
	"time"
)

var Database, _ = gorm.Open("mysql", "root:123456@/testing?charset=utf8&parseTime=True&loc=Local")

func init() {
	Database.AutoMigrate(&TestTable{})
}

type TestTable struct {
	ID           *int64 `gorm:"primary_key;auto_increment" json:"id"`
	Name, Family string
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

func (v *TestTable) NewOne() interface{} {
	return new(TestTable)
}
func (v *TestTable) NewSlice() interface{} {
	return new([]TestTable)
}

func TestCRUDer_GetOrNewTransaction(t *testing.T) {
	var TestTableCRUD = InitCRUDer(Database, new(TestTable))
	var TestTableCRUD2 = InitCRUDer(Database, new(TestTable))
	var tr = TestTableCRUD.GetOrNewTransaction(123456)
	//Database.Create(&TestTable{
	tr.DB.Create(&TestTable{
		Name:   "Saeed",
		Family: "Falsafin",
	})

	Database.Create(&TestTable{
		Name:   "Ali",
		Family: "Falsafin",
	})

	var tr12 = TestTableCRUD.GetOrNewTransaction(123456)

	tr12.DB.Create(&TestTable{
		Name:   "Saeede",
		Family: "Falsafin",
	})

	var tr2 = TestTableCRUD.GetOrNewTransaction(321321)

	tr2.DB.Create(&TestTable{
		Name:   "NADER",
		Family: "Falsafin",
	})
	var tr3 = TestTableCRUD2.GetOrNewTransaction(321321)

	tr3.DB.Create(&TestTable{
		Name:   "NADER",
		Family: "Haderi",
	})

	//tr12.DB.Commit()
	//tr3.DB.Commit()
	tr.DB.Rollback()
	tr2.DB.Commit()
}

func TestID(t *testing.T) {
	fmt.Println(strconv.ParseUint("1590206877993462", 10, 64))
	fmt.Println(strconv.ParseUint("5a648e41509bc", 16, 64))

}

var w = new(sync.WaitGroup)

func TestTick(t *testing.T) {

	w.Add(1)
	go tt(time.After(time.Second * 2))

	fmt.Println("salam")
	w.Wait()
}

func tt(c <-chan time.Time) {
	<-c
	fmt.Println("2 Secs later")
	w.Done()
}
