package mysql_test

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"testing"

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

	t.Run("One", func(t *testing.T) {
		t.Parallel()

		cfg, prefix := testutil.GetCfg()
		testUse(t, cfg, prefix)

		t.Cleanup(func() {
			if err := testutil.Cleanup(testutil.DB, prefix); err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("Two", func(t *testing.T) {
		t.Parallel()

		cfg, prefix := testutil.GetCfg()
		testUse(t, cfg, prefix)
		testUse(t, cfg, prefix)

		t.Cleanup(func() {
			if err := testutil.Cleanup(testutil.DB, prefix); err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("Three", func(t *testing.T) {
		t.Parallel()

		cfg, prefix := testutil.GetCfg()
		testUse(t, cfg, prefix)
		testUse(t, cfg, prefix)

		t.Cleanup(func() {
			if err := testutil.Cleanup(testutil.DB, prefix); err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("Four", func(t *testing.T) {
		t.Parallel()

		cfg, prefix := testutil.GetCfg()
		testUse(t, cfg, prefix)
		testUse(t, cfg, prefix)

		t.Cleanup(func() {
			if err := testutil.Cleanup(testutil.DB, prefix); err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("Five", func(t *testing.T) {
		t.Parallel()

		cfg, prefix := testutil.GetCfg()
		testUse(t, cfg, prefix)
		testUse(t, cfg, prefix)
		testUse(t, cfg, prefix)
		testUse(t, cfg, prefix)
		testUse(t, cfg, prefix)

		t.Cleanup(func() {
			if err := testutil.Cleanup(testutil.DB, prefix); err != nil {
				t.Fatal(err)
			}
		})
	})
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
	if err := allocatedDB.QueryRow(fmt.Sprintf("SELECT context_token FROM %s.%s WHERE id = ?", managementDBName, "management"), id).Scan(&contextToken); err != nil {
		t.Fatal(err)
	}

	if contextToken == "" {
		t.Error("management table's context_token is empty")
		return nil
	}

	return teardown
}
