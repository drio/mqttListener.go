package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const version = "0.0.1"
const mqttTopicDefault = "zigbee2mqtt/aquara-door-01"

type jsonPayload struct {
	Contact bool `json:"contact"`
}

func detectPlayer() string {
	player := "/usr/bin/aplay"
	os_type := exec.Command("uname", "-s")
	os_type_output, err := os_type.Output()
	os_type_trimmed := strings.TrimSuffix(string(os_type_output), "\n")
	log.Println(fmt.Sprintf("Os_type: %s", os_type_trimmed))
	if err != nil {
		panic(err)
	}
	if string(os_type_trimmed) == "Darwin" {
		player = "mpv"
	}

	log.Println(fmt.Sprintf("Player: %s", player))
	return player
}

func genPlaySong(player, songPath string) func() {
	return func() {
		log.Println(fmt.Sprintf("Running cmd: %s %s", player, songPath))
		play_cmd := exec.Command(player, songPath)
		play_cmd.Start()
		play_cmd.Wait()
		log.Println(fmt.Sprintf("cmd status: %+v", *play_cmd.ProcessState))
	}
}

func startMQTT(url string, user string, pass string, topic string, playOpen, playClose func()) mqtt.Client {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	hostname, _ := os.Hostname()

	id := hostname + strconv.Itoa(time.Now().Second())
	opts := mqtt.NewClientOptions().AddBroker(url).SetClientID(id).SetCleanSession(true)
	opts.SetUsername(user)
	opts.SetPassword(pass)

	onMessageReceived := (func(client mqtt.Client, msg mqtt.Message) {
		var jp jsonPayload
		err := json.Unmarshal([]byte(msg.Payload()), &jp)
		if err != nil {
			log.Printf("Error processing payload: %s", err)
			return
		}

		if jp.Contact {
			log.Println("Door CLOSED. Playing sound.")
			playClose()
		}

		if !jp.Contact {
			log.Println("Door OPEN. Playing sound.")
			playOpen()
		}
	})

	opts.OnConnect = func(c mqtt.Client) {
		if token := c.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", url)
	}

	return c
}

func main() {
	log.Println(fmt.Sprintf("Starting mqttSoundPlayer version: [%s] ...", version))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	url := flag.String("url", "tcp://192.168.8.180:1883", "MQTT server url")
	user := flag.String("user", "shelly", "MQTT username")
	pass := flag.String("password", "", "MQTT password")
	topic := flag.String("topic", mqttTopicDefault, "MQTT topic to subscribe to")
	soundsDir := flag.String("sounds", "/opt/sounds", "Sounds directory")
	flag.Parse()

	if *pass == "" {
		log.Println("MQTT password not provided. Bailing out.")
		os.Exit(0)
	}

	_, err := os.Stat(*soundsDir)
	if os.IsNotExist(err) {
		log.Fatal("Sound dir does not exist: ", *soundsDir)
	}
	log.Printf("sound dir is: %s", *soundsDir)

	player := detectPlayer()
	playOpen := genPlaySong(player, fmt.Sprintf("%s/coin.wav", *soundsDir))
	playDream := genPlaySong(player, fmt.Sprintf("%s/dream.wav", *soundsDir))
	playClose := genPlaySong(player, fmt.Sprintf("%s/newstart.wav", *soundsDir))
	log.Println("Starting ...")
	playDream()

	startMQTT(*url, *user, *pass, *topic, playOpen, playClose)

	<-c
	log.Println("Bye")
	os.Exit(0)
}
