package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	ApiAddress = "http://localhost:8080"
)

func saveEntity() {
	entity := url.Values{
		"id":       []string{"id123"},
		"type":     []string{"runner"},
		"location": []string{"{\"latitude\": 51.516509, \"longitude\": 0.124615}"},
	}
	rsp, err := http.PostForm(ApiAddress+"/geo/location/save", entity)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp.Body.Close()
	b, _ := ioutil.ReadAll(rsp.Body)

	if rsp.StatusCode != 200 {
		fmt.Println("non 200 response", rsp.StatusCode, rsp.Status, string(b))
		return
	}

	fmt.Printf("Saved entity: %+v\n", entity)
}

func readEntity() {
	rsp, err := http.PostForm(ApiAddress+"/geo/location/read", url.Values{
		"id": []string{"id123"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp.Body.Close()

	b, _ := ioutil.ReadAll(rsp.Body)

	fmt.Printf("Read entity: %+v\n", string(b))
}

func searchForEntities() {
	rsp, err := http.PostForm(ApiAddress+"/geo/location/search", url.Values{
		"center":       []string{`{"latitude": 51.516509, "longitude": 0.124615}`},
		"type":         []string{"runner"},
		"radius":       []string{"500.0"},
		"num_entities": []string{"5"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp.Body.Close()

	b, _ := ioutil.ReadAll(rsp.Body)

	fmt.Printf("Search results: %+v\n", string(b))
}

func main() {
	saveEntity()
	readEntity()
	searchForEntities()
}
