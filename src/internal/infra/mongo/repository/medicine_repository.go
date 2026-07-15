package repository

import (
	"context"
	"errors"
	"sort"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/config"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	appErrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	mongoEntity "github.com/xtsank/mypills-super-service/src/internal/infra/mongo/entity"
	postgresEntity "github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoMedicineRepository struct {
	collection *mongo.Collection
}

func NewMongoMedicineRepository(i do.Injector) (medicine.IMedicineRepository, error) {
	client := do.MustInvoke[*mongo.Client](i)
	cfg := do.MustInvoke[*config.Config](i)

	collection := client.Database(cfg.DBName).Collection("medicines")

	return &MongoMedicineRepository{collection: collection}, nil
}

func (r *MongoMedicineRepository) FindByIllness(ctx context.Context, illnessID uuid.UUID) ([]*medicine.Medicine, error) {
	filter := bson.M{"recommendations": illnessID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}
	defer cursor.Close(ctx)

	var ents []mongoEntity.MedicineEntity
	if err := cursor.All(ctx, &ents); err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	meds := make([]*medicine.Medicine, 0, len(ents))
	for _, ent := range ents {
		med, err := mapMedicineEntity(ent)
		if err != nil {
			return nil, err
		}
		meds = append(meds, med)
	}

	return meds, nil
}

func (r *MongoMedicineRepository) FindByID(ctx context.Context, id uuid.UUID) (*medicine.Medicine, error) {
	var ent mongoEntity.MedicineEntity
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&ent)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, appErrors.ErrMedicineNotFound.WithSource()
		}
		return nil, appErrors.ErrInternal.WithError(err)
	}

	return mapMedicineEntity(ent)
}

func (r *MongoMedicineRepository) ListDosageRules(ctx context.Context) ([]postgresEntity.DosageRuleWithMedicineEntity, error) {
	opts := options.Find().SetProjection(bson.M{"name": 1, "dosages": 1})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}
	defer cursor.Close(ctx)

	var ents []mongoEntity.MedicineEntity
	if err := cursor.All(ctx, &ents); err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	items := make([]postgresEntity.DosageRuleWithMedicineEntity, 0)
	for _, ent := range ents {
		for _, d := range ent.Dosages {
			items = append(items, postgresEntity.DosageRuleWithMedicineEntity{
				ID:                  d.ID,
				MedicineID:          ent.ID,
				MedicineName:        ent.Name,
				ValueFrom:           d.ValueFrom,
				ValueTo:             d.ValueTo,
				DosageType:          d.DosageType,
				DosageValue:         d.DosageValue,
				NumberOfDosesPerDay: d.NumberOfDosesPerDay,
			})
		}
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].MedicineName != items[j].MedicineName {
			return items[i].MedicineName < items[j].MedicineName
		}
		if items[i].ValueFrom != items[j].ValueFrom {
			return items[i].ValueFrom < items[j].ValueFrom
		}
		return items[i].ValueTo < items[j].ValueTo
	})

	return items, nil
}

func (r *MongoMedicineRepository) Create(ctx context.Context, med *medicine.Medicine) error {
	ent := mapMedicineToEntity(med)
	if _, err := r.collection.InsertOne(ctx, ent); err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *MongoMedicineRepository) Update(ctx context.Context, med *medicine.Medicine) error {
	ent := mapMedicineToEntity(med)
	filter := bson.M{"_id": med.ID}

	result, err := r.collection.ReplaceOne(ctx, filter, ent)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	if result.MatchedCount == 0 {
		return appErrors.ErrMedicineNotFound.WithSource()
	}
	return nil
}

func (r *MongoMedicineRepository) Delete(ctx context.Context, id uuid.UUID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *MongoMedicineRepository) UpdateIndications(ctx context.Context, medicineID uuid.UUID, ids []uuid.UUID) error {
	return r.updateArrayField(ctx, medicineID, "recommendations", ids)
}

func (r *MongoMedicineRepository) UpdateContraindications(ctx context.Context, medicineID uuid.UUID, ids []uuid.UUID) error {
	return r.updateArrayField(ctx, medicineID, "contraindications", ids)
}

