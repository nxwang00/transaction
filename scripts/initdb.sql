CREATE TABLE trans (
  id serial NOT NULL,
  origin character varying(30) COLLATE pg_catalog."default" NOT NULL,
  user_id integer NOT NULL,
  amount money NOT NULL,
  op_type character varying(30) COLLATE pg_catalog."default" NOT NULL,
  registered_at timestamp without time zone NOT NULL,
  CONSTRAINT trans_pkey PRIMARY KEY (id)
)