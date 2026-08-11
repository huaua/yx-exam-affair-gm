package bizDo

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 评委评价数据实体
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeExpertsAppraisal struct {
    AppraisalId int64 `gorm:"primaryKey" json:"appraisalId,string,omitempty"` // 主键 评价ID
    ExpertId int64 `json:"expertId,string,omitempty"` // 专家ID
    JudgeTaskId int64 `json:"judgeTaskId,string,omitempty"` // 任务ID
    JudgeTaskName string `json:"judgeTaskName,omitempty"` // 评委任务名称
    TaskTimeRemark string `json:"taskTimeRemark,omitempty"` // 评分日期
    AppraisalContent string `json:"appraisalContent,omitempty"` // 评价内容
    Remark string `json:"remark,omitempty"` // 备注
    basemodel.BaseEntity
}
