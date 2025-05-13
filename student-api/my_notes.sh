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
# DB_HOST=localhost
# DB_USER=postgres
# DB_PASSWORD=ocYHVBzBVQYoQg50iAIE
# DB_NAME=studentdb
# DB_PORT=5432

For DB migration 
1.(https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#linux-deb-package) (install package)
2.migrate -version
3.https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md (Create migrations)
4.
migrate-create-table:
	@migrate create -ext sql -dir internal/db/migrations -seq create_users_table

# (-- after creating this field, we need to change the go model struct also)
migrate-create-table-add-phone:  
	@migrate create -ext sql -dir internal/db/migrations -seq add_phone_to_students  

After creating this file add this code in up.sql :- ALTER TABLE students ADD COLUMN phone VARCHAR(15);
in down.sql :- ALTER TABLE students DROP COLUMN phone;