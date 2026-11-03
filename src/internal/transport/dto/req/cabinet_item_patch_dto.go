package req

type CabinetItemPatchDto struct {
	Quantity          *float64 `json:"quantity"`
	MinimumQuantity   *float64 `json:"minimum_quantity"`
	DateOfManufacture *DateDto `json:"date_of_manufacture"`
}
