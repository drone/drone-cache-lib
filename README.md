# drone-cache-lib

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-cache/status.svg)](http://beta.drone.io/drone-plugins/drone-cache)
[![Coverage Status](https://aircover.co/badges/drone-plugins/drone-cache/coverage.svg)](https://aircover.co/drone-plugins/drone-cache)

drone-cache is a Go client library for creating cache [plugins](http://plugins.drone.io).

Download the package using `go get`:

```bash
go get "github.com/drone/drone-cache-lib"
```

Import the package:

```Go
import "github.com/drone/drone-cache-lib/cache"
```

The drone-cache library provides an interface for a `Storage` backend. When creating a new backend the following interface needs to be filled in.

```Go
type Storage interface {
	Get(p string, dst io.Writer) error
	Put(p string, src io.Reader) error
}
```

To create a `Cache` object using a `Storage` object:

```Go
cache, err := cache.New(storage)
```

To rebuild the cache:

```Go
err := cache.Rebuild(src, dst)
```

To restore the cache:

```Go
err := cache.Restore(src)
```

The drone-cache library currently supports the following file formats for cache storage
* .tar
