--
-- PostgreSQL database dump
--

-- Dumped from database version 10.10
-- Dumped by pg_dump version 12.1

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
-- Name: btc; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA btc;


--
-- Name: eth; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA eth;


SET default_tablespace = '';

--
-- Name: header_cids; Type: TABLE; Schema: btc; Owner: -
--

CREATE TABLE btc.header_cids (
    id integer NOT NULL,
    block_number bigint NOT NULL,
    block_hash character varying(66) NOT NULL,
    parent_hash character varying(66) NOT NULL,
    cid text NOT NULL,
    "timestamp" numeric NOT NULL,
    bits bigint NOT NULL,
    node_id integer NOT NULL
);


--
-- Name: header_cids_id_seq; Type: SEQUENCE; Schema: btc; Owner: -
--

CREATE SEQUENCE btc.header_cids_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: header_cids_id_seq; Type: SEQUENCE OWNED BY; Schema: btc; Owner: -
--

ALTER SEQUENCE btc.header_cids_id_seq OWNED BY btc.header_cids.id;


--
-- Name: transaction_cids; Type: TABLE; Schema: btc; Owner: -
--

CREATE TABLE btc.transaction_cids (
    id integer NOT NULL,
    header_id integer NOT NULL,
    index integer NOT NULL,
    tx_hash character varying(66) NOT NULL,
    cid text NOT NULL,
    segwit boolean NOT NULL,
    witness_hash character varying(66)
);


--
-- Name: transaction_cids_id_seq; Type: SEQUENCE; Schema: btc; Owner: -
--

CREATE SEQUENCE btc.transaction_cids_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: transaction_cids_id_seq; Type: SEQUENCE OWNED BY; Schema: btc; Owner: -
--

ALTER SEQUENCE btc.transaction_cids_id_seq OWNED BY btc.transaction_cids.id;


--
-- Name: tx_inputs; Type: TABLE; Schema: btc; Owner: -
--

CREATE TABLE btc.tx_inputs (
    id integer NOT NULL,
    tx_id integer NOT NULL,
    index integer NOT NULL,
    witness bytea[],
    sig_script bytea NOT NULL,
    outpoint_id integer
);


--
-- Name: tx_inputs_id_seq; Type: SEQUENCE; Schema: btc; Owner: -
--

CREATE SEQUENCE btc.tx_inputs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tx_inputs_id_seq; Type: SEQUENCE OWNED BY; Schema: btc; Owner: -
--

ALTER SEQUENCE btc.tx_inputs_id_seq OWNED BY btc.tx_inputs.id;


--
-- Name: tx_outputs; Type: TABLE; Schema: btc; Owner: -
--

CREATE TABLE btc.tx_outputs (
    id integer NOT NULL,
    tx_id integer NOT NULL,
    index integer NOT NULL,
    value bigint NOT NULL,
    pk_script bytea NOT NULL,
    script_class integer NOT NULL,
    addresses character varying(66)[],
    required_sigs integer NOT NULL
);


--
-- Name: tx_outputs_id_seq; Type: SEQUENCE; Schema: btc; Owner: -
--

CREATE SEQUENCE btc.tx_outputs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tx_outputs_id_seq; Type: SEQUENCE OWNED BY; Schema: btc; Owner: -
--

ALTER SEQUENCE btc.tx_outputs_id_seq OWNED BY btc.tx_outputs.id;


--
-- Name: header_cids; Type: TABLE; Schema: eth; Owner: -
--

CREATE TABLE eth.header_cids (
    id integer NOT NULL,
    block_number bigint NOT NULL,
    block_hash character varying(66) NOT NULL,
    parent_hash character varying(66) NOT NULL,
    cid text NOT NULL,
    td numeric NOT NULL,
    node_id integer NOT NULL
);


--
-- Name: header_cids_id_seq; Type: SEQUENCE; Schema: eth; Owner: -
--

