package testutil

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
)

func CabinetRows(items ...*cabinet_item.CabinetItem) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "user_id", "medicine_id", "date_of_manufacture", "quantity"})
	for _, item := range items {
		rows.AddRow(item.ID.String(), item.UserID.String(), item.MedicineID.String(), item.DateOfManufacture, item.Quantity)
	}
	return rows
}

func UserBaseRows(users ...*user.User) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id", "login", "email", "password", "is_admin", "sex", "weight", "age",
		"is_pregnant", "is_driver", "notify_enabled", "notify_interval_minutes", "last_notified_at",
	})
	for _, value := range users {
		rows.AddRow(
			value.ID.String(), value.Login, value.Email, value.Password, value.IsAdmin, value.Sex,
			value.Weight, value.Age, value.IsPregnant, value.IsDriver, value.Notify.Enabled,
			value.Notify.IntervalMinutes, value.Notify.LastNotifiedAt,
		)
	}
	return rows
}

func MedicineBaseRows(medicines ...*medicine.Medicine) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id", "form_id", "unit_id", "name", "expire_time", "effect_on_driver",
		"effect_on_pregnant", "method_of_application", "is_prescription",
	})
	for _, value := range medicines {
		rows.AddRow(
			value.ID.String(), value.Form.String(), value.Unit.String(), value.Name, value.ExpireTime,
			value.EffectOnDriver, value.EffectOnPregnant, value.MethodOfApplication, value.IsPrescription,
		)
	}
	return rows
}

func SubstanceRows(value *medicine.Medicine) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "medicine_id", "substance_id", "concentration"})
	for _, substance := range value.Substances {
		rowID := uuid.NewSHA1(uuid.Nil, []byte(value.ID.String()+substance.ID.String()))
		rows.AddRow(rowID.String(), value.ID.String(), substance.ID.String(), substance.Concentration)
	}
	return rows
}

func DosageRows(value *medicine.Medicine) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id", "medicine_id", "value_from", "value_to", "dosage_type", "dosage_value", "number_of_doses_per_day",
	})
	for _, dosage := range value.Dosages {
		rows.AddRow(
			dosage.ID.String(), value.ID.String(), dosage.ValueFrom, dosage.ValueTo,
			string(dosage.Type), dosage.DosageValue, dosage.NumberOfDosesPerDay,
		)
	}
	return rows
}

func UUIDRows(column string, ids ...uuid.UUID) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{column})
	for _, id := range ids {
		rows.AddRow(id.String())
	}
	return rows
}
