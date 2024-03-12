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
│   └── controller      // processor for handling 
└── scripts
    ├── db
    └── seed-db
```