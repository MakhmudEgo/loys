CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE GENDER AS ENUM ('male', 'female');


CREATE TABLE IF NOT EXISTS users (
    id          UUID         DEFAULT uuid_generate_v4() PRIMARY KEY,
    password    VARCHAR(100) NOT NULL,
    first_name  VARCHAR(50)  NOT NULL,
    second_name VARCHAR(50)  NOT NULL,
    birthdate   DATE         NOT NULL,
    gender      GENDER  NOT NULL,
    biography   VARCHAR(100) NOT NULL,
    city        VARCHAR(50)  NOT NULL
);
