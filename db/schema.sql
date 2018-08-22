--
-- PostgreSQL database dump
--

-- Dumped from database version 10.5
-- Dumped by pg_dump version 10.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: logs; Type: TABLE; Schema: public; Owner: iannorden
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


ALTER TABLE public.logs OWNER TO iannorden;

--
-- Name: block_stats; Type: VIEW; Schema: public; Owner: iannorden
--

CREATE VIEW public.block_stats AS
 SELECT max(logs.block_number) AS max_block,
    min(logs.block_number) AS min_block
   FROM public.logs;


ALTER TABLE public.block_stats OWNER TO iannorden;

--
-- Name: blocks; Type: TABLE; Schema: public; Owner: iannorden
--

CREATE TABLE public.blocks (
    number bigint,
    gaslimit bigint,
    gasused bigint,
    "time" bigint,
    id integer NOT NULL,
    difficulty bigint,
    hash character varying(66),
    nonce character varying(20),
    parenthash character varying(66),
    size character varying,
    uncle_hash character varying(66),
    eth_node_id integer NOT NULL,
    is_final boolean,
    miner character varying(42),
    extra_data character varying,
    reward double precision,
    uncles_reward double precision,
    eth_node_fingerprint character varying(128) NOT NULL
);


ALTER TABLE public.blocks OWNER TO iannorden;

--
-- Name: blocks_id_seq; Type: SEQUENCE; Schema: public; Owner: iannorden
--

CREATE SEQUENCE public.blocks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.blocks_id_seq OWNER TO iannorden;

--
-- Name: blocks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: iannorden
--

ALTER SEQUENCE public.blocks_id_seq OWNED BY public.blocks.id;


--
-- Name: eth_nodes; Type: TABLE; Schema: public; Owner: iannorden
--

CREATE TABLE public.eth_nodes (
    id integer NOT NULL,
    genesis_block character varying(66),
    network_id numeric,
    eth_node_id character varying(128),
    client_name character varying
);


ALTER TABLE public.eth_nodes OWNER TO iannorden;

--
-- Name: headers; Type: TABLE; Schema: public; Owner: iannorden
--

CREATE TABLE public.headers (
    id integer NOT NULL,
    hash character varying(66),
    block_number bigint,
    raw bytea,
    eth_node_id integer,
    eth_node_fingerprint character varying(128)
);


ALTER TABLE public.headers OWNER TO iannorden;

--
-- Name: headers_id_seq; Type: SEQUENCE; Schema: public; Owner: iannorden
--

CREATE SEQUENCE public.headers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.headers_id_seq OWNER TO iannorden;

--
-- Name: headers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: iannorden
--

ALTER SEQUENCE public.headers_id_seq OWNED BY public.headers.id;


--
-- Name: log_filters; Type: TABLE; Schema: public; Owner: iannorden
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


ALTER TABLE public.log_filters OWNER TO iannorden;

--
-- Name: log_filters_id_seq; Type: SEQUENCE; Schema: public; Owner: iannorden
--

CREATE SEQUENCE public.log_filters_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.log_filters_id_seq OWNER TO iannorden;

--
-- Name: log_filters_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: iannorden
--

ALTER SEQUENCE public.log_filters_id_seq OWNED BY public.log_filters.id;


--
-- Name: logs_id_seq; Type: SEQUENCE; Schema: public; Owner: iannorden
--

CREATE SEQUENCE public.logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.logs_id_seq OWNER TO iannorden;

--
-- Name: logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: iannorden
--

ALTER SEQUENCE public.logs_id_seq OWNED BY public.logs.id;


--
-- Name: nodes_id_seq; Type: SEQUENCE; Schema: public; Owner: iannorden
--

CREATE SEQUENCE public.nodes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.nodes_id_seq OWNER TO iannorden;

--
-- Name: nodes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: iannorden
--

ALTER SEQUENCE public.nodes_id_seq OWNED BY public.eth_nodes.id;


--
-- Name: receipts; Type: TABLE; Schema: public; Owner: iannorden
--

CREATE TABLE public.receipts (
    id integer NOT NULL,
    contract_address character varying(42),
    cumulative_gas_used numeric,
    gas_used numeric,
    state_root character varying(66),
    status integer,
    tx_hash character varying(66),
    block_id integer NOT NULL
);


ALTER TABLE public.receipts OWNER TO iannorden;

--
-- Name: receipts_id_seq; Type: SEQUENCE; Schema: public; Owner: iannorden
--

CREATE SEQUENCE public.receipts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.receipts_id_seq OWNER TO iannorden;

--
-- Name: receipts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: iannorden
--

ALTER SEQUENCE public.receipts_id_seq OWNED BY public.receipts.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: iannorden
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO iannorden;

