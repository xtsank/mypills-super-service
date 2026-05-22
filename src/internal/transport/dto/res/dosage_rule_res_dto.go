package res

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
)

// DosageRuleItemDto используется для удобного выбора правила дозировки в UI.
type DosageRuleItemDto struct {
	ID    uuid.UUID `json:"id"`
	Label string    `json:"label"`
}

func NewDosageRuleItems(items []entity.DosageRuleWithMedicineEntity) []DosageRuleItemDto {
	res := make([]DosageRuleItemDto, 0, len(items))
	for _, item := range items {
		label := fmt.Sprintf("%s: %s %d-%d (%.2f, %d/д)", item.MedicineName, item.DosageType, item.ValueFrom, item.ValueTo, item.DosageValue, item.NumberOfDosesPerDay)
		res = append(res, DosageRuleItemDto{ID: item.ID, Label: label})
	}
	return res
}


