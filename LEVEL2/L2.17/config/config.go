// Package config provides a connection structure with a method returning net.Conn
package config

import (
	"log"
	"net"
	"time"
)

type Connection struct {
	Port    string        // optional, default value: "25"
	Host    string        // mandatory
	Timeout time.Duration // optional
}

func (c *Connection) GetConnection() net.Conn {
	switch c.Timeout {
	case 0:
		conn, err := net.Dial("tcp", c.Host+":"+c.Port)
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		return conn
	default:
		conn, err := net.DialTimeout("tcp", c.Host+":"+c.Port, c.Timeout)
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		return conn
	}
}
