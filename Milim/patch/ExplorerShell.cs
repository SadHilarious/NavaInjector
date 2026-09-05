using HarmonyLib;
using SafeExamBrowser.WindowsApi;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Milim.patch
{
    [HarmonyPatch(typeof(ExplorerShell), nameof(ExplorerShell.HideAllWindows))]
    public class UnhideWindow
    {
        static bool Prefix()
        {
            return false;
        }
    }

    [HarmonyPatch(typeof(ExplorerShell), "KillExplorerShell")]
    public class EdotenseiExplorerShell
    {
        static bool Prefix()
        {
            return false;
        }
    }

    [HarmonyPatch(typeof(ExplorerShell), nameof(ExplorerShell.Terminate))]
    public class KeepExplorerAlive
    {
        static bool Prefix()
        {
            // Terminate() is the PRIMARY explorer kill: it posts the undocumented close
            // message 0x5B4 to the shell (PostCloseMessageToShell) and only falls back to
            // KillExplorerShell (taskkill) after a 3s grace period. Callers:
            // KioskModeOperation.TerminateExplorerShell (session start, config killExplorerShell=true)
            // and MonitoringResponsibility.ApplicationMonitor_ExplorerStarted (client re-kill).
            Log.Write("Blocked ExplorerShell.Terminate");
            return false;
        }
    }

    [HarmonyPatch(typeof(ExplorerShell), nameof(ExplorerShell.Start))]
    public class NoDuplicateExplorer
    {
        static bool Prefix()
        {
            // Since Terminate() never runs, explorer never dies and KioskModeOperation's
            // reversion (RestartExplorerShell) must not spawn a second shell process.
            // Also avoids Start()'s unbounded wait loop when no shell window ever appears.
            Log.Write("Blocked ExplorerShell.Start");
            return false;
        }
    }
}
