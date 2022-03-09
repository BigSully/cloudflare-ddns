package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

func IsIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

func IsIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

func PrintInterfaceAddr(interfaceName string) {
	var (
		err   error
		ief   *net.Interface
		addrs []net.Addr
	)
	if ief, err = net.InterfaceByName(interfaceName); err != nil { // get interface
		return
	}
	if addrs, err = ief.Addrs(); err != nil { // get addresses
		return
	}
	for _, addr := range addrs { // get ipv4 address
		ip := addr.(*net.IPNet).IP
		ipStr := ip.String()
		if ip.IsLinkLocalUnicast() {
			fmt.Printf("%v IsLinkLocalUnicast() \n", ipStr)
			continue
		}
		if ip.IsPrivate() {
			fmt.Printf("%v IsPrivate() \n", ipStr)
			continue
		}
		fmt.Printf("ip: %s, isipv6: %v\n", ipStr, IsIPv6(ipStr))
	}
}

// useful links:
// https://stackoverflow.com/questions/27410764/dial-with-a-specific-address-interface-golang
// https://stackoverflow.com/questions/22751035/golang-distinguish-ipv4-ipv6
func GetInterfaceIpv4Addr(interfaceName string) (addr string, err error) {
	var (
		ief      *net.Interface
		addrs    []net.Addr
		ipv4Addr net.IP
	)
	if ief, err = net.InterfaceByName(interfaceName); err != nil { // get interface
		return
	}
	if addrs, err = ief.Addrs(); err != nil { // get addresses
		return
	}
	for _, addr := range addrs { // get ipv4 address
		if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			break
		}
	}
	if ipv4Addr == nil || ipv4Addr.String() == "" {
		return "", errors.New(fmt.Sprintf("interface %s don't have an ipv4 address\n", interfaceName))
	}
	return ipv4Addr.String(), nil
}

func Getenv(key, def string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return def
	}
	return value
}

func GetPublicIP() (ip string, err error) {
	resp, err := http.Get("https://1.1.1.1/cdn-cgi/trace")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	s := string(body)

	// map the string to a map
	entries := strings.Split(s, "\n")

	m := make(map[string]string)
	for _, e := range entries {
		parts := strings.Split(e, "=")
		key := strings.ToLower(parts[0])
		if len(parts) > 1 {
			m[key] = strings.TrimSpace(parts[1])
		} else {
			m[key] = ""
		}
	}

	return m["ip"], nil
}

func GetPublicIPv6() (ip string, err error) {
	resp, err := http.Get("http://icanhazip.com")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	s := string(body)
	ipv6 := strings.TrimSpace(s)

	return ipv6, nil
}
