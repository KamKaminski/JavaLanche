package javalanche

import (
	"fmt"
	"log"
	"strings"
)

const debugStage = true

// Println provides custom logging
func (s *Stage) Println(args ...any) {
	if debugStage && len(args) > 0 {
		log.Println(append([]any{"Stage:"}, args...)...)
	}
}

// Printf provides with custom logging
func (s *Stage) Printf(format string, args ...any) {
	if debugStage {
		var msg string
		if len(args) > 0 {
			msg = fmt.Sprintf(format, args...)
		} else {
			msg = format
		}
		log.Println("Stage:", msg)
	}
}

// PrintDetails print details of the Stage if Stage debug is enabled
func (s *Stage) PrintDetails(format string, args ...any) {
	if debugStage {
		var msg string
		if len(args) > 0 {
			msg = fmt.Sprintf(format, args...)
		} else {
			msg = format
		}

		if len(msg) > 0 {
			log.Println("Stage:", msg)
		}

		for _, line := range s.Strings() {
			log.Println("Stage:", line)
		}
	}
}

// String joins strings
func (s Stage) String() string {
	return strings.Join(s.Strings(), "\n")
}

// Strings  returns a string representation of each token and node in its pairs field
func (s Stage) Strings() []string {
	var i int

	total := s.Len()
	lines := make([]string, 0, total)

	for _, n := range s.nodes {
		line := fmt.Sprintf("[%v/%v]: %s", i, total, n.Any())
		lines = append(lines, line)
		i++
	}

	return lines
}
