package service

type Uri struct {
	ID string `uri:"id" binding:"required,len=24,alphanum"`
}

type AddReq struct {
	Title string `json:"title" binding:"required,lte=10,alphanum"`
}

type UpdateReq struct {
	Status int32 `json:"status" binding:"oneof=0 1 2"`
}
