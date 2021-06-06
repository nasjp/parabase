package mysql

import (
	"fmt"
)

func (d *ManagementDatabase) defineManagementDB() string {
	return fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", d.Name)
}

func (d *ManagementDatabase) defineManagementTable() string {
	def := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s ( ", d.Name, d.Tbl.Name)
	for _, field := range []*ManagementField{d.Tbl.ID, d.Tbl.InUse, d.Tbl.ContextToken} {
		def += fmt.Sprintf("%s %s %s ,", field.Name, field.Type, field.Constraint)
	}

	def += fmt.Sprintf("PRIMARY KEY ( %s ) )", d.Tbl.ID.Name)

	return def
}

func (d *ManagementDatabase) setupManagementTable(num int) string {
	if num == 0 {
		num = 1
	}

	var unionQuery string

	for id := 1; id <= num; id++ {
		unionQuery += fmt.Sprintf(`
SELECT %d %s, %s %s, %s %s `,
			id, d.Tbl.ID.Name, d.Tbl.InUse.Zero, d.Tbl.InUse.Name, d.Tbl.ContextToken.Zero, d.Tbl.ContextToken.Name)

		if id == num {
			unionQuery += "\n"
			break
		}

		unionQuery += "UNION "
	}

	query := fmt.Sprintf(`
INSERT INTO %s.%s ( %s, %s, %s )
SELECT %s, %s, %s FROM (%s) tmp WHERE NOT EXISTS (SELECT 1 FROM %s.%s)`,
		d.Name, d.Tbl.Name, d.Tbl.ID.Name, d.Tbl.InUse.Name, d.Tbl.ContextToken.Name,
		d.Tbl.ID.Name, d.Tbl.InUse.Name, d.Tbl.ContextToken.Name, unionQuery, d.Name, d.Tbl.Name)

	return query
}

func (d *ManagementDatabase) testDBName(id int) string {
	return d.TestDBBaseName + fmt.Sprintf("_%d", id)
}

func (d *ManagementDatabase) defineTestDB(id int) string {
	return fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", d.testDBName(id))
}

func (d *ManagementDatabase) allocateTestDB(token string) string {
	return fmt.Sprintf(`
UPDATE %s.%s SET %s = '%s'
WHERE %s = ( SELECT %s FROM (
  (SELECT %s FROM %s.%s WHERE %s = %s ORDER BY RAND() LIMIT 1)
  UNION ALL
  (SELECT %s %s FROM DUAL WHERE NOT EXISTS ( SELECT 1 FROM %s.%s WHERE %s = %s LIMIT 1 ))
) tmp )`,
		d.Name, d.Tbl.Name, d.Tbl.ContextToken.Name, token,
		d.Tbl.ID.Name, d.Tbl.ID.Name,
		d.Tbl.ID.Name, d.Name, d.Tbl.Name, d.Tbl.ContextToken.Name, d.Tbl.ContextToken.Zero,
		d.Tbl.ID.Zero, d.Tbl.ID.Name, d.Name, d.Tbl.Name, d.Tbl.ContextToken.Name, d.Tbl.ContextToken.Zero,
	)
}

func (d *ManagementDatabase) checkAllocatedTestDB(token string) string {
	return fmt.Sprintf(
		"SELECT %s FROM %s.%s WHERE %s = '%s'",
		d.Tbl.ID.Name, d.Name, d.Tbl.Name, d.Tbl.ContextToken.Name, token,
	)
}

func (d *ManagementDatabase) freeTestDB(id int) string {
	return fmt.Sprintf(
		"UPDATE %s.%s SET %s = %s WHERE %s = %d",
		d.Name, d.Tbl.Name, d.Tbl.ContextToken.Name, d.Tbl.ContextToken.Zero, d.Tbl.ID.Name, id,
	)
}
