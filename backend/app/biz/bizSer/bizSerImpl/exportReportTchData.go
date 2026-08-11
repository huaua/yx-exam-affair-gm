package bizSerImpl

import (
	"strconv"
	"time"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/sealsee/web-base/public/basemodel"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/ds/query"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

/*
 * desc: 征集老师数据导出
 * author: liyu
 * date: 2024-12-23
 */
func exportReportTchData(q *bizDto.EaTchTaskReportQuery) ([]byte, error) {

	// 获取数据列表
	listAll := make([]*bizDto.EaTchTaskReportResultDTO, 0)

	sql := `select a.report_id, a.tch_task_id, t.tch_task_name, a.task_detail_id, a.task_time, a.tch_id, 
			 	b.tch_name, b.job_no, b.id_card, a.remark, a.dept_id, d.dept_name, a.state, b.sex, b.duties, b.prof_office,
			 	(select count(1) from ea_arrange_tch ar where ar.deleted = 1 and ar.tch_id = a.tch_id) as invigilation_count
			from ea_tch_task_report a 
				left join ea_tch_task_info t on t.tch_task_id = a.tch_task_id
				left join ea_tch_info b on a.tch_id = b.tch_id
				left join ea_dept_info d on d.dept_id = a.dept_id
			where `
	condition, args := (&basemodel.BaseEntityQuery{}).
		AddEq("a.tch_task_id", q.TchTaskId).
		AddEq("a.task_time", q.TaskTime).
		AddEq("a.deleted", 1).AddEq("b.deleted", 1).
		AddEq("a.task_detail_id", q.TaskDetailId).
		AddEq("a.dept_id", q.DeptId).
		AddLikeAll("b.job_no", q.JobNo).
		AddLikeAll("b.tch_name", q.TchName).
		GetConditionArgs()

	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE
	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := query.RawSqlQueryListWithPage[bizDto.EaTchTaskReportResultDTO](pageInfo, sql+condition, args...)
		if len(list) == 0 {
			break
		}
		listAll = append(listAll, list...)
		pageInfo.CurPage++
		firstPage = false
	}

	// 导出数据封装
	// 表头数据切片
	data := [][]interface{}{
		{"任务名称", "姓名", "性别", "工号", "身份证号", "征集日期", "备注", "所属学院", "监考次数"},
	}
	for _, tchReport := range listAll {
		data = append(data, []interface{}{
			tchReport.TchTaskName,
			tchReport.TchName,
			tchReport.Sex,
			tchReport.JobNo,
			tchReport.IdCard,
			tchReport.TaskTime.FormatString(time.DateOnly),
			tchReport.Remark,
			tchReport.DeptName,
			tchReport.InvigilationCount,
		})
	}
	dataRowNum := len(data) - 1
	// 处理生成Excel文件
	f := excelize.NewFile()
	defer f.Close()
	var (
		err          error
		addr         string
		sheet        = "Sheet1"
		size         = 9           // 尺寸A4
		orientation  = "landscape" // portrait 纵向 默认 , landscape 横向
		fitToWidth   = 1           // 将所有列调整为一页,高度自适应
		marginTop    = 0.5
		marginBottom = 0.5
		marginLeft   = 0.3
		marginRight  = 0.3
		dataStyle    int
		titleStyle   int
		top          = excelize.Border{Type: "top", Style: 1, Color: "000000"}
		left         = excelize.Border{Type: "left", Style: 1, Color: "000000"}
		right        = excelize.Border{Type: "right", Style: 1, Color: "000000"}
		bottom       = excelize.Border{Type: "bottom", Style: 1, Color: "000000"}
	)
	// 设置页边距
	if err = f.SetPageMargins(sheet, &excelize.PageLayoutMarginsOptions{
		Top:    &marginTop,
		Bottom: &marginBottom,
		Left:   &marginLeft,
		Right:  &marginRight,
	}); err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 页面布局
	if err := f.SetPageLayout(sheet, &excelize.PageLayoutOptions{
		Size:        &size,
		Orientation: &orientation,
		FitToWidth:  &fitToWidth,
	}); err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 设置页眉页脚
	if err = f.SetHeaderFooter(sheet, &excelize.HeaderFooterOptions{
		OddFooter: "&R&P / &N",
	}); err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 按行赋值
	for r, row := range data {
		if addr, err = excelize.JoinCellName("A", r+1); err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
		if err = f.SetSheetRow(sheet, addr, &row); err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}
	}
	// 表格样式
	// 设置表格样式-边框、居中
	if dataStyle, err = f.NewStyle(&excelize.Style{
		Border:    []excelize.Border{top, left, right, bottom},
		Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "center"},
	}); err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	if err = f.SetCellStyle(sheet, "A1", "I"+strconv.Itoa(dataRowNum+1), dataStyle); err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	// 设置表头-边框、加粗、行高
	if titleStyle, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Border:    []excelize.Border{top, left, right, bottom},
		Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "center"},
	}); err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	if err = f.SetCellStyle(sheet, "A1", "I1", titleStyle); err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	// 输出流
	buffer, err := f.WriteToBuffer()
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return buffer.Bytes(), nil
}
