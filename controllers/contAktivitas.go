package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	 _"github.com/jinzhu/gorm/dialects/mysql"
	 "uepkube-api/db"
	 "uepkube-api/models"
	 "strconv"
	 "uepkube-api/helpers"
)

func GetAktivitas(c echo.Context) error {
	id 		:= c.QueryParam("id")

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	Activity := models.Tbl_activity{}
	q := con
	q = q.Model(&Activity)
	q = q.Preload("Photo")
	q = q.First(&Activity, id)

	r := &models.Jn{Msg: Activity}
	defer con.Close()

	return c.JSON(http.StatusOK, r)
}

func GetPaginateAktivitas(c echo.Context) (err error) {	
	if err := helpers.PaginateAktivitas(c, &r); err != nil {
		return echo.ErrInternalServerError
	}	
	return c.JSON(http.StatusOK, r)
}

func AddAktivitas(c echo.Context) (err error) {
	activity := &models.Activity{}

	if err := c.Bind(activity); err != nil {
		return err
	}

	Activity := &models.Tbl_activity{}
	Activity = activity.Tbl_activity

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Create(&Activity).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	defer con.Close()

	r := &models.Jn{Msg: "Success Store Data"}
	return c.JSON(http.StatusOK, r)
}

func UpdateAktivitas(c echo.Context) (err error) {
	activity := &models.Activity{}

	if err := c.Bind(activity); err != nil {
		return err
	}

	if activity.Id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id")
	}

	Activity := &models.Tbl_activity{}
	Activity = activity.Tbl_activity

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Model(&models.Tbl_activity{}).UpdateColumns(&Activity).Error; err != nil {
		return echo.ErrInternalServerError
	}

	defer con.Close()

	r := &models.Jn{Msg: "Success Update Data"}
	return c.JSON(http.StatusOK, r)
}

func DeleteAktivitas(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("id"))

	if id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	}

	activity := &models.Tbl_activity{}
	activity.Id = id

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Delete(&activity).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	defer con.Close()

	r := &models.Jn{Msg: "Success Delete Data"	}
	return c.JSON(http.StatusOK, r)	
}