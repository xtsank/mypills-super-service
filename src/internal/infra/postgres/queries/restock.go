package queries

type RestockQueries struct {
	ListAll string
	Count   string
	List    string
	Find    string
	Create  string
	Update  string
	Delete  string
}

var Restock = RestockQueries{
	ListAll: `SELECT id,user_id,medicine_id,desired_quantity FROM Restock_Items`,
	Count: `SELECT count(*) FROM Restock_Items r JOIN Medicine m ON m.id=r.medicine_id
		WHERE r.user_id=$1 AND ($2::uuid IS NULL OR r.id=$2)`,
	List: `SELECT r.id,r.user_id,r.medicine_id,m.name AS medicine_name,r.desired_quantity
		FROM Restock_Items r JOIN Medicine m ON m.id=r.medicine_id
		WHERE r.user_id=$1 AND ($2::uuid IS NULL OR r.id=$2)
		ORDER BY r.id LIMIT $3 OFFSET $4`,
	Find: `SELECT r.id,r.user_id,r.medicine_id,m.name AS medicine_name,r.desired_quantity
		FROM Restock_Items r JOIN Medicine m ON m.id=r.medicine_id
		WHERE r.user_id=$1 AND r.id=$2`,
	Create: `INSERT INTO Restock_Items(id,user_id,medicine_id,desired_quantity) VALUES ($1,$2,$3,$4)`,
	Update: `UPDATE Restock_Items SET desired_quantity=$3 WHERE id=$1 AND user_id=$2`,
	Delete: `DELETE FROM Restock_Items WHERE id=$1 AND user_id=$2`,
}
