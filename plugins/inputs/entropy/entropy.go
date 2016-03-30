package entropy

// simple.go

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type Entropy struct {
	Proc string
}

var dfltProc = "/proc/sys/kernel/random/entropy_avail"

func (e *Entropy) Description() string {
	return "uses /proc to collect available entropy"
}

func (e *Entropy) SampleConfig() string {
	return fmt.Sprintf("proc = %s # path to entropy file", dfltProc)
}

func (e *Entropy) Gather(acc telegraf.Accumulator) error {
	ent := 0
	proc := e.Proc

	if proc == "" {
		proc = dfltProc
	}

	if _, err := os.Stat(proc); err != nil {
		return fmt.Errorf("could not stat proc file '%s': %v", proc, err)
	}

	content, err := ioutil.ReadFile(proc)
	if err != nil {
		return fmt.Errorf("failed to read proc file '%s': %v", proc, err)
	}

	ent, err = strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		return fmt.Errorf("expected integer content but found %s: %v", content, err)
	}
	acc.AddFields("entropy", map[string]interface{}{"available": ent}, nil)
	return nil
}

func init() {
	inputs.Add("entropy", func() telegraf.Input { return &Entropy{} })
}
