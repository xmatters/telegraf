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

type Conntrack struct {
	Path  string
	Dir   string
	Files string
}

var dfltPath = "conntrack"
var dfltDirs = []string{"/proc/sys/net/ipv4/netfilter",
	"/proc/sys/net/netfilter"}
var dfltFiles = []string{"ip_conntrack_count,ip_conntrack_max",
	"nf_conntrack_count,nf_conntrack_max"}

func (c *Conntrack) setDefaults() {
	if c.Path == "" {
		c.Path = dfltPath
	}

	if c.Dir == "" {
		c.Dir = dfltDir
	}

	if c.Files == "" {
		c.Files = dfltFiles
	}

}

var dfltProc = "/proc/sys/kernel/random/entropy_avail"

func (c *Conntrack) Description() string {
	// TODO
	return "uses /proc to collect available entropy"
}

func (c *Conntrack) SampleConfig() string {
	// TODO
	return fmt.Sprintf("proc = %s # path to entropy file", dfltProc)
}

func (c *Conntrack) Gather(acc telegraf.Accumulator) error {
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
