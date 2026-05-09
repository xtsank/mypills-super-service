package res

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
)

type AdminResDto struct {
	ID                  uuid.UUID               `json:"id"`
	Name                string                  `json:"name"`
	ExpireTime          int                     `json:"expire_time"`
	IsPrescription      bool                    `json:"is_prescription"`
	MethodOfApplication string                  `json:"method_of_application"`
	EffectOnPregnant    bool                    `json:"effect_on_pregnant"`
	EffectOnDriver      bool                    `json:"effect_on_driver"`
	Form                uuid.UUID               `json:"form_id"`
	Unit                uuid.UUID               `json:"unit_id"`
	Recommendations     []uuid.UUID             `json:"recommendations"`
	Contraindications   []uuid.UUID             `json:"contraindications"`
	Substances          []ActiveSubstanceResDto `json:"substances"`
	Dosages             []DosageRuleResDto      `json:"dosages"`
}

type ActiveSubstanceResDto struct {
	ID            uuid.UUID `json:"id"`
	Concentration float32   `json:"concentration"`
}

type DosageRuleResDto struct {
	ID                  uuid.UUID        `json:"id"`
	ValueFrom           int              `json:"value_from"`
	ValueTo             int              `json:"value_to"`
	Type                DosageTypeResDto `json:"type"`
	DosageValue         float32          `json:"dosage_value"`
	NumberOfDosesPerDay int              `json:"number_of_doses_per_day"`
}

type DosageTypeResDto string

const (
	ByWeight DosageTypeResDto = "weight"
	ByAge    DosageTypeResDto = "age"
)

func NewAdminResDto(m *medicine.Medicine) *AdminResDto {
	substances := make([]ActiveSubstanceResDto, 0, len(m.Substances))
	for _, s := range m.Substances {
		substances = append(substances, ActiveSubstanceResDto{
			ID:            s.ID,
			Concentration: s.Concentration,
		})
	}

	dosages := make([]DosageRuleResDto, 0, len(m.Dosages))
	for _, d := range m.Dosages {
		dosages = append(dosages, DosageRuleResDto{
			ID:                  d.ID,
			ValueFrom:           d.ValueFrom,
			ValueTo:             d.ValueTo,
			Type:                DosageTypeResDto(d.Type),
			DosageValue:         d.DosageValue,
			NumberOfDosesPerDay: d.NumberOfDosesPerDay,
		})
	}

	return &AdminResDto{
		ID:                  m.ID,
		Name:                m.Name,
		ExpireTime:          m.ExpireTime,
		IsPrescription:      m.IsPrescription,
		MethodOfApplication: m.MethodOfApplication,
		EffectOnPregnant:    m.EffectOnPregnant,
		EffectOnDriver:      m.EffectOnDriver,
		Form:                m.Form,
		Unit:                m.Unit,
		Recommendations:     m.Recommendation,
		Contraindications:   m.Contraindications,
		Substances:          substances,
		Dosages:             dosages,
	}
}
