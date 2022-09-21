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
	"sync"
)

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

	r.scenarios = []ZencodeScenario{
		ZencodeScenario{
			scenario: "given",
			statements: []string{"stat 1", "stat 2"},
		},
	}
	r.finished = false
}

func (r *zencodeItemGenerator) next() (item, bool) {
	if r.mtx == nil {
		r.reset()
	}

	r.mtx.Lock()
	defer r.mtx.Unlock()

	i := item{
		scenario: r.scenarios[r.scenarioIndex].scenario,
		statement: r.scenarios[r.scenarioIndex].statements[r.statementIndex],
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
