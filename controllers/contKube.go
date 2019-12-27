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
	
	var tmpPath, urlPath, blobFile, flag, host string
	flag = "KUBE"
	host = c.Request().Host

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
			// q = q.Preload("PeriodsHistory.BantuanPeriods.Usaha", func(q *gorm.DB) *gorm.DB {
			// 	return q.Where("id_kube = ?", id).Preload("JenisUsaha")
			// })
			// q = q.Preload("PeriodsHistory.BantuanPeriods.Usaha.AllProduk.DetailProduk.JenisProduk")
			// q = q.Preload("PeriodsHistory.BantuanPeriods.MonevHistory", func(q *gorm.DB) *gorm.DB {
			// 	return q.Where("id_kube = ?", id)
			// })				
			// q = q.Preload("PeriodsHistory.BantuanPeriods.CreditDebit", func(q *gorm.DB) *gorm.DB {
			// 	return q.Where("id_kube = ?", id)
			// })
			q = q.Preload("Pendamping", func(q *gorm.DB) *gorm.DB {
				return q.Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Select("tbl_pendamping.*,tbl_user.nama")
			})
			q = q.Preload("Photo", func(q *gorm.DB) *gorm.DB {
				return q.Where("id_kube = ?", id)	
			})
			q = q.First(&Kube, id)

			// get kubes_member
			var KubesMember = []string{"ketua", "sekertaris", "bendahara", "anggota1", "anggota2", "anggota3", "anggota4", "anggota5", "anggota6", "anggota7"}
			tmp := []models.Kubes_items{}

			for i, _ := range KubesMember {
				if err := con.Table("tbl_kube t1").Select("t2.*, '" +KubesMember[i] + "' as posisi").Joins("join tbl_user t2 on t2.id_user = t1." + KubesMember[i]).Where("id_kube = ?", id).Scan(&tmp).Error; err != nil { return echo.ErrInternalServerError }
				if len(tmp) != 0 {
					Regions := models.View_address{}
					if err := con.Table("view_address").Select("view_address.*").Where("id_kelurahan = ?", tmp[0].Id_kelurahan).Scan(&Regions).Error; err != nil { return echo.ErrInternalServerError }
					Kube.Items = append(Kube.Items, tmp[0])
					Kube.Region = Regions
					// log.Println("id_kelurahan kube : ", Kube.Region)
				}
			}

			// rename photo
			for i, _ := range Kube.Photo {
					id_photo := Kube.Photo[i].Id

					tmpPath	= fmt.Sprintf(helpers.GoPath + "/src/uepkube-api/static/assets/images/%s_id_%d_photo_id_%d.png", flag,id,id_photo)
					urlPath	= fmt.Sprintf("http://%s/images/%s_id_%d_photo_id_%d.png", host,flag,id,id_photo)
					blobFile = Kube.Photo[i].Files

					if check := CreateFile(tmpPath, blobFile); check == false {
						log.Println("blob is empty : ", check)
					}
				
					Kube.Photo[i].Files = urlPath
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

	// get log post
	helpers.FetchPost(Kube)	

	// validation
	if Kube.Id_pendamping == 0 { 
		return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Id Pendamping") 
	}
	if Kube.Id_jenis_usaha == 0 { 
		return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Id jenis usaha") 
	}
	// if Kube.Id_periods == 0 { 
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Please Fill id_periods ( Bantuan Modal )") 
	// }
	if Kube.Nama_usaha == "" { 
		return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Nama Usaha Modal") 
	}	
	if len(Kube.Items) == 0 { 
		return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Items as Member for Kube") 
	}		

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	kube := &models.Tbl_kube{}
	kube = Kube.Tbl_kube	

	// validation for items as memberKube
	if len(kube.Items) != 0 {
		for i, _ := range kube.Items {
			if kube.Items[i].Nik == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Please Fill NIK") 
			} else {
				// validation NIK
				if len(kube.Items[i].Nik) > 16 || len(kube.Items[i].Nik) < 16 {
					return echo.NewHTTPError(http.StatusBadRequest, "Please fill NIK with 16 Digits")
				}			
			}
			if kube.Items[i].Posisi == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Posisi") 
			}			
			if kube.Items[i].Nama == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Nama") 
			}
		}

		// execute store memberKube
		for o, _ := range kube.Items {
			
			Memberskube := &models.Tbl_user{}
			Memberskube = kube.Items[o].Tbl_user

			if Memberskube.Id_user == 0 {
				// create
				if err := con.Create(&Memberskube).Error; err != nil {return echo.ErrInternalServerError}
			} else if Memberskube.Id_user != 0 {
				// update
				if err := con.Model(&models.Tbl_user{}).UpdateColumns(&Memberskube).Error; err != nil {
					return echo.ErrInternalServerError
				}				
			}

			if kube.Items[o].Posisi == "ketua" { kube.Ketua = Memberskube.Id_user }
			if kube.Items[o].Posisi == "sekertaris" || kube.Items[o].Posisi == "sekretaris" { kube.Sekertaris = Memberskube.Id_user }
			if kube.Items[o].Posisi == "bendahara" { kube.Bendahara = Memberskube.Id_user }
			if kube.Items[o].Posisi == "anggota1" { kube.Anggota1 = Memberskube.Id_user }
			if kube.Items[o].Posisi == "anggota2" { kube.Anggota2 = Memberskube.Id_user }
			if kube.Items[o].Posisi == "anggota3" { kube.Anggota3 = Memberskube.Id_user }
			if kube.Items[o].Posisi == "anggota4" { kube.Anggota4 = Memberskube.Id_user }
			if kube.Items[o].Posisi == "anggota5" { kube.Anggota5 = Memberskube.Id_user }
			if kube.Items[o].Posisi == "anggota6" { kube.Anggota6 = Memberskube.Id_user }
			if kube.Items[o].Posisi == "anggota7" { kube.Anggota7 = Memberskube.Id_user }

		}
	}

	// check kube is exist
    var n_kube []string
	con.Table("tbl_kube").Where("nama_kube = ?", kube.Nama_kube).Pluck("nama_kube", &n_kube)
	if len(n_kube) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Maaf Nama Kube sudah digunakan")
	}	

	// create kube
	// log.Println("kube : ", kube)
	if err := con.Create(&kube).Error; err != nil {return echo.ErrInternalServerError}	

	// store usaha_kube
	// kubeUsaha := &models.Tbl_usaha_uepkube{}
	// kubeUsaha.Id_kube = Kube.Id_kube
	// kubeUsaha.Nama_usaha = Kube.Nama_usaha
	// kubeUsaha.Id_jenis_usaha = Kube.Id_jenis_usaha
	// kubeUsaha.Id_periods = Kube.Id_periods
	// if err := con.Create(&kubeUsaha).Error; err != nil {return echo.ErrInternalServerError}	

	// store bantuan_periods history
	// kubePeriods := &models.Tbl_periods_uepkube{}
	// kubePeriods.Id_kube = Kube.Id_kube
	// kubePeriods.Id_periods = Kube.Id_periods
	// if err := con.Create(&kubePeriods).Error; err != nil {return echo.ErrInternalServerError}	

	// store credit_debit
	creditDebit := &models.Tbl_inventory{}
	creditDebit.Id_kube = Kube.Id_kube
	creditDebit.Debit = 1
	// var nilai []float32
	// con.Table("tbl_bantuan_periods").Where("id = ?", creditDebit.Id_periods).Pluck("bantuan_modal", &nilai)
	// creditDebit.Nilai = nilai[0]
	// creditDebit.Description = fmt.Sprintf("Credit dengan nilai : Rp. %.2f,-", nilai[0])
	if err := con.Create(&creditDebit).Error; err != nil {return echo.ErrInternalServerError}

	// add queue monev_uepkube
	monev := &models.Tbl_monev_final{}
	monev.Id_kube = Kube.Id_kube
	monev.Id_pendamping = Kube.Id_pendamping
	monev.Is_monev = "BELUM"
	// monev.Id_periods = Kube.Id_periods
	monev.Flag = "KUBE"
	if err := con.Create(&monev).Error; err != nil {return echo.ErrInternalServerError}

	defer con.Close()

	r := &models.Jn1{Msg: "Success Store Data", Id: Kube.Id_kube}
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
	Kube := &models.Kube{}

	if err := c.Bind(Kube); err != nil {
		return err
	}

	if Kube.Id_kube == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id")
	}

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	kube := &models.Tbl_kube{}
	kube = Kube.Tbl_kube	

	// update kube
	if err := con.Model(&models.Tbl_kube{}).UpdateColumns(&kube).Error; err != nil {
		return echo.ErrInternalServerError
	}

	// update kube status
	if err := con.Table("tbl_kube").Where("id_kube = ?",kube.Id_kube).UpdateColumn("status", kube.Status).Error; err != nil {return echo.ErrInternalServerError}

	// store usaha_kube
	// kubeUsaha := &models.Tbl_usaha_uepkube{}
	// kubeUsaha.Id_kube = Kube.Id_kube
	// kubeUsaha.Nama_usaha = Kube.Nama_usaha
	// kubeUsaha.Id_jenis_usaha = Kube.Id_jenis_usaha
	// kubeUsaha.Id_periods = Kube.Id_periods
	// if err := con.Model(&models.Tbl_usaha_uepkube{}).Where("id_kube = ?", kubeUsaha.Id_kube).Where("id_jenis_usaha = ?", kubeUsaha.Id_jenis_usaha).Where("id_periods = ?", kubeUsaha.Id_periods).UpdateColumns(&kubeUsaha).Error; err != nil {
	// 	return echo.ErrInternalServerError
	// }

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
		
		// init UepFiles model
		UepFiles := &models.Tbl_uepkube_files{}

		UepFiles.Id_kube = id
		UepFiles.Files   = encoded
		UepFiles.Description = description
		UepFiles.Type = types
		UepFiles.Is_display = is_display

		if err := con.Create(&UepFiles).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}	
		
		// log.Println("encoded : ", encoded)	

	}

	defer con.Close()

	log.Println("Uploads Kube's file to id : ", id)
	r := &models.Jn{Msg: "Success Upload files"}
	return c.JSON(http.StatusOK, r)
	
}