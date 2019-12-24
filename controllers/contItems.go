package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	 "log"
)

// @Summary Paginate -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev
// @Tags Universal-Controller
// @Accept  json
// @Produce  json
// @Param key path string true "Key -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev"
// @Param sample body models.PosPagin true "Paginate ItemsUepKubePendamping| verifikator"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /{key} [post]
func GetPaginateItems(c echo.Context) (err error) {
	key 	:= c.Param("key")

	log.Println("key : ", key)

	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill key")
	}

	switch key {

	case "uep":
		return GetPaginateUep(c)
	case "kube":
		return GetPaginateKube(c)
	case "pendamping":
		return GetPaginatePendamping(c)
	case "verifikator":
		return GetPaginateVerifikator(c)
	case "produk":
		return GetPaginateProduk(c)
	case "pelatihan":
		return GetPaginatePelatihan(c)			
	case "log_aktivitas":
		return GetPaginateAktivitas(c)			
	case "inventaris":
		return GetPaginateInventaris(c)
	case "lap_keu":
		return GetPaginateLapKeu(c)	
	case "monev":
		return GetPaginateMonev(c)	

	default:
		return echo.NewHTTPError(http.StatusBadRequest, "please, choose the right key")

	}

	return nil
}

// @Summary GetDetails -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev
// @Tags Universal-Controller
// @Accept  json
// @Produce  json
// @Param key path string true "Key -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev"
// @Param id query int true "int"
// @Param for query int false "for (int) -> uep : 0 | kube : 1"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /{key} [get]
func GetItems(c echo.Context) error {
	key 	:= c.Param("key")

	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill key")
	}

	switch key {

	case "uep":
		return GetUep(c)
	case "kube":
		return GetKube(c)
	case "pendamping":
		return GetPendamping(c)
	case "verifikator":
		return GetVerifikator(c)
	case "produk":
		return GetProduk(c)
	case "pelatihan":
		return GetPelatihan(c)
	case "log_aktivitas":
		return GetAktivitas(c)
	case "inventaris":
		return GetInventaris(c)
	case "lap_keu":
		return GetLapKeu(c)	
	case "monev":
		return GetMonev(c)	

	default:
		return echo.NewHTTPError(http.StatusBadRequest, "please, choose the right key")
	}

	return nil
}

// @Summary Add -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev
// @Tags Universal-Controller
// @Accept  json
// @Produce  json
// @Param key path string true "Key -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev "
// @Param sample body models.Dummy true "Add ItemsUepKubePendamping"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /add/{key} [post]
func AddItems(c echo.Context) (err error) {
	key 	:= c.Param("key")

	log.Println("key : ", key)

	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill key")
	}

	switch key {

	case "uep":
		return AddUep(c)
	case "kube":
		return AddKube(c)
	case "pendamping":
		return AddPendamping(c)
	case "verifikator":
		return AddVerifikator(c)
	case "produk":
		return AddProduk(c)
	case "pelatihan":
		return AddPelatihan(c)
	case "kehadiran":
		return AddPelatihanKehadiran(c)		
	case "log_aktivitas":
		return AddAktivitas(c)
	case "inventaris":
		return AddInventaris(c)
	case "lap_keu":
		return AddLapKeu(c)
	case "monev":
		return AddMonev(c)

	default:
		return echo.NewHTTPError(http.StatusBadRequest, "please, choose the right key")

	}

	return nil
}

// @Summary Update -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev
// @Tags Universal-Controller
// @Accept  json
// @Produce  json
// @Param key path string true "Key -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev "
// @Param sample body models.Dummy true "Update ItemsUepKubePendamping"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /{key} [put]
func UpdateItems(c echo.Context) (err error) {
	key 	:= c.Param("key")

	log.Println("key : ", key)

	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill key")
	}

	switch key {

	case "uep":
		return UpdateUep(c)
	case "kube":
		return UpdateKube(c)
	case "pendamping":
		return UpdatePendamping(c)
	case "verifikator":
		return UpdateVerifikator(c)
	case "produk":
		return UpdateProduk(c)
	case "pelatihan":
		return UpdatePelatihan(c)
	case "log_aktivitas":
		return UpdateAktivitas(c)
	case "inventaris":
		return UpdateInventaris(c)
	case "lap_keu":
		return UpdateLapKeu(c)	
	case "monev":
		return UpdateMonev(c)	
		
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "please, choose the right key")

	}

	return nil
}

// @Summary Delete -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev
// @Tags Universal-Controller
// @Accept  json
// @Produce  json
// @Param key path string true "Key -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev "
// @Param id path int true "Id -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev "
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /{key}/{id}  [post]
func DeleteItems(c echo.Context) (err error) {
	key 	:= c.Param("key")

	log.Println("key : ", key)

	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill key")
	}

	switch key {

	case "uep": // for uep,pendamping,verifikator
		return DeleteUep(c)
	case "kube":
		return DeleteKube(c)
	// case "pendamping":
	// 	return DeletePendamping(c)
	// case "verifikator":
	// 	return DeleteVerifikator(c)	
	case "produk":
		return DeleteProduk(c)
	case "pelatihan":
		return DeletePelatihan(c)
	case "log_aktivitas":
		return DeleteAktivitas(c)
	case "inventaris":
		return DeleteInventaris(c)				
	case "lap_keu":
		return DeleteLapKeu(c)	
	case "monev":
		return DeleteMonev(c)	
	
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "please, choose the right key")

	}

	return nil
}