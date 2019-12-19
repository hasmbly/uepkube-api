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

func PaginateMonev(c echo.Context, r *models.ResPagin) (err error) {
	u := &models.PosPagin{}
	num := 1

	// GetLoggedUser(c,"roles")

	if err := c.Bind(u); err != nil {
		return err
	}

	var co int = (u.Page - num) * u.Size
	
	PaginateResult, _ := ExecPaginateMonev(u,co,&CountRows)

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

func ExecPaginateMonev(f *models.PosPagin, offset int, count *int64) (ur []models.Tbl_monev_uepkube, err error) {

	Monevs := []models.Tbl_monev_uepkube{}

	con, err := db.CreateCon()
	if err != nil { return ur, echo.ErrInternalServerError }
	con.SingularTable(true)	

	q := con
	// q = q.Table("tbl_pelatihan t1")
	q = q.Model(&Monevs)
	q = q.Limit(int(f.Size))
	q = q.Offset(int(offset))
	q = q.Preload("Category")
	// q = q.Select("t1.id_pelatihan, t1.judul_pelatihan, t1.lokasi_pelatihan, t1.peruntukan, t1.start, t1.instruktur, t1.end, t1.deskripsi")

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
	q = q.Order(fmt.Sprintf("%s %s",f.SortField,f.SortOrder))	
	
	// q = q.Scan(&Monevs)
	q = q.Find(&Monevs)
	q = q.Limit(-1)
	q = q.Offset(-1)

	// get photos
	// if len(Monevs) != 0 {
	// 	for i,_ := range Monevs {
	// 		var pelatihan_photos []models.Tbl_pelatihan_files
	// 		// var account = models.Tbl_account{}

	// 		con.Table("tbl_pelatihan_files").Where("type = 'IMAGE' ").Where("id_pelatihan = ?", Monevs[i].Id_pelatihan).Select("tbl_pelatihan_files.*").Find(&pelatihan_photos)

	// 		for i,_ := range pelatihan_photos {
	// 			ImageBlob := pelatihan_photos[i].Files
	// 			pelatihan_photos[i].Files = "data:image/png;base64," + ImageBlob			
	// 		}
	// 		Monevs[i].Photo = pelatihan_photos
	// 	}
	// }

	if err := q.Count(count).Error; err != nil {
		return ur, err
	}

	// log.Println("result : ", Monevs)

	defer con.Close()
	return Monevs, nil
}	