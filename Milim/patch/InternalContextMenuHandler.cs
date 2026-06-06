using HarmonyLib;
using System;
using System.Linq;
using System.Reflection;

namespace Milim.patch
{
    public static class InternalContextMenuHandler
    {
        private static Harmony _harmony;

        public static void Apply(Harmony harmony)
        {
            _harmony = harmony;

            var clazz = AppDomain.CurrentDomain.GetAssemblies().FirstOrDefault(a => a.GetName().Name == "SafeExamBrowser.Browser");

            if (clazz != null)
            {
                Patch(clazz);
            }
            else
            {
                AppDomain.CurrentDomain.AssemblyLoad += onDllLoad;
            }
        }

        private static void onDllLoad(object sender, AssemblyLoadEventArgs ea)
        {
            if (ea.LoadedAssembly.GetName().Name == "SafeExamBrowser.Browser")
            {
                AppDomain.CurrentDomain.AssemblyLoad -= onDllLoad;
                Patch(ea.LoadedAssembly);
            }
        }

        private static void Patch(Assembly targetAssembly)
        {
            var type = targetAssembly.GetType("SafeExamBrowser.Browser.Handlers.ContextMenuHandler");
            var tMethod = type?.GetMethod("OnBeforeContextMenu", BindingFlags.Public | BindingFlags.Instance);

            if (tMethod == null) return;

            var prefix = typeof(InternalContextMenuHandler).GetMethod(nameof(Prefix), BindingFlags.Static | BindingFlags.NonPublic);

            _harmony.Patch(tMethod, prefix: new HarmonyMethod(prefix));
        }

        private static bool Prefix()
        {
            return false;
        }
    }
}