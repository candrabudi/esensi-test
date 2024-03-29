package migration

import (
	"esensi-test/database"
	"esensi-test/internal/models"
)

var tables = []interface{}{
	&models.User{},
	&models.UserSession{},
	&models.Item{},
	&models.Customer{},
	&models.Invoice{},
	&models.InvoiceDetail{},
}

func Migrate() {
	conn := database.GetConnection() // Get db connection
	conn.AutoMigrate(tables...)      // migrate the tables
}
