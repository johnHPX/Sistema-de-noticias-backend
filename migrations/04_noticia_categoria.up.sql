create table tb_noticia_categoria(
    nid varchar(36) not null,
    cid varchar(36) not null,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_noticia_categoria primary key (nid, cid)
)