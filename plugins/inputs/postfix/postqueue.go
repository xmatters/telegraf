package system

import (
	"bytes"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type Postqueue struct {
	getPostQueueLog func() ([]byte, error)
}

func (k *Postqueue) Description() string {
	return "Get the number of mails in the Postfix queue."
}

func (k *Postqueue) SampleConfig() string { return "" }

// We are interested in the last line from the output of "postqueue -p"
// An example of the last line is "-- 2106 Kbytes in 205 Requests." and we are
// interested in extracting the "205" number.
func (k *Postqueue) Gather(acc telegraf.Accumulator) error {

	out, err := k.getPostQueueLog()
	if err != nil {
		return err
	}

	currentLine := ""
	previousLine := ""
	fields := make(map[string]interface{})

	buffer := bytes.NewBuffer(out)

	// This loop is to find the last line.
	for {
		currentLine, err = buffer.ReadString('\n')
		if err == io.EOF {
			// Break out when we reach the end-of-file
			break
		} else {
			previousLine = currentLine
		}
	}

	// Get the last line. In my own tests, it is possible that current line
	// is blank when EOF is reached then we have to check the previous line.
	lastLine := ""
	if len(previousLine) > 0 && strings.HasPrefix(previousLine, "--") {
		lastLine = previousLine
	} else if len(currentLine) > 0 && strings.HasPrefix(currentLine, "--") {
		lastLine = currentLine
	}

	// Split the last line into word tokens and extract the
	// second last word.
	token := strings.Fields(lastLine)

	var count int
	// Don't try parsing if the output is obviously wrong, since it most likely indicates
	// nothing in the queue
	if !strings.Contains(strings.ToLower(lastLine), "empty") && len(token) >= 2 {
		count, err = strconv.Atoi(token[len(token)-2])
		if err != nil {
			return fmt.Errorf("failed to parse count from last line of output: %s, %v", lastLine, err)
		}
	}

	fields[string("total_count")] = count

	acc.AddFields("postqueue", fields, map[string]string{})
	return nil
}

func init() {
	inputs.Add("postqueue", func() telegraf.Input {
		return &Postqueue{
			getPostQueueLog: func() ([]byte, error) {
				out, err := exec.Command("postqueue", "-p").Output()
				if err != nil {
					return nil, err
				}
				return out, nil
			},
		}
	})
}
