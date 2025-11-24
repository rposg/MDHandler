package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/yuin/goldmark"
)

const outputFilename = "output/output.html"

/* convertMarkdownToHTML converte um markdown local em HTML usando a lib goldmark */
func convertMarkdownToHTML(markdown []byte) (string, error) {
	md := goldmark.New()
	var buf bytes.Buffer

	if err := md.Convert(markdown, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

/* openBrowser tenta abrir o arquivo no navegador padrão, mostrando o conteúdo do html recém convertido*/
func openBrowser(filepath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", filepath)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", filepath)
	case "darwin": // macOS
		cmd = exec.Command("open", filepath)
	default:
		return fmt.Errorf("sistema operacional não suportado..: %s", runtime.GOOS)
	}

	log.Printf("Abrindo o arquivo %s no navegador...", filepath)
	return cmd.Start()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("exemplo de uso: go run main.go <input/arquivo.md>")
	}
	filename := os.Args[1]
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("erro ao ler arquivo: %v", err)
	}

	html, err := convertMarkdownToHTML(content)
	if err != nil {
		log.Fatalf("erro ao converter markdown: %v", err)
	}

	htmlBase := `<!DOCTYPE html>
<html lang="pt">
	<head>
    	<meta charset="UTF-8">
		<!-- dispositivos pequenos -->
    	<meta name="viewport" content="width=device-width, initial-scale=1.0"> 
    	<title>Markdown Convertido</title>
	</head>
	<body>
		%s
	</body>
</html>`

	htmlFinal := fmt.Sprintf(htmlBase, html)
	//fmt.Println(html) ao invés de imprimir, salvar em arquivo
	//err = ioutil.WriteFile(outputFilename, []byte(html), 0644)
	err = ioutil.WriteFile(outputFilename, []byte(htmlFinal), 0644)
	if err != nil {
		log.Fatalf("erro ao escrever arquivo HTML: %v", err)
	}
	log.Printf("Arquivo Html salvo em %s", outputFilename)

	// Tenta abrir o arquivo HTML no navegador do usuário

	if err := openBrowser(outputFilename); err != nil {
		log.Printf(" Warning: Não foi possível abrir o navegador e mostrar o conteúdo do html gerado: %v", err)
	} else {
		log.Println("Navegador aberto com sucesso.")
	}
}
