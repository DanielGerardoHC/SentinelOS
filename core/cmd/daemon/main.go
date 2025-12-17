package main

import (
	"fmt"
	"log"
	"sentinelos/core/internal/config"
	"sentinelos/core/internal/firewall"
	"sentinelos/core/internal/network"
)

func main() {

	path := "/srv/sentinelos/core/internal/config/config.yml"

	raw, err := config.LoadRawConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[OK] Archivo YAML cargado correctamente")

	/*
		fmt.Printf("Zonas: %d\n", len(raw.Zones))
		fmt.Printf("Interfaces: %d\n", len(raw.Interfaces))
		fmt.Printf("Policies: %d\n", len(raw.Policies))
	*/
	fw, err := config.BuildFirewall(raw)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[OK] Model Firewall construido correctamente")
	/*
		fmt.Printf("Zonas: %d\n", len(fw.Zones))
		fmt.Printf("Interfaces: %d\n", len(fw.Interfaces))
		fmt.Printf("Policies: %d\n", len(fw.Policies))
	*/
	routes := network.GenerateRoutesConfig(fw.Routes)
	err = firewall.ApplyRoutes(routes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("========== RUTAS CREADAS ==========")
	fmt.Println(routes)

	rules := firewall.GenerateRules(fw)

	fmt.Println("========== NFTABLES CREADAS ==========")

	rules += firewall.GenerateNATRules(fw)
	fmt.Println("========== NATRULES CREADAS ==========")
	fmt.Println(rules)

	err = firewall.ApplyRules(rules)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[OK] Reglas aplicadas")

}
