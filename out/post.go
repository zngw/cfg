// @Title
// @Description $
// @Author  55
// @Date  2021/12/3
package out

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/zngw/cfg/conf"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type updata struct {
	Name string
	Time int64
	Data string
	Sign string
}

type OutPost struct {
	Type string
	Url  string
	Key  string
}

func (o *OutPost) Init(path string) (err error) {
	o.Type = conf.BuildToPost
	o.Url = conf.Cfg.PostUrl
	o.Key = conf.Cfg.PostKey
	return
}

func (o *OutPost) GetType() (t string) {
	return o.Type
}

// 转JS
func (o *OutPost) OutTo(file string, attr bool, keys *[]string, s *[]map[string]interface{}) (err error) {
	data, err := json.Marshal(*s)
	if err != nil {
		return
	}

	now := time.Now().UnixNano() / 1e6
	up := updata{
		Name: file,
		Time: now,
		Data: string(data),
	}

	if len(o.Key) > 0 {
		h := md5.New()
		h.Write([]byte(file + strconv.FormatInt(now, 10) + o.Key))
		up.Sign = hex.EncodeToString(h.Sum(nil))
	}

	jsonStr, err := json.Marshal(up)

	req, err := http.NewRequest("POST", o.Url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	body, _ := ioutil.ReadAll(resp.Body)
	if code != 200 {
		err = fmt.Errorf("Error: %s, %d", file, code)
		return
	}

	if string(body) != "OK" {
		err = fmt.Errorf("Error: %s, %s", file, string(body))
		return
	}

	return
}
