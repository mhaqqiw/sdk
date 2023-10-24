package entity

type QHardwareInterface struct {
	Index           int        `json:"index"`
	MTU             int        `json:"mtu"`
	Name            string     `json:"name"`
	HardwareAddr    string     `json:"hardware_addr"`
	Flags           string     `json:"flags"`
	Type            int        `json:"type"`
	Desc            string     `json:"desc"`
	InterfaceNumber int        `json:"interface_number"`
	Code            string     `json:"code"`
	Addres          []QAddress `json:"addr"`
}

type QAddress struct {
	IP   string `json:"ip"`
	Mask string `json:"mask"`
	IsV4 bool   `json:"is_v4"`
}