--
-- Name: token_supply; Type: TABLE; Schema: public; Owner: iannorden
--

CREATE TABLE public.token_supply (
    id integer NOT NULL,
    block_id integer NOT NULL,
    supply numeric NOT NULL,
    token_address character varying(66) NOT NULL
);


ALTER TABLE public.token_supply OWNER TO iannorden;

--
-- Name: token_supply_id_seq; Type: SEQUENCE; Schema: public; Owner: iannorden
--

CREATE SEQUENCE public.token_supply_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.token_supply_id_seq OWNER TO iannorden;

--
-- Name: token_supply_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: iannorden
--

ALTER SEQUENCE public.token_supply_id_seq OWNED BY public.token_supply.id;


--
-- Name: token_allowance; Type: TABLE; Schema: public; Owner: iannorden
--

CREATE TABLE public.token_allowance (
    id integer DEFAULT nextval('public.token_supply_id_seq'::regclass) NOT NULL,
    block_id integer NOT NULL,
    allowance numeric NOT NULL,
    token_address character varying(66) NOT NULL,
    token_holder_address character varying(66) NOT NULL,
    token_spender_address character varying(66) NOT NULL
);


ALTER TABLE public.token_allowance OWNER TO iannorden;

--
-- Name: token_allowance_id_seq; Type: SEQUENCE; Schema: public; Owner: iannorden
--

CREATE SEQUENCE public.token_allowance_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.token_allowance_id_seq OWNER TO iannorden;

--
-- Name: token_allowance_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: iannorden
--

ALTER SEQUENCE public.token_allowance_id_seq OWNED BY public.token_allowance.id;


--
-- Name: token_balance; Type: TABLE; Schema: public; Owner: iannorden
--

CREATE TABLE public.token_balance (
    id integer DEFAULT nextval('public.token_supply_id_seq'::regclass) NOT NULL,
    block_id integer NOT NULL,
    balance numeric NOT NULL,
    token_address character varying(66) NOT NULL,
    token_holder_address character varying(66) NOT NULL
);


ALTER TABLE public.token_balance OWNER TO iannorden;

--
-- Name: token_balance_id_seq; Type: SEQUENCE; Schema: public; Owner: iannorden
--

CREATE SEQUENCE public.token_balance_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.token_balance_id_seq OWNER TO iannorden;

--
-- Name: token_balance_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: iannorden
--

ALTER SEQUENCE public.token_balance_id_seq OWNED BY public.token_balance.id;


--
-- Name: transactions; Type: TABLE; Schema: public; Owner: iannorden
--

CREATE TABLE public.transactions (
    id integer NOT NULL,
    hash character varying(66),
    nonce numeric,
    tx_to character varying(66),
    gaslimit numeric,
    gasprice numeric,
    value numeric,
    block_id integer NOT NULL,
    tx_from character varying(66),
    input_data character varying
);


ALTER TABLE public.transactions OWNER TO iannorden;

--
-- Name: transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: iannorden
--

CREATE SEQUENCE public.transactions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.transactions_id_seq OWNER TO iannorden;

--
-- Name: transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: iannorden
--

ALTER SEQUENCE public.transactions_id_seq OWNED BY public.transactions.id;


--
-- Name: watched_contracts; Type: TABLE; Schema: public; Owner: iannorden
--

CREATE TABLE public.watched_contracts (
    contract_id integer NOT NULL,
    contract_hash character varying(66),
    contract_abi json
);


ALTER TABLE public.watched_contracts OWNER TO iannorden;

--
-- Name: watched_contracts_contract_id_seq; Type: SEQUENCE; Schema: public; Owner: iannorden
--

CREATE SEQUENCE public.watched_contracts_contract_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.watched_contracts_contract_id_seq OWNER TO iannorden;

--
-- Name: watched_contracts_contract_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: iannorden
--

ALTER SEQUENCE public.watched_contracts_contract_id_seq OWNED BY public.watched_contracts.contract_id;


--
-- Name: watched_event_logs; Type: VIEW; Schema: public; Owner: iannorden
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


ALTER TABLE public.watched_event_logs OWNER TO iannorden;

--
-- Name: blocks id; Type: DEFAULT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.blocks ALTER COLUMN id SET DEFAULT nextval('public.blocks_id_seq'::regclass);


--
-- Name: eth_nodes id; Type: DEFAULT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.eth_nodes ALTER COLUMN id SET DEFAULT nextval('public.nodes_id_seq'::regclass);


--
-- Name: headers id; Type: DEFAULT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.headers ALTER COLUMN id SET DEFAULT nextval('public.headers_id_seq'::regclass);


--
-- Name: log_filters id; Type: DEFAULT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.log_filters ALTER COLUMN id SET DEFAULT nextval('public.log_filters_id_seq'::regclass);


