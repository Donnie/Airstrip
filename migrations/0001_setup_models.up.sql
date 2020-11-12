CREATE TABLE IF NOT EXISTS accounts (
    id bigint NOT NULL,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
    currency text,
    name text,
    self boolean,
    user_id bigint,
    PRIMARY KEY (id)
);

CREATE SEQUENCE accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE ONLY accounts ALTER COLUMN id SET DEFAULT nextval('accounts_id_seq'::regclass);
ALTER SEQUENCE accounts_id_seq OWNED BY accounts.id;
SELECT pg_catalog.setval('accounts_id_seq', 1, false);

CREATE INDEX IF NOT EXISTS idx_accounts_deleted_at ON accounts(deleted_at);
CREATE INDEX IF NOT EXISTS idx_accounts_name ON accounts(name);
CREATE INDEX IF NOT EXISTS idx_accounts_self ON accounts(self);
CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);


CREATE TABLE IF NOT EXISTS convos (
    id bigint NOT NULL,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
    user_id bigint,
    expect text,
    context_id bigint,
    PRIMARY KEY (id)
);

CREATE SEQUENCE convos_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE ONLY convos ALTER COLUMN id SET DEFAULT nextval('convos_id_seq'::regclass);
ALTER SEQUENCE convos_id_seq OWNED BY convos.id;
SELECT pg_catalog.setval('convos_id_seq', 1, false);

CREATE INDEX IF NOT EXISTS idx_convos_deleted_at ON convos(deleted_at);


CREATE TABLE IF NOT EXISTS records (
    id bigint NOT NULL,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
    account_in_id bigint,
    account_out_id bigint,
    amount bigint,
    date timestamp,
    description text,
    from_date timestamp,
    mandate boolean,
    till_date timestamp,
    user_id bigint,
    PRIMARY KEY (id),
    CONSTRAINT fk_records_account_in FOREIGN KEY (account_in_id) REFERENCES accounts(id),
    CONSTRAINT fk_records_account_out FOREIGN KEY (account_out_id) REFERENCES accounts(id)
);

CREATE SEQUENCE records_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE ONLY records ALTER COLUMN id SET DEFAULT nextval('records_id_seq'::regclass);
ALTER SEQUENCE records_id_seq OWNED BY records.id;
SELECT pg_catalog.setval('records_id_seq', 1, false);

CREATE INDEX IF NOT EXISTS idx_records_account_in_id ON records(account_in_id);
CREATE INDEX IF NOT EXISTS idx_records_account_out_id ON records(account_out_id);
CREATE INDEX IF NOT EXISTS idx_records_date ON records(date);
CREATE INDEX IF NOT EXISTS idx_records_deleted_at ON records(deleted_at);
CREATE INDEX IF NOT EXISTS idx_records_mandate ON records(mandate);
CREATE INDEX IF NOT EXISTS idx_records_user_id ON records(user_id);
