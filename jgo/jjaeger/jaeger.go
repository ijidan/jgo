package jjaeger

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"sync"
)

const LocalAgentHostPort = "192.168.33.10:6831"

//结构体
type JJaeger struct {
	tracer opentracing.Tracer
	closer io.Closer
	span   opentracing.Span
}

//初始化
func (j *JJaeger) init(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: LocalAgentHostPort,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("Error: connot init Jaeger: %v\n", err))
	}
	j.tracer = tracer
	j.closer = closer
	return tracer, closer
}

//创建span
func (j *JJaeger) CreateSpan(spanName string) opentracing.Span {
	j.span = j.tracer.StartSpan(spanName)
	return j.span
}

//实例
var instance *JJaeger
var once sync.Once

//获取单例
func NewJJaeger(service string, spanRootName string) (opentracing.Tracer, io.Closer, opentracing.Span) {
	once.Do(func() {
		instance = &JJaeger{}
		instance.init(service)
		instance.CreateSpan(spanRootName)
	})
	return instance.tracer, instance.closer, instance.span
}
