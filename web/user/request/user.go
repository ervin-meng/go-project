package request

type ListRequest struct {
	Page     uint32 `form:"page" binding:"required"`
	PageSize uint32 `form:"pagesize" binding:"required"`
}
