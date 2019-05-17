package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/ofonimefrancis/onefootball/models"
)

const (
	LIMIT = 10
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
	// var firstName = r[i].Name
	// var secondName = r[j].Name
	// var firstNameLower = strings.ToLower(firstName)
	// var secondNameLower = strings.ToLower(secondName)
	// if firstNameLower == secondNameLower {
	// 	return firstName < secondName
	// }
	//return firstNameLower < secondNameLower
	return r[i].Name < r[j].Name
}

func (r ResultList) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func main() {

	var allPlayers ResultList
	requiredTeams := []string{"Germany", "England", "France", "Spain", "Manchester Utd", "Arsenal", "Chelsea", "Barcelona", "Real Madrid", "FC Bayern Munich"}
	//var reOrderedPlayers ResultList

	for i := 1; i <= LIMIT; i++ {
		teamID := strconv.Itoa(i)
		url := fmt.Sprintf("https://vintagemonster.onefootball.com/api/teams/en/%s.json", teamID)
		resp, err := Get(url)
		if err != nil {
			continue //Skip route with error
		}
		var responseData models.Response
		err = json.NewDecoder(resp.Body).Decode(&responseData)
		if err != nil {
			log.Println("Error parsing response body")
			log.Println(err)
			break
		}
		ok := inRequiredTeams(responseData.Data.Team.Name, requiredTeams)
		if ok {
			//fmt.Printf("Found %s\nPlayers %v\n\n", responseData.Data.Team.Name, responseData.Data.Team.Players)
			allPlayers = append(allPlayers, RetrievePlayers(responseData.Data.Team.Name, responseData.Data.Team.Players)...)
		}
	}

	fmt.Println("[AllPlayers slice]")
	for _, player := range allPlayers {
		fmt.Printf("%s; %s; %s, %s\n", player.Name, player.Age, player.Country, player.Team)
	}

	fmt.Println("[Permutation]")
	for _, player := range permutationOfPlayers(allPlayers) {
		fmt.Printf("%s; %s; %s, %s\n", player.Name, player.Age, player.Country, player.Team)
	}
}

func RetrievePlayers(team string, players []models.Player) ResultList {
	resultList := ResultList{}
	for _, player := range players {
		if player.Name == "" || player.Age == "" {
			continue
		}
		result := Result{Name: player.Name, Age: player.Age, Country: player.Country, Team: team}
		resultList = append(resultList, result)
	}
	return resultList
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
