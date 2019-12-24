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
	// "os"

	// "bufio"
	// "encoding/base64"	
	// "io/ioutil"	
)

var tmpPath, urlPath, blobFile,flag,host string

func PaginatePelatihan(c echo.Context, r *models.ResPagin) (err error) {
	flag = "PELATIHAN"
	host = c.Request().Host

	u := &models.PosPagin{}
	num := 1

	// GetLoggedUser(c,"roles")

	if err := c.Bind(u); err != nil {
		return err
	}

	var co int = (u.Page - num) * u.Size
	
	PaginateResult, _ := ExecPaginatePelatihan(u,co,&CountRows)

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

	// log.Println("Result is : ", PaginateResult)

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

func ExecPaginatePelatihan(f *models.PosPagin, offset int, count *int64) (ur []models.Tbl_pelatihan, err error) {

	// var Pelatihans []models.Tbl_pendamping
	Pelatihans := []models.Tbl_pelatihan{}

	con, err := db.CreateCon()
	if err != nil { return ur, echo.ErrInternalServerError }
	con.SingularTable(true)	

	q := con
	q = q.Model(&Pelatihans)
	q = q.Limit(int(f.Size))
	q = q.Offset(int(offset))
	q = q.Preload("Files", "type = ?", "PDF")
	q = q.Preload("Photo", "type = ?", "IMAGE")

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
	
	q = q.Find(&Pelatihans)
	q = q.Limit(-1)
	q = q.Offset(-1)

	// get photos
	if len(Pelatihans) != 0 {
		for i,_ := range Pelatihans {
			id := Pelatihans[i].Id_pelatihan

			// for Photo
			for x,_ := range Pelatihans[i].Photo {
				id_photo := Pelatihans[i].Photo[x].Id

				tmpPath	= fmt.Sprintf(GoPath + "/src/uepkube-api/static/assets/images/%s_id_%d_photo_id_%d.png", flag,id,id_photo)
				urlPath	= fmt.Sprintf("http://%s/images/%s_id_%d_photo_id_%d.png", host,flag,id,id_photo)
				blobFile = Pelatihans[i].Photo[x].Files

				if check := CreateFile(tmpPath, blobFile); check == false {
					log.Println("blob is empty : ", check)
				}

				Pelatihans[i].Photo[x].Files = urlPath
			}

			// for files pdf
			for y,_ := range Pelatihans[i].Files {
				id_pdf := Pelatihans[i].Files[y].Id

				tmpPath	= fmt.Sprintf(GoPath + "/src/uepkube-api/static/assets/pdf/%s_id_%d_pdf_id_%d.pdf", flag,id,id_pdf)
				urlPath	= fmt.Sprintf("http://%s/pdf/%s_id_%d_pdf_id_%d.pdf", host,flag,id,id_pdf)
				blobFile = Pelatihans[i].Files[y].Files

				if check := CreateFile(tmpPath, blobFile); check == false {
					log.Println("blob is empty : ", check)
				}

				Pelatihans[i].Files[y].Files = urlPath
			}			
		}
	}

	if err := q.Count(count).Error; err != nil {
		return ur, err
	}

	// log.Println("result : ", Pelatihans)

	defer con.Close()
	return Pelatihans, nil
}