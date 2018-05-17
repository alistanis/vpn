package vpn

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"os/exec"

	"github.com/pquerna/otp/totp"
)

func start(network string, u *user.User, mfa bool) error {

	vpnPath := filepath.Join(u.HomeDir, ".vpn")
	vpnUserPath := filepath.Join(vpnPath, "username")

	totpPath := filepath.Join(u.HomeDir, ".totp", "secret")

	var code string
	if mfa {
		totpData, err := ioutil.ReadFile(totpPath)
		if err != nil {
			return err
		}
		code, err = totp.GenerateCode(string(totpData), time.Now())

		fmt.Println(code)
	}

	if _, err := os.Stat(vpnUserPath); err != nil {
		if os.IsNotExist(err) {
			err = ioutil.WriteFile(vpnUserPath, []byte(u.Username), 0600)
			if err != nil {
				return err
			}
		}
	}

	targetDir := filepath.Join(vpnPath, u.Username+"-"+network)

	err := os.Chdir(targetDir)
	if err != nil {
		return err
	}

	command := fmt.Sprintf(`sudo openvpn --daemon "openvpn %s" --remote bastion.%s.nsone.co 1194 --ca ca.crt --key %s.key --cert %s.crt --proto udp --resolv-retry infinite --nobind --user nobody --group nogroup --persist-key --persist-tun --max-routes 1000 --verb 3 --comp-lzo --dev tun --client --remote-cert-tls server --auth-user-pass %s --auth-nocache --reneg-sec 0`, network, network, u.Username, u.Username, vpnUserPath)

	cmd := exec.Command("/bin/bash", "-c", command)

	err = cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func stop(network string) error {
	command := fmt.Sprintf(`ps aux | grep openvpn | grep "%s" | awk '{print $2}' | xargs sudo kill -SIGTERM`, network)
	err := execute(command)
	if err != nil {
		if err.Error() == "signal: terminated" {
			return nil
		}
	}
	return err
}

func otp(u *user.User) error {
	totpPath := filepath.Join(u.HomeDir, ".totp", "secret")

	totpData, err := ioutil.ReadFile(totpPath)
	if err != nil {
		return err
	}

	code, err := totp.GenerateCode(string(totpData), time.Now())
	if err != nil {
		return err
	}
	fmt.Println(code)
	return nil
}

// make this take network for future compatibility
func status(network string) error {
	command := fmt.Sprintf(`ps aux | grep openvpn | grep -v grep`)
	err := execute(command)
	if err != nil {
		if err.Error() == "exit status 1" {
			fmt.Println("vpn is not running")
			return nil
		}
	}
	return err
}

func execute(command string) error {
	cmd := exec.Command("/bin/bash", "-c", command)
	data, err := cmd.Output()
	if len(data) > 0 {
		fmt.Print(string(data))
	}
	return err
}