--
-- Name: logs id; Type: DEFAULT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.logs ALTER COLUMN id SET DEFAULT nextval('public.logs_id_seq'::regclass);


--
-- Name: receipts id; Type: DEFAULT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.receipts ALTER COLUMN id SET DEFAULT nextval('public.receipts_id_seq'::regclass);


--
-- Name: token_supply id; Type: DEFAULT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.token_supply ALTER COLUMN id SET DEFAULT nextval('public.token_supply_id_seq'::regclass);


--
-- Name: transactions id; Type: DEFAULT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.transactions ALTER COLUMN id SET DEFAULT nextval('public.transactions_id_seq'::regclass);


--
-- Name: watched_contracts contract_id; Type: DEFAULT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.watched_contracts ALTER COLUMN contract_id SET DEFAULT nextval('public.watched_contracts_contract_id_seq'::regclass);


--
-- Data for Name: blocks; Type: TABLE DATA; Schema: public; Owner: iannorden
--

COPY public.blocks (number, gaslimit, gasused, "time", id, difficulty, hash, nonce, parenthash, size, uncle_hash, eth_node_id, is_final, miner, extra_data, reward, uncles_reward, eth_node_fingerprint) FROM stdin;
\.


--
-- Data for Name: eth_nodes; Type: TABLE DATA; Schema: public; Owner: iannorden
--

COPY public.eth_nodes (id, genesis_block, network_id, eth_node_id, client_name) FROM stdin;
1	GENESIS	1	2ea672a45c4c7b96e3c4b130b21a22af390a552fd0b3cff96420b4bda26568d470dc56e05e453823f64f2556a6e4460ad1d4d00eb2d8b8fc16fcb1be73e86522	Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9
111	0x456	1		
66	GENESIS	1	x123	geth
42	GENESIS	1	b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845	Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9
70		0	EthNodeFingerprint	
104	0x456	1	x123456	Geth
73		0	Fingerprint	
74		0	FingerprintTwo	
5	GENESIS	1	testNodeId	Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9
69		0		
81		0	NodeFingerprint	
79		0	NodeFingerprintTwo	
\.


--
-- Data for Name: headers; Type: TABLE DATA; Schema: public; Owner: iannorden
--

COPY public.headers (id, hash, block_number, raw, eth_node_id, eth_node_fingerprint) FROM stdin;
304		1	\\x	81	NodeFingerprint
305		3	\\x	81	NodeFingerprint
306		5	\\x	81	NodeFingerprint
\.


--
-- Data for Name: log_filters; Type: TABLE DATA; Schema: public; Owner: iannorden
--

COPY public.log_filters (id, name, from_block, to_block, address, topic0, topic1, topic2, topic3) FROM stdin;
\.


--
-- Data for Name: logs; Type: TABLE DATA; Schema: public; Owner: iannorden
--

COPY public.logs (id, block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data, receipt_id) FROM stdin;
\.


--
-- Data for Name: receipts; Type: TABLE DATA; Schema: public; Owner: iannorden
--

COPY public.receipts (id, contract_address, cumulative_gas_used, gas_used, state_root, status, tx_hash, block_id) FROM stdin;
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: iannorden
--

COPY public.schema_migrations (version, dirty) FROM stdin;
\.


--
-- Data for Name: token_allowance; Type: TABLE DATA; Schema: public; Owner: iannorden
--

COPY public.token_allowance (id, block_id, allowance, token_address, token_holder_address, token_spender_address) FROM stdin;
\.


--
-- Data for Name: token_balance; Type: TABLE DATA; Schema: public; Owner: iannorden
--

COPY public.token_balance (id, block_id, balance, token_address, token_holder_address) FROM stdin;
\.


--
-- Data for Name: token_supply; Type: TABLE DATA; Schema: public; Owner: iannorden
--

COPY public.token_supply (id, block_id, supply, token_address) FROM stdin;
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: iannorden
--

COPY public.transactions (id, hash, nonce, tx_to, gaslimit, gasprice, value, block_id, tx_from, input_data) FROM stdin;
\.


--
-- Data for Name: watched_contracts; Type: TABLE DATA; Schema: public; Owner: iannorden
--

COPY public.watched_contracts (contract_id, contract_hash, contract_abi) FROM stdin;
\.


--
-- Name: blocks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: iannorden
--

SELECT pg_catalog.setval('public.blocks_id_seq', 1902, true);


--
-- Name: headers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: iannorden
--

SELECT pg_catalog.setval('public.headers_id_seq', 306, true);


--
-- Name: log_filters_id_seq; Type: SEQUENCE SET; Schema: public; Owner: iannorden
--

SELECT pg_catalog.setval('public.log_filters_id_seq', 102, true);


--
-- Name: logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: iannorden
--

