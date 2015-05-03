# Packer Redis post processor 
[ ![Codeship Status for motymichaely/packer-post-processor-redis](https://codeship.com/projects/d2460c60-d231-0132-faa9-267aebe4cf02/status?branch=master)](https://codeship.com/projects/77470)

Packer Redis post processor is a [Packer post processor](https://packer.io/docs/extend/post-processor.html) plugin to store builder artifacts metadata into a Redis server for later retrieval.

This project was inspired by [Packer Consul post processor][packer-post-processor-consul].

## Usage

### Installation

You first need to install the plugin on your machine and let packer be aware of it.

Run:

```shell
$ go get github.com/motymichaely/packer-post-processor-redis
$ go install github.com/motymichaely/packer-post-processor-redis
```

Copy the binary file into your packer plugins folder:

```shell
$ mkdir $HOME/.packer.d/plugins
$ cp $GOPATH/bin/packer-post-processor-redis $HOME/.packer.d/plugins
```

Add the post-processor to Packer config file (`~/.packerconfig`):

```json
{
  "post-processors": {
    "redis": "packer-post-processor-redis"
  }
}
```

### Configuration

The configuration for this post-processor is extremely simple. 

`redis_url (string)` - Your access token for the Atlas API. This can be generated on your tokens page. Alternatively you can export your Redis URL as an environmental variable and remove it from the configuration. Note 

`key_prefix (string)` - The prefix of the key to be set, i.e `my_service/images/`. for builders that support multiple regions, the prefix will be appended with the region to the key prefix, i.e `my_service/images/region`.

Add the post-processor to your packer template:

```json
{
    "builders": [{
        "type": "amazon-ebs",
        "region": "us-east-1",
        "ami_regions": ["us-east-1", "us-west-1"],
        ...
        }]
    "post-processors": [
      {
        "type": "redis",
        "redis_url": "redis://localhost:6379/0",
        "key_prefix": "my_service/images/",
        "only": ["amazonebs"]
      }
    ]
}
```

This example would take each artifact that was built by the "amazonebs" builder and set the keys `my_service/images/us-east-1` and `my_service/images/us-west-1` with a value of the image id, i.e `ami_5123412`.

## Building

Build & Test:
```shell
$ make all
```  

Build:  

```shell
$ make bin
```

Test:  

```shell
$ make test
```

## Contribute

See the [CONTRIBUTING](CONTRIBUTING.md) file for contribution guidelines.

## License

Copyright (c) 2015 Moty Michaely

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).

[packer-post-processor-consul]: https://github.com/bhourigan/packer-post-processor-consul