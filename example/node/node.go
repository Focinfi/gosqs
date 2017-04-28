package main

import "github.com/Focinfi/sqs/node"

func main() {
	node.New(":54461", 54461, ":5446").Start()
}
