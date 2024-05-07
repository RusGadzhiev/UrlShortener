.PHONY: grpc
grpc:
	pg_pass="pass" transport_mode="grpc" docker compose up --build

.PHONY: http
http:
	pg_pass="pass" transport_mode="http" docker compose up --build

.DEFAULT_GOAL := http