package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type AdvertsData struct {
	MsgType string            `db:"msg_type"`
	MsgText map[string]string `db:"msg_text"`
}

var db *sql.DB

func InitDatabase() error {
	var err error

	db, err = createConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	err = CreateTable()
	if err != nil {
		return err
	}

	err = GetAdverts()
	if err != nil {
		return err
	}

	return nil
}

func createConnection() (*sql.DB, error) {
	loginURLValues := url.Values{}

	// default params
	loginURLValues.Add("loc", time.Now().Location().String())
	loginURLValues.Add("parseTime", "true")
	loginURLValues.Add("charset", "utf8mb4")

	// override params
	for k, v := range Plugin.Config.Database.Params {
		loginURLValues.Add(k, v)
	}

	dsn := url.URL{
		Host:     fmt.Sprintf("tcp(%s:%d)", Plugin.Config.Database.Host, Plugin.Config.Database.Port),
		User:     url.UserPassword(Plugin.Config.Database.User, Plugin.Config.Database.Pass),
		Path:     Plugin.Config.Database.Base,
		RawQuery: loginURLValues.Encode(),
	}

	db, err := sql.Open("mysql", dsn.String()[2:])
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}

func CreateTable() error {
	queries := []string{ // For N+ queries
		`CREATE TABLE IF NOT EXISTS adverts(
			id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			servers VARCHAR(64) NOT NULL DEFAULT '[]',
			msg_type ENUM('CHAT', 'CENTER', 'ALERT', 'HTML') NOT NULL,
			msg_text VARCHAR(4096) NOT NULL DEFAULT '{\r\n\"en\": \"\",\r\n\"uk\": \"\",\r\n\"ru\": \"\"\r\n}',
			disable TINYINT UNSIGNED NOT NULL DEFAULT '0',
			position INT NOT NULL DEFAULT '0',
			PRIMARY KEY(id)
		) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	}

	for index := range queries {
		_, err := db.Exec(queries[index])
		if err != nil {
			return fmt.Errorf("create table 'adverts' (#%d): %w", index, err)
		}
	}

	return nil
}

func GetAdverts() error {
	query := `
		SELECT msg_type, msg_text
		FROM adverts
		WHERE (
		    JSON_CONTAINS(servers, ?, '$') 
		    OR servers = '' 
		    OR servers = '[]'
		) AND disable = 0
		ORDER BY id, position DESC
	`

	rows, err := db.Query(query, Plugin.Config.ServerId)
	if err != nil {
		return fmt.Errorf("select 'adverts': %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item AdvertsData
		var msgTextRaw string

		err = rows.Scan(&item.MsgType, &msgTextRaw)
		if err != nil {
			return fmt.Errorf("scan 'adverts': %w", err)
		}

		err = json.Unmarshal([]byte(msgTextRaw), &item.MsgText)
		if err != nil {
			return fmt.Errorf("unmarshal msg_text: %w (value: %s)", err, msgTextRaw)
		}

		ReplaceStaticPlaceholders(&item)

		Plugin.Adverts = append(Plugin.Adverts, item)
	}

	err = rows.Err()
	if err != nil {
		return fmt.Errorf("rows 'adverts': %w", err)
	}

	return nil
}
