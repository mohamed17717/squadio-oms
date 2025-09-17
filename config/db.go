package config

import (
	"errors"
	"oms-services/utils"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		utils.FailOnError(errors.New("missing env variable: DATABASE_URL"), "DATABASE_URL not set")
	}

	DB, err = gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	utils.FailOnError(err, "Cannot connect to database")

}
