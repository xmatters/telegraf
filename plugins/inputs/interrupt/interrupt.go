package entropy

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	inputName = "interrupt"
)

type Interrupt struct {
	Proc        string
	Trace       bool
	prevMetrics map[string][]uint64
	prevTime    time.Time
	readProcFile func (s string) (*bytes.Buffer, error)
}

type Timer interface {
	Now() time.Time
}

const (
	dfltProc = "/proc/interrupts"
)

func (i *Interrupt) Description() string {
	return ""
}

var sampleConfig = `
`

func (i *Interrupt) SampleConfig() string {
	return sampleConfig
}

func (i *Interrupt) setDefaults() {
	if i.Proc == "" {
		i.Proc = dfltProc
	}
}

// Counts the interrupts described on a given line and returns the label along with a per-cpu count
// the returned count array will be <= to cpuCount.  It is up to the caller to handle the case
// where len(counts) < cpuCount.
//
// Counts will be returned in left-to-right order by CPU, eg, CPU #4 will correspond to count[3]
func countInterrupt(line string, cpuCount int) (string, []uint64) {
	cols := strings.Fields(line)
	if len(cols) == 0 {
		return "", []uint64{}
	}
	// Any columns beyond the label and CPU counts should be used as a label prefix
	var buffer bytes.Buffer
	r := strings.NewReplacer(
		",", "_",
		" ", "_")
	for i := cpuCount + 1; i < len(cols); i++ {
		if buffer.Len() > 0 {
			buffer.WriteString(".")
		}
		buffer.WriteString(r.Replace(cols[i]))
	}

	if buffer.Len() > 0 {
		buffer.WriteString(".")
	}
	buffer.WriteString(strings.TrimRight(cols[0], ":"))
	labelPrefix := buffer.String()

	counts := make([]uint64, 0, cpuCount)
	// Note: some rows don't have a value for every CPU (eg ERR and MIS)
	for i := 1; i <= cpuCount && i < len(cols); i++ {
		num, err := strconv.ParseUint(cols[i], 10, 64)
		if err != nil {
			log.Printf("Failed to parse number '%s' as count in line '%s': %v", cols[i], line, err)
		}
		counts = append(counts, num)
	}

	return labelPrefix, counts
}

var readProcFile = func (procFile string) (*bytes.Buffer, error) {
	if _, err := os.Stat(procFile); err != nil {
		return nil, fmt.Errorf("failed to stat file '%s': %v", procFile, err)
	}

	contents, err := ioutil.ReadFile(procFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s': %v", procFile, err)
	}

	return bytes.NewBuffer(contents), nil
}

func (i *Interrupt) Gather(acc telegraf.Accumulator) error {
	i.setDefaults()
	contents, err := i.readProcFile(i.Proc)
	//fmt.Printf("Contents: %v\n", contents)
	now := time.Now()

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(contents)

	ok := scanner.Scan()
	if !ok {
		return fmt.Errorf("failed to count CPUs: %v\n", scanner.Err())
	}
	numCpus := len(strings.Fields(scanner.Text()))

	metrics := make(map[string][]uint64)
	for scanner.Scan() {
		labelPrefix, counts := countInterrupt(scanner.Text(), numCpus)
		metrics[labelPrefix] = counts
	}

	fieldsByCpu := make([]map[string]interface{}, numCpus)
	globalFields := make(map[string]interface{})

	fmt.Printf("numCpus: %d\n", numCpus)
	for i := 0; i < numCpus; i++ {
		fieldsByCpu[i] = make(map[string]interface{})
	}

	fmt.Printf("Fields by CPU: %d", len(fieldsByCpu))

	for label, counts := range metrics {
		if len(counts) == numCpus {
			for cpuIdx, count := range counts {
				fieldsByCpu[cpuIdx][label] = i.derivative(label, cpuIdx, count, now)
			}
		} else if len(counts) < numCpus {
			globalFields[label] = i.derivative(label, 0, counts[0], now)
		}
	}

	acc.AddFields(inputName, globalFields, nil)
	for cpuIdx, fields := range fieldsByCpu {
		acc.AddFields(inputName, fields, map[string]string{"cpu": fmt.Sprintf("cpu%d", cpuIdx)})
	}

	i.prevTime = now
	i.prevMetrics = metrics
	return nil
}

func (i *Interrupt) derivative(label string, cpuIdx int, count uint64, now time.Time) uint64 {
	if i.prevTime.IsZero() || i.prevMetrics == nil {
		return 0
	}

	prevCounts, ok := i.prevMetrics[label]
	if !ok || len(prevCounts) <= cpuIdx {
		return 0
	}

	prev := prevCounts[cpuIdx]
	if prev == 0 {
		return 0
	}

	var yDelta uint64

	if prev <= count {
		yDelta = count - prev
	} else {
		// must have rolled over
		yDelta = count + (math.MaxUint64 - prev)
	}

	xDelta := now.Sub(i.prevTime).Seconds()
	if xDelta == 0 {
		return 0
	}

	return uint64(float64(yDelta) / float64(xDelta))
}

func init() {
	input := &Interrupt{
		readProcFile: readProcFile,
	}
	inputs.Add(inputName, func() telegraf.Input { return input })
}
