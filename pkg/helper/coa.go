package helper

import "fmt"

func CoaTableShard(year string) string {
	tableName := fmt.Sprintf("coas_%s", year)

	return tableName
}
