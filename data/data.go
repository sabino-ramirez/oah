package data

import (
	"database/sql"
	"fmt"

	// "log"

	_ "github.com/mattn/go-sqlite3"
	// "github.com/sabino-ramirez/oah/models"
)

var db *sql.DB
var err error

func InitDB(dburl string) error {
	db, err = sql.Open("sqlite3", dburl)
	if err != nil {
		return err
	}
	// log.Println("DB created")
	return db.Ping()
}

func CreateTable() error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS params(tryId INTEGER NOT NULL PRIMARY KEY CHECK (tryId = 1), auth TEXT, orgId INT, projTempId INT );`

	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		return fmt.Errorf("error preparing create table statement: %v", err)
	}

	statement.Exec()
	// log.Println("Created params table")

	insertDefaultSQL := `REPLACE INTO params (tryId, auth, orgId, projTempId) VALUES (1, 1, 1, 1)`
	statement, err = db.Prepare(insertDefaultSQL)
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %v", err)
	}

	statement.Exec()
	// log.Println("default insert successful")

	return nil
}

func UpdateX(key string, value any) error {
	updateSQL := `UPDATE params SET ` + key + ` = ? WHERE tryId = 1`
	statement, err := db.Prepare(updateSQL)
	if err != nil {
		return fmt.Errorf("error preparing update statement: %v", err)
	}

	statement.Exec(value)
	// log.Printf("%v update successful", key)

	return nil
}

// func GetXValue() (models.DbRow, error) {
// 	selectSQL := `SELECT auth, orgId, projTempId FROM params WHERE tryId = 1;`
//
// 	row := db.QueryRow(selectSQL)
// 	params := models.DbRow{}
// 	if err = row.Scan(&params.Auth, &params.OrgId, &params.ProjTempId); err != nil {
// 		return models.DbRow{}, fmt.Errorf("not found: %v", err)
// 	}
//
// 	return params, nil
// }
