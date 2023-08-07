# structdefaults
通过tag来设置struct的默认值

useage:

```shell
go get -u github.com/foxliu/structdefaults@v0.2.0
```

example:

define struct

```go
type LogConfig struct {
	Level      string   `json:"level,omitempty" default:"info"`
	Filename   string   `json:"filename,omitempty" default:"logs/app.log"`
	MaxSize    int      `json:"max_size,omitempty" default:"100"`
	MaxAge     int      `json:"max_age,omitempty" default:"5"`
	MaxBackups int      `json:"max_backups,omitempty" default:"5"`
	TestArry   []string `json:"test_array,omitempty" default:"[\"start\"]"`
}
```

use:

```go
cfg := LogConfig{}
err := setDefaults.SetStructDefaults(&cfg)
log.Printf("Config: %+v", cfg)
```

