package main

import (
	"../logger"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"io/ioutil"
	"os"
)

var Logger *log.Logger

type Args struct {
	A int `msg:"a"`
	B int `msg:"b"`
}
type Reply struct {
	C int `msg:"c"`
}

type Config struct {
	LogPath string `json:"log_path"` //各级别日志路径
	Address string `json:"server_address"`
}

func loadConfig(path string) *Config {
	if len(path) == 0 {
		panic("path of conifg is null.")
	}
	_, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	var cfg Config
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func main() {
	var cfg_path string
	flag.StringVar(&cfg_path, "conf", "../conf/conf.json", "config file path")
	flag.Parse()
	fmt.Println(cfg_path)

	cfg := loadConfig(cfg_path)
	l := logger.GetLogger(cfg.LogPath, "client")
	Logger = l
	Logger.Infof("rpcx client start.%v", cfg)

	d := client.NewPeer2PeerDiscovery("tcp@" + cfg.Address, "")
	xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &Args{
		A: 10,
		B: 20,
	}

	reply := &Reply{}
	err := xclient.Call(context.Background(), "Mul", args, reply)
	if err != nil {
		Logger.Errorf("failed to call: %v", err)
	}

	Logger.Infof("%d * %d = %d", args.A, args.B, reply.C)
}
