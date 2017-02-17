package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestNewConfig(t *testing.T) {
	c, err := NewConfig()

	if err != nil {
		t.Errorf("cannot create new config: %q", err)
	}

	ex, _ := os.Executable()

	if c.Root != path.Dir(ex) {
		t.Errorf("invalid root, expected = %q; got = %q", path.Dir(ex), c.Root)
	}
}

func TestNewConfigFromYaml(t *testing.T) {
	b, err := ioutil.ReadFile("./testdata/test_config.yml")
	if err != nil {
		t.Errorf("cannot read test config: %q", err)
	}

	c, err := NewConfigFromYaml(b)
	if err != nil {
		t.Errorf("cannot create new config: %q", err)
	}

	if c.Listen != "0.0.0.0:6029" {
		t.Errorf(`invalid config for "listen", expected = "0.0.0.0:6029"; got = %q`, c.Listen)
	}

	if len(c.Files) != 1 {
		t.Errorf(`invalid config for "files", expected = ["/foo.html"]; got = %q`, c.Files)
	}
}

func TestConfigIsHTTPS(t *testing.T) {
	b, err := ioutil.ReadFile("./testdata/test_config_https.yml")
	if err != nil {
		t.Errorf("cannot read test config: %q", err)
	}

	c, err := NewConfigFromYaml(b)
	if err != nil {
		t.Errorf("cannot create new config: %q", err)
	}

	if !c.IsHTTPS() {
		t.Error(`invalid config for "https", expected = "true"; got = "false"`)
	}
}

func TestConfigGetConfigForPath(t *testing.T) {
	b, err := ioutil.ReadFile("./testdata/test_config.yml")
	if err != nil {
		t.Errorf("cannot read test config: %q", err)
	}

	c, err := NewConfigFromYaml(b)
	if err != nil {
		t.Errorf("cannot create new config: %q", err)
	}

	f, ok := c.Files["/foo.html"]
	if !ok {
		t.Error(`invalid config for "files", expected "/foo.html" to be present`)
	}

	if !f.Cache {
		t.Error(`invalid config for "files[/foo.html].cache", expected = "true"; got = "false"`)
	}

	if f.Gzip {
		t.Error(`invalid config for "files[/foo.html].gzip", expected = "false"; got = "true"`)
	}

	if len(f.Push) != 1 || f.Push[0] != "/assets/js/app.js" {
		t.Errorf(`invalid config for "files[/foo.html].push", expected = ["/assets/js/app.js"]; got = %q`, f.Push)
	}
}
