package main

import (
	"bytes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
)

// TemplateData contains all the necessary data
// to fill a profile card HTML template
type TemplateData struct {
	User     User
	BasePath string
	Preview  bool
}

func GeneratePDF(data *TemplateData, templatePath string) ([]byte, error) {
	var tmpl *template.Template
	var err error

	// use Go's default HTML template generation tools to generate your HTML
	if tmpl, err = template.ParseFiles(templatePath); err != nil {
		return nil, err
	}

	// apply the parsed HTML template data and keep the result in a Buffer
	var body bytes.Buffer
	if err = tmpl.Execute(&body, data); err != nil {
		return nil, err
	}

	// initialize a wkhtmltopdf generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	// read the HTML page as a PDF page
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(body.Bytes()))

	// enable this if the HTML file contains local references such as images, CSS, etc.
	page.EnableLocalFileAccess.Set(true)
	page.DisableInternalLinks.Set(true)
	// add the page to your generator
	pdfg.AddPage(page)

	// manipulate page attributes as needed
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginTop.Set(0)
	pdfg.Dpi.Set(600)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA5)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationLandscape)

	// magic
	err = pdfg.Create()
	if err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}
