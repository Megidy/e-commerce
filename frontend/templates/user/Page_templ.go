// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Page(nav bool) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Gear Up</title><script src=\"https://unpkg.com/htmx.org@1.9.2\"></script><style>\r\n                body {\r\n                    font-family: Arial, sans-serif;\r\n                    margin: 0;\r\n                    padding: 0;\r\n                    display: flex;\r\n                    flex-direction: column;\r\n                    min-height: 100vh;\r\n                    background: #fff;\r\n                }\r\n                .background-images {\r\n                    position: fixed;\r\n                    top: 0;\r\n                    left: 0;\r\n                    width: 100%;\r\n                    height: 100%;\r\n                    overflow: hidden;\r\n                    z-index: -1;\r\n                }\r\n                .background-images img {\r\n                    position: absolute;\r\n                    width: 50%;\r\n                    height: 50%;\r\n                    object-fit: cover;\r\n                    filter: brightness(80%); \r\n                }\r\n                .img1 {\r\n                    top: 0;\r\n                    left: 0;\r\n                }\r\n                .img2 {\r\n                    top: 0;\r\n                    left: 50%;\r\n                }\r\n                .img3 {\r\n                    top: 50%;\r\n                    left: 0;\r\n                }\r\n                .img4 {\r\n                    top: 50%;\r\n                    left: 50%;\r\n                }\r\n                nav {\r\n                    background-color: #DEA54B;\r\n                    padding: 15px;\r\n                    display: flex;\r\n                    justify-content: center;\r\n                    flex-wrap: wrap;\r\n                }\r\n                nav a {\r\n                    color: #F2CD5D;\r\n                    text-decoration: none;\r\n                    margin: 10px 15px;\r\n                    font-size: 18px;\r\n                    font-weight: bold;\r\n                }\r\n                nav a:hover {\r\n                    color: #3A1772;\r\n                }\r\n                @media (max-width: 600px) {\r\n                    nav {\r\n                        flex-direction: column;\r\n                        align-items: center;\r\n                    }\r\n                    nav a {\r\n                        margin: 10px 0;\r\n                        font-size: 16px;\r\n                        font-weight: bold;\r\n                    }\r\n                }\r\n            </style></head><body><div class=\"background-images\"><img src=\"/static/1.jpg\" class=\"img1\" alt=\"Background Image 1\"> <img src=\"/static/2.jpg\" class=\"img2\" alt=\"Background Image 2\"> <img src=\"/static/3.jpg\" class=\"img3\" alt=\"Background Image 3\"> <img src=\"/static/4.jpg\" class=\"img4\" alt=\"Background Image 4\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if nav {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<nav><a href=\"/products/accessories\">Accessories</a> <a href=\"/products/bicycles\">Bicycles</a> <a href=\"/orders\">Your Orders</a> <a href=\"/overall\">Overall about Bicycles</a> <a href=\"/tech-support\">Tech-Support</a> <a href=\"/user\">Account</a></nav>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
