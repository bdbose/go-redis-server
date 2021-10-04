# go-redis-server

Making a caching server for my [NewsToday](https://newstoday.tech/)

Earlier the response time of the heavy NodeJs Server is around 1s and now as I use redis to cache the response as the response does not change for 18mins.

Using caching response time came to an average of 5ms which is 95% less than the previous response and a huge benefit for the frontend and SEO.

