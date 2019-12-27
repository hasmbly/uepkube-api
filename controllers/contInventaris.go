package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/db"
	"uepkube-api/models"
	"uepkube-api/helpers"
	"strconv"
	"fmt"
	"log"

	"bufio"
	"encoding/base64"	
	"io/ioutil"		
)

func GetInventaris(c echo.Context) error {
	id := c.QueryParam("id")

	For, _	:= strconv.Atoi(c.QueryParam("for"))
	var Field int = int(For) // flag uep : 0 | kube : 1

	var tmpPath, urlPath, blobFile, flag, host string
	flag = "INVENTORY"
	host = c.Request().Host

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	Inventory := models.Tbl_inventory{}


	if Field == 0 {
		q := con
		q = q.Model(&models.Tbl_inventory{})
		q = q.Preload("Photo")
		// q = q.Preload("Pendamping", func(q *gorm.DB) *gorm.DB {
		// 	return q.Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Select("tbl_pendamping.*,tbl_user.nama")
		// })
		q = q.Where("id_uep = ?", id)
		// q = q.First(&Inventory)
		if err := q.First(&Inventory).Error; gorm.IsRecordNotFoundError(err) {
			return echo.ErrNotFound
		} else if err != nil {
			return echo.ErrInternalServerError
		}		
	} else if Field == 1 {
		q := con
		q = q.Model(&models.Tbl_lapkeu{})
		q = q.Preload("Photo")
		q = q.Preload("Pendamping", func(q *gorm.DB) *gorm.DB {
			return q.Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Select("tbl_pendamping.*,tbl_user.nama")
		})
		q = q.Where("id_kube = ?", id)
		if err := q.First(&Inventory).Error; gorm.IsRecordNotFoundError(err) {
			return echo.ErrNotFound
		} else if err != nil {
			return echo.ErrInternalServerError
		}	
	}

	// photo
	for i, _ := range Inventory.Photo {
			id_photo := Inventory.Photo[i].Id

			tmpPath	= fmt.Sprintf(helpers.GoPath + "/src/uepkube-api/static/assets/images/%s_id_%s_photo_id_%d.png", flag,id,id_photo)
			urlPath	= fmt.Sprintf("http://%s/images/%s_id_%s_photo_id_%d.png", host,flag,id,id_photo)
			blobFile = Inventory.Photo[i].Files

			if check := CreateFile(tmpPath, blobFile); check == false {
				log.Println("blob is empty : ", check)
			}
				
			Inventory.Photo[i].Files = urlPath
	}

	// detail
	if Inventory.Id_uep != 0 {
		id := Inventory.Id_uep
		
		User 	:= Tbl_user{}
		q := con
		q = q.Model(&User)
		q = q.Joins("join tbl_uep on tbl_uep.id_uep = tbl_user.id_user")
		q = q.Select("tbl_uep.*, tbl_user.*")
		q = q.Preload("JenisUsaha")
		q = q.Preload("LapkeuHistory", func(q *gorm.DB) *gorm.DB {
			return q.Where("id_uep = ?", id)
		})			
		q = q.Preload("MonevHistory", func(q *gorm.DB) *gorm.DB {
			return q.Where("id_uep = ?", id)
		})	
		q = q.Preload("InventarisHistory", func(q *gorm.DB) *gorm.DB {
			return q.Where("id_uep = ?", id)
		})
		q = q.Preload("PelatihanHistory", func(q *gorm.DB) *gorm.DB {
			return q.Where("id_uep = ?", id)
		})
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

				tmpPath	= fmt.Sprintf(helpers.GoPath + "/src/uepkube-api/static/assets/images/%s_id_%d_photo_id_%d.png", flag,id,id_photo)
				urlPath	= fmt.Sprintf("http://%s/images/%s_id_%d_photo_id_%d.png", host,flag,id,id_photo)
				blobFile = User.Photo[index].Files

				if check := CreateFile(tmpPath, blobFile); check == false {
					log.Println("blob is empty : ", check)
				}
			
				User.Photo[index].Files = urlPath
		}

		Inventory.Detail = User

	} else if Inventory.Id_kube != 0 {

		id := Inventory.Id_kube
		
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

				tmpPath	= fmt.Sprintf(helpers.GoPath + "/src/uepkube-api/static/assets/images/%s_id_%d_photo_id_%d.png", flag,id,id_photo)
				urlPath	= fmt.Sprintf("http://%s/images/%s_id_%d_photo_id_%d.png", host,flag,id,id_photo)
				blobFile = Kube.Photo[index].Files

				if check := CreateFile(tmpPath, blobFile); check == false {
					log.Println("blob is empty : ", check)
				}
			
				Kube.Photo[index].Files = urlPath
		}
		Inventory.Detail = Kube
	}

	r := &models.Jn{Msg: Inventory}
	defer con.Close()

	return c.JSON(http.StatusOK, r)
}

