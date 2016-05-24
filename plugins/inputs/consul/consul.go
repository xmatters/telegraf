package consul

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/telegraf"
	"strings"
	"strconv"
	"log"
)

type Consul struct {
}

func flattenJson(js string) (map[string]interface{}, error) {
	flat := make(map[string]interface{})
	parsed := make(map[string]interface{})
	err := json.Unmarshal([]byte(js), &parsed)
	if err != nil {
		return flat, fmt.Errorf("failed to parse json: %s\n%v", js, err)
	}

	err = flatten(parsed, flat, make([]string, 0))
	if err != nil {
		return flat, fmt.Errorf("failed to flatten map: %v\n%v", parsed, err)
	}

	return flat, nil
}

func flatten(input interface{}, outMap map[string]interface{},
	prefix []string) error {
	fmt.Println("Called")
	switch value := input.(type) {
	case string, int64, float64, bool:
		fmt.Println("Handling primitive")
		outMap[makeKey(prefix...)] = value
	case map[string]interface{}:
		fmt.Println("Handling map")
		for k,v  := range value {
			prefix := append(prefix,k)
			flatten(v, outMap, prefix)
			prefix = prefix[:len(prefix)-1]
		}
	case []interface{}:
		fmt.Println("Handling array")
		for i, v := range value {
			prefix := append(prefix, strconv.Itoa(i))
			flatten(v, outMap, prefix)
			prefix = prefix[:len(prefix)-1]
		}
	case nil:
		fmt.Println("Handling nil")
		// silently ignore, since this is expected for some values and
		// could get noisy
	default:
		fmt.Println("Handling default")
		log.Printf("consul plugin is ignoring value for key: %s, could not " +
			"identify type of %v", makeKey(prefix...), value)
	}

	return nil
}

func makeKey(s ...string) string {
	key := strings.Join(s, ".")
	return strings.ToLower(key);
}


func (c *Consul) Gather(acc telegraf.Accumulator) error {

	return nil
}
