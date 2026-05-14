package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	appErrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
)

type PostgresMedicineRepository struct {
	db *sqlx.DB
}

func NewPostgresMedicineRepository(i do.Injector) (medicine.IMedicineRepository, error) {
	db := do.MustInvoke[*sqlx.DB](i)

	return &PostgresMedicineRepository{db: db}, nil
}

func (r *PostgresMedicineRepository) findBaseByID(ctx context.Context, id uuid.UUID) (*entity.MedicineEntity, error) {
	var ent entity.MedicineEntity
	query := `SELECT * FROM Medicine WHERE id = $1`

	err := r.db.GetContext(ctx, &ent, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, appErrors.ErrInternal.WithError(err)
	}

	return &ent, nil
}

func (r *PostgresMedicineRepository) findBasesByIllness(ctx context.Context, illnessID uuid.UUID) ([]entity.MedicineEntity, error) {
	var ents []entity.MedicineEntity
	query := `
       SELECT m.* FROM Medicine m
       JOIN Recommendations r ON m.id = r.medicine_id
       WHERE r.illness_id = $1`

	err := r.db.SelectContext(ctx, &ents, query, illnessID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	return ents, nil
}

func (r *PostgresMedicineRepository) getSubstances(ctx context.Context, medID uuid.UUID) ([]medicine.ActiveSubstance, error) {
	var ents []entity.MedicineSubstanceEntity
	query := `SELECT * FROM Medicine_Substance WHERE medicine_id = $1`

	err := r.db.SelectContext(ctx, &ents, query, medID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	if len(ents) == 0 {
		return []medicine.ActiveSubstance{}, nil
	}

	res := make([]medicine.ActiveSubstance, len(ents))
	for i, e := range ents {
		res[i] = medicine.ActiveSubstance{ID: e.SubstanceId, Concentration: e.Concentration}
	}
	return res, nil
}

func (r *PostgresMedicineRepository) getDosages(ctx context.Context, medID uuid.UUID) ([]medicine.DosageRule, error) {
	var ents []entity.DosageEntity

	query := `SELECT * FROM Dosage WHERE medicine_id = $1`

	err := r.db.SelectContext(ctx, &ents, query, medID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	if len(ents) == 0 {
		return []medicine.DosageRule{}, nil
	}

	res := make([]medicine.DosageRule, len(ents))
	for i, e := range ents {
		res[i] = medicine.DosageRule{
			ID:                  e.ID,
			ValueFrom:           e.ValueFrom,
			ValueTo:             e.ValueTo,
			Type:                medicine.DosageType(e.DosageType),
			DosageValue:         e.DosageValue,
			NumberOfDosesPerDay: e.NumberOfDosesPerDay,
		}
	}
	return res, nil
}

func (r *PostgresMedicineRepository) getContraindications(ctx context.Context, medID uuid.UUID) ([]uuid.UUID, error) {
	var illnesses []uuid.UUID
	query := `SELECT illness_id FROM Contraindications WHERE medicine_id = $1`

	err := r.db.SelectContext(ctx, &illnesses, query, medID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	if illnesses == nil {
		return []uuid.UUID{}, nil
	}

	return illnesses, nil
}

func (r *PostgresMedicineRepository) getRecommendations(ctx context.Context, medID uuid.UUID) ([]uuid.UUID, error) {
	var illnesses []uuid.UUID
	query := `SELECT illness_id FROM Recommendations WHERE medicine_id = $1`

	err := r.db.SelectContext(ctx, &illnesses, query, medID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	if illnesses == nil {
		return []uuid.UUID{}, nil
	}

	return illnesses, nil
}

func (r *PostgresMedicineRepository) FindByIllness(ctx context.Context, illnessID uuid.UUID) ([]*medicine.Medicine, error) {
	ents, err := r.findBasesByIllness(ctx, illnessID)
	if err != nil {
		return nil, err
	}
	if len(ents) == 0 {
		return []*medicine.Medicine{}, nil
	}

	meds := make([]*medicine.Medicine, len(ents))

	for i, ent := range ents {
		substances, err := r.getSubstances(ctx, ent.ID)
		if err != nil {
			return nil, err
		}

		dosages, err := r.getDosages(ctx, ent.ID)
		if err != nil {
			return nil, err
		}

		contra, err := r.getContraindications(ctx, ent.ID)
		if err != nil {
			return nil, err
		}

		recom, err := r.getRecommendations(ctx, ent.ID)
		if err != nil {
			return nil, err
		}

		m, err := medicine.NewMedicine(
			ent.ID, ent.Name, ent.ExpireTime, ent.IsPrescription,
			ent.MethodOfApplication, ent.EffectOnPregnant, ent.EffectOnDriver,
			ent.FormId, ent.UnitId,
			substances, dosages, contra, recom,
		)
		if err != nil {
			return nil, err
		}
		meds[i] = m
	}

	return meds, nil
}

func (r *PostgresMedicineRepository) FindByID(ctx context.Context, id uuid.UUID) (*medicine.Medicine, error) {
	ent, err := r.findBaseByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ent == nil {
		return nil, appErrors.ErrMedicineNotFound.WithSource()
	}

	substances, err := r.getSubstances(ctx, ent.ID)
	if err != nil {
		return nil, err
	}

	dosages, err := r.getDosages(ctx, ent.ID)
	if err != nil {
		return nil, err
	}

	contra, err := r.getContraindications(ctx, ent.ID)
	if err != nil {
		return nil, err
	}

	recom, err := r.getRecommendations(ctx, ent.ID)
	if err != nil {
		return nil, err
	}

	return medicine.NewMedicine(
		ent.ID, ent.Name, ent.ExpireTime, ent.IsPrescription,
		ent.MethodOfApplication, ent.EffectOnPregnant, ent.EffectOnDriver,
		ent.FormId, ent.UnitId,
		substances, dosages, contra, recom,
	)
}

func (r *PostgresMedicineRepository) insertBase(ctx context.Context, tx *sqlx.Tx, med *medicine.Medicine) error {
	ent := entity.MedicineEntity{
		ID:                  med.ID,
		Name:                med.Name,
		ExpireTime:          med.ExpireTime,
		IsPrescription:      med.IsPrescription,
		MethodOfApplication: med.MethodOfApplication,
		EffectOnPregnant:    med.EffectOnPregnant,
		EffectOnDriver:      med.EffectOnDriver,
		FormId:              med.Form,
		UnitId:              med.Unit,
	}

	query := `INSERT INTO Medicine  (id, name, expire_time, is_prescription, 
                                   method_of_application, effect_on_pregnant, effect_on_driver, 
                                   form_id, unit_id)
              VALUES (:id, :name, :expire_time, :is_prescription, :method_of_application, :effect_on_pregnant,
                      :effect_on_driver, :form_id, :unit_id)`

	_, err := tx.NamedExecContext(ctx, query, ent)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) insertSubstances(ctx context.Context, tx *sqlx.Tx, medID uuid.UUID, substances []medicine.ActiveSubstance) error {
	if len(substances) == 0 {
		return nil
	}

	rows := make([]map[string]interface{}, len(substances))
	for i, s := range substances {
		rows[i] = map[string]interface{}{
			"medicine_id":   medID,
			"substance_id":  s.ID,
			"concentration": s.Concentration,
		}
	}

	query := `INSERT INTO Medicine_Substance (medicine_id, substance_id, concentration) 
              VALUES (:medicine_id, :substance_id, :concentration)`

	_, err := tx.NamedExecContext(ctx, query, rows)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) insertRecommendations(ctx context.Context, tx *sqlx.Tx, medID uuid.UUID, recommendations []uuid.UUID) error {
	if len(recommendations) == 0 {
		return nil
	}

	rows := make([]map[string]interface{}, len(recommendations))
	for i, illnessID := range recommendations {
		rows[i] = map[string]interface{}{
			"medicine_id": medID,
			"illness_id":  illnessID,
		}
	}

	query := `INSERT INTO Recommendations (medicine_id, illness_id) VALUES (:medicine_id, :illness_id)`
	_, err := tx.NamedExecContext(ctx, query, rows)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) insertContraindications(ctx context.Context, tx *sqlx.Tx, medID uuid.UUID, contraindications []uuid.UUID) error {
	if len(contraindications) == 0 {
		return nil
	}

	rows := make([]map[string]interface{}, len(contraindications))
	for i, illnessID := range contraindications {
		rows[i] = map[string]interface{}{
			"medicine_id": medID,
			"illness_id":  illnessID,
		}
	}

	query := `INSERT INTO Contraindications (medicine_id, illness_id) VALUES (:medicine_id, :illness_id)`
	_, err := tx.NamedExecContext(ctx, query, rows)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) insertDosages(ctx context.Context, tx *sqlx.Tx, medID uuid.UUID, dosages []medicine.DosageRule) error {
	if len(dosages) == 0 {
		return nil
	}

	ents := make([]entity.DosageEntity, len(dosages))
	for i, d := range dosages {
		ents[i] = entity.DosageEntity{
			ID:                  d.ID,
			MedicineId:          medID,
			ValueFrom:           d.ValueFrom,
			ValueTo:             d.ValueTo,
			DosageType:          string(d.Type),
			DosageValue:         d.DosageValue,
			NumberOfDosesPerDay: d.NumberOfDosesPerDay,
		}
	}

	query := `INSERT INTO Dosage (id, medicine_id, value_from, value_to, dosage_type, dosage_value, number_of_doses_per_day)
              VALUES (:id, :medicine_id, :value_from, :value_to, :dosage_type, :dosage_value, :number_of_doses_per_day)`

	_, err := tx.NamedExecContext(ctx, query, ents)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) Create(ctx context.Context, med *medicine.Medicine) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	defer tx.Rollback()

	if err := r.insertBase(ctx, tx, med); err != nil {
		return err
	}

	if err := r.insertSubstances(ctx, tx, med.ID, med.Substances); err != nil {
		return err
	}

	if err := r.insertRecommendations(ctx, tx, med.ID, med.Recommendation); err != nil {
		return err
	}

	if err := r.insertContraindications(ctx, tx, med.ID, med.Contraindications); err != nil {
		return err
	}

	if err := r.insertDosages(ctx, tx, med.ID, med.Dosages); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) deleteSubstances(ctx context.Context, tx *sqlx.Tx, medID uuid.UUID) error {
	query := `DELETE FROM Medicine_Substance WHERE medicine_id = $1`
	_, err := tx.ExecContext(ctx, query, medID)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) deleteDosages(ctx context.Context, tx *sqlx.Tx, medID uuid.UUID) error {
	query := `DELETE FROM Dosage WHERE medicine_id = $1`
	_, err := tx.ExecContext(ctx, query, medID)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) deleteRecommendations(ctx context.Context, tx *sqlx.Tx, medID uuid.UUID) error {
	query := `DELETE FROM Recommendations WHERE medicine_id = $1`
	_, err := tx.ExecContext(ctx, query, medID)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) deleteContraindications(ctx context.Context, tx *sqlx.Tx, medID uuid.UUID) error {
	query := `DELETE FROM Contraindications WHERE medicine_id = $1`
	_, err := tx.ExecContext(ctx, query, medID)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) updateBase(ctx context.Context, tx *sqlx.Tx, med *medicine.Medicine) error {
	ent := entity.MedicineEntity{
		ID:                  med.ID,
		Name:                med.Name,
		ExpireTime:          med.ExpireTime,
		IsPrescription:      med.IsPrescription,
		MethodOfApplication: med.MethodOfApplication,
		EffectOnPregnant:    med.EffectOnPregnant,
		EffectOnDriver:      med.EffectOnDriver,
		FormId:              med.Form,
		UnitId:              med.Unit,
	}

	query := `UPDATE Medicine 
              SET name = :name, expire_time = :expire_time, is_prescription = :is_prescription, 
                  method_of_application = :method_of_application, effect_on_pregnant = :effect_on_pregnant, 
                  effect_on_driver = :effect_on_driver, form_id = :form_id, unit_id = :unit_id
              WHERE id = :id`

	_, err := tx.NamedExecContext(ctx, query, ent)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) Update(ctx context.Context, med *medicine.Medicine) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	defer tx.Rollback()

	if err := r.updateBase(ctx, tx, med); err != nil {
		return err
	}

	if err := r.deleteSubstances(ctx, tx, med.ID); err != nil {
		return err
	}
	if err := r.insertSubstances(ctx, tx, med.ID, med.Substances); err != nil {
		return err
	}

	if err := r.deleteDosages(ctx, tx, med.ID); err != nil {
		return err
	}
	if err := r.insertDosages(ctx, tx, med.ID, med.Dosages); err != nil {
		return err
	}

	if err := r.deleteRecommendations(ctx, tx, med.ID); err != nil {
		return err
	}
	if err := r.insertRecommendations(ctx, tx, med.ID, med.Recommendation); err != nil {
		return err
	}

	if err := r.deleteContraindications(ctx, tx, med.ID); err != nil {
		return err
	}
	if err := r.insertContraindications(ctx, tx, med.ID, med.Contraindications); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) deleteBase(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) error {
	query := `delete from Medicine where id = $1`
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	defer tx.Rollback()

	if err := r.deleteSubstances(ctx, tx, id); err != nil {
		return err
	}

	if err := r.deleteDosages(ctx, tx, id); err != nil {
		return err
	}

	if err := r.deleteRecommendations(ctx, tx, id); err != nil {
		return err
	}

	if err := r.deleteContraindications(ctx, tx, id); err != nil {
		return err
	}

	if err := r.deleteBase(ctx, tx, id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) UpdateIndications(ctx context.Context, medicineID uuid.UUID, ids []uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	defer tx.Rollback()

	if err := r.deleteRecommendations(ctx, tx, medicineID); err != nil {
		return err
	}
	if err := r.insertRecommendations(ctx, tx, medicineID, ids); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) UpdateContraindications(ctx context.Context, medicineID uuid.UUID, ids []uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	defer tx.Rollback()

	if err := r.deleteContraindications(ctx, tx, medicineID); err != nil {
		return err
	}
	if err := r.insertContraindications(ctx, tx, medicineID, ids); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) UpdateComposition(ctx context.Context, medicineID uuid.UUID, substances []medicine.ActiveSubstance) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	defer tx.Rollback()

	if err := r.deleteSubstances(ctx, tx, medicineID); err != nil {
		return err
	}
	if err := r.insertSubstances(ctx, tx, medicineID, substances); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) AddDosageRule(ctx context.Context, medicineID uuid.UUID, rule *medicine.DosageRule) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	defer tx.Rollback()

	if err := r.insertDosages(ctx, tx, medicineID, []medicine.DosageRule{*rule}); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresMedicineRepository) DeleteDosageRule(ctx context.Context, ruleID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	defer tx.Rollback()

	query := `delete from Dosage where id = $1`
	if _, err := tx.ExecContext(ctx, query, ruleID); err != nil {
		return appErrors.ErrInternal.WithError(err)
	}

	if err := tx.Commit(); err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}
