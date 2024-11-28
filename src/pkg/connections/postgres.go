package connections

import (
	"dealls-dating-app/src/models"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgreConnection() (pg *gorm.DB, err error) {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("PG_DB_HOST"),
		os.Getenv("PG_DB_USER"),
		os.Getenv("PG_DB_PASSWORD"),
		os.Getenv("PG_DB_NAME"),
		os.Getenv("PG_DB_PORT"),
	)

	var newLog logger.Interface

	if os.Getenv("ENVIRONMENT") != "production" {
		newLog = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Disable color
			},
		)
	}

	pg, err = gorm.Open(postgres.Open(config), &gorm.Config{
		Logger: newLog,
	})
	if err != nil {
		return
	}

	fmt.Println("Success connected to database")

	if err = pg.Debug().AutoMigrate(models.User{}, models.Swipe{}); err != nil {
		return
	}

	if pg.Migrator().HasTable(models.User{}) {
		if err = pg.First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			// Seeding data
			if err = seedUser(pg); err != nil {
				return
			}
			fmt.Println("Success seeding data user")
		}
	}

	return
}

func seedUser(pg *gorm.DB) (err error) {
	type seed struct {
		email, firstName, lastName, sex, password string
	}

	seedData := []seed{
		{"dummy1@mail.com", "Dummy", "1", "male", "Test1234"},
		{"dummy2@mail.com", "Dummy", "2", "male", "Test1234"},
		{"dummy3@mail.com", "Dummy", "3", "male", "Test1234"},
		{"dummy4@mail.com", "Dummy", "4", "male", "Test1234"},
		{"dummy5@mail.com", "Dummy", "5", "male", "Test1234"},
		{"dummy6@mail.com", "Dummy", "6", "male", "Test1234"},
		{"dummy7@mail.com", "Dummy", "7", "male", "Test1234"},
		{"dummy8@mail.com", "Dummy", "8", "male", "Test1234"},
		{"dummy9@mail.com", "Dummy", "9", "male", "Test1234"},
		{"dummy10@mail.com", "Dummy", "10", "male", "Test1234"},
		{"dummy11@mail.com", "Dummy", "11", "female", "Test1234"},
		{"dummy12@mail.com", "Dummy", "12", "female", "Test1234"},
		{"dummy13@mail.com", "Dummy", "13", "female", "Test1234"},
		{"dummy14@mail.com", "Dummy", "14", "female", "Test1234"},
		{"dummy15@mail.com", "Dummy", "15", "female", "Test1234"},
		{"dummy16@mail.com", "Dummy", "16", "female", "Test1234"},
		{"dummy17@mail.com", "Dummy", "17", "female", "Test1234"},
		{"dummy18@mail.com", "Dummy", "18", "female", "Test1234"},
		{"dummy19@mail.com", "Dummy", "19", "female", "Test1234"},
		{"dummy20@mail.com", "Dummy", "20", "female", "Test1234"},
	}

	users := []models.User{}
	for _, data := range seedData {
		users = append(users, models.User{
			Email:     data.email,
			FirstName: data.firstName,
			LastName:  data.lastName,
			Sex:       data.sex,
			Password:  data.password,
		})
	}

	return pg.Model(&models.User{}).Create(&users).Error
}
