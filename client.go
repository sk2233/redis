/*
@author: sk
@date: 2024/5/13
*/
package main

import "net"

type Client struct {
	Address string
}

func (c *Client) ExecuteCmd(cmd *Cmd) *Data {
	conn, err := net.Dial("tcp", c.Address)
	HandleErr(err)
	defer conn.Close()
	WriteAny(conn, cmd)
	res := &Data{}
	ReadAny(conn, res)
	return res
}

func NewClient(address string) *Client {
	return &Client{Address: address}
}
