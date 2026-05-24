# Nava
Nava work against SEB 3.10.1

The goal of this script is to prevent SEB to start new obfuscated desktop and put in KIOSK mode. My method is to patch all necessary function with Harmony C#. You can see all that function inside `NavaInjector/Milim/patch`

## Installation


```bash
git clone https://github.com/seynth/NavaInjector
cd NavaInjector
```

## Step 1 - Build the runner (Nava.dll)

Nava.dll is a payload that have an exported function called `Nava`. After injector inject this Nava.dll to both
SafeExamBrowser.exe and SafeExamBrowser.Client.exe, exported function can then be called from injector by Opening SafeExamBrowser.exe and SafeExamBrowser.Client.exe process. The purpose of this dll to manipulate .NET CLR to call some function i've patched using Harmony C# of Safe Exam Browser. That's why this script will not change browser-exam-key and config-key-hash


```bash
cd Nava
go build -o nava.dll --buildmode=c-shared
```

## Step 2 - Build the payload (Milim.dll)

Milim.dll is a Harmony C# project that have all necessary modified function of Safe Exam Browser, when Milim.dll injected to parent and client process it will do the work like preventing to make obfuscated desktop, whitelisted all blacklisted software, unintercepted the keyboard, and more.

- Open solution file NavaInjector/Milim/Milim.sln in Visual Studio 2022
- Make sure to use `Debug` and `x64` platform configuration
- Right click on Milim project and select build or press `Ctrl + Shift + B`

or

```bash
cd Milim
dotnet build
```

## Step 3 - Build the injector

The injector is entry point of this project. I embed Nava.dll and Milim.dll to this executable to make it standalone executable 

```bash
cd injector
copy ..\Milim\bin\Debug\Milim.dll .\
copy ..\Nava\nava.dll .\
go build -o nava_standalone.exe
```

## Usage

Usage is very simple, you just need to run NavaInjector.exe as Administrator, but make sure to turn off Real-time Protection windows defender.

Wait until it say `[Nava] Starting`, then you can start your Safe Exam Browser by double click .seb file or from exercise/quiz provider. Dont close your terminal!.




