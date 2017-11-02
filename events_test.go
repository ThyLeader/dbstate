package dbstate

import (
	"github.com/bwmarrin/discordgo"
	"testing"
)

func TestGuilds(t *testing.T) {
	g := &discordgo.Guild{
		ID:          "123",
		Name:        "some fun name",
		Large:       true,
		MemberCount: 10,
	}

	err := testState.GuildCreate(0, g)
	AssertFatal(t, err, "failed handling guild create")

	g2, err := testState.Guild(nil, "123")
	AssertFatal(t, err, "failed calling Guild(ID)")

	if g2.ID != g.ID || g2.Name != g.Name || g2.Large != g.Large || g2.MemberCount != g.MemberCount {
		t.Errorf("stored guild was modified: correct(%#v) stored(%#v)", g, g2)
	}

	// Test iterating
	cnt := 0
	err = testState.IterateGuilds(nil, func(g *discordgo.Guild) bool {
		cnt++
		return true
	})
	AssertFatal(t, err, "failed iterating guilds")

	if cnt != 1 {
		t.Errorf("unexpted guild count! Got: %d, Expected: 1", cnt)
	}

	err = testState.GuildDelete("123")
	AssertFatal(t, err, "failed removing guild")

	if _, err := testState.Guild(nil, "123"); err == nil {
		t.Fatal("guild still there after being deleted")
	}
}

func TestGuildMembers(t *testing.T) {
	m := &discordgo.Member{
		User:    &discordgo.User{ID: "123"},
		GuildID: "321",
		Nick:    "some fun name",
		Roles:   []string{"123", "321"},
	}

	err := testState.MemberUpdate(0, nil, m)
	AssertFatal(t, err, "failed handling member update")

	m2, err := testState.GuildMember(nil, "321", "123")
	AssertFatal(t, err, "failed calling GuildMember(gid, uid)")

	if m2.User.ID != m.User.ID || m2.Nick != m.Nick || len(m2.Roles) != len(m.Roles) || m2.GuildID != m.GuildID {
		t.Errorf("stored member was modified: correct(%#v) stored(%#v)", m, m2)
	}

	for i, v := range m.Roles {
		if m2.Roles[i] != v {
			t.Errorf("mismatched roles, Got: %q, Expected: %q", m2.Roles[i], v)
		}
	}

	cnt := 0
	err = testState.IterateGuildMembers(nil, "321", func(g *discordgo.Member) bool {
		cnt++
		return true
	})

	AssertFatal(t, err, "failed iterating members")

	if cnt != 1 {
		t.Errorf("unexpted member count! Got: %d, Expected: 1", cnt)
	}

	err = testState.MemberRemove(0, nil, "321", "123", false)
	AssertFatal(t, err, "failed removing member")

	if _, err = testState.GuildMember(nil, "321", "123"); err == nil {
		t.Fatal("member still there after being removed")
	}
}

func TestGuildChannels(t *testing.T) {
	c := &discordgo.Channel{
		GuildID: "1",
		ID:      "2",
		Type:    discordgo.ChannelTypeGuildText,
	}
	g := &discordgo.Guild{
		ID: "1",
	}

	AssertFatal(t, testState.GuildCreate(0, g), "failed creating guild")
	AssertFatal(t, testState.ChannelCreateUpdate(0, nil, c, true), "failed creating channel")

	fetched, err := testState.Channel(nil, "2")
	AssertFatal(t, err, "failed retrieving channel")

	if fetched.ID != c.ID || fetched.GuildID != c.GuildID || fetched.Type != c.Type {
		t.Errorf("mismatched results, got %#v, expected %#v", fetched, c)
	}

	AssertFatal(t, testState.ChannelDelete(0, nil, "2"), "failed deleting channel")
	if _, err = testState.Channel(nil, "2"); err == nil {
		t.Fatal("channel still there after being removed")
	}
	AssertErr(t, testState.GuildDelete("1"), "failed removing guild")
}
