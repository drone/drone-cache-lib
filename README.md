# drone-cache-lib

[![Build Status](http://beta.drone.io/api/badges/drone/drone-cache-lib/status.svg)](http://beta.drone.io/drone/drone-cache-lib)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![Go Doc](https://godoc.org/github.com/drone/drone-cache-lib?status.svg)](http://godoc.org/github.com/drone/drone-cache-lib)
[![Go Report](https://goreportcard.com/badge/github.com/drone/drone-cache-lib)](https://goreportcard.com/report/github.com/drone/drone-cache-lib)

A Go client library for creating cache [plugins](http://plugins.drone.io).

## Usage

### Download the packe

```bash
go get -d github.com/drone/drone-cache-lib
```

### Import the package

```Go
import "github.com/drone/drone-cache-lib/cache"
```

### Create a `Cache` object

```Go
cache, err := cache.New(storage)
```

### To rebuild the cache

```Go
err := cache.Rebuild(src, dst)
```

### To restore the cache

```Go
err := cache.Restore(src)
```

### Supported archive formats

* .tar
