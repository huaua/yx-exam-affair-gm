package cst

// 专家类型
var ExpertTypeMap = map[string]int{
	"本科":  1,
	"研究生": 2,
	"博士":  3,
	"附中": 4,
}

func ExpertLabelIsValid(expertType string) bool {
	if _, ok := ExpertTypeMap[expertType]; ok {
		return true
	}
	return false
}

func ExpertValueIsValid(expertType int) bool {
	for _, v := range ExpertTypeMap {
		if v == expertType {
			return true
		}
	}
	return false
}

func GetExpertTypeLabelByValue(expertType int) string {
	for label, v := range ExpertTypeMap {
		if v == expertType {
			return label
		}
	}
	return ""
}
