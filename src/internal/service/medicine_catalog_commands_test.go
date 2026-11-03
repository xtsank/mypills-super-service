package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
)

func TestMedicinePatchKeepsUnchangedClinicalData(t *testing.T) {
	ruleID := uuid.New()
	current := &medicine.Medicine{
		ID: uuid.New(), Name: "Лекарство", ExpireTime: 12, MethodOfApplication: "Внутрь",
		Form: uuid.New(), Unit: uuid.New(),
		Recommendation: []uuid.UUID{uuid.New()},
		Substances:     []medicine.ActiveSubstance{{ID: uuid.New(), Concentration: 10}},
		Dosages:        []medicine.DosageRule{{ID: ruleID, Type: medicine.ByAge, ValueFrom: 18, ValueTo: 120, DosageValue: 1, NumberOfDosesPerDay: 2}},
	}
	expireTime := 24
	updated, err := patchCatalogMedicine(current, command.MedicinePatchCmd{ID: current.ID, ExpireTime: &expireTime})
	require.NoError(t, err)
	require.Equal(t, 24, updated.ExpireTime)
	require.Equal(t, current.Recommendation, updated.Recommendation)
	require.Equal(t, current.Substances, updated.Substances)
	require.Equal(t, ruleID, updated.Dosages[0].ID)
}
