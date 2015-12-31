package handler

import (
	"encoding/json"
	"strconv"
	"time"

	log "github.com/golang/glog"
	proto "github.com/micro/geo-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/server"
	api "github.com/micro/micro/api/proto"

	"golang.org/x/net/context"
)

var (
	topic = "geo.location"
)

func extractValue(pair *api.Pair) string {
	if pair == nil {
		return ""
	}
	if len(pair.Values) == 0 {
		return ""
	}
	return pair.Values[0]
}

func (l *Location) Save(ctx context.Context, req *api.Request, rsp *api.Response) error {
	var latlon map[string]float64
	err := json.Unmarshal([]byte(extractValue(req.Post["location"])), &latlon)
	if err != nil {
		return errors.BadRequest(server.DefaultOptions().Name+".search", "invalid location")
	}

	unix, _ := strconv.ParseInt(extractValue(req.Post["timestamp"]), 10, 64)

	entity := &proto.Entity{
		Id:   extractValue(req.Post["id"]),
		Type: extractValue(req.Post["type"]),
		Location: &proto.Point{
			Latitude:  latlon["latitude"],
			Longitude: latlon["longitude"],
			Timestamp: time.Unix(unix, 0).Unix(),
		},
	}

	if len(entity.Id) == 0 {
		return errors.BadRequest(server.DefaultOptions().Name+".save", "ID cannot be blank")
	}

	p := client.NewPublication(topic, entity)

	if err := client.Publish(ctx, p); err != nil {
		log.Errorf("Error publishing to topic %s: %v", topic, err)
		return errors.InternalServerError(server.DefaultOptions().Name+".save", err.Error())
	}

	log.Infof("Publishing entity ID %s", entity.Id)
	rsp.StatusCode = 200
	rsp.Body = `{}`
	return nil
}
