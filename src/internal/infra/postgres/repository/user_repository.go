package repository

import (
	"context"
	"database/sql"
	"errors"

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

	return user.NewUser(
		ent.ID,
		ent.Login,
		ent.Password,
		ent.IsAdmin,
		ent.Sex,
		ent.Weight,
		ent.Age,
		ent.IsPregnant,
		ent.IsDriver,
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

	return user.NewUser(
		ent.ID,
		ent.Login,
		ent.Password,
		ent.IsAdmin,
		ent.Sex,
		ent.Weight,
		ent.Age,
		ent.IsPregnant,
		ent.IsDriver,
		illnesses,
		allergies,
	)
}

func (r *PostgresUserRepository) insertBase(ctx context.Context, tx *sqlx.Tx, u *user.User) error {
	ent := entity.UserEntity{
		ID:         u.ID,
		Login:      u.Login,
		Password:   u.Password,
		IsAdmin:    u.IsAdmin,
		Sex:        u.Sex,
		Weight:     u.Weight,
		Age:        u.Age,
		IsPregnant: u.IsPregnant,
		IsDriver:   u.IsDriver,
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
		Password:   u.Password,
		IsAdmin:    u.IsAdmin,
		Sex:        u.Sex,
		Weight:     u.Weight,
		Age:        u.Age,
		IsPregnant: u.IsPregnant,
		IsDriver:   u.IsDriver,
	}

		  query := queries.User.UpdateUser

	_, err := tx.NamedExecContext(ctx, query, ent)
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
