package queries

type DictionaryQueries struct {
	ListIllnesses  string
	ListSubstances string
	ListForms      string
	ListUnits      string
	ListMedicines  string
}

var Dictionary = DictionaryQueries{
	ListIllnesses:  `SELECT id, name FROM Illness ORDER BY name`,
	ListSubstances: `SELECT id, name FROM Substance ORDER BY name`,
	ListForms:      `SELECT id, name FROM Form ORDER BY name`,
	ListUnits:      `SELECT id, name FROM Unit ORDER BY name`,
	ListMedicines:  `SELECT id, name FROM Medicine ORDER BY name`,
}

