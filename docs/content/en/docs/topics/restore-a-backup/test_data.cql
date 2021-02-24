CREATE KEYSPACE medusa_test WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
USE medusa_test;
CREATE TABLE users (email text primary key, name text, state text);
insert into users (email, name, state) values ('alice@example.com', 'Alice Smith', 'TX');
insert into users (email, name, state) values ('bob@example.com', 'Bob Jones', 'VA');
insert into users (email, name, state) values ('carol@example.com', 'Carol Jackson', 'CA');
insert into users (email, name, state) values ('david@example.com', 'David Yang', 'NV');
