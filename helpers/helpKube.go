package helpers

import (
	// "net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/models"
	"uepkube-api/db"
	"github.com/fatih/structs"
	"fmt"
	"log"
	"strconv"
	// "github.com/ulule/paging"
	"math"
)

func SetMemberNameKube(s *models.ShowKube, Kube models.Tbl_kube) error {
	/*prepare DB*/
	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)
	/*
	begin:find member name of Kube
	 */
	kv := structs.Values(Kube)
	kr := kv[4:15]
	
	var tmp []string
	res := make([]string, len(kr))
	ints := make([]int, len(kr))
	
	for i,d := range kr {
		nf := fmt.Sprintf("%+v",d)
		ints[i],_ = strconv.Atoi(nf)
	}

    for i,d := range ints {
		if err := con.Table("tbl_user").Where(&models.Tbl_user{Id_user:d}).Pluck("nama", &tmp).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
    	if d == 0 { tmp[0] = "" }
		res[i] = tmp[0]
    }

    // get alamat from ketua kube (first man)
    var alamat []*string
	if err := con.Table("tbl_user").Where(&models.Tbl_user{Id_user: ints[0]}).Pluck("alamat", &alamat).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

    // get lat from ketua kube (first man)
    var lat []*string
	if err := con.Table("tbl_user").Where(&models.Tbl_user{Id_user: ints[0]}).Pluck("lat", &lat).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

    // get long from ketua kube (first man)
    var lng []*string
	if err := con.Table("tbl_user").Where(&models.Tbl_user{Id_user: ints[0]}).Pluck("lng", &lng).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

    // get photo kube
    var photo []models.Tbl_uepkube_files

	// if err := con.Table("tbl_uepkube_photo").Where(&models.Tbl_uepkube_photo{Id_kube: Kube.Id_kube}).Find(&photo).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// for i,_ := range photo {

	// 		if photo[i].Photo != "" {
	// 			ImageBlob := photo[i].Photo
	// 			photo[i].Photo = "data:image/png;base64," + ImageBlob	
	// 		}

	// 	}

    /*
	end:find member name of Kube
	 */	
	*s = models.ShowKube{
		Id_kube: 		Kube.Id_kube,
		Nama_kube: 		Kube.Nama_kube,
		Nama_usaha: 	Kube.Nama_usaha,
		Alamat: 		alamat[0],
		Lat: 			lat[0],
		Lng: 			lng[0],
		Items: models.Items{
			Ketua:			res[0],
			Sekertaris:		res[1],
			Bendahara:  	res[2],
			Anggota1:		res[3],
			Anggota2:		res[4],
			Anggota3:		res[5],
			Anggota4:		res[6],
			Anggota5:		res[7],
			Anggota6: 		res[8],
			Anggota7:		res[9],
			Pendamping:		res[10],
		},
		Photo:			photo,
		Flag:			"KUBE",
	}
	defer con.Close()
	return err
}

