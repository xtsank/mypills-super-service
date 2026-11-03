package medicine

import (
	"strings"

	"github.com/google/uuid"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
)

func NewCatalogMedicine(id uuid.UUID, name string, expireTime int, isPrescription bool,
	methodOfApplication string, effectOnPregnant, effectOnDriver bool, form, unit uuid.UUID,
	substances []ActiveSubstance, dosages []DosageRule,
	contraindications, recommendations []uuid.UUID) (*Medicine, error) {
	name = strings.TrimSpace(name)
	methodOfApplication = strings.TrimSpace(methodOfApplication)
	if len([]rune(name)) == 0 || len([]rune(name)) > 200 || methodOfApplication == "" ||
		form == uuid.Nil || unit == uuid.Nil || len(substances) == 0 || len(dosages) == 0 ||
		!uniqueCatalogIDs(contraindications) || !uniqueCatalogIDs(recommendations) {
		return nil, apperrors.ErrInvalidInput.WithSource()
	}
	substanceIDs := make([]uuid.UUID, 0, len(substances))
	for _, item := range substances {
		substanceIDs = append(substanceIDs, item.ID)
	}
	if !uniqueCatalogIDs(substanceIDs) {
		return nil, apperrors.ErrInvalidInput.WithSource()
	}
	for _, item := range dosages {
		if item.Type != ByAge && item.Type != ByWeight {
			return nil, apperrors.ErrInvalidInput.WithSource()
		}
	}
	return NewMedicine(id, name, expireTime, isPrescription, methodOfApplication,
		effectOnPregnant, effectOnDriver, form, unit, substances, dosages,
		contraindications, recommendations)
}

func uniqueCatalogIDs(values []uuid.UUID) bool {
	seen := make(map[uuid.UUID]bool, len(values))
	for _, id := range values {
		if id == uuid.Nil || seen[id] {
			return false
		}
		seen[id] = true
	}
	return true
}
