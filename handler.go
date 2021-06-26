package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func EditHandler(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	est, err := OneEstimationByID(id)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "edit", est)
}

func CreateHandler(c echo.Context) error {

	est := new(Estimation)
	est.Init()

	return c.Render(http.StatusOK, "edit", est)
}

func SaveHandler(c echo.Context) error {

	est := new(Estimation)

	if err := c.Bind(est); err != nil {
		return err
	}

	id := est.ID

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/edit/%d", id))
}

func ErrorHandler(c echo.Context, err error) error {
	log.Error(err)
	return c.String(http.StatusInternalServerError, "Internal Server Error")
}
