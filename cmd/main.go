package main

import (
	"fmt"
	"udo_mass/pkg/api"
	"udo_mass/pkg/logger"
)

func main() {
	//formula1 := "KNO3"
	//formula2 := "KH2PO4"
	//formula3 := "K2SO4"
	//formula4 := "FeSO4*7H2O

	fmt.Println("// -------------------------------------------------------------------------")

	// экземпляр api
	httpServer := api.New()
	err := httpServer.Start()
	if err != nil {
		logger.Fatal("Ошибка при запуске сервера:", err)
	}

	api.GraceShutdown(httpServer)

}
