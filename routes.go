package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"sync"
	"zgo.at/zcache"
)

// mutex is used to ensure only one request accesses the activityMap at one time. Alternative solutions are possible
var mutex = &sync.Mutex{}

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

	// Obtain lock and use activity look-up key
	mutex.Lock()
	activityCache, found := activityMap[activityKey]

	if !found {

		// Since activity has not been previously used, create new cache instance
		activityCache = zcache.New(activityTTLDuration, activityTTLCheckInterval)
		activityMap[activityKey] = activityCache
	}

	// Create new cache entry using provided value. Utilizing UUID to ensure each cache entry is unique
	mutex.Unlock()
	activityCache.Set(uuid.NewString(), activity.Value, zcache.DefaultExpiration)

	// Using SendString operation rather than normal c.JSON to avoid placing word "OK" in the response body
	c.Set("Content-Type", "application/json")
	return c.SendString("{}")
}

// reportActivity returns recent service traffic metrics.
func reportActivity(c *fiber.Ctx) error {

	// Obtain activity key from path parameter. Due to routing matching, it is impossible for the activityKey to be empty
	activityKey := c.Params("key")

	// Use struct to hold activity count since it can be returned as JSON by Go-Fiber
	activity := activityMetric{}

	// Attempt to find matching cache using provided activity key
	activityCache, found := activityMap[activityKey]

	if found {

		log.Println("Found matching key entry for activityKey:", activityKey)

		// Sum all activity counts present in the cache
		for _, cacheKey := range activityCache.Keys() {

			cacheEntry, isPresent := activityCache.Get(cacheKey)

			if isPresent {
				activity.Value += cacheEntry.(int)
			}
		}

		log.Println("Completed summing activity counts for activityKey:", activityKey)

	} else {
		log.Println("Unable to find matching key entry for activityKey:", activityKey)
	}

	return c.JSON(activity)
}
