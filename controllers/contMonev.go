package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/db"
	"uepkube-api/models"
	"uepkube-api/helpers"
	// "strconv"
	"fmt"
	"log"
	// "time"
)

func GetMonev(c echo.Context) error {
	id 		:= c.QueryParam("id")

	var tmpPath, urlPath, blobFile, flag, host string
	flag = "MONEV"
	host = c.Request().Host

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	Monev := models.Tbl_monev_final{}
	q := con
	q = q.Model(&Monev)
	q = q.Preload("Category")
	q = q.Preload("Pendamping")
	// q = q.First(&Monev, id)
	if err := q.First(&Monev, id).Error; gorm.IsRecordNotFoundError(err) {
		return echo.ErrNotFound
	} else if err != nil {
		return echo.ErrInternalServerError
	}
	// Detail Uep_kube

	// Data Monev
	MonevResult := []models.Data_monev{}
	var FieldId string
	var ValueId int
	var TblType string

	if Monev.Id_uep == 0 { 
		FieldId = "id_kube"
		TblType = "_kube"
		ValueId = Monev.Id_kube 
	}
	
	if Monev.Id_kube == 0 {
		FieldId = "id_uep"
		TblType = "_uep"	
		ValueId = Monev.Id_uep 
	}
	
	q2 := con
	q2 = q2.Table("tbl_monev_calculate t1")
	q2 = q2.Select("t1.skor_total, t2.skor_indikator, t3.bobot, t4.nama_aspek, t5.nama_dimensi")
	q2 = q2.Joins("join tbl_indikator" + TblType + " t2 on t2.id_indikator = t1.id_indikator")
	q2 = q2.Joins("join tbl_kriteria" + TblType + " t3 on t3.id_kriteria = t2.id_kriteria")
	q2 = q2.Joins("join tbl_aspek" + TblType + " t4 on t4.id_aspek = t3.id_aspek")
	q2 = q2.Joins("join tbl_dimensi_uepkube t5 on t5.id_dimensi = t4.id_dimensi")
	q2 = q2.Where(FieldId + " = ?", ValueId)
	q2 = q2.Scan(&MonevResult)

	// data_monev
	for i,_ := range MonevResult {
		Monev.Data_monev = append(Monev.Data_monev, &MonevResult[i])
	}

	// detail
	if Monev.Id_uep != 0 {
		id := Monev.Id_uep
		
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

		Monev.Detail = User

	} else if Monev.Id_kube != 0 {

		id := Monev.Id_kube
		
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
		Monev.Detail = Kube
	}

	r := &models.Jn{Msg: Monev}
	defer con.Close()

	return c.JSON(http.StatusOK, r)
}

func GetPaginateMonev(c echo.Context) (err error) {	
	if err := helpers.PaginateMonev(c, &rMonev); err != nil {
		return echo.ErrInternalServerError
	}	
	return c.JSON(http.StatusOK, rMonev)
	// return nil
}

func GetPaginatePkt(c echo.Context) (err error) {	
	if err := helpers.PaginatePkt(c, &r); err != nil {
		return echo.ErrInternalServerError
	}	
	return c.JSON(http.StatusOK, r)
}

