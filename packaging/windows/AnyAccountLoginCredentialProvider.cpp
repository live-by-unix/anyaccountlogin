/*
 * AnyAccountLogin Credential Provider Implementation
 */

#include "AnyAccountLoginCredentialProvider.h"
#include <shlwapi.h>
#include <strsafe.h>
#include <fstream>

// Field IDs
#define FID_FLASH_DRIVE_PATH 1
#define FID_PASSWORD 2
#define FID_SUBMIT_BUTTON 3
#define FID_LARGE_TEXT 4

AnyAccountLoginCredentialProvider::AnyAccountLoginCredentialProvider()
    : _cRef(1), _cpus(CPUS_INVALID), _pcpe(NULL), _upAdviseContext(0), _fReenumerateRequested(false)
{
    DllAddRef();
}

AnyAccountLoginCredentialProvider::~AnyAccountLoginCredentialProvider()
{
    if (_pcpe != NULL)
    {
        _pcpe->Release();
    }
    DllRelease();
}

IFACEMETHODIMP AnyAccountLoginCredentialProvider::QueryInterface(REFIID riid, void** ppv)
{
    static const QITAB qit[] = {
        QITABENT(AnyAccountLoginCredentialProvider, ICredentialProvider),
        {0},
    };
    return QISearch(this, qit, riid, ppv);
}

IFACEMETHODIMP_(ULONG) AnyAccountLoginCredentialProvider::AddRef()
{
    return InterlockedIncrement(&_cRef);
}

IFACEMETHODIMP_(ULONG) AnyAccountLoginCredentialProvider::Release()
{
    LONG cRef = InterlockedDecrement(&_cRef);
    if (!cRef)
    {
        delete this;
    }
    return cRef;
}

IFACEMETHODIMP AnyAccountLoginCredentialProvider::SetUsageScenario(CREDENTIAL_PROVIDER_USAGE_SCENARIO cpus, DWORD dwFlags)
{
    _cpus = cpus;
    return S_OK;
}

IFACEMETHODIMP AnyAccountLoginCredentialProvider::SetSerialization(const CREDENTIAL_PROVIDER_CREDENTIAL_SERIALIZATION* pcpcs)
{
    return E_NOTIMPL;
}

IFACEMETHODIMP AnyAccountLoginCredentialProvider::GetFieldDescriptor(DWORD dwFieldID, CREDENTIAL_PROVIDER_FIELD_DESCRIPTOR** ppcpfd)
{
    return E_NOTIMPL;
}

IFACEMETHODIMP AnyAccountLoginCredentialProvider::GetCredentialCount(DWORD* pdwCount, DWORD* pdwDefault, DWORD* pdwCountAtLogon)
{
    *pdwCount = 1;
    *pdwDefault = 0;
    *pdwCountAtLogon = 1;
    return S_OK;
}

IFACEMETHODIMP AnyAccountLoginCredentialProvider::GetCredentialAt(DWORD dwIndex, ICredentialProviderCredential** ppcpc)
{
    if (dwIndex != 0)
    {
        return E_INVALIDARG;
    }

    AnyAccountLoginCredential* pCredential = new AnyAccountLoginCredential();
    if (pCredential == NULL)
    {
        return E_OUTOFMEMORY;
    }

    *ppcpc = pCredential;
    return S_OK;
}

IFACEMETHODIMP AnyAccountLoginCredentialProvider::Advise(ICredentialProviderEvents* pcpe, UINT_PTR upAdviseContext)
{
    _pcpe = pcpe;
    _pcpe->AddRef();
    _upAdviseContext = upAdviseContext;
    return S_OK;
}

IFACEMETHODIMP AnyAccountLoginCredentialProvider::UnAdvise()
{
    if (_pcpe != NULL)
    {
        _pcpe->Release();
        _pcpe = NULL;
    }
    _upAdviseContext = 0;
    return S_OK;
}

IFACEMETHODIMP AnyAccountLoginCredentialProvider::GetFieldState(DWORD dwFieldID, CREDENTIAL_PROVIDER_FIELD_STATE* pcpfs,
    CREDENTIAL_PROVIDER_FIELD_INTERACTIVE_STATE* pcpfis)
{
    return E_NOTIMPL;
}

// AnyAccountLoginCredential implementation
AnyAccountLoginCredential::AnyAccountLoginCredential()
    : _cRef(1), _pcpce(NULL), _fChecked(false)
{
    DllAddRef();
}

AnyAccountLoginCredential::~AnyAccountLoginCredential()
{
    if (_pcpce != NULL)
    {
        _pcpce->Release();
    }
    DllRelease();
}

