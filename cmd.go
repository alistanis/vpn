package vpn

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"os/exec"

	"encoding/json"

	"github.com/pquerna/otp/totp"
)

func start(network string, u *user.User) error {

	networkPath, vpnUserPath := getPaths(u, network)
	configPath := filepath.Join(networkPath, "config.json")

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	c := &config{}

	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}

	var code string
	if c.MFASecret != "" {
		code, err = totp.GenerateCode(c.MFASecret, time.Now())
		if err != nil {
			return err
		}
	}

	vpnFileData := u.Username + "\n"
	if code != "" {
		vpnFileData += code
	}

	err = ioutil.WriteFile(vpnUserPath, []byte(vpnFileData), 0600)
	if err != nil {
		return err
	}

	err = os.Chdir(networkPath)
	if err != nil {
		return err
	}

	//bastion.%s.nsone.co 1194

	command := fmt.Sprintf(`sudo openvpn --daemon "openvpn %s" --remote %s --ca ca.crt --key %s.key --cert %s.crt --proto udp --resolv-retry infinite --nobind --user nobody --group nogroup --persist-key --persist-tun --max-routes 1000 --verb 3 --comp-lzo --dev tun --client --remote-cert-tls server --auth-user-pass %s --auth-nocache --reneg-sec 0`, network, c.URLString(), u.Username, u.Username, vpnUserPath)
	fmt.Println(command)
	cmd := exec.Command("/bin/bash", "-c", command)

	return cmd.Run()
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

func otp(u *user.User, network string) error {
	networkPath, _ := getPaths(u, network)
	configPath := filepath.Join(networkPath, "config.json")

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	c := &config{}
	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}

	code, err := totp.GenerateCode(c.MFASecret, time.Now())
	if err != nil {
		return err
	}
	fmt.Println(code)
	return nil
}

func status(network string) error {
	command := fmt.Sprintf(`ps aux | grep openvpn | grep -v grep | grep %s`, network)
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

func configure(network string, c *config, u *user.User) error {
	if c.Port == -1 {
		return errors.New("must provide port")
	}

	if c.URL == "" {
		return errors.New("must provide a url")
	}

	paths := sjoin(getPaths(u, network))

	for _, p := range paths {
		err := createIfNotExist(p)
		if err != nil {
			return err
		}
	}

	networkPath, _ := getPaths(u, network)
	configPath := filepath.Join(networkPath, "config.json")

	configData, err := json.Marshal(c)
	err = ioutil.WriteFile(configPath, configData, 0600)
	if err != nil {
		return err
	}

	fmt.Printf("config data has been written to %s\n", configPath)

	return nil
}

func getPaths(u *user.User, network string) (networkPath, vpnUserPath string) {
	vpnPath := filepath.Join(u.HomeDir, ".vpn")
	vpnUserPath = filepath.Join(vpnPath, "username")

	networkPath = filepath.Join(vpnPath, u.Username+"-"+network)
	return
}

func sjoin(s ...string) []string {
	return s
}

func createIfNotExist(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0600)
		}
		return err
	}
	return nil
}
