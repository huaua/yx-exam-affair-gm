package bizDto

import "github.com/sealsee/web-base/public/basemodel"

/*
 * desc: 评委评价数据查询实体
 * author: init code
 * date: 2025/06/11
 */
type EaJudgeExpertsAppraisalQuery struct {
    AppraisalId int64 `gorm:"primaryKey" form:"appraisalId" json:"appraisalId,string,omitempty"` // 主键 评价ID
    ExpertId int64 `form:"expertId" json:"expertId,string,omitempty"` // 专家ID
    JudgeTaskId int64 `form:"judgeTaskId" json:"judgeTaskId,string,omitempty"` // 任务ID
    JudgeTaskName string `form:"judgeTaskName" json:"judgeTaskName,omitempty"` // 评委任务名称
    TaskTimeRemark string `form:"taskTimeRemark" json:"taskTimeRemark,omitempty"` // 评分日期
    AppraisalContent string `form:"appraisalContent" json:"appraisalContent,omitempty"` // 评价内容
    Remark string `form:"remark" json:"remark,omitempty"` // 备注
    basemodel.BaseEntityQuery
}