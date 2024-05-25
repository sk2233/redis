/*
@author: sk
@date: 2024/5/14
*/
package main

const (
	CmdGet    = "get"
	CmdSet    = "set"
	CmdDel    = "del"
	CmdKeys   = "keys"
	CmdZAdd   = "zadd"   // zadd key score name
	CmdZRem   = "zrem"   // zrem key name
	CmdZScore = "zscore" // zscore key name
	CmdZQuery = "zquery" // zquery key score name offset limit
	CmdExpire = "expire" // expire key seconds  //  seconds = -1  取消超时
	CmdTTL    = "ttl"    // ttl key
)

const (
	SkipDep = 8
)
