package influx

import (
	"encoding/json"
	"time"

	"github.com/Cepave/open-falcon-backend/common/model"
	cutils "github.com/Cepave/open-falcon-backend/common/utils"
	"github.com/Cepave/open-falcon-backend/modules/consumer/g"
	log "github.com/Sirupsen/logrus"
	"github.com/influxdata/influxdb/client/v2"
)

var conn client.Client = nil

func Send(msg []byte) {
	for conn == nil {
		// Make client
		var err error
		conn, err = client.NewHTTPClient(client.HTTPConfig{
			Addr:     g.Config().Influx.Addr,
			Username: g.Config().Influx.User,
			Password: g.Config().Influx.Pass,
		})

		if err != nil {
			log.Println("Influx create http client Error: ", err)
		}
	}
	transfer(msg)
}

func transfer(msg []byte) {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  g.Config().Influx.DB,
		Precision: "s",
	})

	if err != nil {
		log.Errorf("Influx create batch points error: %v", err)
	}

	// Create a point and add to batch from msg
	v := model.MetricValueExtend{}
	err = json.Unmarshal(msg, &v)

	tags := cutils.DictedTagstring(v.Tags)
	tags["endpoint"] = v.Endpoint
	fields := cutils.DictedFieldstring(v.Fields)
	fields["value"] = v.Value

	pt, err := client.NewPoint(
		v.Metric,
		tags,
		fields,
		time.Unix(v.Timestamp, 0))

	if err != nil {
		log.Printf("Influx marshalling Error: %v", err)
	}

	log.Debugf("send datapoint %s to influx database", pt)
	bp.AddPoint(pt)

	// Write the batch
	err = conn.Write(bp)
	if err != nil {
		log.Errorf("Influx write to database error: %v", err)
		duration, version, err := conn.Ping(time.Duration(10) * time.Second)
		log.Printf("Influx ping result -> duration: %v, version: %v", duration, version)
		if err != nil {
			// if client ping error, then close connection.
			log.Errorf("Influx ping database error: %v", err)
			conn.Close()
			conn = nil
		}
	}
}
