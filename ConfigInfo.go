package main

type ConfigInfo struct {
	LocalPort int         `json:"localPort"`
	Protocol  string      `json:"protocol"`
	Index     int         `json:"index"`
	Vmess     []VmessInfo `json:"vmess"`
	SubURL    string      `json:"subUrl"`
}

type VmessInfo struct {
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	ID   string `json:"id"`
	Aid  string `json:"aid"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Host string `json:"host"`
	TLS  string `json:"tls"`
}
