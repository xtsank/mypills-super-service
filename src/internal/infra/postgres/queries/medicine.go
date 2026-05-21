package queries

type MedicineQueries struct {
	SelectByID              string
	SelectByIllness         string
	SelectSubstances        string
	SelectDosages           string
	SelectContraindications string
	SelectRecommendations   string
	InsertMedicine          string
	InsertSubstances        string
	InsertRecommendations   string
	InsertContraindications string
	InsertDosages           string
	DeleteSubstances        string
	DeleteDosages           string
	DeleteRecommendations   string
	DeleteContraindications string
	UpdateMedicine          string
	DeleteMedicine          string
	DeleteDosageRule        string
}

var Medicine = MedicineQueries{
	SelectByID:      `SELECT * FROM Medicine WHERE id = $1`,
	SelectByIllness: `SELECT m.* FROM Medicine m
       JOIN Recommendations r ON m.id = r.medicine_id
       WHERE r.illness_id = $1`,
	SelectSubstances:        `SELECT * FROM Medicine_Substance WHERE medicine_id = $1`,
	SelectDosages:           `SELECT * FROM Dosage WHERE medicine_id = $1`,
	SelectContraindications: `SELECT illness_id FROM Contraindications WHERE medicine_id = $1`,
	SelectRecommendations:   `SELECT illness_id FROM Recommendations WHERE medicine_id = $1`,
	InsertMedicine: `INSERT INTO Medicine  (id, name, expire_time, is_prescription, 
                                   method_of_application, effect_on_pregnant, effect_on_driver, 
                                   form_id, unit_id)
              VALUES (:id, :name, :expire_time, :is_prescription, :method_of_application, :effect_on_pregnant,
                      :effect_on_driver, :form_id, :unit_id)`,
	InsertSubstances: `INSERT INTO Medicine_Substance (medicine_id, substance_id, concentration) 
              VALUES (:medicine_id, :substance_id, :concentration)`,
	InsertRecommendations:   `INSERT INTO Recommendations (medicine_id, illness_id) VALUES (:medicine_id, :illness_id)`,
	InsertContraindications: `INSERT INTO Contraindications (medicine_id, illness_id) VALUES (:medicine_id, :illness_id)`,
	InsertDosages: `INSERT INTO Dosage (id, medicine_id, value_from, value_to, dosage_type, dosage_value, number_of_doses_per_day)
              VALUES (:id, :medicine_id, :value_from, :value_to, :dosage_type, :dosage_value, :number_of_doses_per_day)`,
	DeleteSubstances:        `DELETE FROM Medicine_Substance WHERE medicine_id = $1`,
	DeleteDosages:           `DELETE FROM Dosage WHERE medicine_id = $1`,
	DeleteRecommendations:   `DELETE FROM Recommendations WHERE medicine_id = $1`,
	DeleteContraindications: `DELETE FROM Contraindications WHERE medicine_id = $1`,
	UpdateMedicine: `UPDATE Medicine 
              SET name = :name, expire_time = :expire_time, is_prescription = :is_prescription, 
                  method_of_application = :method_of_application, effect_on_pregnant = :effect_on_pregnant, 
                  effect_on_driver = :effect_on_driver, form_id = :form_id, unit_id = :unit_id
              WHERE id = :id`,
	DeleteMedicine:   `delete from Medicine where id = $1`,
	DeleteDosageRule: `delete from Dosage where id = $1`,
}

