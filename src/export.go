package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func ExportEstimation(est *Estimation) (string, error) {

	n := time.Now().Format("20060102_150405")
	outFileName := fmt.Sprintf("files/%s-%s-%s", est.ClientName, est.EstimationName, n)
	outFileExcel := outFileName + ".xlsx"
	outFilePdf := outFileName + ".pdf"

	f, err := excelize.OpenFile("files/template.xlsx")
	if err != nil {
		return outFilePdf, err
	}

	rowNum := calcRowNum(est)

	if rowNum <= 26 {
		createSinglePage(f, est)
	}

	if err := f.SaveAs(outFileExcel); err != nil {
		return outFilePdf, err
	}

	var command string
	var approot string
	env := os.Getenv("ENV")
	if env == "LOCAL" {
		command = "C:\\LibreOffice\\program\\soffice"
		approot = "C:\\Users\\s.mine\\dev\\shimo_app\\"
	} else {
		command = "/usr/bin/soffice"
		approot = "/app/"
	}

	err = exec.Command(command, "--headless", "--convert-to", "pdf", "--outdir", approot+"files", approot+outFileExcel).Run()
	if err != nil {
		return outFilePdf, err
	}

	return outFilePdf, err
}

func calcRowNum(est *Estimation) int {

	ret := 0

	for i, group := range est.Groups {
		ret += len(group.Items)
		if i != len(est.Groups)-1 {
			ret += 1
		}
	}
	return ret
}

func createSinglePage(f *excelize.File, est *Estimation) error {

	page := "page"

	f.SetSheetName("template1", page)

	// client
	if err := f.SetCellValue(page, "A5", est.ClientName); err != nil {
		return err
	}

	// total on head
	if err := f.SetCellValue(page, "B10", est.Total); err != nil {
		return err
	}

	// subtotal
	if err := f.SetCellValue(page, "J40", est.SubTotal); err != nil {
		return err
	}

	// tax
	if err := f.SetCellValue(page, "J41", est.Tax); err != nil {
		return err
	}

	// total on bottom
	if err := f.SetCellValue(page, "J43", est.Total); err != nil {
		return err
	}

	row, err := setGroupValue(f, page, est.Groups, 14)
	if err != nil {
		return err
	}

	for i := row; i < 40; i++ {
		if err := f.RemoveRow(page, row); err != nil {
			return err
		}
	}

	f.DeleteSheet("template1")
	f.DeleteSheet("template2")
	f.DeleteSheet("template3")
	f.DeleteSheet("template4")

	return nil
}

func setGroupValue(f *excelize.File, sn string, groups []Group, row int) (int, error) {

	for i, group := range groups {

		groupNameCell := fmt.Sprintf("A%d", row)

		// group name
		if err := f.SetCellValue(sn, groupNameCell, group.Name); err != nil {
			return row, err
		}

		for _, item := range group.Items {

			nameCell := fmt.Sprintf("C%d", row)
			amountCell := fmt.Sprintf("G%d", row)
			unitCell := fmt.Sprintf("H%d", row)
			unitPriceCell := fmt.Sprintf("I%d", row)
			priceCell := fmt.Sprintf("J%d", row)

			// name
			if err := f.SetCellValue(sn, nameCell, item.Name); err != nil {
				return row, err
			}

			// amount
			if err := f.SetCellValue(sn, amountCell, item.Amount); err != nil {
				return row, err
			}

			// unit
			if err := f.SetCellValue(sn, unitCell, item.Unit); err != nil {
				return row, err
			}

			// unit price
			if err := f.SetCellValue(sn, unitPriceCell, item.UnitPrice); err != nil {
				return row, err
			}

			// price
			if err := f.SetCellValue(sn, priceCell, item.Price); err != nil {
				return row, err
			}

			row += 1
		}

		if i != len(groups)-1 {

			mgFrom := fmt.Sprintf("A%d", row)
			mgTo := fmt.Sprintf("J%d", row)

			if err := f.MergeCell(sn, mgFrom, mgTo); err != nil {
				return row, err
			}
			row += 1
		}
	}

	return row, nil
}
