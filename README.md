# nba-client
A client for NBA stats.

## About
`nba-client` is a CLI tool for processing NBA game, player and statistical data.

## Getting Started

### Layout
```
├── .gitignore
├── README.md
├── cmd                     // entrypoints for cli
├── pkg
│   ├── balldontlie         // api pkg for balldontlie
│   ├── db                  // postgres client
│   └── controller          // processor for handling menu requests
└── scripts
    └── db                  // docker-compose for postgres and adminer
```

### Commands

* `api` - pings balldontlie api 
* `initdb` - initializes db connection and tables
* `games` - fetches game data based on flags and saves to db 
    * Flag - `date` - in `YYYY-MM-DD` format
    * Flag - `start`/`end` - specify a date range, in `YYYY-MM-DD` format
* `syncStats` - fetch player stats for games in db, save to db
* `stats` root command for the following subcommands:
    * `nba60` - returns a list of all players that hit the NBA 60 criteria (sum of points, reb, ast, blk, stl > 60)
