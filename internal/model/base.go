package model
 
// BaseModel 基础模型: 包含ID、创建时间、更新时间和软删除时间
type BaseModel struct {
	ID        uint64 `json:"id" gorm:"primarykey"`                   // 主键ID
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime:milli"` // 毫秒时间戳
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime:milli"` // 毫秒时间戳
	DeletedAt int64  `json:"deleted_at"`                             // 软删除时间戳，0 表示未删除
}

// PageRequest 分页请求: 包含页码和每页数量
type PageRequest struct {
	Page int `json:"page" form:"page" binding:"min=1"`
	Size int `json:"size" form:"size" binding:"min=1,max=100"`
}

func (p *PageRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.Size
}

func (p *PageRequest) GetLimit() int {
	if p.Size <= 0 {
		p.Size = 10
	}
	if p.Size > 100 {
		p.Size = 100
	}
	return p.Size
}
