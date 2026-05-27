package queries

type UserQueries struct {
	ExistsByLogin    string
	FindByLogin      string
	FindByID         string
	SelectNotifyEnabled string
	SelectIllnesses  string
	SelectAllergies  string
	InsertUser       string
	InsertIllness    string
	InsertAllergy    string
	DeleteIllnesses  string
	DeleteAllergies  string
	UpdateUser       string
	UpdateNotify     string
}

var User = UserQueries{
	ExistsByLogin:    `select exists(select 1 from Users where login = $1)`,
	FindByLogin:      `SELECT * FROM Users WHERE login = $1`,
	FindByID:         `SELECT * FROM Users WHERE id = $1`,
	SelectNotifyEnabled: `SELECT * FROM Users WHERE notify_enabled = true`,
	SelectIllnesses:  `SELECT illness_id FROM User_Illness WHERE user_id = $1`,
	SelectAllergies:  `SELECT substance_id FROM User_Substance WHERE user_id = $1`,
	InsertUser: `INSERT INTO Users (id, login, email, password, is_admin, sex, weight, age, is_pregnant, is_driver, notify_enabled, notify_interval_minutes, last_notified_at)
              VALUES (:id, :login, :email, :password, :is_admin, :sex, :weight, :age, :is_pregnant, :is_driver, :notify_enabled, :notify_interval_minutes, :last_notified_at)`,
	InsertIllness:    `INSERT INTO User_Illness (user_id, illness_id) VALUES (:user_id, :illness_id)`,
	InsertAllergy:    `INSERT INTO User_Substance (user_id, substance_id) VALUES (:user_id, :substance_id)`,
	DeleteIllnesses:  `DELETE FROM User_Illness WHERE user_id = $1`,
	DeleteAllergies:  `DELETE FROM User_Substance WHERE user_id = $1`,
	UpdateUser: `UPDATE Users 
              SET email = :email, sex = :sex, weight = :weight, age = :age, 
                  is_pregnant = :is_pregnant, is_driver = :is_driver,
                  notify_enabled = :notify_enabled, notify_interval_minutes = :notify_interval_minutes,
                  last_notified_at = :last_notified_at
              WHERE id = :id`,
	UpdateNotify: `UPDATE Users
              SET notify_enabled = :notify_enabled, notify_interval_minutes = :notify_interval_minutes,
                  last_notified_at = :last_notified_at
              WHERE id = :id`,
}