CREATE SEQUENCE eth.header_cids_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: header_cids_id_seq; Type: SEQUENCE OWNED BY; Schema: eth; Owner: -
--

ALTER SEQUENCE eth.header_cids_id_seq OWNED BY eth.header_cids.id;


--
-- Name: receipt_cids; Type: TABLE; Schema: eth; Owner: -
--

CREATE TABLE eth.receipt_cids (
    id integer NOT NULL,
    tx_id integer NOT NULL,
    cid text NOT NULL,
    contract character varying(66),
    topic0s character varying(66)[],
    topic1s character varying(66)[],
    topic2s character varying(66)[],
    topic3s character varying(66)[]
);


--
-- Name: receipt_cids_id_seq; Type: SEQUENCE; Schema: eth; Owner: -
--

CREATE SEQUENCE eth.receipt_cids_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: receipt_cids_id_seq; Type: SEQUENCE OWNED BY; Schema: eth; Owner: -
--

ALTER SEQUENCE eth.receipt_cids_id_seq OWNED BY eth.receipt_cids.id;


--
-- Name: state_cids; Type: TABLE; Schema: eth; Owner: -
--

CREATE TABLE eth.state_cids (
    id integer NOT NULL,
    header_id integer NOT NULL,
    state_key character varying(66) NOT NULL,
    leaf boolean NOT NULL,
    cid text NOT NULL
);


--
-- Name: state_cids_id_seq; Type: SEQUENCE; Schema: eth; Owner: -
--

CREATE SEQUENCE eth.state_cids_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: state_cids_id_seq; Type: SEQUENCE OWNED BY; Schema: eth; Owner: -
--

ALTER SEQUENCE eth.state_cids_id_seq OWNED BY eth.state_cids.id;


--
-- Name: storage_cids; Type: TABLE; Schema: eth; Owner: -
--

CREATE TABLE eth.storage_cids (
    id integer NOT NULL,
    state_id integer NOT NULL,
    storage_key character varying(66) NOT NULL,
    leaf boolean NOT NULL,
    cid text NOT NULL
);


--
-- Name: storage_cids_id_seq; Type: SEQUENCE; Schema: eth; Owner: -
--

CREATE SEQUENCE eth.storage_cids_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: storage_cids_id_seq; Type: SEQUENCE OWNED BY; Schema: eth; Owner: -
--

ALTER SEQUENCE eth.storage_cids_id_seq OWNED BY eth.storage_cids.id;


--
-- Name: transaction_cids; Type: TABLE; Schema: eth; Owner: -
--

CREATE TABLE eth.transaction_cids (
    id integer NOT NULL,
    header_id integer NOT NULL,
    tx_hash character varying(66) NOT NULL,
    index integer NOT NULL,
    cid text NOT NULL,
    dst character varying(66) NOT NULL,
    src character varying(66) NOT NULL
);


--
-- Name: transaction_cids_id_seq; Type: SEQUENCE; Schema: eth; Owner: -
--

CREATE SEQUENCE eth.transaction_cids_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: transaction_cids_id_seq; Type: SEQUENCE OWNED BY; Schema: eth; Owner: -
--

ALTER SEQUENCE eth.transaction_cids_id_seq OWNED BY eth.transaction_cids.id;


--
-- Name: uncle_cids; Type: TABLE; Schema: eth; Owner: -
--

CREATE TABLE eth.uncle_cids (
    id integer NOT NULL,
    header_id integer NOT NULL,
    block_hash character varying(66) NOT NULL,
    parent_hash character varying(66) NOT NULL,
    cid text NOT NULL
);


--
-- Name: uncle_cids_id_seq; Type: SEQUENCE; Schema: eth; Owner: -
--

CREATE SEQUENCE eth.uncle_cids_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: uncle_cids_id_seq; Type: SEQUENCE OWNED BY; Schema: eth; Owner: -
--

