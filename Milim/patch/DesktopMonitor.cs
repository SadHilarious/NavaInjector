using HarmonyLib;
using SafeExamBrowser.WindowsApi.Contracts;
using SafeExamBrowser.WindowsApi.Desktops;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Milim.patch
{
    [HarmonyPatch(typeof(DesktopMonitor), "Timer_Elapsed")]
    public class UnMonitorDesktop
    {
        static bool Prefix()
        {
            return false;
        }
    }

    [HarmonyPatch(typeof(DesktopMonitor), nameof(DesktopMonitor.Start))]
    public class Start_Hook
    {
        static bool Prefix()
        {
            return false;
        }
    }
}
