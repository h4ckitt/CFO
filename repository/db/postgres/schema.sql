--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Ubuntu 14.5-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.5 (Ubuntu 14.5-0ubuntu0.22.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: cfo; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE cfo WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE = 'en_US.UTF-8';


ALTER DATABASE cfo OWNER TO postgres;

\connect cfo

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    id integer NOT NULL,
    userid integer NOT NULL,
    category character varying(128) NOT NULL
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.categories ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: spending; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.spending (
    messageid integer NOT NULL,
    userid integer NOT NULL,
    amount integer NOT NULL,
    category character varying(128) NOT NULL,
    createdat timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    note text
);


ALTER TABLE public.spending OWNER TO postgres;

--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.categories (id, userid, category) OVERRIDING SYSTEM VALUE VALUES (1, 0, 'Food 🍝');
INSERT INTO public.categories (id, userid, category) OVERRIDING SYSTEM VALUE VALUES (2, 0, 'Bills 💰');
INSERT INTO public.categories (id, userid, category) OVERRIDING SYSTEM VALUE VALUES (3, 0, 'Transportation 🚖');
INSERT INTO public.categories (id, userid, category) OVERRIDING SYSTEM VALUE VALUES (4, 0, 'Entertainment 🎮');
INSERT INTO public.categories (id, userid, category) OVERRIDING SYSTEM VALUE VALUES (5, 0, 'Shopping 🛍');
INSERT INTO public.categories (id, userid, category) OVERRIDING SYSTEM VALUE VALUES (6, 0, 'Social 👥');
INSERT INTO public.categories (id, userid, category) OVERRIDING SYSTEM VALUE VALUES (7, 0, 'Others 💸');
INSERT INTO public.categories (id, userid, category) OVERRIDING SYSTEM VALUE VALUES (8, 0, 'Savings 🤑');
INSERT INTO public.categories (id, userid, category) OVERRIDING SYSTEM VALUE VALUES (9, 0, 'Devices 📱');
INSERT INTO public.categories (id, userid, category) OVERRIDING SYSTEM VALUE VALUES (10, 0, 'Health 🩹');


--
-- Data for Name: spending; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categories_id_seq', 10, true);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: spending spending_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.spending
    ADD CONSTRAINT spending_pkey PRIMARY KEY (messageid);


--
-- PostgreSQL database dump complete
--

