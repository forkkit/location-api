package handler

import (
	"code.google.com/p/go.net/context"
	"encoding/json"
	common "github.com/myodc/geo-srv/proto"
	search "github.com/myodc/geo-srv/proto/location/search"
	"github.com/myodc/go-micro/client"
	"github.com/myodc/go-micro/errors"
	"github.com/myodc/go-micro/server"
	api "github.com/myodc/micro/api/proto"
	"strconv"
)

func (l *Location) Search(ctx context.Context, req *api.Request, rsp *api.Response) error {
	radius, _ := strconv.ParseFloat(extractValue(req.Post["radius"]), 64)
	typ := extractValue(req.Post["type"])
	entities, _ := strconv.ParseInt(extractValue(req.Post["num_entities"]), 10, 64)

	var latlon map[string]float64
	err := json.Unmarshal([]byte(extractValue(req.Post["center"])), &latlon)
	if err != nil {
		return errors.BadRequest(server.Name+".search", "invalid center point")
	}

	if len(typ) == 0 {
		return errors.BadRequest(server.Name+".search", "type cannot be blank")
	}

	if entities == 0 {
		return errors.BadRequest(server.Name+".search", "num_entities must be greater than 0")
	}

	request := client.NewRequest("go.micro.srv.geo", "Location.Search", &search.Request{
		Center: &common.Location{
			Latitude:  latlon["latitude"],
			Longitude: latlon["longitude"],
		},
		Radius:      radius,
		NumEntities: entities,
		Type:        typ,
	})

	response := &search.Response{}

	err = client.Call(request, response)
	if err != nil {
		return errors.InternalServerError(server.Name+".search", "could not retrieve results")
	}

	b, _ := json.Marshal(response.Entities)
	rsp.StatusCode = 200
	rsp.Body = string(b)
	return nil
}
