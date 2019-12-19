package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	 _"github.com/jinzhu/gorm/dialects/mysql"
	 "uepkube-api/models"
	 "uepkube-api/db"
	 "strconv"
	 "uepkube-api/helpers"
	 "log"
	 "fmt"

	"bufio"
	"encoding/base64"	
	"io/ioutil"			 
)
/*@Summary GetKubeById
@Tags Kube-Controller
@Accept  json
@Produce  json
@Param id query int true "int"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /kube [get]*/
func GetKube(c echo.Context) error {
	/*prepare DB*/
	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)	

	var val string
	ShowKubes := models.ShowKube{}
	var tempo []interface{}
	
	r := &models.Jn{}

	/*check if query key -> "val"*/
	qk := c.QueryParams()
	for k,v := range qk {
		if k == "val" {

			/*find kube by Nama_kube:*/
			val = v[0]
			Kubes 	:= []models.Tbl_kube{}
			if err := con.Where("nama_kube LIKE ?", "%" + val + "%").Find(&Kubes).Error; gorm.IsRecordNotFoundError(err)  {
				return echo.NewHTTPError(http.StatusNotFound, "Kube Not Found")
			}
			for i,_ := range Kubes {
				helpers.SetMemberNameKube(&ShowKubes, Kubes[i])
				tempo = append(tempo, ShowKubes)
			}			
			r.Msg = tempo

		} else if k == "id" {

			/*find kube by Id Kube :*/
			val = v[0]
			id,_ := strconv.Atoi(val)
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

			// get kubes_member
			var KubesMember = []string{"ketua", "sekertaris", "bendahara", "anggota1", "anggota2", "anggota3", "anggota4", "anggota5", "anggota6", "anggota7"}
			tmp := []models.Kubes_items{}

			for i, _ := range KubesMember {
			
				if err := con.Table("tbl_kube t1").Select("t2.id_user, t2.nama, t2.nik, '" +KubesMember[i] + "' as posisi").Joins("join tbl_user t2 on t2.id_user = t1." + KubesMember[i]).Where("id_kube = ?", id).Scan(&tmp).Error; err != nil { return echo.ErrInternalServerError }

				if len(tmp) != 0 {
					Kube.Items = append(Kube.Items, tmp[0])
				}

			}

			r.Msg = Kube
		}
	}




	defer con.Close()
	return c.JSON(http.StatusOK, r)
}

/*@Summary GetPaginateKube
@Tags Kube-Controller
@Accept  json
@Produce  json
@Param kube body models.PosPagin true "Show Kube Paginate"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /kube [post]*/
func GetPaginateKube(c echo.Context) (err error) {	
	if err := helpers.PaginateKube(c, &r); err != nil {
		return err
	}	
	return c.JSON(http.StatusOK, r)
}

/*@Summary AddKube
@Tags Kube-Controller
@Accept  json
@Produce  json
@Param kube body models.Kube true "Add Kube"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /kube/add [post]*/
func AddKube(c echo.Context) (err error) {
	Kube := &models.Kube{}

	if err := c.Bind(Kube); err != nil {
		return err
	}

	// validation
	if Kube.Id_pendamping == 0 { 
		return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Id Pendamping") 
	}
	if Kube.Id_jenis_usaha == 0 { 
		return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Id jenis usaha") 
	}
	if Kube.Id_periods == 0 { 
		return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Bantuan Modal") 
	}
	if Kube.Nama_usaha == "" { 
		return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Nama Usaha Modal") 
	}	

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// create kube
	kube := &models.Tbl_kube{}
	kube = Kube.Tbl_kube
	if err := con.Create(&kube).Error; err != nil {return echo.ErrInternalServerError}	

	// store bantuan_periods history
	kubePeriods := &models.Tbl_periods_uepkube{}
	kubePeriods.Id_kube = Kube.Id_kube
	kubePeriods.Id_periods = Kube.Id_periods
	if err := con.Create(&kubePeriods).Error; err != nil {return echo.ErrInternalServerError}	

	// store credit_debit
	creditDebit := &models.Tbl_credit_debit{}
	creditDebit.Id_kube = Kube.Id_kube
	creditDebit.Debit = 1
	creditDebit.Id_periods = Kube.Id_periods
	var nilai []float32
	con.Table("tbl_bantuan_periods").Where("id = ?", creditDebit.Id_periods).Pluck("bantuan_modal", &nilai)
	creditDebit.Nilai = nilai[0]
	creditDebit.Deskripsi = fmt.Sprintf("Credit dengan nilai : Rp. %.2f,-", nilai[0])
	if err := con.Create(&creditDebit).Error; err != nil {return echo.ErrInternalServerError}

	// add queue monev_uepkube
	monev := &models.Tbl_monev_uepkube{}
	monev.Id_kube = Kube.Id_kube
	monev.Id_pendamping = Kube.Id_pendamping
	monev.Is_monev = "BELUM"
	monev.Id_periods = Kube.Id_periods
	if err := con.Create(&monev).Error; err != nil {return echo.ErrInternalServerError}

	defer con.Close()

	r := &models.Jn{Msg: "Success Store Data"}
	return c.JSON(http.StatusOK, r)
}

/*@Summary UpdateKube
@Tags Kube-Controller
@Accept  json
@Produce  json
@Param kube body models.Tbl_kube true "Update Kube"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /kube [put]*/
func UpdateKube(c echo.Context) (err error) {
	kube := &models.Tbl_kube{}

	if err := c.Bind(kube); err != nil {
		return err
	}

	if kube.Id_kube == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id")
	}

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Model(&models.Tbl_kube{}).UpdateColumns(&kube).Error; err != nil {
		return echo.ErrInternalServerError
	}

	if err := con.Table("tbl_kube").Where("id_kube = ?",kube.Id_kube).UpdateColumn("status", kube.Status).Error; err != nil {return echo.ErrInternalServerError}

	defer con.Close()

	r := &models.Jn{Msg: "Success Update Data"}
	return c.JSON(http.StatusOK, r)
}

/*@Summary DeleteKube
@Tags Kube-Controller
@Accept  json
@Produce  json
@Param id path int true "Delete Kube by id"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /kube/{id} [post]*/
func DeleteKube(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("id"))

	if id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	}

	kube := &models.Tbl_kube{}
	kube.Id_kube = id

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// delete kube
	if err := con.Delete(&kube).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// delete kube_photo
	if err := con.Where("id_kube = ?", kube.Id_kube).Delete(models.Tbl_kube_photo{}).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	defer con.Close()

	r := &models.Jn{Msg: "Success Delete Data"	}
	return c.JSON(http.StatusOK, r)	
}


func UploadKubeFiles(c echo.Context) (err error) {
	id, _ 		:= strconv.Atoi(c.QueryParam("id"))
	is_display, _ 	:= strconv.Atoi(c.QueryParam("is_display"))

	if id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	}

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["photo"]

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
		
		// init photo model
		photo := &models.Tbl_kube_photo{}

		photo.Id_kube = id
		photo.Is_display = is_display
		photo.Photo   = encoded

		if err := con.Create(&photo).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
		
	}
	
	defer con.Close()

	log.Println("Uploads Kube's file to id : ", id)
	r := &models.Jn{Msg: "Success Upload files"}
	return c.JSON(http.StatusOK, r)
	
}