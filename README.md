# icecli

## Installation
Grab the appropriate archive from the releases section of this repo, extract it, move the included binary somewhere on your $PATH and make it executable.

## Usage
In order to use the tool you have to be on a German ICE train and connected to the onboard WiFi SSID `WIFIonICE`.

```bash
# Get help
$ icecli --help

# Show your status information about your train
$ icecli status

# Show information about the train trip
$ icecli trip

# Show information about the stops on the trip
$ icecli trip stops

# Override destination information with your stop
$ icecli trip -d "Berlin Hbf (tief)"

# Change output format to CSV
$ icecli trip -d "Berlin Hbf (tief)" -o csv

# Change output format to CSV and filter for only ETA to your destination override
$ icecli trip -d "Berlin Hbf (tief)" -o csv -f "ARRIVING"
```
