package redis

import (
	"context"
	"fmt"
	"strings"
	"net/url"

	"github.com/garyburd/redigo/redis"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

var builtins = map[string]string{
	"mitchellh.amazonebs":                 "amazonebs",
	"mitchellh.amazon.instance":           "amazoninstance",
	"packer.googlecompute":                "googlecompute",
	"packer.docker":                       "docker",
	"packer.post-processor.docker-import": "docker",
	"packer.post-processor.docker-tag":    "docker",
	"Azure.ResourceManagement.VMImage":    "azure-arm",
}

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	RedisUrl  string `mapstructure:"redis_url"`
	KeyPrefix string `mapstructure:"key_prefix"`
	ImageId   string `mapstructure:"image_id"`

	ctx interpolate.Context
}

type PostProcessor struct {
	client redis.Conn

	config Config
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate: true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return err
	}

	required := map[string]*string{
		"redis_url":  &p.config.RedisUrl,
		"key_prefix": &p.config.KeyPrefix,
	}

	var errs *packer.MultiError
	for key, ptr := range required {
		if *ptr == "" {
			errs = packer.MultiErrorAppend(
				errs, fmt.Errorf("%s must be set", key))
		}
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (p *PostProcessor) PostProcess(ctx context.Context, ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, bool, error) {
	_, ok := builtins[artifact.BuilderId()]
	if !ok {
		return nil, true, false, fmt.Errorf(
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
			return artifact, true, false, fmt.Errorf(
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
			return artifact, true, false, fmt.Errorf(
				"Error connecting to Redis: %s", err)
		}

		if len(auth) > 0 {
			_, err = p.client.Do("AUTH", auth)

			if err != nil {
				return artifact, true, false, fmt.Errorf(
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
			if artifact.BuilderId() == "packer.post-processor.docker-import" {
				// packer.post-processor.docker-import image ID includes the SHA256 part. Ignore it.
				redis_key = fmt.Sprintf("%s", p.config.KeyPrefix)
			} else {
				region := parts[0]
				redis_key = fmt.Sprintf("%s/%s", p.config.KeyPrefix, region)
			}
			
			if len(p.config.ImageId) > 0 {
				image_id = p.config.ImageId
			} else { 
				image_id = parts[1]
			}
			
		} else {
			if len(p.config.ImageId) > 0 {
				image_id = p.config.ImageId
			} else { 
				image_id = parts[0]
			}

			redis_key = fmt.Sprintf("%s", p.config.KeyPrefix)
		}

		ui.Message(fmt.Sprintf("Setting key %s with value %s", redis_key, image_id))

		if _, err := p.client.Do("SET", redis_key, image_id); err != nil {
			return artifact, true, false, err
		}
	}
	return artifact, true, false, nil
}
