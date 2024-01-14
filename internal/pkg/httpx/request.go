package httpx

type Page struct {
	PageNum  int `json:"num" form:"pn"`  // page number
	PageSize int `json:"size" form:"ps"` // page size
}
