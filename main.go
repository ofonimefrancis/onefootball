package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/ofonimefrancis/onefootball/models"
)

const (
	LIMIT = 8000
)

type Result struct {
	Name    string `json:"name"`
	Age     string `json:"age"`
	Country string `json:"country"`
	Team    string `json:"team"`
}

type ResultList []Result

func (r ResultList) Len() int { return len(r) }
func (r ResultList) Less(i, j int) bool {
	return r[i].Name < r[j].Name
}

var allPlayers ResultList

func (r ResultList) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func main() {
	goGroup := new(sync.WaitGroup)

	for i := 1; i <= LIMIT; i++ {
		go getTeamPlayers(i, goGroup)
	}

	goGroup.Add(LIMIT)
	goGroup.Wait()

	for _, player := range permutationOfPlayers(allPlayers) {
		fmt.Printf("%s; %s; %s, %s\n", player.Name, player.Age, player.Country, player.Team)
	}
}

func getTeamPlayers(teamID int, goGroup *sync.WaitGroup) {
	teamId := strconv.Itoa(teamID)
	url := fmt.Sprintf("https://vintagemonster.onefootball.com/api/teams/en/%s.json", teamId)
	resp, err := Get(url)
	if err != nil {
		return
	}
	var responseData models.Response
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		log.Println("Error parsing response body")
		log.Println(err)
		return
	}
	ok := responseData.InRequiredTeam()
	if ok {
		players := responseData.Data.Team.Players
		for _, player := range players {
			result := Result{Name: player.Name, Age: player.Age, Country: player.Country, Team: responseData.Data.Team.Name}
			allPlayers = append(allPlayers, result)
		}
	}

	goGroup.Done()
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

func inRequiredTeams(team string, requiredTeams []string) bool {
	for _, requiredTeam := range requiredTeams {
		if requiredTeam == team {
			return true
		}
	}
	return false
}
