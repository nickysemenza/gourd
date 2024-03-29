package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

var (
	latexError = fmt.Errorf("latex error")
)

func (a *API) GetLatexByRecipeId(c echo.Context, recipeId string) error {
	ctx := c.Request().Context()
	res, err := a.Latex(ctx, recipeId)
	if err != nil {
		if errors.Is(err, latexError) {
			return c.Blob(http.StatusInternalServerError, "text/plain", []byte(err.Error()))
		}
		return handleErr(c, err)
	}
	return c.Blob(http.StatusOK, "application/pdf", res)

}
func (a *API) Latex(ctx context.Context, id string) ([]byte, error) {
	apiR, err := a.recipeById(ctx, id)
	if err != nil {
		return nil, err
	}

	// from https://github.com/anaseto/gofrundis/blob/master/escape/escape.go
	var latexEscapes = []string{
		// "{", "\\{",
		// "}", "\\}",
		// "[", "[",
		// "]", "]",
		"%", `\%`,
		"&", `\&`,
		"$", `\$`,
		"#", `\#`,
		"_", `\_`,
		"'", `\textquotesingle{}`,
		`"`, `\textquotesingle{}`,
		"^", `\^{}`,
		`\`, `\textbackslash{}`,
		// todo: get vulgar frac mapping from rs?
		"¾", `3/4`,
		"¼", `1/4`,
		"⅛", `1/8`,
		// "~", "\\~{}",
		// string('\xa0'), "~"
	}

	var latexEscaper = strings.NewReplacer(latexEscapes...)

	// LaTeX escapes special LaTeX characters.
	esc := func(text string) string {
		return latexEscaper.Replace(text)
	}

	u := func(amounts []Amount) string {

		if s := firstAmount(amounts, false); s != nil {
			return s.Unit
		}
		return ""

	}
	aa := func(amounts []Amount) string {
		if s := firstAmount(amounts, false); s != nil {
			return fmt.Sprintf("%.2f", s.Value)
		}
		return ""

	}
	g := func(amounts []Amount) string {
		if s := firstAmount(amounts, true); s != nil {
			return fmt.Sprintf("%.2f", s.Value)
		}
		return ""
	}
	n := func(si SectionIngredient) string {
		if si.Ingredient != nil {
			return si.Ingredient.Ingredient.Name
		}
		if si.Recipe != nil {
			return fmt.Sprintf(`\bf{%s}`, si.Recipe.Name)
		}
		return ""
	}
	adj := func(si SectionIngredient) string {
		if si.Adjective != nil {
			return esc(*si.Adjective)
		}
		return ""
	}

	const templateText = `
\documentclass{article}
\usepackage{multirow}
\usepackage{booktabs}
\usepackage{graphicx}
\usepackage[margin=0.5in]{geometry}
\title{ {{.Detail.Name}}}
\author{todo}
\begin{document}

\maketitle

\section{Recipe}
\begin{table}[htbp]
\centering
\resizebox{\textwidth}{!}{%
\begin{tabular}{|l|rllll|l|}
\hline
\multirow{2}{*}{section} & \multicolumn{5}{c|}{ingredient}  & \multirow{2}{*}{instruction} \\ \cline{2-6}
							& \multicolumn{1}{l|}{amount} & \multicolumn{1}{l|}{unit}  & \multicolumn{1}{l|}{grams}  & \multicolumn{1}{l|}{name}  & adj             &                              \\ \hline
{{range $i, $s := .Detail.Sections -}}
	{{range $j, $ing := $s.Ingredients -}}
		{{if eq $j 0 -}}
			\multirow{ {{$s.Ingredients | len}} }{*}{A}
			{{end -}}
			& \multicolumn{1}{r|}{ {{$ing.Amounts| a }} }     & \multicolumn{1}{l|}{ {{$ing.Amounts| u }} } & \multicolumn{1}{l|}{ {{$ing.Amounts| g }} }   & \multicolumn{1}{l|}{ {{$ing | n}} }   &      \multicolumn{1}{l|}{ {{$ing | adj}} }           
			{{if eq $j 0 -}}
				& \multirow{ {{$s.Ingredients | len}} }{*}{\parbox[t]{10cm}{ {{ $s.Instructions | foo }} }}           \\ 
			{{else -}}
				& \\ 
			{{end -}}
			{{ $length := len $s.Ingredients }} {{if (isLast $j $length)}}
			\hline
		{{else -}}
			\cline{2-6}
		{{end -}}
	{{end -}}
{{end -}}
\end{tabular}%
}
\end{table}
\end{document}
`

	foo := func(si []SectionInstruction) string {
		i := []string{}
		for _, x := range si {
			i = append(i, esc(x.Instruction))
		}
		return strings.Join(i, `\\ `)

	}

	funcMap := template.FuncMap{

		"foo": foo,
		"u":   u,
		"adj": adj,
		"g":   g,
		"a":   aa,
		"n":   n,
		"esc": esc,
		"isLast": func(index int, len int) bool {
			return index+1 == len
		},
	}

	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("titleTest").Funcs(funcMap).Parse(templateText)
	if err != nil {
		return nil, err
	}

	file, err := os.CreateTemp("", id+".*.tex")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	log.Println(file.Name())

	// Run the template to verify the output.
	err = tmpl.Execute(file, apiR)
	if err != nil {
		return nil, err
	}

	dir, err := os.MkdirTemp("", "text-")
	if err != nil {
		return nil, err
	}
	log.Println(dir)

	viper.SetDefault("PDFLATEX_BINARY", "pdflatex")
	viper.AutomaticEnv()
	binary := viper.GetString("PDFLATEX_BINARY")

	ctx, span := a.tracer.Start(ctx, "latexCommand")
	defer span.End()

	jobName := "gourd"

	cmd := exec.CommandContext(ctx, binary, "-jobname="+jobName, "-output-directory", dir, file.Name())
	cmd.Dir = dir
	// cmd.Stdin = strings.NewReader(document)

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("error running %s: %w", binary, err)
	}
	err = cmd.Wait()
	if err != nil {
		// The actual error is useless, do provide a better one.
		output, err := os.ReadFile(path.Join(dir, jobName+".log"))
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("failed to generate latex: %s: %w", string(output), latexError)
	}
	outFile := path.Join(dir, jobName+".pdf")
	log.Println(outFile)

	output, err := os.ReadFile(outFile)
	if err != nil {
		return nil, err
	}

	// Clean up the temp directory.
	_ = os.RemoveAll(dir)

	return output, nil
}
