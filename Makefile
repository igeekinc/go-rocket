export GOOS=linux
export GOARCH=arm

all: rocket-tele rocket-tracker

rocket-tele:
	cd cmd/rocket-tele; go build

rocket-tracker:
	cd cmd/rocket-tracker; go build
