package vpn

/*

ported from this function, which was so generously provided by Che Ruisi-Besares: https://github.com/CheRuisiBesares

vpn ()
{
    local target=${2:-lga08};
    case $1 in
        up | start)
            ( echo $USER > ~/.vpn/username;
            chmod 0600 ~/.vpn/username;
            cd ~/.vpn/$USER-${target}/ && sudo openvpn --daemon "openvpn ${target}" --remote bastion.${target}.nsone.co 1194 --ca ca.crt --key $USER.key --cert $USER.crt --proto udp --resolv-retry infinite --nobind --user nobody --group nogroup --persist-key --persist-tun --max-routes 1000 --verb 3 --comp-lzo --dev tun --client --remote-cert-tls server --auth-user-pass ~/.vpn/username --auth-nocache --reneg-sec 0 )
        ;;
        down | stop)
            ps aux | grep openvpn | grep "${target}" | awk '{print $2}' | xargs sudo kill -SIGTERM
        ;;
        restart)
            vpn down $target && vpn up $target
        ;;
        ls | status)
            ps aux | grep openvpn | grep -v grep
        ;;
        *)
            echo "must be up or down" && return 1
        ;;
    esac
}
*/
