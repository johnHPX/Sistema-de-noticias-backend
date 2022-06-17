create table if not exists tb_conteudo(
    id varchar(36) not null,
    subtitulo varchar(255) not null,
    texto text not null,
    noticia_nid varchar(36) not null,
    created_at timestamp not null default now(),
    updated_at timestamp,
    deleted_at timestamp,
    constraint pk_conteudo primary key (id),
    constraint fk_pk_conteudo_0 foreign key (noticia_nid) references tb_noticia(id)
)