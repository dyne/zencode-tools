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
	"bufio"
	"strings"
	_ "embed"
	"sync"
)

//go:embed introspection.txt
var introspection string

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

	r.scenarios = []ZencodeScenario{}
	sc := bufio.NewScanner(strings.NewReader(introspection))
	scenario := ""
	statements := []string{}
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			if(len(statements) > 0) {
				r.scenarios = append(r.scenarios, ZencodeScenario {
					scenario: scenario,
					statements: statements,
				})
				scenario = ""
				statements = []string{}
			}
		} else if scenario == "" {
			scenario = line
		} else {
			statements = append(statements, line)
		}
	}

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
	i := item{
		scenario: r.scenarios[r.scenarioIndex].scenario,
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
