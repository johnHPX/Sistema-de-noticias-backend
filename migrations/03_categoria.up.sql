create table if not exists tb_categoria(
    id varchar(36) not null,
    kind varchar(100) not null,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_categoria primary key (id)
)