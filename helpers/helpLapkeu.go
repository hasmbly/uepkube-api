package helpers

import (
	// "net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/db"
	"uepkube-api/models"
	"log"
	"math"	
	"fmt"
)

func PaginateLapkeu(c echo.Context, r *models.ResPagin) (err error) {
	flag = "LAPKEU"
	host = c.Request().Host
	
	u := &models.PosPagin{}
	num := 1

	// GetLoggedUser(c,"roles")
	
	if err := c.Bind(u); err != nil {
		return err
	}

	var co int = (u.Page - num) * u.Size
	
	PaginateResult, _ := ExecPaginateLapkeu(u,co,&CountRows)

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

func ExecPaginateLapkeu(f *models.PosPagin, offset int, count *int64) (ur []models.Tbl_lapkeu_uepkube, err error) {

	Lapkeu := []models.Tbl_lapkeu_uepkube{}

	con, err := db.CreateCon()
	if err != nil { return ur, echo.ErrInternalServerError }
	con.SingularTable(true)	

	q := con
	q = q.Model(&Lapkeu)
	q = q.Limit(int(f.Size))
	q = q.Offset(int(offset))
	q = q.Preload("Pendamping")
	q = q.Preload("Photo")

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
	
	q = q.Find(&Lapkeu)
	q = q.Limit(-1)
	q = q.Offset(-1)

	//getDetail
	if len(Lapkeu) != 0 {

		for i, _ := range Lapkeu {

			if Lapkeu[i].Id_uep != 0 {
				id := Lapkeu[i].Id_uep
				
				User 	:= Tbl_user{}
				q := con
				q = q.Model(&User)
				q = q.Joins("join tbl_uep on tbl_uep.id_uep = tbl_user.id_user")
				q = q.Select("tbl_uep.*, tbl_user.*")
				q = q.Preload("JenisUsaha")
				q = q.Preload("PeriodsHistory.BantuanPeriods.Usaha", func(q *gorm.DB) *gorm.DB {
					return q.Where("id_uep = ?", id).Preload("JenisUsaha")
				})
				q = q.Preload("PeriodsHistory.BantuanPeriods.Usaha.AllProduk.DetailProduk.JenisProduk")
				q = q.Preload("PeriodsHistory.BantuanPeriods.CreditDebit", func(q *gorm.DB) *gorm.DB {
					return q.Where("id_uep = ?", id)
				})
				q = q.Preload("Pendamping")
				q = q.Preload("Kelurahan")
				q = q.Preload("Kecamatan")
				q = q.Preload("Kabupaten")
				q = q.Preload("Photo", func(q *gorm.DB) *gorm.DB {
					return q.Where("id_uep = ?", id)
				})
				q = q.First(&User, id)

				for index, _ := range User.Photo {
						id_photo := User.Photo[index].Id

						tmpPath	= fmt.Sprintf(GoPath + "/src/uepkube-api/static/assets/images/%s_id_%d_photo_id_%d.png", flag,id,id_photo)
						urlPath	= fmt.Sprintf("http://%s/images/%s_id_%d_photo_id_%d.png", host,flag,id,id_photo)
						blobFile = User.Photo[index].Files

						if check := CreateFile(tmpPath, blobFile); check == false {
							log.Println("blob is empty : ", check)
						}
					
						User.Photo[index].Files = urlPath
				}

				Lapkeu[i].Detail = User

			} else if Lapkeu[i].Id_kube != 0 {

				id := Lapkeu[i].Id_kube
				
				Kube 	:= models.Tbl_kube{}
				q := con
				q = q.Model(&Kube)
				q = q.Preload("JenisUsaha")
				q = q.Preload("PeriodsHistory.BantuanPeriods.Usaha", func(q *gorm.DB) *gorm.DB {
					return q.Where("id_kube = ?", id).Preload("JenisUsaha")
				})
				q = q.Preload("PeriodsHistory.BantuanPeriods.Usaha.AllProduk.DetailProduk.JenisProduk")			
				q = q.Preload("PeriodsHistory.BantuanPeriods.CreditDebit", func(q *gorm.DB) *gorm.DB {
					return q.Where("id_kube = ?", id)
				})
				q = q.Preload("Pendamping")
				q = q.Preload("Photo", func(q *gorm.DB) *gorm.DB {
					return q.Where("id_kube = ?", id)	
				})
				q = q.First(&Kube, id)

				for index, _ := range Kube.Photo {
						id_photo := Kube.Photo[index].Id

						tmpPath	= fmt.Sprintf(GoPath + "/src/uepkube-api/static/assets/images/%s_id_%d_photo_id_%d.png", flag,id,id_photo)
						urlPath	= fmt.Sprintf("http://%s/images/%s_id_%d_photo_id_%d.png", host,flag,id,id_photo)
						blobFile = Kube.Photo[index].Files

						if check := CreateFile(tmpPath, blobFile); check == false {
							log.Println("blob is empty : ", check)
						}
					
						Kube.Photo[index].Files = urlPath
				}

				Lapkeu[i].Detail = Kube

			}
		}
	}

	// photo
	if len(Lapkeu) != 0 {
		for i, _ := range Lapkeu {
			id := Lapkeu[i].Id
			if len(Lapkeu[i].Photo) != 0 {
				for x, _ := range Lapkeu[i].Photo {

					if Lapkeu[i].Photo[x].Type == "IMAGE" {

						id_photo := Lapkeu[i].Photo[x].Id
						
						tmpPath	= fmt.Sprintf(GoPath + "/src/uepkube-api/static/assets/images/%s_id_%d_photo_id_%d.png", flag,id,id_photo)
						urlPath	= fmt.Sprintf("http://%s/images/%s_id_%d_photo_id_%d.png", host,flag,id,id_photo)
						blobFile = Lapkeu[i].Photo[x].Files

						if check := CreateFile(tmpPath, blobFile); check == false {
							log.Println("blob is empty : ", check)
						}
						Lapkeu[i].Photo[x].Files = urlPath
					}
				}
			}
		}
	}

	if err := q.Count(count).Error; err != nil {
		return ur, err
	}

	defer con.Close()
	return Lapkeu, nil
}