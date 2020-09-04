--
-- PostgreSQL database dump
--

-- Dumped from database version 12.4
-- Dumped by pg_dump version 12.4

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
-- Name: accounts; Type: TABLE; Schema: public; Owner: airstrip
--

CREATE TABLE public.accounts (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text
);


ALTER TABLE public.accounts OWNER TO airstrip;

--
-- Name: accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: airstrip
--

CREATE SEQUENCE public.accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.accounts_id_seq OWNER TO airstrip;

--
-- Name: accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: airstrip
--

ALTER SEQUENCE public.accounts_id_seq OWNED BY public.accounts.id;


--
-- Name: convos; Type: TABLE; Schema: public; Owner: airstrip
--

CREATE TABLE public.convos (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id bigint,
    expect text,
    context_id bigint
);


ALTER TABLE public.convos OWNER TO airstrip;

--
-- Name: convos_id_seq; Type: SEQUENCE; Schema: public; Owner: airstrip
--

CREATE SEQUENCE public.convos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.convos_id_seq OWNER TO airstrip;

--
-- Name: convos_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: airstrip
--

ALTER SEQUENCE public.convos_id_seq OWNED BY public.convos.id;


--
-- Name: records; Type: TABLE; Schema: public; Owner: airstrip
--

CREATE TABLE public.records (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    account_id bigint,
    amount bigint,
    currency text,
    date timestamp with time zone,
    description text,
    form text,
    from_date timestamp with time zone,
    till_date timestamp with time zone,
    type text,
    user_id bigint
);


ALTER TABLE public.records OWNER TO airstrip;

--
-- Name: records_id_seq; Type: SEQUENCE; Schema: public; Owner: airstrip
--

CREATE SEQUENCE public.records_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.records_id_seq OWNER TO airstrip;

--
-- Name: records_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: airstrip
--

ALTER SEQUENCE public.records_id_seq OWNED BY public.records.id;


--
-- Name: accounts id; Type: DEFAULT; Schema: public; Owner: airstrip
--

ALTER TABLE ONLY public.accounts ALTER COLUMN id SET DEFAULT nextval('public.accounts_id_seq'::regclass);


--
-- Name: convos id; Type: DEFAULT; Schema: public; Owner: airstrip
--

ALTER TABLE ONLY public.convos ALTER COLUMN id SET DEFAULT nextval('public.convos_id_seq'::regclass);


--
-- Name: records id; Type: DEFAULT; Schema: public; Owner: airstrip
--

ALTER TABLE ONLY public.records ALTER COLUMN id SET DEFAULT nextval('public.records_id_seq'::regclass);


--
-- Data for Name: accounts; Type: TABLE DATA; Schema: public; Owner: airstrip
--

COPY public.accounts (id, created_at, updated_at, deleted_at, name) FROM stdin;
\.


--
-- Data for Name: convos; Type: TABLE DATA; Schema: public; Owner: airstrip
--

COPY public.convos (id, created_at, updated_at, deleted_at, user_id, expect, context_id) FROM stdin;
\.


--
-- Data for Name: records; Type: TABLE DATA; Schema: public; Owner: airstrip
--

COPY public.records (id, created_at, updated_at, deleted_at, account_id, amount, currency, date, description, form, from_date, till_date, type, user_id) FROM stdin;
\.


--
-- Name: accounts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: airstrip
--

SELECT pg_catalog.setval('public.accounts_id_seq', 1, false);


--
-- Name: convos_id_seq; Type: SEQUENCE SET; Schema: public; Owner: airstrip
--

SELECT pg_catalog.setval('public.convos_id_seq', 1, false);


--
-- Name: records_id_seq; Type: SEQUENCE SET; Schema: public; Owner: airstrip
--

SELECT pg_catalog.setval('public.records_id_seq', 1, false);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: airstrip
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: convos convos_pkey; Type: CONSTRAINT; Schema: public; Owner: airstrip
--

ALTER TABLE ONLY public.convos
    ADD CONSTRAINT convos_pkey PRIMARY KEY (id);


--
-- Name: records records_pkey; Type: CONSTRAINT; Schema: public; Owner: airstrip
--

ALTER TABLE ONLY public.records
    ADD CONSTRAINT records_pkey PRIMARY KEY (id);


--
-- Name: records fk_records_account; Type: FK CONSTRAINT; Schema: public; Owner: airstrip
--

ALTER TABLE ONLY public.records
    ADD CONSTRAINT fk_records_account FOREIGN KEY (account_id) REFERENCES public.accounts(id);


--
-- PostgreSQL database dump complete
--

