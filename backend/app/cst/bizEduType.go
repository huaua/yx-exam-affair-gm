package cst

// 学历
var EduTypeMap = map[string]int{
	"本科": 1,
	"硕士": 2,
	"博士": 3,
}

func EduLabelIsValid(eduType string) bool {
	if _, ok := EduTypeMap[eduType]; ok {
		return true
	}
	return false
}

func EduValueIsValid(eduType int) bool {
	for _, v := range EduTypeMap {
		if v == eduType {
			return true
		}
	}
	return false
}

func GetEduTypeLabelByValue(eduType int) string {
	for label, v := range EduTypeMap {
		if v == eduType {
			return label
		}
	}
	return ""
}
