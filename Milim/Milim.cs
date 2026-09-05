using HarmonyLib;
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Text;
using System.Threading.Tasks;

namespace Milim
{
    internal static class Log
    {
        // Minimal diagnostic logger: one line per event, appended to %TEMP%\milim.log.
        // Used to reconstruct the patch timeline of a run (e.g. whether PatchAll landed
        // before SEB's kiosk mode attempted to terminate the explorer shell).
        internal static void Write(string message)
        {
            try
            {
                var path = System.IO.Path.Combine(System.IO.Path.GetTempPath(), "milim.log");
                var line = System.DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss.fff")
                    + " [" + System.Diagnostics.Process.GetCurrentProcess().ProcessName + "] " + message;
                System.IO.File.AppendAllText(path, line + "\r\n");
            }
            catch
            {
                // Never let logging break the host process.
            }
        }
    }

    public class InitMilim
    {
        public static int Nava(string unused)
        {
            var milim = new Harmony("lord.of.wrath");
            milim.PatchAll(Assembly.GetExecutingAssembly());
            Log.Write("PatchAll complete (harmony id 'lord.of.wrath')");

            if (Process.GetCurrentProcess().ProcessName.ToLower() == "safeexambrowser.client") {
                patch.InternalContextMenuHandler.Apply(milim);
            }

            return 7;
        }
    }
}