IFACEMETHODIMP AnyAccountLoginCredential::QueryInterface(REFIID riid, void** ppv)
{
    static const QITAB qit[] = {
        QITABENT(AnyAccountLoginCredential, ICredentialProviderCredential),
        {0},
    };
    return QISearch(this, qit, riid, ppv);
}

IFACEMETHODIMP_(ULONG) AnyAccountLoginCredential::AddRef()
{
    return InterlockedIncrement(&_cRef);
}

IFACEMETHODIMP_(ULONG) AnyAccountLoginCredential::Release()
{
    LONG cRef = InterlockedDecrement(&_cRef);
    if (!cRef)
    {
        delete this;
    }
    return cRef;
}

IFACEMETHODIMP AnyAccountLoginCredential::Advise(ICredentialProviderCredentialEvents* pcpce)
{
    _pcpce = pcpce;
    _pcpce->AddRef();
    return S_OK;
}

IFACEMETHODIMP AnyAccountLoginCredential::UnAdvise()
{
    if (_pcpce != NULL)
    {
        _pcpce->Release();
        _pcpce = NULL;
    }
    return S_OK;
}

IFACEMETHODIMP AnyAccountLoginCredential::SetSelected(BOOL* pbAutoLogon)
{
    *pbAutoLogon = FALSE;
    return S_OK;
}

IFACEMETHODIMP AnyAccountLoginCredential::SetDeselected()
{
    return S_OK;
}

IFACEMETHODIMP AnyAccountLoginCredential::GetFieldState(DWORD dwFieldID, CREDENTIAL_PROVIDER_FIELD_STATE* pcpfs,
    CREDENTIAL_PROVIDER_FIELD_INTERACTIVE_STATE* pcpfis)
{
    *pcpfs = CPFS_DISPLAY_IN_SELECTED_TILE;
    *pcpfis = CPFIS_NONE;
    return S_OK;
}

IFACEMETHODIMP AnyAccountLoginCredential::GetStringValue(DWORD dwFieldID, PWSTR* ppwsz)
{
    HRESULT hr = S_OK;
    switch (dwFieldID)
    {
    case FID_FLASH_DRIVE_PATH:
        hr = SHStrDupW(_pszFlashDrivePath.c_str(), ppwsz);
        break;
    case FID_PASSWORD:
        hr = SHStrDupW(_pszPassword.c_str(), ppwsz);
        break;
    default:
        hr = E_INVALIDARG;
        break;
    }
    return hr;
}

IFACEMETHODIMP AnyAccountLoginCredential::GetBitmapValue(DWORD dwFieldID, HBITMAP* phbmp)
{
    return E_NOTIMPL;
}

IFACEMETHODIMP AnyAccountLoginCredential::GetSubmitButtonValue(DWORD dwFieldID, DWORD* pdwAdjacentTo)
{
    if (dwFieldID == FID_SUBMIT_BUTTON)
    {
        *pdwAdjacentTo = FID_PASSWORD;
        return S_OK;
    }
    return E_INVALIDARG;
}

IFACEMETHODIMP AnyAccountLoginCredential::GetCheckboxValue(DWORD dwFieldID, BOOL* pbChecked, PWSTR* ppwszLabel)
{
    return E_NOTIMPL;
}

IFACEMETHODIMP AnyAccountLoginCredential::GetComboBoxValueCount(DWORD dwFieldID, DWORD* pcItems, DWORD* pdwSelectedItem)
{
    return E_NOTIMPL;
}

IFACEMETHODIMP AnyAccountLoginCredential::GetComboBoxValueAt(DWORD dwFieldID, DWORD dwItem, PWSTR* ppwszItem)
{
    return E_NOTIMPL;
}

IFACEMETHODIMP AnyAccountLoginCredential::SetStringValue(DWORD dwFieldID, PCWSTR pwsz)
{
    HRESULT hr = S_OK;
    switch (dwFieldID)
    {
    case FID_FLASH_DRIVE_PATH:
        _pszFlashDrivePath = pwsz;
        break;
    case FID_PASSWORD:
        _pszPassword = pwsz;
        break;
    default:
        hr = E_INVALIDARG;
        break;
    }
    return hr;
}

