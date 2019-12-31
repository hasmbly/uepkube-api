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
	"time"
	"strconv"
)

type Tbl_pendamping struct {
	*models.Tbl_pendamping
	Nama string `json:"nama"`
}

type Tbl_user struct {
	*models.Tbl_user
	*models.Tbl_uep
}

var SudahMonev, BelumMonev int

func PaginateMonev(c echo.Context, r *models.ResPaginMonev) (err error) {
	flag = "MONEV"
	host = c.Request().Host

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

	*r = models.ResPaginMonev{
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
		SudahMonev:SudahMonev,
		BelumMonev:BelumMonev,
	}
	return err
}

func ExecPaginateMonev(f *models.PosPagin, offset int, count *int64) (ur []models.Tbl_monev_final, err error) {

	Monevs := []models.Tbl_monev_final{}

	con, err := db.CreateCon()
	if err != nil { return ur, echo.ErrInternalServerError }
	con.SingularTable(true)	

	q := con
	q = q.Model(&Monevs)
	q = q.Limit(int(f.Size))
	q = q.Offset(int(offset))
	q = q.Preload("Category")
	q = q.Preload("Pendamping", func(q *gorm.DB) *gorm.DB {
		return q.Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Select("tbl_pendamping.*,tbl_user.nama")
	})
	// q = q.Preload("Periods")

	for i,_ := range f.Filters {
		k := f.Filters[i].Key
		o := f.Filters[i].Operation
		v := f.Filters[i].Value

		if k == "periods" {
			if o == ":" {
				var id_periods []int
				if v != "" {
					 con.Table("tbl_bantuan_periods").Where("start_date like ?", "%"+v+"%").Pluck("id", &id_periods)					
					} else if v == "" {
						var CurrYear string
						year, _ , _ := time.Now().Date()
						CurrYear = strconv.Itoa(year)
						log.Println("currYear : ", CurrYear)
					 	con.Table("tbl_bantuan_periods").Where("start_date like ?", "%"+ CurrYear +"%").Pluck("id", &id_periods)						
					}

				 if len(id_periods) != 0 {
					q = q.Where("id_periods = ?", id_periods[0])
				 } else {
					// q = q.Where("id_periods = ?", id_periods[0])
					continue
				 }
			} else {
				continue
			}
		}

		if o == "LIKE" || o == "like" {
			if k == "periods" { continue }
			if v == "" { continue }
			q = q.Where(fmt.Sprintf("%s %s",k,o) + "?", "%"+v+"%")
		} else if o == ":" {
			if k == "periods" { continue }
			if v == "" {
				continue 
			} else {	
			 	q = q.Where(fmt.Sprintf("%s ",k) + "=" + "?", v) 
			}
		}
	}
	q = q.Order(fmt.Sprintf("%s %s",f.SortField,f.SortOrder))	
	
	q = q.Find(&Monevs)
	q = q.Limit(-1)
	q = q.Offset(-1)

	//getDetail
	if len(Monevs) != 0 {

		for i, _ := range Monevs {

			if Monevs[i].Is_monev == "BELUM" { BelumMonev = BelumMonev+1 }
			if Monevs[i].Is_monev == "SUDAH" { SudahMonev = SudahMonev+1 }

			if Monevs[i].Id_uep != 0 {
				id := Monevs[i].Id_uep
				
				User 	:= Tbl_user{}
				q := con
				q = q.Model(&User)
				q = q.Joins("join tbl_uep on tbl_uep.id_uep = tbl_user.id_user")
				q = q.Select("tbl_uep.*, tbl_user.*")
				q = q.Preload("Region")
				q = q.Preload("JenisUsaha")
				q = q.Preload("LapkeuHistory")
				q = q.Preload("MonevHistory")
				q = q.Preload("InventarisHistory")
				q = q.Preload("PelatihanHistory")
				q = q.Preload("Region")
				q = q.Preload("Pendamping", func(q *gorm.DB) *gorm.DB {
					return q.Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Select("tbl_pendamping.*,tbl_user.nama")
				})
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

				Monevs[i].Detail = User

			} else if Monevs[i].Id_kube != 0 {

				id := Monevs[i].Id_kube
				
				Kube 	:= models.Tbl_kube{}
				q := con
				q = q.Model(&Kube)
				q = q.Preload("JenisUsaha")
				q = q.Preload("LapkeuHistory", func(q *gorm.DB) *gorm.DB {
					return q.Where("id_kube = ?", id)
				})
				q = q.Preload("MonevHistory", func(q *gorm.DB) *gorm.DB {
					return q.Where("id_kube = ?", id)
				})	
				q = q.Preload("InventarisHistory", func(q *gorm.DB) *gorm.DB {
					return q.Where("id_kube = ?", id)
				})
				q = q.Preload("PelatihanHistory", func(q *gorm.DB) *gorm.DB {
					return q.Where("id_kube = ?", id)
				})	
				q = q.Preload("Pendamping", func(q *gorm.DB) *gorm.DB {
					return q.Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Select("tbl_pendamping.*,tbl_user.nama")
				})
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

				// Region
				// get region ketua
				var id_kelurahan []string
				con.Table("tbl_user").Where("id_user = ?", Kube.Ketua).Pluck("id_kelurahan", &id_kelurahan)
				// get region from view_address
				con.Table("view_address").Where("id_kelurahan = ?", id_kelurahan[0]).Scan(&Kube.Region)
				
				Monevs[i].Detail = Kube
			}
		}
	}

	if err := q.Count(count).Error; err != nil {
		return ur, err
	}

	// log.Println("result : ", Monevs)

	defer con.Close()
	return Monevs, nil
}	

