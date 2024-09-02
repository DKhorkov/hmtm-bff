package main

import (
	"fmt"
	"hmtm_bff/configs"
)

func main() {
	var config = configs.GetConfig()

	fmt.Println(config.Graphql.Port)
}
