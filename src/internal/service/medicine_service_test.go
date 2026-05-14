package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	svcErrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
)

type mockUserRepoForMed struct {
	u *user.User
}

func (m *mockUserRepoForMed) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	return m.u, nil
}

// unused
func (m *mockUserRepoForMed) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	panic("unused")
}
func (m *mockUserRepoForMed) Create(ctx context.Context, u *user.User) error { panic("unused") }
func (m *mockUserRepoForMed) FindByLogin(ctx context.Context, login string) (*user.User, error) {
	panic("unused")
}
func (m *mockUserRepoForMed) Update(ctx context.Context, u *user.User) error { panic("unused") }

type mockMedRepo struct {
	meds []*medicine.Medicine
}

func (m *mockMedRepo) FindByIllness(ctx context.Context, illnessID uuid.UUID) ([]*medicine.Medicine, error) {
	return m.meds, nil
}

// unused stubs
func (m *mockMedRepo) Create(ctx context.Context, med *medicine.Medicine) error { panic("unused") }
func (m *mockMedRepo) FindByID(ctx context.Context, id uuid.UUID) (*medicine.Medicine, error) {
	panic("unused")
}
func (m *mockMedRepo) Update(ctx context.Context, med *medicine.Medicine) error { panic("unused") }
func (m *mockMedRepo) Delete(ctx context.Context, id uuid.UUID) error           { panic("unused") }
func (m *mockMedRepo) UpdateIndications(ctx context.Context, medID uuid.UUID, ids []uuid.UUID) error {
	panic("unused")
}
func (m *mockMedRepo) UpdateContraindications(ctx context.Context, medID uuid.UUID, ids []uuid.UUID) error {
	panic("unused")
}
func (m *mockMedRepo) UpdateComposition(ctx context.Context, medID uuid.UUID, subs []medicine.ActiveSubstance) error {
	panic("unused")
}
func (m *mockMedRepo) AddDosageRule(ctx context.Context, medID uuid.UUID, rule *medicine.DosageRule) error {
	panic("unused")
}
func (m *mockMedRepo) DeleteDosageRule(ctx context.Context, ruleID uuid.UUID) error { panic("unused") }

type mockCabinetRepoForMed struct {
	items []*cabinet_item.CabinetItem
}

func (m *mockCabinetRepoForMed) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*cabinet_item.CabinetItem, error) {
	return m.items, nil
}

// unused
func (m *mockCabinetRepoForMed) FindExistingCabinetItem(ctx context.Context, userID, medID uuid.UUID, dom time.Time) (*cabinet_item.CabinetItem, error) {
	panic("unused")
}
func (m *mockCabinetRepoForMed) Update(ctx context.Context, item *cabinet_item.CabinetItem) error {
	panic("unused")
}
func (m *mockCabinetRepoForMed) Save(ctx context.Context, item *cabinet_item.CabinetItem) error {
	panic("unused")
}
func (m *mockCabinetRepoForMed) Delete(ctx context.Context, id uuid.UUID) error { panic("unused") }
func (m *mockCabinetRepoForMed) FindById(ctx context.Context, id uuid.UUID) (*cabinet_item.CabinetItem, error) {
	panic("unused")
}

func TestMedicineService_Select_UserNotFound(t *testing.T) {
	s := &MedicineService{userRepo: &mockUserRepoForMed{u: nil}}
	_, err := s.Select(context.Background(), &command.SelectMedicineCmd{UserID: uuid.New()})
	if !errors.Is(err, svcErrors.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestMedicineService_Select_Success(t *testing.T) {
	uid := uuid.New()
	u := &user.User{ID: uid, Weight: 70, Age: 30}

	// medicine safe for user and has dosage by weight
	med := &medicine.Medicine{
		ID: uuid.New(), Name: "m",
		ExpireTime: 12,
		Substances: nil,
		Dosages:    []medicine.DosageRule{{ID: uuid.New(), ValueFrom: 50, ValueTo: 100, Type: medicine.ByWeight, DosageValue: 1.5, NumberOfDosesPerDay: 2}},
	}

	// cabinet item not expired and quantity >0
	dom := time.Now().AddDate(0, -1, 0)
	item := &cabinet_item.CabinetItem{ID: uuid.New(), MedicineID: med.ID, UserID: uid, DateOfManufacture: dom, Quantity: 3}

	s := &MedicineService{userRepo: &mockUserRepoForMed{u: u}, medicineRepo: &mockMedRepo{meds: []*medicine.Medicine{med}}, cabinetRepo: &mockCabinetRepoForMed{items: []*cabinet_item.CabinetItem{item}}}

	dto, err := s.Select(context.Background(), &command.SelectMedicineCmd{UserID: uid, IllnessID: uuid.Nil})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if dto == nil {
		t.Fatalf("expected dto")
	}
	if len(dto.Recommendations) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(dto.Recommendations))
	}
	if dto.Recommendations[0].Dosage == 0 {
		t.Fatalf("expected dosage computed")
	}
}
