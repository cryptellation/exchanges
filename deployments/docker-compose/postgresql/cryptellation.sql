CREATE USER cryptellation;
ALTER USER cryptellation PASSWORD 'cryptellation';
ALTER USER cryptellation CREATEDB;

CREATE DATABASE exchanges;
GRANT ALL PRIVILEGES ON DATABASE exchanges TO cryptellation;
\c exchanges postgres
GRANT ALL ON SCHEMA public TO cryptellation;