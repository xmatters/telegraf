// +build linux

package sysstat

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

var (
	firstTimestamp time.Time
	execCommand    = exec.Command // execCommand is used to mock commands in tests.
	dfltActivities = []string{"DISK"}
)

const parseInterval = 1 // parseInterval is the interval (in seconds) where the parsing of the binary file takes place.

type Sysstat struct {
	// Sadc represents the path to the sadc collector utility.
	Sadc string `toml:"sadc_path"`

	// Sadf represents the path to the sadf cmd.
	Sadf string `toml:"sadf_path"`

	// Activities is a list of activities that are passed as argument to the
	// collector utility (e.g: DISK, SNMP etc...)
	// The more activities that are added, the more data is collected.
	Activities []string

	// Options is a map of options.
	//
	// The key represents the actual option that the Sadf command is called with and
	// the value represents the description for that option.
	//
	// For example, if you have the following options map:
	//    map[string]string{"-C": "cpu", "-d": "disk"}
	// The Sadf command is run with the options -C and -d to extract cpu and
	// disk metrics from the collected binary file.
	//
	// If Group is false (see below), each metric will be prefixed with the corresponding description
	// and represents itself a measurement.
	//
	// If Group is true, metrics are grouped to a single measurement with the corresponding description as name.
	Options map[string]string

	// Group determines if metrics are grouped or not.
	Group bool

	// DeviceTags adds the possibility to add additional tags for devices.
	DeviceTags map[string][]map[string]string `toml:"device_tags"`
	tmpFile    string
	interval   int
}

func (*Sysstat) Description() string {
	return "Sysstat metrics collector"
}

var sampleConfig = `
  ## Path to the sadc command.
  #
  ## Common Defaults:
  ##   Debian/Ubuntu: /usr/lib/sysstat/sadc
  ##   Arch:          /usr/lib/sa/sadc
  ##   RHEL/CentOS:   /usr/lib64/sa/sadc
  sadc_path = "/usr/lib/sa/sadc" # required
  #
  #
  ## Path to the sadf command, if it is not in PATH
  # sadf_path = "/usr/bin/sadf"
  #
  #
  ## Activities is a list of activities, that are passed as argument to the
  ## sadc collector utility (e.g: DISK, SNMP etc...)
  ## The more activities that are added, the more data is collected.
  # activities = ["DISK"]
  #
  #
  ## Group metrics to measurements.
  ##
  ## If group is false each metric will be prefixed with a description
  ## and represents itself a measurement.
  ##
  ## If Group is true, corresponding metrics are grouped to a single measurement.
  # group = true
  #
  #
  ## Options for the sadf command. The values on the left represent the sadf options and
  ## the values on the right their description (wich are used for grouping and prefixing metrics).
  ##
  ## Run 'sar -h' or 'man sar' to find out the supported options for your sysstat version.
  [inputs.sysstat.options]
	-C = "cpu"
	-B = "paging"
	-b = "io"
	-d = "disk"             # requires DISK activity
	"-n ALL" = "network"
	"-P ALL" = "per_cpu"
	-q = "queue"
	-R = "mem"
	-r = "mem_util"
	-S = "swap_util"
	-u = "cpu_util"
	-v = "inode"
	-W = "swap"
	-w = "task"
  #	-H = "hugepages"        # only available for newer linux distributions
  #	"-I ALL" = "interrupts" # requires INT activity
  #
  #
  ## Device tags can be used to add additional tags for devices. For example the configuration below
  ## adds a tag vg with value rootvg for all metrics with sda devices.
  # [[inputs.sysstat.device_tags.sda]]
  #  vg = "rootvg"
`

func (*Sysstat) SampleConfig() string {
	return sampleConfig
}

func (s *Sysstat) Gather(acc telegraf.Accumulator) error {
	if s.interval == 0 {
		if firstTimestamp.IsZero() {
			firstTimestamp = time.Now()
		} else {
			s.interval = int(time.Since(firstTimestamp).Seconds())
		}
	}
	ts := time.Now().Add(time.Duration(s.interval) * time.Second)
	if err := s.collect(); err != nil {
		return err
	}
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(s.Options)*2)
	for option := range s.Options {
		wg.Add(1)
		go func(acc telegraf.Accumulator, option string) {
			defer wg.Done()
			if err := s.parse(acc, option, ts); err != nil {
				errorChannel <- err
			}
		}(acc, option)
	}
	wg.Wait()
	close(errorChannel)

	errorStrings := []string{}
	for err := range errorChannel {
		errorStrings = append(errorStrings, err.Error())
	}

	if _, err := os.Stat(s.tmpFile); err == nil {
		if err := os.Remove(s.tmpFile); err != nil {
			errorStrings = append(errorStrings, err.Error())
		}
	}

	if len(errorStrings) == 0 {
		return nil
	}
	return errors.New(strings.Join(errorStrings, "\n"))
}

