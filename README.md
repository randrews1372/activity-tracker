# activity-tracker
Tracks active visitors and provides summary reporting via REST API.

Getting Started:

Activity Tracker is a go-lang based application that accepts traffic on the default 3000 port. The application is also Go Modules based. Please ensure you have synchronized all required modules before attempting use.

To insert new activity metrics, please use the POST operation as shown below. Please note that summary reporting is based on matching the provided activity key.

curl --location --request POST 'localhost:3000/metric/some-activity-key' \
--header 'Content-Type: application/json' \
--data-raw '{"value": 30}'

To obtain the summary data for a given activity key over the past 1 hr, please use the GET operation as shown below. Please note activity metrics that are older than 1 hr are automatically evicted to reduce the application's memory footprint.

curl --location --request GET 'localhost:3000/metric/some-activity-key/sum'

Included are tests which verify functionality for initial state, adding values, obtaining summary counts, and activity eviction. The BDD test also serves as a unit test which can provide code coverage. To run the tests, please use the following command:

go test -cover