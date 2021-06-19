package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
	"zgo.at/zcache"
)

// activityMap uses an in-memory cache since it automatically provides TTL eviction and quick in-memory access
var activityMap map[string]*zcache.Cache

// appServer holds a reference to the application server so it can be shutdown if needed
var appServer *fiber.App

// appServerInitWaitDuration is the time increment used while waiting for the app server to complete startup
const appServerInitWaitDuration = time.Millisecond * 250

// appServerListenAddress is the host and port for the application server to use for binding
const appServerListenAddress = "localhost:3000"

// isInitializationComplete will set to true after API routing initialization is finished
var isInitializationComplete bool

// main serves as the application entry point.
func main() {

	// Reset init flag to false to be safe
	isInitializationComplete = false

	// Setup activity map to store captured values
	activityMap = make(map[string]*zcache.Cache)

	// Use GoFiber to host REST API traffic
	appServer = fiber.New(fiber.Config{
		// Make sure open client connections do not hinder app server from shutting down gracefully
		ReadTimeout: time.Second * 1,
	})

	// Activity metric routes
	routes := appServer.Group("/metric/:key")
	{
		routes.Post("", captureActivity)
		routes.Get("/sum", reportActivity)
	}

	// Start listening for traffic
	isInitializationComplete = true
	log.Println("Starting application server at:", appServerListenAddress)
	err := appServer.Listen(appServerListenAddress)

	if err != nil {
		log.Panicln("Unable to start activity-tracker application at: "+appServerListenAddress, err)
	}
}

// shutdown stops the application server if it is running.
func shutdown() {

	log.Println("Entering shutdown")

	if appServer != nil {

		log.Println("Shutting down application server")
		err := appServer.Shutdown()

		if err != nil {
			log.Println("Error received when shutting down application server: " + err.Error())
		}
	}

	log.Println("Exiting shutdown")
}

// waitForInitialization ensures app server is fully started before returning.
func waitForInitialization() {

	log.Println("Entering waitForInitialization")

	for !isInitializationComplete {
		time.Sleep(appServerInitWaitDuration)
	}

	// Wait an additional period to be certain app server has started listening
	time.Sleep(appServerInitWaitDuration)
	log.Println("Exiting waitForInitialization")
}
