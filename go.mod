module github.com/xBlaz3kx/ChargePi-go

go 1.16

replace (
	github.com/xBlaz3kx/ChargePi-go/components/cache => ./components/cache
	github.com/xBlaz3kx/ChargePi-go/components/connector => ./components/connector
	github.com/xBlaz3kx/ChargePi-go/components/connector-manager => ./components/connector-manager
	github.com/xBlaz3kx/ChargePi-go/components/hardware => ./components/hardware
	github.com/xBlaz3kx/ChargePi-go/components/hardware/display => ./components/hardware/display
	github.com/xBlaz3kx/ChargePi-go/components/hardware/display/i18n => ./components/hardware/display/i18n
	github.com/xBlaz3kx/ChargePi-go/components/hardware/indicator => ./components/hardware/indicator
	github.com/xBlaz3kx/ChargePi-go/components/hardware/power-meter => ./components/hardware/power-meter
	github.com/xBlaz3kx/ChargePi-go/components/hardware/reader => ./components/hardware/reader
	github.com/xBlaz3kx/ChargePi-go/components/scheduler => ./components/scheduler
	github.com/xBlaz3kx/ChargePi-go/components/settings => ./components/settings
	github.com/xBlaz3kx/ChargePi-go/components/settings/conf-manager => ./components/settings/conf-manager
	github.com/xBlaz3kx/ChargePi-go/components/settings/settings-manager => ./components/settings/settings-manager

	github.com/xBlaz3kx/ChargePi-go/data => ./data
	github.com/xBlaz3kx/ChargePi-go/data/auth => ./data/auth
	github.com/xBlaz3kx/ChargePi-go/data/ocpp => ./data/ocpp
	github.com/xBlaz3kx/ChargePi-go/data/session => ./data/session
	github.com/xBlaz3kx/ChargePi-go/data/settings => ./data/settings
)

require (
	github.com/Graylog2/go-gelf v0.0.0-20170811154226-7ebf4f536d8f
	github.com/agrison/go-commons-lang v0.0.0-20200208220349-58e9fcb95174
	github.com/avast/retry-go v3.0.0+incompatible
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/clausecker/nfc/v2 v2.1.4
	github.com/d2r2/go-hd44780 v0.0.0-20181002113701-74cc28c83a3e
	github.com/d2r2/go-i2c v0.0.0-20191123181816-73a8a799d6bc
	github.com/d2r2/go-logger v0.0.0-20210606094344-60e9d1233e22 // indirect
	github.com/go-co-op/gocron v1.6.0
	github.com/kkyr/fig v0.3.0
	github.com/lorenzodonini/ocpp-go v0.14.0
	github.com/nicksnyder/go-i18n/v2 v2.1.2
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/reactivex/rxgo/v2 v2.5.0
	github.com/rpi-ws281x/rpi-ws281x-go v1.0.8
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/teivah/onecontext v1.3.0 // indirect
	github.com/warthog618/gpiod v0.6.0
	golang.org/x/text v0.3.3
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
