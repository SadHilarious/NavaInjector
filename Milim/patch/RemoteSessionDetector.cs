using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using HarmonyLib;
using SafeExamBrowser.Monitoring;

namespace Milim.patch
{
    [HarmonyPatch(typeof(RemoteSessionDetector), nameof(RemoteSessionDetector.IsRemoteSession))]
    public class RemoteSessionDetectorPatch
    {
        public void Postfix(ref bool __result)
        {
            __result = false;
        }
    }
}
