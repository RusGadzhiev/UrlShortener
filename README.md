## UrlShortener

<br/>

## Технологии

- PostgreSQL в качестве хранилища ссылок
- Развертывание приложения независимо от платформы с помощью Docker
- Разработка API в соответствии с принципами REST
- На транспортном уровне - gRPC и HTTP
- Многоуровневое ведение журнала с использованием пакета Zap
- Трехуровневая архитектура: Transport –> Service –> Storage

Технологии: Go, gRPC, SQL, Git, Docker, Linux, PostgreSQL, SOLID, REST, Zap

## Рабочее дерево
```
UrlShortener
├── cmd
│   └── url_shortener
│       └── main.go
├── config.yaml
├── docker-compose.yaml
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── service
│   │   ├── encoder
│   │   │   └── encoder.go
│   │   └── service.go
│   ├── storage
│   │   └── postgres
│   │       └── postgres.go
│   └── transport
│       ├── gRPC
│       │   ├── gRPCHandler
│       │   │   ├── grpcMiddleware.go
│       │   │   └── gRPCHandler.go
│       │   └── gRPCServer
│       │       └── gRPCServer.go
│       └── http
│           ├── httpHandler
│           │   ├── httpMiddleware.go
│           │   └── htppHandler.go
│           └── httpServer
│               └── httpServer.go
├── Makefile
├── pkg
│   ├── logger
│   │   └── logger.go
│   └── validator
│       └── validator.go
├── proto
│   ├── url_shortener_grpc.pb.go
│   ├── url_shortener.pb.go
│   └── url_shortener.proto
├── README.md
├── TASK.md
└── Dockerfile

```

##  Начало работы

0. Установите Go, Docker и прочее

1. Клонируйте репозиторий

```bash
git clone https://github.com/RusGadzhiev/UrlShortener
```   

2. Запустите контейнеры, используя следующие команды:
```
# http server
 make http
```

```
# grpc server
 make grpc
```

## Contact  

Гаджиев Руслан - [@RusGadzhiev](https://t.me/driveinyourheart)

