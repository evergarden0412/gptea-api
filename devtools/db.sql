create table if not exists users(
    id text primary key,
    created_at timestamptz not null default now() 
);

create table if not exists user_credentials(
    user_id text references users(id) on delete cascade not null,
    credential_type text not null,
    credential_id text not null,
    created_at timestamptz not null default now(),
    primary key (credential_type, credential_id)
);

create table if not exists refresh_tokens(
    user_id text references users(id) on delete cascade primary key, 
    token_id text not null,
    created_at timestamptz not null default now() 
);

create table if not exists chats(
    id text primary key,
    user_id text references users(id) on delete cascade not null,
    name text not null,
    created_at timestamptz not null default now() 
);

create table if not exists messages(
    chat_id text references chats(id) on delete cascade not null,
    seq integer not null,
    content text not null,
    role text not null,
    created_at timestamptz not null default now(), 
    unique (chat_id, seq)
);

create table if not exists scrapbooks(
    id text primary key,
    user_id text references users(id) on delete cascade not null,
    name text not null,
    is_default boolean not null default false,
    created_at timestamptz not null default now() 
);

create table if not exists scraps(
    id text primary key,
    message_chat_id text not null,
    message_seq integer not null,
    memo text not null,
    created_at timestamptz not null default now(),
    foreign key (message_chat_id, message_seq) references messages(chat_id, seq) on delete cascade,
    unique(message_chat_id, message_seq)
);

create table if not exists scraps_scrapbooks(
    scrap_id text references scraps(id) on delete cascade not null,
    scrapbook_id text references scrapbooks(id) on delete cascade not null,
    created_at timestamptz not null default now(),
    primary key (scrap_id, scrapbook_id)
);