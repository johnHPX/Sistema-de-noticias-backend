package conteudo

type Conteudo struct {
	Subtitulo string
	Texto     string
}

type Entity struct {
	CID     string
	Contudo Conteudo
	NID     string
}
