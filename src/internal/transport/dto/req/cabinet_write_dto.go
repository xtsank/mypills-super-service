package req

type CabinetWriteDto struct {
	Name string `json:"name" required:"true"`
}
