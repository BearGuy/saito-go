package main

import "github.com/bearguy/saito-go/cmd/saito"

func main() {
	saito := saito.InitSaito()
	saito.Mempool.Bundle()
}
