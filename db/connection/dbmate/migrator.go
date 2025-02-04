package dbmate

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
)

func Migrate(ctx context.Context, dbDns string, debug bool) {
	log.Printf("Migrating database: %v", dbDns)
	u, err := url.Parse(dbDns)
	if err != nil {
		panic(err)
	}
	db := dbmate.New(u)

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("err", err)
	}
	log.Printf("Current working directory: %s", dir)

	if err := db.CreateAndMigrate(); err != nil {
		panic(err)
	}
	fmt.Println("DB migration done...!")
}
