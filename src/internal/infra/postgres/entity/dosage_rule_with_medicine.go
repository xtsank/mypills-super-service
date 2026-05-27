package entity

import "github.com/google/uuid"

type DosageRuleWithMedicineEntity struct {
	ID                  uuid.UUID `db:"id"`
	MedicineID          uuid.UUID `db:"medicine_id"`
	MedicineName        string    `db:"medicine_name"`
	ValueFrom           int       `db:"value_from"`
	ValueTo             int       `db:"value_to"`
	DosageType          string    `db:"dosage_type"`
	DosageValue         float32   `db:"dosage_value"`
	NumberOfDosesPerDay int       `db:"number_of_doses_per_day"`
}

