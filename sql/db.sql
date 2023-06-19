-- auto-generated definition
create table snippets
(
    id      serial       not null
        constraint snippets_pkey
            primary key,
    title   varchar(100) not null,
    content text         not null,
    created timestamp    not null,
    expires timestamp    not null
);

alter table snippets
    owner to postgres;

create index idx_snippets_created
    on snippets (created);
