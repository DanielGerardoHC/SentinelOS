package firewall

import (
	"os"
	"os/exec"
)

func ApplyInterfacesConfig(interfaces string) error {

	cmd := exec.Command("/bin/sh", "-c", interfaces)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
