package log

import (
	"log"
	"net"
	"regexp"
	"strings"
	"x-ui-monitor/internal/usecase"

	"github.com/hpcloud/tail"
)

// func TailLogFile(filePath string, ) {
// 	t, err := tail.TailFile(filePath, tail.Config{
// 		Follow: true,
// 		ReOpen: true, // 🔥 survive log rotation
// 		Logger: tail.DiscardingLogger,
// 	})
// 	if err != nil {
// 		log.Fatalf("failed to tail log file: %v", err)
// 	}

// 	ipv4Regex := regexp.MustCompile(`from (\d+\.\d+\.\d+\.\d+):`)
// 	ipv6Regex := regexp.MustCompile(`from \[([a-fA-F0-9:]+)\]:`)

// 	for line := range t.Lines {
// 		text := line.Text

// 		if matches := ipv4Regex.FindStringSubmatch(text); len(matches) > 1 {
// 			userUsecase.AddUser(matches[1])
// 		} else if matches := ipv6Regex.FindStringSubmatch(text); len(matches) > 1 {
// 			userUsecase.AddUser(matches[1])
// 		}
// 	}
// }

func TailLogFile(filePath string, userUsecase *usecase.UserUsecase) {
	log.Println("👀 Start watching log file ...")
	t, err := tail.TailFile(filePath, tail.Config{
		Follow: true,
		ReOpen: true, // <- THIS is key to surviving rotations!
		Logger: tail.DiscardingLogger,
	})

	if err != nil {
		log.Fatalf("Failed to tail file: %v", err)
	}

	for line := range t.Lines {
		processLogLine(line.Text, userUsecase)
	}
}

func processLogLine(line string, userUsecase *usecase.UserUsecase) {
	log.Printf(line)
	fromRegex := regexp.MustCompile(`from (\[?[a-fA-F0-9:.]+\]?):\d+ accepted .* \[(.*?)\]`)
	matches := fromRegex.FindStringSubmatch(line)
	if len(matches) < 3 {
		return
	}

	clientIP := strings.Trim(matches[1], "[]")
	inboundInfo := matches[2]

	inboundParts := strings.Split(inboundInfo, " ")
	if len(inboundParts) < 1 {
		return
	}
	inboundTag := strings.Split(inboundParts[0], ">>")[0]
	inboundTag = strings.TrimSpace(inboundTag)

	if isLocalOrPrivateIP(clientIP) {
		return
	}

	log.Printf("✔️ Accepted IP: %s on inbound %s\n", clientIP, inboundTag)
	userUsecase.AddUser(inboundTag, clientIP)
}

func isLocalOrPrivateIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return true
	}
	if parsedIP.IsLoopback() {
		return true
	}
	privateBlocks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"fc00::/7",
	}
	for _, block := range privateBlocks {
		_, subnet, _ := net.ParseCIDR(block)
		if subnet.Contains(parsedIP) {
			return true
		}
	}
	return false
}
