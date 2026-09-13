package testutil

import (
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	"github.com/xtsank/mypills-super-service/src/internal/test_util/testid"
)

type UserBuilder struct {
	value user.User
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{value: user.User{
		ID:       testid.DefaultUser,
		Login:    "test-user",
		Email:    "test@example.com",
		Password: "password-hash",
		Weight:   70,
		Age:      30,
		Notify:   &user.NotifyInfo{},
	}}
}

func (b *UserBuilder) WithID(id uuid.UUID) *UserBuilder {
	b.value.ID = id
	return b
}

func (b *UserBuilder) WithWeight(weight int) *UserBuilder {
	b.value.Weight = weight
	return b
}

func (b *UserBuilder) WithPregnancy(value bool) *UserBuilder {
	b.value.IsPregnant = value
	return b
}

func (b *UserBuilder) WithDriver(value bool) *UserBuilder {
	b.value.IsDriver = value
	return b
}

func (b *UserBuilder) WithIllnesses(ids ...uuid.UUID) *UserBuilder {
	b.value.Illnesses = append([]uuid.UUID(nil), ids...)
	return b
}

func (b *UserBuilder) WithAllergies(ids ...uuid.UUID) *UserBuilder {
	b.value.Allergies = append([]uuid.UUID(nil), ids...)
	return b
}

func (b *UserBuilder) Build() *user.User {
	result := b.value
	result.Illnesses = append([]uuid.UUID(nil), b.value.Illnesses...)
	result.Allergies = append([]uuid.UUID(nil), b.value.Allergies...)
	if b.value.Notify != nil {
		notify := *b.value.Notify
		result.Notify = &notify
	}
	return &result
}

type MedicineBuilder struct {
	value medicine.Medicine
}

func NewMedicineBuilder() *MedicineBuilder {
	return &MedicineBuilder{value: medicine.Medicine{
		ID:                  testid.DefaultMedicine,
		Name:                "Test medicine",
		ExpireTime:          24,
		MethodOfApplication: "oral",
		Form:                testid.DefaultForm,
		Unit:                testid.DefaultUnit,
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

func (b *MedicineBuilder) WithExpireTime(months int) *MedicineBuilder {
	b.value.ExpireTime = months
	return b
}

func (b *MedicineBuilder) WithPregnancyRisk(value bool) *MedicineBuilder {
	b.value.EffectOnPregnant = value
	return b
}

func (b *MedicineBuilder) WithDriverRisk(value bool) *MedicineBuilder {
	b.value.EffectOnDriver = value
	return b
}

func (b *MedicineBuilder) WithSubstances(values ...medicine.ActiveSubstance) *MedicineBuilder {
	b.value.Substances = append([]medicine.ActiveSubstance(nil), values...)
	return b
}

func (b *MedicineBuilder) WithDosages(values ...medicine.DosageRule) *MedicineBuilder {
	b.value.Dosages = append([]medicine.DosageRule(nil), values...)
	return b
}

func (b *MedicineBuilder) WithContraindications(ids ...uuid.UUID) *MedicineBuilder {
	b.value.Contraindications = append([]uuid.UUID(nil), ids...)
	return b
}

func (b *MedicineBuilder) WithRecommendations(ids ...uuid.UUID) *MedicineBuilder {
	b.value.Recommendation = append([]uuid.UUID(nil), ids...)
	return b
}

func (b *MedicineBuilder) Build() *medicine.Medicine {
	result := b.value
	result.Substances = append([]medicine.ActiveSubstance(nil), b.value.Substances...)
	result.Dosages = append([]medicine.DosageRule(nil), b.value.Dosages...)
	result.Contraindications = append([]uuid.UUID(nil), b.value.Contraindications...)
	result.Recommendation = append([]uuid.UUID(nil), b.value.Recommendation...)
	return &result
}

type CabinetItemBuilder struct {
	value cabinet_item.CabinetItem
}

func NewCabinetItemBuilder() *CabinetItemBuilder {
	return &CabinetItemBuilder{value: cabinet_item.CabinetItem{
		ID:                testid.DefaultCabinetItem,
		UserID:            testid.DefaultUser,
		MedicineID:        testid.DefaultMedicine,
		DateOfManufacture: time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC),
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

func (b *CabinetItemBuilder) WithManufactureDate(date time.Time) *CabinetItemBuilder {
	b.value.DateOfManufacture = date
	return b
}

func (b *CabinetItemBuilder) WithQuantity(quantity float32) *CabinetItemBuilder {
	b.value.Quantity = quantity
	return b
}

func (b *CabinetItemBuilder) Build() *cabinet_item.CabinetItem {
	result := b.value
	return &result
}
