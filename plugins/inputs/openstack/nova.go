package openstack

import (
    "fmt"

    "github.com/influxdata/telegraf"
    "github.com/influxdata/telegraf/plugins/inputs"
)

type Nova struct {
    Authurl  string
    Username string
    Password string
}

func (n *Nova) Description() string {
    return "A plugin to gather Openstack Nova metrics"
}

func (n *Nova) setDefaults() {
        if n.Authurl == "" {
                fmt.Println("No authurl given, exiting")
        }
        if n.Username == "" {
                fmt.Println("No username given, exiting")
        }
        if n.Password == "" {
                fmt.Println("No password given, exiting")
        }
}

func (n *Nova) SampleConfig() string {
    return `
  # The keystone endpoint to authenticate against, can be http or https
  authurl = "http://keystone.url"
  ## Username to authenticate with
  username = "admin"
  ## Password to authenticate with
  password = "admin"
`
}

func (n *Nova) Gather(acc telegraf.Accumulator) error {
    n.setDefaults()
    token, tenant_id := getToken(n.Authurl,n.Username,n.Password)

    fmt.Println(token, tenant_id)

    return nil
}

func init() {
    inputs.Add("nova", func() telegraf.Input { return &Nova{} })
}
