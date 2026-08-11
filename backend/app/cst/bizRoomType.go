package cst

/*
 * desc: 房间类型
 * author: 蓝鲸
 * date: 2024/08/15
 */
type RoomType int

const (
	BIZ_TYPE_GROUP  RoomType = 1 // 按组
	BIZ_TYPE_PERSON RoomType = 2 // 按位
)

var RoomTypeMap = map[RoomType]string{
	BIZ_TYPE_GROUP:  "按组",
	BIZ_TYPE_PERSON: "按位",
}
