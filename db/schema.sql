--
-- PostgreSQL database dump
--

-- Dumped from database version 11.5
-- Dumped by pg_dump version 11.5

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

SET default_with_oids = false;

--
-- Name: addresses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.addresses (
    id integer NOT NULL,
    address character varying(42),
    hashed_address character varying(66)
);


--
-- Name: addresses_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.addresses_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: addresses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.addresses_id_seq OWNED BY public.addresses.id;


--
-- Name: checked_headers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.checked_headers (
    id integer NOT NULL,
    header_id integer NOT NULL
);


--
-- Name: checked_headers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.checked_headers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: checked_headers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.checked_headers_id_seq OWNED BY public.checked_headers.id;


--
-- Name: eth_nodes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.eth_nodes (
    id integer NOT NULL,
    client_name character varying,
    genesis_block character varying(66),
    network_id numeric,
    eth_node_id character varying(128)
);


--
-- Name: eth_nodes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.eth_nodes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: eth_nodes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.eth_nodes_id_seq OWNED BY public.eth_nodes.id;


--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.goose_db_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.goose_db_version_id_seq OWNED BY public.goose_db_version.id;


--
-- Name: header_sync_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.header_sync_logs (
    id integer NOT NULL,
    header_id integer NOT NULL,
    address integer NOT NULL,
    topics bytea[],
    data bytea,
    block_number bigint,
    block_hash character varying(66),
    tx_hash character varying(66),
    tx_index integer,
    log_index integer,
    raw jsonb,
    transformed boolean DEFAULT false NOT NULL
);


--
-- Name: header_sync_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.header_sync_logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: header_sync_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.header_sync_logs_id_seq OWNED BY public.header_sync_logs.id;


--
-- Name: header_sync_receipts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.header_sync_receipts (
    id integer NOT NULL,
    transaction_id integer NOT NULL,
    header_id integer NOT NULL,
    contract_address_id integer NOT NULL,
    cumulative_gas_used numeric,
    gas_used numeric,
    state_root character varying(66),
    status integer,
    tx_hash character varying(66),
    rlp bytea
);


--
-- Name: header_sync_receipts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.header_sync_receipts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: header_sync_receipts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.header_sync_receipts_id_seq OWNED BY public.header_sync_receipts.id;


--
-- Name: header_sync_transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.header_sync_transactions (
    id integer NOT NULL,
    header_id integer NOT NULL,
    hash character varying(66) NOT NULL,
    gas_limit numeric,
    gas_price numeric,
    input_data bytea,
    nonce numeric,
    raw bytea,
    tx_from character varying(44),
    tx_index integer,
    tx_to character varying(44),
    value numeric
);


--
-- Name: header_sync_transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.header_sync_transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: header_sync_transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.header_sync_transactions_id_seq OWNED BY public.header_sync_transactions.id;


--
-- Name: headers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.headers (
    id integer NOT NULL,
    hash character varying(66) NOT NULL,
    block_number bigint NOT NULL,
    raw jsonb,
    block_timestamp numeric,
    check_count integer DEFAULT 0 NOT NULL,
    eth_node_id integer NOT NULL
);


--
-- Name: headers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.headers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: headers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.headers_id_seq OWNED BY public.headers.id;


--
-- Name: queued_storage; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.queued_storage (
    id integer NOT NULL,
    diff_id bigint NOT NULL
);


--
-- Name: queued_storage_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.queued_storage_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: queued_storage_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.queued_storage_id_seq OWNED BY public.queued_storage.id;


--
-- Name: storage_diff; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.storage_diff (
    id bigint NOT NULL,
    block_height bigint,
    block_hash bytea,
    hashed_address bytea,
    storage_key bytea,
    storage_value bytea
);


--
-- Name: storage_diff_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.storage_diff_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: storage_diff_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.storage_diff_id_seq OWNED BY public.storage_diff.id;


--
-- Name: watched_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.watched_logs (
    id integer NOT NULL,
    contract_address character varying(42),
    topic_zero character varying(66)
);


--
-- Name: watched_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.watched_logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: watched_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.watched_logs_id_seq OWNED BY public.watched_logs.id;


--
-- Name: addresses id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.addresses ALTER COLUMN id SET DEFAULT nextval('public.addresses_id_seq'::regclass);


--
-- Name: checked_headers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers ALTER COLUMN id SET DEFAULT nextval('public.checked_headers_id_seq'::regclass);


--
-- Name: eth_nodes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes ALTER COLUMN id SET DEFAULT nextval('public.eth_nodes_id_seq'::regclass);


--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


--
-- Name: header_sync_logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_logs ALTER COLUMN id SET DEFAULT nextval('public.header_sync_logs_id_seq'::regclass);


--
-- Name: header_sync_receipts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_receipts ALTER COLUMN id SET DEFAULT nextval('public.header_sync_receipts_id_seq'::regclass);


--
-- Name: header_sync_transactions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_transactions ALTER COLUMN id SET DEFAULT nextval('public.header_sync_transactions_id_seq'::regclass);


--
-- Name: headers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers ALTER COLUMN id SET DEFAULT nextval('public.headers_id_seq'::regclass);


--
-- Name: queued_storage id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage ALTER COLUMN id SET DEFAULT nextval('public.queued_storage_id_seq'::regclass);


--
-- Name: storage_diff id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.storage_diff ALTER COLUMN id SET DEFAULT nextval('public.storage_diff_id_seq'::regclass);


--
-- Name: watched_logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_logs ALTER COLUMN id SET DEFAULT nextval('public.watched_logs_id_seq'::regclass);


