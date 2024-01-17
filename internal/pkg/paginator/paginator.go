package paginator

import (
	"gorm.io/gorm"
)

const (
	maxPageSize = 500
	minPageSize = 10
)

func NewPageReply(num int64, size int64, total int64) *PageRely {
	reply := PageRely{
		Size:  size,
		Total: total,
	}
	pages := total / size
	if total%size != 0 {
		pages++
	}
	if num > pages {
		num = pages
	}
	reply.Current = num
	return &reply
}

// 标准分页结构体，接收最原始的DO
// 建议在外部再建一个字段一样的结构体，用以将DO转换成DTO或VO
type Page[T any] struct {
	CurrentPage int   `json:"currentPage"`
	PageSize    int   `json:"pageSize"`
	Total       int64 `json:"total"`
	Data        []T   `json:"data"`
}

func NewPage[T any](currentPage, pageSize int) *Page[T] {
	return &Page[T]{
		CurrentPage: currentPage,
		PageSize:    pageSize,
	}
}

func NewPageFromReq[T any](req *PageRequest) *Page[T] {
	return &Page[T]{
		CurrentPage: int(req.Num),
		PageSize:    int(req.Size),
	}
}

// 各种查询条件先在query设置好后再放进来
func (page *Page[T]) SelectPages(query *gorm.DB) error {
	var model T
	err := query.Model(&model).Count(&page.Total).Error
	if err != nil {
		return err
	}
	if page.Total == 0 {
		page.Data = []T{}
		return nil
	}

	return query.Model(&model).Scopes(Paginate(page)).Find(&page.Data).Error
}

func Paginate[T any](page *Page[T]) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page.CurrentPage <= 0 {
			page.CurrentPage = 1
		}
		switch {
		case page.PageSize > maxPageSize:
			page.PageSize = maxPageSize // 限制一下分页大小
		case page.PageSize < minPageSize:
			page.PageSize = minPageSize
		}
		pages := page.Total / int64(page.PageSize)
		if page.Total%int64(page.PageSize) != 0 {
			pages++
		}
		p := page.CurrentPage
		if page.CurrentPage > int(pages) {
			p = int(pages)
		}
		size := page.PageSize
		offset := int((p - 1) * size)

		return db.Offset(offset).Limit(int(size))
	}
}
