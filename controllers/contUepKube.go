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
	 "log"
	 // "strconv"
	 // "fmt"
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
		q1 = q1.Joins("join tbl_jenis_usaha on tbl_jenis_usaha.id_usaha = tbl_uep.id_jenis_usaha")		
		q1 = q1.Select("tbl_user.id_user, tbl_user.nik, tbl_user.nama, tbl_user.alamat, tbl_user.lat, tbl_user.lng, tbl_uep.*, tbl_jenis_usaha.jenis_usaha")

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
	// if err := helpers.PaginateProduk(c, &r); err != nil {
	// 	return err
	// }
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

type ChartDashBoard struct {
	HasilMonev *HasilMonev `json:"hasil_monev"`
	Persebaran *Persebaran `json:"persebaran"`
	JenisUsaha *JenisUsaha `json:"jenis_usaha"`
}

type Labels struct {
	Labels 	string 	`json:"labels"`
	Data 	[]int 	`json:"data"`
}

type HasilMonev struct {
	Labels 	[]string 		`json:"labels"`
	Uep 	interface{} 	`json:"uep"`
	Kube 	interface{}		`json:"kube"`
}

type Persebaran struct {
	Labels 	[]string 		`json:"labels"`
	Uep 	interface{}		`json:"uep"`
	Kube 	interface{}		`json:"kube"`
}

type JenisUsaha struct {
	Labels 	[]string 		`json:"labels"`
	Uep 	interface{}		`json:"uep"`
	Kube 	interface{}		`json:"kube"`
}

// @Summary GetChartDasboard
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/chart_dashboard [get]
func GetChartDasboard(c echo.Context) (err error) {
	// AllHasilMonev := []HasilMonev{}
	ChartDashBoard := &ChartDashBoard{}

	HasilMonev 	:= &HasilMonev{}
	Persebaran 	:= &Persebaran{}
	JenisUsaha 	:= &JenisUsaha{}

	For 		:= []string{"UEP", "KUBE"}		

	var (
		yearsM []string
	)

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// HASIL_MONEV
	LabelsM := []string{"Sangat Berhasil", "Cukup Berhasil", "Kurang Berhasil"}
	for i, _ := range LabelsM {
		HasilMonev.Labels = append(HasilMonev.Labels, LabelsM[i])
	}

	// PERSEBARAN
	LabelsP := []string{"Jakarta Pusat", "Jakarta Utara", "Jakarta Barat", "Jakarta Selatan", "Jakarta Timur", "Kep. Seribu"}
	LabelP_Id := []string{"3171", "3172", "3173", "3174", "3175", "3101"}
	for i, _ := range LabelsP {
		Persebaran.Labels = append(Persebaran.Labels, LabelsP[i])
	}

	// JENIS_USAHA
	var LabelsJ []string
	var LabelsJ_Id []string
	// get jenis_usaha
	jUsaha := models.Tbl_jenis_usaha{}
	if err := con.Find(&jUsaha).Pluck("jenis_usaha", &LabelsJ).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	// get id_jenis_usaha
	if err := con.Find(&jUsaha).Pluck("id_usaha", &LabelsJ_Id).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	for i, _ := range LabelsJ {
		JenisUsaha.Labels = append(JenisUsaha.Labels, LabelsJ[i])
	}

	// get all years for HASIL_MONEV
	// MonevYears := make(map[string]interface{})
	if err := con.Model(&models.Tbl_monev_final{}).Where("flag = ?", "UEP").Pluck("distinct YEAR(created_at) as created_at", &yearsM).Group("created_at").Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// get all years for PERSEBARAN UEP
	var yearsP_UEP []string 
	if err := con.Model(&models.Tbl_uep{}).Pluck("distinct YEAR(created_at) as created_at", &yearsP_UEP).Group("created_at").Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	// get all years for PERSEBARAN KUBE
	var yearsP_KUBE []string 
	if err := con.Model(&models.Tbl_kube{}).Pluck("distinct YEAR(created_at) as created_at", &yearsP_KUBE).Group("created_at").Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	log.Println("yearsP_KUBE : ", yearsP_KUBE)

	for f, _ := range For {
		var ResultM []interface{}

		// HASIL_MONEV per Tahun
		for x, _ := range yearsM {
			Labels 		:= Labels{}
			var dataM []int

			log.Println("years : ", yearsM[x])
			for i := 1; i <= len(LabelsM); i++ {

				var id_category []int
				if err := con.Model(&models.Tbl_monev_final{}).Where("flag = ?", For[f]).Where("created_at like ?", "%"+ yearsM[x] +"%").Where("id_category = ?", i).Pluck("id_category", &id_category).Error; err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
				dataM = append(dataM, len(id_category))
			}

			Labels.Labels 	= yearsM[x]
			Labels.Data 	= dataM

			ResultM = append(ResultM, Labels)
		}

		if For[f] == "UEP" { 
			HasilMonev.Uep = ResultM 
		}
		if For[f] == "KUBE" { 
			HasilMonev.Kube = ResultM 
		}

		// PERSEBEARAN_MONEV UEP
		for x, _ := range yearsP_UEP {
			var ResultP []interface{}
			var dataP []int
			
			var ResultJ []interface{}
			var dataJ []int
			
			Labels 		:= Labels{}

			// count persebaran UEP by id_kabupaten
			for i,_ := range LabelP_Id {
				var count int
				if err := con.Table("tbl_user").Where("id_kabupaten = ?", LabelP_Id[i]).Count(&count).Error; err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
				dataP = append(dataP, count)
			}

			Labels.Labels 	= yearsP_UEP[x]
			Labels.Data 	= dataP

			ResultP = append(ResultP, Labels)
			Persebaran.Uep = ResultP

			// count jenis_usaha UEP by id_jenis_usaha
			for i,_ := range LabelsJ_Id {
				var count int
				if err := con.Table("tbl_uep").Where("id_jenis_usaha = ?", LabelsJ_Id[i]).Count(&count).Error; err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
				dataJ = append(dataJ, count)
			}		

			Labels.Labels 	= yearsP_UEP[x]
			Labels.Data 	= dataJ

			ResultJ = append(ResultJ, Labels)
			JenisUsaha.Uep = ResultJ
		}

		// PERSEBEARAN_MONEV KUBE
		for x, _ := range yearsP_KUBE {
			var ResultP []interface{}
			var dataP []int
			
			var ResultJ []interface{}
			var dataJ []int			

			Labels 		:= Labels{}

			// count persebaran KUBE by id_kabupaten
			for i,_ := range LabelP_Id {
				var count int
				// var id_ketua []int
				if err := con.Table("tbl_kube t1").Select("t1.id_ketua,t2.id_kabupaten").Joins("join tbl_user t2 on t2.id_user = t1.ketua").Where("id_kabupaten = ?", LabelP_Id[i]).Count(&count).Error; err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
				dataP = append(dataP, count)
			}

			Labels.Labels 	= yearsP_KUBE[x]
			Labels.Data 	= dataP

			ResultP = append(ResultP, Labels)
			Persebaran.Kube = ResultP

			// count jenis_usaha KUBE by id_jenis_usaha
			for i,_ := range LabelsJ_Id {
				var count int
				if err := con.Table("tbl_kube").Where("id_jenis_usaha = ?", LabelsJ_Id[i]).Count(&count).Error; err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
				dataJ = append(dataJ, count)
			}		

			Labels.Labels 	= yearsP_UEP[x]
			Labels.Data 	= dataJ

			ResultJ = append(ResultJ, Labels)
			JenisUsaha.Kube = ResultJ
		}

	}

	ChartDashBoard.HasilMonev = HasilMonev
	ChartDashBoard.Persebaran = Persebaran
	ChartDashBoard.JenisUsaha = JenisUsaha

	defer con.Close()
	return c.JSON(http.StatusOK, ChartDashBoard)
}

