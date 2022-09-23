/* Software Tools to work with Zenroom (https://dev.zenroom.org)
*
* Copyright (C) 2022 Dyne.org foundation
* Originally written as example code in Bubblewrap
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU Affero General Public License as
* published by the Free Software Foundation, either version 3 of the
* License, or (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU Affero General Public License for more details.
*
* You should have received a copy of the GNU Affero General Public License
* along with this program.  If not, see <https://www.gnu.org/licenses/>.
*
*/

package main

import (
	_ "embed"
	"sync"
	"strings"
	"encoding/json"
)

type ZenStatements struct {
	Given []string         `json:"given"`
	When map[string][]string `json:"when"`
	Then []string          `json:"then"`
}
//go:embed load_statements.lua
var loadStatementsScript string

//go:embed default_statements.json
var defaultStatements string

type ZencodeScenario struct {
	scenario string
	statements []string
}
type zencodeItemGenerator struct {
	finished       bool
	scenarios      []ZencodeScenario
	scenarioIndex  int
	statementIndex int
	mtx            *sync.Mutex
	shuffle        *sync.Once
}

func (r *zencodeItemGenerator) reset() {
	r.mtx = &sync.Mutex{}
	r.shuffle = &sync.Once{}

	stdin := strings.NewReader(loadStatementsScript)
	res, err := ZenroomExec(stdin, ZenInput{})
	var statementsJson string
	if err {
		statementsJson = defaultStatements
	} else {
		statementsJson = res.Output
	}


	var zen ZenStatements

	json.Unmarshal([]byte(statementsJson), &zen)
	r.scenarios = []ZencodeScenario{
		ZencodeScenario {
			scenario: "given",
			statements: zen.Given,
		},
	}
	for k, v := range zen.When {
		r.scenarios = append(r.scenarios, ZencodeScenario {
			scenario: k,
			statements: v,
		})
	}
	r.scenarios = append(r.scenarios,
		ZencodeScenario {
			scenario: "then",
			statements: zen.Then,
		},
	)

	r.finished = false
}

func (r *zencodeItemGenerator) count() int {
	if r.mtx == nil {
		r.reset()
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var count int = 0
	for i := 0; i < len(r.scenarios); i++ {
		count = count + len(r.scenarios[i].statements)
	}
	return count;
}

func (r *zencodeItemGenerator) next() (item, bool) {
	if r.mtx == nil {
		r.reset()
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	scenario := r.scenarios[r.scenarioIndex].scenario
	var begin string;
	if scenario == "given" || scenario == "then" {
		begin = strings.Title(scenario)
	} else {
		begin = "When"
	}

	if scenario == "default" || scenario == "given" || scenario == "then" {
		scenario = ""
	}
	begin = begin + " I "
	i := item{
		scenario: scenario,
		statement: begin + " " + r.scenarios[r.scenarioIndex].statements[r.statementIndex],
	}
	finished := r.finished

	r.statementIndex++
	if r.statementIndex >= len(r.scenarios[r.scenarioIndex].statements) {
		r.statementIndex = 0
		r.scenarioIndex++
		if r.scenarioIndex >= len(r.scenarios) {
			r.scenarioIndex = 0;
			r.finished = true
		}
	}

	return i, finished
}
