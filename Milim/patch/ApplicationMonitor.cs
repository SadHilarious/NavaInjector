using HarmonyLib;
using SafeExamBrowser.Monitoring.Applications;
using SafeExamBrowser.Settings.Applications;
using SafeExamBrowser.WindowsApi.Contracts;
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Milim.patch
{
    [HarmonyPatch(typeof(ApplicationMonitor), "SystemEvent_WindowChanged")]
    class SystemEvent_WindowChanged
    {
        static bool Prefix(ApplicationMonitor __instance, IntPtr handle)
        {
            var _progs = new List<string> {
                "snippingtool",
                "chrome",
                ""
            };

            var This = Traverse.Create(__instance);
            var nativeMethods = This.Field("nativeMethods").GetValue<INativeMethods>();
            string winTitle = nativeMethods.GetWindowTitle(handle);

            foreach (var prog in _progs)
            {
                if (winTitle.ToLower().Replace(" ", "").Contains(prog))
                {
                    return false;
                }
            }

            return true;
        }
    }

    [HarmonyPatch(typeof(ApplicationMonitor), "AddForTermination")]
    class Edotensei
    {
        static bool Prefix(ApplicationMonitor __instance, IProcess process)
        {
            var thiz = Traverse.Create(__instance);
            var blacklist = thiz.Field("blacklist").GetValue<IList<BlacklistApplication>>();
            var changeAbleName = process.Name?.Replace(" ", "").ToLower() ?? string.Empty;
            var originalName = process.OriginalName?.Replace(" ", "").ToLower() ?? string.Empty;

            foreach (var b in blacklist)
            {
                var bOri = b.OriginalName?.Replace(" ", "").ToLower() ?? string.Empty;
                var bExe = b.ExecutableName?.Replace(" ", "").ToLower() ?? string.Empty;

                if (!string.IsNullOrEmpty(bOri) && bOri == originalName)
                {
                    return false;
                }
                else if (!string.IsNullOrEmpty(bExe) && bExe == changeAbleName)
                {
                    return false;
                }
            }

 
            return true;
        }
    }

    [HarmonyPatch(typeof(ApplicationMonitor), "IsWhitelisted")]
    class IsWhitelistedPatch
    {
        static bool Prefix(ApplicationMonitor __instance, ref bool __result, IProcess process, out Guid? __state)
        {
            __state = Guid.NewGuid();
            __result = true;
            return false;

        }
    }

    [HarmonyPatch(typeof(ApplicationMonitor), "Timer_Elapsed")]
    public class CallPerSecond
    {
        static bool Prefix(ApplicationMonitor __instance, ref IList<IProcess> ___processes)
        {

            ___processes = ___processes
                .Where(proc => !proc.OriginalName.ToLower().Replace(" ", "").Contains("milim"))
                .Where(proc => !proc.OriginalName.ToLower().Replace(" ", "").Contains("chrome"))
                .Where(proc => !proc.OriginalName.ToLower().Replace(" ", "").Contains("snipping"))
                .Where(proc => !proc.OriginalName.ToLower().Replace(" ", "").Contains("dllhost"))
                .ToList();
            return true;
        }
    }
}
