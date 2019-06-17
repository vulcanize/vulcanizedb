--
-- PostgreSQL database dump
--

-- Dumped from database version 11.3
-- Dumped by pg_dump version 11.3

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
-- Name: logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.logs (
    id integer NOT NULL,
    block_number bigint,
    address character varying(66),
    tx_hash character varying(66),
    index bigint,
    topic0 character varying(66),
    topic1 character varying(66),
    topic2 character varying(66),
    topic3 character varying(66),
    data text,
    receipt_id integer
);


--
-- Name: block_stats; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.block_stats AS
 SELECT max(logs.block_number) AS max_block,
    min(logs.block_number) AS min_block
   FROM public.logs;


--
-- Name: blocks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.blocks (
    id integer NOT NULL,
    difficulty bigint,
    extra_data character varying,
    gas_limit bigint,
    gas_used bigint,
    hash character varying(66),
    miner character varying(42),
    nonce character varying(20),
    number bigint,
    parent_hash character varying(66),
    reward numeric,
    uncles_reward numeric,
    size character varying,
    "time" bigint,
    is_final boolean,
    uncle_hash character varying(66),
    eth_node_id integer NOT NULL,
    eth_node_fingerprint character varying(128) NOT NULL
);


--
-- Name: blocks_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.blocks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: blocks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.blocks_id_seq OWNED BY public.blocks.id;


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
-- Name: full_sync_receipts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.full_sync_receipts (
    id integer NOT NULL,
    contract_address character varying(42),
    cumulative_gas_used numeric,
    gas_used numeric,
    state_root character varying(66),
    status integer,
    tx_hash character varying(66),
    block_id integer NOT NULL
);


--
-- Name: full_sync_receipts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.full_sync_receipts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: full_sync_receipts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.full_sync_receipts_id_seq OWNED BY public.full_sync_receipts.id;


--
-- Name: full_sync_transactions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.full_sync_transactions (
    id integer NOT NULL,
    block_id integer NOT NULL,
    gas_limit numeric,
    gas_price numeric,
    hash character varying(66),
    input_data bytea,
    nonce numeric,
    raw bytea,
    tx_from character varying(66),
    tx_index integer,
    tx_to character varying(66),
    value numeric
);


