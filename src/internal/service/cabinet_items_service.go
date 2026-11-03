package service

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

type ICabinetItemsService interface {
	List(ctx context.Context, cmd *command.ListCabinetItemsCmd) (*res.CabinetItemsPageDto, error)
	Create(ctx context.Context, cmd *command.CreateCabinetItemCmd) (*res.CabinetItemDto, error)
	Update(ctx context.Context, cmd *command.UpdateCabinetItemCmd) (*res.CabinetItemDto, error)
	Delete(ctx context.Context, cmd *command.DeleteCabinetItemCmd) error
}

type CabinetItemsService struct {
	cabinets  cabinet.ICabinetRepository
	items     cabinet.ICabinetItemRepository
	medicines medicine.IMedicineRepository
}

func NewCabinetItemsService(i do.Injector) (ICabinetItemsService, error) {
	return &CabinetItemsService{
		cabinets:  do.MustInvoke[cabinet.ICabinetRepository](i),
		items:     do.MustInvoke[cabinet.ICabinetItemRepository](i),
		medicines: do.MustInvoke[medicine.IMedicineRepository](i),
	}, nil
}

func (s *CabinetItemsService) ensureOwner(ctx context.Context, userID, cabinetID uuid.UUID) error {
	if err := s.cabinets.EnsureLegacy(ctx, userID); err != nil {
		return err
	}
	item, err := s.cabinets.Find(ctx, userID, cabinetID)
	if err != nil {
		return err
	}
	if item == nil {
		return apperrors.ErrCabinetNotFound.WithSource()
	}
	return nil
}

func (s *CabinetItemsService) List(ctx context.Context, cmd *command.ListCabinetItemsCmd) (*res.CabinetItemsPageDto, error) {
	if err := cabinet.ValidateExpiryFilter(cmd.ExpiresWithinDays, cmd.Expired); err != nil {
		return nil, err
	}
	if err := s.ensureOwner(ctx, cmd.UserID, cmd.CabinetID); err != nil {
		return nil, err
	}
	items, err := s.items.List(ctx, cmd.UserID, cmd.CabinetID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayKey := today.Format("2006-01-02")
	filtered := make([]cabinet.Item, 0, len(items))
	for _, item := range items {
		if cmd.ID != nil && item.ID != *cmd.ID {
			continue
		}
		if cmd.LowStock && item.Quantity > item.MinimumQuantity {
			continue
		}
		expires := item.DateOfManufacture.AddDate(0, item.ExpireMonths, 0)
		expiryKey := expires.Format("2006-01-02")
		if cmd.Expired && expiryKey >= todayKey {
			continue
		}
		if cmd.ExpiresWithinDays != nil && (expiryKey < todayKey || expiryKey > today.AddDate(0, 0, *cmd.ExpiresWithinDays).Format("2006-01-02")) {
			continue
		}
		filtered = append(filtered, item)
	}
	sort.Slice(filtered, func(i, j int) bool {
		a := filtered[i].DateOfManufacture.AddDate(0, filtered[i].ExpireMonths, 0)
		b := filtered[j].DateOfManufacture.AddDate(0, filtered[j].ExpireMonths, 0)
		if a.Equal(b) {
			return filtered[i].ID.String() < filtered[j].ID.String()
		}
		return a.Before(b)
	})
	page := cabinet.ItemPage{Total: len(filtered), Limit: cmd.Limit, Offset: cmd.Offset}
	start := min(cmd.Offset, len(filtered))
	end := min(start+cmd.Limit, len(filtered))
	page.Items = filtered[start:end]
	return res.NewCabinetItemsPageDto(page), nil
}

func (s *CabinetItemsService) Create(ctx context.Context, cmd *command.CreateCabinetItemCmd) (*res.CabinetItemDto, error) {
	if err := s.ensureOwner(ctx, cmd.UserID, cmd.CabinetID); err != nil {
		return nil, err
	}
	med, err := s.medicines.FindByID(ctx, cmd.MedicineID)
	if err != nil {
		return nil, err
	}
	if med == nil {
		return nil, apperrors.ErrMedicineNotFound.WithSource()
	}
	item, err := cabinet.NewItem(uuid.New(), cmd.UserID, cmd.CabinetID, cmd.MedicineID,
		med.Name, med.ExpireTime, cmd.DateOfManufacture, cmd.Quantity, cmd.MinimumQuantity)
	if err != nil {
		return nil, err
	}
	if err := s.items.Create(ctx, *item); err != nil {
		return nil, err
	}
	return res.NewCabinetItemDto(*item), nil
}

func (s *CabinetItemsService) Update(ctx context.Context, cmd *command.UpdateCabinetItemCmd) (*res.CabinetItemDto, error) {
	if err := s.ensureOwner(ctx, cmd.UserID, cmd.CabinetID); err != nil {
		return nil, err
	}
	item, err := s.items.Find(ctx, cmd.UserID, cmd.CabinetID, cmd.ItemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, apperrors.ErrCabinetItemNotFound.WithSource()
	}
	updated, err := cabinet.NewUpdatedItem(item, cmd.Quantity, cmd.MinimumQuantity, cmd.DateOfManufacture)
	if err != nil {
		return nil, err
	}
	if err := s.items.Update(ctx, *updated); err != nil {
		return nil, err
	}
	return res.NewCabinetItemDto(*updated), nil
}

func (s *CabinetItemsService) Delete(ctx context.Context, cmd *command.DeleteCabinetItemCmd) error {
	if err := s.ensureOwner(ctx, cmd.UserID, cmd.CabinetID); err != nil {
		return err
	}
	item, err := s.items.Find(ctx, cmd.UserID, cmd.CabinetID, cmd.ItemID)
	if err != nil {
		return err
	}
	if item == nil {
		return apperrors.ErrCabinetItemNotFound.WithSource()
	}
	return s.items.Delete(ctx, cmd.UserID, cmd.CabinetID, cmd.ItemID)
}
