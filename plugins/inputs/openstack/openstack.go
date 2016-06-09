package openstack

import (
    "fmt"

    "github.com/influxdata/telegraf"
    "github.com/influxdata/telegraf/plugins/inputs"
)

type Openstack struct {
    Authurl  string
    Username string
    Password string
    Cluster string
    Nova bool
}

func (o *Openstack) Description() string {
    return "A wrapper to gather Openstack metrics from requested api endpoint"
}

func (o *Openstack) setDefaults() {
        if o.Authurl == "" {
                fmt.Println("No authurl given, exiting")
        }
        if o.Username == "" {
                fmt.Println("No username given, exiting")
        }
        if o.Password == "" {
                fmt.Println("No password given, exiting")
        }
        if o.Cluster == "" {
                fmt.Println("No cluster given, exiting")
        }
}

func (o *Openstack) SampleConfig() string {
    return `
  # The keystone endpoint to authenticate against, can be http or https
  authurl = "http://keystone.url"
  ## Username to authenticate with
  username = "admin"
  ## Password to authenticate with
  password = "admin"
  ## Enable Nova stats collection
  nova = true
`
}

func (o *Openstack) Gather(acc telegraf.Accumulator) error {
    if o.Nova == true {
        nova := gatherNova()
    }

    return nil
}

func init() {
    inputs.Add("openstack", func() telegraf.Input { return &Openstack{} })
}
