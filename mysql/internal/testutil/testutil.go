package testutil

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/nasjp/parabase"
	"github.com/nasjp/parabase/mysql"
)

var (
	DB = setupDB()
)

func setupDB() *sql.DB {
	cfg, _ := GetCfg()
	db, err := sql.Open(cfg.DriverName, cfg.DataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

func Cleanup(db *sql.DB, prefix string) error {
	dbNames := []string{
		fmt.Sprintf("%s_management", prefix),
		fmt.Sprintf("%s_test_1", prefix),
		fmt.Sprintf("%s_test_2", prefix),
		fmt.Sprintf("%s_test_3", prefix),
		fmt.Sprintf("%s_test_4", prefix),
		fmt.Sprintf("%s_test_5", prefix),
	}

	for _, dbName := range dbNames {
		if _, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)); err != nil {
			return err
		}
	}

	return nil
}

func GetCfg() (*parabase.Config, string) {
	managementDatabase, prefix := randomManagementDatabase()
	return &parabase.Config{
		DegreeOfParallelism: 5,
		DriverName:          "mysql",
		DataSourceName:      "root:password@tcp(db:3306)/",
		Timeout:             time.Second * 3,
		ManagementDatabase:  managementDatabase,
	}, prefix
}

func randomManagementDatabase() (*mysql.ManagementDatabase, string) {
	out := deepCopyManagementDatabase(mysql.DefaultManagentDB)

	prefix := randomName()
	out.Name = prefix + "_management"
	out.TestDBBaseName = prefix + "_test"

	return out, prefix
}

func deepCopyManagementDatabase(d *mysql.ManagementDatabase) *mysql.ManagementDatabase {
	return &mysql.ManagementDatabase{
		Name: d.Name,
		Tbl: &mysql.ManagementTable{
			Name: d.Tbl.Name,
			ID: &mysql.ManagementField{
				Name:       d.Tbl.ID.Name,
				Type:       d.Tbl.ID.Type,
				Constraint: d.Tbl.ID.Constraint,
				Zero:       d.Tbl.ID.Zero,
			},
			InUse: &mysql.ManagementField{
				Name:       d.Tbl.InUse.Name,
				Type:       d.Tbl.InUse.Type,
				Constraint: d.Tbl.InUse.Constraint,
				Zero:       d.Tbl.InUse.Zero,
			},
			ContextToken: &mysql.ManagementField{
				Name:       d.Tbl.ContextToken.Name,
				Type:       d.Tbl.ContextToken.Type,
				Constraint: d.Tbl.ContextToken.Constraint,
				Zero:       d.Tbl.ContextToken.Zero,
			},
		},
		TestDBBaseName: d.TestDBBaseName,
	}
}

const (
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	databaseNameLen = 20
)

func randomName() string {
	buf := make([]byte, databaseNameLen)
	max := new(big.Int)
	max.SetInt64(int64(len(letterBytes)))

	for i := range buf {
		r, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}

		buf[i] = letterBytes[r.Int64()]
	}

	return string(buf)
}
