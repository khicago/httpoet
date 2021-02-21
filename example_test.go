package httpoet_test

import (
	"fmt"

	"github.com/khicago/httpoet"
)

func ExampleGet() {
	poet := httpoet.New("https://www.your_awsome_site.com")
	result := poet.Get("relative/path")
	body, err := result.Body()
	if err == nil {
		panic(err)
	}
	fmt.Println(string(body))
}

func ExampleResult() {
	poet := httpoet.New("https://www.your_awsome_site.com")
	result := poet.Get("relative/path")
	v := httpoet.D{}
	if err := result.ParseJson(v); err == nil {
		panic(err)
	}
	fmt.Println(v)
}

func ExampleHeader() {
	h := httpoet.H{"key": "val"}
	poet := httpoet.New("https://www.your_awsome_site.com").SetBaseH(h)

	hs := httpoet.Hs{"key": {"v0", "v1"}}
	result := poet.Get("relative/path", httpoet.OSetHeaders(hs))

	v := httpoet.D{}
	if err := result.ParseJson(v); err == nil {
		panic(err)
	}
	fmt.Println(v)
}
