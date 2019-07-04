package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"vgontakte/vgontakte"
)

var (
	parameterNotFound = func(param string) error {
		return fmt.Errorf("parameter \"%v\" is not found in config", param)
	}
	couldNotConvertTo = func(param string, value string, t string) error {
		return fmt.Errorf("parameter \"%v\" (value is \"%v\") could not be converted to %v", param, value, t)
	}
	claConfig *claconfig
)

func GetCommandLineArgsConfig() vgontakte.Config {
	if claConfig == nil {
		claConfig = &claconfig{}
	}
	return claConfig
}

type claconfig struct {
	args map[string]string
}

func (c *claconfig) GetInt(name string) (int, error) {
	c.ensureArgs()
	if val, ok := c.args[name]; ok {
		res, err := strconv.Atoi(val)
		if err != nil {
			return 0, couldNotConvertTo(name, val, "int")
		}
		return res, nil
	}
	return 0, parameterNotFound(name)
}

func (c *claconfig) GetString(name string) (string, error) {
	if val, ok := c.args[name]; ok {
		return val, nil
	}
	return "", parameterNotFound(name)
}

func (c *claconfig) ensureArgs() {
	if c.args == nil {
		c.args = make(map[string]string, 0)
		for _, v := range os.Args[1:] {
			if strings.Contains(v, "=") {
				val := strings.TrimLeft(v, "-")
				name := string([]rune(val)[:strings.Index(val, "=")])
				value := string([]rune(val)[strings.Index(val, "=")+1:])
				c.args[name] = value
			}
		}
	}

}
