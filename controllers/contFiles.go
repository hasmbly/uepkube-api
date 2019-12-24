package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	// "github.com/jinzhu/gorm"
	 // _"github.com/jinzhu/gorm/dialects/mysql"
	 // "uepkube-api/models"
	 // "uepkube-api/db"
	//  "regexp"
	//  "uepkube-api/helpers"
	 "log"	 
)

// @Summary UploadsFiles -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev
// @Tags UploadFiles-Controller
// @Accept  mpfd
// @Produce  mpfd
// @Param key path string true "Key (string) -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev"
// @Param id query int true "id (int)"
// @Param files formData file true "Uploads Files"
// @Param description formData string false "Uploads Files"
// @Param type formData string true "enums : 'IMAGE', 'PDF' "
// @Param is_display query int false "int (int) -> 0 | 1"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /upload/files/{key} [post]
func UploadFiles(c echo.Context) (err error) {
	key 	:= c.Param("key")

	log.Println("Uploads File to : ", key)

	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill key")
	}

	switch key {
		case "uep":
			return UploadUepFiles(c)
		case "kube":
			return UploadKubeFiles(c)
		case "log_aktivitas":
			return UploadAktivitasFiles(c)			
		case "produk":
			return UploadProdukFiles(c)
		case "pelatihan":
			return UploadPelatihanFiles(c)					
		case "lap_keu":
			return UploadLapKeuFiles(c)
		case "inventaris":
			return UploadInventarisFiles(c)
		case "pendamping":
			return UploadPendampingFiles(c)
		case "verifikator":
			return UploadVerifikatorFiles(c)
		default:
			return echo.NewHTTPError(http.StatusBadRequest, "please, choose the right key")
	}

	return nil
	
}

// @Summary UploadsFiles -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev
// @Tags DownloadFiles-Controller
// @Accept  json
// @Produce  json
// @Param key path string true "Key (string) -> uep | kube | pendamping | verifikator | produk | pelatihan | log_aktivitas | inventaris | lap_keu | kehadiran | monev"
// @Param id query int true "id (int)"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /download/files/{key} [get]
func DownloadFiles(c echo.Context) (err error) {
	key 	:= c.Param("key")

	log.Println("Download File to : ", key)

	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill key")
	}

	switch key {
		case "pelatihan":
			return DownloadPelatihanFiles(c)					
		// case "uep":
		// 	return UploadUepFiles(c)
		// case "kube":
		// 	return UploadKubeFiles(c)
		// case "log_aktivitas":
		// 	return UploadAktivitasFiles(c)			
		// case "produk":
		// 	return UploadProdukFiles(c)
		// case "pelatihan":
		// 	return UploadPelatihanFiles(c)					
		// case "lap_keu":
		// 	return UploadLapKeuFiles(c)
		// case "inventaris":
		// 	return UploadInventarisFiles(c)
		// case "pendamping":
		// 	return UploadPendampingFiles(c)
		// case "verifikator":
		// 	return UploadVerifikatorFiles(c)
		default:
			return echo.NewHTTPError(http.StatusBadRequest, "please, choose the right key")
	}

	return nil
	
}