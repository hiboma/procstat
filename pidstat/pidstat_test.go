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

	Describe(t, "pidstat.Get(1)", func() {
		pidstat := NewPidstat("../misc/proc/")
		stat, _ := pidstat.GetAll()
		It("return map[string] of map[string]interfae{}", func() {
			Expect(stat["1"]["comm"]).To(Equal, "(init)")
			Expect(stat["1"]["ppid"]).To(Equal, 0)
			Expect(stat["1"]["state"]).To(Equal, "S")
		})
	})
}
