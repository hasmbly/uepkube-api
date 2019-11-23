package helpers

import (
	"uepkube-api/models"
	"github.com/astaxie/beego/utils/pagination"
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

	// paginate produks Uep
	PaginProductsUep = []*models.PaginateProduksUep{}

	// paginate produks Kube
	PaginProductsKube = []*models.PaginateProduksKube{}		

	// usaha-produk
	UsahaProducts = []*models.Tbl_usaha_produk{}

	// uep-produk
	UepProduks = []*models.Tbl_usaha_produk_uep{}

	// kube-produk
	KubeProduks = []*models.Tbl_usaha_produk_kube{}

	// kubes
	kubes = []*models.Tbl_kube{}

	// pelatihan
	Pelatihan = []*models.Tbl_pelatihan{}

	// response paginate
	r models.ResPagin = models.ResPagin{}

	// count
	CountRows int64
)