BIN = falcon-agent aggregator graph hbs judge nodata query sender task transfer
TARGET = open-falcon

all: $(BIN)
	mkdir -p bin
	go build -o open-falcon

falcon-agent:
	go build -o ./bin/$@ github.com/cepave/open-falcon/modules/agent
aggregator:
	go build -o ./bin/$@ github.com/cepave/open-falcon/modules/aggregator 
graph:
	go build -o ./bin/$@ github.com/cepave/open-falcon/modules/graph
hbs:
	go build -o ./bin/$@ github.com/cepave/open-falcon/modules/hbs
judge:
	go build -o ./bin/$@ github.com/cepave/open-falcon/modules/judge
#links:
#	go build -o ./bin/$@ github.com/cepave/open-falcon/modules/links
nodata:
	go build -o ./bin/$@ github.com/cepave/open-falcon/modules/nodata
query:
	go build -o ./bin/$@ github.com/cepave/open-falcon/modules/query
sender:
	go build -o ./bin/$@ github.com/cepave/open-falcon/modules/sender
task:
	go build -o ./bin/$@ github.com/cepave/open-falcon/modules/task
transfer:
	go build -o ./bin/$@ github.com/cepave/open-falcon/modules/transfer

clean:
	rm -rf ./bin
	rm -rf ./$(TARGET)
