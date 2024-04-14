package database

import (
	"avito-app/models"
	"fmt"
	"log"
	"strconv"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg models.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode, cfg.TimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Banner{}); err != nil {
		log.Fatal(err)
	}
	DB = db
}

func Migration() {
	p, _ := bcrypt.GenerateFromPassword([]byte("12345"), 14)
	DB.Create(&models.User{
		Login:    "admin",
		Password: p,
		Role:     "admin",
	})
	p, _ = bcrypt.GenerateFromPassword([]byte("123"), 14)
	DB.Create(&models.User{
		Login:    "user",
		Password: p,
		Role:     "user",
	})
}

func MigrationData() {
	FeatureIDS := [10]int{1, 1, 3, 4, 4, 4, 7, 8, 8, 10}
	TagIDS := [10][]int64{{1, 2}, {1, 3}, {2, 3}, {4, 6}, {5, 7}, {5, 8}, {5, 9}, {8, 1}, {7, 5}, {6, 9}}
	for index, value := range FeatureIDS {
		is_active := true
		if (index+1)%2 == 0 {
			is_active = false
		}
		content := fmt.Sprintf(`{"title": "%s", "text": "%s", "url": "%s"}`, "Title "+strconv.Itoa(index), "Text "+strconv.Itoa(index), "Url "+strconv.Itoa(index))
		DB.Create(&models.Banner{
			IsActive:  is_active,
			FeatureID: uint(value),
			TagIds:    pq.Int64Array(TagIDS[index]),
			Content:   datatypes.JSON([]byte(content)),
		})
	}
}
