all: rocket-tele rocket-tracker

rocket-tele:
	cd cmd/rocket-tele; go install

rocket-tracker:
	cd cmd/rocket-tracker; go install
