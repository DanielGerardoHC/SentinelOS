package system

import (
	"fmt"
	"os"
	"os/exec"
)

const sysctlFile = "/etc/sysctl.d/99-sentinelos.conf"

func ApplySysctl(interfaces []string) error {

	content := `net.ipv4.ip_forward=1
net.ipv4.conf.all.rp_filter=0
net.ipv4.conf.default.rp_filter=0
`

	if err := os.WriteFile(sysctlFile, []byte(content), 0644); err != nil {
		return err
	}

	// Aplicar global
	exec.Command("sysctl", "--system").Run()

	// Forzar por interfaz
	for _, iface := range interfaces {
		cmd := exec.Command(
			"sysctl",
			fmt.Sprintf("net.ipv4.conf.%s.rp_filter=0", iface),
		)
		cmd.Run()
	}

	return nil
}