ALTER SEQUENCE eth.uncle_cids_id_seq OWNED BY eth.uncle_cids.id;


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
-- Name: blocks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.blocks (
    key text NOT NULL,
    data bytea NOT NULL
);


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
    hash character varying(66),
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
    hash character varying(66),
    block_number bigint,
    raw jsonb,
    block_timestamp numeric,
    check_count integer DEFAULT 0 NOT NULL,
    node_id integer NOT NULL,
    eth_node_fingerprint character varying(128)
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
-- Name: nodes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.nodes (
    id integer NOT NULL,
    client_name character varying,
    genesis_block character varying(66),
    network_id character varying,
    node_id character varying(128)
);


--
-- Name: nodes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.nodes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: nodes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.nodes_id_seq OWNED BY public.nodes.id;


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
    id integer NOT NULL,
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
    AS integer
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
-- Name: header_cids id; Type: DEFAULT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.header_cids ALTER COLUMN id SET DEFAULT nextval('btc.header_cids_id_seq'::regclass);


--
-- Name: transaction_cids id; Type: DEFAULT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.transaction_cids ALTER COLUMN id SET DEFAULT nextval('btc.transaction_cids_id_seq'::regclass);


--
-- Name: tx_inputs id; Type: DEFAULT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.tx_inputs ALTER COLUMN id SET DEFAULT nextval('btc.tx_inputs_id_seq'::regclass);


--
-- Name: tx_outputs id; Type: DEFAULT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.tx_outputs ALTER COLUMN id SET DEFAULT nextval('btc.tx_outputs_id_seq'::regclass);


--
-- Name: header_cids id; Type: DEFAULT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.header_cids ALTER COLUMN id SET DEFAULT nextval('eth.header_cids_id_seq'::regclass);


--
-- Name: receipt_cids id; Type: DEFAULT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.receipt_cids ALTER COLUMN id SET DEFAULT nextval('eth.receipt_cids_id_seq'::regclass);


--
-- Name: state_cids id; Type: DEFAULT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.state_cids ALTER COLUMN id SET DEFAULT nextval('eth.state_cids_id_seq'::regclass);


--
-- Name: storage_cids id; Type: DEFAULT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.storage_cids ALTER COLUMN id SET DEFAULT nextval('eth.storage_cids_id_seq'::regclass);


--
-- Name: transaction_cids id; Type: DEFAULT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.transaction_cids ALTER COLUMN id SET DEFAULT nextval('eth.transaction_cids_id_seq'::regclass);


--
-- Name: uncle_cids id; Type: DEFAULT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.uncle_cids ALTER COLUMN id SET DEFAULT nextval('eth.uncle_cids_id_seq'::regclass);


--
-- Name: addresses id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.addresses ALTER COLUMN id SET DEFAULT nextval('public.addresses_id_seq'::regclass);


--
-- Name: checked_headers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers ALTER COLUMN id SET DEFAULT nextval('public.checked_headers_id_seq'::regclass);


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
-- Name: nodes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.nodes ALTER COLUMN id SET DEFAULT nextval('public.nodes_id_seq'::regclass);


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
-- Name: header_cids header_cids_block_number_block_hash_key; Type: CONSTRAINT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.header_cids
    ADD CONSTRAINT header_cids_block_number_block_hash_key UNIQUE (block_number, block_hash);


--
-- Name: header_cids header_cids_pkey; Type: CONSTRAINT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.header_cids
    ADD CONSTRAINT header_cids_pkey PRIMARY KEY (id);


--
-- Name: transaction_cids transaction_cids_pkey; Type: CONSTRAINT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.transaction_cids
    ADD CONSTRAINT transaction_cids_pkey PRIMARY KEY (id);


--
-- Name: transaction_cids transaction_cids_tx_hash_key; Type: CONSTRAINT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.transaction_cids
    ADD CONSTRAINT transaction_cids_tx_hash_key UNIQUE (tx_hash);


--
-- Name: tx_inputs tx_inputs_pkey; Type: CONSTRAINT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.tx_inputs
    ADD CONSTRAINT tx_inputs_pkey PRIMARY KEY (id);


