package helpers

import (
	// "net/http"
	"github.com/labstack/echo"
	// "github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/models"
	"uepkube-api/db"
	//"strconv"
	"log"
	// "fmt"
	// "github.com/ulule/paging"
	"math"
)

func PaginateProduk(c echo.Context, r *models.ResPagin) (err error) {
	u := &models.PosPagin{}
	num := 1

	// GetLoggedUser(c,"roles")

	if err := c.Bind(u); err != nil {
		return err
	}

	var co int = (u.Page - num) * u.Size
	// offset := strconv.FormatInt(int64(co), 10)
	// log.Println("Ori Offset is : ", offset)

	// limit := strconv.FormatInt(int64(u.Size), 10)
	// log.Println("Ori size is : ", limit)



	// url := "http://localhost:9000/api/v1/produk?limit="+limit+"&offset="+offset

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// store, err := paging.NewGORMStore(con, "tbl_usaha_produk",&PaginProducts)
	// if err != nil {
	//         e := fmt.Sprintf("%v",err)
	//         log.Println(e)
	//         return echo.NewHTTPError(http.StatusInternalServerError, e)	        
	// }
	// options := paging.NewOptions()
	// request, _ := http.NewRequest("GET", url, nil)

	// log.Println("Request Filter is :", u.Filters)
	
	// Filters := make([]map[string]string, len(u.Filters))
	// for i,_ := range u.Filters {
	// 	fields := map[string]string{
	// 		"key": u.Filters[i].Key,
	// 		"operation": u.Filters[i].Operation,
	// 		"value": u.Filters[i].Value,
	// 	}
	//     Filters[i] = fields
	// }

	// var fakeFilters []map[string]string

	// paginator,_ := paging.NewOffsetPaginator(store, request, options, u.SortField, u.SortOrder, fakeFilters)
	// errp := paginator.PageProducts()
	// if errp != nil {
	//         e := fmt.Sprintf("%v",errp)
	//         log.Println(e)
	//         return echo.NewHTTPError(http.StatusInternalServerError, e)
	// }

	// log.Println("Products : ", PaginProducts)
	// 
	

	// PaginProducts = make([]*models.PaginateProduks, GetCount)
	
	ok, _ := PaginateProd(u,co,&CountRows)

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

	log.Println("Result is : ", ok)

	*r = models.ResPagin{
		Content:ok,
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
	// defer con.Close()
	return err
}

func PaginateProd(u *models.PosPagin, offset int, count *int64) (ur []*models.PaginateProduks, err error) {

	// var UsahaProductsTemp []models.Tbl_usaha_produk
	var UsahaProducts []models.Tbl_usaha_produk
	var Products []*models.PaginateProduks

	con, err := db.CreateCon()
	if err != nil { return ur, echo.ErrInternalServerError }
	con.SingularTable(true)	

	up := con
	up = up.Table("tbl_usaha_produk")
	up = up.Limit(int(u.Size))
	up = up.Offset(int(offset))
	if len(u.Filters) != 0 {
		var id_produk []int

		x1 := con
		x1 = x1.Table("tbl_produk")
		x1 = x1.Where("nama_produk like ?", "%"+u.Filters[0].Value+"%")
		x1 = x1.Pluck("id_produk", &id_produk)

		if len(id_produk) != 0 {

				up = up.Where("id_produk = ?", id_produk[0])				
				for i,_ := range id_produk {				
					up = up.Or("id_produk = ?", id_produk[i])
				}
		}		
	}	
	up = up.Scan(&UsahaProducts)
	up = up.Limit(-1)
	up = up.Offset(-1)

	if err := up.Count(count).Error; err != nil {
		return ur, err
	}


	log.Println("final_usaha_produk : ", UsahaProducts)

	tmp := make([]*models.PaginateProduks, len(UsahaProducts) )

		for i := range UsahaProducts{
			var id_uep = UsahaProducts[i].Id_uep
			var id_kube = UsahaProducts[i].Id_kube
			
				q := con
				q = q.Table("tbl_usaha_produk t1")
			if id_uep != 0 {
				q = q.Select(
					"t1.id,t2.nama,t2.alamat,t2.no_hp,t3.nama_produk,t3.deskripsi,t4.jenis_usaha")
				q = q.Joins("join tbl_user t2 on t2.id_user = t1.id_uep")
				q = q.Joins("join tbl_jenis_usaha t4 on t4.id_usaha = t1.id_usaha")
				q = q.Joins("join tbl_produk t3 on t3.id_produk = t1.id_produk")
				q = q.Where("t1.id_uep = ?", id_uep)
			}else {
				q = q.Select(
					"t1.id,tbl_kube.nama_kube as nama,tbl_user.alamat,tbl_user.no_hp,tbl_produk.nama_produk,tbl_produk.deskripsi,tbl_jenis_usaha.jenis_usaha")
				q = q.Joins("join tbl_kube on tbl_kube.id_kube = t1.id_kube")
				q = q.Joins("join tbl_user on tbl_user.id_user = tbl_kube.ketua")
				q = q.Joins("join tbl_jenis_usaha on tbl_jenis_usaha.id_usaha = t1.id_usaha")
				q = q.Joins("join tbl_produk on tbl_produk.id_produk = t1.id_produk")
				q = q.Where("t1.id_kube = ?", id_kube)
			}
			q = q.Scan(&Products)
			tmp[i] = Products[0]
	 	}

	log.Println("Final Result is : ", tmp)

	defer con.Close()
	return tmp, nil
}