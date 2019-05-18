package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/ofonimefrancis/onefootball/models"
)

const (
	LIMIT = 5000
)

var URL = "https://vintagemonster.onefootball.com/api/teams/en/%s.json"

//Result Reperesents the information that should be outputted on
type Result struct {
	Name    string   `json:"name"`
	Age     string   `json:"age"`
	Country string   `json:"country"`
	Team    []string `json:"team"`
}

//ResultList - Implements the Sort interface for sorting players name alphabetically
type ResultList []Result

func (r ResultList) Len() int      { return len(r) }
func (r ResultList) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r ResultList) Less(i, j int) bool {
	return r[i].Name < r[j].Name
}

var allPlayers ResultList
var allWithoutDuplicates = make(map[string]Result)

func main() {
	goGroup := new(sync.WaitGroup)

	for i := 1; i <= LIMIT; i++ {
		go getTeamPlayers(i, goGroup)
	}

	goGroup.Add(LIMIT)
	goGroup.Wait()

	for _, player := range permutationOfPlayers(allPlayers) {
		fmt.Printf("%s; %s; %s\n", player.Name, player.Age, strings.Join(player.Team, ", "))
	}

}

//getTeamPlayers - Given a team
func getTeamPlayers(teamID int, goGroup *sync.WaitGroup) {
	defer goGroup.Done()
	teamIDInt := strconv.Itoa(teamID)
	url := fmt.Sprintf(URL, teamIDInt)
	resp, err := Get(url)
	if err != nil {
		return
	}
	var responseData models.Response
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		log.Println("Error parsing response body")
		return
	}
	ok := responseData.InRequiredTeam()
	if ok {
		teamName := responseData.Data.Team.Name
		players := responseData.Data.Team.Players
		for _, player := range players {
			result, ok := allWithoutDuplicates[player.Name]
			if ok {
				result.Team = append(result.Team, teamName)
				index := getPlayerIndex(player.Name)
				allWithoutDuplicates[player.Name] = result
				allPlayers[index] = result
			} else {
				newResult := Result{Name: player.Name, Age: player.Age, Team: []string{teamName}}
				allWithoutDuplicates[player.Name] = newResult
				allPlayers = append(allPlayers, newResult)
			}
		}
	}
}

func getPlayerIndex(currentPlayer string) int {
	for index, player := range allPlayers {
		if player.Name == currentPlayer {
			return index
		}
	}
	return -1
}

func permutationOfPlayers(resultList ResultList) ResultList {
	rl := make(ResultList, len(resultList))
	i := 0
	for _, player := range resultList {
		if len(player.Name) == 0 || len(player.Team) == 0 {
			continue
		}
		rl[i] = Result{Name: player.Name, Age: player.Age, Team: player.Team}
		i++
	}
	sort.Sort(rl)
	return rl
}

//Get - Makes a Get request
func Get(url string) (*http.Response, error) {
	return http.Get(url)
}