func AddMonev(c echo.Context) (err error) {
	// monev := &models.Monev{}
	// MonevFinal := &models.Tbl_monev_final{}	

	// if err := c.Bind(monev); err != nil {
	// 	return err
	// }

	// // get log post
	// helpers.FetchPost(monev)	

	// // validation
	// if monev.Id_uep == 0 && monev.Id_kube == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id_uep or id_kube")
	// }	

	// if monev.Id_pendamping == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id_pendamping")
	// }

	// if len(monev.Id_indikator) == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id_indikator")
	// }	

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// var TblType string
	// // var FieldId string
	// // var ValueId int

	// if monev.Id_kube == 0 { 
	// 	// FieldId = "id_uep"
	// 	// ValueId = monev.Id_uep 
	// 	TblType = "_uep"
	// 	MonevFinal.Id_uep = monev.Id_uep
	// 	// get id_monev_final
	// 	var id []int
	// 	if err := con.Table("tbl_monev_final").Where("id_uep = ?", monev.Id_uep).Pluck("id", &id).Error; err != nil {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, err)
	// 	}
	// 	MonevFinal.Id = id[0]
	// 	MonevFinal.Flag = "UEP"
	// }
	// if monev.Id_uep == 0 { 
	// 	// FieldId = "id_kube"
	// 	// ValueId = monev.Id_kube 
	// 	MonevFinal.Id_kube = monev.Id_kube	
	// 	// get id_monev_final
	// 	var id []int
	// 	if err := con.Table("tbl_monev_final").Where("id_kube = ?", monev.Id_kube).Pluck("id", &id).Error; err != nil {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, err)
	// 	}
	// 	MonevFinal.Id = id[0]		
	// 	MonevFinal.Flag = "KUBE"
	// 	TblType = "_kube"
	// }

	// var IdKriteria []int
	
	// var bobot []int
	// var skor []int

	// var TotalSkor []int
	// var SumTotal int

	// // get skor & bobot from indikator
	// for i, _ := range monev.Id_indikator {

	// 	var bobotKriteria []int
	// 	var skorIndikator []int

	// 	con.Table("tbl_indikator" + TblType).Where("id_indikator = ?", monev.Id_indikator[i]).Pluck("skor_indikator", &skorIndikator)

	// 	if len(skorIndikator) != 0 {
	// 		skor = append(skor, skorIndikator[0])
	// 	} else {
	// 		return echo.NewHTTPError(http.StatusBadRequest, "Maaf Indikator tidak ditemukan")
	// 	}

	// 	con.Table("tbl_indikator" + TblType).Where("id_indikator = ?", monev.Id_indikator[i]).Pluck("id_kriteria", &IdKriteria)

	// 	if len(IdKriteria) != 0 {
	// 		con.Table("tbl_kriteria" + TblType).Where("id_kriteria = ?", IdKriteria[0]).Pluck("bobot", &bobotKriteria)

	// 		if len(bobotKriteria) != 0 {
	// 			bobot = append(bobot, bobotKriteria[0])
	// 		}
			
	// 	} else {
	// 		return echo.NewHTTPError(http.StatusBadRequest, "Maaf Kriteria tidak ditemukan")	
	// 	}



	// }

	// log.Println("skor : ", skor)
	// log.Println("bobot : ", bobot)

	// // calculate total_skor
	// for i, _ := range skor {
	// 	total := skor[i] * bobot[i]
	// 	TotalSkor = append(TotalSkor, total)
	// }

	// for x, _ := range TotalSkor {
		
	// 	Monev := &models.Tbl_monev_calculate{}	
	// 	if monev.Id_kube == 0 { Monev.Id_uep = monev.Id_uep }
	// 	if monev.Id_uep == 0 { Monev.Id_kube = monev.Id_kube }
	// 	Monev.Id_indikator = monev.Id_indikator[x]
	// 	Monev.Skor_total = TotalSkor[x]

	// 	log.Println("MonevResult : ", Monev)
	// 	// store total_skor to Tbl_monev_result_uepkube
	// 	if err := con.Create(&Monev).Error; err != nil {
	// 		return echo.ErrInternalServerError
	// 	}
	// }

	// // calculate sum_total
	// for i, _ := range TotalSkor {
	// 	SumTotal = TotalSkor[i] + SumTotal
	// }

	// // store final monev
	// log.Println("SumTotal : ", SumTotal)
	// MonevFinal.Sum_total = SumTotal
	// var id_category int
	// var desc_category string
	
	// if SumTotal >= 301 && SumTotal <= 400 {
	// 	id_category = 1
	// 	desc_category = "Sangat Berhasil"
	// } else if SumTotal >= 201 && SumTotal <= 300 {
	// 	id_category = 2
	// 	desc_category = "Berhasil"
	// } else if SumTotal >= 100 && SumTotal <= 200 {
	// 	id_category = 3
	// 	desc_category = "Kurang Berhasil"
	// } else {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Maaf Hasil Monev anda Belum berhasil, Total Monev : " + strconv.Itoa(SumTotal))
	// }
	
	// MonevFinal.Id_category = id_category
	// MonevFinal.Id_pendamping = monev.Id_pendamping
	// MonevFinal.Is_monev = "SUDAH"
	
	// // MonevFinal.Id_periods = id_periods[0]

	// log.Println("monevFInal : ", MonevFinal)

	// // in the utils function create thread when app start to check current for adding monev queue

	// // check if already do monev in the past
	// var created_at []string
	// var CurrYear string
	// year, _ , _ := time.Now().Date()
	// CurrYear = strconv.Itoa(year)

	// if MonevFinal.Id_uep != 0 {
	// 	var flag []string
	// 	if err := con.Table("tbl_monev_final").Where("created_at like ?". "%"+CurrYear+"%").Where("id_uep = ?", MonevFinal.Id_uep).Pluck("flag", &flag).Error; err != nil {
	// 		return echo.NewHTTPError(http.StatusInternalServerError, err)
	// 	}
	// 	if flag[0] == "SUDAH" {
	// 		// check if already do monev this year


	// 		if err := con.Table("tbl_monev_final").Where("created_at like ?", "%"+CurrYear+"%").Pluck("created_at", &created_at).Error; err != nil {
	// 			return echo.NewHTTPError(http.StatusInternalServerError, err)
	// 		}

	// 		if len(created_at) != 0 {
	// 			return echo.NewHTTPError(http.StatusBadRequest, "Maaf sudah pernah dilakukan Monev tahun ini")
	// 		} else if len(created_at) == 0 {

	// 			// if err := con.Create(&MonevFinal).Error; err != nil {
	// 			// 	return echo.NewHTTPError(http.StatusInternalServerError, err)
	// 			// }
	// 		}

	// 	} else if flag[0] == "BELUM" {
	// 		log.Println("Monev : ", MonevFinal)
	// 		log.Println("Monev : ", flag[0])

	// 	}
	// }

	// 	// update if exist : BELUM
	// 	if err := con.Model(&MonevFinal).UpdateColumns(&MonevFinal).Error; err != nil {
	// 		return echo.ErrInternalServerError
	// 	}

	// 	// store new
	// 	// if err := con.Create(&MonevFinal).Error; err != nil {
	// 	// 	return echo.NewHTTPError(http.StatusInternalServerError, err)
	// 	// }				

	// // store final monev
	// // if err := con.Create(&MonevFinal).Error; err != nil {return echo.ErrInternalServerError}
	// // if err := con.Save(&MonevFinal).Error; err != nil {return echo.ErrInternalServerError}
	// // if err := con.Model(&MonevFinal).UpdateColumns(&MonevFinal).Error; err != nil {
	// // 	return echo.ErrInternalServerError
	// // }

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Store Data, Total Skor Monev : " + strconv.Itoa(SumTotal) + ", Kategori : " + desc_category}
	// return c.JSON(http.StatusOK, r)
	return nil
}

func UpdateMonev(c echo.Context) (err error) {
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

func DeleteMonev(c echo.Context) (err error) {
	// id, _ := strconv.Atoi(c.Param("id"))

	// if id == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	// }

	// monev := &models.Tbl_monev{}
	// monev.Id = id

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// if err := con.Delete(&monev).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Delete Data"	}
	// return c.JSON(http.StatusOK, r)	
	return nil
}