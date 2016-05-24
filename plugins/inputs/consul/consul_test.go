package consul
import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// TODO:
// - populate advertise addrs section of sampleSelf

func TestFlattenJson(t *testing.T) {
	flat, err := flattenJson(sampleSelf)
	assert.NoError(t, err)
	//fmt.Printf("%v", flat)
	for k,v := range map[string]interface{} {
		"config.acldatacenter": "DC1",
		"config.bootstrapexpect": float64(0),
		"config.dnsconfig.allowstale": true,
		"config.dnsconfig.servicettl.*": float64(60000000000),
		"config.dnsconfig.servicettl.service-one": float64(86400000000000),
		"coord.vec.0": float64(-0.028174781623480347),
		"coord.vec.4": float64(-0.01798695632072015),
		"coord.vec.7": float64(0.02062575913622964),
		"member.tags.vsn_min": "1",
	} {
		assert.Equal(t, v, flat[k])
	}

	assert.NotContains(t, flat, "config.advertiseaddrs.rpc")
}

var sampleSelf = `
{
    "Config": {
        "ACLDatacenter": "DC1",
        "ACLDefaultPolicy": "allow",
        "ACLDownPolicy": "extend-cache",
        "ACLTTL": 30000000000,
        "ACLTTLRaw": "",
        "Addresses": {
            "DNS": "",
            "HTTP": "127.0.0.1",
            "HTTPS": "",
            "RPC": ""
        },
        "AdvertiseAddr": "192.168.0.10",
        "AdvertiseAddrWan": "192.168.0.10",
        "AdvertiseAddrs": {
            "RPC": null,
            "RPCRaw": "",
            "SerfLan": null,
            "SerfLanRaw": "",
            "SerfWan": null,
            "SerfWanRaw": ""
        },
        "AtlasEndpoint": "",
        "AtlasInfrastructure": "",
        "AtlasJoin": false,
        "BindAddr": "0.0.0.0",
        "Bootstrap": false,
        "BootstrapExpect": 0,
        "CAFile": "",
        "CertFile": "",
        "CheckUpdateInterval": 300000000000,
        "ClientAddr": "0.0.0.0",
        "DNSConfig": {
            "AllowStale": true,
            "EnableTruncate": false,
            "MaxStale": 30000000000,
            "NodeTTL": 30000000000,
            "OnlyPassing": false,
            "ServiceTTL": {
                "*": 60000000000,
                "service-one": 86400000000000,
                "service-two": 300000000000,
                "service-three": 300000000000,
                "service-four": 604800000000000
            }
        },
        "DNSRecursor": "",
        "DNSRecursors": [],
        "DataDir": "/opt/consul",
        "Datacenter": "dc1",
        "DevMode": false,
        "DisableAnonymousSignature": false,
        "DisableCoordinates": false,
        "DisableRemoteExec": false,
        "DisableUpdateCheck": true,
        "DogStatsdAddr": "",
        "DogStatsdTags": null,
        "Domain": "consul.",
        "EnableDebug": false,
        "EnableSyslog": false,
        "EnableUi": false,
        "HTTPAPIResponseHeaders": null,
        "KeyFile": "",
        "LeaveOnTerm": false,
        "LogLevel": "INFO",
        "NodeName": "service-one-instance1",
        "PidFile": "",
        "Ports": {
            "DNS": 8600,
            "HTTP": 8500,
            "HTTPS": 8443,
            "RPC": 8400,
            "SerfLan": 8301,
            "SerfWan": 8302,
            "Server": 8300
        },
        "Protocol": 2,
        "Reap": null,
        "RejoinAfterLeave": false,
        "RetryIntervalRaw": "",
        "RetryIntervalWanRaw": "",
        "RetryJoin": [
            "consul.service.dc1.consul"
        ],
        "RetryJoinWan": [],
        "RetryMaxAttempts": 0,
        "RetryMaxAttemptsWan": 0,
        "Revision": "c933efde50d25395c7b5a42167578fda603d43d8+CHANGES",
        "Server": false,
        "ServerName": "",
        "SessionTTLMin": 0,
        "SessionTTLMinRaw": "",
        "SkipLeaveOnInt": false,
        "StartJoin": [],
        "StartJoinWan": [],
        "StatsdAddr": "127.0.0.1:8125",
        "StatsiteAddr": "",
        "StatsitePrefix": "consul",
        "SyslogFacility": "LOCAL0",
        "UiDir": "",
        "UnixSockets": {
            "Grp": "",
            "Perms": "",
            "Usr": ""
        },
        "VerifyIncoming": false,
        "VerifyOutgoing": false,
        "VerifyServerHostname": false,
        "Version": "0.6.3",
        "VersionPrerelease": "",
        "Watches": null
    },
    "Coord": {
        "Adjustment": -0.0021454290888388906,
        "Error": 1.1263614142494565,
        "Height": 5.4525869491170155e-05,
        "Vec": [
            -0.028174781623480347,
            -0.019054033364910915,
            0.029057685564317294,
            -0.0013155463084249335,
            -0.01798695632072015,
            0.01732762316164008,
            -0.013734912098256707,
            0.02062575913622964
        ]
    },
    "Member": {
        "Addr": "192.168.0.10",
        "DelegateCur": 4,
        "DelegateMax": 4,
        "DelegateMin": 2,
        "Name": "service-one-instanec1",
        "Port": 8301,
        "ProtocolCur": 2,
        "ProtocolMax": 3,
        "ProtocolMin": 1,
        "Status": 1,
        "Tags": {
            "build": "0.6.3:c933efde",
            "dc": "dc1",
            "role": "node",
            "vsn": "2",
            "vsn_max": "3",
            "vsn_min": "1"
        }
    }
}
`

