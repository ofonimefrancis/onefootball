# onefootball assessment
Given API endpoint: https://vintagemonster.onefootball.com/api/teams/en/%team_id%.json (%team_id% must be replaced with a real number)
Using the API endpoint, find the following teams by name:
1. Germany
2. England
3. France
4. Spain
5. Manchester Utd
6. Arsenal
7. Chelsea
8. Barcelona
9. Real Madrid
10. FC Bayern Munich


Extract all the players from the given teams and render to stdout the information about players alphabetically ordered by name.
Each player entry should contain the following information: full name; age; list of teams.

Example:
1. Alexander Mustermann; 25; France, Manchester Utd
2. Brad Exampleman; 30; Arsenal
3. ...

## Installation
Execute the program by running `go run main.go` in the root of the project directory.