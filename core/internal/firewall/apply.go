package firewall

import (
	"os"
	"os/exec"
)

func ApplyRules(rules string) error {

	tmpFile := "/run/sentinelos.nft"

	err := os.WriteFile(tmpFile, []byte(rules), 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("nft", "-f", tmpFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
