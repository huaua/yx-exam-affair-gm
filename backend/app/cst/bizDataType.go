package cst

type DataType int

const (
	BIZ_ROOM_BASE_INFO    DataType = 1 // 教室基础信息
	BIZ_TCH_BASE_INFO     DataType = 2 // 教职工基础信息
	BIZ_ARRANGE_PROF_INFO DataType = 3 // 专业方向信息
	BIZ_JUDGE_EXPERTS     DataType = 4 // 评委专家库信息
)

var DataTypeMap = map[DataType]string{
	BIZ_ROOM_BASE_INFO:    "教室基础信息",
	BIZ_TCH_BASE_INFO:     "教职工基础信息",
	BIZ_ARRANGE_PROF_INFO: "专业方向信息",
}
