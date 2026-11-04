package testutil

import (
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
)

type MedicineBuilder struct {
	value medicine.Medicine
}

func NewMedicineBuilder() *MedicineBuilder {
	return &MedicineBuilder{value: medicine.Medicine{
		ID:                  MedicineA,
		Name:                "Base remedy",
		ExpireTime:          24,
		IsPrescription:      true,
		MethodOfApplication: "oral",
		EffectOnDriver:      true,
		Form:                FormA,
		Unit:                UnitA,
	}}
}

func (b *MedicineBuilder) WithID(id uuid.UUID) *MedicineBuilder {
	b.value.ID = id
	return b
}

func (b *MedicineBuilder) WithName(name string) *MedicineBuilder {
	b.value.Name = name
	return b
}

func (b *MedicineBuilder) WithSubstances(substances []medicine.ActiveSubstance) *MedicineBuilder {
	b.value.Substances = append([]medicine.ActiveSubstance(nil), substances...)
	return b
}

func (b *MedicineBuilder) WithDosages(dosages []medicine.DosageRule) *MedicineBuilder {
	b.value.Dosages = append([]medicine.DosageRule(nil), dosages...)
	return b
}

func (b *MedicineBuilder) WithRecommendations(ids []uuid.UUID) *MedicineBuilder {
	b.value.Recommendation = append([]uuid.UUID(nil), ids...)
	return b
}

func (b *MedicineBuilder) Build() *medicine.Medicine {
	med := b.value
	med.Substances = append([]medicine.ActiveSubstance(nil), b.value.Substances...)
	med.Dosages = append([]medicine.DosageRule(nil), b.value.Dosages...)
	med.Recommendation = append([]uuid.UUID(nil), b.value.Recommendation...)
	med.Contraindications = append([]uuid.UUID(nil), b.value.Contraindications...)
	return &med
}

type CabinetItemBuilder struct {
	value cabinet_item.CabinetItem
}

func NewCabinetItemBuilder() *CabinetItemBuilder {
	return &CabinetItemBuilder{value: cabinet_item.CabinetItem{
		ID: CabinetItemID, UserID: CabinetUserID, MedicineID: CabinetMedicineID,
		DateOfManufacture: time.Date(2025, time.February, 10, 0, 0, 0, 0, time.UTC),
		Quantity:          1,
	}}
}

func (b *CabinetItemBuilder) WithID(id uuid.UUID) *CabinetItemBuilder {
	b.value.ID = id
	return b
}

func (b *CabinetItemBuilder) WithUserID(id uuid.UUID) *CabinetItemBuilder {
	b.value.UserID = id
	return b
}

func (b *CabinetItemBuilder) WithMedicineID(id uuid.UUID) *CabinetItemBuilder {
	b.value.MedicineID = id
	return b
}

func (b *CabinetItemBuilder) WithDateOfManufacture(date time.Time) *CabinetItemBuilder {
	b.value.DateOfManufacture = date
	return b
}

func (b *CabinetItemBuilder) WithQuantity(quantity float32) *CabinetItemBuilder {
	b.value.Quantity = quantity
	return b
}

func (b *CabinetItemBuilder) Build() *cabinet_item.CabinetItem {
	item := b.value
	return &item
}
