package main

import (
	"fmt"
	"go.code.as/writeas.v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	checkInterval = 60 * time.Second
)

func main() {

	// BEGIN CONFIGURATION //////////////////////////////////////////
	// TODO: replace this with the post ID you want to have self-destruct
	id := "8wyuivd67md9mmjf"
	// TODO: replace these with your Write.as login credentials
	username := "USERNAME"
	password := "PASSWORD"
	// END CONFIGURATION ////////////////////////////////////////////

	// Begin the work
	log.Println("Starting up...")

	// Log in to Write.as
	c := writeas.NewClient()
	_, err := c.LogIn(username, password)
	if err != nil {
		log.Printf("Unable to log in: %s", err)
		os.Exit(1)
	} else {
		log.Println("Authenticated.")
	}
	// Automatically invalidate auth token when finished
	defer logout(c)

	// Clean up when we receive an interrupt
	qc := make(chan os.Signal, 2)
	signal.Notify(qc, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-qc
		log.Printf("Packing up without any destruction.")
		logout(c)
		os.Exit(0)
	}()

	// Check for the post on a regular interval, starting now
	ticker := time.NewTicker(checkInterval)
	quit := make(chan struct{})
	err = checkPost(c, id, quit)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	for {
		select {
		case <-ticker.C:
			err = checkPost(c, id, quit)
			if err != nil {
				log.Printf("%s", err)
				return
			}
		case <-quit:
			log.Printf("Packing up.")
			ticker.Stop()
			return
		}
	}
}

func checkPost(c *writeas.Client, id string, quit chan struct{}) error {
	log.Println("Inspecting post...")
	p, err := c.GetPost(id)
	if err != nil {
		return fmt.Errorf("Unable to fetch post: %s", err)
	}
	log.Printf("Post %s has %d view(s)", id, p.Views)

	if p.Views > 0 {
		log.Println("Preparing to self-destruct...")
		err = c.DeletePost(id, "")
		if err != nil {
			return fmt.Errorf("Unable to delete post: %s", err)
		}
		log.Println("BOOM!")
		close(quit)
	}

	return nil
}

func logout(c *writeas.Client) {
	c.LogOut()
	log.Println("Logged out.")
}
