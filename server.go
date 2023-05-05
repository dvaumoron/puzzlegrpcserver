/*
 *
 * Copyright 2023 puzzlegrpcserver authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package puzzlegrpcserver

import (
	"context"
	"net"
	"os"

	"github.com/dvaumoron/puzzletelemetry"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	pb "google.golang.org/grpc/health/grpc_health_v1"
)

const grpcKey = "PuzzleGRPCServer"

var _ grpc.ServiceRegistrar = GRPCServer{}

type GRPCServer struct {
	inner    *grpc.Server
	listener net.Listener
	Logger   *otelzap.Logger
	tp       *trace.TracerProvider
	Tracer   oteltrace.Tracer
}

func Make(serviceName string, version string, opts ...grpc.ServerOption) GRPCServer {
	logger, tp := puzzletelemetry.Init(serviceName, version)

	lis, err := net.Listen("tcp", ":"+os.Getenv("SERVICE_PORT"))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	augmentedOpts := make([]grpc.ServerOption, 0, len(opts)+2)
	augmentedOpts = append(augmentedOpts, grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	augmentedOpts = append(augmentedOpts, grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()))
	augmentedOpts = append(augmentedOpts, opts...)

	grpcServer := grpc.NewServer(augmentedOpts...)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", pb.HealthCheckResponse_SERVING)
	pb.RegisterHealthServer(grpcServer, healthServer)

	tracer := tp.Tracer(grpcKey)
	return GRPCServer{inner: grpcServer, listener: lis, Logger: logger, tp: tp, Tracer: tracer}
}

func (s GRPCServer) RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.inner.RegisterService(desc, impl)
}

func (s GRPCServer) Start() {
	tp := s.tp
	ctx := context.Background()
	_, startSpan := s.Tracer.Start(ctx, "start")
	s.Logger.InfoContext(ctx, "Listening", zap.String("address", s.listener.Addr().String()))
	startSpan.End()

	err := s.inner.Serve(s.listener)
	if err2 := tp.Shutdown(context.Background()); err2 != nil {
		_, stopSpan := s.Tracer.Start(ctx, "shutdown")
		s.Logger.WarnContext(ctx, "Failed to shutdown telemetry", zap.Error(err2))
		stopSpan.End()
	}
	if err != nil {
		s.Logger.Fatal("Failed to serve", zap.Error(err))
	}
}
