package handler

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"code.google.com/p/go.net/context"
	log "github.com/golang/glog"
	"github.com/myodc/geo-api/domain"
	"github.com/myodc/go-micro/broker"
	"github.com/myodc/go-micro/errors"
	"github.com/myodc/go-micro/server"
	api "github.com/myodc/micro/api/proto"
)

var (
	topic = "geo.location"
	once  sync.Once
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
	once.Do(func() {
		broker.Init()
		broker.Connect()
	})

	var latlon map[string]float64
	err := json.Unmarshal([]byte(extractValue(req.Post["location"])), &latlon)
	if err != nil {
		return errors.BadRequest(server.Name+".search", "invalid location")
	}

	unix, _ := strconv.ParseInt(extractValue(req.Post["timestamp"]), 10, 64)

	entity := &domain.Entity{
		ID:        extractValue(req.Post["id"]),
		Type:      extractValue(req.Post["type"]),
		Latitude:  latlon["latitude"],
		Longitude: latlon["longitude"],
		Timestamp: time.Unix(unix, 0).Unix(),
	}

	if len(entity.ID) == 0 {
		return errors.BadRequest(server.Name+".save", "ID cannot be blank")
	}

	data, err := json.Marshal(entity)
	if err != nil {
		log.Errorf("Error marshalling entity: %v", err)
		return errors.InternalServerError(server.Name+".save", err.Error())
	}

	if err := broker.Publish(ctx, topic, data); err != nil {
		log.Errorf("Error publishing to topic %s: %v", topic, err)
		return errors.InternalServerError(server.Name+".save", err.Error())
	}

	log.Infof("Publishing entity ID %s", entity.ID)
	rsp.StatusCode = 200
	rsp.Body = `{}`
	return nil
}
