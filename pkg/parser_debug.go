package javalanche

import "log"

const doDebug = false

// PrintDetails prints logs of parser
func (p *Parser) PrintDetails() {
	if !doDebug {
		return
	}

	// tokens
	l := len(p.tokens)
	switch l {
	case 0:
		log.Printf("%s:%s: %s",
			"Parser", "tokens", "-empty-")
	default:
		for i, t := range p.tokens {
			log.Printf("%s:%s: [%v/%v] %#v",
				"Parser", "tokens", i, l, t)
		}
	}

	// result
	log.Printf("%s:%s: %#v", "Parser", "result", p.result)
}
