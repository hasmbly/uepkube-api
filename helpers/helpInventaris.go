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

func PaginateInventory(c echo.Context, r *models.ResPagin) (err error) {
	u := &models.PosPagin{}
	num := 1

	// GetLoggedUser(c,"roles")

	if err := c.Bind(u); err != nil {
		return err
	}

	var co int = (u.Page - num) * u.Size
	
	PaginateResult, _ := ExecPaginateInventory(u,co,&CountRows)

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

func ExecPaginateInventory(f *models.PosPagin, offset int, count *int64) (ur []models.Tbl_inventory, err error) {

	Inventory := []models.Tbl_inventory{}

	con, err := db.CreateCon()
	if err != nil { return ur, echo.ErrInternalServerError }
	con.SingularTable(true)	

	q := con
	q = q.Model(&Inventory)
	q = q.Limit(int(f.Size))
	q = q.Offset(int(offset))

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
	
	q = q.Find(&Inventory)
	q = q.Limit(-1)
	q = q.Offset(-1)

	//getDetail
	if len(Inventory) != 0 {

		for i, _ := range Inventory {

			if Inventory[i].Id_uep != 0 {
				id := Inventory[i].Id_uep
				
				User 	:= Tbl_user{}
				q := con
				q = q.Model(&User)
				q = q.Joins("join tbl_uep on tbl_uep.id_uep = tbl_user.id_user")
				q = q.Select("tbl_uep.*, tbl_user.*")
				q = q.Preload("JenisUsaha")
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
				
				Inventory[i].Detail = User

			} else if Inventory[i].Id_kube != 0 {

				id := Inventory[i].Id_kube
				
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

				Inventory[i].Detail = Kube

			}
		}
	}

	// get photos
	// if len(Inventory) != 0 {
	// 	for i,_ := range Inventory {
	// 		var pelatihan_photos []models.Tbl_pelatihan_files
	// 		// var account = models.Tbl_account{}

	// 		con.Table("tbl_inventory_files").Where("type = 'IMAGE' ").Where("id = ?", Inventory[i].Id_pelatihan).Select("tbl_inventory_files.*").Find(&pelatihan_photos)

	// 		for i,_ := range pelatihan_photos {
	// 			ImageBlob := pelatihan_photos[i].Files
	// 			pelatihan_photos[i].Files = "data:image/png;base64," + ImageBlob			
	// 		}
	// 		Inventory[i].Photo = pelatihan_photos
	// 	}
	// }
	
	// get photos
	// if len(Inventory) != 0 {
	// 	for i,_ := range Inventory {
	// 		var name []string
	// 		if Inventory[i].Id_uep != 0 {
	// 			con.Table("tbl_user t1").Where("t1.id_user = ?", Inventory[i].Id_uep).Select("t1.nama,t2.jenis_usaha,t3.bantuan_modal,t1.created_at").Joins("join tbl_uep t4 on t4.id_uep = t1.id_user").Joins("join tbl_jenis_usaha t2 on t2.id_usaha = t4.id_jenis_usaha").Joins("join tbl_periods_uepkube t5 on t5.id_uep = t4.id_uep").Joins("join tbl_bantuan_periods t3 on t3.id = t5.id_periods").Scan(&Inventory[i])
	// 			// Inventory[i].NamaUepKube = "UEP " + name[0]
	// 		} else if Inventory[i].Id_kube != 0 {
	// 			con.Table("tbl_kube").Where("id_kube = ?", Inventory[i].Id_kube).Pluck("nama_kube", &name)
	// 			Inventory[i].NamaUepKube = name[0]
	// 		}
	// 	}
	// }

	if err := q.Count(count).Error; err != nil {
		return ur, err
	}

	// log.Println("result : ", Inventory)

	defer con.Close()
	return Inventory, nil
}