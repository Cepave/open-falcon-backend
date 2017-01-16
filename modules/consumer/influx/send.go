package influx

import (
	"encoding/json"
	"time"

	cutils "github.com/Cepave/common/utils"
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
		log.Println("Influx create batch points error: ", err)
	}

	// Create a point and add to batch from msg
	v := metricValue{}
	err = json.Unmarshal(msg, &v)

	tags := cutils.DictedTagstring(v.Tags)
	tags["endpoint"] = v.Endpoint
	fields := map[string]interface{}{
		"idle":   10.1,
		"system": 53.3,
		"user":   46.6,
	}
	fields["value"] = v.Value

	pt, err := client.NewPoint(
		v.Metric,
		tags,
		fields,
		time.Unix(v.Timestamp, 0))

	if err != nil {
		log.Println("Influx marshalling Error: ", err)
	}

	bp.AddPoint(pt)

	// Write the batch
	err = conn.Write(bp)
	if err != nil {
		log.Println("Influx write to database error: ", err)
		// if client cannot write, close it.
		conn.Close()
		conn = nil
	}
}
