package gosqljson

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func QueryDbToArrayJson(db *sql.DB, theCase string, sqlStatement string, sqlParams ...interface{}) (string, error) {
	data, err := QueryDbToArray(db, theCase, sqlStatement, sqlParams...)
	jsonString, err := json.Marshal(data)
	return string(jsonString), err
}

func QueryDbToMapJson(db *sql.DB, theCase string, sqlStatement string, sqlParams ...interface{}) (string, error) {
	data, err := QueryDbToMap(db, theCase, sqlStatement, sqlParams...)
	jsonString, err := json.Marshal(data)
	return string(jsonString), err
}

func QueryDbToArray(db *sql.DB, theCase string, sqlStatement string, sqlParams ...interface{}) ([][]string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	var results [][]string
	if strings.HasPrefix(strings.ToUpper(sqlStatement), "SELECT") {
		rows, err := db.Query(sqlStatement, sqlParams...)
		if err != nil {
			fmt.Println("Error executing: ", sqlStatement)
			return results, err
		}
		cols, _ := rows.Columns()
		if theCase == "lower" {
			colsLower := make([]string, len(cols))
			for i, v := range cols {
				colsLower[i] = strings.ToLower(v)
			}
			results = append(results, colsLower)
		} else if theCase == "upper" {
			colsUpper := make([]string, len(cols))
			for i, v := range cols {
				colsUpper[i] = strings.ToUpper(v)
			}
			results = append(results, colsUpper)
		} else if theCase == "camel" {
			colsCamel := make([]string, len(cols))
			for i, v := range cols {
				colsCamel[i] = toCamel(v)
			}
			results = append(results, colsCamel)
		}

		rawResult := make([][]byte, len(cols))

		dest := make([]interface{}, len(cols)) // A temporary interface{} slice
		for i, _ := range rawResult {
			dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
		}

		for rows.Next() {
			result := make([]string, len(cols))
			rows.Scan(dest...)
			for i, raw := range rawResult {
				if raw == nil {
					result[i] = ""
				} else {
					result[i] = string(raw)
				}
			}
			results = append(results, result)
		}
	}
	return results, nil
}

func QueryDbToMap(db *sql.DB, theCase string, sqlStatement string, sqlParams ...interface{}) ([]map[string]string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	var results []map[string]string
	if strings.HasPrefix(strings.ToUpper(sqlStatement), "SELECT ") {
		rows, err := db.Query(sqlStatement, sqlParams...)
		if err != nil {
			fmt.Println("Error executing: ", sqlStatement)
			return results, err
		}
		cols, _ := rows.Columns()
		colsLower := make([]string, len(cols))
		colsCamel := make([]string, len(cols))

		if theCase == "lower" {
			for i, v := range cols {
				colsLower[i] = strings.ToLower(v)
			}
		} else if theCase == "upper" {
			for i, v := range cols {
				cols[i] = strings.ToUpper(v)
			}
		} else if theCase == "camel" {
			for i, v := range cols {
				colsCamel[i] = toCamel(v)
			}
		}

		rawResult := make([][]byte, len(cols))

		dest := make([]interface{}, len(cols)) // A temporary interface{} slice
		for i, _ := range rawResult {
			dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
		}

		for rows.Next() {
			result := make(map[string]string, len(cols))
			rows.Scan(dest...)
			for i, raw := range rawResult {
				if raw == nil {
					if theCase == "lower" {
						result[colsLower[i]] = ""
					} else if theCase == "upper" {
						result[cols[i]] = ""
					} else if theCase == "camel" {
						result[colsCamel[i]] = ""
					}
				} else {
					if theCase == "lower" {
						result[colsLower[i]] = string(raw)
					} else if theCase == "upper" {
						result[cols[i]] = string(raw)
					} else if theCase == "camel" {
						result[colsCamel[i]] = string(raw)
					}
				}
			}
			results = append(results, result)
		}
	}
	return results, nil
}

func ExecDb(db *sql.DB, sqlStatement string, sqlParams ...interface{}) (int64, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	sqlUpper := strings.ToUpper(sqlStatement)
	if strings.HasPrefix(sqlUpper, "UPDATE ") ||
		strings.HasPrefix(sqlUpper, "INSERT ") ||
		strings.HasPrefix(sqlUpper, "DELETE FROM ") {
		result, err := db.Exec(sqlStatement, sqlParams...)
		if err != nil {
			fmt.Println("Error executing: ", sqlStatement)
			fmt.Println(err)
			return 0, err
		}
		return result.RowsAffected()
	}
	return 0, errors.New(fmt.Sprint("Invalid SQL:", sqlStatement))
}

func ExecTx(tx *sql.Tx, sqlStatement string, sqlParams ...interface{}) (int64, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	sqlUpper := strings.ToUpper(sqlStatement)
	if strings.HasPrefix(sqlUpper, "UPDATE ") ||
		strings.HasPrefix(sqlUpper, "INSERT ") ||
		strings.HasPrefix(sqlUpper, "DELETE FROM ") {
		result, err := tx.Exec(sqlStatement, sqlParams...)
		if err != nil {
			fmt.Println("Error executing: ", sqlStatement)
			fmt.Println(err)
			return 0, err
		}
		return result.RowsAffected()
	}
	return 0, errors.New(fmt.Sprint("Invalid SQL:", sqlStatement))
}

func toCamel(s string) (ret string) {
	s = strings.ToLower(s)
	a := strings.Split(s, "_")
	for i, v := range a {
		if i == 0 {
			ret += v
		} else {
			f := strings.ToUpper(string(v[0]))
			n := string(v[1:])
			ret += fmt.Sprint(f, n)
		}
	}
	return
}
