package main

import (
	"fmt"
	"log"
	"net/http"
	"sentinelos/core/internal/api/handlers"
	"sentinelos/core/internal/auth"
	"sentinelos/core/internal/config"
	"sentinelos/core/internal/firewall"
	"sentinelos/core/internal/network"
	"sentinelos/core/internal/system"
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

	/*  *********************************************************** */
	interfaces := network.GenerateInterfacesConfig(fw.Interfaces)
	err = firewall.ApplyInterfacesConfig(interfaces)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("========== CONFIG INTERFACES CREADAS ==========")
	fmt.Println(interfaces)

	/*  *********************************************************** */

	vlans := network.GenerateVlansConfig(fw.Vlans)
	err = firewall.ApplyVlansConfig(vlans)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("========== CONFIG VLANS CREADAS ==========")
	fmt.Println(vlans)

	/*  *********************************************************** */

	routes := network.GenerateRoutesConfig(fw.Routes)
	err = firewall.ApplyRoutes(routes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("========== RUTAS CREADAS ==========")
	fmt.Println(routes)

	/*  *********************************************************** */
	var ifaceNames []string
	for name := range fw.Interfaces {
		ifaceNames = append(ifaceNames, name)
	}
	err = system.ApplySysctl(ifaceNames)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("========== SYSCTL CONFIGURADO ==========")

	/*  *********************************************************** */

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

	/*  *********************************************************** */

	dnsmasqConfig := network.GenerateDnsmasqConfig(fw.DHCPConfigs)
	err = firewall.ApplyDHCP(dnsmasqConfig)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("========== CONFIG DHCP CREADA ==========")
	fmt.Println(dnsmasqConfig)

	fmt.Println("[OK] Configuración DHCP aplicada")

	/*  *********************************************************** */

	// 1️⃣ Cargar usuarios
	users, err := auth.LoadUsers("/srv/sentinelos/core/internal/auth/users.yml")
	if err != nil {
		log.Fatalf("error loading users: %v", err)
	}

	// 2️⃣ Crear AuthService
	authService := auth.NewAuthService(users)

	// 3️⃣ Router básico
	mux := http.NewServeMux()

	// 4️⃣ Ruta de login
	mux.HandleFunc("/api/login", handlers.LoginHandler(authService))

	// 5️⃣ Levantar servidor
	log.Println("SentinelOS API listening on :8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