--
-- Name: tx_inputs tx_inputs_tx_id_index_key; Type: CONSTRAINT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.tx_inputs
    ADD CONSTRAINT tx_inputs_tx_id_index_key UNIQUE (tx_id, index);


--
-- Name: tx_outputs tx_outputs_pkey; Type: CONSTRAINT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.tx_outputs
    ADD CONSTRAINT tx_outputs_pkey PRIMARY KEY (id);


--
-- Name: tx_outputs tx_outputs_tx_id_index_key; Type: CONSTRAINT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.tx_outputs
    ADD CONSTRAINT tx_outputs_tx_id_index_key UNIQUE (tx_id, index);


--
-- Name: header_cids header_cids_block_number_block_hash_key; Type: CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.header_cids
    ADD CONSTRAINT header_cids_block_number_block_hash_key UNIQUE (block_number, block_hash);


--
-- Name: header_cids header_cids_pkey; Type: CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.header_cids
    ADD CONSTRAINT header_cids_pkey PRIMARY KEY (id);


--
-- Name: receipt_cids receipt_cids_pkey; Type: CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.receipt_cids
    ADD CONSTRAINT receipt_cids_pkey PRIMARY KEY (id);


--
-- Name: state_cids state_cids_header_id_state_key_key; Type: CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.state_cids
    ADD CONSTRAINT state_cids_header_id_state_key_key UNIQUE (header_id, state_key);


--
-- Name: state_cids state_cids_pkey; Type: CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.state_cids
    ADD CONSTRAINT state_cids_pkey PRIMARY KEY (id);


--
-- Name: storage_cids storage_cids_pkey; Type: CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.storage_cids
    ADD CONSTRAINT storage_cids_pkey PRIMARY KEY (id);


--
-- Name: storage_cids storage_cids_state_id_storage_key_key; Type: CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.storage_cids
    ADD CONSTRAINT storage_cids_state_id_storage_key_key UNIQUE (state_id, storage_key);


--
-- Name: transaction_cids transaction_cids_header_id_tx_hash_key; Type: CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.transaction_cids
    ADD CONSTRAINT transaction_cids_header_id_tx_hash_key UNIQUE (header_id, tx_hash);


--
-- Name: transaction_cids transaction_cids_pkey; Type: CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.transaction_cids
    ADD CONSTRAINT transaction_cids_pkey PRIMARY KEY (id);


--
-- Name: uncle_cids uncle_cids_header_id_block_hash_key; Type: CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.uncle_cids
    ADD CONSTRAINT uncle_cids_header_id_block_hash_key UNIQUE (header_id, block_hash);


--
-- Name: uncle_cids uncle_cids_pkey; Type: CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.uncle_cids
    ADD CONSTRAINT uncle_cids_pkey PRIMARY KEY (id);


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
-- Name: blocks blocks_key_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT blocks_key_key UNIQUE (key);


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
-- Name: header_sync_transactions header_sync_transactions_header_id_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_transactions
    ADD CONSTRAINT header_sync_transactions_header_id_hash_key UNIQUE (header_id, hash);


--
-- Name: header_sync_transactions header_sync_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.header_sync_transactions
    ADD CONSTRAINT header_sync_transactions_pkey PRIMARY KEY (id);


--
-- Name: headers headers_block_number_hash_eth_node_fingerprint_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_block_number_hash_eth_node_fingerprint_key UNIQUE (block_number, hash, eth_node_fingerprint);


--
-- Name: headers headers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_pkey PRIMARY KEY (id);


--
-- Name: nodes node_uc; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.nodes
    ADD CONSTRAINT node_uc UNIQUE (genesis_block, network_id, node_id);


--
-- Name: nodes nodes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.nodes
    ADD CONSTRAINT nodes_pkey PRIMARY KEY (id);


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
-- Name: header_sync_receipts_header; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX header_sync_receipts_header ON public.header_sync_receipts USING btree (header_id);


