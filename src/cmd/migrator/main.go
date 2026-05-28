package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	mongoEntity "github.com/xtsank/mypills-super-service/src/internal/infra/mongo/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type config struct {
	DBUser     string
	DBPassword string
	DBName     string

	PostgresHost string
	PostgresPort string
	MongoHost    string
	MongoPort    string
}

type dictionaryItem struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type pgMedicine struct {
	ID                  uuid.UUID `db:"id"`
	FormID              uuid.UUID `db:"form_id"`
	UnitID              uuid.UUID `db:"unit_id"`
	Name                string    `db:"name"`
	ExpireTime          int       `db:"expire_time"`
	EffectOnDriver      bool      `db:"effect_on_driver"`
	EffectOnPregnant    bool      `db:"effect_on_pregnant"`
	MethodOfApplication string    `db:"method_of_application"`
	IsPrescription      bool      `db:"is_prescription"`
}

type pgDosage struct {
	ID                  uuid.UUID `db:"id"`
	MedicineID          uuid.UUID `db:"medicine_id"`
	ValueFrom           int       `db:"value_from"`
	ValueTo             int       `db:"value_to"`
	DosageType          string    `db:"dosage_type"`
	DosageValue         float32   `db:"dosage_value"`
	NumberOfDosesPerDay int       `db:"number_of_doses_per_day"`
}

type pgMedicineSubstance struct {
	MedicineID    uuid.UUID `db:"medicine_id"`
	SubstanceID   uuid.UUID `db:"substance_id"`
	Concentration float32   `db:"concentration"`
}

type pgMedicineIllness struct {
	MedicineID uuid.UUID `db:"medicine_id"`
	IllnessID  uuid.UUID `db:"illness_id"`
}

type pgUser struct {
	ID                    uuid.UUID  `db:"id"`
	Login                 string     `db:"login"`
	Email                 string     `db:"email"`
	Password              string     `db:"password"`
	IsAdmin               bool       `db:"is_admin"`
	Sex                   bool       `db:"sex"`
	Weight                int        `db:"weight"`
	Age                   int        `db:"age"`
	IsPregnant            bool       `db:"is_pregnant"`
	IsDriver              bool       `db:"is_driver"`
	NotifyEnabled         bool       `db:"notify_enabled"`
	NotifyIntervalMinutes int        `db:"notify_interval_minutes"`
	LastNotifiedAt        *time.Time `db:"last_notified_at"`
}

type pgUserIllness struct {
	UserID    uuid.UUID `db:"user_id"`
	IllnessID uuid.UUID `db:"illness_id"`
}

type pgUserSubstance struct {
	UserID      uuid.UUID `db:"user_id"`
	SubstanceID uuid.UUID `db:"substance_id"`
}

