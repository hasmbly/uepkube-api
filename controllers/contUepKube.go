package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	 _"github.com/jinzhu/gorm/dialects/mysql"
	 "uepkube-api/models"
	 "uepkube-api/db"
	 "regexp"
	 "uepkube-api/helpers"
	 // "log"
	 // "fmt"
	 // "time"
)

// @Summary Get Uep (byNik) or Get Kube (byName)
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Param val query string true "string"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/uepkube [get]
func GetUepKube(c echo.Context) error {
	val 	:= c.QueryParam("val")
	Uep 	:= []models.ShowUep{}

	re := regexp.MustCompile("[0-9]+")
	errr := (re.FindAllString(val, -1))

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	/*query user*/
	if errr != nil {
		q1 := con
		q1 = q1.Table("tbl_user")
		q1 = q1.Where("nik like ?", "%"+val+"%")
		q1 = q1.Joins("join tbl_uep on tbl_uep.id_uep = tbl_user.id_user")
		q1 = q1.Joins("join tbl_usaha_produk on tbl_usaha_produk.id_uep = tbl_uep.id_uep")
		q1 = q1.Joins("join tbl_jenis_usaha on tbl_jenis_usaha.id_usaha = tbl_usaha_produk.id_usaha")		
		q1 = q1.Select("tbl_user.id_user, tbl_user.nik, tbl_user.nama, tbl_user.alamat, tbl_user.lat, tbl_user.lng, tbl_user.photo, tbl_uep.*, tbl_jenis_usaha.jenis_usaha")

		if err := q1.Scan(&Uep).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}	
	} else if errr == nil {
		return GetKube(c)
	}

	if len(Uep) != 0 {
		for i,_ := range Uep {
			if Uep[i].Photo != "" {
				ImageBlob := Uep[i].Photo
				Uep[i].Photo = "data:image/png;base64," + ImageBlob	
			}
		}
	}
	
	for i,_ := range Uep {
		Uep[i].Flag = "UEP"
	}	

	r := &models.Jn{Msg: Uep}

	defer con.Close()
	return c.JSON(http.StatusOK, r)
}

