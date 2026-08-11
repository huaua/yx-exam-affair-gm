package cst

const (
	NOT_START  = "未开始"
	EXTRACTING = "抽取中"
	FINISHED   = "已完成"
)

// 专家类型
var JudgeTaskStatusMap = map[string]int{
	NOT_START:  1,
	EXTRACTING: 2,
	FINISHED:   3,
}

func GetJudgeTaskStatusLabelByValue(status int) string {
	for label, v := range JudgeTaskStatusMap {
		if v == status {
			return label
		}
	}
	return ""
}
