// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
// templates/active_task.templ

package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/ze674/EZLine/internal/models"
	"strconv"
)

func ActiveTask(task models.Task, isScanning bool, packer string) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"bg-white shadow-md rounded-lg p-6\"><div class=\"flex justify-between items-center mb-6\"><h2 class=\"text-2xl font-bold\">Выбранное задание #")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(task.ID))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/active_task.templ`, Line: 12, Col: 100}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</h2><div class=\"flex space-x-2\"><a href=\"/tasks\" class=\"bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded\">К списку заданий</a><form method=\"post\" action=\"/tasks/finish\"><button type=\"submit\" class=\"bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded\">Завершить</button></form></div></div><div class=\"bg-blue-50 rounded-lg p-6 mb-6\"><h3 class=\"text-xl font-semibold mb-4\">Информация о задании</h3><div class=\"grid grid-cols-1 md:grid-cols-2 gap-4\"><div><p class=\"font-semibold\">Продукт:</p><p>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(task.ProductName)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/active_task.templ`, Line: 30, Col: 40}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</p></div><div><p class=\"font-semibold\">Дата:</p><p>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(task.Date)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/active_task.templ`, Line: 34, Col: 33}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</p></div><div><p class=\"font-semibold\">Номер партии:</p><p>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(task.BatchNumber)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/active_task.templ`, Line: 38, Col: 40}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "</p></div><div><p class=\"font-semibold\">Статус:</p><div class=\"mt-1\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if task.Status == "новое" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "<span class=\"bg-blue-100 text-blue-800 py-1 px-2 rounded-full\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 string
			templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(task.Status)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/active_task.templ`, Line: 44, Col: 103}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if task.Status == "в работе" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "<span class=\"bg-yellow-100 text-yellow-800 py-1 px-2 rounded-full\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var7 string
			templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(task.Status)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/active_task.templ`, Line: 46, Col: 107}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if task.Status == "завершено" {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, "<span class=\"bg-green-100 text-green-800 py-1 px-2 rounded-full\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var8 string
			templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(task.Status)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/active_task.templ`, Line: 48, Col: 105}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, "</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "<span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var9 string
			templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(task.Status)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/active_task.templ`, Line: 50, Col: 46}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 13, "</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 14, "</div></div></div></div><!-- Добавляем блок для управления сканированием --><div class=\"bg-gray-100 p-6 rounded-lg mt-4\"><h3 class=\"text-xl font-semibold mb-4\">Управление сканированием</h3><div class=\"flex items-center justify-between\"><div><p class=\"text-gray-600\">Статус сканирования:</p>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if isScanning {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 15, "<span class=\"bg-green-100 text-green-800 py-1 px-2 rounded-full\">Активно</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 16, "<span class=\"bg-red-100 text-red-800 py-1 px-2 rounded-full\">Остановлено</span>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 17, "</div><div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if isScanning {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 18, "<form method=\"post\" action=\"/scanning/stop\"><button type=\"submit\" class=\"bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded\">Остановить сканирование</button></form>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 19, "<form method=\"post\" action=\"/scanning/start\"><button type=\"submit\" class=\"bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded\">Начать сканирование</button></form>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 20, "</div></div></div></div><div class=\"bg-gray-100 p-6 rounded-lg mt-4\"><h3 class=\"text-xl font-semibold mb-4\">Управление упаковщиком</h3><div><p class=\"text-gray-600 mb-2\">Текущий упаковщик: <span class=\"font-semibold\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var10 string
		templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(packer)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/active_task.templ`, Line: 90, Col: 112}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 21, "</span></p><form method=\"post\" action=\"/packer/change\" class=\"flex items-center space-x-2\"><input type=\"text\" name=\"packer\" placeholder=\"Новый упаковщик\" class=\"border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500\" required> <button type=\"submit\" class=\"bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded\">Применить</button></form></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
