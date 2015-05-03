package main

import "github.com/mitchellh/packer/packer/plugin"
import "github.com/motymichaely/packer-post-processor-redis/redis"

func main() {
  server, err := plugin.Server()
  if err != nil {
    panic(err)
  }
  server.RegisterPostProcessor(new(redis.PostProcessor))
  server.Serve()
}