package main

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"os/exec"
)

func saveV2rayConfigAndRun(filePath string, confObj *ConfigInfo) {
	if confObj.Index <= 0 || confObj.Index > (len(confObj.Vmess)+1) {
		color.Red("Error index number.")
		return
	}
	activeVm := confObj.Vmess[confObj.Index-1]
	confContent := fmt.Sprintf(`{
  "log": {
    "access": "",
    "error": "",
    "loglevel": ""
  },
  "inbound": {
    "port": %d,
    "listen": "127.0.0.1",
    "protocol": "%s",
    "settings": {
      "auth": "noauth",
      "udp": true,
      "ip": "127.0.0.1",
      "clients": null
    },
    "streamSettings": null
  },
  "outbound": {
    "tag": "agentout",
    "protocol": "vmess",
    "settings": {
      "vnext": [
        {
          "address": "%s",
          "port": %s,
          "users": [
            {
              "id": "%s",
              "alterId": %s,
              "email": "t@t.tt",
              "security": "aes-128-gcm"
            }
          ]
        }
      ],
      "servers": null
    },
    "streamSettings": {
      "network": "tcp",
      "security": "%s",
      "tlsSettings": null,
      "tcpSettings": null,
      "kcpSettings": null,
      "wsSettings": null,
      "httpSettings": null
    },
    "mux": {
      "enabled": true
    }
  },
  "inboundDetour": null,
  "outboundDetour": [
    {
      "protocol": "freedom",
      "settings": {
        "response": null
      },
      "tag": "direct"
    },
    {
      "protocol": "blackhole",
      "settings": {
        "response": {
          "type": "http"
        }
      },
      "tag": "blockout"
    }
  ],
  "dns": {
    "servers": [
      "8.8.8.8",
      "8.8.4.4",
      "localhost"
    ]
  },
  "routing": {
    "strategy": "rules",
    "settings": {
      "domainStrategy": "IPIfNonMatch",
      "rules": [
        {
          "type": "field",
          "port": null,
          "outboundTag": "direct",
          "ip": [
            "0.0.0.0/8",
            "10.0.0.0/8",
            "100.64.0.0/10",
            "127.0.0.0/8",
            "169.254.0.0/16",
            "172.16.0.0/12",
            "192.0.0.0/24",
            "192.0.2.0/24",
            "192.168.0.0/16",
            "198.18.0.0/15",
            "198.51.100.0/24",
            "203.0.113.0/24",
            "::1/128",
            "fc00::/7",
            "fe80::/10"
          ],
          "domain": null
        }
      ]
    }
  }
}`, confObj.LocalPort, confObj.Protocol,
		activeVm.Add, activeVm.Port, activeVm.ID, activeVm.Aid, activeVm.TLS)

	err := ioutil.WriteFile(filePath, []byte(confContent), 0644)
	if err != nil {
		color.Red("Save v2ray config file failed.")
		return
	}
	color.Blue("Run %s at port:%d", v2rayBinaryName, confObj.LocalPort)
	cmd := exec.Command(v2rayBinaryName)
	cmd.Stdout = os.Stdout
	cmd.Start()
}
