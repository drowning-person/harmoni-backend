package object

type ObjectType int8

const (
	ObjectTypePost    ObjectType = iota + 1 // 帖子
	ObjectTypeComment                       // 评论
	ObjectTypeTag                           // 话题
)
