package vpn

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/pquerna/otp/totp"
)

func start(network string, u *user.User) error {

	vpnPath := filepath.Join(u.HomeDir, ".vpn")
	vpnUserPath := filepath.Join(vpnPath, "username")

	totpPath := filepath.Join(u.HomeDir, ".totp", "secret")

	if _, err := os.Stat(vpnUserPath); err != nil {
		if os.IsNotExist(err) {
			err = ioutil.WriteFile(vpnUserPath, []byte(u.Username), 0600)
			if err != nil {
				return err
			}
		}
	}

	totpData, err := ioutil.ReadFile(totpPath)
	if err != nil {
		return err
	}

	code, err := totp.GenerateCode(string(totpData), time.Now())

	fmt.Println(code)

	targetDir := filepath.Join(vpnPath, u.Username+"-"+network)

	err = os.Chdir(targetDir)
	if err != nil {
		return err
	}

	command := fmt.Sprintf(`sudo openvpn --daemon "openvpn %s" --remote bastion.%s.nsone.co 1194 --ca ca.crt --key %s.key --cert %s.crt --proto udp --resolv-retry infinite --nobind --user nobody --group nogroup --persist-key --persist-tun --max-routes 1000 --verb 3 --comp-lzo --dev tun --client --remote-cert-tls server --auth-user-pass %s --auth-nocache --reneg-sec 0`, network, network, u.Username, u.Username, vpnUserPath)

	//cmd := exec.Command("/bin/bash", "-c", command)
	fmt.Println(command)

	return nil
}