--
-- Name: addresses addresses_address_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_address_key UNIQUE (address);


--
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


--
-- Name: checked_headers checked_headers_header_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers
    ADD CONSTRAINT checked_headers_header_id_key UNIQUE (header_id);


--
-- Name: checked_headers checked_headers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers
    ADD CONSTRAINT checked_headers_pkey PRIMARY KEY (id);


--
-- Name: eth_nodes eth_nodes_genesis_block_network_id_eth_node_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes
    ADD CONSTRAINT eth_nodes_genesis_block_network_id_eth_node_id_key UNIQUE (genesis_block, network_id, eth_node_id);


--
-- Name: eth_nodes eth_nodes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes
    ADD CONSTRAINT eth_nodes_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: header_sync_logs header_sync_logs_header_id_tx_index_log_index_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_logs
    ADD CONSTRAINT header_sync_logs_header_id_tx_index_log_index_key UNIQUE (header_id, tx_index, log_index);


--
-- Name: header_sync_logs header_sync_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_logs
    ADD CONSTRAINT header_sync_logs_pkey PRIMARY KEY (id);


--
-- Name: header_sync_receipts header_sync_receipts_header_id_transaction_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_receipts
    ADD CONSTRAINT header_sync_receipts_header_id_transaction_id_key UNIQUE (header_id, transaction_id);


--
-- Name: header_sync_receipts header_sync_receipts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_receipts
    ADD CONSTRAINT header_sync_receipts_pkey PRIMARY KEY (id);


--
-- Name: header_sync_transactions header_sync_transactions_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_transactions
    ADD CONSTRAINT header_sync_transactions_hash_key UNIQUE (hash);


--
-- Name: header_sync_transactions header_sync_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_transactions
    ADD CONSTRAINT header_sync_transactions_pkey PRIMARY KEY (id);


--
-- Name: headers headers_block_number_hash_eth_node_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_block_number_hash_eth_node_id_key UNIQUE (block_number, hash, eth_node_id);


--
-- Name: headers headers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_pkey PRIMARY KEY (id);


--
-- Name: queued_storage queued_storage_diff_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage
    ADD CONSTRAINT queued_storage_diff_id_key UNIQUE (diff_id);


--
-- Name: queued_storage queued_storage_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage
    ADD CONSTRAINT queued_storage_pkey PRIMARY KEY (id);


--
-- Name: storage_diff storage_diff_block_height_block_hash_hashed_address_storage_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.storage_diff
    ADD CONSTRAINT storage_diff_block_height_block_hash_hashed_address_storage_key UNIQUE (block_height, block_hash, hashed_address, storage_key, storage_value);


--
-- Name: storage_diff storage_diff_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.storage_diff
    ADD CONSTRAINT storage_diff_pkey PRIMARY KEY (id);


--
-- Name: watched_logs watched_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_logs
    ADD CONSTRAINT watched_logs_pkey PRIMARY KEY (id);


--
-- Name: header_sync_logs_address; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX header_sync_logs_address ON public.header_sync_logs USING btree (address);


--
-- Name: header_sync_logs_transaction; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX header_sync_logs_transaction ON public.header_sync_logs USING btree (tx_hash);


--
-- Name: header_sync_logs_untransformed; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX header_sync_logs_untransformed ON public.header_sync_logs USING btree (transformed) WHERE (transformed IS FALSE);


--
-- Name: header_sync_receipts_contract_address; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX header_sync_receipts_contract_address ON public.header_sync_receipts USING btree (contract_address_id);


--
-- Name: header_sync_receipts_transaction; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX header_sync_receipts_transaction ON public.header_sync_receipts USING btree (transaction_id);


--
-- Name: header_sync_transactions_header; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX header_sync_transactions_header ON public.header_sync_transactions USING btree (header_id);


--
-- Name: headers_block_number; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX headers_block_number ON public.headers USING btree (block_number);


--
-- Name: headers_check_count; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX headers_check_count ON public.headers USING btree (check_count);


--
-- Name: headers_eth_node; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX headers_eth_node ON public.headers USING btree (eth_node_id);


--
-- Name: checked_headers checked_headers_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers
    ADD CONSTRAINT checked_headers_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: header_sync_logs header_sync_logs_address_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_logs
    ADD CONSTRAINT header_sync_logs_address_fkey FOREIGN KEY (address) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: header_sync_logs header_sync_logs_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_logs
    ADD CONSTRAINT header_sync_logs_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: header_sync_logs header_sync_logs_tx_hash_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_logs
    ADD CONSTRAINT header_sync_logs_tx_hash_fkey FOREIGN KEY (tx_hash) REFERENCES public.header_sync_transactions(hash) ON DELETE CASCADE;


--
-- Name: header_sync_receipts header_sync_receipts_contract_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_receipts
    ADD CONSTRAINT header_sync_receipts_contract_address_id_fkey FOREIGN KEY (contract_address_id) REFERENCES public.addresses(id) ON DELETE CASCADE;


--
-- Name: header_sync_receipts header_sync_receipts_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_receipts
    ADD CONSTRAINT header_sync_receipts_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: header_sync_receipts header_sync_receipts_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_receipts
    ADD CONSTRAINT header_sync_receipts_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES public.header_sync_transactions(id) ON DELETE CASCADE;


--
-- Name: header_sync_transactions header_sync_transactions_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_transactions
    ADD CONSTRAINT header_sync_transactions_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: headers headers_eth_node_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_eth_node_id_fkey FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- Name: queued_storage queued_storage_diff_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage
    ADD CONSTRAINT queued_storage_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id);


--
-- PostgreSQL database dump complete
--

