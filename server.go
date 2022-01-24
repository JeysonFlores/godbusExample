package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/jeysonflores/dbustest/pkg/datamanager"
	_ "github.com/mattn/go-sqlite3"
)

type Palette struct {
	bus   *dbus.Conn
	model *datamanager.Palette
}

func (p *Palette) AnotherMethod(name string, id int, juj int64, array []string, another_id int) (map[string]int, *dbus.Error) {
	return nil, nil
}

func (p *Palette) Ping() (string, *dbus.Error) {
	p.EmitPingedSignal("Pong")
	return "Pong", nil
}

func (p *Palette) Insert(name string, desc string) (string, *dbus.Error) {
	err := p.model.Insert(name, desc)
	if (err) != nil {
		return "", dbus.MakeFailedError(errors.New("something went wrong"))
	}
	return "palette inserted", nil
}

func (p *Palette) Ping3(id int) (string, *dbus.Error) {
	result, err := p.model.GetById(id)
	if (err) != nil {
		return "", dbus.MakeFailedError(errors.New(result))
	}

	return result, nil
}

func (p *Palette) EmitPingedSignal(name string) {
	fmt.Println("Pinged Signal called")
	p.bus.Emit(p.GetObjectPath(), p.GetInterfacePath()+".Pinged", string(name))
}

func (p *Palette) GetObjectPath() dbus.ObjectPath {
	return dbus.ObjectPath("/com/github/jeysonflores/DBusTest/Palette")
}

func (p *Palette) GetInterfacePath() string {
	return "com.github.jeysonflores.DBusTest.Palette"
}

func (p *Palette) GetIntroData() string {
	return `<node>
				<interface name="com.github.jeysonflores.DBusTest.Palette">
					<method name="Ping">
						<arg name="result" direction="out" type="s"/>
					</method>
					<method name="Insert">
						<arg name="name" direction="in" type="s"/>
						<arg name="desc" direction="in" type="s"/>
						<arg name="result" direction="out" type="s"/>
					</method>
					<method name="Ping3">
						<annotation name="org.freedesktop.DBus.Method.Async" value="server" />
						<arg name="id" direction="in" type="i"/>
						<arg name="result" direction="out" type="s"/>
					</method>
					<signal name="Pinged">
						<arg name="param" type="s"/>
					</signal>
				</interface>
				<interface name="org.freedesktop.DBus.Introspectable">
					<method name="Introspect">
							<arg name="out" direction="out" type="s"/>
					</method>
				</interface>
			</node>
		`
}

func (p *Palette) RegisterToBus() error {
	p.bus.Export(p, p.GetObjectPath(), p.GetInterfacePath())
	p.bus.Export(introspect.Introspectable(p.GetIntroData()), p.GetObjectPath(), "org.freedesktop.DBus.Introspectable")
	return nil
}

func NewPalette(bus *dbus.Conn, db *sql.DB) (*Palette, error) {
	palModel := &datamanager.Palette{
		Con: db,
	}
	palModel.CreateTable()

	return &Palette{
		bus,
		palModel,
	}, nil
}

/*type Palette struct {
	bus  *dbus.Conn
	conn *sql.DB
}

//Methods
func (p *Palette) Ping() (string, *dbus.Error) {
	return "Pong", nil
}

func (p *Palette) Insert() (string, *dbus.Error) {
	return "Pong2", nil
}

func (p *Palette) Ping3(id int) (string, *dbus.Error) {
	palette := &datamanager.Palette{
		Con: p.conn,
	}
	result, err := palette.GetById(id)
	if (err) != nil {
		return "", dbus.MakeFailedError(errors.New("something went wrong"))
	}

	return result, nil
}

//Signals
func (p *Palette) EmitPingedSignal(param string) {
	p.bus.Emit(p.GetObjectPath(), p.GetInterfacePath()+".Pinged", param)
}

//Utils
func (p *Palette) GetObjectPath() dbus.ObjectPath {
	return dbus.ObjectPath("com/github/jeysonflores/DBusTest/Palette")
}

func (p *Palette) GetInterfacePath() string {
	return "com.github.jeysonflores.DBusTest.Palette"
}

func (p *Palette) GetIntroData() string {
	return `<node>
				<interface name="com.github.jeysonflores.DBusTest.Palette">
					<method name="Ping">
						<arg name="result" direction="out" type="s"/>
					</method>
					<method name="Insert">
						<arg name="result" direction="out" type="s"/>
					</method>
					<method name="Ping3">
						<arg name="id" direction="in" type="i"/>
						<arg name="result" direction="out" type="s"/>
					</method>
					<signal name="Pinged">
						<arg name="param" direction="out" type="s"/>
					</signal>
				</interface>
				<interface name="org.freedesktop.DBus.Introspectable">
					<method name="Introspect">
							<arg name="out" direction="out" type="s"/>
					</method>
				</interface>
			</node>`
}

func (p *Palette) RegisterToBus() error {
	p.bus.Export(p, p.GetObjectPath(), p.GetInterfacePath())
	fmt.Println("object exported")
	p.bus.Export(introspect.Introspectable(p.GetIntroData()), p.GetObjectPath(), "org.freedesktop.DBus.Introspectable")
	return nil
}*/

func main() {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		panic(err)
	}

	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}

	iface, _ := NewPalette(conn, db)

	iface.RegisterToBus()

	reply, err := conn.RequestName("com.github.jeysonflores.DBusTest",
		dbus.NameFlagDoNotQueue)

	if err != nil {
		panic(err)
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}

	fmt.Println("Listening on com.github.jeysonflores.DBusTest ...")
	select {}
}
