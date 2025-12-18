package firewall

import (
	"os"
	"os/exec"
)

func ApplyVlansConfig(vlans string) error {

	cmd := exec.Command("/bin/sh", "-c", vlans)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
