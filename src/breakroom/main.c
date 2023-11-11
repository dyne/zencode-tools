/* This file is part of Zenroom (https://zenroom.org)
 *
 * Copyright (C) 2023 Dyne.org foundation
 * designed, written and maintained by Denis Roio <jaromil@dyne.org>
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

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <errno.h>
#include <ctype.h>

#include "zenroom.h"
#include "bestline.h"

#define HISTFILE "breakroom_history.txt"

extern void load_file(char *dst, FILE *fd);

int breakpoint = 0;

char script[MAX_ZENCODE];
char scriptfile[MAX_STRING];

char data[MAX_FILE];
char datafile[MAX_STRING];

char keys[MAX_FILE];
char keysfile[MAX_STRING];

char extra[MAX_FILE];
char extrafile[MAX_STRING];

char conf[MAX_CONFIG];

bool progmatch(const char *buf, const char *word) {
  int len = strlen(buf);
  bool match = 0;
  register int i;
  for(i=0;i<len;i++) {
	if(buf[i] == word[i]) match = true;
	else match = false;
  }
  return(match);
}

void completion(const char *buf, bestlineCompletions *lc) {
  if( progmatch(buf, "help") )   bestlineAddCompletion(lc, "help");
  if( progmatch(buf, "load") )   bestlineAddCompletion(lc, "load");
  if( progmatch(buf, "break") )  bestlineAddCompletion(lc, "break");
  if( progmatch(buf, "heap") )   bestlineAddCompletion(lc, "heap");
  if( progmatch(buf, "trace") )  bestlineAddCompletion(lc, "trace");
  if( progmatch(buf, "source") ) bestlineAddCompletion(lc, "source");
}

char *hints(const char *buf, const char **ansi1, const char **ansi2) {
  if(buf[2] != '\0') return "";
  if(!strcmp(buf,"/help")) return "";
  if(!strcmp(buf,"/break")) return " line number";

    if (!strcmp(buf,"When")) {
        *ansi1 = "\033[35m"; /* magenta foreground */
        *ansi2 = "\033[39m"; /* reset foreground */
        return " I create ";
    }

    if (!strcmp(buf,"When I create")) {
        *ansi1 = "\033[35m"; /* magenta foreground */
        *ansi2 = "\033[39m"; /* reset foreground */
        return " keyring | random | public key ... ";
    }

    if (buf[0] == '\0') {
        *ansi1 = "\033[35m"; /* magenta foreground */
        *ansi2 = "\033[39m"; /* reset foreground */
        return " help | list | break | clear | run | heap | trace | conf | script | data | keys | extra";
    }

    return NULL;
}

void set(char *key, char *value) { // set a key/value pair in breakroom conf
  fprintf(stderr, "%s %s\n", key, value);
  // fprintf(stdout, "set %s %s\n", key, value);
  fflush(stderr);
}
#define LOAD(cmd) if(!strcmp(tok,cmd)) { tok=strtok(NULL, " "); fd=fopen(tok, "rb"); if(!fd) { fprintf(stdout, "%s\n",strerror(errno)); free(line); continue; }; set(cmd, tok); exitcode=0; break; }
#define SETINT(cmd) if(!strcmp(tok,cmd)) { tok=strtok(NULL, " "); bool ok=true; int l=strlen(tok); for(int i=0;i<l;i++) { if(!isdigit(tok[i])) ok=false; } if(!ok) { fprintf(stdout, "Not an integer: %s\n",tok); free(line); continue; }; set(cmd, tok); exitcode=0; break; }


int main(int argc, char **argv) {
  char *line = NULL;
  char *tok = NULL;
  FILE *fd = NULL;
  int exitcode = 1;
  /* Set the completion callback. This will be called every time the
   * user uses the <tab> key. */
  bestlineSetCompletionCallback(completion);
  bestlineSetHintsCallback(hints);

  /* Load history from file. The history file is just a plain text file
   * where entries are separated by newlines. */
  bestlineHistoryLoad(HISTFILE); /* Load the history at startup */

  /* Now this is the main loop of the typical bestline-based application.
   * The call to bestline() will block as long as the user types something
   * and presses enter.
   *
   * The typed string is returned as a malloc() allocated string by
   * bestline, so the user needs to free() it. */
    
  while((line = bestlineWithHistory("breakroom> ", HISTFILE)) != NULL) {
	tok = strtok(line, " ");
	if(!tok) continue;
	LOAD("script");
	LOAD("keys");
	LOAD("data");
	LOAD("extra");
	SETINT("break");
	fprintf(stdout,"Unknown command: %s\n", line);
	free(line);
  }
  if(line) free(line);
  bestlineHistorySave(HISTFILE); /* Save the history on disk. */
  exit(exitcode);
}
