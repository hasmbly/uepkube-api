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
	 _"log"
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
		q1 = q1.Select("tbl_user.*, tbl_uep.*")

		if err := q1.Scan(&Uep).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

		/*find kube by Ketua:id*/
		// if err := con.Where(&models.Tbl_kube{Ketua:id}).First(&Kube).Error; gorm.IsRecordNotFoundError(err) {
		// 	log.Println("is kube exist : ", Kube)
		// 	r := &models.Jn{Msg: models.U{User, R}}
		// 	defer con.Close()
		// 	return c.JSON(http.StatusOK, r)
		// }		

	} else if errr == nil {
		return GetKube(c)
	}

	if len(Uep) != 0 {
		var tempo []string
		for i,_ := range Uep {
			q2 := con
			q2 = q2.Table("tbl_uep")
			q2 = q2.Joins("join tbl_usaha_produk on tbl_usaha_produk.id_uep = tbl_uep.id_uep")
			q2 = q2.Joins("join tbl_jenis_usaha on tbl_jenis_usaha.id_usaha = tbl_usaha_produk.id_usaha")
			q2 = q2.Where("tbl_uep.id_uep = ?", Uep[i].Id_user)

			q2 = q2.Pluck("tbl_jenis_usaha.jenis_usaha", &tempo)
		
		Uep[i].Jenis_usaha = tempo[0]
		// log.Println("n_pendamping : ", tempo)

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
	// return nil
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