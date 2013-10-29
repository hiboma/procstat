package pidstat

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Stat is like 'map[pid:1 ppid:0 comm:(zsh-4.3.17) ... ]]'
// int and string is combined with
type Stat map[string]interface{}
type Pidstat struct {
	basedir string
}

// field keys is taken from `$ man 5 proc`
var field_keys []string = []string{
	/* 0  */ "pid", "comm", "state", "ppid", "pgrp",
	/* 5  */ "session", "tty_nr", "tpgid", "flags", "minflt",
	/* 10 */ "cminflt", "majflt", "cmajflt", "utime", "stime",
	/* 15 */ "cutime", "cstime", "priority", "nice", "num_threads",
	/* 20 */ "itrealvalue", "starttime", "vsize", "rss", "rsslim",
	/* 25 */ "startcode", "endcode", "startstack", "kstkesp", "kstkeip",
	/* 30 */ "signal", "blocked", "sigignore", "sigcatch", "wchan",
	/* 35 */ "nswap", "cnswap", "exit_signal", "processor", "rt_priority",
	/* 40 */ "policy", "delayacct_blkio_ticks", "guest_time", "cguest_time",
}

func NewPidstat(basedir string) *Pidstat {
	if basedir == "" {
		basedir = "/proc"
	}
	return &Pidstat{basedir: basedir}
}

func (self *Pidstat) Get(pid string) (stat Stat, err error) {

	file, err := os.Open(fmt.Sprintf("%s/%s/stat", self.basedir, pid))
	if err != nil {
		return
	}

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return
	}

	chunks := strings.Split(line, " ")
	stat = make(map[string]interface{})
	for i, value := range chunks {
		key := field_keys[i]
		if key == "comm" || key == "state" {
			stat[key] = value
		} else {
			value_int, err := strconv.Atoi(value)
			if err != nil {
				break
			}
			stat[key] = value_int
		}
	}
	return
}

func (self *Pidstat) GetAll() (stats map[string]Stat, err error) {

	r, err := regexp.Compile(`^\d+$`)
	if err != nil {
		return
	}

	f, err := os.Open(self.basedir)
	if err != nil {
		return
	}

	entries, err := f.Readdirnames(-1)
	if err != nil {
		return
	}

	stats = make(map[string]Stat)
	for _, pid := range entries {
		if r.MatchString(pid) {
			stat, err := self.Get(pid)
			if err != nil && err != io.EOF {
				continue
			}
			stats[pid] = stat
		}
	}

	return
}

func (self *Pidstat) Grep(stats map[string]Stat, args ...interface{}) (filterd map[string]Stat) {

	var filter func(st Stat) bool

	/* any better way ?? */
	switch len(args) {
	case 1:
		filter = args[0].(func(st Stat) bool)
	case 2:
		filter = self.compileFilter(args[0].(string), args[1].(string))
	default:
		panic("too many arguments for Grep()")
	}

	filterd = make(map[string]Stat)
	for pid, stat := range stats {
		if filter(stat) {
			filterd[pid] = stat
		}
	}

	fmt.Printf("%v", filterd)

	return
}

func (self *Pidstat) compileFilter(field string, condition string) (compiled func(st Stat) bool) {

	chunks := strings.Split(condition, ":")
	if len(chunks) != 2 {
		panic(fmt.Errorf("specified condition format is invalid: %s\n", condition))
	}

	/* noeq, eq, gt, lt */
	operator := chunks[0]
	value := chunks[1]
	value_int, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}

	switch operator {
	case "eq":
		compiled = func(st Stat) bool {
			if field == "comm" || field == "state" {
				return st[field] == value
			} else {
				return st[field] == value_int
			}
		}
	case "gt":
		compiled = func(st Stat) bool {
			if field == "comm" || field == "state" {
				return false
			} else {
				return st[field].(int) > value_int
			}
		}
	case "lt":
		compiled = func(st Stat) bool {
			if field == "comm" || field == "state" {
				return false
			} else {
				return st[field].(int) < value_int
			}
		}
	case "ne":
		compiled = func(st Stat) bool {
			if field == "comm" || field == "state" {
				return st[field] != value
			} else {
				return st[field].(int) != value_int
			}
		}
	default:
		err = fmt.Errorf("unknown operator: %s", chunks[0])
	}

	return
}
