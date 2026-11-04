package service

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet"
	"github.com/xtsank/mypills-super-service/src/internal/domain/dictionary"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

type IRecommendationService interface {
	List(ctx context.Context, cmd *command.ListRecommendationsCmd) (*res.RecommendationResultDto, error)
}

type RecommendationService struct {
	users      user.IUserRepository
	medicines  medicine.IMedicineRepository
	cabinets   cabinet.ICabinetRepository
	items      cabinet.ICabinetItemRepository
	dictionary dictionary.IRepository
}

func NewRecommendationService(i do.Injector) (IRecommendationService, error) {
	return &RecommendationService{
		users:      do.MustInvoke[user.IUserRepository](i),
		medicines:  do.MustInvoke[medicine.IMedicineRepository](i),
		cabinets:   do.MustInvoke[cabinet.ICabinetRepository](i),
		items:      do.MustInvoke[cabinet.ICabinetItemRepository](i),
		dictionary: do.MustInvoke[dictionary.IRepository](i),
	}, nil
}

func (s *RecommendationService) List(ctx context.Context, cmd *command.ListRecommendationsCmd) (*res.RecommendationResultDto, error) {
	result := make([]medicine.Recommendation, 0)
	profile, err := s.users.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, apperrors.ErrUserNotFound.WithSource()
	}
	if err := s.cabinets.EnsureLegacy(ctx, cmd.UserID); err != nil {
		return nil, err
	}
	if cmd.CabinetID != nil {
		owned, err := s.cabinets.Find(ctx, cmd.UserID, *cmd.CabinetID)
		if err != nil {
			return nil, err
		}
		if owned == nil {
			return nil, apperrors.ErrCabinetNotFound.WithSource()
		}
	}
	medicines, err := s.medicines.FindByIllness(ctx, cmd.IllnessID)
	if err != nil {
		return nil, err
	}
	safe := make(map[uuid.UUID]*medicine.Medicine, len(medicines))
	for _, item := range medicines {
		if item.IsSafeFor(profile) {
			safe[item.ID] = item
		}
	}
	units, err := s.dictionary.ListUnits(ctx)
	if err != nil {
		return nil, err
	}
	unitNames := make(map[uuid.UUID]string, len(units))
	for _, item := range units {
		unitNames[item.ID] = item.Name
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayKey := today.Format("2006-01-02")
	for offset := 0; ; offset += 100 {
		page, err := s.cabinets.List(ctx, cmd.UserID, cmd.CabinetID, 100, offset)
		if err != nil {
			return nil, err
		}
		for _, group := range page.Items {
			stock, err := s.items.List(ctx, cmd.UserID, group.ID)
			if err != nil {
				return nil, err
			}
			for _, item := range stock {
				med := safe[item.MedicineID]
				if med == nil || item.Quantity <= 0 || item.DateOfManufacture.AddDate(0, item.ExpireMonths, 0).Format("2006-01-02") < todayKey {
					continue
				}
				dosage, frequency := med.CalculateDosage(profile)
				if dosage <= 0 || frequency <= 0 || item.Quantity < float64(dosage) {
					continue
				}
				result = append(result, medicine.Recommendation{
					MedicineID: med.ID, CabinetID: group.ID, CabinetItemID: item.ID,
					Name: med.Name, MethodOfApplication: med.MethodOfApplication,
					Dosage: dosage, UnitName: unitNames[med.Unit], FrequencyPerDay: frequency,
					QuantityInCabinet: item.Quantity,
				})
			}
		}
		if offset+len(page.Items) >= page.Total {
			break
		}
	}
	sort.Slice(result, func(i, j int) bool {
		a, b := result[i], result[j]
		if a.Name != b.Name {
			return a.Name < b.Name
		}
		if a.MedicineID != b.MedicineID {
			return a.MedicineID.String() < b.MedicineID.String()
		}
		return a.CabinetItemID.String() < b.CabinetItemID.String()
	})
	return res.NewRecommendationResultDto(result), nil
}