// @Summary GetPaginateProduk
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Param produk body models.PosPagin true "Show Produk UepKube"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/uepkube/produk [post]
func GetPaginateProdukUepKube(c echo.Context) (err error) {
	if err := helpers.PaginateProduk(c, &r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

// @Summary GetPaginatePelatihanUepKube
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Param produk body models.PosPagin true "Show Pelatihan UepKube"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/uepkube/pelatihan [post]
func GetPaginatePelatihanUepKube(c echo.Context) (err error) {
	if err := helpers.PaginatePelatihan(c, &r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
	// return nil
}

// @Summary GeAllBantuanPeriods
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/bantuan_periods [get]
func GeAllBantuanPeriods(c echo.Context) (err error) {
	Periods := []models.Tbl_bantuan_periods{}
	// timeFormat := "2006-01-02 15:04:05"

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	/*query user*/
	if err := con.Find(&Periods).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
	
	// var rPeriods []RPeriods

	// if len(Periods) != 0 {
	// 	for i,_ := range Periods {

	// 		rPeriods = append(rPeriods, Periods[i])

	// 		ts := Periods[i].Start_date
	// 		te := Periods[i].End_date

	// 		rPeriods[i].Start_date = ts.Format(timeFormat)
	// 		rPeriods[i].End_date = te.Format(timeFormat)

	// 		log.Println(Periods[i].Start_date)
	// 		// log.Println(Periods[i].End_date)
	// 	}
	// }

	r := &models.Jn{Msg: Periods}

	defer con.Close()
	return c.JSON(http.StatusOK, r)
}

// @Summary GeAllMemberPelatihan
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Param id_pendamping query int true "int"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/member_pelatihan [get]
func GeAllMemberPelatihan(c echo.Context) (err error) {
	id_pendamping 	:= c.QueryParam("id_pendamping")
	MemberUep 			:= []models.MemberPelatihan{}
	MemberKube 			:= []models.MemberPelatihan{}
	var Result []interface{}

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// get uep base on id_pendamping
	if err := con.Table("tbl_uep t1").Select("t2.id_user, t2.nama, 'UEP' as flag ").Joins("join tbl_user t2 on t2.id_user = t1.id_uep").Where("id_pendamping = ?", id_pendamping).Scan(&MemberUep).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	if len(MemberUep) != 0 {
		for i,_ := range MemberUep {
			Result = append(Result, MemberUep[i])
		}
	}

	// get kube base on id_pendamping
	// if err := con.Table("tbl_kube t1").Select("t1.id_kube as id, t1.nama_kube as nama, 'kube' as flag ").Where("id_pendamping = ?", id_pendamping).Scan(&MemberKube).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// if len(MemberKube) != 0 {
	// 	for i,_ := range MemberKube {
	// 		Result = append(Result, MemberKube[i])
	// 	}
	// }	

	// get member pendamping from uep
	var KubesMember = []string{"ketua", "sekertaris", "bendahara", "anggota1", "anggota2", "anggota3", "anggota4", "anggota5", "anggota6", "anggota7"}

	for i,_ := range KubesMember {
		if err := con.Table("tbl_kube t1").Select("t2.id_user, t2.nama, 'KUBE' as flag ").Joins("join tbl_user t2 on t2.id_user = t1." + KubesMember[i]).Where("id_pendamping = ?", id_pendamping).Scan(&MemberKube).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

		if len(MemberKube) != 0 {
			for i,_ := range MemberKube {
				Result = append(Result, MemberKube[i])
			}
		}
	}

	r := &models.Jn{Msg: Result}

	defer con.Close()
	return c.JSON(http.StatusOK, r)
}

// @Summary GetAllFaq
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/faq [get]
func GeAllFaq(c echo.Context) (err error) {
	Faq := []models.Tbl_faq{}

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	/*query user*/
	if err := con.Find(&Faq).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
	
	r := &models.Jn{Msg: Faq}

	defer con.Close()
	return c.JSON(http.StatusOK, r)
}

// @Summary GeAllPendamping
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/pendamping [get]
func GeAllPendamping(c echo.Context) (err error) {
	Pendampings := []models.CustomPendamping{}

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	/*query user*/
	if err := con.Table("tbl_pendamping t1").Select("t1.*, t2.nama").Joins("join tbl_user t2 on t2.id_user = t1.id_pendamping").Scan(&Pendampings).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
	
	// check how many pendamping in each uep/kube

	r := &models.Jn{Msg: Pendampings}

	defer con.Close()
	return c.JSON(http.StatusOK, r)
}

// @Summary GeAllJenisUsaha
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/jenis_usaha [get]
func GeAllJenisUsaha(c echo.Context) (err error) {
	JenisUsaha := []models.Tbl_jenis_usaha{}

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	/*query user*/
	if err := con.Find(&JenisUsaha).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
	
	r := &models.Jn{Msg: JenisUsaha}

	defer con.Close()
	return c.JSON(http.StatusOK, r)
}

// @Summary GeAllUser
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Param nik query string true "string"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/users [get]
func GeAllUser(c echo.Context) (err error) {
	nik 	:= c.QueryParam("nik")
	// User 	:= models.Tbl_user{}
	Users 	:= []models.Tbl_user{}

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	/*query user*/
	q := con
	q = q.Model(&Users)
	q = q.Preload("Kelurahan")
	q = q.Preload("Kecamatan")
	q = q.Preload("Kabupaten")
	q = q.Where("nik like ?", "%"+nik+"%")
	q = q.Find(&Users)
	
	// Chceck in Uep
	for i, _ := range Users {
		var id_uep []int
		q1 := con
		q1 = q1.Table("tbl_uep")
		q1 = q1.Where("id_uep = ?", Users[i].Id_user)
		q1 = q1.Pluck("id_uep", &id_uep)
		
		if len(id_uep) != 0 {
			Users[i].Flag = "UEP"
			continue
		}

		var id_kube_members []int
		var KubesMember = []string{"ketua", "sekertaris", "bendahara", "anggota1", "anggota2", "anggota3", "anggota4", "anggota5", "anggota6", "anggota7"}		
		for o, _ := range KubesMember {
			q2 := con
			q2 = q2.Table("tbl_kube")
			q2 = q2.Where(KubesMember[o] + " = ?", Users[i].Id_user)
			q2 = q2.Pluck(KubesMember[o], &id_kube_members)

			if len(id_kube_members) != 0 {
				Users[i].Flag = "KUBE"
			}
		}
	}

	r := &models.Jn{Msg: Users}

	defer con.Close()
	return c.JSON(http.StatusOK, r)
}

// @Summary GetAllAddress
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Param region query string true "string"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/address [get]
func GeAllAddress(c echo.Context) (err error) {
	region 	:= c.QueryParam("region")
	address := []models.View_address{}

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	/*query user*/
	if err := con.Model(&models.View_address{}).Where("region like ?", "%"+region+"%").Find(&address).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
	
	r := &models.Jn{Msg: address}

	defer con.Close()
	return c.JSON(http.StatusOK, r)
}

// @Summary GeAllMonevItems
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Param flag query string true "string"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/monev_items [get]
func GeAllMonevItems(c echo.Context) (err error) {
	flag 	:= c.QueryParam("flag")

	Dimensi := []models.Tbl_dimensi_uepkube{}
	// PertanyaanUep := []models.Tbl_indikator_uep{}
	// PertanyaanKube := []models.Tbl_indikator_kube{}

	r := &models.Jn{}

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if flag == "uep" {
		q := con
		q = q.Model(&Dimensi)
		q = q.Preload("Aspek_uep.Kriteria_uep.Indikator_uep")
		q = q.Find(&Dimensi)
	} else if flag == "kube" {
		q := con
		q = q.Model(&Dimensi)
		q = q.Preload("Aspek_kube.Kriteria_kube.Indikator_kube")
		q = q.Find(&Dimensi)
	}
	
	r.Msg = Dimensi

	defer con.Close()
	return c.JSON(http.StatusOK, r)
}


// @Summary GeAllUepKubeDetail
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/persebaran [get]
func GeAllUepKubeDetail(c echo.Context) (err error) {
	
	Uep 	:= []models.ShowUep{}
	Kube 	:= []models.Tbl_kube{}
	ShowKubes := models.ShowKube{}
	var tempo []interface{}	

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	/*query uep*/
	q1 := con
	q1 = q1.Table("tbl_user")
	q1 = q1.Joins("join tbl_uep on tbl_uep.id_uep = tbl_user.id_user")
	q1 = q1.Select("tbl_user.*")

	if err := q1.Find(&Uep).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	/*query kube*/
	if err := con.Find(&Kube).Error; gorm.IsRecordNotFoundError(err)  {
		return echo.NewHTTPError(http.StatusNotFound, "Kube Not Found")
	}		
		
	for i,_ := range Kube {
		
		helpers.SetMemberNameKube(&ShowKubes, Kube[i])

		tempo = append(tempo, ShowKubes)
	}	

	for i,_ := range Uep {
		Uep[i].Flag = "UEP"
		tempo = append(tempo, Uep[i])
	}

	r := &models.Jn{Msg: tempo}

	defer con.Close()
	return c.JSON(http.StatusOK, r)

}