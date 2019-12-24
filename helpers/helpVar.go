package helpers

import (
	"uepkube-api/models"
	"github.com/astaxie/beego/utils/pagination"
	 "os"
)

var (
	// paginator
	paginator = &pagination.Paginator{}

	// ueps
	ueps = []*models.Tbl_uep{}
	U = []*models.U{}

	// produk
	produks = []*models.Tbl_produk{}

	// paginate produks
	PaginProducts = []models.PaginateProduks{}

	// usaha-produk
	UsahaProducts = []*models.Tbl_usaha_produk{}

	// kubes
	kubes = []*models.Tbl_kube{}

	// pelatihan
	Pelatihan = []*models.PaginatePelatihan{}

	// response paginate
	r models.ResPagin = models.ResPagin{}

	// count
	CountRows int64

	GoPath string = os.Getenv("GOPATH")	
)