CREATE SCHEMA IF NOT EXISTS messenger AUTHORIZATION postgres;

CREATE SEQUENCE IF NOT EXISTS messenger.messages_id_seq INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;

ALTER SEQUENCE messenger.messages_id_seq OWNER TO postgres;

CREATE TABLE IF NOT EXISTS messenger.messages (
    id bigint NOT NULL DEFAULT nextval('messenger.messages_id_seq' :: regclass),
    body text COLLATE pg_catalog."default" NOT NULL,
    "timestamp" timestamp with time zone NOT NULL,
    CONSTRAINT messages_pkey PRIMARY KEY (id)
) TABLESPACE pg_default;

ALTER TABLE
    IF EXISTS messenger.messages OWNER to postgres;