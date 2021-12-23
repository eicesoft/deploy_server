package controller

//var _ Handler = (*handler)(nil)

const (
	DefaultPageSize = 10
	DefaultPage     = 1
)

type IdRequest struct {
	Id int32 `form:"id" query:"id" json:"id" uri:"id" binding:"required" :"id"`
}

type SimpleResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//
//type handler struct {
//	Db db.Repo
//}
//
//type Handler interface {
//	GetDBReader(ctx context.Context) *gorm.DB
//	GetDBWriter(ctx context.Context) *gorm.DB
//}
//
//func (h handler) GetDBWriter(ctx context.Context) *gorm.DB {
//	dbr := h.Db.GetDbW()
//	dbr.WithContext(ctx)
//	return dbr
//}
//
//func (h handler) GetDBReader(ctx context.Context) *gorm.DB {
//	dbr := h.Db.GetDbR()
//	dbr.WithContext(ctx)
//	return dbr
//}
