package openstack

import (
    "encoding/json"
    "fmt"

    "github.com/influxdata/telegraf"
)

type Nova struct {
    Authurl   string
    Tenant_id string
    Token     string
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

func (n *Nova) Gather(acc telegraf.Accumulator) error {
    url := fmt.Sprintf("%s:8774/v2/%s/os-hypervisors",  n.Authurl, n.Tenant_id)

    hypervisor_list := getData(url,n.Token)

    parsed := make(map[string][]hypervisors)
    json.Unmarshal([]byte(hypervisor_list), &parsed)
    for _, h := range parsed["hypervisors"] {
        url := fmt.Sprintf("%s:8774/v2/%s/os-hypervisors/%v",  n.Authurl, n.Tenant_id, h.Id)
        hypervisor_detail_str := getData(url,n.Token)
        parsed := hypervisor{}
        json.Unmarshal([]byte(hypervisor_detail_str), &parsed)
        fmt.Println(parsed.Hypervisor.HostIp)
    }

    return nil
}

func init() {
}
