package parser

import (
	"encoding/xml"
	"strings"
)

type Replay struct {
	ClientVersion string       `xml:"ClientVersion"`
	ReplaySteps   []ReplayStep `xml:"ReplayStep"`
}

type ReplayStep struct {
	GameInfos                GameInfos `xml:"GameInfos"`
	RulesEventWaitingRequest interface{}
	BoardState               BoardState
	RulesEventGameFinished   RulesEventGameFinished
}

type GameInfos struct {
	ID           string         `xml:"Id"`
	CoachesInfos CoachesInfos   `xml:"CoachesInfos"`
	StadiumName  string         `xml:"NameStadium"`
	StadiumLevel int            `xml:"LevelStadium"`
	League       RowLeague      `xml:"RowLeague"`
	Competition  RowCompetition `xml:"RowCompetition"`
}

type CoachesInfos struct {
	CoachInfos []CoachInfos `xml:"CoachInfos"`
}

type CoachInfos struct {
	Login  string `xml:"Login"`
	UserID string `xml:"UserId"`
}

type RowLeague struct {
	Name            string `xml:"Name"`
	RegisteredTeams int    `xml:"NbRegisteredTeams"`
}

type RowCompetition struct {
	Name  string `xml:"Name"`
	Round int    `xml:"CurrentRound"`
}

type BoardState struct {
	ListTeams ListTeams `xml:"ListTeams"`
	Weather   int       `xml:"Meteo"`
}

type ListTeams struct {
	TeamState []Team `xml:"TeamState"`
}

type Team struct {
	RerollNumber        int
	TeamRerollAvailable int
	Name                string `xml:"Data>Name"`
	Value               int    `xml:"Data>Value"`
	Rerolls             int    `xml:"Data>Reroll"`
	Race                Race   `xml:"Data>IdRace"`
	ListPitchPlayers    PitchPlayers
}

type RulesEventGameFinished struct {
	Coaches []CoachResult `xml:"MatchResult>CoachResults>CoachResult"`
}

type MatchResult struct {
	CoachResults CoachResults
}

type CoachResults struct {
	Coaches []CoachResult `xml:"CoachResult"`
}

type CoachResult struct {
	TeamResult TeamResult
}

type TeamResult struct {
	PopularityBeforeMatch      int
	TeamValue                  int    `xml:"TeamData>Value"`
	Name                       string `xml:"TeamData>Name"`
	Cheerleaders               string `xml:"TeamData>Cheerleaders"`
	Supporters                 int    `xml:"NoSupporters"`
	CashBeforeMatch            int    `xml:"CashBeforeMatch"`
	CashEarnedBeforeConcession int
	WinningsDice               int
	CashEarned                 int
	Popularity                 string         `xml:"TeamData>Popularity"`
	PlayerResults              []PlayerResult `xml:"PlayerResults>PlayerResult"`
}

type PlayerResult struct {
}

type Race string

func (r *Race) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var rr string
	if err := d.DecodeElement(&rr, &start); err != nil {
		return err
	}

	if race, ok := RaceMapping[rr]; ok {
		rc := Race(race)
		*r = rc
	}
	return nil
}

type PitchPlayers struct {
	PlayerStates []PlayerState `xml:"PlayerState"`
}

type PlayerState struct {
	Name     string   `xml:"Data>Name"`
	Type     string   `xml:"Data>IdPlayerTypes"`
	Movement int      `xml:"Data>Ma"`
	Agility  int      `xml:"Data>Ag"`
	Armor    int      `xml:"Data>Av"`
	Strength int      `xml:"Data>St"`
	Skills   []string `xml:"Data>ListSkills"`
}

type rawData struct {
	Name     string `xml:"Data>Name"`
	Type     string `xml:"Data>IdPlayerTypes"`
	Movement int    `xml:"Data>Ma"`
	Agility  int    `xml:"Data>Ag"`
	Armor    int    `xml:"Data>Av"`
	Strength int    `xml:"Data>St"`
	Skills   string `xml:"Data>ListSkills"`
}

func (ps *PlayerState) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw rawData
	if err := d.DecodeElement(&raw, &start); err != nil {
		return err
	}

	ps.Name = raw.Name
	if tp, ok := PlayerTypesMapping[raw.Type]; ok {
		ps.Type = tp
	}
	ps.Movement = raw.Movement
	ps.Agility = raw.Agility
	ps.Armor = raw.Armor
	ps.Strength = raw.Strength

	if raw.Skills != "" {
		skills := strings.Split(strings.Trim(raw.Skills, "()"), ",")
		skillList := make([]string, 0)
		for _, skill := range skills {
			if skillName, ok := SkillMapping[skill]; ok {
				skillList = append(skillList, skillName)
			}
		}
		ps.Skills = skillList
	}

	return nil
}
