package dao

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type DBHost struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

var Dsn = make(map[string]*DBHost, 1)
var DBGroup *sql.DB

func InitDataBase() {
	dsngame := Dsn["group"].Username + ":" + Dsn["group"].Password + "@tcp" + "(" + Dsn["group"].Host + ":" + Dsn["group"].Port + ")" + "/" + Dsn["group"].DBName + "?charset=utf8mb4&loc=Local"
	c, err := sql.Open("mysql", dsngame)
	if err != nil {
		log.Fatal(err)
	}
	DBGroup = c

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err = DBGroup.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	DBGroup.SetConnMaxLifetime(9 * time.Second)
}
