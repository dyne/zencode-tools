# Zencode tools

Set of terminal based tools (TUI) to ease development of Zencode

One can start cloning this and installing golang apt-get install
golang should suffice, min 1.16 is OK, debian latest packages 1.18.

This repository is currently an import of some [bubbletea examples](https://github.com/charmbracelet/bubbletea), a nice library for TUI.

It will provide two tools: **zenexplorer** and **zendebug**.

## Zenexplorer

Zenexplorer is a TUI for zencode code completion, it does fuzzy matching of zencode statements as you type, in a terminal. it shows the statement, the scenario it belongs (as description below in gray) and when hit enter it copies the statement inside the clipboard.

Inside zencode-tools type `make` then run `./zenexplorer` and you will see.

Code is in [](src/zenexplorer/main.go)

Zenexplorer is very simple and just needs an array from zencode statement introspection, perhaps we can produce the latest JSON array from a github action.

## Zendebug

Zendebug is a more complex TUI: a triple pane editor (may be more, may zoom into one) supposed to contain zencode, data and stderr/out. The key Ctrl-enter should execute zencode in the left pane with data from center pane and show stderr/out in the right pane. For a start, this doesn't work yet (search for TODO).

Next we will need to be able to set breakpoints in zencode lines (`textarea.Prompt` chars) and then execute zenroom with a special configuration break=number which Zenroom core will implement soon: Zenroom will break at a specific zencode line and print out a json representation of the HEAP at that moment.

Then Zendebug will have keys for step or next and run to restart: every step zenroom will be re-executed so code can change during a debug session without problems.

# License

Software Tools to work with [Zenroom](https://dev.zenroom.org)

Copyright (C) 2022 Dyne.org foundation

Originally written as example code in Bubblewrap 

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public
License along with this program.  If not, see
[gnu.org/licenses](https://www.gnu.org/licenses/).

