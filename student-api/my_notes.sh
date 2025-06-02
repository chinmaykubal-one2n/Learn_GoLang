go mod init student-api
go get -u github.com/gin-gonic/gin
go get -u github.com/google/uuid


go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/joho/godotenv


# install PostgreSQL
https://www.postgresql.org/download/linux/ubuntu/

# install pgadmin
https://www.pgadmin.org/download/pgadmin-4-apt/

# how to setup the PostgreSQL
# (then create the db server in the pgadmin and then right click on Databases and create db studentdb)
sudo -i -u postgres
psql
ALTER USER postgres WITH PASSWORD 'your_new_password';
\q
exit

# .env content
# DB_HOST=
# DB_USER=
# DB_PASSWORD=
# DB_NAME=
# DB_PORT=
# DB_URL=
# JWT_SECRET=

# SERVICE_NAME=student-api
# INSECURE_MODE=true
# OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317

#OTLP_ENDPOINT=localhost:4318

For DB migration 
1.(https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#linux-deb-package) (install package)
2.migrate -version
3.https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md (Create migrations)
4.
migrate-create-table:
	@migrate create -ext sql -dir internal/db/migrations -seq create_users_table # then write schema in it.


migrate-up:
	@migrate -path internal/db/migrations -database $(DB_URL) up

migrate-down:
	@migrate -path internal/db/migrations -database $(DB_URL) down 1

migrate-status:
	@migrate -path internal/db/migrations -database $(DB_URL) version

migrate-force:
	@migrate -path internal/db/migrations -database $(DB_URL) force VERSION


	

# (-- after creating this field, we need to change the go model struct also)
migrate-create-table-add-phone:  
	@migrate create -ext sql -dir internal/db/migrations -seq add_phone_to_students  

After creating this file add this code in up.sql :- ALTER TABLE students ADD COLUMN phone VARCHAR(15);
in down.sql :- ALTER TABLE students DROP COLUMN phone;


# DB migration Atlas with GORM
# (https://atlasgo.io/)
1. curl -sSf https://atlasgo.sh | sh 
2. go get -u ariga.io/atlas-provider-gorm
3. Create an Atlas Configuration File atlas.hcl
# COMMANDS SECTION
4. atlas migrate diff --env gorm # add new field in struct and then run this command  
5. atlas migrate apply --env gorm
6. atlas migrate down --env gorm  
7. atlas migrate status --env gorm  
8. atlas migrate rm 20250513115014 --env gorm


# gin-jwt middleware (https://github.com/appleboy/gin-jwt)
1. 
export GO111MODULE=on
go get github.com/appleboy/gin-jwt/v2


# swggar docs
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
swag init

# opentelemetry setup (tracing+metrics) with signoze :- https://signoz.io/blog/opentelemetry-gin/
# logs https://signoz.io/docs/logs-management/send-logs/zap-to-signoz/#requirements-1


# Docker stuff
docker build -t student-api-go:multistage .


docker run -d -p 8090:8090 \
  --env-file .env.docker \
  --add-host=host.docker.internal:host-gateway \
  --network signoz-net \
  student-api-go:multistage