package helpers

import (
	// "net/http"
	"github.com/labstack/echo"
	// "github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/db"
	"uepkube-api/models"
	"log"
	"math"	
	"fmt"
)

func PaginateAktivitas(c echo.Context, r *models.ResPagin) (err error) {
	flag = "AKTIVITAS"
	host = c.Request().Host		

	u := &models.PosPagin{}
	num := 1

	// GetLoggedUser(c,"roles")

	if err := c.Bind(u); err != nil {
		return err
	}

	var co int = (u.Page - num) * u.Size
	
	PaginateResult, _ := ExecPaginateAktivitas(u,co,&CountRows)

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

func ExecPaginateAktivitas(f *models.PosPagin, offset int, count *int64) (ur []models.Tbl_activity, err error) {

	Aktivitas := []models.Tbl_activity{}

	con, err := db.CreateCon()
	if err != nil { return ur, echo.ErrInternalServerError }
	con.SingularTable(true)	

	q := con
	q = q.Model(&Aktivitas)
	q = q.Limit(int(f.Size))
	q = q.Offset(int(offset))
	q = q.Preload("Photo")

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
	q = q.Order(fmt.Sprintf("%s %s",f.SortField,f.SortOrder))
	
	q = q.Find(&Aktivitas)
	q = q.Limit(-1)
	q = q.Offset(-1)
	
	// photo
	if len(Aktivitas) != 0 {
		for i, _ := range Aktivitas {
			id := Aktivitas[i].Id
			if len(Aktivitas[i].Photo) != 0 {
				for x, _ := range Aktivitas[i].Photo {

					if Aktivitas[i].Photo[x].Type == "IMAGE" {

						id_photo := Aktivitas[i].Photo[x].Id
						
						tmpPath	= fmt.Sprintf(GoPath + "/src/uepkube-api/static/assets/images/%s_id_%d_photo_id_%d.png", flag,id,id_photo)
						urlPath	= fmt.Sprintf("http://%s/images/%s_id_%d_photo_id_%d.png", host,flag,id,id_photo)
						blobFile = Aktivitas[i].Photo[x].Files

						if check := CreateFile(tmpPath, blobFile); check == false {
							log.Println("blob is empty : ", check)
						}
						Aktivitas[i].Photo[x].Files = urlPath
					}
				}
			}
		}
	}

	if err := q.Count(count).Error; err != nil {
		return ur, err
	}

	// log.Println("result : ", Aktivitas)

	defer con.Close()
	return Aktivitas, nil
}