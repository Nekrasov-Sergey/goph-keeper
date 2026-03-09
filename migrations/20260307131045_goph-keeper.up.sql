create table if not exists users
(
    id                 bigint generated always as identity primary key,
    login              text  not null unique,
    password_hash      text  not null,
    encrypted_user_key bytea not null
);

do
$$
    begin
        create type secret_type as enum ('LoginPassword', 'Text', 'Binary', 'BankCard');
    exception
        when duplicate_object then null;
    end
$$;

create table if not exists secrets
(
    id             bigint generated always as identity primary key,
    name           text        not null,
    type           secret_type not null,
    encrypted_data bytea       not null,
    metadata       text,
    user_id        bigint      not null references users (id),
    created_at     timestamptz not null default now(),
    updated_at     timestamptz,
    constraint secrets_name_user_id_key unique (name, user_id)
);