func GetPaginateInventaris(c echo.Context) (err error) {	
	if err := helpers.PaginateInventory(c, &r); err != nil {
		return echo.ErrInternalServerError
	}	
	return c.JSON(http.StatusOK, r)
}

func AddInventaris(c echo.Context) (err error) {
	inventory := &models.Inventory{}
	
	if err := c.Bind(inventory); err != nil {
		return err
	}

	// get log post
	helpers.FetchPost(inventory)

	// validation
	if inventory.Id_uep == 0 && inventory.Id_kube == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id_uep or id_kube")
	}	

	if inventory.Id_pendamping == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id_pendamping")
	}	

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	creditDebit := &models.Tbl_inventory{}
	
	if inventory.Id_uep == 0 { 
		creditDebit.Id_kube = inventory.Id_kube
	}

	if inventory.Id_kube == 0 { 
		creditDebit.Id_uep = inventory.Id_uep
	}

	Inventory := &models.Tbl_inventory{}
	Inventory = inventory.Tbl_inventory

	if err := con.Create(&Inventory).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// store creditDebit
	// creditDebit.Credit = 1
	// if len(id_periods) != 0 { creditDebit.Id_periods = id_periods[0] }
	// creditDebit.Nilai = inventory.Harga
	// creditDebit.Deskripsi = fmt.Sprintf("Debit dengan nilai : Rp. %.2f,-", inventory.Harga)
	// if err := con.Create(&creditDebit).Error; err != nil { return echo.ErrInternalServerError }

	defer con.Close()

	r := &models.Jn{Msg: "Success Store Data"}
	return c.JSON(http.StatusOK, r)
}

func UpdateInventaris(c echo.Context) (err error) {
	// inventaris := &models.Tbl_inventaris{}

	// if err := c.Bind(inventaris); err != nil {
	// 	return err
	// }

	// if inventaris.Id_inventaris == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id")
	// }

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// if err := con.Model(&models.Tbl_inventaris{}).UpdateColumns(&inventaris).Error; err != nil {
	// 	return echo.ErrInternalServerError
	// }

	// if err := con.Table("tbl_inventaris").Where("id_inventaris = ?",inventaris.Id_inventaris).UpdateColumn("status", inventaris.Status).Error; err != nil {return echo.ErrInternalServerError}

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Update Data"}
	// return c.JSON(http.StatusOK, r)
	return nil	
}

func DeleteInventaris(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("id"))

	if id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	}

	inventory := &models.Tbl_inventory{}
	inventory.Id = id

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Delete(&inventory).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	defer con.Close()

	r := &models.Jn{Msg: "Success Delete Data"	}
	return c.JSON(http.StatusOK, r)	
}

func UploadInventarisFiles(c echo.Context) (err error) {
	// query
	id, _ 			:= strconv.Atoi(c.QueryParam("id"))
	is_display, _ 	:= strconv.Atoi(c.QueryParam("is_display"))

	// formValue
	description := c.FormValue("description")
	types 		:= c.FormValue("type")

	if id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	}

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["files"]

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	log.Println("files : ", len(files))

	for _,f := range files {

		src, err := f.Open()
		if err != nil {
			return err
		}
		defer src.Close()

	    // Read entire JPG into byte slice.
	    reader := bufio.NewReader(src)
	    content, _ := ioutil.ReadAll(reader)

	    // Encode as base64.
	    encoded := base64.StdEncoding.EncodeToString(content)
		
		// execute
		InventoryFiles := &models.Tbl_inventory_files{}

		InventoryFiles.Id_inventory = id
		InventoryFiles.Files = encoded
		InventoryFiles.Description = description
		InventoryFiles.Type = types
		InventoryFiles.Is_display = is_display

		if err := con.Create(&InventoryFiles).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
	}

	defer con.Close()

	log.Println("Uploads Inventory's file to id : ", id)
	r := &models.Jn{Msg: "Success Upload files"}	
	return c.JSON(http.StatusOK, r)	
	
}