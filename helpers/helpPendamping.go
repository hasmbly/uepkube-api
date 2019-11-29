package helpers

import (
	// "net/http"
	"github.com/labstack/echo"
	// "github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/db"
	"uepkube-api/models"
	"log"
	"math"	
	"fmt"	
)

func PaginatePendamping(c echo.Context, r *models.ResPagin) (err error) {
	u := &models.PosPagin{}
	num := 1

	// GetLoggedUser(c,"roles")

	if err := c.Bind(u); err != nil {
		return err
	}

	var co int = (u.Page - num) * u.Size
	
	PaginateResult, _ := PaginatePen(u,co,&CountRows)

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

func PaginatePen(f *models.PosPagin, offset int, count *int64) (ur []models.PaginatePendamping, err error) {

	// var Pendampings []models.Tbl_pendamping
	Pendampings := []models.PaginatePendamping{}

	con, err := db.CreateCon()
	if err != nil { return ur, echo.ErrInternalServerError }
	con.SingularTable(true)	

	q := con
	q = q.Table("tbl_pendamping t1")
	q = q.Limit(int(f.Size))
	q = q.Offset(int(offset))
	q = q.Select("t1.id_pendamping, t1.jenis_pendamping, t1.periode, t2.nama, t2.nik")
	q = q.Joins("join tbl_user t2 on t2.id_user = t1.id_pendamping")

	for i,_ := range f.Filters {
		k := f.Filters[i].Key
		o := f.Filters[i].Operation
		v := f.Filters[i].Value

		if o == "LIKE" || o == "like" {
			if v == "" { continue }
			if k == "roles_name" || k == "username" {
				q = q.Select("t1.id_pendamping, t1.jenis_pendamping, t1.periode, t2.nama, t2.nik, t3.username, t4.roles_name")
				q = q.Joins("join tbl_account t3 on t3.id_user = t1.id_pendamping")
				q = q.Joins("join tbl_roles t4 on t4.id = t3.id_roles")				
				q = q.Where(fmt.Sprintf("%s %s",k,o) + "?", "%"+v+"%")
			}
			q = q.Where(fmt.Sprintf("%s %s",k,o) + "?", "%"+v+"%")
		} else if o == ":" {
			if v == "" { 
				continue 
			} else {
				if k == "roles_name" || k == "username" {
					q = q.Select("t1.id_pendamping, t1.jenis_pendamping, t1.periode, t2.nama, t2.nik, t3.username, t4.roles_name")
					q = q.Joins("join tbl_account t3 on t3.id_user = t1.id_pendamping")
					q = q.Joins("join tbl_roles t4 on t4.id = t3.id_roles")				
				 	q = q.Where(fmt.Sprintf("%s ",k) + "=" + "?", v) 
				}			
			 	q = q.Where(fmt.Sprintf("%s ",k) + "=" + "?", v) 
			}
		}
	}
	q = q.Order(fmt.Sprintf("t1.%s %s",f.SortField,f.SortOrder))	
	
	q = q.Scan(&Pendampings)
	q = q.Limit(-1)
	q = q.Offset(-1)

	if len(Pendampings) != 0 {
		for i,_ := range Pendampings {
			var roles_name []string
			var account = models.Tbl_account{}

			con.Table("tbl_account").Where("id_user = ?", Pendampings[i].Id_pendamping).Select("username, id_roles").Find(&account)
			Pendampings[i].Username = account.Username

			if account.Id_roles == 0 {
				continue
			} else {

			con.Table("tbl_roles").Where("id = ?", account.Id_roles).Pluck("roles_name", &roles_name)
			}
			Pendampings[i].Roles_name = roles_name[0]
			
			log.Println("Username : ", account.Username)
			log.Println("Id_roles : ", account.Id_roles)
			log.Println("Roles_name : ", roles_name[0])
		}
	}

	if err := q.Count(count).Error; err != nil {
		return ur, err
	}

	// log.Println("result : ", Pendampings)

	defer con.Close()
	return Pendampings, nil
}