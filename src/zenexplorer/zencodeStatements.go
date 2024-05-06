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
	"encoding/json"
	"strings"
	"sync"
)

type ZenStatements struct {
	Given   []string            `json:"given"`
	When  	map[string][]string `json:"when"`
	If    	map[string][]string `json:"if"`
	Foreach []string            `json:"foreach"`
	Then  	[]string            `json:"then"`
	mtx   	*sync.Mutex
}

//go:embed load_statements.lua
var loadStatementsScript string

//go:embed default_statements.json
var defaultStatements string


func (z *ZenStatements) reset() {
	z.mtx = &sync.Mutex{}

	z.mtx.Lock()
	defer z.mtx.Unlock()

	stdin := strings.NewReader(loadStatementsScript)
	res, err := ZenroomExec(stdin, ZenInput{})
	var statementsJson string
	if err {
		statementsJson = defaultStatements
	} else {
		statementsJson = res.Output
	}

	json.Unmarshal([]byte(statementsJson), z)
}

func (z *ZenStatements) count() int {
	if z.mtx == nil {
		z.reset()
	}

	z.mtx.Lock()
	defer z.mtx.Unlock()

	var count int = 0
	for _, v := range z.When {
		count = count + len(v)
	}
	for _, v := range z.If {
		count = count + len(v)
	}
	return count + len(z.Given) + len(z.Then) + len(z.Foreach)
}

