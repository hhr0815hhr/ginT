package logic

type Pagination struct {
	Page  int `json:"page,omitempty" form:"page"`   // 当前页
	Limit int `json:"limit,omitempty" form:"limit"` // 每页数量
}
