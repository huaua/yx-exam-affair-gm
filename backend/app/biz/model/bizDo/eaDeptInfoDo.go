package bizDo

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 学院基础数据实体
 * author: 蓝鲸
 * date: 2024/08/14
 */
type EaDeptInfo struct {
	DeptId     int64  `gorm:"primaryKey" json:"deptId,string,omitempty"` // 主键 学院id
	DeptCode   string `json:"deptCode,omitempty"`                        // 学院代码
	DeptName   string `json:"deptName,omitempty"`                        // 学院名称
	CampusName string `json:"campusName,omitempty"`                      // 所在校区
	Remark     string `json:"remark,omitempty"`                          // 备注
	basemodel.BaseEntity
}
