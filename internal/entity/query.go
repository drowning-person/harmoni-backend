package entity

type PageCond struct {
	Page     int64 `query:"page"`
	PageSize int64 `query:"page_size"`
}
