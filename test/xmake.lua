add_rules("mode.release", "mode.debug")

add_requires("gtk+3", "glib")

target("dbus")
    set_kind("binary")
    add_rules("vala")
    add_files("dbus.vala")
    add_packages("gtk+3", "glib")
    add_values("vala.packages", "gtk+-3.0")