--
-- Name: full_sync_transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.full_sync_transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: full_sync_transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.full_sync_transactions_id_seq OWNED BY public.full_sync_transactions.id;


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
-- Name: header_sync_receipts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.header_sync_receipts (
    id integer NOT NULL,
    transaction_id integer NOT NULL,
    header_id integer NOT NULL,
    contract_address character varying(42),
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
    eth_node_id integer NOT NULL,
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
-- Name: log_filters; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.log_filters (
    id integer NOT NULL,
    name character varying NOT NULL,
    from_block bigint,
    to_block bigint,
    address character varying(66),
    topic0 character varying(66),
    topic1 character varying(66),
    topic2 character varying(66),
    topic3 character varying(66),
    CONSTRAINT log_filters_from_block_check CHECK ((from_block >= 0)),
    CONSTRAINT log_filters_name_check CHECK (((name)::text <> ''::text)),
    CONSTRAINT log_filters_to_block_check CHECK ((to_block >= 0))
);


--
-- Name: log_filters_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.log_filters_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: log_filters_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.log_filters_id_seq OWNED BY public.log_filters.id;


--
-- Name: logs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.logs_id_seq OWNED BY public.logs.id;


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

ALTER SEQUENCE public.nodes_id_seq OWNED BY public.eth_nodes.id;


--
-- Name: queued_storage; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.queued_storage (
    id integer NOT NULL,
    block_height bigint,
    block_hash bytea,
    contract bytea,
    storage_key bytea,
    storage_value bytea
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
-- Name: uncles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.uncles (
    id integer NOT NULL,
    hash character varying(66) NOT NULL,
    block_id integer NOT NULL,
    reward numeric NOT NULL,
    miner character varying(42) NOT NULL,
    raw jsonb,
    block_timestamp numeric,
    eth_node_id integer NOT NULL,
    eth_node_fingerprint character varying(128)
);


--
-- Name: uncles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.uncles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: uncles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.uncles_id_seq OWNED BY public.uncles.id;


--
-- Name: watched_contracts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.watched_contracts (
    contract_id integer NOT NULL,
    contract_abi json,
    contract_hash character varying(66)
);


--
-- Name: watched_contracts_contract_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.watched_contracts_contract_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: watched_contracts_contract_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.watched_contracts_contract_id_seq OWNED BY public.watched_contracts.contract_id;


--
-- Name: watched_event_logs; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.watched_event_logs AS
 SELECT log_filters.name,
    logs.id,
    logs.block_number,
    logs.address,
    logs.tx_hash,
    logs.index,
    logs.topic0,
    logs.topic1,
    logs.topic2,
    logs.topic3,
    logs.data,
    logs.receipt_id
   FROM ((public.log_filters
     CROSS JOIN public.block_stats)
     JOIN public.logs ON ((((logs.address)::text = (log_filters.address)::text) AND (logs.block_number >= COALESCE(log_filters.from_block, block_stats.min_block)) AND (logs.block_number <= COALESCE(log_filters.to_block, block_stats.max_block)))))
  WHERE ((((log_filters.topic0)::text = (logs.topic0)::text) OR (log_filters.topic0 IS NULL)) AND (((log_filters.topic1)::text = (logs.topic1)::text) OR (log_filters.topic1 IS NULL)) AND (((log_filters.topic2)::text = (logs.topic2)::text) OR (log_filters.topic2 IS NULL)) AND (((log_filters.topic3)::text = (logs.topic3)::text) OR (log_filters.topic3 IS NULL)));


--
-- Name: blocks id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blocks ALTER COLUMN id SET DEFAULT nextval('public.blocks_id_seq'::regclass);


--
-- Name: checked_headers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers ALTER COLUMN id SET DEFAULT nextval('public.checked_headers_id_seq'::regclass);


--
-- Name: eth_nodes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes ALTER COLUMN id SET DEFAULT nextval('public.nodes_id_seq'::regclass);


--
-- Name: full_sync_receipts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_receipts ALTER COLUMN id SET DEFAULT nextval('public.full_sync_receipts_id_seq'::regclass);


--
-- Name: full_sync_transactions id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_transactions ALTER COLUMN id SET DEFAULT nextval('public.full_sync_transactions_id_seq'::regclass);


--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


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
-- Name: log_filters id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.log_filters ALTER COLUMN id SET DEFAULT nextval('public.log_filters_id_seq'::regclass);


--
-- Name: logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logs ALTER COLUMN id SET DEFAULT nextval('public.logs_id_seq'::regclass);


--
-- Name: queued_storage id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage ALTER COLUMN id SET DEFAULT nextval('public.queued_storage_id_seq'::regclass);


--
-- Name: uncles id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles ALTER COLUMN id SET DEFAULT nextval('public.uncles_id_seq'::regclass);


--
-- Name: watched_contracts contract_id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_contracts ALTER COLUMN contract_id SET DEFAULT nextval('public.watched_contracts_contract_id_seq'::regclass);


--
-- Name: blocks blocks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT blocks_pkey PRIMARY KEY (id);


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
-- Name: blocks eth_node_id_block_number_uc; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT eth_node_id_block_number_uc UNIQUE (number, eth_node_id);


--
-- Name: eth_nodes eth_node_uc; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes
    ADD CONSTRAINT eth_node_uc UNIQUE (genesis_block, network_id, eth_node_id);


--
-- Name: full_sync_receipts full_sync_receipts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_receipts
    ADD CONSTRAINT full_sync_receipts_pkey PRIMARY KEY (id);


--
-- Name: full_sync_transactions full_sync_transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_transactions
    ADD CONSTRAINT full_sync_transactions_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


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
-- Name: headers headers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_pkey PRIMARY KEY (id);


--
-- Name: logs logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT logs_pkey PRIMARY KEY (id);


--
-- Name: log_filters name_uc; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.log_filters
    ADD CONSTRAINT name_uc UNIQUE (name);


--
-- Name: eth_nodes nodes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.eth_nodes
    ADD CONSTRAINT nodes_pkey PRIMARY KEY (id);


--
-- Name: queued_storage queued_storage_block_height_block_hash_contract_storage_key_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage
    ADD CONSTRAINT queued_storage_block_height_block_hash_contract_storage_key_key UNIQUE (block_height, block_hash, contract, storage_key, storage_value);


--
-- Name: queued_storage queued_storage_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.queued_storage
    ADD CONSTRAINT queued_storage_pkey PRIMARY KEY (id);


--
-- Name: uncles uncles_block_id_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles
    ADD CONSTRAINT uncles_block_id_hash_key UNIQUE (block_id, hash);


--
-- Name: uncles uncles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles
    ADD CONSTRAINT uncles_pkey PRIMARY KEY (id);


--
-- Name: watched_contracts watched_contracts_contract_hash_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_contracts
    ADD CONSTRAINT watched_contracts_contract_hash_key UNIQUE (contract_hash);


--
-- Name: watched_contracts watched_contracts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.watched_contracts
    ADD CONSTRAINT watched_contracts_pkey PRIMARY KEY (contract_id);


--
-- Name: block_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX block_id_index ON public.full_sync_transactions USING btree (block_id);


--
-- Name: headers_block_number; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX headers_block_number ON public.headers USING btree (block_number);


--
-- Name: node_id_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX node_id_index ON public.blocks USING btree (eth_node_id);


--
-- Name: number_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX number_index ON public.blocks USING btree (number);


--
-- Name: tx_from_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tx_from_index ON public.full_sync_transactions USING btree (tx_from);


--
-- Name: tx_to_index; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX tx_to_index ON public.full_sync_transactions USING btree (tx_to);


--
-- Name: full_sync_receipts blocks_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_receipts
    ADD CONSTRAINT blocks_fk FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: checked_headers checked_headers_header_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.checked_headers
    ADD CONSTRAINT checked_headers_header_id_fkey FOREIGN KEY (header_id) REFERENCES public.headers(id) ON DELETE CASCADE;


--
-- Name: full_sync_transactions full_sync_transactions_block_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.full_sync_transactions
    ADD CONSTRAINT full_sync_transactions_block_id_fkey FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


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
-- Name: blocks node_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT node_fk FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- Name: logs receipts_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT receipts_fk FOREIGN KEY (receipt_id) REFERENCES public.full_sync_receipts(id) ON DELETE CASCADE;


--
-- Name: uncles uncles_block_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles
    ADD CONSTRAINT uncles_block_id_fkey FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: uncles uncles_eth_node_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.uncles
    ADD CONSTRAINT uncles_eth_node_id_fkey FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