// collect collects sysstat data with the collector utility sadc. It runs the following command:
//     Sadc -S <Activity1> -S <Activity2> ... <collectInterval> 2 tmpFile
// The above command collects system metrics during <collectInterval> and saves it in binary form to tmpFile.
func (s *Sysstat) collect() error {
	options := []string{}
	for _, act := range s.Activities {
		options = append(options, "-S", act)
	}
	s.tmpFile = path.Join("/tmp", fmt.Sprintf("sysstat-%d", time.Now().Unix()))
	collectInterval := s.interval - parseInterval // collectInterval has to be smaller than the telegraf data collection interval

	if collectInterval < 0 { // If true, interval is not defined yet and Gather is run for the first time.
		collectInterval = 1 // In that case we only collect for 1 second.
	}

	options = append(options, strconv.Itoa(collectInterval), "2", s.tmpFile)
	cmd := execCommand(s.Sadc, options...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run command %s: %s", strings.Join(cmd.Args, " "), string(out))
	}
	return nil
}

// parse runs Sadf on the previously saved tmpFile:
//    Sadf -p -- -p <option> tmpFile
// and parses the output to add it to the telegraf.Accumulator acc.
func (s *Sysstat) parse(acc telegraf.Accumulator, option string, ts time.Time) error {
	cmd := execCommand(s.Sadf, s.sadfOptions(option)...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("running command '%s' failed: %s", strings.Join(cmd.Args, " "), err)
	}

	r := bufio.NewReader(stdout)
	csv := csv.NewReader(r)
	csv.Comma = '\t'
	csv.FieldsPerRecord = 6
	var measurement string
	// groupData to accumulate data when Group=true
	type groupData struct {
		tags   map[string]string
		fields map[string]interface{}
	}
	m := make(map[string]groupData)
	for {
		record, err := csv.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		device := record[3]
		value, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return err
		}

		tags := map[string]string{}
		if device != "-" {
			tags["device"] = device
			if addTags, ok := s.DeviceTags[device]; ok {
				for _, tag := range addTags {
					for k, v := range tag {
						tags[k] = v
					}
				}

			}
		}

		if s.Group {
			measurement = s.Options[option]
			if _, ok := m[device]; !ok {
				m[device] = groupData{
					fields: make(map[string]interface{}),
					tags:   make(map[string]string),
				}
			}
			g, _ := m[device]
			if len(g.tags) == 0 {
				for k, v := range tags {
					g.tags[k] = v
				}
			}
			g.fields[escape(record[4])] = value
		} else {
			measurement = s.Options[option] + "_" + escape(record[4])
			fields := map[string]interface{}{
				"value": value,
			}
			acc.AddFields(measurement, fields, tags, ts)
		}

	}
	if s.Group {
		for _, v := range m {
			acc.AddFields(measurement, v.fields, v.tags, ts)
		}
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command %s failed with %s", strings.Join(cmd.Args, " "), err)
	}
	return nil
}

// sadfOptions creates the correct options for the sadf utility.
func (s *Sysstat) sadfOptions(activityOption string) []string {
	options := []string{
		"-p",
		"--",
		"-p",
	}

	opts := strings.Split(activityOption, " ")
	options = append(options, opts...)
	options = append(options, s.tmpFile)

	return options
}

// escape removes % and / chars in field names
func escape(dirty string) string {
	var fieldEscaper = strings.NewReplacer(
		`%`, "pct_",
		`/`, "_per_",
	)
	return fieldEscaper.Replace(dirty)
}

func init() {
	s := Sysstat{
		Group:      true,
		Activities: dfltActivities,
	}
	sadf, _ := exec.LookPath("sadf")
	if len(sadf) > 0 {
		s.Sadf = sadf
	}
	inputs.Add("sysstat", func() telegraf.Input {
		return &s
	})
}
