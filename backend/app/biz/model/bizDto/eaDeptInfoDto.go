package bizDto

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 学院基础数据查询实体
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaDeptInfoQuery struct {
	DeptId     int64  `gorm:"primaryKey" form:"deptId" json:"deptId,string,omitempty"` // 主键 学院id
	DeptCode   string `form:"deptCode" json:"deptCode,omitempty"`                      // 学院代码
	DeptName   string `form:"deptName" json:"deptName,omitempty"`                      // 学院名称
	CampusName string `form:"campusName" json:"campusName,omitempty"`                  // 所在校区
	Remark     string `form:"remark" json:"remark,omitempty"`                          // 备注
	basemodel.BaseEntityQuery
}
