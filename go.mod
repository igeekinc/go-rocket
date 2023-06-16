module github.com/igeekinc/go-rocket

go 1.13

require (
	fyne.io/fyne/v2 v2.3.5
	github.com/adrianmo/go-nmea v1.7.0
	github.com/bskari/go-lsm303 v0.0.0-20200927082938-3432d22cb4f1
	github.com/igeekinc/go-gps-i2c v0.0.0-20221002190116-ae4bfee6c7f5
	github.com/jacobsa/go-serial v0.0.0-20180131005756-15cf729a72d4
	github.com/pkg/errors v0.9.1
	periph.io/x/conn/v3 v3.7.0
	periph.io/x/devices/v3 v3.7.0
	periph.io/x/host/v3 v3.8.0
)

replace github.com/bskari/go-lsm303 => ../../bskari/go-lsm303

replace periph.io/x/conn/v3 => ../../../periph.io/x/conn

replace github.com/igeekinc/go-gps-i2c => ../go-gps-i2c
