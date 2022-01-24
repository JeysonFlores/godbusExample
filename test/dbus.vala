
[DBus (name = "com.github.jeysonflores.DBusTest.Palette")]
interface DBusTest.Regex : Object {
    public abstract string ping () throws GLib.Error;
    public abstract string ping2 () throws GLib.Error;
    public abstract string ping3 (string name) throws GLib.Error;

    public signal void pinged (string result);
}


public class DBusTest.Application : Gtk.Application {
    DBusTest.Regex palette_interface = null;
    public Application () {
        Object (
            application_id: "com.github.jeysonflores.xpressr",
            flags: ApplicationFlags.FLAGS_NONE
        );
    }

    protected override void activate () {
        try {
            palette_interface = Bus.get_proxy_sync (BusType.SESSION, "com.github.jeysonflores.DBusTest",
                                                                "/com/github/jeysonflores/DBusTest/Palette");
            
                                                                
            
            print (palette_interface.ping ());
            print ("\n");

            print (palette_interface.ping2 ());
            print ("\n");

            print (palette_interface.ping3 ("DBus from Vala"));
            print ("\n");

            
        } catch (GLib.Error e) {
            warning ("⚠️ Error trying to connect to the DBus service: [%s]\n", e.message);
        }
        palette_interface.pinged.connect((result) => {
            print (result);
            print ("\n");
        });
        var main_window = new Gtk.ApplicationWindow (this);
        main_window.show_all ();
    }

    public static int main (string[] args) {
        return new Application ().run (args);
    }
}