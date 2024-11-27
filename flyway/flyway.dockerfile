FROM flyway/flyway:latest

COPY db_users /flyway/sql/db_users

RUN mkdir -p /flyway/conf

COPY flyway_users.conf /flyway/conf/flyway_users.conf

ENTRYPOINT ["sh", "-c", "flyway migrate -configFiles=/flyway/conf/flyway_users.conf"]
