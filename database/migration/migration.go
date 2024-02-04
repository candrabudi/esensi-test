package migration

import (
	"esensi-test/database"
)

var tables = []interface{}{}

func Migrate() {
	conn := database.GetConnection() // Get db connection
	conn.AutoMigrate(tables...)      // migrate the tables
}
