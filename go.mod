module github.com/igeekinc/go-rocket

go 1.13

require (
	fyne.io/fyne/v2 v2.1.3
	github.com/adrianmo/go-nmea v1.2.0
	github.com/bskari/go-lsm303 v0.0.0-20200927082938-3432d22cb4f1
	github.com/d2r2/go-i2c v0.0.0-20191123181816-73a8a799d6bc
	github.com/d2r2/go-logger v0.0.0-20210606094344-60e9d1233e22 // indirect
	github.com/jacobsa/go-serial v0.0.0-20180131005756-15cf729a72d4
	github.com/pkg/errors v0.9.1
	periph.io/x/conn/v3 v3.6.10
	periph.io/x/devices/v3 v3.6.13
	periph.io/x/host/v3 v3.7.2
	periph.io/x/periph v3.6.8+incompatible // indirect
)

replace github.com/bskari/go-lsm303 => ../../bskari/go-lsm303

replace periph.io/x/conn/v3 => ../../../periph.io/x/conn
