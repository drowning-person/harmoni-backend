package action

type Action int8

const (
	ActionNo    Action = iota
	ActionReply        // 回复
	ActionAt           // at
	ActionLike         // 点赞
)

func (a Action) String() string {
	switch a {
	case ActionNo:
		return ""
	case ActionReply:
		return "回复"
	case ActionAt:
		return "提到"
	case ActionLike:
		return "喜欢"
	}
	return "unknown"
}
