package charm_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
	. "launchpad.net/gocheck"
	"launchpad.net/juju/go/charm"
	"launchpad.net/goyaml"
	"os"
	"path/filepath"
)

func Test(t *testing.T) {
	TestingT(t)
}

type S struct{}

var _ = Suite(&S{})

func (s *S) TestParseId(c *C) {
	namespace, name, rev, err := charm.ParseId("local:mysql-21")
	c.Assert(err, IsNil)
	c.Assert(namespace, Equals, "local")
	c.Assert(name, Equals, "mysql")
	c.Assert(rev, Equals, 21)

	namespace, name, rev, err = charm.ParseId("local:mysql-cluster-21")
	c.Assert(err, IsNil)
	c.Assert(namespace, Equals, "local")
	c.Assert(name, Equals, "mysql-cluster")
	c.Assert(rev, Equals, 21)

	_, _, _, err = charm.ParseId("foo")
	c.Assert(err, Matches, `Missing charm namespace: "foo"`)

	_, _, _, err = charm.ParseId("local:foo-x")
	c.Assert(err, Matches, `Missing charm revision: "local:foo-x"`)
}

func checkDummy(c *C, f charm.Charm, path string) {
	c.Assert(f.Meta().Name, Equals, "dummy")
	c.Assert(f.Config().Options["title"].Default, Equals, "My Title")
	switch f := f.(type) {
	case *charm.Bundle:
		c.Assert(f.Path, Equals, path)
	case *charm.Dir:
		c.Assert(f.Path, Equals, path)
		_, err := os.Stat(filepath.Join(path, "src", "hello.c"))
		c.Assert(err, IsNil)
	}
}

type YamlHacker map[interface{}]interface{}

func ReadYaml(r io.Reader) YamlHacker {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	m := make(map[interface{}]interface{})
	err = goyaml.Unmarshal(data, m)
	if err != nil {
		panic(err)
	}
	return YamlHacker(m)
}

func (yh YamlHacker) Reader() io.Reader {
	data, err := goyaml.Marshal(yh)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(data)
}