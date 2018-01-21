package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// CMD is the child mining process
var CMD *exec.Cmd

// Set some colors up for use by anything
var cTitle = color.New(color.Bold, color.FgBlue)
var cLog = color.New(color.FgGreen)
var cMiner = color.New(color.FgWhite)
var cMediaminer = color.New(color.FgCyan)
var cError = color.New(color.FgRed)

// TomlConfig defines the toml config file
type TomlConfig struct {
	Worker        Worker
	Miningpoolhub Miningpoolhub
	Hodlpool      Hodlpool
	Output        Output
	Wallets       Wallets
}

// Worker reads worker config from toml
type Worker struct {
	Name string
}

// Hodlpool reads hodlpool config settings
type Hodlpool struct {
	UbiqPort   int    `toml:"ubiq_port"`
	UbiqServer string `toml:"ubiq_server"`
}

// Miningpoolhub reads miningpoolhub config from toml
type Miningpoolhub struct {
	Username string
	Ports    Ports
	Servers  Servers
}

// Ports for miningpoolhub
type Ports struct {
	EquihashAuto        int `toml:"equihash_auto"`
	EquihashZencash     int `toml:"equihash_zencash"`
	EquihashZclassic    int `toml:"equihash_zclassic"`
	EquihashBitcoingold int `toml:"equihash_bitcoingold"`
	EquihashZcash       int `toml:"equihash_zcash"`
}

// MphServers for miningpoolhub
type Servers struct {
	Equihash string
}

// Output reads output config from toml
type Output struct {
	Mediaminer bool
	Miner      bool
}

// Wallets
type Wallets struct {
	Ubiq string
}

var config TomlConfig

func main() {
	// Read config file
	_, err := os.Stat("config.toml")
	if err != nil {
		log.Fatal("Can't find config file: config.toml")
	}

	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	// EXAMPLE: Override a template
	cli.AppHelpTemplate = `
	NAME:
		{{.Name}} - {{.Usage}}
	
	Donations:
		BTC:
		LTC:
		ETH:
		BCH:
		DASH:

	USAGE:
		{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
		{{if len .Authors}}
	AUTHOR:
		{{range .Authors}}{{ . }}{{end}}
		{{end}}{{if .Commands}}
	COMMANDS:
	{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
	GLOBAL OPTIONS:
		{{range .VisibleFlags}}{{.}}
		{{end}}{{end}}{{if .Copyright }}
	COPYRIGHT:
		{{.Copyright}}
		{{end}}{{if .Version}}
	VERSION:
		{{.Version}}
		{{end}}
	`

	app := cli.NewApp()
	app.Name = "mediaminer"
	app.Usage = "wrapper for crypto miners"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "miner, m",
			Value: "dstm",
			Usage: "miner to mine with",
		},
		cli.StringFlag{
			Name:  "algorithm, a",
			Value: "equihash",
			Usage: "algorithm to mine with",
		},
		cli.StringFlag{
			Name:  "coin, c",
			Value: "auto",
			Usage: "coin to mine",
		},
		cli.StringFlag{
			Name:  "pool, p",
			Value: "miningpoolhub",
			Usage: "Pool to use",
		},
	}

	app.Action = func(c *cli.Context) error {
		miner := c.String("miner")
		algo := c.String("algorithm")
		coin := c.String("coin")
		pool := c.String("pool")

		cTitle.Println("mediaminer v0.0.1")
		cMediaminer.Println("Trying to start a new miner using:")
		cMediaminer.Println("Pool: ", pool)
		cMediaminer.Println("Miner: ", miner)
		cMediaminer.Println("Algorithm: ", algo)
		cMediaminer.Println("Coin: ", coin)

		username := ""
		worker := config.Worker.Name
		server := ""
		port := 0
		wallet := ""

		// MiningPoolHub
		if pool == "miningpoolhub" {
			server = config.Miningpoolhub.Servers.Equihash
			username = config.Miningpoolhub.Username

			switch coin {
			case "zclassic":
				port = config.Miningpoolhub.Ports.EquihashZclassic
			case "bitcoingold":
				port = config.Miningpoolhub.Ports.EquihashBitcoingold
			case "zcash":
				port = config.Miningpoolhub.Ports.EquihashZcash
			case "zencash":
				port = config.Miningpoolhub.Ports.EquihashZencash
			default:
				port = config.Miningpoolhub.Ports.EquihashAuto
			}
		}

		// Hodlpool
		if pool == "hodlpool" {
			switch coin {
			case "ubiq":
				port = config.Hodlpool.UbiqPort
				server = config.Hodlpool.UbiqServer
				wallet = config.Wallets.Ubiq
			}
		}

		switch miner {
		case "dstm":
			Dstm(server, port, username, worker)
		case "claymoreethdcr":
			Ethdcr(server, port, wallet, worker)
		}
		return nil
	}

	app.Run(os.Args)

	// Test load Dstm
	//Dstm(config.Miningpoolhub.Servers.Equihash, config.Miningpoolhub.Ports.EquihashAuto)
}

// Gets the datetime right now
func now() string {
	return time.Now().Format("02-Jan-2006 15:04:05")
}

// Post is structure for an api call
type Post struct {
	time string
	gpu  string
	temp string
	rate float64
	eff  string
	avg  string
	is   string
}

// Posts the line to an api
func doPost(data Post) {
	_ = data
}

// Outputs a line to the screen
func doMinerOutput(ln string) {
	if config.Output.Miner {
		cMiner.Printf(ln)
	}
}

// Outputs a line to the screen
func doMediaMinerOutput(ln string) {
	if config.Output.Mediaminer {
		cMediaminer.Printf(ln)
	}
}

// Outputs a log line to the screen
func doLog(ln string) {
	cLog.Printf(ln)
}
