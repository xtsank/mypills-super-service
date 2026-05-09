package medicine

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
)

type Medicine struct {
	ID                  uuid.UUID
	Name                string
	ExpireTime          int
	IsPrescription      bool
	MethodOfApplication string
	EffectOnPregnant    bool
	EffectOnDriver      bool
	Form                uuid.UUID
	Unit                uuid.UUID

	Recommendation    []uuid.UUID
	Contraindications []uuid.UUID

	Substances []ActiveSubstance
	Dosages    []DosageRule
}

type ActiveSubstance struct {
	ID            uuid.UUID
	Concentration float32
}

type DosageRule struct {
	ID                  uuid.UUID
	ValueFrom           int
	ValueTo             int
	Type                DosageType
	DosageValue         float32
	NumberOfDosesPerDay int
}

type DosageType string

const (
	ByWeight DosageType = "weight"
	ByAge    DosageType = "age"
)

func NewMedicine(
	id uuid.UUID,
	name string,
	expireTime int,
	isPrescription bool,
	methodOfApplication string,
	effectOnPregnant bool,
	effectOnDriver bool,
	form uuid.UUID,
	unit uuid.UUID,
	substances []ActiveSubstance,
	dosages []DosageRule,
	contraindications []uuid.UUID,
	recommendations []uuid.UUID,
) (*Medicine, error) {
	if name == "" {
		return nil, errors.ErrEmptyName
	}

	if expireTime <= 0 {
		return nil, errors.ErrExpireTimeTooLow
	}

	for _, s := range substances {
		if s.Concentration <= 0 {
			return nil, errors.ErrInvalidConcentration
		}
	}
	for _, d := range dosages {
		if d.ValueFrom < 0 || d.ValueTo < d.ValueFrom {
			return nil, errors.ErrInvalidDosageRange
		}
		if d.DosageValue <= 0 {
			return nil, errors.ErrInvalidDosageValue
		}
		if d.NumberOfDosesPerDay <= 0 {
			return nil, errors.ErrInvalidNumDoses
		}
	}

	return &Medicine{
		ID:                  id,
		Name:                name,
		ExpireTime:          expireTime,
		IsPrescription:      isPrescription,
		MethodOfApplication: methodOfApplication,
		EffectOnPregnant:    effectOnPregnant,
		EffectOnDriver:      effectOnDriver,
		Form:                form,
		Unit:                unit,
		Substances:          substances,
		Dosages:             dosages,
		Contraindications:   contraindications,
		Recommendation:      recommendations,
	}, nil
}

func (med *Medicine) IsSafeFor(u *user.User) bool {
	if u.IsPregnant && med.EffectOnPregnant {
		return false
	}

	if u.IsDriver && med.EffectOnDriver {
		return false
	}

	if med.hasContraindications(u.Illnesses) {
		return false
	}

	if med.hasAllergies(u.Allergies) {
		return false
	}

	return true
}

func (med *Medicine) hasContraindications(illnesses []uuid.UUID) bool {
	for _, userIllness := range illnesses {
		for _, medicineContr := range med.Contraindications {
			if userIllness == medicineContr {
				return true
			}
		}
	}

	return false
}

func (med *Medicine) hasAllergies(allergies []uuid.UUID) bool {
	for _, userAllergy := range allergies {
		for _, medicineSubstance := range med.Substances {
			if userAllergy == medicineSubstance.ID {
				return true
			}
		}
	}

	return false
}

func (med *Medicine) CalculateDosage(u *user.User) (float32, int) {
	for _, rule := range med.Dosages {
		switch rule.Type {
		case ByWeight:
			if u.Weight >= rule.ValueFrom && u.Weight <= rule.ValueTo {
				return rule.DosageValue, rule.NumberOfDosesPerDay
			}

		case ByAge:
			if u.Age >= rule.ValueFrom && u.Age <= rule.ValueTo {
				return rule.DosageValue, rule.NumberOfDosesPerDay
			}
		}
	}

	return 0, 0
}
