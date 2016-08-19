package mayday

import (
	"bytes"
	"github.com/coreos/go-systemd/dbus"
	godbus "github.com/godbus/dbus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJournalHeader(t *testing.T) {
	content := bytes.NewBufferString("testd daemon log")
	jnl := SystemdJournal{name: "testd", content: content}
	assert.Equal(t, jnl.Name(), "testd")

	hdr := jnl.Header() // header is generated from name and content
	assert.Equal(t, hdr.Name, "/journals/testd.log")
	assert.Equal(t, hdr.Size, int64(content.Len()))
}

func TestListJournals(t *testing.T) {
	getJournals = func() ([]dbusStatus, error) {

		statuses := []dbusStatus{
			{
				unit: dbus.UnitStatus{Name: "testd"},
				property: &dbus.Property{"testd",
					godbus.MakeVariant("/usr/lib64/systemd/system/testd.service")}},
			{
				unit: dbus.UnitStatus{Name: "examd"},
				property: &dbus.Property{"examd",
					godbus.MakeVariant("/usr/lib64/systemd/system/examd.service")}},
			{
				unit: dbus.UnitStatus{Name: "notaservice"},
				property: &dbus.Property{"notaservice",
					godbus.MakeVariant("/usr/lib/systemd/system/umount.target")},
			}}

		return statuses, nil
	}

	journals, err := ListJournals()
	assert.Nil(t, err)
	assert.Len(t, journals, 2)
	assert.Equal(t, journals[0].Name(), "testd")
	assert.Equal(t, journals[1].Name(), "examd")
}
