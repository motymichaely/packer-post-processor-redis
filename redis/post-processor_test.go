package redis

import (
	"bytes"
	"context"

	"github.com/hashicorp/packer/builder/amazon/ebs"
	"github.com/hashicorp/packer/builder/amazon/instance"
	"github.com/hashicorp/packer/packer"
	"github.com/rafaeljusto/redigomock"

	"testing"
)

func testConfig() map[string]interface{} {
	return map[string]interface{}{}
}

func testPP(t *testing.T) *PostProcessor {
	var p PostProcessor
	if err := p.Configure(testConfig()); err != nil {
		t.Fatalf("err: %s", err)
	}

	return &p
}

func testUi() *packer.BasicUi {
	return &packer.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	}
}

func TestPostProcessor_ImplementsPostProcessor(t *testing.T) {
	var _ packer.PostProcessor = new(PostProcessor)
}

func TestPostProcessor_PostProcess(t *testing.T) {
	conn := redigomock.NewConn()
	p := &PostProcessor{client: conn}
	if err := p.Configure(validDefaults()); err != nil {
		t.Fatalf("err: %s", err)
	}

	artifact := &packer.MockArtifact{
		BuilderIdValue: ebs.BuilderId,
		IdValue:        "us-east-1:ami_12345",
	}

	conn.Command("SET", "my_prefix/us-east-1", "ami_12345").Expect("ok")

	result, keep, forceOverride, err := p.PostProcess(context.Background(), testUi(), artifact)
	if result != artifact {
		t.Fatal("should not return given artifact")
	}
	if keep == false {
		t.Fatal("should keep")
	}
	if forceOverride == true {
		t.Fatal("should not force override")
	}
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcessor_PostProcess_invalidDefaults(t *testing.T) {
	p := &PostProcessor{client: redigomock.NewConn()}
	if err := p.Configure(invalidDefaults()); err == nil {
		t.Fatalf("should error for missing required configuration: %s", err)
	}
}

func TestPostProcessor_PostProcess_validDefaults(t *testing.T) {
	p := &PostProcessor{client: redigomock.NewConn()}
	if err := p.Configure(validDefaults()); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcessor_PostProcess_amazonebs(t *testing.T) {
	conn := redigomock.NewConn()
	p := &PostProcessor{client: conn}
	if err := p.Configure(validDefaults()); err != nil {
		t.Fatalf("err: %s", err)
	}
	artifact := &packer.MockArtifact{
		BuilderIdValue: ebs.BuilderId,
		IdValue:        "us-east-1:ami_12345",
	}

	conn.Command("SET", "my_prefix/us-east-1", "ami_12345").Expect("ok")

	result, keep, forceOverride, err := p.PostProcess(context.Background(), testUi(), artifact)
	if result != artifact {
		t.Fatal("should not return given artifact")
	}
	if keep == false {
		t.Fatal("should keep")
	}
	if forceOverride == true {
		t.Fatal("should not force override")
	}
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcessor_PostProcess_amazoninstance(t *testing.T) {
	conn := redigomock.NewConn()
	p := &PostProcessor{client: conn}
	if err := p.Configure(validDefaults()); err != nil {
		t.Fatalf("err: %s", err)
	}
	artifact := &packer.MockArtifact{
		BuilderIdValue: instance.BuilderId,
		IdValue:        "us-east-1:ami_12345",
	}

	conn.Command("SET", "my_prefix/us-east-1", "ami_12345").Expect("ok")

	result, keep, forceOverride, err := p.PostProcess(context.Background(),testUi(), artifact)
	if result != artifact {
		t.Fatal("should not return given artifact")
	}
	if keep == false {
		t.Fatal("should keep")
	}
	if forceOverride == true {
		t.Fatal("should not force override")
	}
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcessor_PostProcess_amazoninstanceMultiRegion(t *testing.T) {
	conn := redigomock.NewConn()
	p := &PostProcessor{client: conn}
	if err := p.Configure(validDefaults()); err != nil {
		t.Fatalf("err: %s", err)
	}
	artifact := &packer.MockArtifact{
		BuilderIdValue: instance.BuilderId,
		IdValue:        "us-east-1:ami_12345,us-west-1:ami_12345",
	}

	conn.Command("SET", "my_prefix/us-east-1", "ami_12345").Expect("ok")
	conn.Command("SET", "my_prefix/us-west-1", "ami_12345").Expect("ok")

	result, keep, forceOverride, err := p.PostProcess(context.Background(), testUi(), artifact)
	if result != artifact {
		t.Fatal("should not return given artifact")
	}
	if keep == false {
		t.Fatal("should keep")
	}
	if forceOverride == true {
		t.Fatal("should not force override")
	}
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcessor_PostProcess_googlecompute(t *testing.T) {
	conn := redigomock.NewConn()
	p := &PostProcessor{client: conn}
	if err := p.Configure(validDefaults()); err != nil {
		t.Fatalf("err: %s", err)
	}
	artifact := &packer.MockArtifact{
		BuilderIdValue: "packer.googlecompute",
		IdValue:        "image-name-12345",
	}

	conn.Command("SET", "my_prefix", "image-name-12345").Expect("ok")

	result, keep, forceOverride, err := p.PostProcess(context.Background(), testUi(), artifact)
	if result != artifact {
		t.Fatal("should not return given artifact")
	}
	if keep == false {
		t.Fatal("should keep")
	}
	if forceOverride == true {
		t.Fatal("should not force override")
	}
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcessor_PostProcess_docker(t *testing.T) {
	conn := redigomock.NewConn()
	p := &PostProcessor{client: conn}
	if err := p.Configure(validDefaults()); err != nil {
		t.Fatalf("err: %s", err)
	}
	artifact := &packer.MockArtifact{
		BuilderIdValue: "packer.docker",
		IdValue:        "sha256:image-name-12345",
	}

	conn.Command("SET", "my_prefix/sha256", "image-name-12345").Expect("ok")

	result, keep, forceOverride, err := p.PostProcess(context.Background(), testUi(), artifact)
	if result != artifact {
		t.Fatal("should not return given artifact")
	}
	if keep == false {
		t.Fatal("should keep")
	}
	if forceOverride == true {
		t.Fatal("should not force override")
	}
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestPostProcessor_PostProcess_docker_import(t *testing.T) {
	conn := redigomock.NewConn()
	p := &PostProcessor{client: conn}
	if err := p.Configure(validDefaults()); err != nil {
		t.Fatalf("err: %s", err)
	}
	artifact := &packer.MockArtifact{
		BuilderIdValue: "packer.post-processor.docker-import",
		IdValue:        "sha256:image-name-12345",
	}

	conn.Command("SET", "my_prefix", "image-name-12345").Expect("ok")

	result, keep, forceOverride, err := p.PostProcess(context.Background(), testUi(), artifact)
	if result != artifact {
		t.Fatal("should not return given artifact")
	}
	if keep == false {
		t.Fatal("should keep")
	}
	if forceOverride == true {
		t.Fatal("should not force override")
	}
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func validDefaults() map[string]interface{} {
	return map[string]interface{}{
		"redis_url":  ":6379",
		"key_prefix": "my_prefix",
	}
}

func invalidDefaults() map[string]interface{} {
	return map[string]interface{}{}
}
