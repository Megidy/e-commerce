// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/Megidy/e-commerce/frontend/templates/user"

func LoadAddAccessoryPage(message string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"swap\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templates.Page(true).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<html><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Add Accessory</title><script src=\"https://unpkg.com/htmx.org@1.9.2\"></script></head><body><form hx-post=\"/products/accessory/action/add/confirm\" hx-swap=\"outerHTML\" hx-target=\".swap\"><div class=\"form-group\"><label for=\"id\">id</label> <input type=\"text\" id=\"id\" name=\"id\" required></div><div class=\"form-group\"><label for=\"name\">name</label> <input type=\"text\" id=\"name\" name=\"name\" required></div><div class=\"form-group\"><label for=\"description\">description</label> <input type=\"text\" id=\"description\" name=\"description\" required></div><div class=\"form-group\"><label for=\"quantity\">quantity</label> <input type=\"text\" id=\"quantity\" name=\"quantity\" required></div><div class=\"form-group\"><label for=\"price\">price</label> <input type=\"text\" id=\"price\" name=\"price\" required></div><div class=\"form-group\"><label for=\"category\">category</label> <input type=\"text\" id=\"category\" name=\"category\" required></div><div class=\"form-group\"><label for=\"name\">image</label> <input type=\"text\" id=\"image\" name=\"image\" required></div><button type=\"submit\" class=\"btn\">Add Accessory</button><div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(message)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `frontend/templates/product/AddAccessoryPage.templ`, Line: 46, Col: 33}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></form></body></html></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate