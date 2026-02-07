package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	lcdDev  = "/dev/ttyS1"
	refresh = 5 * time.Second
)

// cpu

func getCPU() int {
	data, _ := os.ReadFile("/proc/stat")
	fields := strings.Fields(strings.Split(string(data), "\n")[0])

	var total, idle int
	for i := 1; i < len(fields); i++ {
		val, _ := strconv.Atoi(fields[i])
		total += val
		if i == 4 {
			idle = val
		}
	}

	time.Sleep(120 * time.Millisecond)

	data2, _ := os.ReadFile("/proc/stat")
	fields2 := strings.Fields(strings.Split(string(data2), "\n")[0])

	var total2, idle2 int
	for i := 1; i < len(fields2); i++ {
		val, _ := strconv.Atoi(fields2[i])
		total2 += val
		if i == 4 {
			idle2 = val
		}
	}

	return 100 * (total2 - total - (idle2 - idle)) / (total2 - total)
}

// ram

func getRAM() int {
	f, _ := os.Open("/proc/meminfo")
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var total, free int

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal") {
			fmt.Sscanf(line, "MemTotal: %d kB", &total)
		}
		if strings.HasPrefix(line, "MemAvailable") {
			fmt.Sscanf(line, "MemAvailable: %d kB", &free)
		}
	}

	return 100 * (total - free) / total
}

// disk

func getDisk() int {
	out, _ := exec.Command("df", "/").Output()
	lines := strings.Split(string(out), "\n")
	fields := strings.Fields(lines[1])
	percent, _ := strconv.Atoi(strings.TrimSuffix(fields[4], "%"))
	return percent
}

// temp

func getTemp() (string, bool) {
	data, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return "", false
	}
	val, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return fmt.Sprintf("%dC", val/1000), true
}

// gateway 

func getGateway() string {
	out, err := exec.Command("ip", "route").Output()
	if err != nil {
		return "N/A"
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "default") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				return fields[2]
			}
		}
	}
	return "N/A"
}

// uptime

func getUptime() string {
	data, _ := os.ReadFile("/proc/uptime")
	secs, _ := strconv.ParseFloat(strings.Fields(string(data))[0], 64)

	d := int(secs) / 86400
	h := (int(secs) % 86400) / 3600
	m := (int(secs) % 3600) / 60

	return fmt.Sprintf("%dd%02dh%02dm", d, h, m)
}

// dhcp leases

func getDHCPLeases() int {
	data, err := os.ReadFile("/var/lib/misc/dnsmasq.leases")
	if err != nil {
		return 0
	}
	lines := strings.Split(string(data), "\n")
	count := 0
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			count++
		}
	}
	return count
}



func StartLCDStatus(stop <-chan struct{}) {

	fd, err := os.OpenFile(lcdDev, os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer fd.Close()

	view := 0

	for {
		select {
		case <-stop:
			return
		default:
		}

        // reset de lcd
		fmt.Fprint(fd, "\r\r")
		fmt.Fprint(fd, "\r\n\r\n\r\n\r\n\r\n")

		line0 := "                " // 16 chars

		if view%2 == 0 {

			line1 := fmt.Sprintf("CPU:%02d%%", getCPU())
			line2 := fmt.Sprintf("RAM:%02d%%", getRAM())
			line3 := fmt.Sprintf("DISK:%02d%%", getDisk())
			line4 := "SENTINELOS"

			fmt.Fprintf(fd, "%-16s\n", line0)
			fmt.Fprintf(fd, "%-16s\n", line1)
			fmt.Fprintf(fd, "%-16s\n", line2)
			fmt.Fprintf(fd, "%-16s\n", line3)
			fmt.Fprintf(fd, "%-16s\n\n", line4)

		} else {

			temp, ok := getTemp()
			var line1 string
			if ok {
				line1 = "TEMP:" + temp
			} else {
				line1 = fmt.Sprintf("DHCP:%d", getDHCPLeases())
			}

			line2 := "GW:" + getGateway()
			line3 := "UP:" + getUptime()
			line4 := "SENTINELOS"

			fmt.Fprintf(fd, "%-16s\n", line0)
			fmt.Fprintf(fd, "%-16s\n", line1)
			fmt.Fprintf(fd, "%-16s\n", line2)
			fmt.Fprintf(fd, "%-16s\n", line3)
			fmt.Fprintf(fd, "%-16s\n\n", line4)
		}

		view++
		time.Sleep(refresh)
	}
}