var sampleChecks = `
{
    "service:service-one": {
        "CheckID": "service:service-one",
        "CreateIndex": 0,
        "ModifyIndex": 0,
        "Name": "Service 'service-one' check",
        "Node": "service-one-instance1",
        "Notes": "",
        "Output": "HTTP GET http://192.168.0.10:80/service-one-status: 200 OK Output: <!DOCTYPE html>\n<html>\n  <head>\n    <title>200 OK</title>\n  </head>\n  <body>\n    <h1>Error 200 OK</h1>\n    <p>OK</p>\n   </body>\n</html>\n",
        "ServiceID": "service-one",
        "ServiceName": "service-one",
        "Status": "passing"
    },
    "service:service-two": {
        "CheckID": "service:service-two",
        "CreateIndex": 0,
        "ModifyIndex": 0,
        "Name": "Service 'service-two' check",
        "Node": "prd-service-one-instance1-ooqmc4dksfb3",
        "Notes": "",
        "Output": "HTTP GET http://192.168.0.20:80/service-two-status: 200 OK Output: <!DOCTYPE html>\n<html>\n  <head>\n    <title>200 OK</title>\n  </head>\n  <body>\n    <h1>Error 200 OK</h1>\n    <p>OK</p>\n   </body>\n</html>\n",
        "ServiceID": "service-two",
        "ServiceName": "service-two",
        "Status": "passing"
    }
}`

var sampleServices = `
{
    "service-one": {
        "Address": "192.168.0.10",
        "CreateIndex": 0,
        "EnableTagOverride": false,
        "ID": "service-one",
        "ModifyIndex": 0,
        "Port": 80,
        "Service": "service-one",
        "Tags": [
            "prd"
        ]
    },
    "service-two": {
        "Address": "192.168.0.20",
        "CreateIndex": 0,
        "EnableTagOverride": false,
        "ID": "repo",
        "ModifyIndex": 0,
        "Port": 80,
        "Service": "repo",
        "Tags": [
            "prd"
        ]
    }
}`