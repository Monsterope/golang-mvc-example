package databases

import (
	"fmt"
	"monsterloveshop/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	DB *gorm.DB
}

func (con *DatabaseConfig) Connect() {
	db_host := config.GetEnv("db.host")
	db_port := config.GetEnv("db.port")
	db_name := config.GetEnv("db.database")
	db_user := config.GetEnv("db.username")
	db_pass := config.GetEnv("db.password")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", db_user, db_pass, db_host, db_port, db_name)

	var err error
	con.DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic("cannot connect database please")
	}
}

func (con *DatabaseConfig) AutoMigrate() {
	// con.DB.AutoMigrate(models.User{})
	// con.DB.AutoMigrate(models.Category{})
	// con.DB.AutoMigrate(models.Product{})
}
