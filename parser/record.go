package parser

import (
	"fmt"

	"github.com/google/uuid"
)

type Record struct {
	ID   uuid.UUID
	Home TeamStats
	Away TeamStats
}

type TeamStats struct {
	Name         string
	Cheerleaders int
	Supporters   int
	Popularity   int
	Race         Race

	InflictedInjuries          int
	SustainedKO                int
	OccupationOwn              int
	Value                      int
	CoachName                  string
	WinningsDice               int
	Score                      int
	InflictedTackles           int
	PossessionBall             int
	CashEarnedBeforeConcession int
	CashBeforeMatch            int
	InflictedCasualties        int
	PopularityBeforeMatch      int
	OccupationTheir            int
	SustainedTackles           int
	InflictedMetersRunning     int
	MVP                        string
	PopularityGain             int
	InflictedTouchdowns        int
	SustainedCasualties        int
	SustainedInjuries          int
	CashEarned                 int
	InflictedKO                int
	NbSupporters               int

	PlayerResults []PlayerResult
}

func NewRecordFromReplay(replay Replay) Record {
	finished := replay.ReplaySteps[len(replay.ReplaySteps)-1].RulesEventGameFinished
	stats := finished.Statistics
	homeTeam := finished.Coaches[0].TeamResult
	awayTeam := finished.Coaches[1].TeamResult

	home := TeamStats{
		Name:                       stats.TeamHomeName,
		Cheerleaders:               homeTeam.Cheerleaders,
		Supporters:                 homeTeam.Supporters,
		Popularity:                 homeTeam.Popularity,
		Race:                       homeTeam.Race,
		InflictedInjuries:          stats.HomeInflictedInjuries,
		SustainedKO:                stats.HomeSustainedKO,
		OccupationOwn:              stats.HomeOccupationOwn,
		Value:                      stats.HomeValue,
		CoachName:                  stats.CoachHomeName,
		WinningsDice:               stats.HomeWinningsDice,
		Score:                      stats.HomeScore,
		InflictedTackles:           stats.HomeInflictedTackles,
		PossessionBall:             stats.HomePossessionBall,
		CashEarnedBeforeConcession: stats.HomeCashEarnedBeforeConcession,
		CashBeforeMatch:            stats.HomeCashBeforeMatch,
		InflictedCasualties:        stats.HomeInflictedCasualties,
		PopularityBeforeMatch:      stats.HomePopularityBeforeMatch,
		OccupationTheir:            stats.HomeOccupationTheir,
		SustainedTackles:           stats.HomeSustainedTackles,
		InflictedMetersRunning:     stats.HomeInflictedMetersRunning,
		MVP:                        "",
		PopularityGain:             stats.HomePopularityGain,
		InflictedTouchdowns:        stats.HomeInflictedTouchdowns,
		SustainedCasualties:        stats.HomeSustainedCasualties,
		SustainedInjuries:          stats.HomeSustainedInjuries,
		CashEarned:                 stats.HomeCashEarned,
		InflictedKO:                stats.HomeInflictedKO,
		NbSupporters:               stats.HomeNbSupporters,
		PlayerResults:              homeTeam.PlayerResults,
	}

	for _, player := range homeTeam.PlayerResults {
		if player.MVP {
			home.MVP = player.Name
		}
	}

	away := TeamStats{
		Name:                       stats.TeamAwayName,
		Cheerleaders:               awayTeam.Cheerleaders,
		Supporters:                 awayTeam.Supporters,
		Popularity:                 awayTeam.Popularity,
		Race:                       awayTeam.Race,
		InflictedInjuries:          stats.AwayInflictedInjuries,
		SustainedKO:                stats.AwaySustainedKO,
		OccupationOwn:              stats.AwayOccupationOwn,
		Value:                      stats.AwayValue,
		CoachName:                  stats.CoachAwayName,
		WinningsDice:               stats.AwayWinningsDice,
		Score:                      stats.AwayScore,
		InflictedTackles:           stats.AwayInflictedTackles,
		PossessionBall:             stats.AwayPossessionBall,
		CashEarnedBeforeConcession: stats.AwayCashEarnedBeforeConcession,
		CashBeforeMatch:            stats.AwayCashBeforeMatch,
		InflictedCasualties:        stats.AwayInflictedCasualties,
		PopularityBeforeMatch:      stats.AwayPopularityBeforeMatch,
		OccupationTheir:            stats.AwayOccupationTheir,
		SustainedTackles:           stats.AwaySustainedTackles,
		InflictedMetersRunning:     stats.AwayInflictedMetersRunning,
		MVP:                        "",
		PopularityGain:             stats.AwayPopularityGain,
		InflictedTouchdowns:        stats.AwayInflictedTouchdowns,
		SustainedCasualties:        stats.AwaySustainedCasualties,
		SustainedInjuries:          stats.AwaySustainedInjuries,
		CashEarned:                 stats.AwayCashEarned,
		InflictedKO:                stats.AwayInflictedKO,
		NbSupporters:               stats.AwayNbSupporters,
		PlayerResults:              awayTeam.PlayerResults,
	}

	for _, player := range awayTeam.PlayerResults {
		if player.MVP {
			away.MVP = player.Name
		}
	}

	id := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(fmt.Sprintf("%s-%s:%s-%s", home.Name, home.CoachName, away.Name, away.CoachName)))

	return Record{
		ID:   id,
		Home: home,
		Away: away,
	}
}
