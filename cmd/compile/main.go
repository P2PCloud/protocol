package main

import (
	"github.com/sirupsen/logrus"

	"github.com/p2pcloud/protocol/implementations/evm"
)

func main() {
	err := evm.CompileContracts(true)
	if err != nil {
		panic(err)
	}
	logrus.Println("Contracts compiled")
}
