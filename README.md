# vpn
Handle vpn's with totp

# usage

```shell
vpn help
start a vpn client connection

Usage:
  vpn [command]

Available Commands:
  configure   creates the necessary configuration files for vpn + totp - ca.crt, user.key, and user.crt must be provided separately
  help        Help about any command
  restart     restarts a vpn on the given network
  start       starts a client connection on the specified network. Aliases "up"
  status      reports the status of the vpn on the given network. Aliases "ls, stat"
  stop        stops a client connection on the specified network. Aliases "down kill"
  totp        generates a totp password for use with a mfa vpn

Flags:
  -h, --help   help for vpn

Use "vpn [command] --help" for more information about a command.
```

## configuring

```shell
vpn help configure

configure will create configuration files for this vpn.

vpn expects user certs/keys to live in $HOMEDIR/.vpn/$VPNUSERNAME-$NETWORK, so /Users/ccooper/.vpn/ccooper-iad or $HOMEDIR/.vpn/ccooper-default

this command will create a config.json file in each network directory.

configure will create these directories if they do not exist, but it is the user's responsibility to place the correct ca.crt, username.key, and username.crt in each directory

to be clear, username.key and username.crt should be your actual vpn username, so in my case they would be ccooper.key and ccooper.crt

ca.crt should be called ca.crt exactly

Usage:
  vpn configure [flags] [network]

Flags:
  -h, --help                help for configure
  -m, --mfa-secret string   the mfa secret key to generate one time passwords
  -p, --port int            the port for this remote (default -1)
  -u, --url string          the url for this remote
```

## start
Note: Currently start simply starts the VPN daemon and exits - in the future, I intend to add
support for watching the vpn daemon and restarting as necessary. That would probably require running
the entire process as root initially, so some additional configuration would be necessary.

```shell
vpn help start
starts a client connection on the specified network. Aliases "up"

Usage:
  vpn start [flags] [network]

Aliases:
  start, up

Flags:
  -h, --help   help for start
```

## stop

```shell
stops a client connection on the specified network. Aliases "down, kill"

Usage:
  vpn stop [flags] [network]

Aliases:
  stop, down, kill

Flags:
  -h, --help   help for stop
```

## status

```shell
reports the status of the vpn on the given network. Aliases "ls, stat"

Usage:
  vpn status [network] [flags]

Aliases:
  status, ls, stat

Flags:
  -h, --help   help for status
```

## restart

```shell
vpn help restart
restarts a vpn on the given network

Usage:
  vpn restart [network] [flags]

Flags:
  -h, --help   help for restart
```

## totp
Note: If you would just like totp command line support, that is also provided.

```shell
vpn help totp
generates a totp password for use with a mfa vpn (or any mfa)

Usage:
  vpn totp [flags]

Flags:
  -h, --help   help for totp
```