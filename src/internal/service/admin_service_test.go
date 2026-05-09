package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	svcErrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	req "github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
)

type mockMedicineRepo struct {
	created   *medicine.Medicine
	findRes   *medicine.Medicine
	createErr error
	updateErr error

	deletedID uuid.UUID

	updIndMedID uuid.UUID
	updIndIDs   []uuid.UUID

	updContraMedID uuid.UUID
	updContraIDs   []uuid.UUID

	compMedID uuid.UUID
	compSubs  []medicine.ActiveSubstance

	addedRuleMedID uuid.UUID
	addedRule      *medicine.DosageRule

	deletedRuleID uuid.UUID
}

func (m *mockMedicineRepo) Create(ctx context.Context, med *medicine.Medicine) error {
	m.created = med
	return m.createErr
}
func (m *mockMedicineRepo) FindByID(ctx context.Context, id uuid.UUID) (*medicine.Medicine, error) {
	return m.findRes, nil
}
func (m *mockMedicineRepo) Update(ctx context.Context, med *medicine.Medicine) error {
	return m.updateErr
}
func (m *mockMedicineRepo) Delete(ctx context.Context, id uuid.UUID) error {
	m.deletedID = id
	return nil
}
func (m *mockMedicineRepo) UpdateIndications(ctx context.Context, medID uuid.UUID, ids []uuid.UUID) error {
	m.updIndMedID = medID
	m.updIndIDs = ids
	return nil
}
func (m *mockMedicineRepo) UpdateContraindications(ctx context.Context, medID uuid.UUID, ids []uuid.UUID) error {
	m.updContraMedID = medID
	m.updContraIDs = ids
	return nil
}
func (m *mockMedicineRepo) UpdateComposition(ctx context.Context, medID uuid.UUID, subs []medicine.ActiveSubstance) error {
	m.compMedID = medID
	m.compSubs = subs
	return nil
}
func (m *mockMedicineRepo) AddDosageRule(ctx context.Context, medID uuid.UUID, rule *medicine.DosageRule) error {
	m.addedRuleMedID = medID
	m.addedRule = rule
	return nil
}
func (m *mockMedicineRepo) DeleteDosageRule(ctx context.Context, ruleID uuid.UUID) error {
	m.deletedRuleID = ruleID
	return nil
}
func (m *mockMedicineRepo) FindByIllness(ctx context.Context, illnessID uuid.UUID) ([]*medicine.Medicine, error) {
	return nil, nil
}

func TestAdminService_AddMedicine(t *testing.T) {
	repo := &mockMedicineRepo{}
	s := &AdminService{medicineRepo: repo}

	cmd := command.AddMedicineCmd{Name: "m", ExpireTime: 1, Substances: nil, Dosages: nil}
	dto, err := s.AddMedicine(context.Background(), &cmd)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if dto == nil {
		t.Fatalf("expected dto")
	}
	if repo.created == nil || repo.created.Name != "m" {
		t.Fatalf("created not set")
	}
}

func TestAdminService_UpdateMedicine_NotFound(t *testing.T) {
	repo := &mockMedicineRepo{findRes: nil}
	s := &AdminService{medicineRepo: repo}

	_, err := s.UpdateMedicine(context.Background(), &command.UpdateMedicineCmd{ID: uuid.New()})
	if !errors.Is(err, svcErrors.ErrMedicineNotFound) {
		t.Fatalf("expected ErrMedicineNotFound, got %v", err)
	}
}

func TestAdminService_UpdateMedicine_Success(t *testing.T) {
	med := &medicine.Medicine{ID: uuid.New(), ExpireTime: 1}
	repo := &mockMedicineRepo{findRes: med}
	s := &AdminService{medicineRepo: repo}

	exptime := 5
	dto, err := s.UpdateMedicine(context.Background(), &command.UpdateMedicineCmd{ID: med.ID, ExpireTime: &exptime})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if dto == nil {
		t.Fatalf("expected dto")
	}
	if med.ExpireTime != 5 {
		t.Fatalf("expire not updated")
	}
}

func TestAdminService_RemoveMedicine(t *testing.T) {
	id := uuid.New()
	repo := &mockMedicineRepo{}
	s := &AdminService{medicineRepo: repo}

	err := s.RemoveMedicine(context.Background(), &command.RemoveMedicineCmd{ID: id})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if repo.deletedID != id {
		t.Fatalf("expected deleted id to be %v, got %v", id, repo.deletedID)
	}
}

func TestAdminService_UpdateIndicationsAndContra(t *testing.T) {
	medID := uuid.New()
	ids := []uuid.UUID{uuid.New(), uuid.New()}
	repo := &mockMedicineRepo{}
	s := &AdminService{medicineRepo: repo}

	err := s.UpdateIndications(context.Background(), command.UpdateLinksCmd{MedicineID: medID, IDs: ids})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if repo.updIndMedID != medID || !reflect.DeepEqual(repo.updIndIDs, ids) {
		t.Fatalf("indications not set")
	}

	err = s.UpdateContraindications(context.Background(), command.UpdateLinksCmd{MedicineID: medID, IDs: ids})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if repo.updContraMedID != medID || !reflect.DeepEqual(repo.updContraIDs, ids) {
		t.Fatalf("contraindications not set")
	}
}

func TestAdminService_UpdateComposition(t *testing.T) {
	medID := uuid.New()
	subs := []*req.ActiveSubstanceDto{{ID: uuid.New(), Concentration: 1.2}}
	repo := &mockMedicineRepo{}
	s := &AdminService{medicineRepo: repo}

	err := s.UpdateComposition(context.Background(), command.UpdateCompositionCmd{MedicineID: medID, Substances: subs})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if repo.compMedID != medID {
		t.Fatalf("comp med id not set")
	}
	if len(repo.compSubs) != 1 || repo.compSubs[0].Concentration != subs[0].Concentration {
		t.Fatalf("composition not transferred")
	}
}

func TestAdminService_AddAndRemoveDosageRule(t *testing.T) {
	medID := uuid.New()
	dto := &req.DosageRuleDto{ValueFrom: 1, ValueTo: 2, Type: req.ByWeight, DosageValue: 0.5, NumberOfDosesPerDay: 1}
	repo := &mockMedicineRepo{}
	s := &AdminService{medicineRepo: repo}

	err := s.AddDosageRule(context.Background(), command.AddDosageRuleCmd{MedicineID: medID, Dosage: dto})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if repo.addedRuleMedID != medID || repo.addedRule == nil {
		t.Fatalf("expected added rule")
	}

	ruleID := uuid.New()
	err = s.DeleteDosageRule(context.Background(), command.RemoveDosageRuleCmd{RuleID: ruleID})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if repo.deletedRuleID != ruleID {
		t.Fatalf("expected deleted rule id set")
	}
}