type pgCabinetItem struct {
	ID                uuid.UUID `db:"id"`
	UserID            uuid.UUID `db:"user_id"`
	MedicineID        uuid.UUID `db:"medicine_id"`
	DateOfManufacture time.Time `db:"date_of_manufacture"`
	Quantity          float32   `db:"quantity"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	direction := strings.ToLower(strings.TrimSpace(os.Getenv("MIGRATION_DIRECTION")))
	if direction == "" {
		direction = "pg2mongo"
	}
	cleanup := parseBoolEnv("MIGRATION_CLEANUP", true)

	pg, err := sqlx.Open("pgx", postgresConnString(cfg))
	if err != nil {
		log.Fatalf("postgres open error: %v", err)
	}
	defer pg.Close()

	if err := pg.Ping(); err != nil {
		log.Fatalf("postgres ping error: %v", err)
	}

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(mongoConnString(cfg)))
	if err != nil {
		log.Fatalf("mongo connect error: %v", err)
	}
	defer func() {
		_ = mongoClient.Disconnect(context.Background())
	}()

	ctx := context.Background()
	mongoDB := mongoClient.Database(cfg.DBName)

	switch direction {
	case "pg2mongo":
		if err := migratePgToMongo(ctx, pg, mongoDB); err != nil {
			log.Fatalf("migration error: %v", err)
		}
		log.Println("migration finished")
		if cleanup {
			if err := cleanupPostgres(pg); err != nil {
				log.Fatalf("postgres cleanup error: %v", err)
			}
			log.Println("postgres cleanup finished")
		}
	case "mongo2pg":
		if err := migrateMongoToPostgres(ctx, pg, mongoDB); err != nil {
			log.Fatalf("migration error: %v", err)
		}
		log.Println("migration finished")
		if cleanup {
			if err := cleanupMongo(ctx, mongoDB); err != nil {
				log.Fatalf("mongo cleanup error: %v", err)
			}
			log.Println("mongo cleanup finished")
		}
	default:
		log.Fatalf("unknown MIGRATION_DIRECTION: %s", direction)
	}
}

func loadConfig() (config, error) {
	_ = godotenv.Load()

	cfg := config{
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_NAME"),
		PostgresHost: os.Getenv("POSTGRES_HOST"),
		PostgresPort: os.Getenv("POSTGRES_PORT"),
		MongoHost:    os.Getenv("MONGO_HOST"),
		MongoPort:    os.Getenv("MONGO_PORT"),
	}

	missing := []string{}
	if cfg.DBUser == "" {
		missing = append(missing, "DB_USER")
	}
	if cfg.DBPassword == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if cfg.DBName == "" {
		missing = append(missing, "DB_NAME")
	}
	if cfg.PostgresHost == "" {
		missing = append(missing, "POSTGRES_HOST")
	}
	if cfg.PostgresPort == "" {
		missing = append(missing, "POSTGRES_PORT")
	}
	if cfg.MongoHost == "" {
		missing = append(missing, "MONGO_HOST")
	}
	if cfg.MongoPort == "" {
		missing = append(missing, "MONGO_PORT")
	}
	if len(missing) > 0 {
		return config{}, fmt.Errorf("missing env: %s", missing)
	}
	return cfg, nil
}

func postgresConnString(cfg config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
}

func mongoConnString(cfg config) string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
		cfg.DBUser, cfg.DBPassword, cfg.MongoHost, cfg.MongoPort, cfg.DBName)
}

func parseBoolEnv(key string, fallback bool) bool {
	val := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if val == "" {
		return fallback
	}
	return val == "1" || val == "true" || val == "yes" || val == "y"
}

func migratePgToMongo(ctx context.Context, pg *sqlx.DB, mongoDB *mongo.Database) error {
	log.Println("migrating dictionaries")
	if err := migrateDictionary(ctx, pg, mongoDB, "form", "forms"); err != nil {
		return err
	}
	if err := migrateDictionary(ctx, pg, mongoDB, "unit", "units"); err != nil {
		return err
	}
	if err := migrateDictionary(ctx, pg, mongoDB, "substance", "substances"); err != nil {
		return err
	}
	if err := migrateDictionary(ctx, pg, mongoDB, "illness", "illnesses"); err != nil {
		return err
	}

	log.Println("migrating medicines")
	if err := migrateMedicines(ctx, pg, mongoDB); err != nil {
		return err
	}

	log.Println("migrating users")
	if err := migrateUsers(ctx, pg, mongoDB); err != nil {
		return err
	}

	log.Println("migrating cabinet items")
	if err := migrateCabinetItems(ctx, pg, mongoDB); err != nil {
		return err
	}

	return nil
}

func migrateDictionary(ctx context.Context, pg *sqlx.DB, db *mongo.Database, table, collection string) error {
	var items []dictionaryItem
	query := fmt.Sprintf("SELECT id, name FROM %s ORDER BY name", table)
	if err := pg.Select(&items, query); err != nil {
		return err
	}

	docs := make([]interface{}, 0, len(items))
	for _, item := range items {
		docs = append(docs, mongoEntity.DictionaryItemEntity{ID: item.ID, Name: item.Name})
	}

	return replaceCollection(ctx, db.Collection(collection), docs, collection)
}

func migrateMedicines(ctx context.Context, pg *sqlx.DB, db *mongo.Database) error {
	var meds []pgMedicine
	if err := pg.Select(&meds, "SELECT id, form_id, unit_id, name, expire_time, effect_on_driver, effect_on_pregnant, method_of_application, is_prescription FROM medicine"); err != nil {
		return err
	}

	var dosages []pgDosage
	if err := pg.Select(&dosages, "SELECT id, medicine_id, value_from, value_to, dosage_type, dosage_value, number_of_doses_per_day FROM dosage"); err != nil {
		return err
	}

	var subs []pgMedicineSubstance
	if err := pg.Select(&subs, "SELECT medicine_id, substance_id, concentration FROM medicine_substance"); err != nil {
		return err
	}

	var recs []pgMedicineIllness
	if err := pg.Select(&recs, "SELECT medicine_id, illness_id FROM recommendations"); err != nil {
		return err
	}

	var cons []pgMedicineIllness
	if err := pg.Select(&cons, "SELECT medicine_id, illness_id FROM contraindications"); err != nil {
		return err
	}

	dosageMap := make(map[uuid.UUID][]mongoEntity.Dosage)
	for _, d := range dosages {
		dosageMap[d.MedicineID] = append(dosageMap[d.MedicineID], mongoEntity.Dosage{
			ID:                  d.ID,
			ValueFrom:           d.ValueFrom,
			ValueTo:             d.ValueTo,
			DosageType:          d.DosageType,
			DosageValue:         d.DosageValue,
			NumberOfDosesPerDay: d.NumberOfDosesPerDay,
		})
	}

	subMap := make(map[uuid.UUID][]mongoEntity.Substance)
	for _, s := range subs {
		subMap[s.MedicineID] = append(subMap[s.MedicineID], mongoEntity.Substance{
			SubstanceID:   s.SubstanceID,
			Concentration: s.Concentration,
		})
	}

	recMap := make(map[uuid.UUID][]uuid.UUID)
	for _, r := range recs {
		recMap[r.MedicineID] = append(recMap[r.MedicineID], r.IllnessID)
	}

	conMap := make(map[uuid.UUID][]uuid.UUID)
	for _, c := range cons {
		conMap[c.MedicineID] = append(conMap[c.MedicineID], c.IllnessID)
	}

	docs := make([]interface{}, 0, len(meds))
	for _, med := range meds {
		docs = append(docs, mongoEntity.MedicineEntity{
			ID:                  med.ID,
			FormID:              med.FormID,
			UnitID:              med.UnitID,
			Name:                med.Name,
			ExpireTime:          med.ExpireTime,
			EffectOnDriver:      med.EffectOnDriver,
			EffectOnPregnant:    med.EffectOnPregnant,
			MethodOfApplication: med.MethodOfApplication,
			IsPrescription:      med.IsPrescription,
			Substances:          subMap[med.ID],
			Dosages:             dosageMap[med.ID],
			Contraindications:   conMap[med.ID],
			Recommendations:     recMap[med.ID],
		})
	}

	return replaceCollection(ctx, db.Collection("medicines"), docs, "medicines")
}

func migrateUsers(ctx context.Context, pg *sqlx.DB, db *mongo.Database) error {
	var users []pgUser
	if err := pg.Select(&users, "SELECT id, login, email, password, is_admin, sex, weight, age, is_pregnant, is_driver, notify_enabled, notify_interval_minutes, last_notified_at FROM users"); err != nil {
		return err
	}

	var userIllness []pgUserIllness
	if err := pg.Select(&userIllness, "SELECT user_id, illness_id FROM user_illness"); err != nil {
		return err
	}

	var userSubstances []pgUserSubstance
	if err := pg.Select(&userSubstances, "SELECT user_id, substance_id FROM user_substance"); err != nil {
		return err
	}

	illnessMap := make(map[uuid.UUID][]uuid.UUID)
	for _, item := range userIllness {
		illnessMap[item.UserID] = append(illnessMap[item.UserID], item.IllnessID)
	}

	allergyMap := make(map[uuid.UUID][]uuid.UUID)
	for _, item := range userSubstances {
		allergyMap[item.UserID] = append(allergyMap[item.UserID], item.SubstanceID)
	}

	docs := make([]interface{}, 0, len(users))
	for _, u := range users {
		docs = append(docs, mongoEntity.UserEntity{
			ID:                    u.ID,
			Login:                 u.Login,
			Email:                 u.Email,
			Password:              u.Password,
			IsAdmin:               u.IsAdmin,
			Sex:                   u.Sex,
			Weight:                u.Weight,
			Age:                   u.Age,
			IsPregnant:            u.IsPregnant,
			IsDriver:              u.IsDriver,
			NotifyEnabled:         u.NotifyEnabled,
			NotifyIntervalMinutes: u.NotifyIntervalMinutes,
			LastNotifiedAt:        u.LastNotifiedAt,
			Illnesses:             illnessMap[u.ID],
			Allergies:             allergyMap[u.ID],
		})
	}

	return replaceCollection(ctx, db.Collection("users"), docs, "users")
}

func migrateCabinetItems(ctx context.Context, pg *sqlx.DB, db *mongo.Database) error {
	var items []pgCabinetItem
	if err := pg.Select(&items, "SELECT id, user_id, medicine_id, date_of_manufacture, quantity FROM user_medicine"); err != nil {
		return err
	}

	docs := make([]interface{}, 0, len(items))
	for _, item := range items {
		docs = append(docs, mongoEntity.CabinetItemEntity{
			ID:                item.ID,
			UserID:            item.UserID,
			MedicineID:        item.MedicineID,
			DateOfManufacture: item.DateOfManufacture,
			Quantity:          item.Quantity,
		})
	}

	return replaceCollection(ctx, db.Collection("cabinet_items"), docs, "cabinet_items")
}

func replaceCollection(ctx context.Context, collection *mongo.Collection, docs []interface{}, label string) error {
	if _, err := collection.DeleteMany(ctx, bson.M{}); err != nil {
		return err
	}
	if len(docs) == 0 {
		log.Printf("%s collection is empty", label)
		return nil
	}
	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}

func migrateMongoToPostgres(ctx context.Context, pg *sqlx.DB, mongoDB *mongo.Database) error {
	forms, err := fetchDictionary(ctx, mongoDB, "forms")
	if err != nil {
		return err
	}
	units, err := fetchDictionary(ctx, mongoDB, "units")
	if err != nil {
		return err
	}
	substances, err := fetchDictionary(ctx, mongoDB, "substances")
	if err != nil {
		return err
	}
	illnesses, err := fetchDictionary(ctx, mongoDB, "illnesses")
	if err != nil {
		return err
	}
	medicines, err := fetchMedicines(ctx, mongoDB)
	if err != nil {
		return err
	}
	users, err := fetchUsers(ctx, mongoDB)
	if err != nil {
		return err
	}
	items, err := fetchCabinetItems(ctx, mongoDB)
	if err != nil {
		return err
	}

	tx, err := pg.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := truncatePostgres(tx); err != nil {
		return err
	}

	if err := insertDictionaries(tx, "form", forms); err != nil {
		return err
	}
	if err := insertDictionaries(tx, "unit", units); err != nil {
		return err
	}
	if err := insertDictionaries(tx, "substance", substances); err != nil {
		return err
	}
	if err := insertDictionaries(tx, "illness", illnesses); err != nil {
		return err
	}

	if err := insertMedicines(tx, medicines); err != nil {
		return err
	}
	if err := insertUsers(tx, users); err != nil {
		return err
	}
	if err := insertCabinetItems(tx, items); err != nil {
		return err
	}

	return tx.Commit()
}

func fetchDictionary(ctx context.Context, db *mongo.Database, collection string) ([]mongoEntity.DictionaryItemEntity, error) {
	cursor, err := db.Collection(collection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []mongoEntity.DictionaryItemEntity
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func fetchMedicines(ctx context.Context, db *mongo.Database) ([]mongoEntity.MedicineEntity, error) {
	cursor, err := db.Collection("medicines").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []mongoEntity.MedicineEntity
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func fetchUsers(ctx context.Context, db *mongo.Database) ([]mongoEntity.UserEntity, error) {
	cursor, err := db.Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []mongoEntity.UserEntity
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func fetchCabinetItems(ctx context.Context, db *mongo.Database) ([]mongoEntity.CabinetItemEntity, error) {
	cursor, err := db.Collection("cabinet_items").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []mongoEntity.CabinetItemEntity
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func insertDictionaries(tx *sqlx.Tx, table string, items []mongoEntity.DictionaryItemEntity) error {
	query := fmt.Sprintf("INSERT INTO %s (id, name) VALUES ($1, $2)", table)
	for _, item := range items {
		if _, err := tx.Exec(query, item.ID, item.Name); err != nil {
			return err
		}
	}
	return nil
}

func insertMedicines(tx *sqlx.Tx, meds []mongoEntity.MedicineEntity) error {
	for _, med := range meds {
		if _, err := tx.Exec(
			"INSERT INTO medicine (id, form_id, unit_id, name, expire_time, effect_on_driver, effect_on_pregnant, method_of_application, is_prescription) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)",
			med.ID,
			med.FormID,
			med.UnitID,
			med.Name,
			med.ExpireTime,
			med.EffectOnDriver,
			med.EffectOnPregnant,
			med.MethodOfApplication,
			med.IsPrescription,
		); err != nil {
			return err
		}

		for _, d := range med.Dosages {
			if _, err := tx.Exec(
				"INSERT INTO dosage (id, medicine_id, value_from, value_to, dosage_type, dosage_value, number_of_doses_per_day) VALUES ($1,$2,$3,$4,$5,$6,$7)",
				d.ID,
				med.ID,
				d.ValueFrom,
				d.ValueTo,
				d.DosageType,
				d.DosageValue,
				d.NumberOfDosesPerDay,
			); err != nil {
				return err
			}
		}

		for _, s := range med.Substances {
			if _, err := tx.Exec(
				"INSERT INTO medicine_substance (medicine_id, substance_id, concentration) VALUES ($1,$2,$3)",
				med.ID,
				s.SubstanceID,
				s.Concentration,
			); err != nil {
				return err
			}
		}

		for _, id := range med.Recommendations {
			if _, err := tx.Exec(
				"INSERT INTO recommendations (medicine_id, illness_id) VALUES ($1,$2)",
				med.ID,
				id,
			); err != nil {
				return err
			}
		}

		for _, id := range med.Contraindications {
			if _, err := tx.Exec(
				"INSERT INTO contraindications (medicine_id, illness_id) VALUES ($1,$2)",
				med.ID,
				id,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func insertUsers(tx *sqlx.Tx, users []mongoEntity.UserEntity) error {
	for _, u := range users {
		if _, err := tx.Exec(
			"INSERT INTO users (id, login, email, password, is_admin, sex, weight, age, is_pregnant, is_driver, notify_enabled, notify_interval_minutes, last_notified_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)",
			u.ID,
			u.Login,
			u.Email,
			u.Password,
			u.IsAdmin,
			u.Sex,
			u.Weight,
			u.Age,
			u.IsPregnant,
			u.IsDriver,
			u.NotifyEnabled,
			u.NotifyIntervalMinutes,
			u.LastNotifiedAt,
		); err != nil {
			return err
		}

		for _, id := range u.Illnesses {
			if _, err := tx.Exec("INSERT INTO user_illness (user_id, illness_id) VALUES ($1,$2)", u.ID, id); err != nil {
				return err
			}
		}
		for _, id := range u.Allergies {
			if _, err := tx.Exec("INSERT INTO user_substance (user_id, substance_id) VALUES ($1,$2)", u.ID, id); err != nil {
				return err
			}
		}
	}
	return nil
}

func insertCabinetItems(tx *sqlx.Tx, items []mongoEntity.CabinetItemEntity) error {
	for _, item := range items {
		if _, err := tx.Exec(
			"INSERT INTO user_medicine (id, user_id, medicine_id, date_of_manufacture, quantity) VALUES ($1,$2,$3,$4,$5)",
			item.ID,
			item.UserID,
			item.MedicineID,
			item.DateOfManufacture,
			item.Quantity,
		); err != nil {
			return err
		}
	}
	return nil
}

func truncatePostgres(tx *sqlx.Tx) error {
	query := `
		TRUNCATE TABLE
			user_medicine,
			user_substance,
			user_illness,
			contraindications,
			recommendations,
			medicine_substance,
			dosage,
			medicine,
			users,
			illness,
			substance,
			unit,
			form
		RESTART IDENTITY CASCADE;
	`
	_, err := tx.Exec(query)
	return err
}

func cleanupMongo(ctx context.Context, db *mongo.Database) error {
	collections := []string{"cabinet_items", "users", "medicines", "illnesses", "substances", "units", "forms"}
	for _, name := range collections {
		if _, err := db.Collection(name).DeleteMany(ctx, bson.M{}); err != nil {
			return err
		}
	}
	return nil
}

func cleanupPostgres(pg *sqlx.DB) error {
	tx, err := pg.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := truncatePostgres(tx); err != nil {
		return err
	}
	return tx.Commit()
}