func (r *MongoMedicineRepository) UpdateComposition(ctx context.Context, medicineID uuid.UUID, substances []medicine.ActiveSubstance) error {
	subs := make([]mongoEntity.Substance, len(substances))
	for i, s := range substances {
		subs[i] = mongoEntity.Substance{SubstanceID: s.ID, Concentration: s.Concentration}
	}

	return r.updateArrayField(ctx, medicineID, "substances", subs)
}

func (r *MongoMedicineRepository) AddDosageRule(ctx context.Context, medicineID uuid.UUID, rule *medicine.DosageRule) error {
	filter := bson.M{"_id": medicineID}
	update := bson.M{
		"$push": bson.M{
			"dosages": mongoEntity.Dosage{
				ID:                  rule.ID,
				ValueFrom:           rule.ValueFrom,
				ValueTo:             rule.ValueTo,
				DosageType:          string(rule.Type),
				DosageValue:         rule.DosageValue,
				NumberOfDosesPerDay: rule.NumberOfDosesPerDay,
			},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	if result.MatchedCount == 0 {
		return appErrors.ErrMedicineNotFound.WithSource()
	}
	return nil
}

func (r *MongoMedicineRepository) DeleteDosageRule(ctx context.Context, ruleID uuid.UUID) error {
	filter := bson.M{"dosages.id": ruleID}
	update := bson.M{
		"$pull": bson.M{
			"dosages": bson.M{"id": ruleID},
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *MongoMedicineRepository) updateArrayField(ctx context.Context, medicineID uuid.UUID, field string, value interface{}) error {
	filter := bson.M{"_id": medicineID}
	update := bson.M{"$set": bson.M{field: value}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	if result.MatchedCount == 0 {
		return appErrors.ErrMedicineNotFound.WithSource()
	}
	return nil
}

func mapMedicineEntity(ent mongoEntity.MedicineEntity) (*medicine.Medicine, error) {
	substances := make([]medicine.ActiveSubstance, len(ent.Substances))
	for i, s := range ent.Substances {
		substances[i] = medicine.ActiveSubstance{ID: s.SubstanceID, Concentration: s.Concentration}
	}

	dosages := make([]medicine.DosageRule, len(ent.Dosages))
	for i, d := range ent.Dosages {
		dosages[i] = medicine.DosageRule{
			ID:                  d.ID,
			ValueFrom:           d.ValueFrom,
			ValueTo:             d.ValueTo,
			Type:                medicine.DosageType(d.DosageType),
			DosageValue:         d.DosageValue,
			NumberOfDosesPerDay: d.NumberOfDosesPerDay,
		}
	}

	return medicine.NewMedicine(
		ent.ID,
		ent.Name,
		ent.ExpireTime,
		ent.IsPrescription,
		ent.MethodOfApplication,
		ent.EffectOnPregnant,
		ent.EffectOnDriver,
		ent.FormID,
		ent.UnitID,
		substances,
		dosages,
		ent.Contraindications,
		ent.Recommendations,
	)
}

func mapMedicineToEntity(med *medicine.Medicine) mongoEntity.MedicineEntity {
	substances := make([]mongoEntity.Substance, len(med.Substances))
	for i, s := range med.Substances {
		substances[i] = mongoEntity.Substance{SubstanceID: s.ID, Concentration: s.Concentration}
	}

	dosages := make([]mongoEntity.Dosage, len(med.Dosages))
	for i, d := range med.Dosages {
		dosages[i] = mongoEntity.Dosage{
			ID:                  d.ID,
			ValueFrom:           d.ValueFrom,
			ValueTo:             d.ValueTo,
			DosageType:          string(d.Type),
			DosageValue:         d.DosageValue,
			NumberOfDosesPerDay: d.NumberOfDosesPerDay,
		}
	}

	return mongoEntity.MedicineEntity{
		ID:                  med.ID,
		FormID:              med.Form,
		UnitID:              med.Unit,
		Name:                med.Name,
		ExpireTime:          med.ExpireTime,
		EffectOnDriver:      med.EffectOnDriver,
		EffectOnPregnant:    med.EffectOnPregnant,
		MethodOfApplication: med.MethodOfApplication,
		IsPrescription:      med.IsPrescription,
		Substances:          substances,
		Dosages:             dosages,
		Contraindications:   med.Contraindications,
		Recommendations:     med.Recommendation,
	}
}
