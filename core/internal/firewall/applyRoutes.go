package firewall

import (
	"os"
	"os/exec"
)

func ApplyRoutes(routes string) error {

	tmpFile := "/run/sentinelos.routes"

	err := os.WriteFile(tmpFile, []byte(routes), 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("/bin/sh", "-c", routes)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
