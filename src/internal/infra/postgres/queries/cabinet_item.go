package queries

type CabinetItemQueries struct {
	SelectByUserID      string
	SelectExistingByKey string
	SelectByID          string
	UpdateQuantity      string
	InsertItem          string
	DeleteByID          string
}

var CabinetItem = CabinetItemQueries{
	SelectByUserID:      `SELECT * FROM User_Medicine WHERE user_id = $1`,
	SelectExistingByKey: `SELECT * FROM User_Medicine 
              WHERE user_id = $1 AND medicine_id = $2 AND date_of_manufacture = $3::date`,
	SelectByID:          `SELECT * FROM User_Medicine WHERE id = $1`,
	UpdateQuantity:      `UPDATE User_Medicine SET quantity = :quantity WHERE id = :id`,
	InsertItem: `INSERT INTO User_Medicine (id, user_id, medicine_id, date_of_manufacture, quantity)
              VALUES (:id, :user_id, :medicine_id, :date_of_manufacture, :quantity)`,
	DeleteByID: `DELETE FROM User_Medicine WHERE id = $1`,
}

