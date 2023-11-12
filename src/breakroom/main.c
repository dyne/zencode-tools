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

#define CMPL(name) if(progmatch(buf, name)) bestlineAddCompletion(lc, name)

void completion(const char *buf, bestlineCompletions *lc) {
  CMPL("run");
  CMPL("list");
  CMPL("script");
  CMPL("keys");
  CMPL("data");
  CMPL("extra");
  CMPL("break");
  CMPL("clear");
  CMPL("conf");
}

char *hints(const char *buf, const char **ansi1, const char **ansi2) {
    if (buf[0] == '\0') {
        *ansi1 = "\033[35m"; /* magenta foreground */
        *ansi2 = "\033[39m"; /* reset foreground */
        return " run | script | keys | data | extra | break | clear | conf";
    }
    return NULL;
}

void set(char *key, char *value) { // set a key/value pair in breakroom conf
  fprintf(stderr, "%s %s\n", key, value);
  fflush(stderr);
}
#define CMD(cmd) if(!strcmp(tok,cmd)) { tok=strtok(NULL, " "); set(cmd, tok); exitcode=0; break; }
#define SETINT(cmd) if(!strcmp(tok,cmd)) { tok=strtok(NULL, " "); bool ok=true; int l=strlen(tok); for(int i=0;i<l;i++) { if(!isdigit(tok[i])) ok=false; } if(!ok) { fprintf(stdout, "Not an integer: %s\n",tok); free(line); continue; }; set(cmd, tok); exitcode=0; break; }
#define SETCONF(cmd) if(!strcmp(tok,cmd)) { tok=strtok(NULL, " "); bool ok=true; int l=strlen(tok); for(int i=0;i<l;i++) { if(! (isalnum(tok[i])||tok[i]==','||tok[i]=='='||tok[i]==' ')) ok=false; } if(!ok) { fprintf(stdout, "Invalid configuration string: %s\n",tok); free(line); continue; }; set(cmd, tok); exitcode=0; break; }

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
	CMD("run");
	CMD("list");
	CMD("script");
	CMD("keys");
	CMD("data");
	CMD("extra");
	CMD("clear");
	SETINT("break");
	CMD("conf");
	CMD("trace");
	CMD("heap");
	CMD("codec"); CMD("schema"); CMD("given");
	fprintf(stdout,"Unknown command: %s\n", line);
	free(line);
  }
  if(line) free(line);
  bestlineHistorySave(HISTFILE); /* Save the history on disk. */
  exit(exitcode);
}
