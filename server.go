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
	"net"
	"os"

	"github.com/dvaumoron/puzzlelogger"
	"github.com/dvaumoron/puzzlelogger/grpclogger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health"
	pb "google.golang.org/grpc/health/grpc_health_v1"
)

type GRPCServer interface {
	grpc.ServiceRegistrar
	Start()
}

type server struct {
	grpcServer *grpc.Server
	listener   net.Listener
	Logger     *zap.Logger
}

func New(opts ...grpc.ServerOption) GRPCServer {
	logger := puzzlelogger.New()

	lis, err := net.Listen("tcp", ":"+os.Getenv("SERVICE_PORT"))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	grpclog.SetLoggerV2(grpclogger.New(logger))
	grpcServer := grpc.NewServer(opts...)

	healthServer := health.NewServer()
	healthServer.SetServingStatus("", pb.HealthCheckResponse_SERVING)
	pb.RegisterHealthServer(grpcServer, healthServer)

	return server{grpcServer: grpcServer, listener: lis, Logger: logger}
}

func (s server) RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.grpcServer.RegisterService(desc, impl)
}

func (s server) Start() {
	s.Logger.Info("Listening", zap.String("address", s.listener.Addr().String()))
	if err := s.grpcServer.Serve(s.listener); err != nil {
		s.Logger.Fatal("Failed to serve", zap.Error(err))
	}
}
