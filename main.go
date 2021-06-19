package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"time"
	"zgo.at/zcache"
)

// activityMetric is used to parse incoming traffic metrics info during activity capture.
type activityMetric struct {
	Value int `json:"value"`
}

// Use in-memory cache since it automatically provides TTL eviction and quick in-memory access.
var activityMap map[string]*zcache.Cache

// main serves as the application entry point.
func main() {

	// Setup activity map to store captured values
	activityMap = make(map[string]*zcache.Cache)

	// Use GoFiber to host REST API traffic
	app := fiber.New()

	// Activity metric routes
	routes := app.Group("/metric/:key")
	{
		routes.Post("", captureActivity)
		routes.Get("/sum", reportActivity)
	}

	// Start listening for traffic
	err := app.Listen(":3000")

	if err != nil {
		log.Panicln("Unable to start activity-tracker application on port 3000", err)
	}
}

// captureActivity records provided service traffic metrics.
func captureActivity(c *fiber.Ctx) error {

	// Use struct to hold parse JSON activity value
	activity := activityMetric{}

	// Parse JSON body into activity metric
	if err := c.BodyParser(&activity); err != nil {
		return err
	}

	// Obtain activity key from path parameter
	activityKey := c.Params("key")

	if len(activityKey) == 0 {
		return fiber.NewError(fiber.StatusPreconditionFailed, "Invalid activity used as path parameter")
	}

	activityCache, found := activityMap[activityKey]

	if !found {
		activityCache = zcache.New(time.Hour*1, time.Minute*1)
		activityMap[activityKey] = activityCache
	}

	activityCache.Set(uuid.NewString(), activity.Value, zcache.DefaultExpiration)

	c.Set("Content-Type", "application/json")
	return c.SendString("{}")
}

// reportActivity returns service traffic metrics for the last hour.
func reportActivity(c *fiber.Ctx) error {

	// Obtain activity key from path parameter
	activityKey := c.Params("key")

	if len(activityKey) == 0 {
		return fiber.NewError(fiber.StatusPreconditionFailed, "Invalid activity used as path parameter")
	}

	activity := activityMetric{}

	activityCache, found := activityMap[activityKey]

	if found {

		for _, cacheKey := range activityCache.Keys() {

			cacheEntry, isPresent := activityCache.Get(cacheKey)

			if isPresent {
				activity.Value += cacheEntry.(int)
			}
		}
	}

	return c.JSON(activity)
}
