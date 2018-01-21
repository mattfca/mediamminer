package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"strconv"
)

// Dstm Runs DSTM Miner
func Dstm(server string, port int, username string, worker string) {

	// Set regex values for parsing dstm miner output
	var dstmx = regexp.MustCompile(`GPU(?s).*Sol\/s(?s).*Sol\/W(?s).*Avg(?s).*I\/s(?s).*`)
	var gpux = regexp.MustCompile(`(GPU\d)`)
	var tempx = regexp.MustCompile(`(\d*)C`)
	var ratex = regexp.MustCompile(`Sol\/s:\ (\S*)`)
	var effx = regexp.MustCompile(`Sol\/W:\ (\S*)`)
	var avgratex = regexp.MustCompile(`Avg:\ (\S*)`)
	var isx = regexp.MustCompile(`I\/s:\ (\S*)`)

	// Configure arguments to run with miner
	args := fmt.Sprintf("--dev --server %s --port %d --user %s.%s --pass x", server, port, username, worker)
	doLog(fmt.Sprintf("Starting DSTM Miner with %s\n", args))

	// Start the miner
	CMD = exec.Command("miners/zm_0.5.7_win/zm.exe", strings.Split(args, " ")...)

	stdout, err := CMD.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	CMD.Start()

	scanner := bufio.NewScanner(stdout)

	badRate := 0

	for scanner.Scan() {
		ln := scanner.Text()
		doMinerOutput(fmt.Sprintf("%s\n", ln))

		if dstmx.MatchString(ln) {
			gpu := gpux.FindStringSubmatch(ln)[1]
			temp := tempx.FindStringSubmatch(ln)[1]
			rate, err := strconv.ParseFloat(ratex.FindStringSubmatch(ln)[1], 64)
			if err != nil {
				log.Fatal(err)
			}
			eff := effx.FindStringSubmatch(ln)[1]
			avgrate := avgratex.FindStringSubmatch(ln)[1]
			is := isx.FindStringSubmatch(ln)[1]

			doMediaMinerOutput(fmt.Sprintf("%s - %s Temp: %s Sol/s: %f Sol/W: %s Avg: %s I/s: %s\n", now(), gpu, temp, rate, eff, avgrate, is))

			if rate == 0 {
				if badRate >= 3 {
					_ = CMD.Process.Kill()
					Dstm(server, port, username, worker)
				}
				badRate++
			}

			doPost(Post{
				time: now(),
				gpu:  gpu,
				temp: temp,
				rate: rate,
				eff:  eff,
				avg:  avgrate,
				is:   is,
			})
		}
	}

	CMD.Wait()
}
