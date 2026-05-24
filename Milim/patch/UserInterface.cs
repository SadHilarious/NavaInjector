using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using HarmonyLib;
using SafeExamBrowser.UserInterface.Shared.Utilities;

namespace Milim.patch
{
    [HarmonyPatch(typeof(WindowExtensions), nameof(WindowExtensions.ExcludeFromCapture))]
    public class ExcludeFormCapture
    {
        static bool Prefix()
        {
            return false;
        }
    }
}
