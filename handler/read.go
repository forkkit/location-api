package handler

import (
	"code.google.com/p/go.net/context"
	"encoding/json"

	read "github.com/myodc/geo-srv/proto/location/read"
	"github.com/myodc/go-micro/client"
	"github.com/myodc/go-micro/errors"
	"github.com/myodc/go-micro/server"
	api "github.com/myodc/micro/api/proto"
)

type Location struct{}

func (l *Location) Read(ctx context.Context, req *api.Request, rsp *api.Response) error {
	id := extractValue(req.Post["id"])

	if len(id) == 0 {
		return errors.BadRequest(server.Name+".read", "Require Id")
	}

	request := client.NewRequest("go.micro.srv.geo", "Location.Read", &read.Request{
		Id: id,
	})

	response := &read.Response{}

	err := client.Call(request, response)
	if err != nil {
		return errors.InternalServerError(server.Name+".read", "failed to read location")
	}

	b, _ := json.Marshal(response.Entity)

	rsp.StatusCode = 200
	rsp.Body = string(b)

	return nil
}
