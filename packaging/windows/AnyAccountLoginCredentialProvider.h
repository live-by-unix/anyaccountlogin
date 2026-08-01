/*
 * AnyAccountLogin Credential Provider for Windows
 * 
 * This credential provider integrates with Windows Credential Provider
 * framework to provide flash drive-based authentication at the login screen.
 */

#pragma once

#include <windows.h>
#include <credentialprovider.h>
#include <string>

class AnyAccountLoginCredentialProvider : public ICredentialProvider
{
public:
    AnyAccountLoginCredentialProvider();
    virtual ~AnyAccountLoginCredentialProvider();

    // IUnknown methods
    IFACEMETHODIMP QueryInterface(REFIID riid, void** ppv);
    IFACEMETHODIMP_(ULONG) AddRef();
    IFACEMETHODIMP_(ULONG) Release();

    // ICredentialProvider methods
    IFACEMETHODIMP SetUsageScenario(CREDENTIAL_PROVIDER_USAGE_SCENARIO cpus, DWORD dwFlags);
    IFACEMETHODIMP SetSerialization(const CREDENTIAL_PROVIDER_CREDENTIAL_SERIALIZATION* pcpcs);
    IFACEMETHODIMP GetFieldDescriptor(DWORD dwFieldID, CREDENTIAL_PROVIDER_FIELD_DESCRIPTOR** ppcpfd);
    IFACEMETHODIMP GetCredentialCount(DWORD* pdwCount, DWORD* pdwDefault, DWORD* pdwCountAtLogon);
    IFACEMETHODIMP GetCredentialAt(DWORD dwIndex, ICredentialProviderCredential** ppcpc);
    IFACEMETHODIMP Advise(ICredentialProviderEvents* pcpe, UINT_PTR upAdviseContext);
    IFACEMETHODIMP UnAdvise();
    IFACEMETHODIMP GetFieldState(DWORD dwFieldID, CREDENTIAL_PROVIDER_FIELD_STATE* pcpfs, 
                                  CREDENTIAL_PROVIDER_FIELD_INTERACTIVE_STATE* pcpfis);

protected:
    LONG _cRef;
    CREDENTIAL_PROVIDER_USAGE_SCENARIO _cpus;
    ICredentialProviderEvents* _pcpe;
    UINT_PTR _upAdviseContext;
    bool _fReenumerateRequested;
};

class AnyAccountLoginCredential : public ICredentialProviderCredential
{
public:
    AnyAccountLoginCredential();
    virtual ~AnyAccountLoginCredential();

    // IUnknown methods
    IFACEMETHODIMP QueryInterface(REFIID riid, void** ppv);
    IFACEMETHODIMP_(ULONG) AddRef();
    IFACEMETHODIMP_(ULONG) Release();

    // ICredentialProviderCredential methods
    IFACEMETHODIMP Advise(ICredentialProviderCredentialEvents* pcpce);
    IFACEMETHODIMP UnAdvise();
    IFACEMETHODIMP SetSelected(BOOL* pbAutoLogon);
    IFACEMETHODIMP SetDeselected();
    IFACEMETHODIMP GetFieldState(DWORD dwFieldID, CREDENTIAL_PROVIDER_FIELD_STATE* pcpfs,
                                  CREDENTIAL_PROVIDER_FIELD_INTERACTIVE_STATE* pcpfis);
    IFACEMETHODIMP GetStringValue(DWORD dwFieldID, PWSTR* ppwsz);
    IFACEMETHODIMP GetBitmapValue(DWORD dwFieldID, HBITMAP* phbmp);
    IFACEMETHODIMP GetSubmitButtonValue(DWORD dwFieldID, DWORD* pdwAdjacentTo);
    IFACEMETHODIMP GetCheckboxValue(DWORD dwFieldID, BOOL* pbChecked, PWSTR* ppwszLabel);
    IFACEMETHODIMP GetComboBoxValueCount(DWORD dwFieldID, DWORD* pcItems, DWORD* pdwSelectedItem);
    IFACEMETHODIMP GetComboBoxValueAt(DWORD dwFieldID, DWORD dwItem, PWSTR* ppwszItem);
    IFACEMETHODIMP SetStringValue(DWORD dwFieldID, PCWSTR pwsz);
    IFACEMETHODIMP GetSerialization(CREDENTIAL_PROVIDER_GET_SERIALIZATION_RESPONSE* pcpgsr,
                                    CREDENTIAL_PROVIDER_CREDENTIAL_SERIALIZATION* pcpcs,
                                    PWSTR* ppwszOptionalStatusText,
                                    CREDENTIAL_PROVIDER_STATUS_ICON* pcpsiOptionalStatusIcon);
    IFACEMETHODIMP Result(HRESULT hrStatus, const CREDENTIAL_PROVIDER_FIELD_STATE* pcpfs);
    IFACEMETHODIMP GetOptions(PCREDENTIAL_PROVIDER_OPTIONS pcpo);

private:
    LONG _cRef;
    ICredentialProviderCredentialEvents* _pcpce;
    std::wstring _pszFlashDrivePath;
    std::wstring _pszPassword;
    bool _fChecked;

    bool ValidateFlashDriveAuthentication(const std::wstring& flashDrivePath, const std::wstring& password);
};

class AnyAccountLoginCredentialProviderFactory : public IClassFactory
{
public:
    static HRESULT CreateInstance(REFIID riid, void** ppv);

    // IUnknown methods
    IFACEMETHODIMP QueryInterface(REFIID riid, void** ppv);
    IFACEMETHODIMP_(ULONG) AddRef();
    IFACEMETHODIMP_(ULONG) Release();

    // IClassFactory methods
    IFACEMETHODIMP CreateInstance(IUnknown* pUnkOuter, REFIID riid, void** ppv);
    IFACEMETHODIMP LockServer(BOOL fLock);

private:
    AnyAccountLoginCredentialProviderFactory();
    virtual ~AnyAccountLoginCredentialProviderFactory();
    LONG _cRef;
};
