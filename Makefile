.PHONY: compile
compile: ## Compile the proto file.
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative client/common/proto/auth.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative server/common/proto/auth.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative auth_server/common/proto/auth.proto

.PHONY: dep-install
dep-install: ## Build and run server.
	cd auth_server && go mod tidy && cd ..
	cd server && go mod tidy && cd ..
	cd client && go mod tidy && cd ..


.PHONY: serve
server: ## Build and run server.
	docker-compose up -d

.PHONY: watch
watch: ## watch server.
	docker-compose logs -f --tail 50 app-dev
