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
    public class InitMilim
    {
        public static int Nava(string unused)
        {
            var milim = new Harmony("lord.of.wrath");
            milim.PatchAll(Assembly.GetExecutingAssembly());

            if (Process.GetCurrentProcess().ProcessName.ToLower() == "safeexambrowser.client") {
                patch.InternalContextMenuHandler.Apply(milim);
            }

            return 7;
        }
    }
}
