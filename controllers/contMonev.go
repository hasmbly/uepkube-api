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
)

func GetMonev(c echo.Context) error {
	id 		:= c.QueryParam("id")

	var tmpPath, urlPath, blobFile, flag, host string
	flag = "MONEV"
	host = c.Request().Host

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	Monev := models.Tbl_monev_uepkube{}
	q := con
	q = q.Model(&Monev)
	q = q.Preload("Category")
	q = q.Preload("Pendamping")
	q = q.First(&Monev, id)

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
	q2 = q2.Table("tbl_monev_result_uepkube t1")
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
	if err := helpers.PaginateMonev(c, &r); err != nil {
		return echo.ErrInternalServerError
	}	
	return c.JSON(http.StatusOK, r)
	// return nil
}

func AddMonev(c echo.Context) (err error) {
	monev := &models.Monev{}
	MonevFinal := &models.Tbl_monev_uepkube{}	

	if err := c.Bind(monev); err != nil {
		return err
	}

	// get log post
	helpers.FetchPost(monev)	

	// validation
	if monev.Id_uep == 0 && monev.Id_kube == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id_uep or id_kube")
	}	

	if monev.Id_pendamping == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id_pendamping")
	}

	if len(monev.Id_indikator) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id_indikator")
	}	

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	var TblType string
	var FieldId string
	var ValueId int

	if monev.Id_kube == 0 { 
		FieldId = "id_uep"
		TblType = "_kube"
		MonevFinal.Id_uep = monev.Id_uep	
		ValueId = monev.Id_uep 
	}
	if monev.Id_uep == 0 { 
		FieldId = "id_kube"
		MonevFinal.Id_kube = monev.Id_kube	
		TblType = "_uep"
		ValueId = monev.Id_kube 
	}

	var IdKriteria []int
	
	var bobot []int
	var skor []int

	var TotalSkor []int
	var SumTotal int

	// get skor & bobot from indikator
	for i, _ := range monev.Id_indikator {

		var bobotKriteria []int
		var skorIndikator []int

		con.Table("tbl_indikator" + TblType).Where("id_indikator = ?", monev.Id_indikator[i]).Pluck("skor_indikator", &skorIndikator)
		skor = append(skor, skorIndikator[0])

		con.Table("tbl_indikator" + TblType).Where("id_indikator = ?", monev.Id_indikator[i]).Pluck("id_kriteria", &IdKriteria)
		con.Table("tbl_kriteria" + TblType).Where("id_kriteria = ?", IdKriteria[0]).Pluck("bobot", &bobotKriteria)
		bobot = append(bobot, bobotKriteria[0])

	}

	log.Println("skor : ", skor)
	log.Println("bobot : ", bobot)

	// calculate total_skor
	for i, _ := range skor {
		total := skor[i] * bobot[i]
		TotalSkor = append(TotalSkor, total)
	}

	for i, _ := range TotalSkor {
		

		Monev := &models.Tbl_monev_result_uepkube{}	
		if monev.Id_kube == 0 { Monev.Id_uep = monev.Id_uep }
		if monev.Id_uep == 0 { Monev.Id_kube = monev.Id_kube }
		Monev.Id_indikator = monev.Id_indikator[i]
		Monev.Skor_total = TotalSkor[i]

		log.Println("MonevResult : ", Monev)
		// store total_skor to Tbl_monev_result_uepkube
		if err := con.Create(&Monev).Error; err != nil {return echo.ErrInternalServerError}		
	}

	// calculate sum_total
	for i, _ := range TotalSkor {
		SumTotal = TotalSkor[i] + SumTotal
	}

	// store final monev
	log.Println("SumTotal : ", SumTotal)
	MonevFinal.Sum_total = SumTotal
	var id_category int
	var desc_category string
	
	if SumTotal >= 301 && SumTotal <= 400 {
		id_category = 1
		desc_category = "Sangat Berhasil"
	} else if SumTotal >= 201 && SumTotal <= 300 {
		id_category = 2
		desc_category = "Berhasil"
	} else if SumTotal >= 100 && SumTotal <= 200 {
		id_category = 3
		desc_category = "Kurang Berhasil"
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, "Maaf Hasil Monev anda Belum berhasil, Total Monev : " + strconv.Itoa(SumTotal))
	}
	
	MonevFinal.Id_category = id_category
	MonevFinal.Id_pendamping = monev.Id_pendamping
	MonevFinal.Is_monev = "SUDAH"

	// get id_periods
	var id_periods []int
	con.Table("tbl_periods_uepkube").Where(FieldId + " = ?", ValueId).Pluck("id_periods", &id_periods)
	
	MonevFinal.Id_periods = id_periods[0]

	log.Println("monevFInal : ", MonevFinal)

	// store final monev
	if err := con.Create(&MonevFinal).Error; err != nil {return echo.ErrInternalServerError}

	defer con.Close()

	r := &models.Jn{Msg: "Success Store Data, Total Skor Monev : " + strconv.Itoa(SumTotal) + ", Kategori : " + desc_category}
	return c.JSON(http.StatusOK, r)
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