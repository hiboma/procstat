package pidstat

import (
	. "github.com/r7kamura/gospel"
	"testing"
)

func TestPidstat(t *testing.T) {
	Describe(t, "pidstat.Get(1)", func() {
		pidstat := NewPidstat("../misc/proc/")
		stat, _ := pidstat.Get("1")
		It("return map[string]interfae{}", func() {
			Expect(stat["comm"]).To(Equal, "(init)")
			Expect(stat["ppid"]).To(Equal, 0)
			Expect(stat["state"]).To(Equal, "S")
		})
	})

	Describe(t, "pidstat.GetAll()", func() {
		pidstat := NewPidstat("../misc/proc/")
		stat, _ := pidstat.GetAll()
		It("return map[string] of map[string]interfae{}", func() {
			Expect(stat["1"]["comm"]).To(Equal, "(init)")
			Expect(stat["1"]["ppid"]).To(Equal, 0)
			Expect(stat["1"]["state"]).To(Equal, "S")
		})
	})

	Describe(t, "pidstat.Grep()", func() {
		pidstat := NewPidstat("../misc/proc/")
		stats, _ := pidstat.GetAll()

		It("pid == 1 should return map", func() {
			f := func(st Stat) bool {
				return st["pid"] == 1
			}
			filterd := pidstat.Grep(stats, f)
			Expect(filterd["1"]["pid"]).To(Equal, 1)
		})

		It("pid != 1 should not return map", func() {
			f := func(st Stat) bool {
				return st["pid"] != 1
			}
			filterd := pidstat.Grep(stats, f)
			Expect(len(filterd)).To(Equal, 0)
		})

		It("'pid eq:1' should return map", func() {
			filterd := pidstat.Grep(stats, "pid", "eq:1")
			Expect(filterd["1"]["pid"]).To(Equal, 1)
		})

		It("'pid ne:1' should not return map", func() {
			filterd := pidstat.Grep(stats, "pid", "ne:1")
			Expect(len(filterd)).To(Equal, 0)
		})

		It("'rss lt:300' should return map", func() {
			filterd := pidstat.Grep(stats, "rss", "lt:300")
			Expect(filterd["1"]["pid"]).To(Equal, 1)
		})

		It("'rss gt:300' should not return map", func() {
			filterd := pidstat.Grep(stats, "rss", "gt:300")
			Expect(len(filterd)).To(Equal, 0)
		})
	})
}
