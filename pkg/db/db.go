package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-multierror"
	_ "github.com/lib/pq"
	"github.com/rajwalgautam/nba-client/pkg/balldontlie"
)

type DB interface {
	BulkSetStats(stats []balldontlie.Stats) error
	SetPlayerStats(s balldontlie.Stats) error
	SetGame(g balldontlie.Game) error
	SetApiKey(k string) error
	ApiKey() (string, error)
}

type nbaDB struct {
	pgdb *sql.DB
	ppm  perPlayerMutex
}

func New() (DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "example", "nba-client", "nba-stats")
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	err = setupTables(db)
	return &nbaDB{pgdb: db, ppm: make(perPlayerMutex)}, err
}

func (n *nbaDB) BulkSetStats(stats []balldontlie.Stats) error {
	var me error
	for _, s := range stats {
		err := n.SetPlayerStats(s)
		if err != nil {
			me = multierror.Append(me, err)
		}
	}
	return me
}

func (n *nbaDB) SetPlayerStats(s balldontlie.Stats) error {
	n.ppm.Lock(s.Player.Id)
	defer n.ppm.Unlock(s.Player.Id)
	gs, err := n.gameStatsByPlayerId(s.Player.Id)
	if err != nil {
		return err
	}
	for _, g := range gs {
		if g.Id == s.Id {
			return nil
		}
	}
	gs = append(gs, balldontlie.PlayerStats{
		Id:       s.Id,
		Points:   s.Points,
		Rebounds: s.Rebounds,
		Assists:  s.Assists,
		Steals:   s.Steals,
		Blocks:   s.Blocks,
		Date:     s.Game.Date,
	})
	b, _ := json.Marshal(gs)
	_, err = n.pgdb.Exec(setPlayerStatsSql, s.Player.Id, fmt.Sprintf("%s %s", s.Player.First, s.Player.Last), s.Player.Position, s.Team.Abbr, b)
	return err
}

func (n *nbaDB) gameStatsByPlayerId(id int) ([]balldontlie.PlayerStats, error) {
	rows, err := n.pgdb.Query(getPlayerGamesSql, id)
	if err != nil {
		return nil, nil
	}
	var stats []balldontlie.PlayerStats
	for rows.Next() {
		var b []byte
		err = rows.Scan(&b)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &stats)
		if err != nil {
			return nil, err
		}
	}
	return stats, nil
}

func (n *nbaDB) playerExists(id int) (bool, error) {
	rows, err := n.pgdb.Query(getPlayerSql, id)
	if err != nil {
		return false, err
	}
	for rows.Next() {
		p := new(balldontlie.Player)
		err = rows.Scan(&p)
		if p != nil {
			return true, nil
		}
	}
	return false, nil
}

func (n *nbaDB) SetGame(g balldontlie.Game) error {
	b, _ := json.Marshal(g)
	_, err := n.pgdb.Exec(setGameSql, g.Id, g.Date, b)
	return err
}

func (n *nbaDB) SetApiKey(k string) error {
	_, err := n.pgdb.Exec(setApiKeySql, k)
	return err
}

func (n *nbaDB) ApiKey() (string, error) {
	rows, err := n.pgdb.Query(getApiKeySql)
	if err != nil {
		return "", err
	}
	var key uuid.UUID
	for rows.Next() {
		err := rows.Scan(&key)
		if err != nil {
			return "", err
		}
		if key.String() != "" {
			break
		}
	}
	return key.String(), nil
}

type perPlayerMutex map[int]*sync.Mutex

var ppmu = &sync.Mutex{}

func (ppm perPlayerMutex) Unlock(id int) {
	ppmu.Lock()
	defer ppmu.Unlock()
	_, exists := ppm[id]
	if exists {
		ppm[id].Unlock()
		delete(ppm, id)
	}
}

func (ppm perPlayerMutex) Lock(id int) {
	ppmu.Lock()
	defer ppmu.Unlock()
	_, exists := ppm[id]
	if !exists {
		ppm[id] = &sync.Mutex{}
		ppm[id].Lock()
	}
}
