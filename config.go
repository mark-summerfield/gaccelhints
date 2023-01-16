// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/mark-summerfield/accelhint"
	"github.com/mark-summerfield/gong"
)

type Config struct {
	dirty    bool
	filename string
	x        int // use getPosition()
	y        int // use getPosition()
	width    int
	height   int
	xdec     int
	ydec     int
	alphabet string
	marker   byte
}

func NewConfig() *Config {
	config := Config{x: 0, y: 0, width: defaultWidth, height: defaultHeight,
		xdec: 0, ydec: 0, alphabet: accelhint.Alphabet,
		marker: accelhint.Marker}
	found, filename := getConfigFilename()
	if found {
		config.filename = filename
		config.load()
	}
	return &config
}

func (me *Config) getPosition() (int, int) {
	return maxInt(0, me.x-me.xdec), maxInt(0, me.y-me.ydec)
}

func (me *Config) load() {
	file, err := os.Open(me.filename)
	if err != nil {
		log.Println("Failed to open configuration file:", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		i := strings.IndexByte(line, ':')
		if i == -1 {
			i = strings.IndexByte(line, '=')
		}
		if i == -1 {
			continue
		}
		key := strings.ToUpper(strings.TrimSpace(line[:i]))
		value := strings.TrimSpace(line[i+1:])
		switch key {
		case "X":
			me.x = gong.StrToInt(value, 0)
		case "Y":
			me.y = gong.StrToInt(value, 0)
		case "WIDTH":
			me.width = gong.StrToInt(value, 640)
		case "HEIGHT":
			me.height = gong.StrToInt(value, 480)
		case "XDEC":
			me.xdec = gong.StrToInt(value, 0)
		case "YDEC":
			me.ydec = gong.StrToInt(value, 0)
		case "ALPHABET":
			me.alphabet = value
		case "MARKER":
			me.marker = value[0]
		default:
			log.Println("unrecognized configuration key: %q", key)
		}
	}
	me.dirty = false
}

func (me *Config) updateGeometry(x, y, width, height int) {
	if x != me.x || y != me.y || width != me.width || height != me.height {
		me.dirty = true
		me.x = x
		me.y = y
		me.width = width
		me.height = height
	}
}

func (me *Config) save() bool {
	file, err := os.OpenFile(me.filename, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Println("Failed to save configuration file:", err)
		return false
	}
	defer file.Close()
	_, err = file.WriteString(me.String())
	return err != nil
}

func (me *Config) String() string {
	return fmt.Sprintf("x: %d\ny: %d\nwidth: %d\nheight: %d\nxdec: %d\n"+
		"ydec: %d\nalphabet: %s\nmarker: %c\n", me.x, me.y, me.width,
		me.height, me.xdec, me.ydec, me.alphabet, me.marker)
}

func getConfigFilename() (bool, string) {
	var filenames []string
	var fallback1 string
	var fallback2 string
	configDir, err := os.UserConfigDir()
	if err == nil {
		fallback1 = path.Join(configDir, configFilename)
		filenames = append(filenames, fallback1)
	}
	homeDir, err := os.UserHomeDir()
	if err == nil {
		fallback2 = path.Join(homeDir, "."+configFilename)
		filenames = append(filenames, fallback2)
		filenames = append(filenames, path.Join(homeDir, "data",
			configFilename))
	}
	if len(filenames) == 0 {
		filenames = append(filenames, configFilename)
	}
	for _, filename := range filenames {
		stat, err := os.Stat(filename)
		if err == nil && !stat.IsDir() {
			return true, filename // found
		}
	}
	if fallback1 != "" {
		return false, fallback1
	}
	if fallback2 != "" {
		return false, fallback2
	}
	return false, configFilename
}
