CREATE SEQUENCE public.car_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE public.car_id_seq OWNER TO postgres;

SET default_tablespace = '';
SET default_table_access_method = heap;

CREATE TABLE public.cars (
                             id integer DEFAULT nextval('public.car_id_seq'::regclass) NOT NULL,
                             user_id integer,
                             car_name character varying(255),
                             city character varying(255),
                             car_type character varying(255),
                             created_at timestamp without time zone,
                             updated_at timestamp without time zone
);

ALTER TABLE public.cars OWNER TO postgres;

SELECT pg_catalog.setval('public.car_id_seq', 1, true);

ALTER TABLE ONLY public.cars
    ADD CONSTRAINT cars_pkey PRIMARY KEY (id);

ALTER TABLE public.cars
    ADD COLUMN active boolean DEFAULT false;