// @Summary GeAllMemberPelatihan
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Param id_pendamping query int true "int"
// @Param id_pelatihan query int true "int"
// @Param peruntukan query string false "(string) -> uep | kube "
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/member_pelatihan [get]
func GeAllMemberPelatihan(c echo.Context) (err error) {
	id_pendamping 	:= c.QueryParam("id_pendamping")
	id_pelatihan 	:= c.QueryParam("id_pelatihan")
	For				:= c.QueryParam("peruntukan")

	MemberUep 			:= []models.MemberPelatihanUep{}
	MemberKube 			:= []models.MemberPelatihanKube{}
	var Result []interface{}

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if For == "UEP" || For == "uep" {
		// get uep base on id_pendamping
		if err := con.Table("tbl_uep t1").Select("t1.nama_usaha, t2.id_user, t2.nama, t3.region, 'UEP' as flag ").Joins("join tbl_user t2 on t2.id_user = t1.id_uep").Joins("join view_address t3 on t3.id_kelurahan = t2.id_kelurahan").Where("id_pendamping = ?", id_pendamping).Scan(&MemberUep).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

		if len(MemberUep) != 0 {
			for i,_ := range MemberUep {
				// check if anggota UEP already in tbl_kehadiran or Invited
				var count int
				if err := con.Table("tbl_kehadiran").Where("id_pendamping = ?", id_pendamping).Where("id_pelatihan = ?", id_pelatihan).Where("id_user = ?", MemberUep[i].Id_user).Count(&count).Error; err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
				log.Println("count : ", count)
				if count != 0 { MemberUep[i].Invited = true }

				Result = append(Result, MemberUep[i])
			}
		}

	} else if For == "KUBE" || For == "kube" {
		// get member pendamping from uep
		var KubesMember = []string{"ketua", "sekertaris", "bendahara", "anggota1", "anggota2", "anggota3", "anggota4", "anggota5", "anggota6", "anggota7"}

		for i,_ := range KubesMember {
			if err := con.Table("tbl_kube t1").Select("t1.nama_usaha, t3.region, t2.id_user, t2.nama, t1.nama_kube as flag ").Joins("join tbl_user t2 on t2.id_user = t1." + KubesMember[i]).Where("id_pendamping = ?", id_pendamping).Joins("join view_address t3 on t3.id_kelurahan = t2.id_kelurahan").Scan(&MemberKube).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

			if len(MemberKube) != 0 {
				for i,_ := range MemberKube {
					// check if anggota UEP already in tbl_kehadiran or Invited
					var count int
					if err := con.Table("tbl_kehadiran").Where("id_pendamping = ?", id_pendamping).Where("id_pelatihan = ?", id_pelatihan).Where("id_user = ?", MemberKube[i].Id_user).Count(&count).Error; err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, err)
					}
					if count != 0 { MemberKube[i].Invited = true }

					Result = append(Result, MemberKube[i])
				}
			}
		}		
	} else if For == "" {
		// get uep base on id_pendamping
		if err := con.Table("tbl_uep t1").Select("t1.nama_usaha, t2.id_user, t2.nama, t3.region, 'UEP' as flag ").Joins("join tbl_user t2 on t2.id_user = t1.id_uep").Joins("join view_address t3 on t3.id_kelurahan = t2.id_kelurahan").Where("id_pendamping = ?", id_pendamping).Scan(&MemberUep).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

		if len(MemberUep) != 0 {
			for i,_ := range MemberUep {
				// check if anggota UEP already in tbl_kehadiran or Invited
				var count int
				if err := con.Table("tbl_kehadiran").Where("id_pendamping = ?", id_pendamping).Where("id_pelatihan = ?", id_pelatihan).Where("id_user = ?", MemberUep[i].Id_user).Count(&count).Error; err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
				if count != 0 { MemberUep[i].Invited = true }

				Result = append(Result, MemberUep[i])
			}
		}

		// get member pendamping from uep
		var KubesMember = []string{"ketua", "sekertaris", "bendahara", "anggota1", "anggota2", "anggota3", "anggota4", "anggota5", "anggota6", "anggota7"}

		for i,_ := range KubesMember {
			if err := con.Table("tbl_kube t1").Select("t1.nama_usaha, t3.region, t2.id_user, t2.nama, t1.nama_kube as flag ").Joins("join tbl_user t2 on t2.id_user = t1." + KubesMember[i]).Where("id_pendamping = ?", id_pendamping).Joins("join view_address t3 on t3.id_kelurahan = t2.id_kelurahan").Scan(&MemberKube).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

			if len(MemberKube) != 0 {
				for i,_ := range MemberKube {
					// check if anggota UEP already in tbl_kehadiran or Invited
					var count int
					if err := con.Table("tbl_kehadiran").Where("id_pendamping = ?", id_pendamping).Where("id_pelatihan = ?", id_pelatihan).Where("id_user = ?", MemberKube[i].Id_user).Count(&count).Error; err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, err)
					}
					if count != 0 { MemberKube[i].Invited = true }
					
					Result = append(Result, MemberKube[i])
				}
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
	// q = q.Preload("Kelurahan")
	// q = q.Preload("Kecamatan")
	// q = q.Preload("Kabupaten")
	q = q.Preload("Region")
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
			Users[i].Is_eligible = false
			// continue
		}

		// log.Println("Flag : ", Users[i].Flag)
		// check in kube
		var id_kube_members []int
		var KubesMember = []string{"ketua", "sekertaris", "bendahara", "anggota1", "anggota2", "anggota3", "anggota4", "anggota5", "anggota6", "anggota7"}		
		for o, _ := range KubesMember {
			q2 := con
			q2 = q2.Table("tbl_kube")
			q2 = q2.Where(KubesMember[o] + " = ?", Users[i].Id_user)
			q2 = q2.Pluck(KubesMember[o], &id_kube_members)

			if len(id_kube_members) != 0 {
				if len(id_uep) != 0 {
					Users[i].Flag = "UEPKUBE"
					Users[i].Is_eligible = false
					// continue
				} else if len(id_uep) == 0 {
					Users[i].Flag = "KUBE"
					Users[i].Is_eligible = false
				}
			} 
		}

		if Users[i].Flag == "" {
			Users[i].Is_eligible = true
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

// @Summary GetAllAddressDetail
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/address_detail [get]
func GetAllAddressDetail(c echo.Context) (err error) {

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	Address := []models.Tbl_kabupaten{}
	q := con
	q = q.Model(&Address)
	q = q.Preload("Kecamatan")
	q = q.Find(&Address)	
	
	r := &models.Jn{Msg: Address}

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

	log.Println("youre in persebaran ")

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