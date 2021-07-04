package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type ExportTmpl interface {
	Export(est *Estimation) (string, error)
	configs(est *Estimation) []sheetConfig
	filename() string
}

type sheetConfig struct {
	fromSheet  string
	duplicate  bool
	maxRow     int
	befmap     map[string]interface{}
	aftmap     map[string]interface{}
	mergeRows  []mergeRow
	deleteRows deleteRow
}

type mergeRow struct {
	from string
	to   string
}

type deleteRow struct {
	start int
	end   int
}

type template1 struct{}
type template2 struct{}

var outfileDir = "outfiles/"

func GetTemplate(est *Estimation) ExportTmpl {

	switch est.TemplateID {
	case 1:
		return template1{}
	}

	return nil

}

func (t template1) Export(est *Estimation) (string, error) {

	n := time.Now().Format("20060102_150405")
	outFileName := fmt.Sprintf("%s-%s-%s", est.ClientName, est.EstimationName, n)
	outFileExcel := outfileDir + outFileName + ".xlsx"
	outFilePdf := outfileDir + outFileName + ".pdf"

	f, err := excelize.OpenFile("files/" + t.filename())
	if err != nil {
		return outFilePdf, err
	}

	if err := apllyExcel(f, t.configs(est)); err != nil {
		return outFilePdf, err
	}

	f.DeleteSheet("template1")
	f.DeleteSheet("template2")
	f.DeleteSheet("template3")
	f.DeleteSheet("template4")

	if err := f.SaveAs(outFileExcel); err != nil {
		return outFilePdf, err
	}

	if err := exportPdf(outFileExcel); err != nil {
		return outFilePdf, err
	}

	return outFilePdf, err
}

func apllyExcel(f *excelize.File, sc []sheetConfig) error {

	for i, sheet := range sc {

		from := sheet.fromSheet
		to := strconv.Itoa(i)

		if sheet.duplicate {
			f.NewSheet(to)
			if err := f.CopySheet(f.GetSheetIndex(from), f.GetSheetIndex(to)); err != nil {
				return err
			}
		} else {
			f.SetSheetName(from, to)
		}

		// 先行して配置しておきたい値を入力
		for cel, val := range sheet.befmap {
			if err := f.SetCellValue(to, cel, val); err != nil {
				return err
			}
		}
		// 後から入力する値を入力
		for cel, val := range sheet.aftmap {
			if err := f.SetCellValue(to, cel, val); err != nil {
				return err
			}
		}
		// 対象となっているセルをマージ
		for _, mgrw := range sheet.mergeRows {
			if err := f.MergeCell(to, mgrw.from, mgrw.to); err != nil {
				return err
			}
		}
		// 対象となっている行を削除
		for i := sheet.deleteRows.start; i < sheet.deleteRows.end; i++ {
			if err := f.RemoveRow(to, sheet.deleteRows.start); err != nil {
				return err
			}
		}

	}
	return nil
}

func exportPdf(convFile string) error {
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
	return exec.Command(command, "--headless", "--convert-to", "pdf", "--outdir", approot+outfileDir, approot+convFile).Run()
}

func calcRowNum(groups []Group) int {
	ret := calcItemCount(groups) + len(groups) - 1
	return ret
}

func calcItemCount(groups []Group) int {
	ret := 0
	for _, group := range groups {
		ret += len(group.Items)
	}
	return ret
}

func calcSubtotal(groups []Group) int {
	ret := 0
	for _, group := range groups {
		for _, item := range group.Items {
			ret += item.Price
		}
	}
	return ret
}

func convGroups(groups []Group, gNameCl, iNameCl, iAmountCl, iUnitCl, iUnitPriceCl, iPriceCl string,
	row, maxrow int) (aftmap map[string]interface{}, mgrw []mergeRow, dlrw deleteRow) {

	aftmap = map[string]interface{}{}

	for i, group := range groups {
		strow := strconv.Itoa(row)
		aftmap[gNameCl+strow] = group.Name

		for _, item := range group.Items {
			strow = strconv.Itoa(row)
			aftmap[iNameCl+strow] = item.Name
			aftmap[iAmountCl+strow] = item.Amount
			aftmap[iUnitCl+strow] = item.Unit
			aftmap[iUnitPriceCl+strow] = item.UnitPrice
			aftmap[iPriceCl+strow] = item.Price
			row += 1
		}

		if i != len(groups)-1 {
			strow = strconv.Itoa(row)
			mgrw = append(mgrw, mergeRow{
				from: gNameCl + strow,
				to:   iPriceCl + strow,
			})
			row += 1
		}
	}

	dlrw.start = row
	dlrw.end = maxrow + 1

	return
}

