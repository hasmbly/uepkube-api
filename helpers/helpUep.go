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

func PaginateUep(c echo.Context, r *models.ResPagin) (err error) {
	u := &models.PosPagin{}
	num := 1

	// GetLoggedUser(c,"roles")

	if err := c.Bind(u); err != nil {
		return err
	}

	var co int = (u.Page - num) * u.Size
	
	PaginateResult, _ := ExecPaginateUep(u,co,&CountRows)

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

func ExecPaginateUep(f *models.PosPagin, offset int, count *int64) (ur []models.PaginateUep, err error) {

	// var Pelatihans []models.Tbl_pendamping
	Ueps := []models.PaginateUep{}

	con, err := db.CreateCon()
	if err != nil { return ur, echo.ErrInternalServerError }
	con.SingularTable(true)	

	q := con
	q = q.Table("tbl_uep t1")
	q = q.Limit(int(f.Size))
	q = q.Offset(int(offset))
	q = q.Select("t1.id_uep, t2.nama, t2.nik, t2.no_kk, t2.alamat, t1.status, t1.created_at")
	q = q.Joins("join tbl_user t2 on t2.id_user = t1.id_uep")
	q = q.Joins("join tbl_jenis_usaha t3 on t3.id_usaha = t1.id_jenis_usaha")

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
	
	q = q.Scan(&Ueps)
	q = q.Limit(-1)
	q = q.Offset(-1)

	// get Pendampings
	if len(Ueps) != 0 {
		for i,_ := range Ueps {
			var id_pendamping []int
			var pendamping models.CustomPendamping

			con.Table("tbl_uep").Where("id_uep = ?", Ueps[i].Id_uep).Pluck("id_pendamping", &id_pendamping)

			for i,_ := range id_pendamping {
				con.Table("tbl_pendamping").Select("tbl_pendamping.*, tbl_user.nama as nama_pendamping").Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Where("id_pendamping = ?", id_pendamping[i]).Find(&pendamping)
			}
				Ueps[i].Pendamping = pendamping
		}
	}

	// get Usaha
	if len(Ueps) != 0 {
		for i,_ := range Ueps {
			var uep_usaha models.UsahaUep
 			// var id_produk []int
			var photos []models.Tbl_usaha_uep_photo

			log.Println("id_uep : ", Ueps[i].Id_uep)
			q := con.Table("tbl_uep t1")
			q = q.Select("t1.id_uep, t1.nama_usaha, t2.id_usaha, t2.jenis_usaha")
			q = q.Joins("join tbl_jenis_usaha t2 on t2.id_usaha = t1.id_jenis_usaha")
			q = q.Where("t1.id_uep = ?", Ueps[i].Id_uep)
			q = q.Scan(&uep_usaha)

			if uep_usaha.Id_usaha != 0 { Ueps[i].Usaha = uep_usaha }

			con.Table("tbl_usaha_uep_photo").Where("id_uep = ?", Ueps[i].Id_uep).Find(&photos)

			for index,_ := range photos {
				ImageBlob := photos[index].Photo
				photos[index].Photo = "data:image/png;base64," + ImageBlob			
				Ueps[i].Usaha.Photo = photos
			}

			// log.Println("photos : ", photos)
			log.Println("usaha : ", Ueps[i].Usaha)
			
		}
	}

	if err := q.Count(count).Error; err != nil {
		return ur, err
	}

	// log.Println("result : ", Pelatihans)

	defer con.Close()
	return Ueps, nil
}