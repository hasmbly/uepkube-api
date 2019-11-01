package helpers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/models"
	"uepkube-api/db"
	"github.com/fatih/structs"
	"fmt"
	"strconv"
	"github.com/ulule/paging"
	"math"
)

func SetMemberNameKube(s *models.Ktype, Kube models.Tbl_kube) error {
	/*prepare DB*/
	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)		
	/*
	begin:find member name of Kube
	 */
	kv := structs.Values(Kube)
	kr := kv[4:15]
	// fmt.Println(kr)
	
	var tmp []string
	res := make([]string, len(kr))
	ints := make([]int, len(kr))
	
	for i,d := range kr {
		// fmt.Println(d)
		nf := fmt.Sprintf("%+v",d)
		// fmt.Println(nf)
		ints[i],_ = strconv.Atoi(nf)
	}
	// fmt.Println(ints)
    for i,d := range ints {
		// fmt.Println(d)
		if err := con.Table("tbl_user").Where(&models.Tbl_user{Id_user:d}).Pluck("nama", &tmp).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
		res[i] = tmp[0]
    }
    /*
	end:find member name of Kube
	 */	
	*s = models.Ktype{
		Id_kube: 		Kube.Id_kube,
		Nama_kube: 		Kube.Nama_kube,
		Jenis_usaha: 	Kube.Jenis_usaha,
		Bantuan_modal: 	Kube.Bantuan_modal,
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
		Photo:			Kube.Photo,
		Status:			Kube.Status,
	}
	defer con.Close()
	return err
}

func PaginateKube(c echo.Context, r *models.ResPagin) (err error) {
	u := &models.PosPagin{}
	num := 1

	if err := c.Bind(u); err != nil {
		return err
	}

	limit := strconv.FormatInt(int64(u.Size), 10)
	var co int = (u.Page - num) * u.Size
	offset := strconv.FormatInt(int64(co), 10)
	url := "http://localhost:9000/api/v1/kube?limit="+limit+"&offset="+offset

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	store, err := paging.NewGORMStore(con, &kubes)
	if err != nil {
	        return err
	}
	options := paging.NewOptions()
	request, _ := http.NewRequest("GET", url, nil)

	test := make([]map[string]string, len(u.Filters))
	for i,_ := range u.Filters {
		fields := map[string]string{
			"key": u.Filters[i].Key,
			"operation": u.Filters[i].Operation,
			"value": u.Filters[i].Value,
		}
	    test[i] = fields
	}

	
	paginator,_ := paging.NewOffsetPaginator(store, request, options, u.SortField, u.SortOrder, test)

	errp := paginator.Page()
	if errp != nil {
	        return errp
	}

	l 	:= paginator.Limit
	o 	:= paginator.Offset
	t 	:= paginator.Count
	f 	:= false
	la 	:= false
	tp 	:= float64(t)/float64(l)
	rtp := math.Round(tp)

	if u.Page == 1 {f = true}
	if u.Page == int(rtp) {la = true}

	uu := make([]*models.PaginateKubes, len(kubes))
	JoinKube(uu)

	*r = models.ResPagin{
		Content:uu,
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

func JoinKube(ur []*models.PaginateKubes) (err error){
	for i := range kubes{
		ur[i] = &models.PaginateKubes{
			kubes[i].Id_kube,
			kubes[i].Nama_kube,
			kubes[i].Jenis_usaha,
			kubes[i].Bantuan_modal,
			kubes[i].Photo,
			kubes[i].Status,		
		}
	 }
	return err
}