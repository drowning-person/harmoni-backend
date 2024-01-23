package v1

func (t ObjectType) Format() string {
	switch t {
	case ObjectType_OBJECT_TYPE_USER:
		return "user"
	case ObjectType_OBJECT_TYPE_POST:
		return "post"
	case ObjectType_OBJECT_TYPE_COMMENT:
		return "comment"
	case ObjectType_OBJECT_TYPE_TAG:
		return "tag"
	}
	return "unspecified"
}

func (o *Object) GetTypeStr() string {
	return o.Type.Format()
}
