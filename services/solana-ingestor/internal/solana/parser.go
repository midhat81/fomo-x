package solana

import "strings"

// ParsedTransaction is the structured result of parsing a raw log event.
type ParsedTransaction struct {
	Signature string
	Success   bool
	Programs  []string
	RawLogs   []string
}

// ParseLogEvent extracts structured info from a raw LogEvent received over
// the websocket subscription.
func ParseLogEvent(evt LogEvent) ParsedTransaction {
	return ParsedTransaction{
		Signature: evt.Signature,
		Success:   evt.Err == nil,
		Programs:  extractProgramIDs(evt.Logs),
		RawLogs:   evt.Logs,
	}
}

// extractProgramIDs pulls program IDs out of "Program <id> invoke [n]" lines.
func extractProgramIDs(logs []string) []string {
	var programs []string
	for _, line := range logs {
		if strings.HasPrefix(line, "Program ") && strings.Contains(line, "invoke") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				programs = append(programs, parts[1])
			}
		}
	}
	return programs
}