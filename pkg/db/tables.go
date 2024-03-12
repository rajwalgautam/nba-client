package db

import (
	"database/sql"
)

var (
	// api key
	createApiKeySql = `create table if not exists apikey (
			key uuid not null,
			PRIMARY KEY (key)
		);`
	setApiKeySql = `
		insert into apikey values ( $1 )
		ON CONFLICT (key) DO NOTHING;`
	getApiKeySql = `
		select key from apikey;`

	// single game table
	createGameSql = `create table if not exists games (
			id int not null,
			date text not null,
			data JSONB,
			PRIMARY KEY(id)
		);`
	setGameSql = `
			insert into games values ( $1, $2, $3 )
			ON CONFLICT (id) DO UPDATE SET (data, date) =
			(excluded.data, excluded.date);`
	getAllGameIdsSql = `select id from games`

	createPlayerSql = `create table if not exists players (
				id int not null,
				name text not null,
				pos text not null,
				team text not null,
				games JSONB,
				PRIMARY KEY(id)
		);`
	setPlayerStatsSql = `
		insert into players (id, name, pos, team, games) values ( $1, $2, $3, $4, $5 )
			ON CONFLICT (id) DO UPDATE SET (name, pos, team, games) =
			(excluded.name, excluded.pos, excluded.team, excluded.games) ;`
	getPlayerSql = `
		select id, name, pos, team, games from players where id=$1;`
	getPlayerGamesSql = `
		select games from players where id=$1`
	getAllPlayerIdsSql = `select id from players`

	tables = []string{createApiKeySql, createGameSql, createPlayerSql}
)

func setupTables(db *sql.DB) error {
	for _, t := range tables {
		_, err := db.Exec(t)
		if err != nil {
			return err
		}
	}
	return nil
}
