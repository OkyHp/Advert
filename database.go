package main

import (
	"encoding/json"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type AdvertsData struct {
	MsgType string            `db:"msg_type"`
	MsgText map[string]string `db:"msg_text"`
}

type Database struct {
	conn *sqlx.DB
}

func LoadAdvert() error {
	MSGDebug("Advert InitDatabase")

	var db Database
	var err error

	db.conn, err = createDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.conn.Close()

	_, err = db.conn.Exec(`SET search_path TO ` + Plugin.Config.Database.Schema + `;`)
	if err != nil {
		return fmt.Errorf("Set schema: %w", err)
	}

	if !Plugin.DatabaseInit {
		err = db.createTable()
		if err != nil {
			return err
		}

		Plugin.DatabaseInit = true
	}

	err = db.getAdverts()
	if err != nil {
		return err
	}

	return nil
}

func createDatabaseConnection() (*sqlx.DB, error) {
	// Build PostgreSQL connection string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Plugin.Config.Database.Host,
		Plugin.Config.Database.Port,
		Plugin.Config.Database.User,
		Plugin.Config.Database.Pass,
		Plugin.Config.Database.Base,
	)

	MSGDebug("Advert DB: host=%s port=%d user=%s dbname=%s",
		Plugin.Config.Database.Host,
		Plugin.Config.Database.Port,
		Plugin.Config.Database.User,
		Plugin.Config.Database.Base)

	dbConn, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	MSGDebug("Advert createDatabaseConnection")

	return dbConn, nil
}

func (db *Database) createTable() error {
	query := `CREATE TABLE IF NOT EXISTS adverts (
			id SERIAL PRIMARY KEY,
			servers VARCHAR(64) NOT NULL DEFAULT '[]',
			msg_type VARCHAR(16) NOT NULL CHECK (msg_type IN ('CHAT', 'CENTER', 'ALERT', 'HTML')),
			msg_text VARCHAR(8192) NOT NULL DEFAULT '{
"en": "",
"uk": "",
"ru": ""
}',
			disable SMALLINT NOT NULL DEFAULT 0,
			position INTEGER NOT NULL DEFAULT 0
		);`

	_, err := db.conn.Exec(query)
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	MSGDebug("Advert createTable")

	return nil
}

func (db *Database) getAdverts() error {
	query := `
		SELECT msg_type, msg_text
		FROM adverts
		WHERE (
		    servers = '' 
		    OR servers = '[]'
		    OR servers::jsonb @> jsonb_build_array($1::integer)
		) AND disable = 0
		ORDER BY id, position DESC
	`

	rows, err := db.conn.Query(query, Plugin.Config.ServerId)
	if err != nil {
		return fmt.Errorf("select 'adverts': %w", err)
	}
	defer rows.Close()

	Plugin.Adverts = []AdvertsData{}
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

	MSGDebug("Advert getAdverts")

	return nil
}
