package main

import (
	"log"

	"sentinelos/core/internal/api"
	"sentinelos/core/internal/config"
	"sentinelos/core/internal/runtime"
	"sentinelos/core/internal/system"
)

func main() {

	path := "/srv/sentinelos/core/internal/config/config.yml"

	// lectura del archivo YAML
	raw, err := config.LoadRawConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	// construir el modelo firewall
	fw, err := config.BuildFirewall(raw)
	if err != nil {
		log.Fatal(err)
	}

	// aplicar runtime completo 
	if err := runtime.ApplyFullRuntime(fw); err != nil {
		log.Fatal(err)
	}

	// guardar como running config
	system.SetFirewall(fw)
    
	/*
	LCD status activar luego
	stopLCD := make(chan struct{}) 
	go utils.StartLCDStatus(stopLCD) 
	fmt.Println("Estado del hardware en lcd actualizado")
     */
	// iniciar API REST
	go api.StartAPIServer()

	// mantener el daemon up
	select {}
}