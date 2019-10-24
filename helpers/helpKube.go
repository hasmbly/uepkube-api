package helpers

import (
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	 _"github.com/jinzhu/gorm/dialects/mysql"
	 "uepkube-api/models"
	 "uepkube-api/db"
	 "fmt"
	 "github.com/fatih/structs"
	 "strconv"
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
	
	var tmp []string
	res := make([]string, len(kr))
	ints := make([]int, len(kr))
	
	for i,d := range kr {
		nf := fmt.Sprintf("%+v",d)
		ints[i],_ = strconv.Atoi(nf)
	}
    for i,d := range ints {
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