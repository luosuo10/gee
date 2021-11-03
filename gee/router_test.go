package gee

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestParsePattern(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		pattern string
		results []string
	}{
		{
			pattern: "/p/:name",
			results:[]string{
				"p",
				":name",
			},
		}, {
			pattern: "/p/*aaa/bbb",
			results: []string {
				"p",
				"*aaa",
			},
		}, {
			pattern: "/p/*",
			results: []string {
				"p",
				"*",
			},
		},
	}

	for _, test := range tests {
		assert.Equal(parsePattern(test.pattern), test.results)
	}
}

func TestGetRoute(t *testing.T) {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)

	Convey("TestGetRoute node(/hello/b/d) is null\n", t, func ()  {
		node, _ := r.getRoute("GET", "/hello/b/d")
		ShouldBeNil(node, nil)
	})

	Convey("TestGetRoute node(/hello/b/c) is not null\n", t, func ()  {
		node, _ := r.getRoute("GET", "/hello/b/c")
		ShouldNotBeNil(node, nil)
	})

	Convey("TestGetRoute node(/hello/geektutu) should match /hello/:name\n", t, func ()  {
		node, _ := r.getRoute("GET", "/hello/geektutu")
		ShouldEqual(node.pattern, "/hello/:name")
	})

	Convey("TestGetRoute node(/hello/geektutu) params name should be equal to 'geektutu'\n", t, func ()  {
		_, params := r.getRoute("GET", "/hello/geektutu")
		ShouldEqual(params["name"], "geektutu")
	})

	Convey("TestGetRoute node(/hello/hel/llo)\n", t, func ()  {
		node, _ := r.getRoute("GET", "/hello/hel/llo")
		ShouldBeNil(node, nil)
	})

	Convey("TestGetRoute node(assets/zc/cz)\n", t, func ()  {
		node, _ := r.getRoute("GET", "/assets/zc/cz")
		ShouldEqual(node.pattern, "/assets/zc/cz")
	})
}