package comment

import "harmoni/app/harmoni/internal/entity/comment"

type CommnetList []*Comment

func (l CommnetList) ToDomain() []*comment.Comment {
	cl := make([]*comment.Comment, 0, len(l))
	for _, comment := range l {
		cl = append(cl, comment.ToDomain())
	}
	return cl
}
