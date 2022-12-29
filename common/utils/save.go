package utils

import (
	"browser-reptile/common/config"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var lockMap = make(map[string]*sync.Mutex)
var lockMapLock = sync.Mutex{}

func SaveData(sockId int, buf []byte, isTo bool) {
	prefix := ""
	suffix := "\n"
	if isTo {
		prefix += ">>>"
		suffix += ">>>"
	} else {
		prefix += "<<<"
		suffix += "<<<"
	}

	prefix += time.Now().String()
	suffix += "over"

	prefix += "\n"
	suffix += "\n"

	fileName := config.LogDir + string(os.PathSeparator) + strconv.Itoa(sockId) + ".bin"
	lock := getLock(fileName)
	lock.Lock()
	defer lock.Unlock()

	if config.CommonConfig.LogTo == config.LogToConsole {
		fmt.Println(fileName)
		fmt.Println(prefix)
		fmt.Println(string(buf))
		fmt.Println(suffix)
	} else if config.CommonConfig.LogTo == config.LogToFile {
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			fmt.Println("open", fileName, "error:", err)
		}
		defer func() {
			_ = file.Close()
		}()

		_, _ = file.Write([]byte(prefix))
		_, _ = file.Write(buf)
		_, _ = file.Write([]byte(suffix))
	}
}

func getLock(name string) *sync.Mutex {
	lockMapLock.Lock()
	defer lockMapLock.Unlock()
	if lock, ok := lockMap[name]; ok {
		return lock
	}
	lockMap[name] = &sync.Mutex{}
	return lockMap[name]
}
