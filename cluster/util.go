package cluster

import (
	log "github.com/chuccp/coke-log"
	"github.com/chuccp/cokePush/message"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func MachineId() string {
	f, err := os.OpenFile(".machineId", os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()
	if err == nil {
		data, err := ioutil.ReadAll(f)
		if err == nil {
			if len(data) == 0 {
				uid := strconv.FormatUint(uint64(message.MsgId()), 36)
				f.Write([]byte(uid))
				return uid
			}
			return strings.TrimSpace(string(data))
		}
	}
	log.Panic("生成机器码错误,请检查程序的读写权限")
	return ""
}
