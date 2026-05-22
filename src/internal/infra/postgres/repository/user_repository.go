package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	appErrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/queries"
)

type PostgresUserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(i do.Injector) (user.IUserRepository, error) {
	db := do.MustInvoke[*sqlx.DB](i)

	return &PostgresUserRepository{db: db}, nil
}

func (r *PostgresUserRepository) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	var exists bool
	query := queries.User.ExistsByLogin

	err := r.db.GetContext(ctx, &exists, query, login)
	if err != nil {
		return false, appErrors.ErrInternal.WithError(err)
	}

	return exists, nil
}

func (r *PostgresUserRepository) findBaseByLogin(ctx context.Context, login string) (*entity.UserEntity, error) {
	var ent entity.UserEntity
	query := queries.User.FindByLogin

	err := r.db.GetContext(ctx, &ent, query, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, appErrors.ErrInternal.WithError(err)
	}

	return &ent, nil
}

func (r *PostgresUserRepository) findBaseByID(ctx context.Context, id uuid.UUID) (*entity.UserEntity, error) {
	var ent entity.UserEntity
	query := queries.User.FindByID

	err := r.db.GetContext(ctx, &ent, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, appErrors.ErrInternal.WithError(err)
	}

	return &ent, nil
}

func (r *PostgresUserRepository) getIllnesses(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var illnesses []uuid.UUID
	query := queries.User.SelectIllnesses

	err := r.db.SelectContext(ctx, &illnesses, query, userID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	if illnesses == nil {
		return []uuid.UUID{}, nil
	}

	return illnesses, nil
}

func (r *PostgresUserRepository) getAllergies(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	var allergies []uuid.UUID
	query := queries.User.SelectAllergies

	err := r.db.SelectContext(ctx, &allergies, query, userID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	if allergies == nil {
		return []uuid.UUID{}, nil
	}

	return allergies, nil
}

func (r *PostgresUserRepository) FindByLogin(ctx context.Context, login string) (*user.User, error) {
	ent, err := r.findBaseByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if ent == nil {
		return nil, appErrors.ErrUserNotFound.WithSource()
	}

	illnesses, err := r.getIllnesses(ctx, ent.ID)
	if err != nil {
		return nil, err
	}

	allergies, err := r.getAllergies(ctx, ent.ID)
	if err != nil {
		return nil, err
	}

	notify := &user.NotifyInfo{
		Enabled:         ent.NotifyEnabled,
		IntervalMinutes: ent.NotifyIntervalMinutes,
		LastNotifiedAt:  ent.LastNotifiedAt,
	}

	return user.NewUser(
		ent.ID,
		ent.Login,
		ent.Email,
		ent.Password,
		ent.IsAdmin,
		ent.Sex,
		ent.Weight,
		ent.Age,
		ent.IsPregnant,
		ent.IsDriver,
		notify,
		illnesses,
		allergies,
	)
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	ent, err := r.findBaseByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ent == nil {
		return nil, appErrors.ErrUserNotFound.WithSource()
	}

	illnesses, err := r.getIllnesses(ctx, ent.ID)
	if err != nil {
		return nil, err
	}

	allergies, err := r.getAllergies(ctx, ent.ID)
	if err != nil {
		return nil, err
	}

	notify := &user.NotifyInfo{
		Enabled:         ent.NotifyEnabled,
		IntervalMinutes: ent.NotifyIntervalMinutes,
		LastNotifiedAt:  ent.LastNotifiedAt,
	}

	return user.NewUser(
		ent.ID,
		ent.Login,
		ent.Email,
		ent.Password,
		ent.IsAdmin,
		ent.Sex,
		ent.Weight,
		ent.Age,
		ent.IsPregnant,
		ent.IsDriver,
		notify,
		illnesses,
		allergies,
	)
}

func (r *PostgresUserRepository) insertBase(ctx context.Context, tx *sqlx.Tx, u *user.User) error {
	ent := entity.UserEntity{
		ID:         u.ID,
		Login:      u.Login,
		Email:      u.Email,
		Password:   u.Password,
		IsAdmin:    u.IsAdmin,
		Sex:        u.Sex,
		Weight:     u.Weight,
		Age:        u.Age,
		IsPregnant: u.IsPregnant,
		IsDriver:   u.IsDriver,
		NotifyEnabled:         u.Notify.Enabled,
		NotifyIntervalMinutes: u.Notify.IntervalMinutes,
		LastNotifiedAt:        u.Notify.LastNotifiedAt,
	}

	query := queries.User.InsertUser

	_, err := tx.NamedExecContext(ctx, query, ent)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresUserRepository) insertIllnesses(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID, illnesses []uuid.UUID) error {
	if len(illnesses) == 0 {
		return nil
	}

	rows := make([]map[string]interface{}, len(illnesses))
	for i, id := range illnesses {
		rows[i] = map[string]interface{}{
			"user_id":    userID,
			"illness_id": id,
		}
	}

	query := queries.User.InsertIllness
	_, err := tx.NamedExecContext(ctx, query, rows)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresUserRepository) insertAllergies(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID, allergies []uuid.UUID) error {
	if len(allergies) == 0 {
		return nil
	}

	rows := make([]map[string]interface{}, len(allergies))
	for i, id := range allergies {
		rows[i] = map[string]interface{}{
			"user_id":      userID,
			"substance_id": id,
		}
	}

	query := queries.User.InsertAllergy
	_, err := tx.NamedExecContext(ctx, query, rows)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, u *user.User) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	defer tx.Rollback()

	if err := r.insertBase(ctx, tx, u); err != nil {
		return err
	}

	if err := r.insertIllnesses(ctx, tx, u.ID, u.Illnesses); err != nil {
		return err
	}

	if err := r.insertAllergies(ctx, tx, u.ID, u.Allergies); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return appErrors.ErrInternal.WithError(err)
	}

	return nil
}

