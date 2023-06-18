package javalanche

import (
	"fmt"
	"log"
)

const doDebug = false

func (p *Parser) Println(args ...any) {
	if doDebug && len(args) > 0 {
		log.Println(append([]any{"Parser:"}, args...)...)
	}
}

func (p *Parser) Printf(format string, args ...any) {
	if doDebug {
		var msg string
		if len(args) > 0 {
			msg = fmt.Sprintf(format, args...)
		} else {
			msg = format
		}
		log.Println("Parser:", msg)
	}
}

// PrintDetails prints logs of parser
func (p *Parser) PrintDetails(msg string, args ...any) {
	if !doDebug {
		return
	}

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	if len(msg) > 0 {
		log.Println("Parser:", msg)
	}

	// stage
	switch {
	case p.stage.IsEmpty():
		log.Printf("%s:%s: %s",
			"Parser", "Stage", "-empty-")
	default:
		for _, line := range p.stage.Strings() {
			log.Println("Parser:", line)
		}
	}

	// result
	log.Printf("%s %s %#v", "Parser:", "result:", p.result)
}
