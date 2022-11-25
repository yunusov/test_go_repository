package main

import (
	"go_diploma/pkg/service"
	"go_diploma/pkg/utils"
)

func main() {
	config := utils.LoadSettings()
	service.Start(config)
}
