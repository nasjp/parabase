package mysql

import (
	"fmt"
	"strings"
)

func dns(dataSourceName string, dbName string) (string, error) {
	splitted := strings.SplitAfterN(dataSourceName, "/", 2)

	if len(splitted) != 2 {
		return "", fmt.Errorf("dns '%s' doesn't contain '/'", dataSourceName)
	}

	sp2 := splitBefore(splitted[len(splitted)-1], "?")

	return strings.Join([]string{splitted[0], dbName, sp2[1]}, ""), nil
}

func splitBefore(str string, sep string) []string {
	if i := strings.Index(str, sep); i >= 0 {
		return []string{str[:i], str[i:]}
	}

	return []string{str, ""}
}
