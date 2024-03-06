package balldontlie

import "fmt"

type Team struct {
	Id         int    `json:"id"`
	Conference string `json:"conference"`
	Division   string `json:"division"`
	FullName   string `json:"full_name"`
	Abbr       string `json:"abbreviation"`
}

type Games []Game

func (gs Games) Print() {
	for _, g := range gs {
		g.Print()
	}
}

type Game struct {
	Id        int    `json:"id"`
	HomeScore int    `json:"home_team_score"`
	AwayScore int    `json:"visitor_team_score"`
	HomeTeam  Team   `json:"home_team"`
	AwayTeam  Team   `json:"visitor_team"`
	Date      string `json:"date"`
}

func (g Game) Print() {
	fmt.Printf("%s (%s): %d\n%s (%s): %d\n\n",
		g.HomeTeam.FullName, g.HomeTeam.Abbr, g.HomeScore,
		g.AwayTeam.FullName, g.AwayTeam.Abbr, g.HomeScore)
}

type gameWrapper struct {
	Data Games `json:"data"`
}

type statsWrapper struct {
	Data []Stats `json:"data"`
}

type Stats struct {
	Id       int    `json:"id"`
	Points   int    `json:"pts"`
	Rebounds int    `json:"reb"`
	Assists  int    `json:"ast"`
	Steals   int    `json:"stl"`
	Blocks   int    `json:"blk"`
	Player   Player `json:"player"`
	Team     Team   `json:"team"`
	Game     Game   `json:"game"`
}

func (s Stats) Empty() bool {
	return s.Points == 0 && s.Rebounds == 0 && s.Assists == 0 && s.Steals == 0 && s.Blocks == 0
}

type Player struct {
	Id       int           `json:"id"`
	First    string        `json:"first_name"`
	Last     string        `json:"last_name"`
	Position string        `json:"position"`
	Team     string        `json:"team"`
	Games    []PlayerStats `json:"-"`
}

type PlayerStats struct {
	Id       int    `json:"id"`
	Points   int    `json:"pts"`
	Rebounds int    `json:"reb"`
	Assists  int    `json:"ast"`
	Steals   int    `json:"stl"`
	Blocks   int    `json:"blk"`
	Date     string `json:"date"`
}
