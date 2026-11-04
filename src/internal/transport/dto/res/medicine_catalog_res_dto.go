package res

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
)

type MedicineCatalogDto struct {
	ID                  uuid.UUID            `json:"id"`
	Name                string               `json:"name"`
	ExpireTime          int                  `json:"expire_time"`
	IsPrescription      bool                 `json:"is_prescription"`
	MethodOfApplication string               `json:"method_of_application"`
	EffectOnPregnant    bool                 `json:"effect_on_pregnant"`
	EffectOnDriver      bool                 `json:"effect_on_driver"`
	FormID              uuid.UUID            `json:"form_id"`
	UnitID              uuid.UUID            `json:"unit_id"`
	Indications         []uuid.UUID          `json:"indications"`
	Contraindications   []uuid.UUID          `json:"contraindications"`
	Substances          []ActiveSubstanceDto `json:"substances"`
	Dosages             []DosageRuleDto      `json:"dosages"`
}

type ActiveSubstanceDto struct {
	ID            uuid.UUID `json:"id"`
	Concentration float32   `json:"concentration"`
}

type DosageRuleDto struct {
	ID                  uuid.UUID           `json:"id"`
	Type                medicine.DosageType `json:"type"`
	ValueFrom           int                 `json:"value_from"`
	ValueTo             int                 `json:"value_to"`
	DosageValue         float32             `json:"dosage_value"`
	NumberOfDosesPerDay int                 `json:"number_of_doses_per_day"`
}

type MedicineCatalogPageDto struct {
	Items  []MedicineCatalogDto `json:"items"`
	Total  int                  `json:"total"`
	Limit  int                  `json:"limit"`
	Offset int                  `json:"offset"`
}

func NewMedicineCatalogDto(item *medicine.Medicine) *MedicineCatalogDto {
	result := &MedicineCatalogDto{
		ID: item.ID, Name: item.Name, ExpireTime: item.ExpireTime,
		IsPrescription: item.IsPrescription, MethodOfApplication: item.MethodOfApplication,
		EffectOnPregnant: item.EffectOnPregnant, EffectOnDriver: item.EffectOnDriver,
		FormID: item.Form, UnitID: item.Unit,
		Indications:       append([]uuid.UUID{}, item.Recommendation...),
		Contraindications: append([]uuid.UUID{}, item.Contraindications...),
		Substances:        make([]ActiveSubstanceDto, 0, len(item.Substances)),
		Dosages:           make([]DosageRuleDto, 0, len(item.Dosages)),
	}
	for _, substance := range item.Substances {
		result.Substances = append(result.Substances, ActiveSubstanceDto{ID: substance.ID, Concentration: substance.Concentration})
	}
	for _, dosage := range item.Dosages {
		result.Dosages = append(result.Dosages, DosageRuleDto{
			ID: dosage.ID, Type: dosage.Type, ValueFrom: dosage.ValueFrom, ValueTo: dosage.ValueTo,
			DosageValue: dosage.DosageValue, NumberOfDosesPerDay: dosage.NumberOfDosesPerDay,
		})
	}
	return result
}

func NewMedicineCatalogPageDto(page medicine.Page) *MedicineCatalogPageDto {
	result := &MedicineCatalogPageDto{Items: make([]MedicineCatalogDto, 0, len(page.Items)), Total: page.Total, Limit: page.Limit, Offset: page.Offset}
	for _, item := range page.Items {
		result.Items = append(result.Items, *NewMedicineCatalogDto(item))
	}
	return result
}
