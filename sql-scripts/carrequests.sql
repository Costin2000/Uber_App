CREATE SEQUENCE public.car_request_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE public.car_request_id_seq OWNER TO postgres;

SET default_tablespace = '';
SET default_table_access_method = heap;

CREATE TABLE public.car_requests (
                                     id integer DEFAULT nextval('public.car_request_id_seq'::regclass) NOT NULL,
                                     user_id integer,
                                     user_name character varying(255),
                                     car_type character varying(255),
                                     car_id integer,
                                     city character varying(255),
                                     address character varying(255),
                                     active boolean DEFAULT true,
                                     rating integer DEFAULT 0,
                                     created_at timestamp without time zone,
                                     updated_at timestamp without time zone
);

ALTER TABLE public.car_requests OWNER TO postgres;

SELECT pg_catalog.setval('public.car_request_id_seq', 1, true);

ALTER TABLE ONLY public.car_requests
    ADD CONSTRAINT car_requests_pkey PRIMARY KEY (id);
