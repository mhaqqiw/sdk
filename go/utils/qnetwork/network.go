package qnetwork

import (
	"errors"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/mhaqqiw/sdk/go/entity"
)

var descMap = map[string]map[string]interface{}{
	"lo": {
		"desc": "Loopback Interface: The lo0 interface, also known as the loopback interface, is a virtual network interface used to establish network connections with the local host (i.e., the same device)",
		"type": 1,
	},
	"gif": {
		"desc": "Generic Tunnel Interface: The gif0 interface is typically used for creating IPv6 over IPv4 tunnels or other generic tunneling purposes. It enables communication between devices or networks that use different network protocols",
		"type": 2,
	},
	"stf": {
		"desc": "6to4 Interface: The stf0 interface, also known as the 6to4 interface, is used for automatic tunneling of IPv6 traffic over an IPv4 network",
		"type": 3,
	},
	"utun": {
		"desc": "User Tunnel Interface: The utun0 interface is often used for VPN (Virtual Private Network) connections, and it's associated with user-space applications",
		"type": 4,
	},
	"llw": {
		"desc": "LLW (Low Latency Wi-Fi) is a technology introduced in macOS for more responsive and low-latency wireless communication",
		"type": 5,
	},
	"anpi": {
		"desc": "ANPI (Advanced Networking Packet Interface) is a technology introduced in macOS for high-speed wireless communication",
		"type": 6,
	},
	"bridge": {
		"desc": "A bridge interface is used to connect multiple network interfaces together, allowing traffic to pass between them",
		"type": 7,
	},
	"ap": {
		"desc": "An access point is a device that connects to a wireless network",
		"type": 8,
	},
	"awdl": {
		"desc": "This interface is associated with Apple Wireless Direct Link (AWDL), a technology used for peer-to-peer communication between Apple devices, such as AirDrop and AirPlay",
		"type": 9,
	},
	"en": {
		"desc": "This interface is associated with Ethernet, a technology used for network communication",
		"type": 10,
	},
}

func splitAlphaNumeric(input string) (string, int, error) {
	// Use regular expressions to separate alphabetic and numeric parts
	re := regexp.MustCompile(`([a-zA-Z]+)([0-9]+)`)
	parts := re.FindStringSubmatch(input)
	if len(parts) == 3 {
		num, err := strconv.Atoi(parts[2])
		if err != nil {
			return "", 0, err
		}
		return parts[1], num, nil
	}
	return "", 0, errors.New("invalid format")
}

func GetNetworkInfo() ([]entity.QHardwareInterface, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []entity.QHardwareInterface
	for _, ifa := range ifas {
		part1, part2, err := splitAlphaNumeric(ifa.Name)
		if err != nil {
			return nil, err
		}
		var addrs []entity.QAddress
		address, err := ifa.Addrs()
		if err != nil {
			return nil, err
		}
		for _, a := range address {
			addr := entity.QAddress{
				IP:   a.(*net.IPNet).IP.String(),
				Mask: a.(*net.IPNet).Mask.String(),
			}
			if a.(*net.IPNet).IP.To4() != nil {
				addr.IsV4 = true
			}
			addrs = append(addrs, addr)
		}
		as = append(as, entity.QHardwareInterface{
			Addres:          addrs,
			Code:            part1,
			Desc:            descMap[part1]["desc"].(string),
			InterfaceNumber: part2,
			Type:            descMap[part1]["type"].(int),
			Index:           ifa.Index,
			MTU:             ifa.MTU,
			Name:            ifa.Name,
			HardwareAddr:    strings.ToUpper(ifa.HardwareAddr.String()),
			Flags:           ifa.Flags.String(),
		})
	}
	return as, nil
}
