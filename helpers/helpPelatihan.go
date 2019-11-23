package helpers

import (
	"net/http"
	"github.com/labstack/echo"
	// "github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/models"
	"uepkube-api/db"
	"strconv"
	"log"
	"fmt"
	"github.com/ulule/paging"
	"math"
)

func PaginatePelatihan(c echo.Context, r *models.ResPagin) (err error) {
	u := &models.PosPagin{}
	num := 1

	if err := c.Bind(u); err != nil {
		return err
	}

	limit := strconv.FormatInt(int64(u.Size), 10)
	var co int = (u.Page - num) * u.Size
	offset := strconv.FormatInt(int64(co), 10)
	url := "http://localhost:9000/api/v1/pelatihan1?limit="+limit+"&offset="+offset

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	store, err := paging.NewGORMStore(con, "tbl_pelatihan", &Pelatihan)
	if err != nil {
	        e := fmt.Sprintf("%v",err)
	        log.Println(e)
	        return echo.NewHTTPError(http.StatusInternalServerError, e)
	}
	options := paging.NewOptions()
	request, _ := http.NewRequest("GET", url, nil)

	filters := make([]map[string]string, len(u.Filters))
	for i,_ := range u.Filters {
		fields := map[string]string{
			"key": u.Filters[i].Key,
			"operation": u.Filters[i].Operation,
			"value": u.Filters[i].Value,
		}
	    filters[i] = fields
	}
	
	paginator,_ := paging.NewOffsetPaginator(store, request, options, u.SortField, u.SortOrder, filters)

	errp := paginator.Page()
	if errp != nil {
	        e := fmt.Sprintf("%v",errp)
	        log.Println(e)
	        return echo.NewHTTPError(http.StatusInternalServerError, e)
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

	*r = models.ResPagin{
		Content:Pelatihan,
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