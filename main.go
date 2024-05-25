/*
@author: sk
@date: 2024/5/11
*/
package main

import (
	"fmt"
	"time"
)

//  https://build-your-own.org/redis/a2_reading   不太好的教程

func main() {
	NewService(":8080").Run()
	time.Sleep(time.Second * 3)
	client := NewClient(":8080")
	fmt.Println(client.ExecuteCmd(&Cmd{
		Cmd:  CmdSet,
		Args: []*Data{NewStrData("name"), NewStrData("你好")},
	}))
	fmt.Println(client.ExecuteCmd(&Cmd{
		Cmd:  CmdExpire,
		Args: []*Data{NewStrData("name"), NewNumData(3)},
	}))
	fmt.Println(client.ExecuteCmd(&Cmd{
		Cmd:  CmdGet,
		Args: []*Data{NewStrData("name")},
	}))
	time.Sleep(time.Second * 4)
	fmt.Println(client.ExecuteCmd(&Cmd{
		Cmd:  CmdGet,
		Args: []*Data{NewStrData("name")},
	}))
	//data := client.ExecuteCmd(&Cmd{
	//	Cmd:  CmdSet,
	//	Args: []*Data{NewStrData("name"), NewStrData("博丽灵梦")},
	//})
	//fmt.Println(data)
	//data = client.ExecuteCmd(&Cmd{
	//	Cmd: CmdKeys,
	//})
	//fmt.Println(data)
	//skipList := NewSkipList()
	//skipList.Add(2233, "test")
	//item := skipList.Get(2233)
	//fmt.Println(item)
	//skipList.Del(2233)
	//list := NewSkipList()
	//list.Add(10, "2233")
	//list.Add(10, "1122")
	//list.Add(5, "www")
	//list.Add(15, "dddd")
	//fmt.Println(list.Get(10, "2233"))
	//fmt.Println(list.Get(5, "www"))
	//fmt.Println(list.Get(15, "ssss"))
	//list.Del(10, "2233")
	//fmt.Println("=============================")
	//fmt.Println(list.Get(10, "2233"))
	//fmt.Println(list.Get(10, "1122"))
	//fmt.Println(list.Get(5, "www"))
	//fmt.Println(list.Get(15, "ssss"))
	//heap := NewHeap()
	//heap.Add("test", 100)
	//heap.Rem("test")
	//heap.Pop()
	//heap.Peek()
	//heap.IsEmpty()
}
