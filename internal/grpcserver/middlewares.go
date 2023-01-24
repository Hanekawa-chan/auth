package grpcserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"time"
)

func recoveryHandler(ctx context.Context, p interface{}) (err error) {
	log.Error().
		Str("panic", fmt.Sprintf("%+v", p)).
		Str("ctx", fmt.Sprintf("%+v", ctx)).
		Msg("PANIC")
	if err, ok := p.(error); ok {
		return fmt.Errorf("panic: %w", err)
	}
	return fmt.Errorf("panic: %+v", p)
}

func logIncomingRequestsMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	requestJSON, _ := json.Marshal(req)
	result, err := handler(ctx, req)
	responseJSON, _ := json.Marshal(result)
	var logEvent *zerolog.Event
	if err != nil {
		logEvent = log.Error().Str("error", fmt.Sprintf("%+v", err))
	} else {
		logEvent = log.Info()
	}
	logEvent.
		Dur("duration", time.Since(start)).
		RawJSON("json_response", responseJSON).
		RawJSON("json_request", requestJSON).
		Str("url", info.FullMethod).
		Str("ctx", fmt.Sprintf("%+v", ctx)).
		Msg("complete")

	return result, err
}
