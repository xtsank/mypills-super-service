package queries

type CabinetItemsQueries struct {
	List            string
	Find            string
	LockDuplicate   string
	ExistsDuplicate string
	Create          string
	Update          string
	Delete          string
}

var CabinetItems = CabinetItemsQueries{
	List: `SELECT um.id,um.user_id,um.cabinet_id,um.medicine_id,
		m.name AS medicine_name,m.expire_time AS expire_months,
		um.date_of_manufacture,um.quantity,um.minimum_quantity
		FROM User_Medicine um JOIN Medicine m ON m.id=um.medicine_id
		WHERE um.user_id=$1 AND um.cabinet_id=$2`,
	Find: `SELECT um.id,um.user_id,um.cabinet_id,um.medicine_id,
		m.name AS medicine_name,m.expire_time AS expire_months,
		um.date_of_manufacture,um.quantity,um.minimum_quantity
		FROM User_Medicine um JOIN Medicine m ON m.id=um.medicine_id
		WHERE um.user_id=$1 AND um.cabinet_id=$2 AND um.id=$3`,
	LockDuplicate: `SELECT pg_advisory_xact_lock(hashtextextended($1,0))`,
	ExistsDuplicate: `SELECT EXISTS(SELECT 1 FROM User_Medicine
		WHERE cabinet_id=$1 AND medicine_id=$2 AND date_of_manufacture=$3)`,
	Create: `INSERT INTO User_Medicine(id,user_id,cabinet_id,medicine_id,date_of_manufacture,quantity,minimum_quantity)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
	Update: `UPDATE User_Medicine SET quantity=$4,minimum_quantity=$5,date_of_manufacture=$6
		WHERE id=$1 AND user_id=$2 AND cabinet_id=$3
		AND NOT EXISTS (SELECT 1 FROM User_Medicine other WHERE other.cabinet_id=$3
			AND other.medicine_id=$7 AND other.date_of_manufacture=$6 AND other.id<>$1)`,
	Delete: `DELETE FROM User_Medicine WHERE id=$1 AND user_id=$2 AND cabinet_id=$3`,
}
