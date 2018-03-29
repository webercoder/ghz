package grpcannon

import (
	"context"
	"fmt"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)

type rpcStatsTagKey string

const rpcStatsID = rpcStatsTagKey("grpcannon_id")

// StatsHandler is for gRPC stats
type statsHandler struct {
	results chan *Result
	data    map[string]*statsData
}

type statsData struct {
	begin      time.Time
	outHeader  time.Time
	outPayload time.Time
	inHeader   time.Time
	inTrailer  time.Time
	inPayload  time.Time
}

// HandleConn handle the connection
func (c *statsHandler) HandleConn(ctx context.Context, cs stats.ConnStats) {
	// no-op
}

// TagConn exists to satisfy gRPC stats.Handler.
func (c *statsHandler) TagConn(ctx context.Context, cti *stats.ConnTagInfo) context.Context {
	// no-op
	return ctx
}

// HandleRPC implements per-RPC tracing and stats instrumentation.
func (c *statsHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	switch rs.(type) {
	case *stats.Begin:
		id, ok := ctx.Value(rpcStatsID).(string)
		if ok {
			rpcStats := rs.(*stats.Begin)
			c.data[id].begin = rpcStats.BeginTime
		}
	case *stats.OutHeader:
		id, ok := ctx.Value(rpcStatsID).(string)
		if ok {
			c.data[id].outHeader = time.Now()
		}
	case *stats.OutPayload:
		id, ok := ctx.Value(rpcStatsID).(string)
		if ok {
			rpcStats := rs.(*stats.OutPayload)
			c.data[id].outPayload = rpcStats.SentTime
		}
	case *stats.InHeader:
		id, ok := ctx.Value(rpcStatsID).(string)
		if ok {
			c.data[id].inHeader = time.Now()
		}
	case *stats.InTrailer:
		id, ok := ctx.Value(rpcStatsID).(string)
		if ok {
			c.data[id].inTrailer = time.Now()
		}
	case *stats.InPayload:
		id, ok := ctx.Value(rpcStatsID).(string)
		if ok {
			rpcStats := rs.(*stats.InPayload)
			c.data[id].inPayload = rpcStats.RecvTime
		}
	case *stats.End:
		id, ok := ctx.Value(rpcStatsID).(string)
		if ok {
			fmt.Printf("DATA: %+v\n", c.data[id])
		}
		rpcStats := rs.(*stats.End)
		end := time.Now()
		duration := end.Sub(rpcStats.BeginTime)

		var st string
		if rpcStats.Error != nil {
			s, ok := status.FromError(rpcStats.Error)
			if ok {
				st = s.Code().String()
			}
		}

		c.results <- &Result{rpcStats.Error, st, duration}
	}
}

// TagRPC implements per-RPC context management.
func (c *statsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	if info == nil {
		return ctx
	}

	idValue, err := uuid.GenerateUUID()
	if err == nil {
		c.data[idValue] = &statsData{}
		ctx = context.WithValue(ctx, rpcStatsID, idValue)

	}

	return ctx
}
