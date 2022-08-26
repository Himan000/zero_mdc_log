package zero_mdc_log

import "github.com/rs/zerolog"

type MDCHook struct{}

// 将MDC的<key,value>放到fileds
func (m MDCHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	if len(_globalFields) == 0 {
		return
	}

	fields := make(map[string]interface{})
	for _, field := range _globalFields {
		fields[field] = MDC().GetString(field)
	}

	e.Fields(fields)
}
