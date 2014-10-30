package daemon

import (
	"os"
	"os/exec"
	"strings"
)

func IptablesMasq(not, src string) error {
	//  iptables -t nat -A POSTROUTING  -s 172.16.2.0/24 ! -d  172.16.0.0/12 -j MASQUERADE
	cmdline := "iptables -t nat -A POSTROUTING  -s " + src + " ! -d " + not + " -j MASQUERADE"

	arg := strings.Fields(cmdline)
	cmd := exec.Command(arg[0], arg[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
