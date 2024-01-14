package object

type ObjectType int8

const (
	ObjectTypePost    ObjectType = iota + 1 // 帖子
	ObjectTypeComment                       // 评论
	ObjectTypeTag                           // 话题
)

func (o ObjectType) String() string {
	switch o {
	case ObjectTypePost:
		return "帖子"
	case ObjectTypeComment:
		return "评论"
	case ObjectTypeTag:
		return "话题"
	}
	return "unknown"
}