func (r *PostgresUserRepository) deleteIllnesses(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID) error {
	query := queries.User.DeleteIllnesses
	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresUserRepository) deleteAllergies(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID) error {
	query := queries.User.DeleteAllergies
	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresUserRepository) updateBase(ctx context.Context, tx *sqlx.Tx, u *user.User) error {
	ent := entity.UserEntity{
		ID:         u.ID,
		Login:      u.Login,
		Email:      u.Email,
		Password:   u.Password,
		IsAdmin:    u.IsAdmin,
		Sex:        u.Sex,
		Weight:     u.Weight,
		Age:        u.Age,
		IsPregnant: u.IsPregnant,
		IsDriver:   u.IsDriver,
		NotifyEnabled:         u.Notify.Enabled,
		NotifyIntervalMinutes: u.Notify.IntervalMinutes,
		LastNotifiedAt:        u.Notify.LastNotifiedAt,
	}

	query := queries.User.UpdateUser

	_, err := tx.NamedExecContext(ctx, query, ent)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresUserRepository) FindNotifyEnabled(ctx context.Context) ([]*user.User, error) {
	var ents []entity.UserEntity
	query := queries.User.SelectNotifyEnabled

	err := r.db.SelectContext(ctx, &ents, query)
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	users := make([]*user.User, 0, len(ents))
	for _, ent := range ents {
		notify := &user.NotifyInfo{
			Enabled:         ent.NotifyEnabled,
			IntervalMinutes: ent.NotifyIntervalMinutes,
			LastNotifiedAt:  ent.LastNotifiedAt,
		}
		u, err := user.NewUser(
			ent.ID,
			ent.Login,
			ent.Email,
			ent.Password,
			ent.IsAdmin,
			ent.Sex,
			ent.Weight,
			ent.Age,
			ent.IsPregnant,
			ent.IsDriver,
			notify,
			nil,
			nil,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *PostgresUserRepository) UpdateNotify(ctx context.Context, id uuid.UUID, enabled bool, intervalMinutes int, lastNotifiedAt *time.Time) error {
	ent := entity.UserEntity{
		ID:                     id,
		NotifyEnabled:         enabled,
		NotifyIntervalMinutes: intervalMinutes,
		LastNotifiedAt:        lastNotifiedAt,
	}

	query := queries.User.UpdateNotify

	_, err := r.db.NamedExecContext(ctx, query, ent)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, u *user.User) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	defer tx.Rollback()

	if err := r.updateBase(ctx, tx, u); err != nil {
		return err
	}

	if err := r.deleteIllnesses(ctx, tx, u.ID); err != nil {
		return err
	}
	if err := r.insertIllnesses(ctx, tx, u.ID, u.Illnesses); err != nil {
		return err
	}

	if err := r.deleteAllergies(ctx, tx, u.ID); err != nil {
		return err
	}
	if err := r.insertAllergies(ctx, tx, u.ID, u.Allergies); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return appErrors.ErrInternal.WithError(err)
	}

	return nil
}
