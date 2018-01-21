package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// Ethdcr is Claymore EthDcrMiner64
func Ethdcr(server string, port int, wallet string, worker string) {
	// Set regex

	// Configure arguments to run with miner
	args := fmt.Sprintf(" -epool %s:%d -ewal %s -eworker %s -epsw x -esm 0 -allcoins 1", server, port, wallet, worker)
	doLog(fmt.Sprintf("Starting EthDcrMiner64 Miner with %s\n", args))

	// Start the miner
	CMD = exec.Command("miners/Claymore_Dual_Ethereum_v10.2/EthDcrMiner64.exe", strings.Split(args, " ")...)

	stdout, err := CMD.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	CMD.Start()

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		ln := scanner.Text()
		doMinerOutput(fmt.Sprintf("%s\n", ln))
	}

	CMD.Wait()
}
