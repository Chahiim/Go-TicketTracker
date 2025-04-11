CREATE TABLE IF NOT EXISTS item (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    itemName text NOT NULL,
    desc text NOT NULL
);