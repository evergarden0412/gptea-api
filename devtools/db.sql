create table if not exists users(
    id text primary key,
    created_at timestamptz default now()
);

create table if not exists user_credentials(
    user_id text references users(id) on delete cascade not null,
    credential_type text not null,
    credential_id text not null,
    created_at timestamptz default now() not null,
    primary key (credential_type, credential_id)
);

create table if not exists refresh_tokens(
    user_id text references users(id) on delete cascade primary key, 
    token_id text not null,
    created_at timestamptz default now() not null
);

create table if not exists chats(
    id text primary key,
    user_id text references users(id) on delete cascade not null,
    name text not null,
    created_at timestamptz default now() not null
);

create table if not exists messages(
    id text primary key,
    chat_id text references chats(id) on delete cascade not null,
    seq integer not null,
    content text not null,
    role text not null,
    created_at timestamptz default now() not null,
    unique (chat_id, seq)
);

create table if not exists scrapbooks(
    id text primary key,
    user_id text references users(id) on delete cascade not null,
    name text not null,
    created_at timestamptz default now() not null
);

create table if not exists scraps(
    scrapbook_id text references scrapbooks(id) on delete cascade not null,
    message_id text references messages(id) on delete cascade not null,
    created_at timestamptz default now() not null,
    primary key (scrapbook_id, message_id)
);
