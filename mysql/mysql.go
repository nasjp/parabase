package mysql

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/nasjp/parabase"
)

var DefaultManagentDB = &ManagementDatabase{
	Name: "parabase_management",
	Tbl: &ManagementTable{
		Name: "management",
		ID: &ManagementField{
			Name:       "id",
			Type:       "INT",
			Constraint: "NOT NULL",
			Zero:       "0",
		},
		InUse: &ManagementField{
			Name:       "in_use",
			Type:       "TINYINT(1)",
			Constraint: "NOT NULL",
			Zero:       "0",
		},
		ContextToken: &ManagementField{
			Name:       "context_token",
			Type:       "VARCHAR(255)",
			Constraint: "NOT NULL",
			Zero:       "''",
		},
	},
	TestDBBaseName: "test_db",
}

type ManagementDatabase struct {
	Name           string
	Tbl            *ManagementTable
	TestDBBaseName string
}

type ManagementTable struct {
	Name         string
	ID           *ManagementField
	InUse        *ManagementField
	ContextToken *ManagementField
}

type ManagementField struct {
	Name       string
	Type       string
	Constraint string
	Zero       string
}

func (d *ManagementDatabase) Connect(dbName string, cfg *parabase.Config) (*sql.DB, error) {
	dns, err := dns(cfg.DataSourceName, dbName)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(cfg.DriverName, dns)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (d *ManagementDatabase) Setup(db *sql.DB, cfg *parabase.Config) error {
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, d.defineManagementDB()); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, d.defineManagementTable()); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, d.setupManagementTable(cfg.DegreeOfParallelism)); err != nil {
		return err
	}

	for id := 1; id <= cfg.DegreeOfParallelism; id++ {
		if _, err := db.ExecContext(ctx, d.defineTestDB(id)); err != nil {
			return err
		}
	}

	return nil
}

func tryFreeDB(db *sql.DB, d *ManagementDatabase, id int, cnt int) error {
	if cnt == 5 {
		return errors.New("can't exec %s by Deadlock")
	}
	if _, err := db.ExecContext(context.Background(), d.freeTestDB(id)); err != nil {
		sqlErr := &mysql.MySQLError{}
		if errors.As(err, &sqlErr) {
			// Error 1213: Deadlock found when trying to get lock; try restarting transaction
			if sqlErr.Number == 1213 {
				tryFreeDB(db, d, id, cnt+1)
			}
		}
	}

	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

func (d *ManagementDatabase) Get(db *sql.DB, cfg *parabase.Config) (*sql.DB, func(testing.TB), error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	getTeardownFunc := func(id int, allocatedDB *sql.DB) func(testing.TB) {
		return func(t testing.TB) {
			if err := tryFreeDB(db, d, id, 1); err != nil {
				t.Fatal(err)
			}
		}
	}

	token := uuid.New().String()

	for {
		select {
		case <-ctx.Done():
			return nil, nil, errors.New("allocated deadline exceeded")
		default:
			result, err := db.ExecContext(ctx, d.allocateTestDB(token))
			if err != nil {
				return nil, nil, err
			}

			if _, err := result.RowsAffected(); err != nil {
				return nil, nil, err
			}

			var id int
			err = db.QueryRowContext(ctx, d.checkAllocatedTestDB(token)).Scan(&id)
			if errors.Is(err, sql.ErrNoRows) {
				time.Sleep(1000 * time.Millisecond)
				continue
			}

			if err != nil {
				return nil, nil, err
			}

			allocatedDB, err := d.Connect(d.testDBName(id), cfg)
			if err != nil {
				return nil, nil, err
			}

			return allocatedDB, getTeardownFunc(id, allocatedDB), nil
		}
	}
}
