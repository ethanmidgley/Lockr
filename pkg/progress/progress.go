package progress

import "fmt"

// Monitor is a fake io reader to monitor the progress of the
type Monitor struct {
	Message         string
	CompleteMessage string
	Iteration       int
}

func (p *Monitor) Write(src []byte) (int, error) {
	n := len(src)

	prog := []string{`|`, `|`, `|`, `|`, `/`, `/`, `/`, `/`, `-`, `-`, `-`, `-`, `\`, `\`, `\`, `\`}

	fmt.Printf("\r\033[2K%s%s", p.Message, prog[p.Iteration%16])
	p.Iteration++

	return n, nil
}

// Complete is a method that can be called once reading has finished
func (p *Monitor) Complete() {
	fmt.Printf("\r\033[2K%s\n", p.CompleteMessage)
}
