package helpers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/models"
	"uepkube-api/db"
	"strconv"
	"log"
	"github.com/ulule/paging"
	"math"
)

func PaginateUep(c echo.Context, r *models.ResPagin) (err error) {
	u := &models.PosPagin{}
	num := 1

	if err := c.Bind(u); err != nil {
		return err
	}

	limit := strconv.FormatInt(int64(u.Size), 10)
	var co int = (u.Page - num) * u.Size
	offset := strconv.FormatInt(int64(co), 10)
	url := "http://localhost:9000/api/v1/uep?limit="+limit+"&offset="+offset

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	store, err := paging.NewGORMStore(con, &ueps)
	if err != nil {
	        log.Fatal(err)
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
	        log.Fatal(errp)
	}
	l := paginator.Limit
	o := paginator.Offset
	t := paginator.Count
	f := false
	la := false
	tp := float64(t)/float64(l)
	rtp := math.Round(tp)

	if u.Page == 1 {f = true}
	if u.Page == int(rtp) {la = true}

	uu := make([]*models.U, len(ueps))
	JoinUep(uu)
	
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

func JoinUep(ur []*models.U) (err error){
	/*prepare DB*/
	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	User := make([]models.Tbl_user, len(ueps))
	R := make([]models.CustU, len(ueps))

	for i := range ueps{
		id := ueps[i].Id_uep
		if err := con.Where(&models.Tbl_user{Id_user:id}).First(&User).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
		// find uep by id + join pendaming-user
		if err := con.Table("tbl_uep").Select("tbl_uep.bantuan_modal, tbl_uep.status,tbl_uep.id_pendamping, tbl_user.nama").Joins("join tbl_user on tbl_user.id_user = tbl_uep.id_pendamping").Where(&models.Tbl_uep{Id_uep:id}).Scan(&R).Error; gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "Uep Not Found")
		}
		ur[i] = &models.U{User[0], R[0]}
	 }
	defer con.Close()
	return err
}