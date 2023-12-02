package action

type Action int8

const (
	ActionNo    Action = iota
	ActionReply        // 回复
	ActionAt           // at
	ActionLike         // 点赞
)
