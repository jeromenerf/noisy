package main

import (
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/dddpaul/golang-evdev/evdev"
	"github.com/spf13/viper"
)

func main() {

	var (
		dev    *evdev.InputDevice
		events []evdev.InputEvent
		err    error
	)

	pkgpath := path.Join(os.Getenv("GOPATH"), "src", "github.com", "jeromenerf", "noisy")

	viper.SetDefault("device", "/dev/input/event0")
	viper.SetDefault("location", path.Join(pkgpath, "samples"))
	viper.SetDefault("sounds", map[string]string{"hi": "hi.mp3", "lo": "lo.mp3", "default": "default.mp3"})

	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/noisy")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	device := viper.GetString("device")
	location := viper.GetString("location")
	sounds := viper.GetStringMapString("sounds")

	dev, err = evdev.Open(device)
	if err != nil {
		log.Fatalf("unable to open input device")
	}

	log.Printf("Listening for events ...\n")

	ch := make(chan string, 12)
	go worker(ch, location, sounds)

	for {
		events, err = dev.Read()
		for _, ev := range events {
			var codeName string
			code := int(ev.Code)
			evType := int(ev.Type)
			if m, ok := evdev.ByEventType[evType]; ok {
				codeName = m[code]
			}
			if evType == evdev.EV_KEY && ev.Value == 1 {
				ch <- codeName
			}
		}
	}
}

func worker(ch <-chan string, location string, sounds map[string]string) {
	key := "default"

	for codename := range ch {
		switch codename {
		case "KEY_ENTER",
			"KEY_SPACE",
			"KEY_BACKSPACE",
			"KEY_DELETE",
			"KEY_ESC":
			key = "lo"

		case "KEY_TAB",
			"KEY_RIGHT",
			"KEY_LEFT",
			"KEY_UP",
			"KEY_DOWN",
			"KEY_RIGHTSHIFT",
			"KEY_LEFTSHIFT",
			"KEY_CAPSLOCK",
			"KEY_WAKEUP",
			"KEY_LEFTALT",
			"KEY_RIGHTALT",
			"KEY_RIGHTCTRL",
			"KEY_LEFTCTRL",
			"KEY_LEFTMETA",
			"KEY_COMPOSE":
			key = "hi"

		default:
			key = "default"

		}
		sound := path.Join(location, sounds[key])
		go func(sound string) {
			_, err := exec.Command("mpg123", "-q", sound).CombinedOutput()
			if err != nil {
				log.Println(err)
			}
		}(sound)
	}
}
