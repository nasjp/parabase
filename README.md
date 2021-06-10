[![ci](https://github.com/nasjp/parabase/actions/workflows/ci.yml/badge.svg)](https://github.com/nasjp/parabase/actions/workflows/ci.yml)

# parabase

```go
package foo_test

import (
	"testing"
	"time"

	"github.com/nasjp/parabase"
	"github.com/nasjp/parabase/mysql"
)

func TestFoo(t *testing.T) {
	db := parabase.Use(&parabase.Config{
		DegreeOfParallelism: 5,
		DriverName:          "mysql",
		DataSourceName:      "root:password@tcp(db:3306)/",
		Timeout:             time.Second * 3,
		ManagementDatabase:  mysql.DefaultManagentDB,
	})
}
```
