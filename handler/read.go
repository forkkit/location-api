package handler

import (
	"encoding/json"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/server"
	api "github.com/micro/micro/api/proto"
	loc "github.com/microhq/location-srv/proto/location"

	"golang.org/x/net/context"
)

type Location struct{}

func (l *Location) Read(ctx context.Context, req *api.Request, rsp *api.Response) error {
	id := extractValue(req.Post["id"])

	if len(id) == 0 {
		return errors.BadRequest(server.DefaultOptions().Name+".read", "Require Id")
	}

	request := client.NewRequest("go.micro.srv.location", "Location.Read", &loc.ReadRequest{
		Id: id,
	})

	response := &loc.ReadResponse{}

	err := client.Call(ctx, request, response)
	if err != nil {
		return errors.InternalServerError(server.DefaultOptions().Name+".read", "failed to read location")
	}

	b, _ := json.Marshal(response.Entity)

	rsp.StatusCode = 200
	rsp.Body = string(b)

	return nil
}