func (t template1) configs(est *Estimation) []sheetConfig {

	sheetConfigs := []sheetConfig{}
	doneGroups := []Group{}

	for {

		if len(sheetConfigs) == 0 {
			if calcRowNum(est.Groups) < 26 {
				sc, dg := t.singlePage(est)
				sheetConfigs = append(sheetConfigs, sc)
				doneGroups = append(doneGroups, dg...)
			} else {
				sc, dg := t.multiPage(est, doneGroups)
				sheetConfigs = append(sheetConfigs, sc)
				doneGroups = append(doneGroups, dg...)
			}
		} else {
			sc, dg := t.multiPage(est, doneGroups)
			sheetConfigs = append(sheetConfigs, sc)
			doneGroups = append(doneGroups, dg...)
		}

		if calcRowNum(est.Groups)-calcRowNum(doneGroups) < 1 {
			break
		}
	}

	return sheetConfigs
}

func (template1) filename() string {
	return "template.xlsx"
}

func (template1) singlePage(est *Estimation) (sheetConfig, []Group) {

	befmap := map[string]interface{}{
		"A5":  est.ClientName,
		"B10": est.Total,
		"J40": est.SubTotal,
		"J41": est.Tax,
		"J43": est.Total,
	}

	aftmap, mergerows, deleterows := convGroups(est.Groups, "A", "C", "G", "H", "I", "J", 14, 39)

	return sheetConfig{
		fromSheet:  "template1",
		duplicate:  false,
		maxRow:     40,
		befmap:     befmap,
		aftmap:     aftmap,
		mergeRows:  mergerows,
		deleteRows: deleterows,
	}, est.Groups
}

func (t template1) multiPage(est *Estimation, dones []Group) (sheetConfig, []Group) {

	var fromSheet string
	var duplicate bool
	var maxRow int
	var befmap map[string]interface{}
	var aftmap map[string]interface{}
	var startRow int
	var mergeRows []mergeRow
	var deleteRows deleteRow

	var max int
	var subtotalCell string

	row := calcRowNum(est.Groups)
	donerow := calcRowNum(dones)

	if len(dones) == 0 {
		fromSheet = "template2"
		duplicate = false
		maxRow = 42
		befmap = map[string]interface{}{
			"A5":  est.ClientName,
			"B10": est.Total,
		}
		startRow = 14
		max = 29
		subtotalCell = "J43"
	} else if row-donerow <= 34 {
		fromSheet = "template4"
		duplicate = true
		maxRow = 36
		befmap = map[string]interface{}{
			"J38": est.Tax,
			"J40": est.Total,
		}
		startRow = 3
		max = 34
		subtotalCell = "J37"
	} else {
		fromSheet = "template3"
		duplicate = true
		maxRow = 39
		befmap = map[string]interface{}{}
		startRow = 3
		max = 37
		subtotalCell = "J40"
	}

	targetGroups := targetGroup(est, dones, max)
	subtotal := calcSubtotal(targetGroups)
	befmap[subtotalCell] = subtotal

	aftmap, mergeRows, deleteRows = convGroups(targetGroups, "A", "C", "G", "H", "I", "J", startRow, maxRow)

	return sheetConfig{
		fromSheet:  fromSheet,
		duplicate:  duplicate,
		maxRow:     maxRow,
		befmap:     befmap,
		aftmap:     aftmap,
		mergeRows:  mergeRows,
		deleteRows: deleteRows,
	}, targetGroups

}

func targetGroup(est *Estimation, dones []Group, max int) []Group {

	targets := []Group{}

	for i, group := range est.Groups {

		if checkDone(group, dones) {
			continue
		}

		// 絶対収まりきらない場合は改ページを許容する
		if len(group.Items) > max {
			est.SeparateGroup(i, max-calcItemCount(targets))
			targets = append(targets, est.Groups[i])
			break
		}

		if calcRowNum(targets)+len(group.Items)+1 < max {
			targets = append(targets, group)
		} else {
			break
		}
	}

	return targets
}

func checkDone(gr Group, dones []Group) bool {

	for _, done := range dones {
		// TODO: この処理のせいで、2回以上のセパレートはできなくなってる
		if gr.ID == done.ID {
			return true
		}
	}
	return false
}