func PaginateKube(c echo.Context, r *models.ResPagin) (err error) {
	flag = "KUBE"
	host = c.Request().Host		

	u := &models.PosPagin{}
	num := 1

	// GetLoggedUser(c,"roles")

	if err := c.Bind(u); err != nil {
		return err
	}

	var co int = (u.Page - num) * u.Size
	
	PaginateResult, _ := ExecPaginateKube(u,co,&CountRows)

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

func ExecPaginateKube(f *models.PosPagin, offset int, count *int64) (ur []models.Tbl_kube, err error) {

	kubes := []models.Tbl_kube{}

	con, err := db.CreateCon()
	if err != nil { return ur, echo.ErrInternalServerError }
	con.SingularTable(true)	

	q := con
	q = q.Model(&kubes)
	q = q.Limit(int(f.Size))
	q = q.Offset(int(offset))
	q = q.Preload("JenisUsaha")
	q = q.Preload("LapkeuHistory")
	q = q.Preload("MonevHistory")
	q = q.Preload("InventarisHistory")
	q = q.Preload("PelatihanHistory")
	q = q.Preload("Pendamping", func(q *gorm.DB) *gorm.DB {
		return q.Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Select("tbl_pendamping.*,tbl_user.nama")
	})

	// q = q.Preload("Region")
	// q = q.Preload("Kelurahan")
	// q = q.Preload("Kecamatan")
	// q = q.Preload("Kabupaten")	

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
	q = q.Find(&kubes)
	q = q.Limit(-1)
	q = q.Offset(-1)

	// get kubes_member
	if len(kubes) != 0 {
		for i, _ := range kubes {
			id := kubes[i].Id_kube
			var KubesMember = []string{"ketua", "sekertaris", "bendahara", "anggota1", "anggota2", "anggota3", "anggota4", "anggota5", "anggota6", "anggota7"}
			tmp := []models.Kubes_items{}

			for x, _ := range KubesMember {
			
				con.Table("tbl_kube t1").Select("t2.*, '" + KubesMember[x] + "' as posisi").Joins("join tbl_user t2 on t2.id_user = t1." + KubesMember[x]).Where("id_kube = ?", id).Scan(&tmp)

				if len(tmp) != 0 {
				// get region from ketua
				Regions := models.View_address{}
				if err := con.Table("view_address").Select("view_address.*").Where("id_kelurahan = ?", tmp[0].Id_kelurahan).Scan(&Regions).Error; err != nil { 
					log.Println("err : ", err)
					continue
				}								
					kubes[i].Region = Regions
					kubes[i].Items = append(kubes[i].Items, tmp[0])
				}
			}	
		}		
	}

	// // get Pendampings
	// if len(kubes) != 0 {
	// 	for i,_ := range kubes {
	// 		var id_pendamping []int
	// 		var pendamping models.CustomPendamping

	// 		con.Table("tbl_kube").Where("id_kube = ?", kubes[i].Id_kube).Pluck("id_pendamping", &id_pendamping)

	// 		if len(id_pendamping) != 0 {

	// 			for i,_ := range id_pendamping {
	// 				con.Table("tbl_pendamping").Select("tbl_pendamping.*, tbl_user.nama").Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Where("id_pendamping = ?", id_pendamping[i]).Find(&pendamping)
	// 			}
	// 				kubes[i].Pendamping = pendamping
	// 		}
	// 	}
	// }

	// // get Usaha
	// if len(kubes) != 0 {
	// 	for i,_ := range kubes {
	// 		var kube_usaha models.UsahaKube
	// 		photos := []models.Tbl_uepkube_files{}

	// 		// log.Println("id_kube : ", kubes[i].Id_kube)
			
	// 		q := con.Table("tbl_kube t1")
	// 		q = q.Select("t1.id_kube, t1.nama_usaha, t2.id_usaha, t2.jenis_usaha")
	// 		q = q.Joins("join tbl_jenis_usaha t2 on t2.id_usaha = t1.id_jenis_usaha")
	// 		q = q.Where("t1.id_kube = ?", kubes[i].Id_kube)
	// 		q = q.Scan(&kube_usaha)

	// 		if kube_usaha.Id_usaha != 0 { kubes[i].Usaha = kube_usaha }

	// 		con.Table("tbl_uepkube_files").Where("id_kube = ?", kubes[i].Id_kube).Find(&photos)

	// 		// exec for files
	// 		id := kubes[i].Id_kube
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
	// 				kubes[i].Usaha.Photo = append(kubes[i].Usaha.Photo, photos[index])
	// 			}

	// 		}			
			
	// 	}
	// }

	// // get hitory_periods
	// if len(kubes) != 0 {
	// 	for i,_ := range kubes {
	// 		history_periods := []*models.Tbl_periods_uepkube{}
	// 		con.Table("tbl_periods_uepkube").Select("*").Where("id_kube = ?", kubes[i].Id_kube).Scan(&history_periods)
			
	// 		if len(history_periods) != 0 {
	// 			for index, _ := range history_periods {
	// 				kubes[i].PeriodsHistory = append(kubes[i].PeriodsHistory, history_periods[index])
	// 			}
	// 		}
	// 	}
	// }		

	// // get bantuan_periods
	// if len(kubes) != 0 {
	// 	for i,_ := range kubes {
	// 		bantuan_periods := models.Tbl_bantuan_periods{}
			
	// 		if len(kubes[i].PeriodsHistory) != 0 {
	// 			for index, _ := range kubes[i].PeriodsHistory {
	// 				con.Table("tbl_bantuan_periods").Select("*").Where("id = ?", kubes[i].PeriodsHistory[index].Id_periods).Scan(&bantuan_periods)

	// 					kubes[i].PeriodsHistory[index].BantuanPeriods = &bantuan_periods
	// 			}
	// 		}
	// 	}
	// }	

	// // get credit_debit
	// if len(kubes) != 0 {
	// 	for i,_ := range kubes {
	// 		credit_debit := []*models.Tbl_inventaris{}

	// 		con.Table("tbl_credit_debit").Select("*").Where("id_kube = ?", kubes[i].Id_kube).Scan(&credit_debit)
			
	// 		if len(credit_debit) != 0 {
	// 			for indexDebit, _ := range credit_debit {
	// 				for indexPeriods, _ := range kubes[i].PeriodsHistory {
	// 					kubes[i].PeriodsHistory[indexPeriods].BantuanPeriods.CreditDebit = append(kubes[i].PeriodsHistory[indexPeriods].BantuanPeriods.CreditDebit, credit_debit[indexDebit])
	// 				}
	// 			}
	// 		}
	// 	}
	// }	

	if err := q.Count(count).Error; err != nil {
		return ur, err
	}

	defer con.Close()
	return kubes, nil
}