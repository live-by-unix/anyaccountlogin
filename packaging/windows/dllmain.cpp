/*
 * DLL entry point for AnyAccountLogin Credential Provider
 */

#include "AnyAccountLoginCredentialProvider.h"
#include <windows.h>

static LONG g_cRef = 1;

HRESULT DllCanUnloadNow()
{
    return g_cRef == 1 ? S_OK : S_FALSE;
}

HRESULT DllGetClassObject(REFCLSID rclsid, REFIID riid, LPVOID* ppv)
{
    if (rclsid != CLSID_AnyAccountLoginCredentialProvider)
    {
        return CLASS_E_CLASSNOTAVAILABLE;
    }

    return AnyAccountLoginCredentialProviderFactory::CreateInstance(riid, ppv);
}

HRESULT DllRegisterServer()
{
    HRESULT hr = S_OK;

    // Register the credential provider
    HKEY hKey = NULL;
    hr = RegCreateKeyExW(HKEY_LOCAL_MACHINE, 
                         L"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\Credential Providers\\{CLSID}",
                         0, NULL, 0, KEY_WRITE, NULL, &hKey, NULL);
    
    if (SUCCEEDED(hr))
    {
        const wchar_t* dllPath = L"C:\\Program Files\\AnyAccountLogin\\AnyAccountLoginCredentialProvider.dll";
        RegSetValueExW(hKey, NULL, 0, REG_SZ, (const BYTE*)dllPath, (wcslen(dllPath) + 1) * sizeof(wchar_t));
        RegCloseKey(hKey);
    }

    return hr;
}

HRESULT DllUnregisterServer()
{
    // Unregister the credential provider
    RegDeleteKeyW(HKEY_LOCAL_MACHINE, 
                  L"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Authentication\\Credential Providers\\{CLSID}");
    return S_OK;
}

STDAPI DllRegisterServer()
{
    return DllRegisterServer();
}

STDAPI DllUnregisterServer()
{
    return DllUnregisterServer();
}

void DllAddRef()
{
    InterlockedIncrement(&g_cRef);
}

void DllRelease()
{
    InterlockedDecrement(&g_cRef);
}

BOOL APIENTRY DllMain(HMODULE hModule, DWORD ul_reason_for_call, LPVOID lpReserved)
{
    switch (ul_reason_for_call)
    {
    case DLL_PROCESS_ATTACH:
        DisableThreadLibraryCalls(hModule);
        break;
    case DLL_PROCESS_DETACH:
        break;
    }
    return TRUE;
}
