using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using HarmonyLib;
using SafeExamBrowser.Monitoring.Keyboard;

namespace Milim.patch
{
    [HarmonyPatch(typeof(KeyboardInterceptor), "Start")]
    static class KeyboardInterceptorPatch
    {
        static bool Prefix()
        {
            return false;
        }
    }
}
