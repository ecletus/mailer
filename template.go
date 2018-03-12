package mailer

import (
	"github.com/moisespsena/template/html/template"
	"github.com/qor/qor"
)

// Template email template
type Template struct {
	Name    string
	Layout  string
	Data    interface{}
	Context *qor.Context
	funcs   *template.FuncValues
}

// Funcs set template's funcs
func (tmpl Template) Funcs(funcMap... template.FuncMap) Template {
	if tmpl.funcs == nil {
		tmpl.funcs = &template.FuncValues{}
	}
	err := tmpl.funcs.Append(funcMap...)
	if err != nil {
		panic(err)
	}
	return tmpl
}

// FuncsValues set template's funcs
func (tmpl Template) FuncsValues(funcValues... *template.FuncValues) Template {
	tmpl.funcs.AppendValues(funcValues...)
	return tmpl
}

// Render render template
func (mailer Mailer) Render(t Template) Email {
	var email Email
	var tmpl = mailer.Config.Render.Template().FuncValues(t.funcs)

	if t.Layout != "" {
		if result, err := tmpl.Layout(t.Layout+".text").Render(t.Name+".text", t.Data, t.Context); err == nil {
			email.Text = string(result)
		}

		if result, err := tmpl.Layout(t.Layout+".html").Render(t.Name+".html", t.Data, t.Context); err == nil {
			email.HTML = string(result)
		} else if result, err := tmpl.Layout(t.Layout).Render(t.Name, t.Data, t.Context); err == nil {
			email.HTML = string(result)
		}
	} else {
		if result, err := tmpl.Render(t.Name+".text", t.Data, t.Context); err == nil {
			email.Text = string(result)
		}

		if result, err := tmpl.Render(t.Name+".html", t.Data, t.Context); err == nil {
			email.HTML = string(result)
		} else if result, err := tmpl.Render(t.Name, t.Data, t.Context); err == nil {
			email.HTML = string(result)
		}
	}

	return email
}
