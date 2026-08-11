package bizSerImpl

import (
	"github.com/sealsee/web-base/public/ds/page"
	"strconv"
	"time"
	"yx-exam-affair-gm/app/biz/model/bizDto"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

/*
 * desc: 通知书打印数据导出
 * author: 蓝鲸
 * date: 2024-09-30
 */
func exportArrangeData(q *bizDto.EaArrangeTchQuery) ([]byte, error) {

	// 获取数据列表
	listAll := make([]*bizDto.EaArrangeProfResultDto, 0)
	pageInfo := page.NewPage()
	pageInfo.PageSize = page.EXPORT_PAGE_SIZE

	firstPage := true
	for pageInfo.CurPage <= pageInfo.TotalPage || firstPage {
		list := eaArrangeTchService.ListEaArrangeTchByProf(q, pageInfo)
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
		{"学院名称", "专业", "方向", "考试时间", "考场人数", "考场", "编号", "监考老师数量", "监考老师"},
	}
	for _, arrangeTchInfo := range listAll {
		data = append(data, []interface{}{
			arrangeTchInfo.DeptName,
			arrangeTchInfo.ProfName,
			arrangeTchInfo.DirectionName,
			time.Time(arrangeTchInfo.TaskDay).Format("2006-01-02"),
			arrangeTchInfo.StuNum,
			arrangeTchInfo.RoomName,
			arrangeTchInfo.ArrangeNoStr,
			arrangeTchInfo.TchNum,
			arrangeTchInfo.ArrangeTchName,
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
