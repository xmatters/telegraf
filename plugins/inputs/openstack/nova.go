package openstack

import (
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"

    "github.com/influxdata/telegraf"
    "github.com/influxdata/telegraf/plugins/inputs"
)

type Nova struct {
    Authurl  string
    Username string
    Password string
}

type hypervisors struct {
    Status             string
    State              string
    Id                 int
    HypervisorHostname string `json:"hypervisor_hostname"`
}

type hypervisor_detail struct {
    Status               string
    Service              string
    VcpusUsed            int    `json:"vcpus_used"`
    HypervisorYype       string `json:"hypervisor_type"`
    LocalGbUsed          int    `json:"local_gb_used"`
    Vcpus                int
    HypervisorHostname   string `json:"hypervisor_hostname"`
    MemoryMbUsed         int    `json:"memory_mb_used"`
    MemoryMb             int    `json:"memory_mb"`
    CurrentWorkload      int    `json:"current_workload"`
    State                string
    HostIp               string `json:"host_ip"`
    CpuInfo              string `json:"cpu_info"`
    RunningVms           int    `json:"running_vms"`
    FreeDiskGb           int    `json:"free_disk_gb"`
    HypervisorVersion    int    `json:"hypervisor_version"`
    DiskAvailableLeast   int    `json:"disk_available_least"`
    LocalGb              int    `json:"disk_available_least"`
    FreeMamMb            int    `json:"free_ram_mb"`
    Id                   int
}

type hypervisor struct {
  Hypervisor hypervisor_detail
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

func getData(auth_url string, token string) (payload []byte) {

    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
    }

    req, _ := http.NewRequest("GET", auth_url, nil)
    req.Header.Set("X-Auth-Token", token)
    req.Header.Set("Accept", "application/json")

    client := &http.Client{Transport: tr}
    resp, _ := client.Do(req)
    defer resp.Body.Close()

    payload, _ = ioutil.ReadAll(resp.Body)
    return

}

func (n *Nova) Gather(acc telegraf.Accumulator) error {
    n.setDefaults()
    token, tenant_id := getToken(n.Authurl,n.Username,n.Password)
    url := fmt.Sprintf("%s:8774/v2/%s/os-hypervisors",  n.Authurl, tenant_id)

    hypervisor_list := getData(url,token)

    parsed := make(map[string][]hypervisors)
    json.Unmarshal([]byte(hypervisor_list), &parsed)
    for _, h := range parsed["hypervisors"] {
        url := fmt.Sprintf("%s:8774/v2/%s/os-hypervisors/%v",  n.Authurl, tenant_id, h.Id)
        hypervisor_detail_str := getData(url,token)
        parsed := hypervisor{}
        json.Unmarshal([]byte(hypervisor_detail_str), &parsed)
        fmt.Println(parsed.Hypervisor.HostIp)
    }

    return nil
}

func init() {
    inputs.Add("nova", func() telegraf.Input { return &Nova{} })
}
