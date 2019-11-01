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

	// kubes
	kubes = []*models.Tbl_kube{}

	r models.ResPagin = models.ResPagin{}
)