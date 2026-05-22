package queries

type CabinetItemQueries struct {
	SelectByUserID      string
	SelectByUserIDWithName string
	SelectExistingByKey string
	SelectByID          string
	SelectExpiredByUserID string
	UpdateQuantity      string
	InsertItem          string
	DeleteByID          string
}

var CabinetItem = CabinetItemQueries{
	SelectByUserID:      `SELECT * FROM User_Medicine WHERE user_id = $1`,
  SelectByUserIDWithName: `SELECT um.id, um.user_id, um.medicine_id, m.name AS medicine_name,
			  um.date_of_manufacture, um.quantity,
			  (um.date_of_manufacture + (m.expire_time || ' months')::interval) AS expires_at
			  FROM User_Medicine um
			  JOIN Medicine m ON m.id = um.medicine_id
			  WHERE um.user_id = $1`,
	SelectExistingByKey: `SELECT * FROM User_Medicine 
              WHERE user_id = $1 AND medicine_id = $2 AND date_of_manufacture = $3::date`,
	SelectByID:          `SELECT * FROM User_Medicine WHERE id = $1`,
	SelectExpiredByUserID: `SELECT um.id, um.user_id, um.medicine_id, m.name AS medicine_name, um.date_of_manufacture, um.quantity,
              (um.date_of_manufacture + (m.expire_time || ' months')::interval) AS expires_at
              FROM User_Medicine um
              JOIN Medicine m ON m.id = um.medicine_id
              WHERE um.user_id = $1
                AND (um.date_of_manufacture + (m.expire_time || ' months')::interval) < CURRENT_DATE`,
	UpdateQuantity:      `UPDATE User_Medicine SET quantity = :quantity WHERE id = :id`,
	InsertItem: `INSERT INTO User_Medicine (id, user_id, medicine_id, date_of_manufacture, quantity)
              VALUES (:id, :user_id, :medicine_id, :date_of_manufacture, :quantity)`,
	DeleteByID: `DELETE FROM User_Medicine WHERE id = $1`,
}
