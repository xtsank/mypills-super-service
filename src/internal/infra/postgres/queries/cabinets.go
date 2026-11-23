package queries

type CabinetQueries struct {
	EnsureHome string
	ListAll    string
	Count      string
	List       string
	Find       string
	Create     string
	Rename     string
	Delete     string
}

var Cabinet = CabinetQueries{
	EnsureHome: `INSERT INTO Cabinets(id,user_id,name)
        VALUES ($1, $2, 'Домашняя аптечка')
        ON CONFLICT (user_id,name) DO UPDATE SET name=EXCLUDED.name
        RETURNING id`,
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
