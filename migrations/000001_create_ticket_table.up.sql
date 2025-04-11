CREATE TABLE IF NOT EXISTS ticket (
    id bigserial PRIMART KEY,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    cname text  NOT NULL,
    iname text NOT NULL,
    quantity int NOT NULL
);