package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

type Layer string

const (
	Default  Layer = "Default"
	Homerow  Layer = "Homerow"
	NavLeft  Layer = "NavLeft"
	NavRight Layer = "NavRight"
	Qwerty   Layer = "Qwerty"
	Mouse    Layer = "Mouse"
	System   Layer = "System"
	Numpad   Layer = "Numpad"
)

var stringToLayer = map[string]Layer{
	"default":   Default,
	"homerow":   Homerow,
	"nav_left":  NavLeft,
	"nav_right": NavRight,
	"qwerty":    Qwerty,
	"mouse":     Mouse,
	"system":    System,
	"numpad":    Numpad,
}

type State string

const (
	Idle     State = "Idle"
	Info     State = "Info"
	Good     State = "Good"
	Warning  State = "Warning"
	Critical State = "Critical"
)

type Block struct {
	Icon  string `json:"icon"`
	State State  `json:"state"`
	Text  Layer  `json:"text"`
}

func defaultBlock() ([]byte, error) {
	block := Block{
		Icon:  "weather_thunder",
		State: Good,
		Text:  Default,
	}
	bytes, err := json.Marshal(&block)
	return bytes, err
}

func isNotExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return true
	}
	return false
}

func read(path string) error {
	if isNotExist(path) {
		bytes, err := defaultBlock()
		if err != nil {
			return err
		}
		fmt.Println(string(bytes))
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))
	return nil
}

func write(path string, layer Layer) error {
	block := Block{
		Icon:  "weather_thunder",
		State: Good,
		Text:  layer,
	}

	bytes, err := json.Marshal(&block)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, bytes, 0666); err != nil {
		return err
	}

	return nil
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	path := path.Join(home, ".layer.json")

	args := os.Args
	if len(args) > 1 {
		layer := stringToLayer[args[1]]
		if err := write(path, layer); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := read(path); err != nil {
			log.Fatal(err)
		}
	}
}
