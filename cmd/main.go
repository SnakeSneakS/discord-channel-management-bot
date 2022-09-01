package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/snakesneaks/discord-channel-management-bot/driver"
	"github.com/snakesneaks/discord-channel-management-bot/driver/db"
)

func Run() {
	log.Println("starting...")

	env, err := driver.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	db.Init(env)

	s := startDiscordSession(env)
	defer endDiscordSession(s)

	//gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")
}
