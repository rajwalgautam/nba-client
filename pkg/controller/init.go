package controller

import (
	"github.com/rajwalgautam/nba-client/pkg/balldontlie"
	"github.com/rajwalgautam/nba-client/pkg/db"
)

type Processor struct {
	Datastore db.DB
	Api       balldontlie.Client
}

func New() (*Processor, error) {
	d, err := db.New()
	if err != nil {
		return nil, err
	}
	return NewWithDB(d)
}

func NewWithDB(d db.DB) (*Processor, error) {
	k, err := d.ApiKey()
	if err != nil {
		return nil, err
	}
	bdl := balldontlie.New(k)
	if err = bdl.Ping(); err != nil {
		return nil, err
	}
	return &Processor{
		Datastore: d,
		Api:       bdl,
	}, nil
}
