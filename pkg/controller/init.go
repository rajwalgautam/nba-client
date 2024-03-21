package controller

import (
	"github.com/rajwalgautam/nba-client/pkg/balldontlie"
	"github.com/rajwalgautam/nba-client/pkg/db"
)

var (
	statsQueue = make(chan balldontlie.Stats)
	gameQueue  = make(chan balldontlie.Game)

	datastore db.DB
	api       balldontlie.Client
)

type Controller struct {
	Datastore db.DB
	Api       balldontlie.Client
}

func New() (*Controller, error) {
	d, err := db.New()
	if err != nil {
		return nil, err
	}
	return NewWithDB(d)
}

func NewWithDB(d db.DB) (*Controller, error) {
	k, err := d.ApiKey()
	if err != nil {
		return nil, err
	}
	bdl := balldontlie.New(k)
	if err = bdl.Ping(); err != nil {
		return nil, err
	}
	return &Controller{
		Datastore: d,
		Api:       bdl,
	}, nil
}
