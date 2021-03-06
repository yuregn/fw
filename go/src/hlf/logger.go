package hlf

import (
	"fmt"
	"time"
)

//Logger log mechanism interface
type Logger interface {
	Child(string) Logger
	To(string) Logger
	Ntf(string, ...interface{})
	Inf(string, ...interface{})
	Err(string, ...interface{})
	Wrn(string, ...interface{})
	Dbg(string, ...interface{})
	Trc(string, ...interface{})
}

func createLogger(id string, target string, parent Logger) Logger {
	if target == "" {
		target = _conf.DefaultFile
	}

	lg := logger{
		id:     id,
		target: target,
	}
	var islogger bool
	lg.parent, islogger = parent.(*logger)
	if !islogger {
		lg.parent = nil
	}
	lg.loadConf()
	return &lg
}

//MakeLogger init a logger
func MakeLogger(id string) Logger {
	return createLogger(id, "", _defaultLogger)
}

type logger struct {
	id     string
	conf   loggerConf
	parent *logger
	target string
}

func (me *logger) loadConf() {
	id := me.id
	if id == "" {
		id = "_"
	}
	conf, found := _conf.Loggers[id]
	if found {
		me.conf = conf
	} else if me.parent != nil {
		me.conf = me.parent.conf
	} else {
		me.conf = _defaultLogConf
	}
}

func (me *logger) Child(id string) Logger {
	return createLogger(id, "", me)
}

func (me *logger) To(target string) Logger {
	return createLogger(me.id, target, me.parent)
}

func (me *logger) Ntf(format string, a ...interface{}) {
	me.print(_LV_NOTIFICATION, format, a...)
}

func (me *logger) Inf(format string, a ...interface{}) {
	me.print(_LV_INFO, format, a...)
}

func (me *logger) Err(format string, a ...interface{}) {
	me.print(_LV_ERROR, format, a...)
}

func (me *logger) Wrn(format string, a ...interface{}) {
	me.print(_LV_WARNING, format, a...)
}

func (me *logger) Dbg(format string, a ...interface{}) {
	me.print(_LV_DEBUG, format, a...)
}

func (me *logger) Trc(format string, a ...interface{}) {
	me.print(_LV_TRACE, format, a...)
}

func (me *logger) print(lv logLevel, format string, a ...interface{}) {
	if me.conf.Lv < lv {
		return
	}

	text := me.formati(nil, lv, format, a...)
	if me.conf.ToConsole {
		me.send2console(text)
	}

	if me.conf.ToFile {
		text = me.formati(me, lv, format, a...)
		me.send2file(me.target, text)

		for log := me.parent; log != nil; log = log.parent {
			clv := me.findAppliedLv(log)
			if clv >= lv {
				text := me.formati(log, lv, format, a...)
				log.send2file(me.target, text)
			}
		}
	}
}

func (me *logger) findAppliedLv(ancestor *logger) logLevel {
	var lv logLevel = _LV_UNKNOWN
	found := false

	for log := me; log != ancestor && log.parent != nil; log = log.parent {
		lv, found = ancestor.conf.ChildLv[log.id]
		if found {
			break
		}
	}

	if !found {
		lv = ancestor.conf.DefaultChildLv
	}

	return lv
}

func (me *logger) send2console(text string) {
	me.send2Srv("console:", text)
}

func (me *logger) send2file(target string, text string) {
	me.send2Srv(me.getFileTarget(target), text)
}

func (me *logger) send2Srv(target string, text string) {
	li := logItem{
		target: target,
		text:   text,
	}
	_logSrvCh <- li
}

func (me *logger) toPrefix() string {
	if me.id == "" {
		return " "
	}
	return "[" + me.id + "] "
}

func (me *logger) formati(parent *logger, lv logLevel, format string, a ...interface{}) string {
	text := ""
	text += lv.toPrefix()
	text += logTime()
	text += me.indent(parent)
	text += me.toPrefix()
	text += fmt.Sprintf(format, a...)
	text += "\n"
	return text
}

func (me *logger) getIndent(parent *logger) int {
	indent := 0
	for log := me; log != parent && log.parent != nil; log = log.parent {
		indent += _conf.Indent
	}
	return indent
}

func (me *logger) indent(parent *logger) string {
	indent := me.getIndent(parent)
	indents := ""
	for i := 0; i < indent; i++ {
		indents += " "
	}
	return indents
}

func (me *logger) getPath() string {
	path := ""
	for l := me; l != nil; l = l.parent {
		if l.id != "" {
			path = l.id + "/" + path
		}
	}
	path = _logRoot + path
	return path
}

func (me *logger) getFileTarget(f string) string {
	return "file:" + me.getPath() + f
}

func logTime() string {
	return "[" + time.Now().Format(time.RFC3339) + "]"
}
