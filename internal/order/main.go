package main

import (
	"log"

	"github.com/loveRyujin/gorder/common/config"
	"github.com/spf13/viper"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		panic(err)
	}
}

func main() {
	log.Println(viper.Get("order"))
}
