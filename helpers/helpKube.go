package helpers

import (
	// "net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/models"
	"uepkube-api/db"
	"github.com/fatih/structs"
	"fmt"
	// "log"
	"strconv"
	// "github.com/ulule/paging"
	"math"
)

func SetMemberNameKube(s *models.ShowKube, Kube models.Tbl_kube) error {
	/*prepare DB*/
	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)		
	/*
	begin:find member name of Kube
	 */
	kv := structs.Values(Kube)
	kr := kv[4:15]
	
	var tmp []string
	res := make([]string, len(kr))
	ints := make([]int, len(kr))
	
	for i,d := range kr {
		nf := fmt.Sprintf("%+v",d)
		ints[i],_ = strconv.Atoi(nf)
	}

    for i,d := range ints {
		if err := con.Table("tbl_user").Where(&models.Tbl_user{Id_user:d}).Pluck("nama", &tmp).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
    	if d == 0 { tmp[0] = "" }
		res[i] = tmp[0]
    }

    // get alamat from ketua kube (first man)
    var alamat []*string
	if err := con.Table("tbl_user").Where(&models.Tbl_user{Id_user: ints[0]}).Pluck("alamat", &alamat).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

    // get lat from ketua kube (first man)
    var lat []*string
	if err := con.Table("tbl_user").Where(&models.Tbl_user{Id_user: ints[0]}).Pluck("lat", &lat).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

    // get long from ketua kube (first man)
    var lng []*string
	if err := con.Table("tbl_user").Where(&models.Tbl_user{Id_user: ints[0]}).Pluck("lng", &lng).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

    // get photo kube
    var photo []models.Tbl_kube_photo
	if err := con.Table("tbl_kube_photo").Where(&models.Tbl_kube_photo{Id_kube: Kube.Id_kube}).Find(&photo).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	for i,_ := range photo {

			if photo[i].Photo != "" {
				ImageBlob := photo[i].Photo
				photo[i].Photo = "data:image/png;base64," + ImageBlob	
			}

		}	

    /*
	end:find member name of Kube
	 */	
	*s = models.ShowKube{
		Id_kube: 		Kube.Id_kube,
		Nama_kube: 		Kube.Nama_kube,
		Nama_usaha: 	Kube.Nama_usaha,
		Alamat: 		alamat[0],
		Lat: 			lat[0],
		Lng: 			lng[0],
		Items: models.Items{
			Ketua:			res[0],
			Sekertaris:		res[1],
			Bendahara:  	res[2],
			Anggota1:		res[3],
			Anggota2:		res[4],
			Anggota3:		res[5],
			Anggota4:		res[6],
			Anggota5:		res[7],
			Anggota6: 		res[8],
			Anggota7:		res[9],
			Pendamping:		res[10],
		},
		Photo:			photo,
		Flag:			"KUBE",
	}
	defer con.Close()
	return err
}

func PaginateKube(c echo.Context, r *models.ResPagin) (err error) {
	u := &models.PosPagin{}
	num := 1

	// GetLoggedUser(c,"roles")

	if err := c.Bind(u); err != nil {
		return err
	}

	var co int = (u.Page - num) * u.Size
	
	PaginateResult, _ := ExecPaginateKube(u,co,&CountRows)

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

func ExecPaginateKube(f *models.PosPagin, offset int, count *int64) (ur []models.PaginateKubes, err error) {

	// var Pelatihans []models.Tbl_pendamping
	kubes := []models.PaginateKubes{}

	con, err := db.CreateCon()
	if err != nil { return ur, echo.ErrInternalServerError }
	con.SingularTable(true)	

	q := con
	q = q.Table("tbl_kube t1")
	q = q.Limit(int(f.Size))
	q = q.Offset(int(offset))
	q = q.Select("t1.id_kube, t1.nama_kube, t1.bantuan_modal, t1.status, t1.created_at")
	// q = q.Joins("join tbl_user t2 on t2.id_user = t1.id_uep")
	q = q.Joins("join tbl_jenis_usaha t3 on t3.id_usaha = t1.id_jenis_usaha")

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
	q = q.Order(fmt.Sprintf("t1.%s %s",f.SortField,f.SortOrder))	
	
	q = q.Scan(&kubes)
	q = q.Limit(-1)
	q = q.Offset(-1)

	// get Pendampings
	if len(kubes) != 0 {
		for i,_ := range kubes {
			var id_pendamping []int
			var pendamping models.CustomPendamping

			con.Table("tbl_kube").Where("id_kube = ?", kubes[i].Id_kube).Pluck("id_pendamping", &id_pendamping)

			if len(id_pendamping) != 0 {

				for i,_ := range id_pendamping {
					con.Table("tbl_pendamping").Select("tbl_pendamping.*, tbl_user.nama as nama_pendamping").Joins("join tbl_user on tbl_user.id_user = tbl_pendamping.id_pendamping").Where("id_pendamping = ?", id_pendamping[i]).Find(&pendamping)
				}
					kubes[i].Pendamping = pendamping
			}
		}
	}

	// get Usaha
	if len(kubes) != 0 {
		for i,_ := range kubes {
			var kube_usaha models.UsahaKube
 			// var id_produk []int
			var photos []models.Tbl_uepkube_photo

			// log.Println("id_kube : ", kubes[i].Id_kube)
			
			q := con.Table("tbl_kube t1")
			q = q.Select("t1.id_kube, t1.nama_usaha, t2.id_usaha, t2.jenis_usaha")
			q = q.Joins("join tbl_jenis_usaha t2 on t2.id_usaha = t1.id_jenis_usaha")
			q = q.Where("t1.id_kube = ?", kubes[i].Id_kube)
			q = q.Scan(&kube_usaha)

			if kube_usaha.Id_usaha != 0 { kubes[i].Usaha = kube_usaha }

			con.Table("tbl_usaha_kube_photo").Where("id_kube = ?", kubes[i].Id_kube).Find(&photos)

			for index,_ := range photos {
				ImageBlob := photos[index].Photo
				photos[index].Photo = "data:image/png;base64," + ImageBlob			
				kubes[i].Usaha.Photo = photos
			}

			// log.Println("photos : ", photos)
			// log.Println("usaha : ", kubes[i].Usaha)
			
		}
	}

	if err := q.Count(count).Error; err != nil {
		return ur, err
	}

	// log.Println("result : ", Pelatihans)

	defer con.Close()
	return kubes, nil
}