package testutil

const sqliteEnableForeignKeys = "PRAGMA foreign_keys = ON"

var sqliteMedicineSchema = []string{
	`CREATE TABLE Medicine (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		expire_time INTEGER NOT NULL,
		is_prescription BOOLEAN NOT NULL,
		method_of_application TEXT NOT NULL,
		effect_on_pregnant BOOLEAN NOT NULL,
		effect_on_driver BOOLEAN NOT NULL,
		form_id TEXT NOT NULL,
		unit_id TEXT NOT NULL
	)`,
	`CREATE TABLE Medicine_Substance (
		id TEXT PRIMARY KEY DEFAULT (
			lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-' ||
			lower(hex(randomblob(2))) || '-' || lower(hex(randomblob(2))) || '-' ||
			lower(hex(randomblob(6)))
		),
		medicine_id TEXT NOT NULL REFERENCES Medicine(id),
		substance_id TEXT NOT NULL,
		concentration REAL NOT NULL
	)`,
	`CREATE TABLE Dosage (
		id TEXT PRIMARY KEY,
		medicine_id TEXT NOT NULL REFERENCES Medicine(id),
		value_from INTEGER NOT NULL,
		value_to INTEGER NOT NULL,
		dosage_type TEXT NOT NULL,
		dosage_value REAL NOT NULL,
		number_of_doses_per_day INTEGER NOT NULL
	)`,
	`CREATE TABLE Recommendations (
		medicine_id TEXT NOT NULL REFERENCES Medicine(id),
		illness_id TEXT NOT NULL
	)`,
	`CREATE TABLE Contraindications (
		medicine_id TEXT NOT NULL REFERENCES Medicine(id),
		illness_id TEXT NOT NULL
	)`,
}

var sqliteCabinetSchema = []string{
	`CREATE TABLE User_Medicine (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		medicine_id TEXT NOT NULL,
		date_of_manufacture TIMESTAMP NOT NULL,
		quantity REAL NOT NULL
	)`,
}
