package openstack

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "io/ioutil"
    "net/http"
    )

func getToken(auth_url string, user string, pass string) (token, tenant_id string) {

    tr := &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
    }

    url := auth_url + ":5000/v2.0/tokens"

    var jsonStr = []byte(`{"auth":{"passwordCredentials":{"username": "` + user + `","password": "` + pass + `"},"tenantName": "admin"}}`)
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{Transport: tr}
    resp, _ := client.Do(req)
    defer resp.Body.Close()

    jsonDataFromHttp, _ := ioutil.ReadAll(resp.Body)

    var y map[string]interface{}
    json.Unmarshal([]byte(jsonDataFromHttp), &y)

    token = y["access"].(map[string]interface{})["token"].(map[string]interface{})["id"].(string)
    tenant_id = y["access"].(map[string]interface{})["token"].(map[string]interface{})["tenant"].(map[string]interface{})["id"].(string)

    return

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
