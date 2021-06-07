package mysql_test

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nasjp/parabase"
	"github.com/nasjp/parabase/mysql/internal/testutil"
)

func TestSetupAndTeardown(t *testing.T) {
	t.Parallel()
	cfg, prefix := testutil.GetCfg()
	if err := cfg.ManagementDatabase.Setup(testutil.DB, cfg); err != nil {
		t.Fatal(err)
	}

	if err := testutil.Cleanup(testutil.DB, prefix); err != nil {
		t.Fatal(err)
	}
}

func TestUse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		process func(t *testing.T) string
	}{
		{
			"One",
			func(t *testing.T) string {
				cfg, prefix := testutil.GetCfg()
				testTeardown := testUse(t, cfg, prefix)

				testTeardown(t)

				return prefix
			},
		},
		{
			"Two",
			func(t *testing.T) string {
				cfg, prefix := testutil.GetCfg()
				testTeardown1 := testUse(t, cfg, prefix)
				testTeardown2 := testUse(t, cfg, prefix)

				testTeardown2(t)
				testTeardown1(t)

				return prefix
			},
		},
		{
			"Three",
			func(t *testing.T) string {
				cfg, prefix := testutil.GetCfg()
				testTeardown1 := testUse(t, cfg, prefix)
				testTeardown2 := testUse(t, cfg, prefix)
				testTeardown3 := testUse(t, cfg, prefix)

				testTeardown3(t)
				testTeardown2(t)
				testTeardown1(t)

				return prefix
			},
		},
		{
			"Four",
			func(t *testing.T) string {
				cfg, prefix := testutil.GetCfg()
				testTeardown1 := testUse(t, cfg, prefix)
				testTeardown2 := testUse(t, cfg, prefix)
				testTeardown3 := testUse(t, cfg, prefix)
				testTeardown4 := testUse(t, cfg, prefix)

				testTeardown4(t)
				testTeardown3(t)
				testTeardown2(t)
				testTeardown1(t)

				return prefix
			},
		},
		{
			"Five",
			func(t *testing.T) string {
				cfg, prefix := testutil.GetCfg()
				testTeardown1 := testUse(t, cfg, prefix)
				testTeardown2 := testUse(t, cfg, prefix)
				testTeardown3 := testUse(t, cfg, prefix)
				testTeardown4 := testUse(t, cfg, prefix)
				testTeardown5 := testUse(t, cfg, prefix)

				testTeardown5(t)
				testTeardown4(t)
				testTeardown3(t)
				testTeardown2(t)
				testTeardown1(t)

				return prefix
			},
		},
		{
			"SixWaitFree",
			func(t *testing.T) string {
				cfg, prefix := testutil.GetCfg()
				testTeardown1 := testUse(t, cfg, prefix)
				testTeardown2 := testUse(t, cfg, prefix)
				testTeardown3 := testUse(t, cfg, prefix)
				testTeardown4 := testUse(t, cfg, prefix)
				testTeardown5 := testUse(t, cfg, prefix)

				var wg sync.WaitGroup

				wg.Add(1)

				go func() {
					defer wg.Done()
					testTeardown6 := testUse(t, cfg, prefix)

					testTeardown6(t)
				}()

				time.Sleep(time.Second * 2)

				testTeardown5(t)
				testTeardown4(t)
				testTeardown3(t)
				testTeardown2(t)
				testTeardown1(t)

				wg.Wait()

				return prefix
			},
		},
		{
			"SevenWaitFree",
			func(t *testing.T) string {
				cfg, prefix := testutil.GetCfg()
				testTeardown1 := testUse(t, cfg, prefix)
				testTeardown2 := testUse(t, cfg, prefix)
				testTeardown3 := testUse(t, cfg, prefix)
				testTeardown4 := testUse(t, cfg, prefix)
				testTeardown5 := testUse(t, cfg, prefix)

				var wg sync.WaitGroup

				wg.Add(1)

				go func() {
					defer wg.Done()
					testTeardown6 := testUse(t, cfg, prefix)

					testTeardown6(t)
				}()

				wg.Add(1)

				go func() {
					defer wg.Done()
					testTeardown7 := testUse(t, cfg, prefix)

					testTeardown7(t)
				}()

				time.Sleep(time.Second * 2)

				testTeardown5(t)
				testTeardown4(t)
				testTeardown3(t)
				testTeardown2(t)
				testTeardown1(t)

				wg.Wait()

				return prefix
			},
		},
		{
			"EightWaitFree",
			func(t *testing.T) string {
				cfg, prefix := testutil.GetCfg()
				testTeardown1 := testUse(t, cfg, prefix)
				testTeardown2 := testUse(t, cfg, prefix)
				testTeardown3 := testUse(t, cfg, prefix)
				testTeardown4 := testUse(t, cfg, prefix)
				testTeardown5 := testUse(t, cfg, prefix)

				var wg sync.WaitGroup

				wg.Add(1)

				go func() {
					defer wg.Done()
					testTeardown6 := testUse(t, cfg, prefix)

					testTeardown6(t)
				}()

				wg.Add(1)

				go func() {
					defer wg.Done()
					testTeardown7 := testUse(t, cfg, prefix)

					testTeardown7(t)
				}()

				wg.Add(1)

				go func() {
					defer wg.Done()
					testTeardown8 := testUse(t, cfg, prefix)

					testTeardown8(t)
				}()

				time.Sleep(time.Second * 2)

				testTeardown5(t)
				testTeardown4(t)
				testTeardown3(t)
				testTeardown2(t)
				testTeardown1(t)

				wg.Wait()

				return prefix
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			prefix := tt.process(t)

			t.Cleanup(func() {
				if err := testutil.Cleanup(testutil.DB, prefix); err != nil {
					t.Fatal(err)
				}
			})
		})
	}
}

func testUse(t *testing.T, cfg *parabase.Config, prefix string) func(t testing.TB) {
	allocatedDB, teardown, err := parabase.Use(cfg)
	if err != nil {
		t.Fatal(err)
	}

	dbName := &sql.NullString{}
	if err := allocatedDB.QueryRow("SELECT DATABASE()").Scan(dbName); err != nil {
		t.Fatal(err)
	}

	if !dbName.Valid {
		t.Errorf("database is not selected")
		return nil
	}

	id, err := strconv.Atoi(dbName.String[strings.LastIndex(dbName.String, "_")+1:])
	if err != nil {
		t.Fatal(err)
	}

	managementDBName := fmt.Sprintf("%s_management", dbName.String[:strings.Index(dbName.String, "_")])

	var contextToken string

	if err := testutil.DB.QueryRow(fmt.Sprintf("SELECT context_token FROM %s.%s WHERE id = ?", managementDBName, "management"), id).Scan(&contextToken); err != nil {
		t.Fatal(err)
	}

	if contextToken == "" {
		t.Error("management table's context_token is empty")
		return nil
	}

	testTeardown := func(t testing.TB) {
		teardown(t)

		rows, err := testutil.DB.Query(fmt.Sprintf("SELECT 1 FROM %s.%s WHERE context_token = ?", managementDBName, "management"), contextToken)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()

		var unfree bool

		for rows.Next() {
			unfree = true
		}

		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}

		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}

		if unfree {
			t.Error("management database is not free")
			return
		}
	}

	return testTeardown
}
