package res

type SuccessResDTO struct {
	Status string `json:"status"`
}

func NewSuccessResDto() *SuccessResDTO {
	return &SuccessResDTO{Status: "success"}
}
