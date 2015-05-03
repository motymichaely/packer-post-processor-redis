package redis

import (
  "fmt"

  "github.com/garyburd/redigo/redis"
  "github.com/mitchellh/packer/common"
  "github.com/mitchellh/packer/packer"

  "net/url"

  "strings"
)

var builtins = map[string]string{
  "mitchellh.amazonebs": "amazonebs",
  "mitchellh.amazon.instance": "amazoninstance",
  "packer.googlecompute": "googlecompute",
}

type Config struct {
  common.PackerConfig `mapstructure:",squash"`

  RedisUrl  string `mapstructure:"redis_url"`
  KeyPrefix string `mapstructure:"key_prefix"`

  tpl *packer.ConfigTemplate
}

type PostProcessor struct {
  client redis.Conn

  config Config
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
  _, err := common.DecodeConfig(&p.config, raws...)
  if err != nil {
    return err
  }

  p.config.tpl, err = packer.NewConfigTemplate()
  if err != nil {
    return err
  }
  p.config.tpl.UserVars = p.config.PackerUserVars

  // Accumulate any errors
  errs := new(packer.MultiError)

  // Process templates
  templates := map[string]*string{
    "redis_url":  &p.config.RedisUrl,
    "key_prefix": &p.config.KeyPrefix,
  }

  for n, ptr := range templates {
    var err error
    *ptr, err = p.config.tpl.Process(*ptr, nil)
    if err != nil {
      errs = packer.MultiErrorAppend(
        errs, fmt.Errorf("Error processing %s: %s", n, err))
    }
  }

  required := map[string]*string{
    "redis_url":  &p.config.RedisUrl,
    "key_prefix": &p.config.KeyPrefix,
  }

  for key, ptr := range required {
    if *ptr == "" {
      errs = packer.MultiErrorAppend(
        errs, fmt.Errorf("%s must be set", key))
    }
  }


  if len(errs.Errors) > 0 {
    return errs
  }

  return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
  _, ok := builtins[artifact.BuilderId()]
  if !ok {
    return artifact, false, fmt.Errorf(
      "Unsupported artifact type: %s", artifact.BuilderId())
  }

  ui.Say("Putting build artifacts into Redis...")

  ui.Message("Opening Redis connection...")

  // If no client is set, then we create a new one
  if p.client == nil {
    err := fmt.Errorf(
        "Error connecting to Redis")

    redisURL, err := url.Parse(p.config.RedisUrl)

    if err != nil {
      return artifact, false, fmt.Errorf(
        "Error parsing RedisUrl: %s", err)
    }

    auth := ""

    if redisURL.User != nil {
      if password, ok := redisURL.User.Password(); ok {
        auth = password
      }
    }

    p.client, err = redis.Dial(
      "tcp",
      redisURL.Host)
    if err != nil {
      return artifact, false, fmt.Errorf(
        "Error connecting to Redis: %s", err)
    }

    if len(auth) > 0 {
      _, err = p.client.Do("AUTH", auth)

      if err != nil {
        return artifact, false, fmt.Errorf(
          "Error connecting to Redis: %s", err)
      }
    }

    defer func() {
      ui.Message("Closing Redis connection...")
      if err := p.client.Close(); err != nil {
        ui.Error(fmt.Sprintf("Error closing Redis connection: %s", err))
      }
      }()
  }
  
  for _, regions := range strings.Split(artifact.Id(), ",") {
    
    parts := strings.Split(regions, ":")

    redis_key := ""
    image_id := ""

    if len(parts) == 2 {
      region :=   parts[0]
      image_id = parts[1]

      redis_key = fmt.Sprintf("%s/%s", p.config.KeyPrefix, region)
    } else {
      image_id = parts[0]

      redis_key = fmt.Sprintf("%s", p.config.KeyPrefix)
    }

    ui.Message(fmt.Sprintf("Setting key %s with value %s", redis_key, image_id))

    if _, err := p.client.Do("SET", redis_key, image_id); err != nil {
      return artifact, false, err
    }
  }
  return artifact, false, nil
}