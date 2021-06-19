package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"time"
	"zgo.at/zcache"
)

// captureActivity records provided service traffic metrics.
func captureActivity(c *fiber.Ctx) error {

	// Use struct to hold parse JSON activity value
	activity := activityMetric{}

	// Parse JSON body into activity metric
	if err := c.BodyParser(&activity); err != nil {
		return err
	}

	// Obtain activity key from path parameter. Due to routing matching, it is impossible for the activityKey to be empty
	activityKey := c.Params("key")

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

	// Obtain activity key from path parameter. Due to routing matching, it is impossible for the activityKey to be empty
	activityKey := c.Params("key")

	activity := activityMetric{}

	activityCache, found := activityMap[activityKey]

	if found {

		for _, cacheKey := range activityCache.Keys() {

			cacheEntry, isPresent := activityCache.Get(cacheKey)

			if isPresent {
				activity.Value += cacheEntry.(int)
			}
		}
	} else {
		log.Println("Unable to find matching key entry for activityKey:", activityKey)
	}

	return c.JSON(activity)
}
