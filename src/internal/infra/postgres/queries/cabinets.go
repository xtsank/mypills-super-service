package queries

type CabinetQueries struct {
	EnsureLegacy         string
	AssignLegacyItems    string
	EnsureAllLegacy      string
	AssignAllLegacyItems string
	ListAll              string
	Count                string
	List                 string
	Find                 string
	Create               string
	Rename               string
	Delete               string
}

var Cabinet = CabinetQueries{
	EnsureLegacy: `INSERT INTO Cabinets(id,user_id,name)
		SELECT gen_random_uuid(), $1, 'Домашняя аптечка'
		WHERE EXISTS (SELECT 1 FROM User_Medicine WHERE user_id=$1 AND cabinet_id IS NULL)
		ON CONFLICT (user_id,name) DO NOTHING`,
	AssignLegacyItems: `UPDATE User_Medicine um SET cabinet_id=c.id
		FROM Cabinets c WHERE um.user_id=$1 AND um.cabinet_id IS NULL
		AND c.user_id=um.user_id AND c.name='Домашняя аптечка'`,
	EnsureAllLegacy: `INSERT INTO Cabinets(id,user_id,name)
		SELECT gen_random_uuid(),u.user_id,'Домашняя аптечка'
		FROM (SELECT DISTINCT user_id FROM User_Medicine WHERE cabinet_id IS NULL) u
		ON CONFLICT (user_id,name) DO NOTHING`,
	AssignAllLegacyItems: `UPDATE User_Medicine um SET cabinet_id=c.id
		FROM Cabinets c WHERE um.cabinet_id IS NULL AND um.user_id=c.user_id AND c.name='Домашняя аптечка'`,
	ListAll: `SELECT id,user_id,name FROM Cabinets`,
	Count:   `SELECT count(*) FROM Cabinets WHERE user_id=$1 AND ($2::uuid IS NULL OR id=$2)`,
	List: `SELECT id,user_id,name FROM Cabinets
		WHERE user_id=$1 AND ($2::uuid IS NULL OR id=$2)
		ORDER BY name,id LIMIT $3 OFFSET $4`,
	Find:   `SELECT id,user_id,name FROM Cabinets WHERE id=$1 AND user_id=$2`,
	Create: `INSERT INTO Cabinets(id,user_id,name) VALUES ($1,$2,$3)`,
	Rename: `UPDATE Cabinets SET name=$3 WHERE id=$1 AND user_id=$2`,
	Delete: `DELETE FROM Cabinets c WHERE c.id=$1 AND c.user_id=$2
		AND NOT EXISTS (SELECT 1 FROM User_Medicine WHERE cabinet_id=c.id)`,
}
