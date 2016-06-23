/* pg_check_cert - Panic when your postgresql certificate is about to expire
 * Copyright (C) 2016  Chris Aumann
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <host:port>\n", os.Args[0])
		os.Exit(1)
	}

	pgHost := os.Args[1]

	var err error
	cn := &conn{}

	// TODO:
	cn.c, err = net.Dial("tcp", pgHost)
	if err != nil {
		panic(err)
	}

	w := cn.writeBuf(0)
	w.int32(80877103)
	cn.sendStartupPacket(w)

	b := cn.scratch[:1]
	_, err = io.ReadFull(cn.c, b)
	if err != nil {
		panic(err)
	}

	if b[0] != 'S' {
		panic("SSL not supported")
	}

	tlsConf := tls.Config{}
	tlsConf.InsecureSkipVerify = true
	client := tls.Client(cn.c, &tlsConf)

	expiresIn := cn.expiresIn(client, &tlsConf)
	fmt.Println(expiresIn)
}

// Check how many days the certificate is still valid
func (cn *conn) expiresIn(client *tls.Conn, tlsConf *tls.Config) int {
	err := client.Handshake()
	if err != nil {
		panic(err)
	}
	certs := client.ConnectionState().PeerCertificates

	expiresIn := certs[0].NotAfter.Sub(time.Now())

	// Convert days to int
	return int(expiresIn.Hours() / 24)
}