SELECT pg_catalog.setval('public.logs_id_seq', 290, true);


--
-- Name: nodes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: iannorden
--

SELECT pg_catalog.setval('public.nodes_id_seq', 1770, true);


--
-- Name: receipts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: iannorden
--

SELECT pg_catalog.setval('public.receipts_id_seq', 153, true);


--
-- Name: token_allowance_id_seq; Type: SEQUENCE SET; Schema: public; Owner: iannorden
--

SELECT pg_catalog.setval('public.token_allowance_id_seq', 1, false);


--
-- Name: token_balance_id_seq; Type: SEQUENCE SET; Schema: public; Owner: iannorden
--

SELECT pg_catalog.setval('public.token_balance_id_seq', 1, false);


--
-- Name: token_supply_id_seq; Type: SEQUENCE SET; Schema: public; Owner: iannorden
--

SELECT pg_catalog.setval('public.token_supply_id_seq', 400, true);


--
-- Name: transactions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: iannorden
--

SELECT pg_catalog.setval('public.transactions_id_seq', 340, true);


--
-- Name: watched_contracts_contract_id_seq; Type: SEQUENCE SET; Schema: public; Owner: iannorden
--

SELECT pg_catalog.setval('public.watched_contracts_contract_id_seq', 102, true);


--
-- Name: blocks blocks_pkey; Type: CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT blocks_pkey PRIMARY KEY (id);


--
-- Name: watched_contracts contract_hash_uc; Type: CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.watched_contracts
    ADD CONSTRAINT contract_hash_uc UNIQUE (contract_hash);


--
-- Name: blocks eth_node_id_block_number_uc; Type: CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT eth_node_id_block_number_uc UNIQUE (number, eth_node_id);


--
-- Name: eth_nodes eth_node_uc; Type: CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.eth_nodes
    ADD CONSTRAINT eth_node_uc UNIQUE (genesis_block, network_id, eth_node_id);


--
-- Name: headers headers_pkey; Type: CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT headers_pkey PRIMARY KEY (id);


--
-- Name: logs logs_pkey; Type: CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT logs_pkey PRIMARY KEY (id);


--
-- Name: log_filters name_uc; Type: CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.log_filters
    ADD CONSTRAINT name_uc UNIQUE (name);


--
-- Name: eth_nodes nodes_pkey; Type: CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.eth_nodes
    ADD CONSTRAINT nodes_pkey PRIMARY KEY (id);


--
-- Name: receipts receipts_pkey; Type: CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT receipts_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: watched_contracts watched_contracts_pkey; Type: CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.watched_contracts
    ADD CONSTRAINT watched_contracts_pkey PRIMARY KEY (contract_id);


--
-- Name: block_id_index; Type: INDEX; Schema: public; Owner: iannorden
--

CREATE INDEX block_id_index ON public.transactions USING btree (block_id);


--
-- Name: block_number_index; Type: INDEX; Schema: public; Owner: iannorden
--

CREATE INDEX block_number_index ON public.blocks USING btree (number);


--
-- Name: node_id_index; Type: INDEX; Schema: public; Owner: iannorden
--

CREATE INDEX node_id_index ON public.blocks USING btree (eth_node_id);


--
-- Name: tx_from_index; Type: INDEX; Schema: public; Owner: iannorden
--

CREATE INDEX tx_from_index ON public.transactions USING btree (tx_from);


--
-- Name: tx_to_index; Type: INDEX; Schema: public; Owner: iannorden
--

CREATE INDEX tx_to_index ON public.transactions USING btree (tx_to);


--
-- Name: transactions blocks_fk; Type: FK CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT blocks_fk FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: receipts blocks_fk; Type: FK CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.receipts
    ADD CONSTRAINT blocks_fk FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: token_supply blocks_fk; Type: FK CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.token_supply
    ADD CONSTRAINT blocks_fk FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: token_balance blocks_fk; Type: FK CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.token_balance
    ADD CONSTRAINT blocks_fk FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: token_allowance blocks_fk; Type: FK CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.token_allowance
    ADD CONSTRAINT blocks_fk FOREIGN KEY (block_id) REFERENCES public.blocks(id) ON DELETE CASCADE;


--
-- Name: headers eth_nodes_fk; Type: FK CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.headers
    ADD CONSTRAINT eth_nodes_fk FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- Name: blocks node_fk; Type: FK CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.blocks
    ADD CONSTRAINT node_fk FOREIGN KEY (eth_node_id) REFERENCES public.eth_nodes(id) ON DELETE CASCADE;


--
-- Name: logs receipts_fk; Type: FK CONSTRAINT; Schema: public; Owner: iannorden
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT receipts_fk FOREIGN KEY (receipt_id) REFERENCES public.receipts(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

