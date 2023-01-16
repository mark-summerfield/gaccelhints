// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"log"
)

func main() {
	log.SetFlags(log.Lmsgprefix)
	config := NewConfig()
	fmt.Println(config)
}
