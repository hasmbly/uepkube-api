package helpers

import (
	// "net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/db"
	"uepkube-api/models"
	// "log"
	"math"	
	"fmt"
)

// var tmpPath, urlPath, blobFile,flag,host string
// 

type Tbl_uep struct {
	*models.Tbl_uep
	*models.Tbl_user
}

func PaginateUep(c echo.Context, r *models.ResPagin) (err error) {
	flag = "UEP"
	host = c.Request().Host	

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

func ExecPaginateUep(f *models.PosPagin, offset int, count *int64) (ur []Tbl_uep, err error) {

	// var Pelatihans []models.Tbl_pendamping
	Ueps := []Tbl_uep{}

	con, err := db.CreateCon()
	if err != nil { return ur, echo.ErrInternalServerError }
	con.SingularTable(true)	

	q := con
	q = q.Model(&Ueps)
	q = q.Joins("join tbl_user on tbl_user.id_user = tbl_uep.id_uep")
	q = q.Select("tbl_uep.id_uep,tbl_uep.id_pendamping,tbl_uep.nama_usaha,tbl_uep.id_jenis_usaha,tbl_uep.status,tbl_uep.nama_usaha,tbl_uep.created_at,tbl_uep.updated_at,tbl_uep.bantuan,tbl_user.*")
	q = q.Limit(int(f.Size))
	q = q.Offset(int(offset))
	q = q.Preload("JenisUsaha")
	q = q.Preload("LapkeuHistory")
	q = q.Preload("MonevHistory.Category")
	q = q.Preload("InventarisHistory")
	q = q.Preload("PelatihanHistory")
	q = q.Preload("Region")
	q = q.Preload("Pendamping", func(q *gorm.DB) *gorm.DB {
		return q.Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Select("tbl_pendamping.*,tbl_user.nama")
	})
	q = q.Preload("Kelurahan")
	q = q.Preload("Kecamatan")
	q = q.Preload("Kabupaten")	

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
	q = q.Order(fmt.Sprintf("tbl_uep.%s %s",f.SortField,f.SortOrder))	
	
	q = q.Find(&Ueps)
	q = q.Limit(-1)
	q = q.Offset(-1)

	// // get Pendampings
	// if len(Ueps) != 0 {
	// 	for i,_ := range Ueps {
	// 		var id_pendamping []int
	// 		var pendamping models.CustomPendamping

	// 		con.Table("tbl_uep").Where("id_uep = ?", Ueps[i].Id_uep).Pluck("id_pendamping", &id_pendamping)

	// 		for i,_ := range id_pendamping {
	// 			con.Table("tbl_pendamping").Select("tbl_pendamping.*, tbl_user.nama").Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Where("id_pendamping = ?", id_pendamping[i]).Find(&pendamping)
	// 		}
	// 			Ueps[i].Pendamping = pendamping
	// 	}
	// }

	// // get Usaha
	// if len(Ueps) != 0 {
	// 	for i,_ := range Ueps {
	// 		var uep_usaha models.UsahaUep
	// 		photos := []models.Tbl_uepkube_files{}

	// 		q := con.Table("tbl_uep t1")
	// 		q = q.Select("t1.id_uep, t1.nama_usaha, t2.id_usaha, t2.jenis_usaha")
	// 		q = q.Joins("join tbl_jenis_usaha t2 on t2.id_usaha = t1.id_jenis_usaha")
	// 		q = q.Where("t1.id_uep = ?", Ueps[i].Id_uep)
	// 		q = q.Scan(&uep_usaha)

	// 		if uep_usaha.Id_usaha != 0 { Ueps[i].Usaha = uep_usaha }

	// 		// get usaha_files
	// 		con.Table("tbl_uepkube_files").Where("id_uep = ?", Ueps[i].Id_uep).Find(&photos)

	// 		// exec for files
	// 		id := Ueps[i].Id_uep
	// 		for index,_ := range photos {

	// 			if photos[index].Type == "IMAGE" {

	// 				id_photo := photos[index].Id
					
	// 				tmpPath	= fmt.Sprintf(GoPath + "/src/uepkube-api/static/assets/images/%s_id_%d_photo_id_%d.png", flag,id,id_photo)
	// 				urlPath	= fmt.Sprintf("http://%s/images/%s_id_%d_photo_id_%d.png", host,flag,id,id_photo)
	// 				blobFile = photos[index].Files

	// 				if check := CreateFile(tmpPath, blobFile); check == false {
	// 					log.Println("blob is empty : ", check)
	// 				}
	// 				photos[index].Files = urlPath
	// 				Ueps[i].Usaha.Photo = append(Ueps[i].Usaha.Photo, photos[index])
	// 			}
	// 		}
			
	// 	}
	// }

	// // get hitory_periods
	// if len(Ueps) != 0 {
	// 	for i,_ := range Ueps {
	// 		history_periods := []*models.Tbl_periods_uepkube{}
	// 		con.Table("tbl_periods_uepkube").Select("*").Where("id_uep = ?", Ueps[i].Id_uep).Scan(&history_periods)
			
	// 		if len(history_periods) != 0 {
	// 			for index, _ := range history_periods {
	// 				Ueps[i].PeriodsHistory = append(Ueps[i].PeriodsHistory, history_periods[index])
	// 			}
	// 		}
	// 	}
	// }		

	// // get bantuan_periods
	// if len(Ueps) != 0 {
	// 	for i,_ := range Ueps {
	// 		bantuan_periods := models.Tbl_bantuan_periods{}
			
	// 		if len(Ueps[i].PeriodsHistory) != 0 {
	// 			for index, _ := range Ueps[i].PeriodsHistory {
	// 				con.Table("tbl_bantuan_periods").Select("*").Where("id = ?", Ueps[i].PeriodsHistory[index].Id_periods).Scan(&bantuan_periods)

	// 					Ueps[i].PeriodsHistory[index].BantuanPeriods = &bantuan_periods
	// 			}
	// 		}
	// 	}
	// }	

	// // get credit_debit
	// if len(Ueps) != 0 {
	// 	for i,_ := range Ueps {
	// 		credit_debit := []*models.Tbl_inventaris{}

	// 		con.Table("tbl_credit_debit").Select("*").Where("id_uep = ?", Ueps[i].Id_uep).Scan(&credit_debit)
			
	// 		if len(credit_debit) != 0 {
	// 			for indexDebit, _ := range credit_debit {
	// 				for indexPeriods, _ := range Ueps[i].PeriodsHistory {
	// 					Ueps[i].PeriodsHistory[indexPeriods].BantuanPeriods.CreditDebit = append(Ueps[i].PeriodsHistory[indexPeriods].BantuanPeriods.CreditDebit, credit_debit[indexDebit])
	// 				}
	// 			}
	// 		}
	// 	}
	// }		

	if err := q.Count(count).Error; err != nil {
		return ur, err
	}

	defer con.Close()
	return Ueps, nil
}