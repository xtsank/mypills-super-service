package req

type RestockItemPatchDto struct {
	DesiredQuantity float64 `json:"desired_quantity" required:"true"`
}
