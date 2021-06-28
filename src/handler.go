package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func GetEstimation(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	est, err := OneEstimationByID(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, est)
}

func EditPage(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	_, err = OneEstimationByID(uint(id))
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "edit", nil)
}

func CreatePage(c echo.Context) error {

	est := new(Estimation)
	est.Init()

	err := SaveEstimation(est)
	if err != nil {
		return err
	}

	id := est.ID

	return c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/estimation/%d/edit", id))
}

func DuplicateHandler(c echo.Context) error {

	est := new(Estimation)

	if err := c.Bind(est); err != nil {
		return err
	}

	err := SaveEstimation(est)
	if err != nil {
		return err
	}

	cp := est

	cp.ID = 0
	for i := range cp.Groups {
		for j := range cp.Groups[i].Items {
			cp.Groups[i].Items[j].ID = 0
			cp.Groups[i].Items[j].GroupID = 0
		}
		cp.Groups[i].ID = 0
		cp.Groups[i].EstimationID = 0
	}

	err = SaveEstimation(cp)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &RedirectResponse{
		URL: fmt.Sprintf("/estimation/%d/edit", cp.ID),
	})
}

func SaveHandler(c echo.Context) error {

	est := new(Estimation)

	if err := c.Bind(est); err != nil {
		return err
	}

	err := SaveEstimation(est)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, est)
}

func DeleteItemHandler(c echo.Context) error {

	item := new(Item)

	if err := c.Bind(item); err != nil {
		return err
	}

	err := DeleteObject(item)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "ok")
}

func DeleteGroupHandler(c echo.Context) error {

	group := new(Group)

	if err := c.Bind(group); err != nil {
		return err
	}

	err := DeleteObject(group)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "ok")
}

func ErrorHandler(c echo.Context, err error) error {
	log.Error(err)
	return c.String(http.StatusInternalServerError, "Internal Server Error")
}
