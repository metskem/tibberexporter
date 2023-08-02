### tibberexporter

A prometheus exporter that provides Tibber data (energy prices and more). It uses the [Tibber API](https://developer.tibber.com/docs/guides/getting-started) to get the data.

Tibber prices vary by the hour, and are provided daily at 13:00 (local time) for the next day.  
We initially wanted to use prometheus pushgateway to simply push the 24 or 48 metrics with their timestamp. But unfortunately the pushgateway does not support adding timestamps.  
That's why we ended up with a bit complicated setup, where we have 3 routines running:
* the http server handling the /metrics calls
* the routine that calls the Tibber API every 8 hours
* the routine that loops over the Tibber response and emits the metrics for the current hour

During startup, an API call to Tibber will be done to get the hourly data for today (and tomorrow if it starts between 13.00-23.59).  
Then another routine is started that loops over all observations (24 or 48), and the data for the current time (current hour), will be emitted (kept in memory).  
Make sure to run in the proper timezone, e.g. set environment variable TZ=Europe/Amsterdam

The provided metrics: 
* tibber_priceinfo_total - gauge, "Total price, combined energy and tax"
* tibber_priceinfo_energy - gauge, "Net energy price, without tax"
* tibber_priceinfo_tax - gauge, "Energy tax price"

