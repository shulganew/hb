CREATE USER hb WITH ENCRYPTED PASSWORD '1';
CREATE DATABASE hb;

GRANT ALL PRIVILEGES ON DATABASE hb TO hb;

ALTER DATABASE hb OWNER TO hb;