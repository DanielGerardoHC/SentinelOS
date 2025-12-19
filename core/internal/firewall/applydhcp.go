package firewall

import (
	"os"
	"os/exec"
)

const (
	configDir  = "/etc/sentinelos"
	configFile = "/etc/sentinelos/dnsmasq.conf"
)

func ApplyDHCP(conf string) error {

	// 1. Crear directorio
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// 2. Escribir config
	if err := os.WriteFile(configFile, []byte(conf), 0644); err != nil {
		return err
	}

	// 3. Detener dnsmasq previo (si existe)
	exec.Command("pkill", "dnsmasq").Run()

	// 4. Arrancar dnsmasq gestionado por SentinelOS
	cmd := exec.Command(
		"dnsmasq",
		"--conf-file="+configFile,
		"--log-dhcp",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Start()
}