func PaginatePkt(c echo.Context, r *models.ResPagin) (err error) {

	u := &models.PosPagin{}
	num := 1

	// GetLoggedUser(c,"roles")

	if err := c.Bind(u); err != nil {
		return err
	}

	var co int = (u.Page - num) * u.Size
	
	PaginateResult, _ := ExecPaginatePkt(u,co,&CountRows)

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

func ExecPaginatePkt(f *models.PosPagin, offset int, count *int64) (ur []models.Tbl_monev_final, err error) {

	Pkt := []models.Tbl_monev_final{}

	con, err := db.CreateCon()
	if err != nil { return ur, echo.ErrInternalServerError }
	con.SingularTable(true)	

	q := con
	q = q.Model(&Pkt)
	q = q.Preload("Category")
	q = q.Preload("Pendamping", func(q *gorm.DB) *gorm.DB {
		return q.Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Select("tbl_pendamping.*,tbl_user.nama")
	})

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

	q = q.Where("id_category = ?", 1)
	q = q.Order(fmt.Sprintf("%s %s",f.SortField,f.SortOrder))	
	q = q.Find(&Pkt)
	q = q.Limit(-1)
	q = q.Offset(-1)	

	if len(Pkt) != 0 {

		for i, _ := range Pkt {

			if Pkt[i].Is_monev == "BELUM" { BelumMonev = BelumMonev+1 }
			if Pkt[i].Is_monev == "SUDAH" { SudahMonev = SudahMonev+1 }

			if Pkt[i].Id_uep != 0 {
				id := Pkt[i].Id_uep
				
				User 	:= Tbl_user{}
				q := con
				q = q.Model(&User)
				q = q.Joins("join tbl_uep on tbl_uep.id_uep = tbl_user.id_user")
				q = q.Select("tbl_uep.*, tbl_user.*")
				q = q.Preload("Region")
				q = q.Preload("JenisUsaha")
				q = q.Preload("LapkeuHistory")
				q = q.Preload("MonevHistory")
				q = q.Preload("InventarisHistory")
				q = q.Preload("PelatihanHistory")
				q = q.Preload("Region")
				q = q.Preload("Pendamping", func(q *gorm.DB) *gorm.DB {
					return q.Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Select("tbl_pendamping.*,tbl_user.nama")
				})
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

				Pkt[i].Detail = User

			} else if Pkt[i].Id_kube != 0 {

				id := Pkt[i].Id_kube
				
				Kube 	:= models.Tbl_kube{}
				q := con
				q = q.Model(&Kube)
				q = q.Preload("JenisUsaha")
				q = q.Preload("LapkeuHistory", func(q *gorm.DB) *gorm.DB {
					return q.Where("id_kube = ?", id)
				})
				q = q.Preload("MonevHistory", func(q *gorm.DB) *gorm.DB {
					return q.Where("id_kube = ?", id)
				})	
				q = q.Preload("InventarisHistory", func(q *gorm.DB) *gorm.DB {
					return q.Where("id_kube = ?", id)
				})
				q = q.Preload("PelatihanHistory", func(q *gorm.DB) *gorm.DB {
					return q.Where("id_kube = ?", id)
				})	
				q = q.Preload("Pendamping", func(q *gorm.DB) *gorm.DB {
					return q.Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Select("tbl_pendamping.*,tbl_user.nama")
				})
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

				// Region
				// get region ketua
				var id_kelurahan []string
				con.Table("tbl_user").Where("id_user = ?", Kube.Ketua).Pluck("id_kelurahan", &id_kelurahan)
				// get region from view_address
				con.Table("view_address").Where("id_kelurahan = ?", id_kelurahan[0]).Scan(&Kube.Region)
				
				Pkt[i].Detail = Kube
			}
		}
	}	

	if err := q.Count(count).Error; err != nil {
		return ur, err
	}

	defer con.Close()
	return Pkt, nil
}