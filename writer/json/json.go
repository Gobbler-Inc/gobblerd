package json

import (
	"bbrz/parser"
	"encoding/json"
	"fmt"
	"os"
)

type OutputFormat struct {
	HomeTeam TeamInfo
	AwayTeam TeamInfo
}

type TeamInfo struct {
	Name    string
	Type    string
	Coach   string
	Value   int
	Rerolls int
	Players []Player
}

type Player struct {
	Name     string
	Type     string
	Movement int
	Strength int
	Agility  int
	Armor    int
	Skills   []string
}

func WriteJSON(replay parser.Replay) error {
	firstStep := replay.ReplaySteps[0]
	filename := fmt.Sprintf("%s-%s-%s_v_%s.json",
		firstStep.GameInfos.League.Name,
		firstStep.GameInfos.Competition.Name,
		firstStep.GameInfos.CoachesInfos.CoachInfos[0].Login,
		firstStep.GameInfos.CoachesInfos.CoachInfos[1].Login)

	fp, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to open %s for writing: %w", filename, err)
	}
	defer fp.Close()

	homeTeamInfo := createTeamData(firstStep.GameInfos.CoachesInfos.CoachInfos[0], firstStep.BoardState.ListTeams.TeamState[0])
	awayTeamInfo := createTeamData(firstStep.GameInfos.CoachesInfos.CoachInfos[1], firstStep.BoardState.ListTeams.TeamState[1])

	output := OutputFormat{
		HomeTeam: homeTeamInfo,
		AwayTeam: awayTeamInfo,
	}

	encoder := json.NewEncoder(fp)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("Failed to encode JSON: %w", err)
	}

	return nil
}

func createTeamData(coach parser.CoachInfos, teamState parser.Team) TeamInfo {
	teamInfo := TeamInfo{
		Name:    teamState.Name,
		Type:    string(teamState.Race),
		Coach:   coach.Login,
		Value:   teamState.Value,
		Rerolls: teamState.Rerolls,
		Players: nil,
	}

	playerList := make([]Player, 0)
	for _, player := range teamState.ListPitchPlayers.PlayerStates {
		playerList = append(playerList, Player{
			Name:     player.Name,
			Type:     player.Type,
			Movement: player.Movement,
			Strength: player.Strength,
			Agility:  player.Agility,
			Armor:    player.Armor,
			Skills:   player.Skills,
		})
	}

	teamInfo.Players = playerList

	return teamInfo
}
