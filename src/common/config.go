package common

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"regexp"
)

// 配置
type Config struct {
	Redis Redis
	Rule  RuleConf
}

//REDIS
type Redis struct {
	Address  string
	Port     string
	Password string
	DBIndex  int
}

type RuleConf struct {
	Pattern    []string
	ScanCount  int
	Method     string
	ExpireTime int64
}

/**
 *  新建config
 *  @param path 配置文件的路劲
 *  @return 返回Config的指针
 */
func NewConfig(path string) *Config {
	config := new(Config)
	config.LoadConfig(path)
	return config
}

/**
 *  导入config
 *  @param path 配置文件的路劲
 */
func (self *Config) LoadConfig(path string) {
	configFile, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to open config file '%s': %s\n", path, err)
		return
	}

	fi, _ := configFile.Stat()
	if fi.Size() == 0 {
		log.Printf("config file (%q) is empty, skipping", path)
		return
	}

	buffer := make([]byte, fi.Size())
	_, err = configFile.Read(buffer)

	// 去掉换行
	buffer, err = StripComments(buffer)

	// 替换环境变量
	buffer = []byte(os.ExpandEnv(string(buffer)))

	err = json.Unmarshal(buffer, self)
	if err != nil {
		log.Fatalf("Failed unmarshalling json: %s\n", err)
		return
	}
}

func StripComments(data []byte) ([]byte, error) {
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0) // Windows
	lines := bytes.Split(data, []byte("\n"))                //split to muli lines
	filtered := make([][]byte, 0)

	for _, line := range lines {
		match, err := regexp.Match(`^\s*#`, line)
		if err != nil {
			return nil, err
		}
		if !match {
			filtered = append(filtered, line)
		}
	}

	return bytes.Join(filtered, []byte("\n")), nil
}
