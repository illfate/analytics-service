CREATE TABLE events
(
    client_time DateTime,
    server_time DateTime,
    ip          String,
    device_id   String,
    device_os   String,
    session     String,
    sequence    Int64,
    event       String,
    param_int   Int64,
    param_str   String
) ENGINE = MergeTree()
    ORDER BY client_time;