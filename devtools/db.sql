create table if not exists users(
    id text primary key,
    created_at timestamptz default now()
);

create table if not exists user_credentials(
    user_id text references users(id) not null,
    credential_type text not null,
    credential_id text not null,
    created_at timestamptz default now(),
    primary key (credential_type, credential_id)
);

create table if not exists refresh_tokens(
    user_id text references users(id) not null,
    token_id text not null,
    created_at timestamptz default now(),
    primary key (user_id, token_id)
);

create table if not exists chats(
    id text primary key,
    user_id text references users(id) not null,
    name text,
    created_at timestamptz default now()
);

create table if not exists messages(
    id text primary key,
    chat_id text references chats(id) not null,
    seq integer not null,
    content text,
    role text,
    created_at timestamptz default now(),
    unique (chat_id, seq)
);

create table if not exists scrapbooks(
    id text primary key,
    user_id text references users(id) not null,
    name text,
    created_at timestamptz default now()
);

create table if not exists scraps(
    scrapbook_id text references scrapbooks(id),
    message_id text references messages(id),
    created_at timestamptz default now(),
    primary key (scrapbook_id, message_id)
);
