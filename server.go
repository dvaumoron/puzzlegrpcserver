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
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type GRPCServer interface {
	grpc.ServiceRegistrar
	Start()
}

type server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func New(opts ...grpc.ServerOption) GRPCServer {
	if godotenv.Overload() == nil {
		log.Println("Loaded .env file")
	}

	lis, err := net.Listen("tcp", ":"+os.Getenv("SERVICE_PORT"))
	if err != nil {
		log.Fatal("Failed to listen :", err)
	}

	return server{grpcServer: grpc.NewServer(opts...), listener: lis}
}

func (s server) RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.grpcServer.RegisterService(desc, impl)
}

func (s server) Start() {
	log.Println("Listening at", s.listener.Addr())
	if err := s.grpcServer.Serve(s.listener); err != nil {
		log.Fatal("Failed to serve :", err)
	}
}
