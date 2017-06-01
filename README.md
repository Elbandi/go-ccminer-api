# CCMiner API for Go #

## Installation ##

    # install the library:
    go get github.com/Elbandi/go-ccminer-api

    // Use in your .go code:
    import (
        "github.com/Elbandi/go-ccminer-api"
    )

## API Documentation ##

Full godoc output from the latest code in master is available here:

http://godoc.org/github.com/Elbandi/go-ccminer-api

## Quickstart ##

```go
package main

import (
    "github.com/Elbandi/go-ccminer-api"
    "log"
)

func main() {
    miner := cgminer.New("localhost", 4068)
    summary, err := miner.Summary()
    if err != nil {
        log.Fatalln("Unable to connect to CCMiner: ", err)
    }

    log.Printf("Average Hashrate: %f KH/s\n", summary.Khs)
}
```