package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type Replay struct {
	ReplaySteps []ReplayStep `xml:"ReplayStep"`
}

type ReplayStep struct {
	RulesEventGameFinished RulesEventGameFinished
}

type RulesEventGameFinished struct {
	Coaches    []CoachResult `xml:"MatchResult>CoachResults>CoachResult"`
	Statistics Statistics    `xml:"MatchResult>Row"`
}

type CoachResult struct {
	TeamResult TeamResult
}

type TeamResult struct {
	PopularityBeforeMatch      int
	TeamValue                  int    `xml:"TeamData>Value"`
	Name                       string `xml:"TeamData>Name"`
	Cheerleaders               int    `xml:"TeamData>Cheerleaders"`
	Supporters                 int    `xml:"NoSupporters"`
	CashBeforeMatch            int    `xml:"CashBeforeMatch"`
	CashEarnedBeforeConcession int
	WinningsDice               int
	CashEarned                 int
	Popularity                 int            `xml:"TeamData>Popularity"`
	Race                       Race           `xml:"TeamData>IdRace"`
	PlayerResults              []PlayerResult `xml:"PlayerResults>PlayerResult"`
}

type PlayerResult struct {
	Name                string
	Type                string
	Movement            int
	Agility             int
	Armor               int
	Strength            int
	Skills              []string
	XP                  int
	InflictedTackles    int
	SustainedTackles    int
	InflictedInjuries   int
	SustainedInjuries   int
	InflictedCasualties int
	SustainedCasualties int
	MVP                 bool
	Casualties          []string
}

type rawPlayerData struct {
	Name                string `xml:"PlayerData>Name"`
	Type                string `xml:"PlayerData>IdPlayerTypes"`
	Movement            int    `xml:"PlayerData>Ma"`
	Agility             int    `xml:"PlayerData>Ag"`
	Armor               int    `xml:"PlayerData>Av"`
	Strength            int    `xml:"PlayerData>St"`
	Skills              string `xml:"PlayerData>ListSkills"`
	XP                  int
	InflictedTackles    int `xml:"Statistics>InflictedTackles"`
	SustainedTackles    int `xml:"Statistics>SustainedTackles"`
	InflictedInjuries   int `xml:"Statistics>InflictedInjuries"`
	SustainedInjuries   int `xml:"Statistics>SustainedInjuries"`
	InflictedCasualties int `xml:"Statistics>InflictedCasualties"`
	SustainedCasualties int `xml:"Statistics>SustainedCasualties"`
	MVP                 int `xml:"Statistics>MVP"`
	Casualty1           int
	Casualty2           int
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

func (ps *PlayerResult) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw rawPlayerData
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
	ps.XP = raw.XP
	ps.InflictedInjuries = raw.InflictedInjuries
	ps.SustainedInjuries = raw.SustainedInjuries
	ps.InflictedTackles = raw.InflictedTackles
	ps.SustainedTackles = raw.SustainedTackles
	ps.InflictedCasualties = raw.InflictedCasualties
	ps.SustainedCasualties = raw.SustainedCasualties

	ps.MVP = raw.MVP == 1

	ps.Casualties = make([]string, 0)
	if raw.Casualty1 != 0 {
		ps.Casualties = append(ps.Casualties, CasualtyMapping[fmt.Sprint(raw.Casualty1)])
	}
	if raw.Casualty2 != 0 {
		ps.Casualties = append(ps.Casualties, CasualtyMapping[fmt.Sprint(raw.Casualty2)])
	}

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

type Statistics struct {
	HomeInflictedInjuries          int
	AwayInflictedInjuries          int
	HomeSustainedKO                int
	AwaySustainedKO                int
	HomeOccupationOwn              int
	AwayOccupationOwn              int
	HomeValue                      int
	AwayValue                      int
	CoachHomeName                  string
	CoachAwayName                  string
	HomeWinningsDice               int
	AwayWinningsDice               int
	HomeScore                      int
	AwayScore                      int
	HomeInflictedTackles           int
	AwayInflictedTackles           int
	HomePossessionBall             int
	AwayPossessionBall             int
	HomeCashEarnedBeforeConcession int
	AwayCashEarnedBeforeConcession int
	HomeCashBeforeMatch            int
	AwayCashBeforeMatch            int
	HomeInflictedCasualties        int
	AwayInflictedCasualties        int
	HomePopularityBeforeMatch      int
	AwayPopularityBeforeMatch      int
	HomeOccupationTheir            int
	AwayOccupationTheir            int
	HomeSustainedTackles           int
	AwaySustainedTackles           int
	HomeInflictedMetersRunning     int
	AwayInflictedMetersRunning     int
	HomeMVP                        int
	AwayMVP                        int
	HomePopularityGain             int
	AwayPopularityGain             int
	HomeInflictedTouchdowns        int
	AwayInflictedTouchdowns        int
	HomeSustainedCasualties        int
	AwaySustainedCasualties        int
	HomeSustainedInjuries          int
	AwaySustainedInjuries          int
	HomeCashEarned                 int
	AwayCashEarned                 int
	HomeInflictedKO                int
	AwayInflictedKO                int
	TeamHomeName                   string
	TeamAwayName                   string
	HomeNbSupporters               int
	AwayNbSupporters               int
}

func Parse(r io.Reader) (Record, error) {
	var rr Replay
	decoder := xml.NewDecoder(r)
	err := decoder.Decode(&rr)
	if err != nil {
		return Record{}, err
	}

	record := NewRecordFromReplay(rr)

	return record, nil
}
