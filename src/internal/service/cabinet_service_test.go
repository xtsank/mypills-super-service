package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	svcErrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
)

type mockCabinetRepo struct {
	existing  *cabinet_item.CabinetItem
	saved     *cabinet_item.CabinetItem
	updated   *cabinet_item.CabinetItem
	deletedID uuid.UUID
}

func (m *mockCabinetRepo) FindExistingCabinetItem(ctx context.Context, userID, medID uuid.UUID, dom time.Time) (*cabinet_item.CabinetItem, error) {
	return m.existing, nil
}
func (m *mockCabinetRepo) Update(ctx context.Context, item *cabinet_item.CabinetItem) error {
	m.updated = item
	return nil
}
func (m *mockCabinetRepo) Save(ctx context.Context, item *cabinet_item.CabinetItem) error {
	m.saved = item
	return nil
}
func (m *mockCabinetRepo) Delete(ctx context.Context, id uuid.UUID) error {
	m.deletedID = id
	return nil
}
func (m *mockCabinetRepo) FindById(ctx context.Context, id uuid.UUID) (*cabinet_item.CabinetItem, error) {
	return m.existing, nil
}
func (m *mockCabinetRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*cabinet_item.CabinetItem, error) {
	return nil, nil
}

func TestCabinetService_AddItem_Existing(t *testing.T) {
	ex := &cabinet_item.CabinetItem{ID: uuid.New(), Quantity: 2}
	repo := &mockCabinetRepo{existing: ex}
	s := &CabinetService{cabinetRepo: repo}

	cmd := command.AddItemCmd{UserID: uuid.New(), MedicineID: uuid.New(), DateOfManufacture: ex.DateOfManufacture, Quantity: 3}
	dto, err := s.AddItem(context.Background(), &cmd)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if dto == nil {
		t.Fatalf("expected dto")
	}
	if ex.Quantity != 5 {
		t.Fatalf("quantity not increased")
	}
}

func TestCabinetService_AddItem_New(t *testing.T) {
	repo := &mockCabinetRepo{existing: nil}
	s := &CabinetService{cabinetRepo: repo}

	cmd := command.AddItemCmd{UserID: uuid.New(), MedicineID: uuid.New(), DateOfManufacture: time.Time{}, Quantity: 4}
	dto, err := s.AddItem(context.Background(), &cmd)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if dto == nil {
		t.Fatalf("expected dto")
	}
	if repo.saved == nil {
		t.Fatalf("expected saved item")
	}
}

func TestCabinetService_UpdateQty_NotFound(t *testing.T) {
	repo := &mockCabinetRepo{existing: nil}
	s := &CabinetService{cabinetRepo: repo}

	_, err := s.UpdateQty(context.Background(), &command.UpdateQtyCmd{ID: uuid.New(), Quantity: 1})
	if !errors.Is(err, svcErrors.ErrCabinetItemNotFound) {
		t.Fatalf("expected ErrCabinetItemNotFound, got %v", err)
	}
}

func TestCabinetService_UpdateQty_Success(t *testing.T) {
	ex := &cabinet_item.CabinetItem{ID: uuid.New(), Quantity: 2}
	repo := &mockCabinetRepo{existing: ex}
	s := &CabinetService{cabinetRepo: repo}

	dto, err := s.UpdateQty(context.Background(), &command.UpdateQtyCmd{ID: ex.ID, Quantity: 10})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if dto == nil {
		t.Fatalf("expected dto")
	}
	if ex.Quantity != 10 {
		t.Fatalf("quantity not updated")
	}
}

func TestCabinetService_RemoveItem(t *testing.T) {
	id := uuid.New()
	repo := &mockCabinetRepo{}
	s := &CabinetService{cabinetRepo: repo}

	err := s.RemoveItem(context.Background(), &command.RemoveItemCmd{ID: id})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if repo.deletedID != id {
		t.Fatalf("expected deleted id %v, got %v", id, repo.deletedID)
	}
}
