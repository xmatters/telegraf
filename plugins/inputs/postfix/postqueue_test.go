package system

import (
  //"io/ioutil"
  //"os"
  "testing"

  "github.com/influxdata/telegraf/testutil"
  "github.com/stretchr/testify/assert"
)

func TestExtractingTotalCount(t *testing.T) {

  k := Postqueue {
    getPostQueueLog : func() ([]byte, error) {

      log := `6C2BF14C1893     3334 Wed Mar 30 01:51:53  MAILER-DAEMON
(delivery temporarily suspended: Host or domain name not found. Name service error for name=mail.abc.com type=AAAA: Host not found)
                                         jsmith@abc.com

631FB14C186A    19545 Mon Mar 28 16:51:53  MAILER-DAEMON
(delivery temporarily suspended: Host or domain name not found. Name service error for name=mail.abc.com type=AAAA: Host not found)
                                         jsmith@abc.com

600C614C180F    17108 Mon Mar 28 03:21:53  MAILER-DAEMON
(delivery temporarily suspended: Host or domain name not found. Name service error for name=mail.abc.com type=AAAA: Host not found)
                                         jsmith@abc.com

-- 2106 Kbytes in 205 Requests.`

      out := []byte(log) 
      return out, nil
    },
  }

  acc := testutil.Accumulator{}
  err := k.Gather(&acc)
  assert.NoError(t, err)

  fields := map[string]interface{}{
    "total_count": int32(205),
  }

  acc.AssertContainsFields(t, "postqueue", fields)
}

func TestTotalCountNotExist(t *testing.T) {

  k := Postqueue {
    getPostQueueLog : func() ([]byte, error) {

      log := `6C2BF14C1893     3334 Wed Mar 30 01:51:53  MAILER-DAEMON
(delivery temporarily suspended: Host or domain name not found. Name service error for name=mail.abc.com type=AAAA: Host not found)
                                         jsmith@abc.com

631FB14C186A    19545 Mon Mar 28 16:51:53  MAILER-DAEMON
(delivery temporarily suspended: Host or domain name not found. Name service error for name=mail.abc.com type=AAAA: Host not found)
                                         jsmith@abc.com

600C614C180F    17108 Mon Mar 28 03:21:53  MAILER-DAEMON
(delivery temporarily suspended: Host or domain name not found. Name service error for name=mail.abc.com type=AAAA: Host not found)
                                         jsmith@abc.com`

      out := []byte(log)
      return out, nil
    },
  }

  acc := testutil.Accumulator{}
  err := k.Gather(&acc)
  assert.Error(t, err)
  assert.Contains(t, err.Error(), "The last line does not contain the total count: []")
}

func TestEmptyLog(t *testing.T) {

  k := Postqueue {
    getPostQueueLog : func() ([]byte, error) {

      log := ``

      out := []byte(log)
      return out, nil
    },
  }

  acc := testutil.Accumulator{}
  err := k.Gather(&acc)
  assert.Error(t, err)
  assert.Contains(t, err.Error(), "The last line does not contain the total count: []")
}
