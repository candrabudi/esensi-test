package main

import (
	"esensi-test/database"
	"esensi-test/database/migration"
	"esensi-test/internal/factory"
	"esensi-test/internal/http"
	"esensi-test/pkg/util"
	"flag"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	var m string

	database.CreateConnection()

	flag.StringVar(
		&m,
		"m",
		"none",
		`This flag is used for migration`,
	)

	flag.Parse()

	if m == "migrate" {
		migration.Migrate()
	}

	f := factory.NewFactory() // Database instance initialization
	g := gin.New()

	http.NewHttp(g, f)

	if err := g.Run(":" + util.GetEnv("APP_PORT", "8080")); err != nil {
		log.Fatal("Can't start server.")
	}
}
