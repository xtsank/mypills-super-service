package service

import (
	"context"
	"time"

	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

type IMedicineService interface {
	Select(ctx context.Context, cmd *command.SelectMedicineCmd) (*res.MedicineResDto, error)
}

type MedicineService struct {
	userRepo     user.IUserRepository
	medicineRepo medicine.IMedicineRepository
	cabinetRepo  cabinet_item.ICabinetItemRepository
}

func NewMedicineService(i do.Injector) (IMedicineService, error) {
	userRepo := do.MustInvoke[user.IUserRepository](i)
	medicineRepo := do.MustInvoke[medicine.IMedicineRepository](i)
	cabinetRepo := do.MustInvoke[cabinet_item.ICabinetItemRepository](i)

	return &MedicineService{
		userRepo:     userRepo,
		medicineRepo: medicineRepo,
		cabinetRepo:  cabinetRepo,
	}, nil
}

func (s *MedicineService) filterByUserInfo(allMeds []*medicine.Medicine, u *user.User) []*medicine.Medicine {
	var safeMeds []*medicine.Medicine

	for _, med := range allMeds {
		if med.IsSafeFor(u) {
			safeMeds = append(safeMeds, med)
		}
	}

	return safeMeds
}

func (s *MedicineService) filterByUserCabinet(allMeds []*medicine.Medicine, items []*cabinet_item.CabinetItem) []*medicine.Medicine {
	var available []*medicine.Medicine
	now := time.Now()

	for _, med := range allMeds {
		for _, item := range items {
			if med.ID == item.MedicineID {
				expireDate := item.DateOfManufacture.AddDate(0, med.ExpireTime, 0)

				if expireDate.After(now) && item.Quantity > 0 {
					available = append(available, med)
					break
				}
			}
		}
	}
	return available
}

func (s *MedicineService) buildMedicineRecommendations(
	meds []*medicine.Medicine,
	cabinet []*cabinet_item.CabinetItem,
	u *user.User,
) *res.MedicineResDto {

	result := make([]*res.MedicineRecommendation, 0, len(meds))

	for _, med := range meds {
		dosageValue, frequency := med.CalculateDosage(u)

		var quantity float32
		for _, item := range cabinet {
			if item.MedicineID == med.ID {
				quantity = item.Quantity
				break
			}
		}

		recommendation := &res.MedicineRecommendation{
			ID:                  med.ID,
			Name:                med.Name,
			MethodOfApplication: med.MethodOfApplication,
			Dosage:              dosageValue,
			Frequency:           frequency,
			QuantityInCabinet:   quantity,
		}

		result = append(result, recommendation)
	}

	return res.NewMedicineResDto(result)
}

func (s *MedicineService) Select(ctx context.Context, cmd *command.SelectMedicineCmd) (*res.MedicineResDto, error) {
	u, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.ErrUserNotFound
	}

	allMeds, err := s.medicineRepo.FindByIllness(ctx, cmd.IllnessID)
	if err != nil {
		return nil, err
	}

	safeMeds := s.filterByUserInfo(allMeds, u)

	cabinetItems, err := s.cabinetRepo.FindByUserID(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	availableMeds := s.filterByUserCabinet(safeMeds, cabinetItems)

	return s.buildMedicineRecommendations(availableMeds, cabinetItems, u), nil
}