IFACEMETHODIMP AnyAccountLoginCredential::GetSerialization(CREDENTIAL_PROVIDER_GET_SERIALIZATION_RESPONSE* pcpgsr,
    CREDENTIAL_PROVIDER_CREDENTIAL_SERIALIZATION* pcpcs,
    PWSTR* ppwszOptionalStatusText,
    CREDENTIAL_PROVIDER_STATUS_ICON* pcpsiOptionalStatusIcon)
{
    // Validate flash drive authentication
    if (!ValidateFlashDriveAuthentication(_pszFlashDrivePath, _pszPassword))
    {
        *pcpgsr = CPGSR_RETURN_NO_CREDENTIAL_FINISHED;
        *ppwszOptionalStatusText = L"Invalid flash drive or password";
        *pcpsiOptionalStatusIcon = CPSI_ERROR;
        return S_OK;
    }

    // Return serialization for successful authentication
    *pcpgsr = CPGSR_RETURN_CREDENTIAL_FINISHED;
    return S_OK;
}

IFACEMETHODIMP AnyAccountLoginCredential::Result(HRESULT hrStatus, const CREDENTIAL_PROVIDER_FIELD_STATE* pcpfs)
{
    return S_OK;
}

IFACEMETHODIMP AnyAccountLoginCredential::GetOptions(PCREDENTIAL_PROVIDER_OPTIONS pcpo)
{
    return E_NOTIMPL;
}

bool AnyAccountLoginCredential::ValidateFlashDriveAuthentication(const std::wstring& flashDrivePath, const std::wstring& password)
{
    // Check if flash drive path exists
    DWORD attrs = GetFileAttributesW(flashDrivePath.c_str());
    if (attrs == INVALID_FILE_ATTRIBUTES || !(attrs & FILE_ATTRIBUTE_DIRECTORY))
    {
        return false;
    }

    // Check for PasswordAuthCode.txt
    std::wstring authCodePath = flashDrivePath + L"\\PasswordAuthCode.txt";
    std::ifstream authFile(authCodePath);
    if (!authFile.is_open())
    {
        return false;
    }

    std::string storedAuthCode((std::istreambuf_iterator<char>(authFile)),
                                std::istreambuf_iterator<char>());
    authFile.close();

    // Compare passwords
    std::wstring wideStoredAuthCode(storedAuthCode.begin(), storedAuthCode.end());
    if (password != wideStoredAuthCode)
    {
        return false;
    }

    // Check for PEM key
    std::wstring pemPath = flashDrivePath + L"\\PasswordAuth.pem";
    attrs = GetFileAttributesW(pemPath.c_str());
    if (attrs == INVALID_FILE_ATTRIBUTES)
    {
        return false;
    }

    return true;
}

// Factory implementation
AnyAccountLoginCredentialProviderFactory::AnyAccountLoginCredentialProviderFactory()
    : _cRef(1)
{
    DllAddRef();
}

AnyAccountLoginCredentialProviderFactory::~AnyAccountLoginCredentialProviderFactory()
{
    DllRelease();
}

HRESULT AnyAccountLoginCredentialProviderFactory::CreateInstance(REFIID riid, void** ppv)
{
    AnyAccountLoginCredentialProvider* pProvider = new AnyAccountLoginCredentialProvider();
    if (pProvider == NULL)
    {
        return E_OUTOFMEMORY;
    }

    HRESULT hr = pProvider->QueryInterface(riid, ppv);
    pProvider->Release();
    return hr;
}

IFACEMETHODIMP AnyAccountLoginCredentialProviderFactory::QueryInterface(REFIID riid, void** ppv)
{
    static const QITAB qit[] = {
        QITABENT(AnyAccountLoginCredentialProviderFactory, IClassFactory),
        {0},
    };
    return QISearch(this, qit, riid, ppv);
}

IFACEMETHODIMP_(ULONG) AnyAccountLoginCredentialProviderFactory::AddRef()
{
    return InterlockedIncrement(&_cRef);
}

IFACEMETHODIMP_(ULONG) AnyAccountLoginCredentialProviderFactory::Release()
{
    LONG cRef = InterlockedDecrement(&_cRef);
    if (!cRef)
    {
        delete this;
    }
    return cRef;
}

IFACEMETHODIMP AnyAccountLoginCredentialProviderFactory::CreateInstance(IUnknown* pUnkOuter, REFIID riid, void** ppv)
{
    if (pUnkOuter != NULL)
    {
        return CLASS_E_NOAGGREGATION;
    }
    return CreateInstance(riid, ppv);
}

IFACEMETHODIMP AnyAccountLoginCredentialProviderFactory::LockServer(BOOL fLock)
{
    if (fLock)
    {
        DllAddRef();
    }
    else
    {
        DllRelease();
    }
    return S_OK;
}
