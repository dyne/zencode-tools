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
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

func acceptConnection(listener net.Listener, conn *net.Conn) {
	var err error
	for {
		*conn, err = listener.Accept()
		if err != nil {
			log.Error("Error while listening for connections: ", err)
			return
		}
		for {
			buf := make([]byte, 1024)
			_, err := (*conn).Read(buf[:])
			if err != nil {
				if err == io.EOF {
					log.Info("Closed connection to client")
					*conn = nil
					break
				} else {
					log.Error("Unknown error ", err)
					return
				}
			}
		}
	}
}

func startServer(listener net.Listener, channel chan string) {
	var conn net.Conn = nil

	go acceptConnection(listener, &conn)
	for stmt := range channel {
		if conn != nil {
			_, err := conn.Write([]byte(stmt))
			if err != nil {

			}
		}
	}

}
