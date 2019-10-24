package db
 
import (
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
)
 
func CreateCon() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "kube_dinsos:kube_dinsos2019@/kube_dinsos?charset=utf8&parseTime=True&loc=Local")
	return db, err
}