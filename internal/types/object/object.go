package object

type ObjectType int8

const (
	ObjectTypeUser    ObjectType = iota + 1 // user
	ObjectTypePost                          // post
	ObjectTypeComment                       // comment
	ObjectTypeTag                           // tag
)

func (o ObjectType) String() string {
	switch o {
	case ObjectTypeUser:
		return "user"
	case ObjectTypePost:
		return "post"
	case ObjectTypeComment:
		return "comment"
	case ObjectTypeTag:
		return "tag"
	}
	return "unknown"
}

type Object struct {
	ID   int64
	Type ObjectType
}
