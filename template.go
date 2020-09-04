package mailer

import (
	"github.com/moisespsena-go/os-common"
	"github.com/moisespsena/template/html/template"
	"github.com/ecletus/core"
)

// Template email template
type Template struct {
	Name    string
	Layout  string
	Data    interface{}
	Context *core.Context
	funcs   template.FuncValues
}

// Funcs set template's funcs
func (tmpl Template) Funcs(funcMap... template.FuncMap) Template {
	err := tmpl.funcs.Append(funcMap...)
	if err != nil {
		panic(err)
	}
	return tmpl
}

// FuncsValues set template's funcs
func (tmpl Template) FuncsValues(funcValues... template.FuncValues) Template {
	tmpl.funcs.AppendValues(funcValues...)
	return tmpl
}

// Render render template
func (mailer Mailer) Render(t Template, langs ...string) (email Email, err error) {
	var tmpl = mailer.Config.Render.Template().SetFuncValues(t.funcs)

	if t.Layout != "" {
		var result template.HTML
		if result, err = tmpl.SetLayout(t.Layout+".text").Render(nil, t.Name+".text", t.Data, t.Context, langs...); err == nil {
			email.Text = string(result)
		}
		if result, err = tmpl.SetLayout(t.Layout+".html").Render(nil, t.Name+".html", t.Data, t.Context, langs...); err == nil {
			email.HTML = string(result)
		} else if result, err = tmpl.SetLayout(t.Layout).Render(nil, t.Name, t.Data, t.Context, langs...); err == nil {
			email.HTML = string(result)
		}
	} else {
		var result template.HTML
		if result, err = tmpl.Render(nil, t.Name+".text", t.Data, t.Context, langs...); err == nil {
			email.Text = string(result)
		}

		if result, err = tmpl.Render(nil, t.Name+".html", t.Data, t.Context, langs...); err == nil {
			email.HTML = string(result)
		} else if result, err = tmpl.Render(nil, t.Name, t.Data, t.Context, langs...); err == nil {
			email.HTML = string(result)
		}
	}

	if oscommon.IsNotFound(err) && (email.Text != "" || email.HTML != "") {
		err = nil
	}

	return
}