--
-- Name: header_sync_receipts_transaction; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX header_sync_receipts_transaction ON public.header_sync_receipts USING btree (transaction_id);


--
-- Name: header_sync_transactions_header; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX header_sync_transactions_header ON public.header_sync_transactions USING btree (header_id);


--
-- Name: header_sync_transactions_tx_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX header_sync_transactions_tx_index ON public.header_sync_transactions USING btree (tx_index);


--
-- Name: headers_block_number; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX headers_block_number ON public.headers USING btree (block_number);


--
-- Name: headers_block_timestamp; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX headers_block_timestamp ON public.headers USING btree (block_timestamp);


--
-- Name: header_cids header_cids_node_id_fkey; Type: FK CONSTRAINT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.header_cids
    ADD CONSTRAINT header_cids_node_id_fkey FOREIGN KEY (node_id) REFERENCES public.nodes(id) ON DELETE CASCADE;


--
-- Name: transaction_cids transaction_cids_header_id_fkey; Type: FK CONSTRAINT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.transaction_cids
    ADD CONSTRAINT transaction_cids_header_id_fkey FOREIGN KEY (header_id) REFERENCES btc.header_cids(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: tx_inputs tx_inputs_outpoint_id_fkey; Type: FK CONSTRAINT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.tx_inputs
    ADD CONSTRAINT tx_inputs_outpoint_id_fkey FOREIGN KEY (outpoint_id) REFERENCES btc.tx_outputs(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: tx_inputs tx_inputs_tx_id_fkey; Type: FK CONSTRAINT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.tx_inputs
    ADD CONSTRAINT tx_inputs_tx_id_fkey FOREIGN KEY (tx_id) REFERENCES btc.transaction_cids(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: tx_outputs tx_outputs_tx_id_fkey; Type: FK CONSTRAINT; Schema: btc; Owner: -
--

ALTER TABLE ONLY btc.tx_outputs
    ADD CONSTRAINT tx_outputs_tx_id_fkey FOREIGN KEY (tx_id) REFERENCES btc.transaction_cids(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: header_cids header_cids_node_id_fkey; Type: FK CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.header_cids
    ADD CONSTRAINT header_cids_node_id_fkey FOREIGN KEY (node_id) REFERENCES public.nodes(id) ON DELETE CASCADE;


--
-- Name: receipt_cids receipt_cids_tx_id_fkey; Type: FK CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.receipt_cids
    ADD CONSTRAINT receipt_cids_tx_id_fkey FOREIGN KEY (tx_id) REFERENCES eth.transaction_cids(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: state_cids state_cids_header_id_fkey; Type: FK CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.state_cids
    ADD CONSTRAINT state_cids_header_id_fkey FOREIGN KEY (header_id) REFERENCES eth.header_cids(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: storage_cids storage_cids_state_id_fkey; Type: FK CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.storage_cids
    ADD CONSTRAINT storage_cids_state_id_fkey FOREIGN KEY (state_id) REFERENCES eth.state_cids(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: transaction_cids transaction_cids_header_id_fkey; Type: FK CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.transaction_cids
    ADD CONSTRAINT transaction_cids_header_id_fkey FOREIGN KEY (header_id) REFERENCES eth.header_cids(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


--
-- Name: uncle_cids uncle_cids_header_id_fkey; Type: FK CONSTRAINT; Schema: eth; Owner: -
--

ALTER TABLE ONLY eth.uncle_cids
    ADD CONSTRAINT uncle_cids_header_id_fkey FOREIGN KEY (header_id) REFERENCES eth.header_cids(id) ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;


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
-- Name: headers headers_node_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_node_id_fkey FOREIGN KEY (node_id) REFERENCES public.nodes(id) ON DELETE CASCADE;


--
-- Name: queued_storage queued_storage_diff_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage
    ADD CONSTRAINT queued_storage_diff_id_fkey FOREIGN KEY (diff_id) REFERENCES public.storage_diff(id);


--
-- PostgreSQL database dump complete
--

