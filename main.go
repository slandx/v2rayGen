package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
)

//
const GVersionString = "v0.0.1"
const GConfigFilePath = "GConfig.json"
const V2rayConfigFilePath = "config.json"
const V2rayBinWin32 = "v2ray.exe"
const V2rayBinLinuxOrMac = "v2ray"

var v2rayBinaryName string

func main() {
	if runtime.GOOS == "windows" {
		v2rayBinaryName = V2rayBinWin32
	} else {
		v2rayBinaryName = V2rayBinLinuxOrMac
	}

	var subUrl string
	var addVmess string
	var port int
	var showVersion bool
	flag.StringVar(&subUrl, "s", "", "Subscribe URL")
	flag.StringVar(&addVmess, "a", "", "String starts with 'vmess://...'")
	flag.IntVar(&port, "p", 1080, "Local port")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.Parse()

	eColor := color.New(color.FgRed)

	if showVersion {
		fmt.Println(GVersionString)
		return
	}

	if _, err := os.Stat(v2rayBinaryName); os.IsNotExist(err) {
		eColor.Println("v2ray core not found\nDownload and extract to this folder from https://github.com/v2ray/v2ray-core/releases/latest")
		return
	}

	if _, err := os.Stat(GConfigFilePath); os.IsNotExist(err) && len(subUrl) == 0 && len(addVmess) == 0 {
		eColor.Println("v2rayGen Config file is not exist. Use -s or -a to add a config file first.")
		flag.Usage()
		return
	}

	var confObj ConfigInfo
	confObj, err := readConfig(GConfigFilePath)
	if err != nil {
		//init config
		confObj.SubURL = subUrl
		confObj.LocalPort = port
		confObj.Protocol = "socks"
		confObj.Index = -1
	}

	if len(subUrl) != 0 {
		err := updateBySubscribeUrl(&confObj, subUrl)
		if err != nil {
			eColor.Println(err)
			return
		}
	}
	if len(addVmess) != 0 {
		if len(confObj.Vmess) == 0 {
			vmObj, err := parseVmessUrl(addVmess)
			if err != nil {
				eColor.Println("Invalid vmess URL, skip.")
			} else {
				confObj.Vmess = append(confObj.Vmess, vmObj)
			}
		}
	}

	generateV2rayConfig(&confObj)

}

func readConfig(confPath string) (ConfigInfo, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	confBytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		return ConfigInfo{}, errors.New("Read config file failed")
	}
	var confObj ConfigInfo
	err = json.Unmarshal(confBytes, &confObj)
	if err != nil {
		return ConfigInfo{}, errors.New("Parse config file failed")
	}
	return confObj, nil
}

func updateBySubscribeUrl(conf *ConfigInfo, subUrl string) error {
	//fmt.Println(subUrl)
	eColor := color.New(color.FgRed)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	resp, err := http.Get(subUrl)
	if err != nil {
		return errors.New("Cannot get subscribe URL")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Cannot get content of URL")
	}

	decodeBytes, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return errors.New("Decode failed with subscribe URL content")
	}
	vms := strings.Split(string(decodeBytes), "\n")
	conf.Vmess = make([]VmessInfo, len(vms))
	for idx, num := range vms {
		if len(num) > 8 {
			vmObj, err := parseVmessUrl(num)
			if err != nil {
				eColor.Println(err)
				continue
			}
			conf.Vmess[idx] = vmObj
		}
	}
	configJsonString, _ := json.Marshal(conf)
	err = ioutil.WriteFile(GConfigFilePath, configJsonString, 0644)
	if err != nil {
		return errors.New("Save config file Failed")
	}
	return nil
}
func parseVmessUrl(vmessUrl string) (VmessInfo, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	//fmt.Println(vmessUrl)
	if len(vmessUrl) <= 0 {
		return VmessInfo{}, errors.New("Empty vmess URL string")
	}
	if !strings.HasPrefix(vmessUrl, "vmess://") {
		return VmessInfo{}, errors.New("Invalid vmess URL string")
	}
	vmDec, err := base64.StdEncoding.DecodeString(vmessUrl[8:])
	if err != nil && len(vmDec) > 0 {
		return VmessInfo{}, errors.New("Decode failed with vmess URL content")
	}
	//fmt.Println(string(vmDec))
	var vmObj VmessInfo
	err = json.Unmarshal(vmDec, &vmObj)
	if err != nil {
		return VmessInfo{}, errors.New("Decode json failed")
	}
	return vmObj, nil
}

func generateV2rayConfig(confObj *ConfigInfo) {
	eColor := color.New(color.FgRed)
	gColor := color.New(color.FgHiGreen)
	wColor := color.New(color.FgWhite)
	if len(confObj.Vmess) == 0 {
		eColor.Println("No vmess config is found.")
		return
	}
	for true {
		for idx, vm := range confObj.Vmess {
			if idx == confObj.Index {
				gColor.Printf("%d. %s (active)\n", idx+1, vm.Ps)
			} else {
				wColor.Printf("%d. %s\n", idx+1, vm.Ps)
			}
		}
		wColor.Printf("%d. 退出\n", len(confObj.Vmess)+1)
		wColor.Print("Select active server:")
		var activeIdx int
		_, err := fmt.Scan(&activeIdx)
		if err != nil || activeIdx < 1 || activeIdx > (len(confObj.Vmess)+1) {
			eColor.Println("Error number")
			continue
		}
		if activeIdx == (len(confObj.Vmess) + 1) {
			return
		}
		confObj.Index = activeIdx
		saveV2rayConfigAndRun(V2rayConfigFilePath, confObj)
		break
	}
}
