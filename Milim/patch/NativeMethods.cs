using HarmonyLib;
using SafeExamBrowser.WindowsApi;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Milim.patch
{
    [HarmonyPatch(typeof(NativeMethods), nameof(NativeMethods.EmptyClipboard))]
    public class EmptyClipboard
    {
        static bool Prefix()
        {
            return false;
        }
    }

    [HarmonyPatch(typeof(NativeMethods), nameof(NativeMethods.HideWindow))]
    public class HideWindow
    {
        static bool Prefix()
        {
            return false;
        }
    }

    [HarmonyPatch(typeof(NativeMethods), nameof(NativeMethods.PostCloseMessageToShell))]
    public class PostCloseMessageToShell
    {
        static bool Prefix()
        {
            // Defense in depth: this posts the undocumented WM_CLOSE-like message 0x5B4 to
            // the shell window, which makes explorer.exe exit "gracefully". Only caller is
            // ExplorerShell.Terminate, but keep the weapon itself disabled as well.
            Log.Write("Blocked NativeMethods.PostCloseMessageToShell");
            return false;
        }
    }
}
