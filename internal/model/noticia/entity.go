package noticia

type Entity struct {
	NID    string
	Titulo string
}

type Conteudos struct {
	SubTitulo string `json:"subTitulo"`
	Texto     string `json:"texto"`
}

type NoticiaEntity struct {
	ID        string      `json:"id"`
	Titulo    string      `json:"titulo"`
	Conteudo  []Conteudos `json:"conteudos"`
	Categoria string      `json:"categoria"`
}
