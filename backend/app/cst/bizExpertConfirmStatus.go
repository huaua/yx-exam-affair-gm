package cst

const (
	NOT_CONFIRM = "未应答"
	ATTEND      = "确认参加"
	UNATTEND    = "不参加"
)

// 专家类型
var ExpertConfirmStatusMap = map[string]int{
	NOT_CONFIRM: 1,
	ATTEND:      2,
	UNATTEND:    3,
}

func ExpertConfirmStatusValueIsValid(confirmState int) bool {
	for _, v := range ExpertConfirmStatusMap {
		if v == confirmState {
			return true
		}
	}
	return false
}
