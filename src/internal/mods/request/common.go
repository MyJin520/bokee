package request

type PageReq struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

func (p *PageReq) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

func (p *PageReq) Offset() int {
	return (p.Page - 1) * p.PageSize
}
