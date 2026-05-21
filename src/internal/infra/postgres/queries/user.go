package queries

type UserQueries struct {
	ExistsByLogin   string
	FindByLogin     string
	FindByID        string
	SelectIllnesses string
	SelectAllergies string
	InsertUser      string
	InsertIllness   string
	InsertAllergy   string
	DeleteIllnesses string
	DeleteAllergies string
	UpdateUser      string
}

var User = UserQueries{
	ExistsByLogin:   `select exists(select 1 from Users where login = $1)`,
	FindByLogin:     `SELECT * FROM Users WHERE login = $1`,
	FindByID:        `SELECT * FROM Users WHERE id = $1`,
	SelectIllnesses: `SELECT illness_id FROM User_Illness WHERE user_id = $1`,
	SelectAllergies: `SELECT substance_id FROM User_Substance WHERE user_id = $1`,
	InsertUser: `INSERT INTO Users (id, login, password, is_admin, sex, weight, age, is_pregnant, is_driver)
              VALUES (:id, :login, :password, :is_admin, :sex, :weight, :age, :is_pregnant, :is_driver)`,
	InsertIllness:   `INSERT INTO User_Illness (user_id, illness_id) VALUES (:user_id, :illness_id)`,
	InsertAllergy:   `INSERT INTO User_Substance (user_id, substance_id) VALUES (:user_id, :substance_id)`,
	DeleteIllnesses: `DELETE FROM User_Illness WHERE user_id = $1`,
	DeleteAllergies: `DELETE FROM User_Substance WHERE user_id = $1`,
	UpdateUser: `UPDATE Users 
              SET sex = :sex, weight = :weight, age = :age, 
                  is_pregnant = :is_pregnant, is_driver = :is_driver
              WHERE id = :id`,
}

