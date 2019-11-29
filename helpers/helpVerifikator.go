package helpers

import (
	// "net/http"
	"github.com/labstack/echo"
	// "github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/db"
	"uepkube-api/models"
	// "log"
	"math"	
	"fmt"	
)

func PaginateVerifikator(c echo.Context, r *models.ResPagin) (err error) {
	u := &models.PosPagin{}
	num := 1

	// GetLoggedUser(c,"roles")

	if err := c.Bind(u); err != nil {
		return err
	}

	var co int = (u.Page - num) * u.Size
	
	PaginateResult, _ := PaginateVer(u,co,&CountRows)

	l := int64(u.Size)
	o := int64(co)
	t := CountRows
	f := false
	la := false
	tp := float64(t)/float64(l)
	rtp := math.Ceil(tp)
	if rtp == 0 { rtp = rtp+1 }

	if u.Page == 1 {f = true}
	if u.Page == int(rtp) {la = true}

	// log.Println("Result is : ", PaginateResult)

	*r = models.ResPagin{
		Content:PaginateResult,
		First:f,
		Last:la,
		Number:u.Page,
		NumberOfElement:l,		
		Pageable: models.Pageable{
			Offset:o,
			PageNumber:u.Page,
			PageSize:l,
			Paged:true,
			Unpaged:false,
		},
		Sort: models.Sort{
			Sorted:true,
			Unsorted:false,
		},
		TotalPages:rtp,
		TotalElements:t,		
	}
	return err
}

func PaginateVer(f *models.PosPagin, offset int, count *int64) (ur []models.PaginateVerifikator, err error) {

	// var Verifikators []models.Tbl_verifikator
	Verifikators := []models.PaginateVerifikator{}

	con, err := db.CreateCon()
	if err != nil { return ur, echo.ErrInternalServerError }
	con.SingularTable(true)	

	q := con
	q = q.Table("tbl_account t1")
	q = q.Limit(int(f.Size))
	q = q.Offset(int(offset))
	q = q.Select("t1.id_user, t1.username, t2.nama, t2.nik, t3.roles_name")
	q = q.Joins("join tbl_user t2 on t2.id_user = t1.id_user")
	q = q.Joins("join tbl_roles t3 on t3.id = t1.id_roles")
	q = q.Where("t1.id_roles = ?", 2)

	for i,_ := range f.Filters {
		k := f.Filters[i].Key
		o := f.Filters[i].Operation
		v := f.Filters[i].Value

		if o == "LIKE" || o == "like" {
			if v == "" { continue }
			q = q.Where(fmt.Sprintf("%s %s",k,o) + "?", "%"+v+"%")
		} else if o == ":" {
			if v == "" { 
				continue 
			} else {		
			 	q = q.Where(fmt.Sprintf("%s ",k) + "=" + "?", v) 
			}
		}
	}
	q = q.Order(fmt.Sprintf("t1.%s %s",f.SortField,f.SortOrder))	
	
	q = q.Scan(&Verifikators)
	q = q.Limit(-1)
	q = q.Offset(-1)

	if err := q.Count(count).Error; err != nil {
		return ur, err
	}

	// log.Println("result : ", Verifikators)

	defer con.Close()
	return Verifikators, nil
}