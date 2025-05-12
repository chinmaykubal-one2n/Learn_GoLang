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
